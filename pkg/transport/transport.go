package transport

import "github.com/hashicorp/raft"

// Transport is a transport layer of raft
type Transport interface {
	raft.Transport
}
