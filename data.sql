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


--csv problem
--more optimized query, works pretty fine
SELECT us.user_id, u.email, s.name, s.id,us.active,us.crt_at,us.del_at
FROM user_segment us
         JOIN segment s ON s.id = us.segment_id
         JOIN users u on u.id = us.user_id
WHERE us.user_id = 2
  AND  us.crt_at >= '2023-08-25 15:49:07.758484 +00:00'
  AND (us.del_at IS NULL OR NOW() >= del_at)
ORDER BY us.crt_at;

--pain in ass with x2 JOINs
SELECT us.user_id, u.email, s.name, s.id,'created' as action,us.crt_at AS date
FROM user_segment us
    JOIN segment s ON s.id = us.segment_id
    JOIN users u ON u.id = us.user_id
WHERE u.id = 2
  AND  us.crt_at >= '2023-08-25 15:49:07.758484 +00:00'
UNION ALL
SELECT us.user_id, u.email, s.name, s.id,'deleted' as action,us.del_at as date
FROM user_segment us
    JOIN segment s ON s.id = us.segment_id
    JOIN users u on u.id = us.user_id
WHERE u.id = 2
  AND us.active=FALSE
  AND NOW() >= del_at
ORDER BY date;

--get random users by given percent
SELECT * FROM users --change * to id
ORDER BY RANDOM()
    LIMIT (SELECT COUNT(*) FROM users) * 0.8


SELECT * FROM users
ORDER BY RANDOM()
    LIMIT (SELECT COUNT(*) FROM users) * 0.8;

--with ttl
INSERT INTO user_segment(segment_id, user_id, ACTIVE,del_after)
WITH data AS (SELECT id AS segment_id, 3 AS user_id, TRUE AS active,CAST('2023-08-25 15:49:07.758484 +00:00' AS TIMESTAMPTZ ) as del_after FROM segment WHERE name = 'discount80')
SELECT segment_id, user_id, active,del_after FROM data
WHERE NOT EXISTS (SELECT * FROM user_segment
                  WHERE (user_id = (select user_id from DATA) AND
                         segment_id = (select segment_id from DATA) AND
                         active = TRUE));

--ttl update
INSERT INTO user_segment(segment_id, user_id, ACTIVE,del_after)
WITH data AS (SELECT id AS segment_id, 3 AS user_id, TRUE AS active,
                     CASE WHEN '' = ''
                              THEN NULL
                          ElSE CAST('2023-08-25 15:49:07.758484 +00:00' AS TIMESTAMPTZ ) END as del_after FROM segment WHERE name = 'discount80')
SELECT segment_id, user_id, active,del_after FROM data
WHERE NOT EXISTS (SELECT * FROM user_segment
                  WHERE (user_id = (select user_id from DATA) AND
                         segment_id = (select segment_id from DATA) AND
                         active = TRUE));