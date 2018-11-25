-- +migrate Up
CREATE TABLE users (
                             id SERIAL PRIMARY KEY,
                             uuid uuid NOT NULL UNIQUE,
                             email VARCHAR(320) NOT NULL UNIQUE
);
CREATE index users_uuid ON users (uuid);

-- +migrate Down
DROP TABLE users;