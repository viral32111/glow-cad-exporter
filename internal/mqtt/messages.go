package mqtt

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
				SupplierName string `json:"supplier"`
				MPAN         string `json:"mpan"`

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
				Supplier string `json:"supplier"`
				MPRN     string `json:"mprn"`

				Cumulative           *float64 `json:"cumulative"`
				CumulativeVolume     *float64 `json:"cumulativevol"`
				CumulativeVolumeUnit string   `json:"cumulativevolunits"`

				Day        *float64 `json:"day"`
				DayVol     *float64 `json:"dayvol"`
				Week       *float64 `json:"week"`
				WeekVol    *float64 `json:"weekvol"`
				Month      *float64 `json:"month"`
				MonthVol   *float64 `json:"monthvol"`
				Unit       string   `json:"units"`
				VolumeUnit string   `json:"dayweekmonthvolunits"`

				Price MeterPrice `json:"price"`
			} `json:"import"`
		} `json:"energy"`
	} `json:"gasmeter"`
}
