CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(100)        NOT NULL,
    surname       VARCHAR(100)        NOT NULL,
    email         VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE groups
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE group_members
(
    group_id  INTEGER REFERENCES groups (id) ON DELETE CASCADE,
    user_id   INTEGER REFERENCES users (id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id)
);

CREATE TABLE expenses
(
    id          SERIAL PRIMARY KEY,
    group_id    INTEGER REFERENCES groups (id) ON DELETE CASCADE,
    payer_id    INTEGER        REFERENCES users (id) ON DELETE SET NULL,
    amount      DECIMAL(12, 2) NOT NULL,
    description TEXT           NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE expense_splits
(
    expense_id INTEGER REFERENCES expenses (id) ON DELETE CASCADE,
    user_id    INTEGER REFERENCES users (id) ON DELETE CASCADE,
    amount     DECIMAL(12, 2) NOT NULL,
    PRIMARY KEY (expense_id, user_id)
);
