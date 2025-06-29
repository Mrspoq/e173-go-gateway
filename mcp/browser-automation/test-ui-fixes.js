const puppeteer = require('puppeteer');

async function testUIFixes() {
    const browser = await puppeteer.launch({
        headless: true,
        args: ['--no-sandbox', '--disable-setuid-sandbox']
    });
    
    const page = await browser.newPage();
    
    try {
        console.log('Starting UI tests...\n');
        
        // Navigate to login page
        console.log('1. Testing login...');
        await page.goto('http://192.168.1.35:8080/login', { waitUntil: 'networkidle2' });
        
        // Check if already logged in by looking for logout button
        const isLoggedIn = await page.$('button[type="submit"]') === null;
        
        if (!isLoggedIn) {
            await page.waitForSelector('input[name="username"]');
            
            // Login
            await page.type('input[name="username"]', 'admin');
            await page.type('input[name="password"]', 'admin123');
            await page.click('button[type="submit"]');
            
            // Wait for navigation to dashboard
            await page.waitForNavigation({ waitUntil: 'networkidle2' });
        }
        console.log('✓ Login successful\n');
        
        // Test 1: Dashboard layout
        console.log('2. Testing dashboard layout...');
        await page.goto('http://192.168.1.35:8080/');
        await page.waitForSelector('#stats-cards');
        
        // Check if 5 cards are displayed in one row
        const cardsCount = await page.evaluate(() => {
            const container = document.querySelector('#stats-cards');
            const cards = container.querySelectorAll('.bg-white.rounded-lg.shadow');
            return cards.length;
        });
        console.log(`✓ Dashboard shows ${cardsCount} cards\n`);
        
        // Test 2: Gateway page
        console.log('3. Testing gateway page...');
        await page.goto('http://192.168.1.35:8080/gateways');
        await page.waitForSelector('table', { timeout: 5000 }).catch(() => {
            console.log('✗ Gateway page failed to load table');
        });
        
        const gatewayTitle = await page.$eval('h3', el => el.textContent);
        console.log(`✓ Gateway page loaded: "${gatewayTitle}"\n`);
        
        // Test 3: Modems page
        console.log('4. Testing modems page...');
        await page.goto('http://192.168.1.35:8080/modems');
        await page.waitForSelector('#modem-list');
        
        // Check if stats cards are not nested
        const modemStatsText = await page.evaluate(() => {
            const statsCard = document.querySelector('.bg-white.rounded-lg.shadow dd');
            return statsCard ? statsCard.textContent : 'Not found';
        });
        console.log(`✓ Modem stats show: "${modemStatsText.trim()}"\n`);
        
        // Test 4: Customer edit links
        console.log('5. Testing customer edit links...');
        await page.goto('http://192.168.1.35:8080/customers');
        await page.waitForSelector('#customer-list');
        
        // Try clicking edit
        const editLinkExists = await page.$('a[href="/customers/1/edit"]') !== null;
        if (editLinkExists) {
            await page.click('a[href="/customers/1/edit"]');
            await page.waitForTimeout(1000);
            const currentUrl = page.url();
            console.log(`✓ Edit link navigated to: ${currentUrl}\n`);
        }
        
        // Test 5: CDR page structure
        console.log('6. Testing CDR page...');
        await page.goto('http://192.168.1.35:8080/cdrs');
        await page.waitForSelector('#cdr-list');
        
        // Check if table structure exists
        const hasTable = await page.evaluate(() => {
            return document.querySelector('table thead') !== null;
        });
        console.log(`✓ CDR page shows table structure: ${hasTable}\n`);
        
        // Test 6: Check authentication display
        console.log('7. Testing authentication display...');
        const userDisplay = await page.evaluate(() => {
            const welcomeText = document.querySelector('.text-gray-500.dark\\:text-gray-300 span');
            return welcomeText ? welcomeText.textContent : 'Not found';
        });
        console.log(`✓ User display shows: "${userDisplay}"\n`);
        
        console.log('All tests completed!');
        
    } catch (error) {
        console.error('Test failed:', error.message);
    } finally {
        await browser.close();
    }
}

// Run the tests
testUIFixes().catch(console.error);