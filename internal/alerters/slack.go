package alerters

import (
	"fmt"
	"time"

	core_types "github.com/gnolang/gno/tm2/pkg/bft/types"
	"github.com/gnoverse/gno-alerters/internal/config"
	"github.com/slack-go/slack"
)

type SlackAlerter struct {
	cfg *config.Config

	chainID        string
	slackClient    *slack.Client
	slackChannelID string
	isStalled      bool
}

// NewSlackAlerter
func NewSlackAlerter(cfg *config.Config, token string, channelID string) *SlackAlerter {
	if cfg.ChainID == "" {
		cfg.ChainID = "dev"
	}

	return &SlackAlerter{
		cfg: cfg,

		chainID:        cfg.ChainID,
		slackClient:    slack.New(token),
		slackChannelID: channelID,
		isStalled:      false,
	}
}

func (s *SlackAlerter) NetworkIsStalled(block *core_types.Block) error {
	if s.isStalled {
		return nil
	}
	s.isStalled = true

	msg := fmt.Sprintf(`[%s] Network is stalled on block: %v`, s.chainID, block.Header.Height)
	_, _, err := s.slackClient.PostMessage(
		s.slackChannelID,
		slack.MsgOptionAttachments(
			slack.Attachment{
				Color: "danger",
				Text:  msg,
			}))
	return err
}

func (s *SlackAlerter) RecoverNetworkIsStalled(block *core_types.Block, stalledTime time.Duration) error {
	msg := fmt.Sprintf("[%s] Network has stalled for: %s, recovered at block: %v", s.chainID, stalledTime, block.Header.Height)
	s.isStalled = false

	_, _, err := s.slackClient.PostMessage(
		s.slackChannelID,
		slack.MsgOptionAttachments(
			slack.Attachment{
				Color: "good",
				Text:  msg,
			}))
	return err
}

func (s *SlackAlerter) IsStalled() bool {
	return s.isStalled
}

func (s *SlackAlerter) MissingBlocks(addr string, missedBlocks int) error {
	msg := fmt.Sprintf("[%s] Node %s has missed %d blocks", s.chainID, s.cfg.GetValidatorMoniker(addr), missedBlocks)

	_, _, err := s.slackClient.PostMessage(
		s.slackChannelID,

		slack.MsgOptionAttachments(
			slack.Attachment{
				Color: "danger",
				Text:  msg,
			}))
	return err
}

func (s *SlackAlerter) RecoverMissingBlocks(addr string) error {
	msg := fmt.Sprintf("[%s] Node %s has recovered from missing blocks", s.chainID, s.cfg.GetValidatorMoniker(addr))

	_, _, err := s.slackClient.PostMessage(
		s.slackChannelID,
		slack.MsgOptionAttachments(
			slack.Attachment{
				Color: "good",
				Text:  msg,
			}))
	return err
}
