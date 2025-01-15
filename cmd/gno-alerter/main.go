package main

import (
	"flag"
	"strings"
	"time"

	"github.com/gnolang/gno/gno.land/pkg/gnoclient"
	rpcclient "github.com/gnolang/gno/tm2/pkg/bft/rpc/client"
	"github.com/sirupsen/logrus"

	"github.com/gnoverse/gno-alerters/internal/alerters"
	"github.com/gnoverse/gno-alerters/internal/config"
)

var (
	configFilePath = flag.String("config", "config.toml", "Path to the configuration file")
	debugMode      = flag.Bool("debug", false, "Enable debug mode")
)

func main() {
	flag.Parse()

	cfg, err := config.ParseConfig(*configFilePath)
	if err != nil {
		logrus.WithError(err).Fatal()
	}

	rpcClient, err := rpcclient.NewHTTPClient(cfg.RPC.Endpoint)
	if err != nil {
		logrus.WithError(err).Fatal()
	}

	client := gnoclient.Client{
		RPCClient: rpcClient,
	}

	var alerter alerters.Alerter

	if *debugMode == true {
		logrus.Info("Debug mode enabled")
		alerter = &alerters.MockAlerter{}
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		alerter = alerters.NewSlackAlerter(cfg, cfg.Slack.Token, cfg.Slack.ChannelID)
		logrus.SetLevel(logrus.InfoLevel)
	}

	err = start(cfg, alerter, client)
	if err != nil {
		logrus.WithError(err).Error()
	}
}

type validator struct {
	Name         string
	Addr         string
	MissedBlocks int
	isSigning    bool
}

type validatorsMissedBlocksMap map[string]*validator

func start(cfg *config.Config, alerter alerters.Alerter, client gnoclient.Client) error {
	latestHeight, err := client.LatestBlockHeight()
	if err != nil {
		return err
	}

	validators, err := client.RPCClient.Validators(&latestHeight)
	if err != nil {
		return err
	}
	validatorsMap := make(validatorsMissedBlocksMap, len(validators.Validators))
	for _, v := range validators.Validators {
		validatorsMap[v.Address.String()] = &validator{
			Name:         "", // TODO
			Addr:         v.Address.String(),
			MissedBlocks: 0,

			isSigning: false,
		}
	}

	previousBlock, err := client.Block(latestHeight - 1)
	if err != nil {
		return err
	}

	l := logrus.WithFields(logrus.Fields{
		"block.height": latestHeight,
	})

	stallCounter := 0
	l.Info("Starting gno-alerter")

	// Start the infinite loop
	for {
		time.Sleep(time.Second * 1)

		block, err := client.RPCClient.Block(nil)
		if err != nil {
			if strings.Contains(err.Error(), `unable to call RPC method block, invalid status code`) {
				// RPC IS DOWN
				continue
			}
			l.WithError(err).Error()
			continue
		} else if block.Block.Height == previousBlock.Block.Height {
			stallCounter += 1

			// stallCounter ensure we come here at least 3 times, there is an issue where block time diff
			// is between block_n+1->block_n+2 rather than block_n -> block_n+1
			if stallCounter > 3 && time.Since(block.Block.Time) > cfg.GetStalledPeriod() {
				if err := alerter.NetworkIsStalled(previousBlock.Block); err != nil {
					l.WithError(err).WithFields(logrus.Fields{
						"block.height": block.Block.Header.Height,
						"block.time":   block.Block.Header.Time,
					}).Error("failed to alert: Network is stalled !!")
				}
			}

			continue
		}
		stallCounter = 0

		l = l.WithFields(logrus.Fields{
			"block.height": block.Block.Header.Height,
			"block.time":   block.Block.Header.Time,
		})

		if alerter.IsStalled() {
			stalledTime := block.Block.Time.Sub(previousBlock.Block.Header.Time)
			if err := alerter.RecoverNetworkIsStalled(block.Block, stalledTime); err != nil {
				l.Error("failed to alert: Network is recovered from stalled mode")
			}
		}

		missedValidators := 0

		for addr := range validatorsMap {
			validatorsMap[addr].isSigning = false
		}

		for _, val := range block.Block.LastCommit.Precommits {
			if val != nil {
				addr := val.ValidatorAddress.String()
				validatorsMap[addr].isSigning = true

				if validatorsMap[addr].MissedBlocks >= cfg.Alerts.ConsecutiveMissed[0] {
					if err := alerter.RecoverMissingBlocks(addr); err != nil {
						l.WithError(err).WithFields(logrus.Fields{
							"validator.addr": addr,
						}).Error("failed to alert: recover missing blocks")
					}
					validatorsMap[addr].MissedBlocks = 0
				}

			} else {
				missedValidators += 1
			}
		}

		for addr, val := range validatorsMap {
			if !val.isSigning {
				validatorsMap[addr].MissedBlocks += 1
				l.WithFields(logrus.Fields{
					"validator.addr":         addr,
					"validator.missed_blocs": validatorsMap[addr].MissedBlocks,
				}).Debug()
			}

			if cfg.HasConsecutiveMissed(val.MissedBlocks) {
				if err := alerter.MissingBlocks(addr, val.MissedBlocks); err != nil {
					l.WithError(err).WithFields(logrus.Fields{
						"validator.addr":          addr,
						"validator.missed_blocks": val.MissedBlocks,
					}).Errorf("failed to alert: missing blocks for val: %s, missed blocks: %v", addr, val.MissedBlocks)
				}
			}
		}

		previousBlock = block

		l.Debugf("%d/%d signed block", len(validatorsMap)-missedValidators, len(validatorsMap))
	}
}
