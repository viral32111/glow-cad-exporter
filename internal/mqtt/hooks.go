package mqtt

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"

	mochi "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"

	"github.com/viral32111/glow-cad-exporter/internal/config"
)

type Hooks struct {
	mochi.HookBase

	Config config.Config
}

func (hook *Hooks) ID() string { return "glow" }

func (hook *Hooks) Provides(buffer byte) bool {
	return bytes.Contains([]byte{
		mochi.OnPublished,
	}, []byte{buffer})
}

func (hook *Hooks) OnPublished(mqttClient *mochi.Client, mqttPacket packets.Packet) {
	// glow/000000000000/STATE
	// glow/000000000000/SENSOR/electricitymeter
	// glow/000000000000/SENSOR/gasmeter
	topicParts := strings.Split(mqttPacket.TopicName, "/")
	if len(topicParts) < 3 || topicParts[0] != hook.Config.MQTTTopic {
		return
	}

	switch strings.ToLower(strings.Join(topicParts[2:], "/")) {
		case "state":
			var message StateMessage
			if json.Unmarshal(mqttPacket.Payload, &message) == nil {
				slog.Info("Received state update", "msg", message)

			}

		case "sensor/electricitymeter":
			var message ElectricityMessage
			if json.Unmarshal(mqttPacket.Payload, &message) == nil {
				slog.Info("Received electricity meter update", "msg", message)
			}

		case "sensor/gasmeter":
			var message GasMessage
			if json.Unmarshal(mqttPacket.Payload, &message) == nil {
				slog.Info("Received gas meter update", "msg", message)
			}

		default:
			slog.Warn("Received MQTT message with unknown topic", "topic", mqttPacket.TopicName, "payload", string(mqttPacket.Payload))
	}
}
