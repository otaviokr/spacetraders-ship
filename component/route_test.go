package component_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/otaviokr/spacetraders-ship/component"
)

func TestReadRouteDescription(t *testing.T) {
	useCases := map[string]map[string]interface{}{
		"uc1": {
			"yaml": "route:\n  - station: myTest\n    sell:\n      good01: 135\n    buy:\n      material01: 123\n      material02: 456",
			"expected": &component.Route{
				Route: []component.RouteStop{
					{
						Station: "myTest",
						Sell: map[string]int{
							"good01": 135,
						},
						Buy: map[string]int{
							"material01": 123,
							"material02": 456,
						}}}}}}

	for _, uc := range useCases {
		actual, err := component.ReadRouteDescription(strings.NewReader(fmt.Sprintf("%v", uc["yaml"])))
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, uc["expected"]) {
			t.Fatalf("\n%+v\n\n%+v\n", actual, uc["expected"])
		}
	}
}
