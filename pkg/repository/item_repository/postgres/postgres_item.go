package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/taroxii/vote-api/pkg/entity"
	"github.com/taroxii/vote-api/pkg/repository"
	"github.com/taroxii/vote-api/pkg/utils/logger"
	"go.uber.org/zap"
)

type postgresRepository struct {
	Conn *sql.DB
}

func NewPostgresItemRepository(conn *sql.DB) entity.ItemRepository {
	return &postgresRepository{
		Conn: conn,
	}
}

func (m *postgresRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []entity.Item, err error) {
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
	result = make([]entity.Item, 0)

	for rows.Next() {
		t := entity.Item{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.VoteCount,
		)
		if err != nil {
			logger.Logger.Error("Error", zap.Error(err))
		}
		result = append(result, t)
	}
	return result, nil
}

func (m *postgresRepository) Fetch(ctx context.Context, cursor string, num int64) (res []entity.Item, nextCursor string, err error) {
	query := `SELECT id, name, description, vote_count, created_at FROM items WHERE created_at > $1 ORDER BY vote_count LIMIT $2 `
	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", entity.ErrBadParamInput
	}
	res, err = m.fetch(ctx, query, pq.FormatTimestamp(decodedCursor), num)
	if err != nil {
		return nil, "", err
	}
	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return
}

func (m *postgresRepository) Insert(ctx context.Context, a *entity.Item) (err error) {
	query := `INSERT INTO items (name, description, vote_count, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	var lastId int64
	err = stmt.QueryRowContext(ctx, a.Name, a.Description, a.VoteCount, a.UpdatedAt, a.CreatedAt).Scan(&lastId)

	if err != nil {
		logger.Logger.Error("Can't insert to db", zap.Error(err), zap.Any("values", a))
		return err
	}

	a.ID = lastId
	return
}

func (m *postgresRepository) Delete(ctx context.Context, id int64) (err error) {
	var vote_count int

	tx, err := m.Conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return err
	}

	result := tx.QueryRowContext(ctx, "SELECT vote_count from items where id = $1", id)
	err = result.Scan(&vote_count)

	if err != nil {
		return
	}
	if vote_count != 0 {
		return entity.ErrInvalidOperationConstraint
	}

	query := "DELETE FROM items WHERE id = $1"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", rowsAffected)
		return
	}
	tx.Commit()
	return
}

func (m *postgresRepository) ClearVote(ctx context.Context, id int64) (item entity.Item, err error) {
	tx, err := m.Conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		logger.Logger.Error("Database tx error ", zap.Error(err))
	}
	res, err := tx.ExecContext(ctx, `DELETE FROM user_items WHERE item_id = $1`, id)
	if err != nil {
		return
	}

	affected, err := res.RowsAffected()

	if err != nil {
		logger.Logger.Error("Database error ", zap.Error(err))
		return

	}

	if affected == 0 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", affected)
		return
	}

	res, err = tx.Exec(`UPDATE items SET vote_count = $1, updated_at = now() where id = $2`, 0, id)
	if err != nil {
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", affect)
		return
	}
	row := tx.QueryRowContext(ctx, `SELECT id, name, description, vote_count, created_at, updated_at FROM items WHERE id = $1`, id)

	if err = row.Scan(&item.ID, &item.Name, &item.Description, &item.VoteCount, &item.CreatedAt, &item.UpdatedAt); err != nil {
		return
	}
	tx.Commit()
	return
}

func (m *postgresRepository) Vote(ctx context.Context, id, user_id int64) (item *entity.Item, err error) {
	uid := uuid.New()
	tx, err := m.Conn.Begin()
	if err != nil {
		return
	}

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO user_items (correlation_id, user_id, item_id, created_at, updated_at) values ($1, $2, $3, now(), now())`)
	if err != nil {
		logger.Logger.Error("Error during prepare context", zap.Error(err))
		return
	}
	res, err := stmt.ExecContext(ctx, uid.String(), user_id, id)
	if err != nil {
		if err.(*pq.Error).Code == "23505" {
			err = entity.ErrConflict
		}

		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", affect)
		return
	}
	if err != nil {
		logger.Logger.Error("Can't insert to db", zap.Error(err), zap.Any("correlationId", uid))
		return
	}

	row := tx.QueryRow("SELECT id,name,description,vote_count,created_at,updated_at FROM items WHERE id = $1", id)
	var itq entity.Item
	err = row.Scan(&itq.ID, &itq.Name, &itq.Description, &itq.VoteCount, &itq.CreatedAt, &itq.UpdatedAt)
	if err != nil {
		return
	}
	itq.VoteCount += 1
	res, err = tx.Exec(`UPDATE items SET vote_count = $1, updated_at = now() where id = $2`, itq.VoteCount, id)
	if err != nil {
		return
	}
	affect, err = res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", affect)
		return
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &itq, nil
}

