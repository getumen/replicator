package serializer

//go:generate mockgen -source=$GOFILE -destination=mock$GOPACKAGE/mock_$GOFILE -package=mock$GOPACKAGE

import "github.com/getumen/replicator/pkg/models"

// Serializer serialize and deserialize Command
type Serializer interface {
	Serialize(*models.Command) ([]byte, error)
	Deserialize([]byte) (*models.Command, error)
}
