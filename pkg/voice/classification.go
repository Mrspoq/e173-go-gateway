package voice

import (
    "context"
    "fmt"
    "strings"
    "regexp"
)

// LLMClassifier uses language models to classify call transcripts
type LLMClassifier struct {
    llmProvider        LLMProvider
    confidenceThreshold float64
    spamKeywords       []string
    operatorPatterns   []string
}

// LLMProvider interface for language model integration
type LLMProvider interface {
    Analyze(ctx context.Context, prompt string) (string, error)
    GetModel() string
}

// NewLLMClassifier creates a new LLM-based classifier
func NewLLMClassifier(llmProvider LLMProvider) *LLMClassifier {
    return &LLMClassifier{
        llmProvider:         llmProvider,
        confidenceThreshold: 0.7,
        spamKeywords: []string{
            "congratulations", "win", "prize", "urgent", "act now",
            "limited time", "offer", "discount", "warranty", "insurance",
            "loan", "credit", "debt", "IRS", "tax", "legal action",
            "press 1", "press 2", "automated message", "this is a recording",
        },
        operatorPatterns: []string{
            "sim card blocked", "sim blocked", "line suspended",
            "insufficient credit", "low balance", "recharge required",
            "your number has been", "service suspended", "payment required",
            "voicemail", "leave a message", "mailbox",
        },
    }
}

// ClassifyCall analyzes a transcript to determine call type
func (c *LLMClassifier) ClassifyCall(ctx context.Context, transcript *Transcript) (*Classification, error) {
    // First, do quick keyword-based classification
    quickClass := c.quickClassify(transcript.Text)
    if quickClass != nil && quickClass.Confidence > 0.9 {
        return quickClass, nil
    }
    
    // Build LLM prompt
    prompt := c.buildClassificationPrompt(transcript.Text)
    
    // Get LLM analysis
    response, err := c.llmProvider.Analyze(ctx, prompt)
    if err != nil {
        // Fallback to keyword classification
        if quickClass != nil {
            return quickClass, nil
        }
        return nil, fmt.Errorf("LLM analysis failed: %w", err)
    }
    
    // Parse LLM response
    classification, err := c.parseLLMResponse(response)
    if err != nil {
        // Again, fallback to keyword classification
        if quickClass != nil {
            return quickClass, nil
        }
        return nil, fmt.Errorf("failed to parse LLM response: %w", err)
    }
    
    // Extract keywords from transcript
    classification.Keywords = c.extractKeywords(transcript.Text)
    
    return classification, nil
}

// quickClassify performs fast keyword-based classification
func (c *LLMClassifier) quickClassify(text string) *Classification {
    lowerText := strings.ToLower(text)
    
    // Check for operator messages
    for _, pattern := range c.operatorPatterns {
        if strings.Contains(lowerText, pattern) {
            category := CategoryOperatorIVR
            action := ActionFlagSIM
            reason := "Operator message detected"
            
            if strings.Contains(pattern, "blocked") || strings.Contains(pattern, "suspended") {
                category = CategorySIMBlocked
                reason = "SIM card blocked by operator"
            } else if strings.Contains(pattern, "credit") || strings.Contains(pattern, "balance") {
                category = CategoryLowCredit
                reason = "Low credit detected"
            } else if strings.Contains(pattern, "voicemail") {
                category = CategoryVoicemail
                action = ActionNormalRouting
                reason = "Voicemail system detected"
            }
            
            return &Classification{
                Category:   category,
                Confidence: 0.95,
                Action:     action,
                Reason:     reason,
                RiskScore:  0.8,
            }
        }
    }
    
    // Check for spam keywords
    spamScore := 0
    for _, keyword := range c.spamKeywords {
        if strings.Contains(lowerText, keyword) {
            spamScore++
        }
    }
    
    if spamScore >= 3 {
        return &Classification{
            Category:   CategorySpamRobocall,
            Confidence: float64(spamScore) / 10.0,
            Action:     ActionRouteToAI,
            Reason:     "Multiple spam keywords detected",
            RiskScore:  0.9,
        }
    }
    
    return nil
}

// buildClassificationPrompt creates the LLM prompt
func (c *LLMClassifier) buildClassificationPrompt(transcript string) string {
    return fmt.Sprintf(`Analyze this call transcript and classify it into one of these categories:

Transcript: "%s"

Categories:
- SPAM_ROBOCALL: Automated marketing, scam, or unwanted robocall
- SIM_BLOCKED: Operator message about SIM being blocked or suspended
- VOICEMAIL: Voicemail system or answering machine
- NORMAL_CALL: Regular human conversation
- OPERATOR_IVR: Operator system message (not blocking)
- LOW_CREDIT: Message about insufficient credit or balance

Respond in this exact format:
category: [CATEGORY]
confidence: [0.0-1.0]
action: [route_to_ai|flag_sim|normal_routing|block_call]
reason: [Brief explanation]
risk_score: [0.0-1.0]

Consider these factors:
1. Is it an automated message or human speech?
2. Does it mention account status, credit, or service issues?
3. Is it trying to sell something or requesting action?
4. Does it sound like a legitimate service message?`, transcript)
}

