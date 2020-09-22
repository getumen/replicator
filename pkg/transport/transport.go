package transport

//go:generate mockgen -source=$GOFILE -destination=mock$GOPACKAGE/mock_$GOFILE -package=mock$GOPACKAGE

import "github.com/hashicorp/raft"

// Transport is a transport layer of raft
type Transport interface {
	raft.Transport
}
