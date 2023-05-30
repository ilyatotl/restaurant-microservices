CREATE TABLE users
(
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(50) UNIQUE  NOT NULL,
    email         VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    role          VARCHAR(10)         NOT NULL CHECK (role IN ('customer', 'chef', 'manager')),
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE sessions
(
    id            BIGSERIAL PRIMARY KEY,
    user_id       INT          NOT NULL,
    session_token VARCHAR(255) NOT NULL,
    expires_at    TIMESTAMP    NOT NULL
);