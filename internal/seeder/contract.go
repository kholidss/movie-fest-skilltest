package seeder

import "context"

type Seederer interface {
	AdminData(ctx context.Context)
	GenresData(ctx context.Context)
}
