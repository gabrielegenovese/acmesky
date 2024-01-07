package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Value struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}

type Elem struct {
	Distance    Value  `json:"distance"`
	Duration    Value  `json:"duration"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Status      string `json:"status"`
}

type Row struct {
	Elements []Elem `json:"elements"`
}

type GeoDistance struct {
	Dest   []string `json:"destination_addresses"`
	Origin []string `json:"origin_addresses"`
	Rows   []Row    `json:"rows"`
	Status string   `json:"status"`
}

type ResBody struct {
	Distance string `json:"distance"`
	Value    int    `json:"value"`
	Status   string `json:"status"`
}

func sendError(w http.ResponseWriter) {
	elemBody := ResBody{
		Distance: "",
		Value:    0,
		Status:   "ERROR",
	}
	writeRes(w, elemBody)
}

func writeRes(w http.ResponseWriter, el ResBody) {
	w.Header().Set("Content-Type", "application/json")
	bodybyte, err := json.MarshalIndent(el, "", "")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(bodybyte)
}

func format_query(from string, to string) *url.URL {
	url, err := url.Parse("https://api.distancematrix.ai/maps/api/distancematrix/json")
	if err != nil {
		log.Fatal(err)
	}

	values := url.Query()
	values.Add("origins", from)
	values.Add("destinations", to)
	values.Add("key", os.Getenv("GEOKEY"))
	url.RawQuery = values.Encode()

	return url
}

func calcDistance(w http.ResponseWriter, req *http.Request) {
	from_param := req.URL.Query().Get("from")
	to_param := req.URL.Query().Get("to")

	url := format_query(from_param, to_param)

	resp, err := http.Get(url.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result GeoDistance
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			log.Fatal(err)
		}
		if result.Status == "OK" {
			stringDistance := result.Rows[0].Elements[0].Distance.Text
			valueDistance := result.Rows[0].Elements[0].Distance.Value
			statusDistance := result.Rows[0].Elements[0].Status
			elemBody := ResBody{
				Distance: stringDistance,
				Value:    valueDistance,
				Status:   statusDistance,
			}
			writeRes(w, elemBody)
		} else {
			sendError(w)
		}
	} else {
		sendError(w)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/distance", calcDistance)
	log.Println("Listing for requests at http://localhost:8000/distance")
	log.Fatal(http.ListenAndServe(":8000", nil))
}