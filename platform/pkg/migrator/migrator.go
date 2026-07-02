package migrator

import "context"

type Migrator interface {
	Up(context.Context) error
}
