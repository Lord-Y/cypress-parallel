ALTER TABLE executions ALTER COLUMN spec DROP NOT NULL;
-- ALTER TABLE executions ALTER COLUMN result DROP DEFAULT;
ALTER TABLE executions ALTER COLUMN execution_error_output DROP DEFAULT;
ALTER TABLE executions ALTER COLUMN pod_name DROP DEFAULT;
