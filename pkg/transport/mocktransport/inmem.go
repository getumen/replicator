package mocktransport

import (
	"github.com/getumen/replicator/pkg/transport"
	"github.com/hashicorp/raft"
)

// NewInMemTransport returns in-memory transport
func NewInMemTransport(addr raft.ServerAddress) (raft.ServerAddress, transport.Transport) {
	return raft.NewInmemTransport(addr)
}
