CREATE TABLE IF NOT EXISTS users
(
    id   serial PRIMARY KEY,
    name varchar(255)  not null
)

CREATE TABLE IF NOT EXISTS segments
(
    id   serial PRIMARY KEY,
    slug varchar(255)  not null,
    UNIQUE (slug)
)


CREATE TABLE IF NOT EXISTS user_segments
(
    user_id int,
    segment_id int,
    percent smallint,
    created_at timestamp WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    expired_date timestamp WITHOUT TIME ZONE,
    PRIMARY KEY(user_id, segment_id)
);

CREATE TABLE IF NOT EXISTS segments_history
(
    id   serial PRIMARY KEY,
    user_id int,
    segment_id int,
    operation varchar(20),
    created_at timestamp WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

