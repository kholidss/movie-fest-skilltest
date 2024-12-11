package seeder

import "context"

type Seederer interface {
	AdminData(ctx context.Context)
}
