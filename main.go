package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
)

// корисні копалини
type Minerals struct {
	Coal  float64 `json:"coal"`
	Mazut float64 `json:"mazut"`
	Gas   float64 `json:"gas"`
}

// результати розрахунків
type CalculationResults struct {
	CoalEmissionFactor  float64 `json:"coalEmissionFactor"`
	CoalEmissionValue   float64 `json:"coalEmissionValue"`
	MazutEmissionFactor float64 `json:"mazutEmissionFactor"`
	MazutEmissionValue  float64 `json:"mazutEmissionValue"`
	GasEmissionFactor   float64 `json:"gasEmissionFactor"`
	GasEmissionValue    float64 `json:"gasEmissionValue"`
}

func calculateResults(minerals Minerals) CalculationResults {
	// нижча теплота згоряння робочої маси вугілля
	coalHeatValue := 20.47
	// нижча теплота згоряння робочої маси мазуту
	mazutHeatValue := 39.48

	// частка золи, яка виходить з котла у вигляді леткої золи (вугілля)
	aCoal := 0.8
	// частка золи, яка виходить з котла у вигляді леткої золи (мазут)
	aMazut := 1.0

	// масовий вміст горючих речовин у леткій золі (вугілля)
	flammableSubstancesCoal := 1.5
	// масовий вміст горючих речовин у леткій золі (мазут)
	flammableSubstancesMazut := 0.0

	// масовий вміст золи в паливі на робочу масу, % (вугілля)
	arCoal := 25.2
	// масовий вміст золи в паливі на робочу масу, % (мазут)
	arMazut := 0.15

	// ефективність очищення димових газів від твердих частинок
	n := 0.985

	// емісія твердих частинок (вугілля)
	coalEmissionFactor := (math.Pow(10.0, 6) / coalHeatValue) * aCoal * (arCoal / (100 - flammableSubstancesCoal)) * (1 - n)
	// валовий викид твердих частинок (вугілля)
	coalEmissionValue := math.Pow(10.0, -6) * coalEmissionFactor * coalHeatValue * minerals.Coal

	// емісія твердих частинок (мазут)
	mazutEmissionFactor := (math.Pow(10.0, 6) / mazutHeatValue) * aMazut * (arMazut / (100 - flammableSubstancesMazut)) * (1 - n)
	// валовий викид твердих частинок (мазут)
	mazutEmissionValue := math.Pow(10.0, -6) * mazutEmissionFactor * mazutHeatValue * minerals.Mazut

	// при спалюванні природного газу тверді частинки відсутні, тоді
	// емісія твердих частинок (газ)
	gasEmissionFactor := 0.0
	// валовий викид твердих частинок (газ)
	gasEmissionValue := 0.0

	return CalculationResults{
		CoalEmissionFactor:  coalEmissionFactor,
		CoalEmissionValue:   coalEmissionValue,
		MazutEmissionFactor: mazutEmissionFactor,
		MazutEmissionValue:  mazutEmissionValue,
		GasEmissionFactor:   gasEmissionFactor,
		GasEmissionValue:    gasEmissionValue,
	}
}

func calculatorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var minerals Minerals
	err := json.NewDecoder(r.Body).Decode(&minerals)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	results := calculateResults(minerals)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func main() {
	// статичні файли
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	// маршрути для головної сторінки
	http.HandleFunc("/calculator", calculatorHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	log.Println("Сервер запущено на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
