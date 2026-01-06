package storage

import (
	"github.com/kweaver-ai/operator-hub/operator-app/server/interfaces"
)

type Storage interface {
	Save(instance *interfaces.MCPServerInstance) error
	Get(mcpID string, version int) (*interfaces.MCPServerInstance, error)
	Delete(mcpID string, version int) error
	Exists(mcpID string, version int) bool
}
