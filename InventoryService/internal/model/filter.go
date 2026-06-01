package model

import (
	"github.com/google/uuid"
)

type Filter struct {
	Uuids      []uuid.UUID
	Names      []string
	Categories []Category
	Countries  []string
	Tags       []string
}
