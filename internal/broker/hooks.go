package broker

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"

	"github.com/viral32111/glow-cad-exporter/internal/config"
	"github.com/viral32111/glow-cad-exporter/internal/helpers"
)

type Hooks struct {
	mqtt.HookBase

	Config config.Config
}

func (this *Hooks) ID() string {
	return "glow"
}

func (this *Hooks) Provides(buffer byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnPublished,
	}, []byte{buffer})
}

func (this *Hooks) OnPublished(mqttClient *mqtt.Client, mqttPacket packets.Packet) {
	// glow/000000000000/STATE
	// glow/000000000000/SENSOR/electricitymeter
	// glow/000000000000/SENSOR/gasmeter
	topicParts := strings.Split(mqttPacket.TopicName, "/")
	if len(topicParts) < 3 || topicParts[0] != this.Config.MQTTTopic {
		return
	}

	switch strings.ToLower(strings.Join(topicParts[2:], "/")) {
	case "state":
		var message StateMessage
		if json.Unmarshal(mqttPacket.Payload, &message) == nil {
			slog.Info("Received state update", "msg", message)

			// Normalise MAC addresses
			message.EthernetMAC = strings.ToLower(helpers.NormaliseMACAddress(message.EthernetMAC))
			message.WirelessMAC = strings.ToLower(helpers.NormaliseMACAddress(message.WirelessMAC))

			// Normalise Zigbee EUI-64
			message.EUI = strings.ToLower(message.EUI)

			// Remove leading v from firmware so it matches zigbee format (semver)
			message.FirmwareVersion = strings.TrimPrefix(message.FirmwareVersion, "v")

			OnStateMessage(message)
		}

	case "sensor/electricitymeter":
		var message ElectricityMessage
		if json.Unmarshal(mqttPacket.Payload, &message) == nil {
			slog.Info("Received electricity meter update", "msg", message)

			// Tidy up fields that probably should have been null in the first place
			if *message.ElectricityMeter.Energy.Import.Supplier == "---" {
				message.ElectricityMeter.Energy.Import.Supplier = nil
			}
			if *message.ElectricityMeter.Energy.Import.MPAN == "read pending" {
				message.ElectricityMeter.Energy.Import.MPAN = nil
			}

			OnElectricityMessage(message)
		}

	case "sensor/gasmeter":
		var message GasMessage
		if json.Unmarshal(mqttPacket.Payload, &message) == nil {
			slog.Info("Received gas meter update", "msg", message)

			// Tidy up fields that probably should have been null in the first place
			if *message.GasMeter.Energy.Import.Supplier == "---" {
				message.GasMeter.Energy.Import.Supplier = nil
			}
			if *message.GasMeter.Energy.Import.MPRN == "read pending" {
				message.GasMeter.Energy.Import.MPRN = nil
			}
			if *message.GasMeter.Energy.Import.CumulativeVolumeUnit == "" {
				message.GasMeter.Energy.Import.CumulativeVolumeUnit = nil
			}
			if *message.GasMeter.Energy.Import.VolumeUnit == "" {
				message.GasMeter.Energy.Import.VolumeUnit = nil
			}

			OnGasMessage(message)
		}

	default:
		slog.Warn("Received MQTT message with unknown topic", "topic", mqttPacket.TopicName, "payload", string(mqttPacket.Payload))
	}
}
