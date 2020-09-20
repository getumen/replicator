package handler

import (
	"github.com/getumen/replicator/pkg/models"
	"github.com/getumen/replicator/pkg/store"
)

// Handler applies command log to FSM
type Handler interface {
	Apply(store.Store, *models.Command) error
	ApplyBatch(store.Store, []*models.Command) error
}
