package main

import (
	encjson "encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sort"
)

var CampaniaPharmacies GeoJson

type SearchArgs struct {
	Limit int `json:"limit"`
	Range int `json:"range"`
	CurrentLocation struct {
		Latitude float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"currentLocation"`
}

type Response struct {
	Result string
}

type SearchService struct {}

type GeoJson struct{
	Features []struct{
		Type string `json:"type"`
		Geometry struct {
			Type string `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		Properties struct {
			Descrizione string `json:"descrizione"`
		} `json:"properties"`
	} `json:"features"`
}

type Pharmacy struct {
	Name string `json:"name"`
	Distance int `json:"distance"`
	Location struct {
		Latitude float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
}

type SearchServiceResponse struct {
	Pharmacies []Pharmacy `json:"pharmacies"`
}

func calcDistance(lat1 float64, lat2 float64, lng1 float64, lng2 float64) float64 {
	const EARTHRADIUS float64 = 6378100

	radlat1 := lat1 * (math.Pi / 180)
	radlat2 := lat2 * (math.Pi / 180)
	radlng1 := lng1 * (math.Pi / 180)
	radlng2 := lng2 * (math.Pi / 180)

	lngDelta := radlng1 - radlng2

	a := math.Pow(math.Cos(radlat2) * math.Sin(lngDelta), 2) +
		math.Pow(math.Cos(radlat1) * math.Sin(radlat2) - math.Sin(radlat1) * math.Cos(radlat2) * math.Cos(lngDelta), 2)
	b := math.Sin(radlat1) * math.Sin(radlat2) + math.Cos(radlat1) * math.Cos(radlat2) * math.Cos(lngDelta)

	angle := math.Atan2(math.Sqrt(a), b)
	return angle * EARTHRADIUS
}

func (t *SearchService) NearestPharmacy(r *http.Request, args *SearchArgs, result *SearchServiceResponse) error {

	// Inizializzo la risposta
	var res SearchServiceResponse

	// Ciclo tutte le farmacie campane
	for _, pharmacy := range CampaniaPharmacies.Features {

		// Calcolo la distanza
		distance := int(calcDistance(args.CurrentLocation.Latitude, pharmacy.Geometry.Coordinates[1], args.CurrentLocation.Longitude, pharmacy.Geometry.Coordinates[0]))

		// Se supera il raggio limite impostato non lo inserisco nella risposta
		if distance <= args.Range {

			// Genero un nuovo oggetto di tipo Pharmacy e lo valorizzo
			var tmp = new(Pharmacy)
			tmp.Name = pharmacy.Properties.Descrizione
			tmp.Distance = distance
			tmp.Location.Latitude = pharmacy.Geometry.Coordinates[1]
			tmp.Location.Longitude = pharmacy.Geometry.Coordinates[0]

			res.Pharmacies = append(res.Pharmacies, *tmp)
		}
	}

	// Ordino per distanza crescente
	sort.Slice(res.Pharmacies, func(i, j int) bool {
		return res.Pharmacies[i].Distance < res.Pharmacies[j].Distance
	})
	limit := args.Limit
	if limit > len(res.Pharmacies) {
		limit = len(res.Pharmacies)
	}
	// Restituisco al massimo un numero Limit di farmacie
	res.Pharmacies = res.Pharmacies[:limit]

	*result = res
	return nil
}

// Recupero il JSON delle farmacie e lo inserisco all'interno di CampaniaPharmacies
func downloadCampaniaPharmacies() error {
	resp, err := http.Get("https://dati.regione.campania.it/catalogo/resources/Elenco-Farmacie.geojson")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		encjson.Unmarshal(bodyBytes, &CampaniaPharmacies)
		fmt.Println("Got pharmacies")
	} else {
		fmt.Println("Error retrieving pharmacies")
	}
	return nil
}

func main() {
	downloadCampaniaPharmacies()
	rpcServer := rpc.NewServer()

	rpcServer.RegisterCodec(json.NewCodec(), "application/json")
	rpcServer.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

	search := new(SearchService)

	rpcServer.RegisterService(search, "Search")

	router := mux.NewRouter()
	router.Handle("/rpc", rpcServer)
	http.ListenAndServe(":8081", router)
}