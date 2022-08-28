package component

const (
	InsufficientFuelRegex = "Ship has insufficient fuel for flight plan. You require ([0-9]+) more FUEL"
)

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// {"error":{"message":"Ship has insufficient fuel for flight plan. You require 13 more FUEL","code":3001}}
