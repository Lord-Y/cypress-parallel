CREATE TABLE teams (
  team_id SERIAL PRIMARY KEY,
  team_name varchar(100) NOT NULL,
  date timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);