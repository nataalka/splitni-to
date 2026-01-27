CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(
    id            UUID                NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    name          VARCHAR(100)        NOT NULL,
    surname       VARCHAR(100)        NOT NULL,
    email         VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT,
    created_at    TIMESTAMP WITH TIME ZONE     DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS friendships
(
    user_id1   UUID REFERENCES users (id) ON DELETE CASCADE,
    user_id2   UUID REFERENCES users (id) ON DELETE CASCADE,
    status     VARCHAR(20)              DEFAULT 'PENDING', -- PENDING, ACCEPTED
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id1, user_id2)
);

CREATE TABLE IF NOT EXISTS groups
(
    id          UUID         NOT NULL    DEFAULT uuid_generate_v4() PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT         NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS group_members
(
    group_id  UUID REFERENCES groups (id) ON DELETE CASCADE,
    user_id   UUID REFERENCES users (id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id)
);

CREATE TABLE IF NOT EXISTS expenses
(
    id          UUID           NOT NULL  DEFAULT uuid_generate_v4() PRIMARY KEY,
    group_id    UUID REFERENCES groups (id) ON DELETE CASCADE,
    payer_id    UUID           REFERENCES users (id) ON DELETE SET NULL,
    amount      DECIMAL(12, 2) NOT NULL,
    currency    VARCHAR(3)               DEFAULT 'EUR',
    description TEXT           NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS expense_splits
(
    expense_id UUID REFERENCES expenses (id) ON DELETE CASCADE,
    user_id    UUID REFERENCES users (id) ON DELETE CASCADE,
    amount     DECIMAL(12, 2) NOT NULL,
    PRIMARY KEY (expense_id, user_id)
);
