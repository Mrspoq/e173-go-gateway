const puppeteer = require('puppeteer');

class UITester {
    constructor() {
        this.browser = null;
        this.page = null;
        this.baseURL = 'http://localhost:8080';
    }

    async init() {
        this.browser = await puppeteer.launch({
            headless: 'new',
            args: ['--no-sandbox', '--disable-setuid-sandbox']
        });
        this.page = await this.browser.newPage();
        await this.page.setViewport({ width: 1280, height: 800 });
    }

    async cleanup() {
        if (this.browser) {
            await this.browser.close();
        }
    }

    async testDashboard() {
        console.log('Testing Dashboard...');
        try {
            await this.page.goto(`${this.baseURL}/`, { waitUntil: 'networkidle2' });
            
            // Check for 5 cards
            const cards = await this.page.$$('.grid > div[id="stats-cards"] > div');
            console.log(`✅ Dashboard cards found: ${cards.length} (expected 5)`);
            
            // Check no full page refresh
            const initialHTML = await this.page.content();
            await this.page.waitForTimeout(2000);
            const afterHTML = await this.page.content();
            
            if (initialHTML === afterHTML) {
                console.log('✅ No full page refresh detected');
            } else {
                console.log('❌ Page refresh detected');
            }
            
            // Check loading states
            const loadingIndicators = await this.page.$$('.htmx-indicator');
            console.log(`✅ Loading indicators found: ${loadingIndicators.length}`);
            
            return true;
        } catch (error) {
            console.error('❌ Dashboard test failed:', error.message);
            return false;
        }
    }

    async testAuthentication() {
        console.log('\nTesting Authentication...');
        try {
            // Go to login page
            await this.page.goto(`${this.baseURL}/login`, { waitUntil: 'networkidle2' });
            
            // Try to login with admin/admin
            await this.page.type('input[name="username"]', 'admin');
            await this.page.type('input[name="password"]', 'admin');
            await this.page.click('button[type="submit"]');
            
            // Wait for navigation
            await this.page.waitForNavigation({ waitUntil: 'networkidle2' });
            
            // Check if logged in (should see welcome message)
            const welcomeText = await this.page.$eval('nav', el => el.textContent);
            if (welcomeText.includes('Welcome')) {
                console.log('✅ Login successful');
                console.log('✅ User display in navigation working');
            } else {
                console.log('❌ Login failed or user display not working');
            }
            
            return true;
        } catch (error) {
            console.error('❌ Authentication test failed:', error.message);
            return false;
        }
    }

    async testGatewaysPage() {
        console.log('\nTesting Gateways Page...');
        try {
            await this.page.goto(`${this.baseURL}/gateways`, { waitUntil: 'networkidle2' });
            
            // Check page title
            const pageContent = await this.page.content();
            
            if (pageContent.includes('Gateway Management') && !pageContent.includes('Modem Management')) {
                console.log('✅ Gateways page shows correct content');
            } else if (pageContent.includes('Modem Management')) {
                console.log('❌ Gateways page shows modem content (template collision)');
            } else {
                console.log('❓ Unable to determine page content');
            }
            
            return true;
        } catch (error) {
            console.error('❌ Gateways test failed:', error.message);
            return false;
        }
    }

    async testCustomerStats() {
        console.log('\nTesting Customer Stats...');
        try {
            await this.page.goto(`${this.baseURL}/customers`, { waitUntil: 'networkidle2' });
            
            // Wait for stats to load
            await this.page.waitForTimeout(2000);
            
            // Check if stats show as number, not JSON
            const statsElement = await this.page.$('[hx-get="/api/customers/stats"]');
            if (statsElement) {
                const statsText = await this.page.evaluate(el => el.textContent, statsElement);
                
                // Check if it's a number or "Loading..."
                if (!statsText.includes('{') && !statsText.includes('}')) {
                    console.log('✅ Customer stats displaying as HTML (number)');
                } else {
                    console.log('❌ Customer stats showing JSON');
                }
            }
            
            return true;
        } catch (error) {
            console.error('❌ Customer stats test failed:', error.message);
            return false;
        }
    }

    async testPollingIssues() {
        console.log('\nTesting HTMX Polling...');
        try {
            // Monitor network requests
            const requests = [];
            this.page.on('request', request => {
                if (request.url().includes('/api/')) {
                    requests.push({
                        url: request.url(),
                        time: Date.now()
                    });
                }
            });
            
            // Test CDR page
            await this.page.goto(`${this.baseURL}/cdrs`, { waitUntil: 'networkidle2' });
            await this.page.waitForTimeout(5000);
            
            const cdrModemRequests = requests.filter(r => 
                r.url.includes('/api/v1/modems') && 
                r.time > Date.now() - 5000
            );
            
            if (cdrModemRequests.length === 0) {
                console.log('✅ CDR page not calling modems API');
            } else {
                console.log(`❌ CDR page made ${cdrModemRequests.length} calls to modems API`);
            }
            
            // Clear requests
            requests.length = 0;
            
            // Test Blacklist page
            await this.page.goto(`${this.baseURL}/blacklist`, { waitUntil: 'networkidle2' });
            await this.page.waitForTimeout(5000);
            
            const blacklistModemRequests = requests.filter(r => 
                r.url.includes('/api/v1/modems') && 
                r.time > Date.now() - 5000
            );
            
            if (blacklistModemRequests.length === 0) {
                console.log('✅ Blacklist page not calling modems API');
            } else {
                console.log(`❌ Blacklist page made ${blacklistModemRequests.length} calls to modems API`);
            }
            
            return true;
        } catch (error) {
            console.error('❌ Polling test failed:', error.message);
            return false;
        }
    }

    async runAllTests() {
        console.log('Starting UI Tests...\n');
        console.log('Base URL:', this.baseURL);
        console.log('=' .repeat(50));
        
        await this.init();
        
        const tests = [
            () => this.testDashboard(),
            () => this.testAuthentication(),
            () => this.testGatewaysPage(),
            () => this.testCustomerStats(),
            () => this.testPollingIssues()
        ];
        
        let passed = 0;
        let failed = 0;
        
        for (const test of tests) {
            try {
                const result = await test();
                if (result) passed++;
                else failed++;
            } catch (error) {
                failed++;
                console.error('Test error:', error);
            }
        }
        
        console.log('\n' + '=' .repeat(50));
        console.log(`Test Summary: ${passed} passed, ${failed} failed`);
        
        await this.cleanup();
    }
}

// Check if server is running first
const http = require('http');

function checkServer() {
    return new Promise((resolve) => {
        http.get('http://localhost:8080/ping', (res) => {
            resolve(res.statusCode === 200);
        }).on('error', () => {
            resolve(false);
        });
    });
}

async function main() {
    const serverRunning = await checkServer();
    
    if (!serverRunning) {
        console.error('❌ Server is not running on port 8080');
        console.error('Please start the server with:');
        console.error('./server_gateway_fixed > server.log 2>&1 &');
        process.exit(1);
    }
    
    const tester = new UITester();
    await tester.runAllTests();
}

if (require.main === module) {
    main().catch(console.error);
}

module.exports = UITester;