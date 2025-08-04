-- Database initialization wrapper
-- This file creates the database and then executes the schema files in order

-- Create the database if it doesn't exist
SELECT 'CREATE DATABASE gophkeeper_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'gophkeeper_db')
\gexec

-- Connect to the database
\c gophkeeper_db;

-- Execute schema files in order (PostgreSQL will run them alphabetically)
-- 001_types.sql will be executed first
-- 002_tables.sql will be executed second
