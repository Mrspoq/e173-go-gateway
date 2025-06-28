package sip

import (
    "log"
    
    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/e173-gateway/e173_go_gateway/pkg/validation"
    "github.com/e173-gateway/e173_go_gateway/pkg/repository"
)

// NewBasicSIPServerWithDB creates a SIP server with database-backed WhatsApp validation
func NewBasicSIPServerWithDB(port int, whatsappAPIKey string, dbPool *pgxpool.Pool) *BasicSIPServer {
    // Create WhatsApp cache repository
    cacheRepo := repository.NewSimpleWhatsAppValidationRepository(dbPool)
    
    // Create filter engine with database-backed validation
    filterEng := &FilterEngine{
        blacklistSvc:     &BlacklistService{},
        whatsappAPI:      validation.NewPrivateWhatsAppValidatorDB(whatsappAPIKey, cacheRepo),
        phoneValidator:   validation.NewGooglePhoneValidator("NG"), // Default to Nigeria
        historyAnalyzer:  &CallHistoryAnalyzer{},
        operatorDetector: &OperatorPrefixDetector{},
    }
    
    return &BasicSIPServer{
        port:       port,
        filterEng:  filterEng,
        routingEng: NewRoutingEngine(),
        voiceAI:    NewVoiceAIService(),
        logger:     log.New(log.Writer(), "[SIP] ", log.LstdFlags),
    }
}