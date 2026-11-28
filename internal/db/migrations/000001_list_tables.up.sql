-- List existing tables in the public schema (no schema changes)
SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';