// parseLLMResponse extracts classification from LLM response
func (c *LLMClassifier) parseLLMResponse(response string) (*Classification, error) {
    classification := &Classification{}
    
    // Parse category
    if match := regexp.MustCompile(`category:\s*(\w+)`).FindStringSubmatch(response); len(match) > 1 {
        classification.Category = CallCategory(match[1])
    } else {
        return nil, fmt.Errorf("could not parse category")
    }
    
    // Parse confidence
    if match := regexp.MustCompile(`confidence:\s*([\d.]+)`).FindStringSubmatch(response); len(match) > 1 {
        fmt.Sscanf(match[1], "%f", &classification.Confidence)
    }
    
    // Parse action
    if match := regexp.MustCompile(`action:\s*(\w+)`).FindStringSubmatch(response); len(match) > 1 {
        classification.Action = CallAction(strings.ToUpper(match[1]))
    }
    
    // Parse reason
    if match := regexp.MustCompile(`reason:\s*(.+)`).FindStringSubmatch(response); len(match) > 1 {
        classification.Reason = strings.TrimSpace(match[1])
    }
    
    // Parse risk score
    if match := regexp.MustCompile(`risk_score:\s*([\d.]+)`).FindStringSubmatch(response); len(match) > 1 {
        fmt.Sscanf(match[1], "%f", &classification.RiskScore)
    }
    
    return classification, nil
}

// extractKeywords pulls out significant words from transcript
func (c *LLMClassifier) extractKeywords(text string) []string {
    keywords := []string{}
    lowerText := strings.ToLower(text)
    
    // Check all defined keywords
    allKeywords := append(c.spamKeywords, c.operatorPatterns...)
    for _, keyword := range allKeywords {
        if strings.Contains(lowerText, keyword) {
            keywords = append(keywords, keyword)
        }
    }
    
    return keywords
}

// GetConfidenceThreshold returns the minimum confidence for classification
func (c *LLMClassifier) GetConfidenceThreshold() float64 {
    return c.confidenceThreshold
}

// RuleBasedClassifier provides a simple rule-based fallback classifier
type RuleBasedClassifier struct {
    rules               []ClassificationRule
    confidenceThreshold float64
}

// ClassificationRule defines a pattern-based rule
type ClassificationRule struct {
    Pattern    *regexp.Regexp
    Category   CallCategory
    Action     CallAction
    Confidence float64
    Reason     string
}

// NewRuleBasedClassifier creates a simple rule-based classifier
func NewRuleBasedClassifier() *RuleBasedClassifier {
    return &RuleBasedClassifier{
        confidenceThreshold: 0.6,
        rules: []ClassificationRule{
            {
                Pattern:    regexp.MustCompile(`(?i)(sim.*blocked|line.*suspended|service.*terminated)`),
                Category:   CategorySIMBlocked,
                Action:     ActionFlagSIM,
                Confidence: 0.95,
                Reason:     "SIM blocking message detected",
            },
            {
                Pattern:    regexp.MustCompile(`(?i)(insufficient.*credit|low.*balance|recharge.*required)`),
                Category:   CategoryLowCredit,
                Action:     ActionFlagSIM,
                Confidence: 0.9,
                Reason:     "Low credit message detected",
            },
            {
                Pattern:    regexp.MustCompile(`(?i)(press.*\d|automated.*message|this.*is.*recording)`),
                Category:   CategorySpamRobocall,
                Action:     ActionRouteToAI,
                Confidence: 0.85,
                Reason:     "Automated call pattern detected",
            },
            {
                Pattern:    regexp.MustCompile(`(?i)(voicemail|leave.*message|mailbox)`),
                Category:   CategoryVoicemail,
                Action:     ActionNormalRouting,
                Confidence: 0.9,
                Reason:     "Voicemail system detected",
            },
        },
    }
}

// ClassifyCall using rules
func (r *RuleBasedClassifier) ClassifyCall(ctx context.Context, transcript *Transcript) (*Classification, error) {
    for _, rule := range r.rules {
        if rule.Pattern.MatchString(transcript.Text) {
            return &Classification{
                Category:   rule.Category,
                Confidence: rule.Confidence,
                Action:     rule.Action,
                Reason:     rule.Reason,
                RiskScore:  0.5, // Default risk score for rules
            }, nil
        }
    }
    
    // Default to normal call if no rules match
    return &Classification{
        Category:   CategoryNormalCall,
        Confidence: 0.5,
        Action:     ActionNormalRouting,
        Reason:     "No specific patterns detected",
        RiskScore:  0.1,
    }, nil
}

// GetConfidenceThreshold returns the threshold
func (r *RuleBasedClassifier) GetConfidenceThreshold() float64 {
    return r.confidenceThreshold
}