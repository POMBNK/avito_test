CREATE TABLE users (
    id serial PRIMARY KEY NOT NULL,
    name varchar(50),
    email varchar(50)
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

-- users stubs
INSERT INTO users (name, email) VALUES ('Ilya','Il@ya.ru');
INSERT INTO users (name, email) VALUES ('Vitya','Vit@ya.ru');
INSERT INTO users (name, email) VALUES ('Kolya','Kol@ya.ru');
INSERT INTO users (name, email) VALUES ('Katya','k@ya.ru');
INSERT INTO users (name, email) VALUES ('pbInya','pbIn@ya.ru');

--segment stubs
INSERT INTO segment (name, active) VALUES ('discount30','1');
INSERT INTO segment (name, active) VALUES ('discount50','1');
INSERT INTO segment (name, active) VALUES ('discount80','1');
INSERT INTO segment (name, active) VALUES ('voice_msg','1');