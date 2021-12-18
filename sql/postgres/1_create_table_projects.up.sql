CREATE TABLE projects (
  project_id SERIAL PRIMARY KEY,
  project_name varchar(100) NOT NULL,
  date timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);