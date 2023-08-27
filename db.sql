CREATE TABLE IF NOT EXISTS users
(
    id   serial PRIMARY KEY,
    name varchar(255)  not null
)

CREATE TABLE IF NOT EXISTS segments
(
    id   serial PRIMARY KEY,
    slug varchar(255)  not null
)
