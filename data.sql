CREATE TABLE users (
    id serial PRIMARY KEY NOT NULL
);

CREATE TABLE segment (
    id serial PRIMARY KEY NOT NULL,
    name varchar(100) NOT NULL UNIQUE,
    active bool NOT NULL
);

CREATE TABLE "user_segment" (
    id serial PRIMARY KEY NOT NULL,
    user_id serial,
    segment_id serial,
    active bool NOT NULL,
    crt_at timestamptz NOT NULL DEFAULT (now()),
    del_at timestamptz,
    del_after timestamptz
);

ALTER TABLE user_segment ADD FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE user_segment ADD FOREIGN KEY (segment_id) REFERENCES segment (id);