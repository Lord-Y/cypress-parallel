ALTER TABLE annotations
ADD CONSTRAINT fk_annotations_projects
FOREIGN KEY (project_id)
REFERENCES projects(project_id)
ON DELETE CASCADE
ON UPDATE CASCADE;