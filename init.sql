create schema tasker;
create table tasker.tasker (
id serial PRIMARY KEY,
message VARCHAR(255),
status VARCHAR(255),
created timestamp default CURRENT_TIMESTAMP,
deadline timestamp);
