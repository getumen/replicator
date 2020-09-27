package mockhandler

import (
	"fmt"

	"github.com/getumen/replicator/pkg/handler"
	models "github.com/getumen/replicator/pkg/models"
	store "github.com/getumen/replicator/pkg/store"
)

type handlerImpl struct {
}

// NewKVSHandler returns kvs handler
func NewKVSHandler() handler.Handler {
	return &handlerImpl{}
}

func (handlerImpl) Apply(st store.Store, command *models.Command) error {
	ba := st.CreateBatch()
	switch command.Operation {
	case "PUT":
		ba.Put(command.Key, command.Value)
	case "DELETE":
		ba.Delete(command.Key)
	default:
		return fmt.Errorf("unknown operation: %s", command.Operation)
	}
	return st.Write(ba)
}
func (handlerImpl) ApplyBatch(
	st store.Store,
	commands []*models.Command,
) error {
	ba := st.CreateBatch()
	for _, command := range commands {
		switch command.Operation {
		case "PUT":
			ba.Put(command.Key, command.Value)
		case "DELETE":
			ba.Delete(command.Key)
		default:
			return fmt.Errorf("unknown operation: %s", command.Operation)
		}
	}
	return st.Write(ba)
}
