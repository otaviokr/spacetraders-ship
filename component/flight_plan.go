package component

type FlightPlan struct {
	Details FlightPlanDetails `yaml:"flightPlan"`
	Error   Error             `yaml:"error"`
}

type FlightPlanDetails struct {
	ArrivesAt              string `yaml:"arrivesAt"`
	CreatedAt              string `yaml:"createdAt"`
	Departure              string `yaml:"departure"`
	Destination            string `yaml:"destination"`
	Distance               int    `yaml:"distance"`
	FuelConsumed           int    `yaml:"fuelConsumed"`
	FuelRemaining          int    `yaml:"fuelRemaining"`
	Id                     string `yaml:"id"`
	ShipId                 string `yaml:"shipId"`
	TerminatedAt           string `yaml:"terminatedAt"`
	TimeRemainingInSeconds int    `yaml:"timeRemainingInSeconds"`
}

// {
// 	flightPlan: {
// 	  arrivesAt: '2021-03-28T23:11:50.068Z',
// 	  createdAt: '2021-03-28T23:05:05.078Z',
// 	  departure: 'OE-PM-TR',
// 	  destination: 'OE-CR',
// 	  distance: 46,
// 	  fuelConsumed: 13,
// 	  fuelRemaining: 67,
// 	  id: 'ckmtrssae0402zgopsy39y9z7',
// 	  shipId: 'ckmtrlqpz0109zgopkeqs4m5s',
// 	  terminatedAt: null,
// 	  timeRemainingInSeconds: 365
// 	}
//   }
