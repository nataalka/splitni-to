CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id            UUID                NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    name          VARCHAR(100)        NOT NULL,
    surname       VARCHAR(100)        NOT NULL,
    email         VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT,
    created_at    TIMESTAMP WITH TIME ZONE     DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE groups
(
    id         UUID         NOT NULL    DEFAULT uuid_generate_v4() PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE group_members
(
    group_id  UUID REFERENCES groups (id) ON DELETE CASCADE,
    user_id   UUID REFERENCES users (id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id)
);

CREATE TABLE expenses
(
    id          UUID           NOT NULL  DEFAULT uuid_generate_v4() PRIMARY KEY,
    group_id    UUID REFERENCES groups (id) ON DELETE CASCADE,
    payer_id    UUID           REFERENCES users (id) ON DELETE SET NULL,
    amount      DECIMAL(12, 2) NOT NULL,
    description TEXT           NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE expense_splits
(
    expense_id UUID REFERENCES expenses (id) ON DELETE CASCADE,
    user_id    UUID REFERENCES users (id) ON DELETE CASCADE,
    amount     DECIMAL(12, 2) NOT NULL,
    PRIMARY KEY (expense_id, user_id)
);
