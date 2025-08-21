-- Initialize the database with the starterkit database
-- This script runs when the PostgreSQL container starts for the first time
-- Create the starterkit database
CREATE DATABASE starterkit;
-- Grant all privileges to postgres user on starterkit database
GRANT ALL PRIVILEGES ON DATABASE starterkit TO postgres;