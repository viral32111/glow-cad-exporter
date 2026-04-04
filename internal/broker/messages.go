package broker

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/viral32111/glow-cad-exporter/internal/metrics"
)

type StateMessage struct {
	Timestamp string `json:"timestamp"`

	EthernetMAC string `json:"ethmac"`
	WirelessMAC string `json:"wifistamac"`
	EUI         string `json:"eui"`

	SMETSVersion string `json:"smetsversion"`

	HardwareRevision string `json:"hardware"`
	FirmwareVersion  string `json:"software"`
	ZigbeeVersion    string `json:"zigbee"`

	HAN struct {
		Status string `json:"status"`
		RSSI   int    `json:"rssi"`
		LQI    int    `json:"lqi"`
	} `json:"han"`
}

type MeterPrice struct {
	UnitRate       *float64 `json:"unitrate"`
	StandingCharge *float64 `json:"standingcharge"`
}

type ElectricityMessage struct {
	ElectricityMeter struct {
		Timestamp string `json:"timestamp"`

		Energy struct {
			Import struct {
				Supplier *string `json:"supplier"`
				MPAN     *string `json:"mpan"`

				Cumulative *float64 `json:"cumulative"`
				Day        *float64 `json:"day"`
				Week       *float64 `json:"week"`
				Month      *float64 `json:"month"`
				Unit       string   `json:"units"`

				Price MeterPrice `json:"price"`
			} `json:"import"`

			Export struct {
				Cumulative *float64 `json:"cumulative"`
				Units      string   `json:"units"`
			} `json:"export"`
		} `json:"energy"`

		Power struct {
			Value *float64 `json:"value"`
			Unit  string   `json:"units"`
		} `json:"power"`
	} `json:"electricitymeter"`
}

type GasMessage struct {
	GasMeter struct {
		Timestamp string `json:"timestamp"`

		Energy struct {
			Import struct {
				Supplier *string `json:"supplier"`
				MPRN     *string `json:"mprn"`

				Cumulative           *float64 `json:"cumulative"`
				CumulativeVolume     *float64 `json:"cumulativevol"`
				CumulativeVolumeUnit *string  `json:"cumulativevolunits"`

				Day         *float64 `json:"day"`
				DayVolume   *float64 `json:"dayvol"`
				Week        *float64 `json:"week"`
				WeekVolume  *float64 `json:"weekvol"`
				Month       *float64 `json:"month"`
				MonthVolume *float64 `json:"monthvol"`
				Unit        string   `json:"units"`
				VolumeUnit  *string  `json:"dayweekmonthvolunits"`

				Price MeterPrice `json:"price"`
			} `json:"import"`
		} `json:"energy"`
	} `json:"gasmeter"`
}

func parseISO8601(iso8601 string) (int64, error) {
	_time, err := time.Parse(time.RFC3339, iso8601)
	if err != nil {
		return 0, err
	}

	return _time.Unix(), nil
}

func OnStateMessage(message StateMessage) {
	labels := prometheus.Labels{
		"ethernet": message.EthernetMAC,
		"wifi":     message.WirelessMAC,
		"eui":      message.EUI,
		"hardware": message.HardwareRevision,
		"firmware": message.FirmwareVersion,
		"zigbee":   message.ZigbeeVersion,
	}

	switch message.HAN.Status {
	case "joined":
		metrics.Status.With(labels).Set(1)
	case "rejoining":
		metrics.Status.With(labels).Set(2)
	case "not joined":
		metrics.Status.With(labels).Set(3)
	default:
		metrics.Status.With(labels).Set(-1)
	}

	metrics.RSSI.With(labels).Set(float64(message.HAN.RSSI))
	metrics.LQI.With(labels).Set(float64(message.HAN.LQI))

	// Convert ISO-8601 to Unix timestamp (seconds)
	if timestamp, err := parseISO8601(message.Timestamp); err == nil {
		metrics.Timestamp.Set(float64(timestamp), labels)
	}
}

