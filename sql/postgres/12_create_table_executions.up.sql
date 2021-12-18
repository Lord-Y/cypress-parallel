CREATE TABLE executions (
  execution_id SERIAL PRIMARY KEY,
  project_id INT NOT NULL,
  branch VARCHAR(100) NOT NULL,
  execution_status VARCHAR(20) DEFAULT 'NOT_STARTED',
  uniq_id VARCHAR(10) NOT NULL,
  spec TEXT,
  result json NOT NULL,
  date timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE executions
ADD CONSTRAINT fk_executions_projects
FOREIGN KEY (project_id)
REFERENCES projects(project_id)
ON DELETE CASCADE
ON UPDATE CASCADE;