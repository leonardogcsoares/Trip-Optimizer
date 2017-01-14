package skyscanner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	apiKey   = "le667217464253823367389433542403"
	currency = "GBP"
	// Country code
	market = "GB"
	// Locale
	locale = "en-GB"
)

// Interface TODO
type Interface interface {
}

// API TODO
type API struct {
	sync.RWMutex
	prices map[string]RouteGridPrice
	start  CPStart
}

// NewAPI TODO
func NewAPI() *API {
	return &API{
		prices: make(map[string]RouteGridPrice),
	}
}

// GetRoutePrice TODO
// Unnecessary for now
// func (a *API) GetRoutePrice(startLoc, endLoc Location) (Route, error) {
//
// 	// http://partners.api.skyscanner.net/apiservices/browsequotes/v1.0/GB/GBP/en-GB/
// 	// LON/JFK/2017-02-10?apiKey=prtl6749387986743898559646983194
//
// 	// reqStr := fmt.Sprintf("http://partners.api.skyscanner.net/apiservices/browsequotes/v1.0/GB/GBP/en-GB/"+
// 	// 	"%s/%s/%s?apiKey=%s", startLoc.Name, endLoc.Name, startLoc.Date.String(), apiKey,
// 	// )
//
// 	return Route{}, nil
// }

// GetCheapestPath TODO
func (a *API) GetCheapestPath(cpr CheapestPathRequest) (FlightPath, error) {

	st, err := time.Parse("2006-01-15", cpr.Start.Date)
	if err != nil {
		return FlightPath{}, err
	}

	toVisit := []CPPlace{}
	for _, p := range cpr.Places {
		toVisit = append(toVisit, p)
	}

	for _, p := range cpr.Places {
		quoteGrid, err := a.getMonthPriceRoute(Location{
			Name: cpr.Start.Name,
			Date: fmt.Sprintf("%d-%s", st.Year(), st.Month().String()),
		},
			Location{
				Name: p.Name,
			})
		if err != nil {
			return FlightPath{}, err
		}

		a.prices[fmt.Sprintf("%s-%s-%d-%s", cpr.Start.Name, p.Name, st.Year(), st.Month().String())] = quoteGrid
	}

	a.start = cpr.Start
	a.calculatePrice(cpr.Start.Name, 0, CPPlace{Name: cpr.Start.Name, startDate: st}, make(map[string]CPPlace), toVisit)

	return FlightPath{}, nil
}

func (a *API) calculatePrice(currentPath string, currentPrice int, origin CPPlace, visited map[string]CPPlace, placesToVisit []CPPlace) {

	haveVisited := make(map[string]CPPlace)
	for k, v := range visited {
		haveVisited[k] = v
	}

	if len(placesToVisit) == 0 {
		// check price from origin to a.start
		// submit final price to prices channel <-finalPrice
		return
	}

	for _, destination := range placesToVisit {
		toVisit := []CPPlace{}
		index := 0
		for i, v := range placesToVisit {
			if v.Name == destination.Name {
				index = i
				continue
			}

			toVisit = append(toVisit, v)
		}

		// check chepeast price from origin to destination given origin.startDate in a.prices map
		price := 0

		toVisit = remove(toVisit, index)
		haveVisited[destination.Name] = destination
		for _, newOrigin := range toVisit {
			a.calculatePrice(fmt.Sprintf("%s-%s", currentPath, destination.Name),
				currentPrice+price,
				newOrigin,
				haveVisited,
				toVisit,
			)
		}
	}

}

func remove(slice []CPPlace, s int) []CPPlace {
	return append(slice[:s], slice[s+1:]...)
}

// GetMonthPriceRoute TODO
func (a *API) getMonthPriceRoute(startLoc, endLoc Location) (RouteGridPrice, error) {
	// http://partners.api.skyscanner.net/apiservices/browsegrid/v1.0/{market}/{currency}/{locale}/
	// {originPlace}/{destinationPlace}/{outboundPartialDate}/{inboundPartialDate}?apiKey={apiKey}

	var rgp RouteGridPrice
	// http: //partners.api.skyscanner.net/apiservices/browsequotes/v1.0/GB/GBP/en-GB/LON/JFK/2017-02?apiKey=prtl6749387986743898559646983194
	url := fmt.Sprintf("http://partners.api.skyscanner.net/apiservices/browsequotes/v1.0/GB/GBP/en-GB/"+
		"%s/%s/%s?apiKey=%s",
		startLoc.Name, endLoc.Name, startLoc.Date, apiKey,
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return rgp, err
	}
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return rgp, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&rgp)
	if err != nil {
		return rgp, err
	}

	return rgp, nil
}
