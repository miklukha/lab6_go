package main

import (
	"encoding/json"
	"math"
	"net/http"
)

// компоненти мазуту
type MazutComposition struct {
	CarbonCombustible     float64 `json:"carbonCombustible"`
	HydrogenCombustible   float64 `json:"hydrogenCombustible"`
	OxygenCombustible     float64 `json:"oxygenCombustible"`
	SulfurCombustible     float64 `json:"sulfurCombustible"`
	VanadiumCombustible   float64 `json:"vanadiumCombustible"`
	MoistureContent       float64 `json:"moistureContent"`
	AshDry                float64 `json:"ashDry"`
	HeatingValueCombustible float64 `json:"heatingValueCombustible"`
}

// результати розрахунків
type MazutResults struct {
	WorkingComposition  map[string]float64 `json:"workingComposition"`
	WorkingHeatingValue float64            `json:"workingHeatingValue"`
}

func calculateMazutResults(composition MazutComposition) MazutResults {
	// формула перерахунку складу палива
	conversionFactor := (100.0 - composition.MoistureContent - composition.AshDry) / 100.0

	// перерахунок для кожного компонента
	workingComposition := map[string]float64{
		"C^P": composition.CarbonCombustible * conversionFactor,
		"H^P": composition.HydrogenCombustible * conversionFactor,
		"O^P": composition.OxygenCombustible * conversionFactor,
		"S^P": composition.SulfurCombustible * conversionFactor,
		// віднімання - отримання маси без вологи (суха маса)
		// оскільки ванадій вказується відносно маси без вологи
		// ділення - отримання відсотку
		"V^P": composition.VanadiumCombustible * (100.0 - composition.MoistureContent) / 100.0,
		// віднімання - отримання відсотка маси без вологи
		// оскільки зола вказується відносно маси без вологи
		// ділення - отримання відсотку
		"A^P": composition.AshDry * (100.0 - composition.MoistureContent) / 100.0,
	}

	// розрахунок нижчої теплоти згорання
	workingHeatingValue := composition.HeatingValueCombustible*
		(100.0-composition.MoistureContent-composition.AshDry)/100.0 -
		0.025*composition.MoistureContent

	// Округлюємо результати до 2 знаків після коми
	for k, v := range workingComposition {
		workingComposition[k] = math.Round(v*100) / 100
	}

	return MazutResults{
		WorkingComposition:  workingComposition,
		WorkingHeatingValue: math.Round(workingHeatingValue*100) / 100,
	}
}

// обробка запитів до калькулятора мазуту
func HandleCalculator2(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var composition MazutComposition
		err := json.NewDecoder(r.Body).Decode(&composition)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		results := calculateMazutResults(composition)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
		return
	}

	http.ServeFile(w, r, "templates/calculator2.html")
}

// налаштування маршрутів для калькулятора 2
func SetupCalculator2Routes() {
	http.HandleFunc("/calculator2", HandleCalculator2)
}