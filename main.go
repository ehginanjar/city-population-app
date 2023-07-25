package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
)

type CitySource struct {
	Name       string `json:"city"`
	Population int    `json:"population"`
}

type City struct {
	Source CitySource `json:"_source"`
}

var es *elasticsearch.Client

func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Health check endpoint hit")
	healthStatus := "OK"
	statusCode := http.StatusOK

	// Check the health of Elasticsearch
	_, err := es.Ping()
	if err != nil {
		log.Println("Error pinging Elasticsearch:", err)
		healthStatus = "Elasticsearch not reachable"
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, `{"status": "%s"}`, healthStatus)
}

func insertOrUpdateCityHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Insert or Update City endpoint hit")
	var city City
	err := json.NewDecoder(r.Body).Decode(&city)
	if err != nil {
		log.Println("Error decoding request payload:", err)
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	if city.Source.Name == "" || city.Source.Population == 0 {
		log.Println("City and population must be provided")
		http.Error(w, `{"error": "City and population must be provided"}`, http.StatusBadRequest)
		return
	}

	err = indexCity(city)
	if err != nil {
		log.Println("Error indexing city:", err)
		http.Error(w, `{"error": "Failed to add/update city"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "%s added/updated successfully"}`, city.Source.Name)
}

func getCityPopulationHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Retrieve Population of a City endpoint hit")
	cityName := strings.TrimPrefix(r.URL.Path, "/city/")
	if cityName == "" {
		log.Println("City name must be provided")
		http.Error(w, `{"error": "City name must be provided"}`, http.StatusBadRequest)
		return
	}

	city, err := getCity(cityName)
	if err != nil {
		log.Println("Error retrieving city:", err)
		http.Error(w, `{"error": "City not found"}`, http.StatusNotFound)
		return
	}

	// Print city details in the server logs
	log.Printf("City: %s, Population: %d\n", city.Source.Name, city.Source.Population)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(city)
}

func indexCity(city City) error {
	log.Println("Indexing city:", city.Source.Name)
	body := fmt.Sprintf(`{"city": "%s", "population": %d}`, strings.ReplaceAll(city.Source.Name, `"`, `\"`), city.Source.Population)
	_, err := es.Index("cities", strings.NewReader(body), es.Index.WithDocumentID(city.Source.Name))
	return err
}

func getCity(cityName string) (City, error) {
	log.Println("Retrieving city:", cityName)
	res, err := es.Get("cities", cityName)
	if err != nil {
		log.Println("Error retrieving city:", err)
		return City{}, err
	}

	if res.IsError() {
		log.Println("Error response:", res)
		return City{}, fmt.Errorf("failed to retrieve city: %s", res.String())
	}

	var city City
	if err := json.NewDecoder(res.Body).Decode(&city); err != nil {
		log.Println("Error decoding city:", err)
		return City{}, err
	}

	return city, nil
}

func main() {
	log.Println("Starting the application...")
	cfg := elasticsearch.Config{
		Addresses: []string{"http://elasticsearch:9200"},
	}
	var err error
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/city", insertOrUpdateCityHandler)
	http.HandleFunc("/city/", getCityPopulationHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
