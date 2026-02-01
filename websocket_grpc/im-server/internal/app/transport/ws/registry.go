package ws

import (
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/transport/ws/contract"
	wsregistry "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/transport/ws/registry"
)

// Registry is kept as a package-level alias for convenience.
// The actual contract is defined in ws/contract so business code can depend on it
// without importing the ws implementation package.
type Registry = contract.Registry

func NewRegistry() Registry {
	return wsregistry.NewRegistry()
}
