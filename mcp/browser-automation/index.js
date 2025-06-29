const express = require('express');
const puppeteer = require('puppeteer');
const bodyParser = require('body-parser');
const sharp = require('sharp');
const fs = require('fs').promises;
const path = require('path');

const app = express();
app.use(bodyParser.json({ limit: '50mb' }));

let browser;
let page;

// Initialize browser
async function initBrowser() {
    browser = await puppeteer.launch({
        headless: 'new',
        args: ['--no-sandbox', '--disable-setuid-sandbox']
    });
    page = await browser.newPage();
    await page.setViewport({ width: 1280, height: 800 });
    console.log('Browser initialized');
}

// Navigate to URL
app.post('/navigate', async (req, res) => {
    try {
        const { url } = req.body;
        await page.goto(url, { waitUntil: 'networkidle2' });
        const title = await page.title();
        res.json({ success: true, title });
    } catch (error) {
        res.status(500).json({ success: false, error: error.message });
    }
});

// Take screenshot
app.post('/screenshot', async (req, res) => {
    try {
        const { filename = 'screenshot.png', fullPage = false } = req.body;
        const screenshot = await page.screenshot({ 
            fullPage,
            encoding: 'base64'
        });
        
        // Save screenshot
        const screenshotPath = path.join('/tmp', filename);
        await fs.writeFile(screenshotPath, Buffer.from(screenshot, 'base64'));
        
        // Also create a smaller preview
        const preview = await sharp(Buffer.from(screenshot, 'base64'))
            .resize(400)
            .toBuffer();
        
        res.json({ 
            success: true, 
            path: screenshotPath,
            preview: preview.toString('base64')
        });
    } catch (error) {
        res.status(500).json({ success: false, error: error.message });
    }
});

// Click element
app.post('/click', async (req, res) => {
    try {
        const { selector } = req.body;
        await page.click(selector);
        res.json({ success: true });
    } catch (error) {
        res.status(500).json({ success: false, error: error.message });
    }
});

// Type text
app.post('/type', async (req, res) => {
    try {
        const { selector, text } = req.body;
        await page.type(selector, text);
        res.json({ success: true });
    } catch (error) {
        res.status(500).json({ success: false, error: error.message });
    }
});

// Get element text
app.post('/getText', async (req, res) => {
    try {
        const { selector } = req.body;
        const text = await page.$eval(selector, el => el.textContent);
        res.json({ success: true, text });
    } catch (error) {
        res.status(500).json({ success: false, error: error.message });
    }
});

// Check if element exists
app.post('/exists', async (req, res) => {
    try {
        const { selector } = req.body;
        const exists = await page.$(selector) !== null;
        res.json({ success: true, exists });
    } catch (error) {
        res.status(500).json({ success: false, error: error.message });
    }
});

// Wait for selector
app.post('/waitFor', async (req, res) => {
    try {
        const { selector, timeout = 5000 } = req.body;
        await page.waitForSelector(selector, { timeout });
        res.json({ success: true });
    } catch (error) {
        res.status(500).json({ success: false, error: error.message });
    }
});

// Get page content
app.post('/getContent', async (req, res) => {
    try {
        const content = await page.content();
        res.json({ success: true, content });
    } catch (error) {
        res.status(500).json({ success: false, error: error.message });
    }
});

// Evaluate JavaScript
app.post('/evaluate', async (req, res) => {
    try {
        const { script } = req.body;
        const result = await page.evaluate(script);
        res.json({ success: true, result });
    } catch (error) {
        res.status(500).json({ success: false, error: error.message });
    }
});

// Get visual analysis
app.post('/analyzeUI', async (req, res) => {
    try {
        const analysis = await page.evaluate(() => {
            const issues = [];
            
            // Check for broken images
            document.querySelectorAll('img').forEach(img => {
                if (!img.complete || img.naturalHeight === 0) {
                    issues.push({
                        type: 'broken-image',
                        selector: img.src,
                        message: 'Broken or missing image'
                    });
                }
            });
            
            // Check for empty containers
            document.querySelectorAll('[id$="-list"], [id$="-table"], [id$="-content"]').forEach(el => {
                if (el.children.length === 0 || el.textContent.trim() === '') {
                    issues.push({
                        type: 'empty-container',
                        selector: '#' + el.id,
                        message: 'Empty container that should have content'
                    });
                }
            });
            
            // Check for loading states (excluding HTMX indicators)
            document.querySelectorAll('.animate-pulse, [class*="loading"]').forEach(el => {
                // Skip HTMX indicators which are meant to pulse
                if (el.closest('.htmx-indicator')) return;
                // Skip elements that are just loading placeholders
                if (el.textContent.includes('Loading')) return;
                
                issues.push({
                    type: 'loading-state',
                    selector: el.className,
                    message: 'Element stuck in loading state'
                });
            });
            
            // Check for error messages
            document.querySelectorAll('[class*="error"], [class*="alert-danger"], .text-red-500').forEach(el => {
                if (el.textContent.trim()) {
                    issues.push({
                        type: 'error-message',
                        text: el.textContent.trim(),
                        message: 'Error message displayed'
                    });
                }
            });
            
            return {
                title: document.title,
                url: window.location.href,
                issues: issues,
                stats: {
                    images: document.querySelectorAll('img').length,
                    links: document.querySelectorAll('a').length,
                    forms: document.querySelectorAll('form').length,
                    buttons: document.querySelectorAll('button').length
                }
            };
        });
        
        res.json({ success: true, analysis });
    } catch (error) {
        res.status(500).json({ success: false, error: error.message });
    }
});

// Start server
const PORT = process.env.PORT || 3001;
app.listen(PORT, async () => {
    console.log(`Browser automation MCP running on port ${PORT}`);
    await initBrowser();
});

// Cleanup on exit
process.on('SIGINT', async () => {
    if (browser) await browser.close();
    process.exit();
});
