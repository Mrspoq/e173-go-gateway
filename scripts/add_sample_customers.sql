-- Sample customer data for testing customer management system
-- This will be removed after UI testing is complete

-- Insert sample customers
INSERT INTO customers (
    customer_code, company_name, contact_person, email, phone, 
    address, city, state, country, postal_code,
    account_status, credit_limit, current_balance, monthly_limit,
    timezone, preferred_currency, auto_recharge_enabled,
    auto_recharge_threshold, auto_recharge_amount, notes
) VALUES 
(
    'CUST001', 'TechCorp Solutions', 'John Smith', 'john@techcorp.com', '+1-555-0101',
    '123 Business Ave', 'New York', 'NY', 'USA', '10001',
    'active', 1000.00, 250.75, 500.00,
    'America/New_York', 'USD', true,
    50.00, 100.00, 'Premium customer with auto-recharge'
),
(
    'CUST002', 'Global Communications Ltd', 'Sarah Johnson', 'sarah@globalcomm.com', '+44-20-7946-0958',
    '456 International Blvd', 'London', 'England', 'UK', 'SW1A 1AA',
    'active', 2000.00, 150.25, 1000.00,
    'Europe/London', 'GBP', false,
    NULL, NULL, 'Enterprise customer - manual billing'
),
(
    'CUST003', 'StartupTech Inc', 'Mike Chen', 'mike@startuptech.com', '+1-555-0102',
    '789 Innovation Dr', 'San Francisco', 'CA', 'USA', '94105',
    'active', 500.00, 75.50, 250.00,
    'America/Los_Angeles', 'USD', true,
    25.00, 50.00, 'Growing startup customer'
),
(
    'CUST004', 'Suspended Corp', 'Jane Doe', 'jane@suspended.com', '+1-555-0103',
    '321 Inactive St', 'Chicago', 'IL', 'USA', '60601',
    'suspended', 1000.00, -25.00, 500.00,
    'America/Chicago', 'USD', false,
    NULL, NULL, 'Suspended due to negative balance'
);

-- Insert sample payments for these customers
INSERT INTO payments (
    customer_id, amount, payment_method, payment_reference,
    status, notes, created_by
) VALUES 
(
    (SELECT id FROM customers WHERE customer_code = 'CUST001'),
    100.00, 'credit_card', 'CC-2024-001',
    'completed', 'Monthly auto-recharge', 1
),
(
    (SELECT id FROM customers WHERE customer_code = 'CUST002'),
    200.00, 'bank_transfer', 'BT-2024-001',
    'completed', 'Manual payment received', 1
),
(
    (SELECT id FROM customers WHERE customer_code = 'CUST003'),
    50.00, 'credit_card', 'CC-2024-002',
    'completed', 'Auto-recharge payment', 1
);

-- Update customer current balances based on payments
UPDATE customers SET current_balance = 350.75 WHERE customer_code = 'CUST001';
UPDATE customers SET current_balance = 350.25 WHERE customer_code = 'CUST002';
UPDATE customers SET current_balance = 125.50 WHERE customer_code = 'CUST003';
