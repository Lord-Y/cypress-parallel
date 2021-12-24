ALTER TABLE executions ALTER COLUMN spec SET NOT NULL;
-- ALTER TABLE executions ALTER COLUMN result SET DEFAULT {};
ALTER TABLE executions ALTER COLUMN execution_error_output SET DEFAULT '';
ALTER TABLE executions ALTER COLUMN pod_name SET DEFAULT '';