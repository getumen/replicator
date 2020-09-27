package mockserializer

import (
	"testing"

	"github.com/getumen/replicator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

func TestMessagepack_Unmarshal(t *testing.T) {
	b, err := msgpack.Marshal(
		&models.Command{
			Operation: "PUT",
			Key:       []byte("foo"),
			Value:     []byte("bar"),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	var command models.Command
	err = msgpack.Unmarshal(b, &command)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "PUT", command.Operation)
	assert.Equal(t, []byte("foo"), command.Key)
	assert.Equal(t, []byte("bar"), command.Value)
}
