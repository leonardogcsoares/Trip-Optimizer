package main

import "gitlab.com/trip-optimizer/skyscanner"

func main() {

	api := skyscanner.API{}
	var cpr = skyscanner.CheapestPathRequest{
		Start: skyscanner.CPStart{
			Name: "LON",
			Date: "2017-02-10",
		},
		Places: []skyscanner.CPPlace{
			skyscanner.CPPlace{
				Name: "VIE", // Vienna
				Stay: 2,
			},
			skyscanner.CPPlace{
				Name: "ROME-sky",
				Stay: 4,
			},
			skyscanner.CPPlace{
				Name: "MA-sky",
				Stay: 3,
			},
			skyscanner.CPPlace{
				Name: "AMS-sky",
				Stay: 5,
			},
		},
	}

	_, err := api.GetCheapestPath(cpr)
	if err != nil {
		return
	}

	// _, err := api.GetMonthPriceRoute(
	// 	skyscanner.Location{
	// 		Name: "LON",
	// 		Date: "2017-02",
	// 	},
	// 	skyscanner.Location{
	// 		Name: "JFK",
	// 	},
	// )
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

}
