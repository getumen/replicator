package serializer

import "github.com/getumen/replicator/pkg/models"

// Serializer serialize and deserialize Command
type Serializer interface {
	Serialize(*models.Command) ([]byte, error)
	Deserialize([]byte) (*models.Command, error)
}