func (m *postgresRepository) Update(ctx context.Context, item *entity.Item) (err error) {
	tx, err := m.Conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	results := tx.QueryRowContext(ctx, `SELECT vote_count from items WHERE id = $1`, item.ID)

	err = results.Scan(&item.VoteCount)
	if item.VoteCount > 0 {
		return entity.ErrInvalidOperationConstraint
	}

	if err != nil {
		return
	}

	query := `UPDATE items set name = $1, description= $2, updated_at = now() WHERE id = $3`

	// sql.Row
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, item.Name, item.Description, item.ID)
	if err != nil {
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", affect)
		return
	}
	tx.Commit()

	return
}

//
// func (m *postgresRepository) FindByID(
// 	ctx context.Context,
// 	id int64,
// 	options ...usecase.RepositoryOption,
// ) (entity.Item, error) {
// 	opts := usecase.MergeOptions(options...)
// 	db, _ := repository.ResolveDB(m.Conn, opts.Transaction())
// 	query := repository.ResolveBaseQuery(db, opts)

// 	var item entity.Item
// 	res := query.QueryRowContext(ctx, "SELECT id, name,description, vote_count, created_at, updated_at WHERE id = $1 ", id)
// 	err := res.Scan(&item.ID, &item.Name, &item.Description, &item.VoteCount, &item.CreatedAt, &item.UpdatedAt)
// 	if err != nil {
// 		var zero entity.Item
// 		return zero, entity.ErrNotFound
// 	}
// 	return item, nil
// }

// func (m *postgresRepository) UpdateTx(ctx context.Context, item *entity.Item, uu *entity.UserItems, options ...usecase.RepositoryOption) error {
// 	opts := usecase.MergeOptions(options...)
// 	db, useTx := repository.ResolveDB(m.Conn, opts.Transaction())
// 	f := func(tx *sql.Tx) error {
// 		stmt, err := tx.Prepare("UPDATE items SET name = $1, description = $2, vote_count = $3, created_at = $4, updated_at =$5")
// 		if err != nil {
// 			return err
// 		}
// 		sqlResults, err := stmt.Exec(item.ID, item.Name, item.Description, item.VoteCount)
// 		if err != nil {
// 			return err
// 		}
// 		affect, err := sqlResults.RowsAffected()
// 		if err != nil {
// 			return err
// 		}
// 		if affect == 0 {
// 			return entity.ErrNotFound
// 		}

// 		stmt, err = tx.Prepare("INSERT INTO user_items (correlation_id, user_id, item_id, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)")
// 		if err != nil {
// 			return err
// 		}
// 		sqlResults, err = stmt.Exec(item.ID, item.Name, item.Description, item.VoteCount)
// 		if err != nil {
// 			return err
// 		}
// 		affect, err = sqlResults.RowsAffected()
// 		if err != nil {
// 			return err
// 		}
// 		if affect == 0 {
// 			return entity.ErrNotFound
// 		}
// 		tx.Commit()
// 		return nil
// 	}

// 	if useTx {
// 		tx := repository.ResolveBaseQuery(db, opts)
// 		return f(tx)
// 	} else {
// 		return nil
// 	}
// }
