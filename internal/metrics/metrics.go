package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	METRICS_NAMESPACE = "glow"
)

var (
	STATE_LABELS    = []string{"ethernet", "wifi", "eui", "hardware", "firmware", "zigbee"}
	ELECTRIC_LABELS = []string{"supplier", "mpan"}
	GAS_LABELS      = []string{"supplier", "mprn"}
)

var (
	Status    = NewGauge("status", "Zigbee connection status with the smart meter communications hub", STATE_LABELS)
	RSSI      = NewGauge("rssi", "Zigbee signal strength to the smart meter communications hub", STATE_LABELS)
	LQI       = NewGauge("lqi", "Zigbee link quality score to the smart meter communications hub", STATE_LABELS)
	Timestamp = NewMonotonicGauge("timestamp", "When the last state update was generated", STATE_LABELS)

	ElectricImportCumulative     = NewMonotonicGauge("electric_import_cumulative", "Cumulative electricity imported (in kWh) from the grid", ELECTRIC_LABELS)
	ElectricImportDay            = NewGauge("electric_import_day", "Electricity imported (in kWh) from the grid in the last day", ELECTRIC_LABELS)
	ElectricImportWeek           = NewGauge("electric_import_week", "Electricity imported (in kWh) from the grid in the last week", ELECTRIC_LABELS)
	ElectricImportMonth          = NewGauge("electric_import_month", "Electricity imported (in kWh) from the grid in the last month", ELECTRIC_LABELS)
	ElectricExportCumulative     = NewMonotonicGauge("electric_export_cumulative", "Cumulative electricity exported (in kWh) to the grid", ELECTRIC_LABELS)
	ElectricPower                = NewGauge("electric_power", "Current electricity power usage (in kW)", ELECTRIC_LABELS)
	ElectricTariffStandingCharge = NewGauge("electric_tariff_standing_charge", "Supplier standing charge (in GBP/day)", ELECTRIC_LABELS)
	ElectricTariffUnitRate       = NewGauge("electric_tariff_unit_rate", "Cost (in GBP/kWh)", ELECTRIC_LABELS)
	ElectricTimestamp            = NewMonotonicGauge("electric_timestamp", "When the last electricity meter update was generated", ELECTRIC_LABELS)

	GasCumulative           = NewMonotonicGauge("gas_cumulative", "Cumulative gas consumption (in kWh)", GAS_LABELS)
	GasCumulativeVolume     = NewMonotonicGauge("gas_cumulative_volume", "Cumulative gas consumption (in m3)", GAS_LABELS)
	GasDay                  = NewGauge("gas_day", "Gas consumption (in kWh) in the last day", GAS_LABELS)
	GasDayVolume            = NewGauge("gas_day_volume", "Gas consumption (in m3) in the last day", GAS_LABELS)
	GasWeek                 = NewGauge("gas_week", "Gas consumption (in kWh) in the last week", GAS_LABELS)
	GasWeekVolume           = NewGauge("gas_week_volume", "Gas consumption (in m3) in the last week", GAS_LABELS)
	GasMonth                = NewGauge("gas_month", "Gas consumption (in kWh) in the last month", GAS_LABELS)
	GasMonthVolume          = NewGauge("gas_month_volume", "Gas consumption (in m3) in the last month", GAS_LABELS)
	GasTariffStandingCharge = NewGauge("gas_tariff_standing_charge", "Supplier standing charge (in GBP/day)", GAS_LABELS)
	GasTariffUnitRate       = NewGauge("gas_tariff_unit_rate", "Cost (in GBP/kWh)", GAS_LABELS)
	GasTimestamp            = NewMonotonicGauge("gas_timestamp", "When the last gas meter update was generated", GAS_LABELS)
)

func NewGauge(name string, help string, labels []string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: fmt.Sprintf("%s_%s", METRICS_NAMESPACE, name),
		Help: help,
	}, labels)
}

func Register(registry *prometheus.Registry) {
	collectors := []prometheus.Collector{
		Status,
		RSSI,
		LQI,
		Timestamp.Collector(),

		ElectricImportCumulative.Collector(),
		ElectricImportDay,
		ElectricImportWeek,
		ElectricImportMonth,
		ElectricExportCumulative.Collector(),
		ElectricPower,
		ElectricTariffStandingCharge,
		ElectricTariffUnitRate,
		ElectricTimestamp.Collector(),

		GasCumulative.Collector(),
		GasCumulativeVolume.Collector(),
		GasDay,
		GasDayVolume,
		GasWeek,
		GasWeekVolume,
		GasMonth,
		GasMonthVolume,
		GasTariffStandingCharge,
		GasTariffUnitRate,
		GasTimestamp.Collector(),
	}

	for _, collector := range collectors {
		registry.MustRegister(collector)
	}
}
