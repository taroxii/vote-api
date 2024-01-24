package postgres

import (
	"context"
	"database/sql"

	"github.com/taroxii/vote-api/pkg/entity"
	"github.com/taroxii/vote-api/pkg/utils/logger"
	"go.uber.org/zap"
)

type postgresRepository struct {
	Conn *sql.DB
}

func NewPostgresUserRepository(conn *sql.DB) entity.UserRepository {
	return &postgresRepository{
		Conn: conn,
	}
}

func (m *postgresRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []entity.User, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		errRow := rows.Close()
		if err != nil {
			logger.Logger.Error("Error", zap.Error(errRow))
		}
	}()
	result = make([]entity.User, 0)

	for rows.Next() {
		t := entity.User{}
		err = rows.Scan(
			&t.ID,
			&t.Username,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			logger.Logger.Error("Error", zap.Error(err))
		}
		result = append(result, t)
	}
	return result, nil
}

func (m *postgresRepository) GetById(ctx context.Context, id int64) (res entity.User, err error) {
	query := `SELECT id, username, updated_at, created_at FROM users where id = $1`
	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, entity.ErrNotFound
	}
	return
}

func (m *postgresRepository) GetByUsername(ctx context.Context, username string) (res entity.User, err error) {
	query := `SELECT id, username, created_at, updated_at
  						FROM users WHERE username = $1`

	list, err := m.fetch(ctx, query, username)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, entity.ErrNotFound
	}
	return
}
