package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kholidss/movie-fest-skilltest/internal/entity"
	"github.com/kholidss/movie-fest-skilltest/pkg/database/mysql"
	"github.com/kholidss/movie-fest-skilltest/pkg/helper"
	"github.com/kholidss/movie-fest-skilltest/pkg/tracer"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type genreRepository struct {
	db mysql.Adapter
}

func NewGenreRepository(db mysql.Adapter) GenreRepository {
	return &genreRepository{
		db: db,
	}
}

func (g *genreRepository) Store(ctx context.Context, payload any, opts ...Option) error {
	var (
		err error
		tx  *sql.Tx
	)

	ctx, span := tracer.NewSpan(ctx, "GenreRepo.Store", nil)
	defer span.End()

	opt := &option{}
	for _, f := range opts {
		f(opt)
	}

	if opt.tx != nil {
		tx = opt.tx
	} else {
		tx, err = g.db.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})
		if err != nil {
			tracer.AddSpanError(span, err)
			return err
		}

		defer func() {
			err = tx.Commit()
			if err != nil {
				tracer.AddSpanError(span, err)
				err = errors.Wrap(err, "failed to commit")
			}
		}()
	}

	query, val, err := helper.StructQueryInsertMysql(payload, TableNameGenres, "db", false)

	_, err = tx.ExecContext(
		ctx,
		query,
		val...,
	)
	if err != nil {
		tracer.AddSpanError(span, err)
		return err
	}

	return err
}

func (g *genreRepository) Update(ctx context.Context, payload any, where any, opts ...Option) error {
	var (
		err error
		tx  *sql.Tx
	)

	ctx, span := tracer.NewSpan(ctx, "GenreRepo.Update", nil)
	defer span.End()

	opt := &option{}
	for _, f := range opts {
		f(opt)
	}

	if opt.tx != nil {
		tx = opt.tx
	} else {
		tx, err = g.db.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})
		if err != nil {
			tracer.AddSpanError(span, err)
			return err
		}

		defer func() {
			err = tx.Commit()
			if err != nil {
				err = errors.Wrap(err, "failed to commit")
			}
		}()
	}

	q, vals, err := helper.StructToQueryUpdateMysql(payload, where, TableNameGenres, "db")
	if err != nil {
		tracer.AddSpanError(span, err)
		return err
	}

	_, err = tx.ExecContext(ctx, q, vals...)
	if err != nil {
		tracer.AddSpanError(span, err)
		return err
	}

	return err
}

func (g *genreRepository) FindOne(ctx context.Context, param any, selectColumn []string) (*entity.Genre, error) {
	var (
		dest entity.Genre
	)

	ctx, span := tracer.NewSpan(ctx, "GenreRepo.FindOne", nil)
	defer span.End()

	wq, vals, _, _, err := helper.StructQueryWhereMysql(param, true, "db")
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	err = g.db.QueryRow(ctx, &dest, fmt.Sprintf(DefaultQueryFindOne, helper.SelectCustom(selectColumn), TableNameGenres, wq), vals...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return &dest, nil
}

func (g *genreRepository) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Genre, error) {
	var (
		dest []entity.Genre
	)

	ctx, span := tracer.NewSpan(ctx, "GenreRepo.Finds", nil)
	defer span.End()

	wq, vals, _, _, err := helper.StructQueryWhereMysql(param, true, "db")
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	err = g.db.Query(ctx, &dest, fmt.Sprintf(DefaultQueryFinds, helper.SelectCustom(selectColumns), TableNameGenres, wq), vals...)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return dest, nil
}

func (g *genreRepository) ListMostView(ctx context.Context, meta entity.MetaPagination, selectColumns []string) ([]entity.Genre, int, error) {
	var (
		dest  []entity.Genre
		count int

		offset = helper.PageToOffset(meta.Limit, meta.Page)
	)

	ctx, span := tracer.NewSpan(ctx, "GenreRepo.ListMostView", nil)
	defer span.End()

	q := `
			SELECT %s
				FROM %s
			WHERE is_deleted = false
			ORDER BY view_number DESC
			LIMIT ? OFFSET ?;`

	qCount := `	SELECT COUNT(id)
					FROM %s
				WHERE is_deleted = false
				ORDER BY view_number DESC;`

	gr, _ := errgroup.WithContext(ctx)

	gr.Go(func() error {
		return g.db.Query(
			ctx,
			&dest,
			fmt.Sprintf(q, helper.SelectCustom(selectColumns), TableNameGenres),
			meta.Limit,
			offset,
		)
	})
	gr.Go(func() error {
		return g.db.QueryRow(
			ctx,
			&count,
			fmt.Sprintf(qCount, TableNameGenres),
		)
	})

	err := gr.Wait()
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, 0, err
	}

	return dest, count, nil
}

func (g *genreRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return g.db.BeginTx(ctx, opts)
}
