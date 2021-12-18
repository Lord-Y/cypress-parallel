CREATE TABLE environments (
  environment_id SERIAL PRIMARY KEY,
  key varchar(100) NOT NULL,
  value varchar(100) NOT NULL,
  project_id INT NOT NULL
);