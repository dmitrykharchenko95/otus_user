-- +migrate Up
CREATE TABLE IF NOT EXISTS public.users
(
    id         BIGSERIAL PRIMARY KEY                  NOT NULL,
    username   VARCHAR(256)                           NOT NULL,
    first_name VARCHAR(256)                           NOT NULL,
    last_name  VARCHAR(256)                           NOT NULL,
    email      VARCHAR(256)                           NOT NULL,
    phone      VARCHAR(256)                           NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS public.users;