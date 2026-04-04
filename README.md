# ⚡ Glow CAD Exporter

[![CI](https://github.com/viral32111/glow-cad-exporter/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/viral32111/glow-cad-exporter/actions/workflows/ci.yml)
[![CodeQL](https://github.com/viral32111/glow-cad-exporter/actions/workflows/codeql.yml/badge.svg)](https://github.com/viral32111/glow-cad-exporter/actions/workflows/codeql.yml)
![GitHub tag (with filter)](https://img.shields.io/github/v/tag/viral32111/glow-cad-exporter?label=Latest)
![GitHub repository size](https://img.shields.io/github/repo-size/viral32111/glow-cad-exporter?label=Size)
![GitHub release downloads](https://img.shields.io/github/downloads/viral32111/glow-cad-exporter/total?label=Downloads)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/viral32111/glow-cad-exporter?label=Commits)

This is a Prometheus metrics exporter for the [Hidlebrand Glow Wired CAD](https://shop.glowmarkt.com/products/glow-cad-for-smart-meters).

The Glow CAD *(Consumer Access Device)* is a small device that talks to your smart meter's Communications Hub over its Zigbee HAN *(Home Area Network)*.
It functions the same as a supplier IHD (In-Home Display) to fetch & display your real-time electricity & gas consumption, but instead of a screen *(assuming you're using a [Glow Wired CAD](https://shop.glowmarkt.com/products/glow-cad-for-smart-meters), they do offer [a variant with a display](https://shop.glowmarkt.com/products/display-and-cad-combined-for-smart-meter-customers) to replace your IHD)* it exposes the data over their [Bright Cloud API](https://bitbucket.org/ijosh/brightglowmarkt/src/master/) and/or a local MQTT publisher.

This project takes advantage of the local MQTT functionality to receive the data directly over your local network without needing to interact with their cloud API at all, thus you can safely cut off Internet access to the Glow CAD without losing any functionality.
No MQTT broker (e.g., Mosquitto) is necessary since this project acts as a stub broker to receive the MQTT messages directly from the Glow CAD and expose the relevant data on a Prometheus scrapeable metrics HTTP endpoint.

## 📊 Metrics

Electricity readings are updated every 6-10 seconds, while gas readings are updated every 30 minutes.
This is a fundamental limitation of most smart meter installations (gas meters are battery powered and thus need to have low energy usage) and not the fault of the Glow CAD or this project.**

### 📡 Device Status

* `glow_status`: Zigbee connection status (`0`: connected, `1`: reconnecting, `2`: not connected).
* `glow_rssi`: Zigbee signal strength in dBm to the smart meter's Communications Hub.
* `glow_lqi`: Zigbee signal quality score from 0 to 255.
* `glow_timestamp`: Unix timestamp of the last received state update from the Glow CAD.

#### Labels

* `ethernet`: MAC address of the Glow CAD's Ethernet interface.
* `wifi`: MAC address of the Glow CAD's Wi-Fi interface.
* `eui`: EUI-64 of the Glow CAD's Zigbee radio.
* `hardware`: Hardware revision of the Glow CAD.
* `firmware`: Firmware version of the Glow CAD.
* `zigbee`: Zigbee stack version of the Glow CAD.

<details>
<summary>Click to view the raw JSON message published by the Glow CAD</summary>

**MQTT Topic**: `glow/000000000000/STATE` *(the zeros are the MAC address of the Ethernet or Wi-Fi interface, depending on which one the Glow CAD is using)*

```jsonc
{
	"ethmac": "000000000000", // MAC address of the Glow CAD's Ethernet port
	"wifistamac": "000000000000", // MAC address of the Glow CAD's Wi-Fi interface
	"eui": "70:B3:D5:00:00:00:00:00", // EUI-64 of the Glow CAD's Zigbee radio
	"han": {
		"lqi": 128, // Zigbee link quality indicator (quality score from 0-255)
		"rssi": -68, // Zigbee signal strength in dBm (negative value, closer to 0 is stronger)
		"status": "joined" // Zigbee network status (usually: 'joined', 'rejoining', 'not joined')
	},
	"hardware": "GLOW-ESP-CAD-0v000-CAD-SMETS2", // Glow CAD's hardware revision
	"smetsversion": "SMETS2", // Meter specification type (SMETS1 or SMETS2)
	"software": "v0.0.0", // Glow CAD's firmware version
	"zigbee": "0.0.0", // Glow CAD's Zigbee stack version
	"timestamp": "2026-04-04T09:16:20Z", // ISO-8601 timestamp of when this message was generated
}
```
</details>

### 🔌 Electricity Meter

* `glow_electric_import_cumulative`: Electricity consumption in kWh since the smart meter was commissioned (a.k.a. meter reading).
* `glow_electric_import_day`: Electricity consumption in kWh for the last 24 hours.
* `glow_electric_import_week`: Electricity consumption in kWh for the last 7 days.
* `glow_electric_import_month`: Electricity consumption in kWh for the last 30 days.
* `glow_electric_export_cumulative`: Electricity generation in kWh since the smart meter was commissioned. Only applicable if you sell electricity back to the grid.
* `glow_electric_tariff_standing_charge`: Your supplier's standing charge in GBP/day.
* `glow_electric_tariff_unit_rate`: Your cost per kWh in GBP.
* `glow_electric_power`: Current power consumption (positive value) or generation (negative value) in kW.
* `glow_electric_timestamp`: Unix timestamp of the last received electricity meter update.

#### Labels

* `supplier`: Name of your electricity supplier.
* `mpan`: Meter Point Administration Number of your electricity meter.

<details>
<summary>Click to view the raw JSON message published by the Glow CAD</summary>

**MQTT Topic**: `glow/000000000000/SENSOR/electricitymeter` *(the zeros are the MAC address of the Ethernet or Wi-Fi interface, depending on which one the Glow CAD is using)*

```jsonc
{
	"electricitymeter": {
		"energy": {
			"import": {
				"supplier": "xxxxxxxxxx", // Name of your electricity supplier (e.g., OVO Energy, British Gas, etc.)
				"mpan": "0000000000000", // Meter Point Administration Number of your electricity meter
				"cumulative": 10000.000, // Total imported (consumed) electricity since the meter was commissioned (a.k.a. meter reading)
				"day": 5.000, // Imported (consumed) electricity in the last 24 hours
				"week": 50.000, // Imported (consumed) electricity in the last 7 days
				"month": 100.000, // Imported (consumed) electricity in the last 30 days
				"units": "kWh", // Units for the above energy readings (always kWh)
				"price": {
					"standingcharge": 0.57900, // Your supplier's daily standing charge cost
					"unitrate": 0.24960 // Your supplier's cost per kWh
				},
			},
			"export": {
				"cumulative": 0.001, // Total exported (generated) electricity since the meter was commissioned
				"units": "kWh" // Units for the above energy reading (always kWh)
			},
		},
		"power": {
			"value": 0.476, // Current instantaneous power consumption (positive value) or generation (negative value)
			"units": "kW" // Units for the above power reading (always kW)
		},
		"timestamp": "2026-04-04T09:16:00Z" // ISO-8601 timestamp of when this message was generated
	}
}
```
</details>

### 🌡️ Gas Meter

* `glow_gas_cumulative`: Gas consumption in kWh since the smart meter was commissioned (a.k.a. meter reading).
* `glow_gas_cumulative_volume`: Gas consumption in cubic meters since the smart meter was commissioned (a.k.a. meter reading).
* `glow_gas_day`: Gas consumption in kWh for the last 24 hours.
* `glow_gas_day_volume`: Gas consumption in cubic meters for the last 24 hours.
* `glow_gas_week`: Gas consumption in kWh for the last 7 days.
* `glow_gas_week_volume`: Gas consumption in cubic meters for the last 7 days.
* `glow_gas_month`: Gas consumption in kWh for the last 30 days.
* `glow_gas_month_volume`: Gas consumption in cubic meters for the last 30 days.
* `glow_gas_tariff_standing_charge`: Your supplier's standing charge in GBP/day.
* `glow_gas_tariff_unit_rate`: Your cost per kWh in GBP.
* `glow_gas_timestamp`: Unix timestamp of the last received gas meter update.

#### Labels

* `supplier`: Name of your gas supplier.
* `mprn`: Meter Point Reference Number of your gas meter.

<details>
<summary>Click to view the raw JSON message published by the Glow CAD</summary>

**MQTT Topic**: `glow/000000000000/SENSOR/gasmeter` *(the zeros are the MAC address of the Ethernet or Wi-Fi interface, depending on which one the Glow CAD is using)*

```jsonc
{
	"gasmeter": {
		"energy": {
			"import": {
				"supplier": "---",
				"mprn": "read pending",
				"cumulative": null,
				"cumulativevol": null,
				"cumulativevolunits": "",
				"day": null,
				"dayvol": null,
				"week": null,
				"weekvol": null,
				"month": null,
				"monthvol": null,
				"dayweekmonthvolunits": "",
				"units": "kWh",
				"price": {
					"standingcharge": null,
					"unitrate": null
				},
			}
		},
		"timestamp": "1970-01-01T00:00:00Z"
	}
}
```
</details>

## ⚙️ Configuration

All configuration is via environment variables for easy deployment with Docker. See [`.env.example`](/.env.example).

| Name | Type | Description | Default |
| --- | --- | --- | --- |
| `GLOW_LOG_LEVEL` | String | Minimum level of logs to print (`debug`, `info`, `warn`, `error`). | `info` |
| `GLOW_HTTP_ADDRESS` | String | IP address to listen for HTTP requests on. | `127.0.0.1` |
| `GLOW_HTTP_PORT` | Number | TCP port to listen for HTTP requests on. | `8080` |
| `GLOW_METRICS_PATH` | String | HTTP path to expose Prometheus metrics on. | `/metrics` |
| `GLOW_METRICS_GO_BUILTIN` | Boolean | Whether to include Go runtime metrics. | `false` |
| `GLOW_MQTT_ADDRESS` | String | IP address to listen for MQTT messages on. | `127.0.0.1` |
| `GLOW_MQTT_PORT` | Number | TCP port to listen for MQTT messages on. | `1883` |
| `GLOW_MQTT_USER` | String | Username for authenticating MQTT connections. | |
| `GLOW_MQTT_PASSWORD` | String | Password for authenticating MQTT connections. | |
| `GLOW_MQTT_TOPIC` | String | Prefix for topics published by the Glow CAD. | `glow` |
| `GLOW_REDIS_URL` | String | URL to the Redis server for caching MQTT messages. | `redis://default@127.0.0.1:6379/0` |
| `GLOW_REDIS_KEY_PREFIX` | String | Prefix for Redis keys to avoid collisions. | `glow` |
| `GLOW_REDIS_STATE_TTL` | Duration | How long to cache device state MQTT messages in Redis. | `24h` |
| `GLOW_REDIS_ELECTRIC_TTL` | Duration | How long to cache electricity meter MQTT messages in Redis. | `730h` |
| `GLOW_REDIS_GAS_TTL` | Duration | How long to cache gas meter MQTT messages in Redis. | `730h` |

## 🧰 Development

1. Copy the [`.env.example`](/.env.example) file to `.env` and populate accordingly.
2. Run the application with `go run ./cmd/glow-cad-exporter/main.go`.

## 🛠️ Building

Run the relevant build command below, assuming you're on a Linux 64-bit (glibc) system.

### Linux 64-bit (glibc)

```bash
GOOS=linux GOARCH=amd64 \
	go build -o glow-cad-exporter-v0.0.0-linux-amd64-glibc ./cmd/glow-cad-exporter/main.go
```

### Linux 64-bit (musl)

```bash
CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
	go build -o glow-cad-exporter-v0.0.0-linux-amd64-musl ./cmd/glow-cad-exporter/main.go
```

### Windows 64-bit

```bash
GOOS=windows GOARCH=amd64 \
	go build -o glow-cad-exporter-v0.0.0-windows-amd64.exe ./cmd/glow-cad-exporter/main.go
```

## ⚖️ License

Copyright (C) 2026 [viral32111](https://viral32111.com).

> This program is free software: you can redistribute it and/or modify
> it under the terms of the GNU Affero General Public License as
> published by the Free Software Foundation, either version 3 of the
> License, or (at your option) any later version.

> This program is distributed in the hope that it will be useful,
> but WITHOUT ANY WARRANTY; without even the implied warranty of
> MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
> GNU Affero General Public License for more details.

> You should have received a copy of the GNU Affero General Public License
> along with this program. If not, see https://www.gnu.org/licenses.
