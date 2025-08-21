-- Seed data for the starterkit database
-- This file creates sample users for development and testing

-- Only insert if the users table exists and is empty
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users') THEN
        IF NOT EXISTS (SELECT 1 FROM users LIMIT 1) THEN
            INSERT INTO users (id, email, name, created_at, updated_at) VALUES
                ('123e4567-e89b-12d3-a456-426614174000', 'john.doe@example.com', 'John Doe', NOW() - INTERVAL '30 days', NOW() - INTERVAL '5 days'),
                ('223e4567-e89b-12d3-a456-426614174001', 'jane.smith@example.com', 'Jane Smith', NOW() - INTERVAL '25 days', NOW() - INTERVAL '3 days'),
                ('323e4567-e89b-12d3-a456-426614174002', 'bob.johnson@example.com', 'Bob Johnson', NOW() - INTERVAL '20 days', NOW() - INTERVAL '10 days'),
                ('423e4567-e89b-12d3-a456-426614174003', 'alice.williams@example.com', 'Alice Williams', NOW() - INTERVAL '15 days', NOW() - INTERVAL '1 day'),
                ('523e4567-e89b-12d3-a456-426614174004', 'charlie.brown@example.com', 'Charlie Brown', NOW() - INTERVAL '10 days', NOW());
            
            RAISE NOTICE 'Seed data: 5 users created successfully';
        ELSE
            RAISE NOTICE 'Seed data: Users table already contains data, skipping seed';
        END IF;
    ELSE
        RAISE WARNING 'Seed data: Users table does not exist. Run migrations first.';
    END IF;
END $$;