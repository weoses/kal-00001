CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE permissions (
    id smallint PRIMARY KEY,
    code text UNIQUE NOT NULL
);

INSERT INTO permissions (id, code) VALUES
    (1, 'CREATE'),
    (2, 'UPDATE'),
    (3, 'DELETE'),
    (4, 'PUBLISH');

CREATE TABLE user_permissions (
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permission_id smallint NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);

CREATE TABLE integration_telegram (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    telegram_ids bigint[] NOT NULL,
    answer_unknown boolean NOT NULL DEFAULT false
);

CREATE INDEX integration_telegram_telegram_ids_idx ON integration_telegram USING gin (telegram_ids);

CREATE TABLE integration_webapp_basic (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    username text UNIQUE NOT NULL,
    password_hash text NOT NULL
);
