package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/dmitrykharchenko95/otus_user/internal/database/entity"
)

type (
	Manager interface {
		Add(ctx context.Context, u *entity.User) (int64, error)
		Get(ctx context.Context, id int64) (*entity.User, error)
		Delete(ctx context.Context, id int64) error
		Update(ctx context.Context, u *entity.User) error

		MigrateUp(ctx context.Context) error
		MigrateDown(ctx context.Context) error

		Ping() error
	}
	manager struct {
		db *sql.DB
	}

	Config struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
	}
)

func New(config Config) (Manager, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	var (
		db  *sql.DB
		err error
	)
	if db, err = sql.Open("postgres", config.dsn()); err != nil {
		return nil, fmt.Errorf("open connection error: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("connection to database error: %w", err)
	}

	return &manager{db: db}, nil
}

func (c *Config) dsn() string {
	var (
		dsn = &url.URL{
			Scheme: "postgresql",
			Host:   fmt.Sprintf("%s:%s", c.Host, c.Port),
			Path:   c.Database,
		}
		q = dsn.Query()
	)

	q.Add("sslmode", "disable")
	q.Add("binary_parameters", "yes")
	dsn.RawQuery = q.Encode()

	if c.Username == "" {
		return dsn.String()
	}

	if c.Password == "" {
		dsn.User = url.User(c.Username)

		return dsn.String()
	}

	dsn.User = url.UserPassword(c.Username, c.Password)

	return dsn.String()
}

func (m *manager) Add(ctx context.Context, u *entity.User) (int64, error) {
	var id int64
	if err := m.db.QueryRowContext(
		ctx,
		`
			INSERT INTO users(
				username, first_name, last_name, email, phone
			)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`,
		u.Username,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (m *manager) Get(ctx context.Context, id int64) (*entity.User, error) {
	var u = &entity.User{}
	if err := m.db.QueryRowContext(
		ctx,
		"SELECT id, username, first_name, last_name, email, phone FROM users WHERE id = $1",
		id,
	).Scan(
		&u.Id,
		&u.Username,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Phone,
	); err != nil {
		return nil, err
	}

	return u, nil
}

func (m *manager) Delete(ctx context.Context, id int64) error {
	var _, err = m.db.ExecContext(
		ctx,
		"DELETE FROM users WHERE id = $1",
		id,
	)

	return err
}

func (m *manager) Update(ctx context.Context, u *entity.User) error {
	var setQuery, args = u.GetUpdateQuery()
	log.Println(*u)
	if len(args) == 0 {
		return nil
	}
	var q = fmt.Sprintf("UPDATE public.users %s", setQuery)
	var _, err = m.db.ExecContext(ctx, q, args...)
	return err
}

func (m *manager) Ping() error {
	return m.db.Ping()
}

func (c *Config) validate() error {
	if c.Host == "" {
		return errors.New("empty db host")
	}

	if _, err := strconv.Atoi(c.Port); err != nil {
		return fmt.Errorf("db port parsing error: %w", err)
	}

	if c.Port == "" {
		return errors.New("empty db port")
	}

	if c.Database == "" {
		return errors.New("empty db name")
	}
	return nil
}
