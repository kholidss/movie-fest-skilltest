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

type movieRepository struct {
	db mysql.Adapter
}

func NewMovieRepository(db mysql.Adapter) MovieRepository {
	return &movieRepository{
		db: db,
	}
}

func (m movieRepository) Store(ctx context.Context, payload any, opts ...Option) error {
	var (
		err error
		tx  *sql.Tx
	)

	ctx, span := tracer.NewSpan(ctx, "MovieRepo.Store", nil)
	defer span.End()

	opt := &option{}
	for _, f := range opts {
		f(opt)
	}

	if opt.tx != nil {
		tx = opt.tx
	} else {
		tx, err = m.db.BeginTx(ctx, &sql.TxOptions{
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

	query, val, err := helper.StructQueryInsertMysql(payload, TableNameMovies, "db", false)

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

func (m movieRepository) Update(ctx context.Context, payload any, where any, opts ...Option) error {
	var (
		err error
		tx  *sql.Tx
	)

	ctx, span := tracer.NewSpan(ctx, "MovieRepo.Update", nil)
	defer span.End()

	opt := &option{}
	for _, f := range opts {
		f(opt)
	}

	if opt.tx != nil {
		tx = opt.tx
	} else {
		tx, err = m.db.BeginTx(ctx, &sql.TxOptions{
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

	q, vals, err := helper.StructToQueryUpdateMysql(payload, where, TableNameMovies, "db")
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

func (m *movieRepository) FindOne(ctx context.Context, param any, selectColumn []string) (*entity.Movie, error) {
	var (
		dest entity.Movie
	)

	ctx, span := tracer.NewSpan(ctx, "MovieRepo.FindOne", nil)
	defer span.End()

	wq, vals, _, _, err := helper.StructQueryWhereMysql(param, true, "db")
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	err = m.db.QueryRow(ctx, &dest, fmt.Sprintf(DefaultQueryFindOne, helper.SelectCustom(selectColumn), TableNameMovies, wq), vals...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return &dest, nil
}

func (m *movieRepository) FindOneWithForUpdate(ctx context.Context, param any, opts ...Option) (*entity.Movie, error) {
	var (
		tx  *sql.Tx
		res entity.Movie
	)

	ctx, span := tracer.NewSpan(ctx, "MovieRepo.FindOneWithForUpdate", nil)
	defer span.End()

	wq, vals, _, _, err := helper.StructQueryWhereMysql(param, true, "db")
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	q := `SELECT
			id,
			title,
			genre_ids,
			view_number
			FROM %s %s
			LIMIT 1
			FOR UPDATE;`

	opt := &option{}
	for _, f := range opts {
		f(opt)
	}

	if opt.tx != nil {
		tx = opt.tx
	} else {
		tx, err = m.db.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})
		if err != nil {
			tracer.AddSpanError(span, err)
			return nil, err
		}

		defer func() {
			err = tx.Commit()
			if err != nil {
				tracer.AddSpanError(span, err)
				err = errors.Wrap(err, "failed to commit")
			}
		}()
	}

	err = opt.tx.QueryRowContext(ctx, fmt.Sprintf(q, TableNameMovies, wq), vals...).Scan(
		&res.ID,
		&res.Title,
		&res.GenreIDS,
		&res.ViewNumber,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return &res, nil
}

func (m *movieRepository) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Movie, error) {
	var (
		dest []entity.Movie
	)

	ctx, span := tracer.NewSpan(ctx, "MovieRepo.Finds", nil)
	defer span.End()

	wq, vals, _, _, err := helper.StructQueryWhereMysql(param, true, "db")
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	err = m.db.Query(ctx, &dest, fmt.Sprintf(DefaultQueryFinds, helper.SelectCustom(selectColumns), TableNameMovies, wq), vals...)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return dest, nil
}

func (m *movieRepository) ListMostView(ctx context.Context, meta entity.MetaPagination, selectColumns []string) ([]entity.Movie, int, error) {
	var (
		dest  []entity.Movie
		count int

		offset = helper.PageToOffset(meta.Limit, meta.Page)
	)

	ctx, span := tracer.NewSpan(ctx, "MovieRepo.ListMostView", nil)
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
		return m.db.Query(
			ctx,
			&dest,
			fmt.Sprintf(q, helper.SelectCustom(selectColumns), TableNameMovies),
			meta.Limit,
			offset,
		)
	})
	gr.Go(func() error {
		return m.db.QueryRow(
			ctx,
			&count,
			fmt.Sprintf(qCount, TableNameMovies),
		)
	})

	err := gr.Wait()
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, 0, err
	}

	return dest, count, nil
}

func (m *movieRepository) List(ctx context.Context, meta entity.MetaPagination, selectColumns []string) ([]entity.Movie, int, error) {
	var (
		dest  []entity.Movie
		count int

		offset = helper.PageToOffset(meta.Limit, meta.Page)
	)

	ctx, span := tracer.NewSpan(ctx, "MovieRepo.List", nil)
	defer span.End()

	q := `
			SELECT %s
				FROM %s
			WHERE is_deleted = false
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?;`

	qCount := `	SELECT COUNT(id)
					FROM %s
				WHERE is_deleted = false
				ORDER BY created_at DESC;`

	gr, _ := errgroup.WithContext(ctx)

	gr.Go(func() error {
		return m.db.Query(
			ctx,
			&dest,
			fmt.Sprintf(q, helper.SelectCustom(selectColumns), TableNameMovies),
			meta.Limit,
			offset,
		)
	})
	gr.Go(func() error {
		return m.db.QueryRow(
			ctx,
			&count,
			fmt.Sprintf(qCount, TableNameMovies),
		)
	})

	err := gr.Wait()
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, 0, err
	}

	return dest, count, nil
}

func (m movieRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return m.db.BeginTx(ctx, opts)
}
