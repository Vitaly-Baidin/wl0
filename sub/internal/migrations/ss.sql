CREATE TABLE caches
(
    id         serial PRIMARY KEY,
    key        VARCHAR(250) UNIQUE NOT NULL,
    value      jsonb,
    expiration bigint
);