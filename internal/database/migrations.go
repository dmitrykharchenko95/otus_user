package database

import "context"

func (m *manager) MigrateUp(ctx context.Context) error {
	var _, err = m.db.ExecContext(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS public.users
(
    id            BIGSERIAL PRIMARY KEY                  NOT NULL,
    username      VARCHAR(256)                           NOT NULL,
    first_name    VARCHAR(256)                           NOT NULL,
    last_name     VARCHAR(256)                           NOT NULL,
    email         VARCHAR(256)                           NOT NULL,
    phone         VARCHAR(256)                           NOT NULL,
    salt          VARCHAR(256)                           NOT NULL,
    password_hash VARCHAR(256)                           NOT NULL,

    created_at    TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);`,
	)
	if err != nil {
		return err
	}

	if _, err = m.db.ExecContext(ctx, "CREATE UNIQUE INDEX  username_idx ON public.users (username);"); err != nil {
		return err
	}

	_, err = m.db.ExecContext(ctx, "CREATE UNIQUE INDEX email_idx ON public.users (email);")

	return err
}

func (m *manager) MigrateDown(ctx context.Context) error {
	var _, err = m.db.ExecContext(
		ctx,
		`DROP TABLE IF EXISTS public.users;`,
	)

	return err
}
