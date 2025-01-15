package alerters

import (
	"time"

	core_types "github.com/gnolang/gno/tm2/pkg/bft/types"
	"github.com/sirupsen/logrus"
)

// MockAlerter
type MockAlerter struct {
	isStalled bool
}

func (a MockAlerter) IsStalled() bool { return a.isStalled }

func (a MockAlerter) MissingBlocks(addr string, missedBlocks int) error {
	logrus.WithFields(logrus.Fields{
		"addr":         addr,
		"missed_block": missedBlocks,
	}).Warn("Validator is missing blocks...")
	return nil
}

func (a MockAlerter) RecoverMissingBlocks(addr string) error {
	logrus.WithFields(logrus.Fields{
		"addr": addr,
	}).Info("Validator is signing again")
	return nil
}

func (a *MockAlerter) NetworkIsStalled(block *core_types.Block) error {
	if a.isStalled {
		return nil
	}
	a.isStalled = true

	logrus.WithFields(logrus.Fields{
		"block.height": block.Header.Height,
	}).Warn("No recent blocks")
	return nil
}

func (a *MockAlerter) RecoverNetworkIsStalled(block *core_types.Block, stalledTime time.Duration) error {
	a.isStalled = false
	logrus.WithFields(logrus.Fields{
		"block.height":     block.Header.Height,
		"stalled_duration": stalledTime,
	}).Info("Network recovered")
	return nil
}
