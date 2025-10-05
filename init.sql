-- Create messages table
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    author VARCHAR(255) NOT NULL DEFAULT 'Anonymous',
    sent BOOLEAN NOT NULL DEFAULT false,
    sent_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on created_at for better query performance
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at DESC);

-- Insert 50 unsent messages for scheduler testing
INSERT INTO messages (content, author, sent) VALUES 
    ('Order #1001 has been confirmed', 'E-commerce System', false),
    ('Welcome to our platform!', 'Onboarding System', false),
    ('Your payment was successful', 'Payment Gateway', false),
    ('New product launch: Check out our latest items', 'Marketing Team', false),
    ('Your subscription is about to expire', 'Subscription Service', false),
    ('Password reset requested', 'Security System', false),
    ('Weekly newsletter: Top 10 articles', 'Newsletter Bot', false),
    ('Your order #1002 has been shipped', 'Logistics System', false),
    ('Flash sale: 50% off on selected items', 'Sales Team', false),
    ('Account verification required', 'Security System', false),
    ('New comment on your post', 'Social Platform', false),
    ('Your report is ready for download', 'Analytics System', false),
    ('Meeting reminder: Team standup in 30 minutes', 'Calendar System', false),
    ('System maintenance scheduled for tonight', 'IT Operations', false),
    ('Your friend request was accepted', 'Social Platform', false),
    ('Invoice #2023-001 is now available', 'Billing System', false),
    ('Security alert: New login from unknown device', 'Security System', false),
    ('Your wishlist item is now on sale', 'E-commerce System', false),
    ('New follower: John Doe started following you', 'Social Platform', false),
    ('Backup completed successfully', 'Backup System', false),
    ('Your cart has items waiting for you', 'E-commerce System', false),
    ('Monthly summary: Your activity report', 'Analytics System', false),
    ('Support ticket #12345 has been resolved', 'Support System', false),
    ('New feature available: Dark mode enabled', 'Product Team', false),
    ('Your trial period ends in 3 days', 'Subscription Service', false),
    ('Congratulations! You earned a new badge', 'Gamification System', false),
    ('Server health check: All systems operational', 'Monitoring System', false),
    ('Your download is ready', 'File System', false),
    ('New message from customer support', 'Support System', false),
    ('Price drop alert: Item you viewed is now cheaper', 'E-commerce System', false),
    ('Reminder: Complete your profile for better experience', 'Onboarding System', false),
    ('Your code review has new comments', 'Development Team', false),
    ('Database backup scheduled at 2 AM', 'Database Admin', false),
    ('API rate limit warning: 80% quota used', 'API Gateway', false),
    ('Your post received 100 likes!', 'Social Platform', false),
    ('System update available: Version 2.0', 'Update Service', false),
    ('Your order is out for delivery', 'Logistics System', false),
    ('Fraud alert: Suspicious activity detected', 'Security System', false),
    ('Welcome bonus: 500 points added to your account', 'Rewards System', false),
    ('Your referral code was used by 5 people', 'Referral System', false),
    ('New job opportunity matching your profile', 'Job Board', false),
    ('Weather alert: Heavy rain expected tomorrow', 'Weather Service', false),
    ('Your file upload completed', 'File System', false),
    ('Team member invited you to a workspace', 'Collaboration Tool', false),
    ('Your application was approved', 'Admin System', false),
    ('Cache cleared successfully', 'Cache System', false),
    ('New update from your favorite creator', 'Content Platform', false),
    ('Your survey response was recorded', 'Survey System', false),
    ('Deployment successful: Production environment updated', 'CI/CD Pipeline', false),
    ('Thank you for your feedback!', 'Feedback System', false);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for updated_at
DROP TRIGGER IF EXISTS update_messages_updated_at ON messages;
CREATE TRIGGER update_messages_updated_at
    BEFORE UPDATE ON messages
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();