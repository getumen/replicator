package mockserializer

import (
	"log"

	models "github.com/getumen/replicator/pkg/models"
	"github.com/getumen/replicator/pkg/serializer"
	"github.com/vmihailenco/msgpack/v5"
)

type gobSerializer struct{}

// NewSerializer returns serializer
func NewSerializer() serializer.Serializer {
	return &gobSerializer{}
}

func (gobSerializer) Serialize(command *models.Command) ([]byte, error) {
	return msgpack.Marshal(command)
}

func (gobSerializer) Deserialize(data []byte) (*models.Command, error) {
	var command models.Command
	err := msgpack.Unmarshal(data, &command)
	if err != nil {
		return nil, err
	}
	log.Printf("deserialize: %v from %s with %+v", command, data, err)
	return &command, nil
}
