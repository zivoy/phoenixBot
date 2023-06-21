package nats

import (
	"encoding/json"
	"time"
)

type DiscordVerifyRequest struct {
	DiscordName string `json:"discord_name,omitempty"`
	DiscordID   string `json:"discord_id,omitempty"`
	RSICode     string `json:"code"`
}

type DiscordVerifyResponse struct {
	DiscordID string `json:"discord_id"`
	Error     string `json:"error,omitempty"`
}

func (c *Connector) VerifyUser(discord *DiscordVerifyRequest) (*DiscordVerifyResponse, error) {
	if err := c.verifyConnected(); err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(discord)
	if err != nil {
		return nil, err
	}

	response, err := c.nc.Request("discord.function.verify-rsi", marshal, 500*time.Millisecond)
	//if err == nats.ErrNoResponders || err == nats.ErrTimeout
	if err != nil {
		return nil, err
	}

	data := new(DiscordVerifyResponse)
	err = json.Unmarshal(response.Data, data)

	return data, nil
}
