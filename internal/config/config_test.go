package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	content := `
[rpc]
endpoint = "http://localhost:26657"

[slack]
token      = "xoxb-xxxx"
channel_id = "XXXXXXX"

[alerts]
stalled_period = "30s"
consecutive_missed = [20, 100, 500]
`
	f, err := os.CreateTemp("", "config.toml")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	_, err = f.Write([]byte(content))
	require.NoError(t, err)

	cfg, err := ParseConfig(f.Name())
	require.NoError(t, err)

	require.NotEmpty(t, cfg.RPC.Endpoint)
	require.NotEmpty(t, cfg.Slack.Token)
	require.NotEmpty(t, cfg.Slack.ChannelID)

	require.Len(t, cfg.Alerts.StalledPeriod, 3)

	require.NoError(t, err)
	require.Equal(t, time.Second*30, cfg.GetStalledPeriod())
}
