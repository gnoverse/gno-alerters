package alerters

import (
	"time"

	core_types "github.com/gnolang/gno/tm2/pkg/bft/types"
)

type Alerter interface {
	NetworkIsStalled(block *core_types.Block) error
	RecoverNetworkIsStalled(block *core_types.Block, stalledTime time.Duration) error
	IsStalled() bool

	MissingBlocks(addr string, missedBlocks int) error
	RecoverMissingBlocks(addr string) error
}
