package handler

//go:generate mockgen -source=$GOFILE -destination=mock$GOPACKAGE/mock_$GOFILE -package=mock$GOPACKAGE

import (
	"github.com/getumen/replicator/pkg/models"
	"github.com/getumen/replicator/pkg/store"
)

// Handler applies command log to FSM
type Handler interface {
	Apply(store.Store, *models.Command) error
	ApplyBatch(store.Store, []*models.Command) error
}
