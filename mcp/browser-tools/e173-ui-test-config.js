// E173 Gateway UI Test Configuration
module.exports = {
  baseURL: 'http://192.168.1.35:8080',
  credentials: {
    username: 'admin',
    password: 'admin'
  },
  tests: [
    {
      name: 'Dashboard Layout Test',
      path: '/dashboard',
      checks: [
        { type: 'element', selector: '#stats-cards', description: 'Stats cards container' },
        { type: 'console', level: 'error', description: 'Check for console errors' },
        { type: 'network', url: '/api/stats/', description: 'Stats API calls' }
      ]
    },
    {
      name: 'Gateway Page Test',
      path: '/gateways',
      checks: [
        { type: 'element', selector: '.gateway-card', description: 'Gateway cards' },
        { type: 'button', selector: '.btn-test', description: 'Test button functionality' }
      ]
    },
    {
      name: 'Customer Management Test',
      path: '/customers',
      checks: [
        { type: 'button', selector: '#add-customer-btn', description: 'Add customer button' },
        { type: 'button', selector: '.btn-edit', description: 'Edit buttons' }
      ]
    },
    {
      name: 'CDR Page Test',
      path: '/cdrs',
      checks: [
        { type: 'element', selector: 'table', description: 'CDR table structure' },
        { type: 'element', selector: 'thead', description: 'Table headers' }
      ]
    }
  ]
};
