ALTER TABLE environments
ADD CONSTRAINT fk_environments_projects
FOREIGN KEY (project_id)
REFERENCES projects(project_id)
ON DELETE CASCADE
ON UPDATE CASCADE;