package staging

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/treeverse/lakefs/pkg/db"
	"github.com/treeverse/lakefs/pkg/graveler"
	"github.com/treeverse/lakefs/pkg/logging"
)

type DBManager struct {
	db  db.Database
	log logging.Logger
}

func NewDBManager(db db.Database) *DBManager {
	return &DBManager{
		db:  db,
		log: logging.Default().WithField("service_name", "DB_staging_manager"),
	}
}

func (p *DBManager) Get(ctx context.Context, st graveler.StagingToken, key graveler.Key) (*graveler.Value, error) {
	value := &graveler.Value{}
	err := p.db.Get(ctx, value, "SELECT identity, data FROM graveler_staging_kv WHERE staging_token=$1 AND key=$2", st, key)
	if errors.Is(err, db.ErrNotFound) {
		return nil, graveler.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if value.Identity == nil {
		return nil, nil
	}
	return value, nil
}

func (p *DBManager) Set(ctx context.Context, st graveler.StagingToken, key graveler.Key, value *graveler.Value, overwrite bool) error {
	if value == nil {
		value = new(graveler.Value)
	} else if value.Identity == nil {
		return graveler.ErrInvalidValue
	}
	if !overwrite {
		res, err := p.db.Exec(ctx, `INSERT INTO graveler_staging_kv (staging_token, key, identity, data)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (staging_token, key) DO NOTHING`,
			st, key, value.Identity, value.Data)
		if err != nil {
			return err
		}
		if res.RowsAffected() == 0 {
			return graveler.ErrPreconditionFailed
		}
		return err
	}
	_, err := p.db.Exec(ctx, `INSERT INTO graveler_staging_kv (staging_token, key, identity, data)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (staging_token, key) DO UPDATE
			SET (staging_token, key, identity, data) =
			(excluded.staging_token, excluded.key, excluded.identity, excluded.data)`,
		st, key, value.Identity, value.Data)
	return err
}

func (p *DBManager) Update(_ context.Context, _ graveler.StagingToken, _ graveler.Key, _ graveler.ValueUpdateFunc) error {
	panic("not implemented")
}

func (p *DBManager) DropKey(ctx context.Context, st graveler.StagingToken, key graveler.Key) error {
	_, err := p.db.Exec(ctx, "DELETE FROM graveler_staging_kv WHERE staging_token=$1 AND key=$2", st, key)
	return err
}

// List returns an iterator of staged values on the staging token st
func (p *DBManager) List(ctx context.Context, st graveler.StagingToken, batchSize int) (graveler.ValueIterator, error) {
	return NewDBStagingIterator(ctx, p.db, p.log, st, batchSize), nil
}

func (p *DBManager) DropAsync(ctx context.Context, st graveler.StagingToken) error {
	return p.Drop(ctx, st)
}

func (p *DBManager) Drop(ctx context.Context, st graveler.StagingToken) error {
	_, err := p.db.Exec(ctx, "DELETE FROM graveler_staging_kv WHERE staging_token=$1", st)
	return err
}

func (p *DBManager) DropByPrefix(ctx context.Context, st graveler.StagingToken, prefix graveler.Key) error {
	upperBound := graveler.UpperBoundForPrefix(prefix)
	builder := sq.Delete("graveler_staging_kv").Where(sq.Eq{"staging_token": st}).Where("key >= ?::bytea", prefix)
	if upperBound != nil {
		builder = builder.Where("key < ?::bytea", upperBound)
	}
	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = p.db.Exec(ctx, query, args...)
	return err
}
