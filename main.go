package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/otakakot/sample-go-postgres-assist/pkg/sqlb"
	"github.com/otakakot/sample-go-postgres-assist/pkg/sqlc"
)

func main() {
	// --------------------------------------------
	// 初期化等々
	// --------------------------------------------

	ctx := context.Background()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		panic("DATABASE_URL is required")
	}

	conn, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, conn)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	db := stdlib.OpenDBFromPool(pool)

	defer func() {
		if err := db.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}

	// --------------------------------------------
	// 正常系
	// --------------------------------------------

	// User と Todo を同時に作成
	// User に紐づく Todo を取得
	// User に紐づく Todo を削除

	now := time.Now()

	userID := uuid.New()

	slog.InfoContext(ctx, fmt.Sprintf("gen user id: %s", userID))

	if err := Tx(ctx, db, []func(ctx context.Context, tx *sql.Tx) error{
		func(ctx context.Context, tx *sql.Tx) error {
			entity := sqlb.User{
				ID:        userID,
				Name:      uuid.NewString(),
				CreatedAt: now,
				UpdatedAt: now,
			}

			if err := entity.Insert(ctx, tx, boil.Infer()); err != nil {
				return fmt.Errorf("failed to insert user: %w", err)
			}

			return nil
		},
		func(ctx context.Context, tx *sql.Tx) error {
			entity := sqlb.Todo{
				ID:        uuid.New(),
				CreatedAt: now,
				UpdatedAt: now,
				Title:     uuid.NewString(),
				Completed: false,
				UserID:    userID,
			}

			if err := entity.Insert(ctx, tx, boil.Infer()); err != nil {
				return fmt.Errorf("failed to insert todo: %w", err)
			}

			return nil
		},
		func(ctx context.Context, tx *sql.Tx) error {
			entity := sqlb.Todo{
				ID:        uuid.New(),
				CreatedAt: now,
				UpdatedAt: now,
				Title:     uuid.NewString(),
				Completed: false,
				UserID:    userID,
			}

			if err := entity.Insert(ctx, tx, boil.Infer()); err != nil {
				return fmt.Errorf("failed to insert todo: %w", err)
			}

			return nil
		},
	}); err != nil {
		panic(err)
	}

	{
		user, err := sqlc.New(pool).FindUserByID(ctx, userID)
		if err != nil {
			panic(err)
		}

		slog.InfoContext(ctx, fmt.Sprintf("got user id: %s", user.ID))

	}

	{
		todos, err := sqlc.New(pool).ListTodoByUserID(ctx, userID)
		if err != nil {
			panic(err)
		}

		for _, todo := range todos {
			slog.InfoContext(ctx, fmt.Sprintf("todo: %+v", todo))
		}
	}

	if err := Tx(ctx, db, []func(ctx context.Context, tx *sql.Tx) error{
		func(ctx context.Context, tx *sql.Tx) error {
			if _, err := sqlb.Todos(sqlb.TodoWhere.UserID.EQ(userID)).DeleteAll(ctx, tx); err != nil {
				return fmt.Errorf("failed to delete todo: %w", err)
			}

			return nil
		},
	}); err != nil {
		panic(err)
	}

	{
		todos, err := sqlc.New(pool).ListTodoByUserID(ctx, userID)
		if err != nil {
			panic(err)
		}

		if len(todos) != 0 {
			panic("failed to delete todo")
		}
	}

	// User と Todo を同時に作成
	// User に紐づく Todo を取得
	// User を削除
	userID2 := uuid.New()

	if err := Tx(ctx, db, []func(ctx context.Context, tx *sql.Tx) error{
		func(ctx context.Context, tx *sql.Tx) error {
			entity := sqlb.User{
				ID:        userID2,
				Name:      uuid.NewString(),
				CreatedAt: now,
				UpdatedAt: now,
			}

			if err := entity.Insert(ctx, tx, boil.Infer()); err != nil {
				return fmt.Errorf("failed to insert user: %w", err)
			}

			return nil
		},
		func(ctx context.Context, tx *sql.Tx) error {
			entity := sqlb.Todo{
				ID:        uuid.New(),
				CreatedAt: now,
				UpdatedAt: now,
				Title:     uuid.NewString(),
				Completed: false,
				UserID:    userID2,
			}

			if err := entity.Insert(ctx, tx, boil.Infer()); err != nil {
				return fmt.Errorf("failed to insert todo: %w", err)
			}

			return nil
		},
		func(ctx context.Context, tx *sql.Tx) error {
			entity := sqlb.Todo{
				ID:        uuid.New(),
				CreatedAt: now,
				UpdatedAt: now,
				Title:     uuid.NewString(),
				Completed: false,
				UserID:    userID2,
			}

			if err := entity.Insert(ctx, tx, boil.Infer()); err != nil {
				return fmt.Errorf("failed to insert todo: %w", err)
			}

			return nil
		},
	}); err != nil {
		panic(err)
	}

	{
		user, err := sqlc.New(pool).FindUserByID(ctx, userID2)
		if err != nil {
			panic(err)
		}

		slog.InfoContext(ctx, fmt.Sprintf("got user id: %s", user.ID))
	}

	{
		todos, err := sqlc.New(pool).ListTodoByUserID(ctx, userID2)
		if err != nil {
			panic(err)
		}

		for _, todo := range todos {
			slog.InfoContext(ctx, fmt.Sprintf("todo: %+v", todo))
		}
	}

	if err := Tx(ctx, db, []func(ctx context.Context, tx *sql.Tx) error{
		func(ctx context.Context, tx *sql.Tx) error {
			target := sqlb.User{
				ID: userID2,
			}

			if _, err := target.Delete(ctx, tx); err != nil {
				return fmt.Errorf("failed to delete user: %w", err)
			}

			return nil
		},
	}); err != nil {
		panic(err)
	}

	// User は存在しない
	{
		if _, err := sqlc.New(pool).FindUserByID(ctx, userID2); err != nil {
			_ = Handle(ctx, err)
		}
	}

	// User に紐づく Todo も存在しない
	{
		todos, err := sqlc.New(pool).ListTodoByUserID(ctx, userID2)
		if err != nil {
			panic(err)
		}

		if len(todos) != 0 {
			panic("failed to delete todo")
		}
	}

	// --------------------------------------------
	// 異常系
	// --------------------------------------------

	if err := Tx(ctx, db, []func(ctx context.Context, tx *sql.Tx) error{
		func(ctx context.Context, tx *sql.Tx) error {
			todo := sqlb.Todo{
				ID:        uuid.New(),
				CreatedAt: now,
				UpdatedAt: now,
				Title:     uuid.NewString(),
				Completed: false,
				UserID:    uuid.New(), // 存在しないUserID
			}

			if err := todo.Insert(ctx, tx, boil.Infer()); err != nil {
				return fmt.Errorf("failed to insert todo: %w", err)
			}

			return nil
		},
	}); err != nil {
		_ = Handle(ctx, err)
	}

	if err := Tx(ctx, db, []func(ctx context.Context, tx *sql.Tx) error{
		func(ctx context.Context, tx *sql.Tx) error {
			exists := sqlb.User{
				ID:        userID,
				Name:      uuid.NewString(),
				CreatedAt: now,
				UpdatedAt: now,
			}

			if err := exists.Insert(ctx, db, boil.Infer()); err != nil {
				return fmt.Errorf("failed to insert user: %w", err)
			}

			return nil
		},
	}); err != nil {
		_ = Handle(ctx, err)
	}
}

// Tx executes the given queries in a transaction.
func Tx(
	ctx context.Context,
	db *sql.DB,
	queries []func(ctx context.Context, tx *sql.Tx) error,
) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return Handle(ctx, fmt.Errorf("failed to begin transaction: %w", err))
	}

	for _, query := range queries {
		if err := query(ctx, tx); err != nil {
			if err := tx.Rollback(); err != nil {
				return Handle(ctx, fmt.Errorf("failed to rollback transaction: %w", err))
			}

			return Handle(ctx, fmt.Errorf("failed to execute query: %w", err))
		}
	}

	if err := tx.Commit(); err != nil {
		return Handle(ctx, fmt.Errorf("failed to commit transaction: %w", err))
	}

	return nil
}

// Handle handles the given error.
func Handle(
	ctx context.Context,
	err error,
) error {
	if err == nil {
		return nil
	}

	var target *pgconn.PgError

	if errors.As(err, &target) {
		// ref: https://pkg.go.dev/github.com/jackc/pgerrcode
		if pgerrcode.IsIntegrityConstraintViolation(target.Code) {
			slog.WarnContext(ctx, target.Error())

			return err
		}
	}

	slog.ErrorContext(ctx, err.Error())

	return err
}
