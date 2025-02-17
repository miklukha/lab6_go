package main

import (
	"encoding/json"
	"math"
	"net/http"
)

// компоненти палива
type FuelComposition struct {
	HP float64 `json:"hp"`
	CP float64 `json:"cp"`
	SP float64 `json:"sp"`
	NP float64 `json:"np"`
	OP float64 `json:"op"`
	WP float64 `json:"wp"`
	AP float64 `json:"ap"`
}

// результати розрахунків
type CalculationResults struct {
	DryMassCoefficient          float64            `json:"dryMassCoefficient"`
	CombustibleMassCoefficient  float64            `json:"combustibleMassCoefficient"`
	DryComposition              map[string]float64 `json:"dryComposition"`
	CombustibleComposition      map[string]float64 `json:"combustibleComposition"`
	LowerHeatingValue           float64            `json:"lowerHeatingValue"`
	LowerDryHeatingValue        float64            `json:"lowerDryHeatingValue"`
	LowerCombustibleHeatingValue float64            `json:"lowerCombustibleHeatingValue"`
}

func calculateResults(composition FuelComposition) CalculationResults {
	// кофіцієнт переходу від робочої до сухої маси
	kpc := 100.0 / (100.0 - composition.WP)
	// кофіцієнт переходу від робочої до горючої маси
	kpg := 100.0 / (100.0 - composition.WP - composition.AP)

	// cклад сухої маси палива
	dryComposition := map[string]float64{
		"H^C": composition.HP * kpc,
		"C^C": composition.CP * kpc,
		"S^C": composition.SP * kpc,
		"N^C": composition.NP * kpc,
		"O^C": composition.OP * kpc,
		"A^C": composition.AP * kpc,
	}

	// cклад горючої маси палива
	combustibleComposition := map[string]float64{
		"H^Г": composition.HP * kpg,
		"C^Г": composition.CP * kpg,
		"S^Г": composition.SP * kpg,
		"N^Г": composition.NP * kpg,
		"O^Г": composition.OP * kpg,
	}

	// нижча теплота згоряння для робочої маси
	qph := (339*composition.CP +
		1030*composition.HP -
		108.8*(composition.OP-composition.SP) -
		25*composition.WP) / 1000

	// нижча теплота згоряння для сухої маси
	qch := (qph + 0.025*composition.WP) * (100.0 / (100.0 - composition.WP))

	// нижча теплота згоряння для горючої маси
	qgh := (qph + 0.025*composition.WP) * (100.0 / (100.0 - composition.WP - composition.AP))

	// Округлюємо результати до 2 знаків після коми
	for k, v := range dryComposition {
		dryComposition[k] = math.Round(v*100) / 100
	}
	for k, v := range combustibleComposition {
		combustibleComposition[k] = math.Round(v*100) / 100
	}

	return CalculationResults{
		DryMassCoefficient:          math.Round(kpc*100) / 100,
		CombustibleMassCoefficient:  math.Round(kpg*100) / 100,
		DryComposition:              dryComposition,
		CombustibleComposition:      combustibleComposition,
		LowerHeatingValue:           math.Round(qph*100) / 100,
		LowerDryHeatingValue:        math.Round(qch*100) / 100,
		LowerCombustibleHeatingValue: math.Round(qgh*100) / 100,
	}
}

// обробка запитів до калькулятора 1
func HandleCalculator1(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var composition FuelComposition
		err := json.NewDecoder(r.Body).Decode(&composition)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		results := calculateResults(composition)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
		return
	}

	http.ServeFile(w, r, "templates/calculator1.html")
}

// налаштування маршрутів для калькулятора 1
func SetupCalculator1Routes() {
	http.HandleFunc("/calculator1", HandleCalculator1)
}