func OnElectricityMessage(message ElectricityMessage) {
	labels := prometheus.Labels{"supplier": "", "mpan": ""}
	if message.ElectricityMeter.Energy.Import.Supplier != nil {
		labels["supplier"] = *message.ElectricityMeter.Energy.Import.Supplier
	}
	if message.ElectricityMeter.Energy.Import.MPAN != nil {
		labels["mpan"] = *message.ElectricityMeter.Energy.Import.MPAN
	}

	if message.ElectricityMeter.Energy.Import.Cumulative != nil {
		metrics.ElectricImportCumulative.Set(*message.ElectricityMeter.Energy.Import.Cumulative, labels)
	}
	if message.ElectricityMeter.Energy.Import.Day != nil {
		metrics.ElectricImportDay.With(labels).Set(*message.ElectricityMeter.Energy.Import.Day)
	}
	if message.ElectricityMeter.Energy.Import.Week != nil {
		metrics.ElectricImportWeek.With(labels).Set(*message.ElectricityMeter.Energy.Import.Week)
	}
	if message.ElectricityMeter.Energy.Import.Month != nil {
		metrics.ElectricImportMonth.With(labels).Set(*message.ElectricityMeter.Energy.Import.Month)
	}

	if message.ElectricityMeter.Energy.Import.Price.StandingCharge != nil {
		metrics.ElectricTariffStandingCharge.With(labels).Set(*message.ElectricityMeter.Energy.Import.Price.StandingCharge)
	}
	if message.ElectricityMeter.Energy.Import.Price.UnitRate != nil {
		metrics.ElectricTariffUnitRate.With(labels).Set(*message.ElectricityMeter.Energy.Import.Price.UnitRate)
	}

	if message.ElectricityMeter.Energy.Export.Cumulative != nil {
		metrics.ElectricExportCumulative.Set(*message.ElectricityMeter.Energy.Export.Cumulative, labels)
	}

	if message.ElectricityMeter.Power.Value != nil {
		metrics.ElectricPower.With(labels).Set(*message.ElectricityMeter.Power.Value)
	}

	// Convert ISO-8601 to Unix timestamp (seconds)
	if timestamp, err := parseISO8601(message.ElectricityMeter.Timestamp); err == nil {
		metrics.ElectricTimestamp.Set(float64(timestamp), labels)
	}
}

func OnGasMessage(message GasMessage) {
	labels := prometheus.Labels{"supplier": "", "mprn": ""}
	if message.GasMeter.Energy.Import.Supplier != nil {
		labels["supplier"] = *message.GasMeter.Energy.Import.Supplier
	}
	if message.GasMeter.Energy.Import.MPRN != nil {
		labels["mprn"] = *message.GasMeter.Energy.Import.MPRN
	}

	if message.GasMeter.Energy.Import.Cumulative != nil {
		metrics.GasCumulative.Set(*message.GasMeter.Energy.Import.Cumulative, labels)
	}
	if message.GasMeter.Energy.Import.CumulativeVolume != nil {
		metrics.GasCumulativeVolume.Set(*message.GasMeter.Energy.Import.CumulativeVolume, labels)
	}
	if message.GasMeter.Energy.Import.Day != nil {
		metrics.GasDay.With(labels).Set(*message.GasMeter.Energy.Import.Day)
	}
	if message.GasMeter.Energy.Import.DayVolume != nil {
		metrics.GasDayVolume.With(labels).Set(*message.GasMeter.Energy.Import.DayVolume)
	}
	if message.GasMeter.Energy.Import.Week != nil {
		metrics.GasWeek.With(labels).Set(*message.GasMeter.Energy.Import.Week)
	}
	if message.GasMeter.Energy.Import.WeekVolume != nil {
		metrics.GasWeekVolume.With(labels).Set(*message.GasMeter.Energy.Import.WeekVolume)
	}
	if message.GasMeter.Energy.Import.Month != nil {
		metrics.GasMonth.With(labels).Set(*message.GasMeter.Energy.Import.Month)
	}
	if message.GasMeter.Energy.Import.MonthVolume != nil {
		metrics.GasMonthVolume.With(labels).Set(*message.GasMeter.Energy.Import.MonthVolume)
	}

	if message.GasMeter.Energy.Import.Price.StandingCharge != nil {
		metrics.GasTariffStandingCharge.With(labels).Set(*message.GasMeter.Energy.Import.Price.StandingCharge)
	}
	if message.GasMeter.Energy.Import.Price.UnitRate != nil {
		metrics.GasTariffUnitRate.With(labels).Set(*message.GasMeter.Energy.Import.Price.UnitRate)
	}

	// Convert ISO-8601 to Unix timestamp (seconds)
	if timestamp, err := parseISO8601(message.GasMeter.Timestamp); err == nil {
		metrics.GasTimestamp.Set(float64(timestamp), labels)
	}
}
