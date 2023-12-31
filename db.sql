CREATE TABLE IF NOT EXISTS users
(
    id   serial PRIMARY KEY,
    name varchar(255)  not null
);

CREATE TABLE IF NOT EXISTS segments
(
    id   serial PRIMARY KEY,
    slug varchar(255) not null,
    percent smallint,
    UNIQUE (slug)
);


CREATE TABLE IF NOT EXISTS user_segments
(
    user_id int,
    segment_id int,
    created_at timestamp WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    expired_date timestamp WITHOUT TIME ZONE,
    PRIMARY KEY(user_id, segment_id)
);

CREATE TYPE segment_operation AS ENUM('add', 'delete');

CREATE TABLE IF NOT EXISTS segments_history
(
    id   serial PRIMARY KEY,
    user_id int not null,
    segment_id int not null,
    operation segment_operation not null,
    created_at timestamp WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);
