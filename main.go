package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
)

// вхідні дані
type ElectricalEquipment struct {
	Name               string  `json:"name"`
	EfficiencyFactor   float64 `json:"efficiencyFactor"`
	LoadPowerFactor    float64 `json:"loadPowerFactor"`
	LoadVoltage        float64 `json:"loadVoltage"`
	Quantity           int     `json:"quantity"`
	RatedPower         int     `json:"ratedPower"`
	UtilizationRate    float64 `json:"utilizationRate"`
	ReactivePowerFactor float64 `json:"reactivePowerFactor"`
	CalculatedValues   EquipmentCalculatedValues `json:"calculatedValues"`
}

// для проміжних розрахунків
type EquipmentCalculatedValues struct {
	PowerTotal       float64 `json:"powerTotal"`       // n × Pн
	UtilizationPower float64 `json:"utilizationPower"` // n × Pн × Кв
	ReactivePower    float64 `json:"reactivePower"`    // n × Pн × Кв × tgφ
	SquaredPower     float64 `json:"squaredPower"`     // n × Pн²
	Current          float64 `json:"current"`          // Ip
}

// результати
type WorkshopResults struct {
	GroupUtilizationRate        float64 `json:"groupUtilizationRate"`
	EffectiveEquipmentCount     float64 `json:"effectiveEquipmentCount"`
	EstimatedActivePowerFactor  float64 `json:"estimatedActivePowerFactor"`
	EstimatedActiveLoad         float64 `json:"estimatedActiveLoad"`
	EstimatedReactiveLoad       float64 `json:"estimatedReactiveLoad"`
	FullPower                   float64 `json:"fullPower"`
	EstimatedGroupCurrent       float64 `json:"estimatedGroupCurrent"`

	UtilizationRateAll          float64 `json:"utilizationRateAll"`
	EffectiveEquipmentCountAll  int     `json:"effectiveEquipmentCountAll"`
	EstimatedActivePowerFactorAll float64 `json:"estimatedActivePowerFactorAll"`
	EstimatedActiveLoadTires    float64 `json:"estimatedActiveLoadTires"`
	EstimatedReactiveLoadTires  float64 `json:"estimatedReactiveLoadTires"`
	FullPowerTires              float64 `json:"fullPowerTires"`
	EstimatedGroupCurrentTires  float64 `json:"estimatedGroupCurrentTires"`
}

// для таблиці 6.3 Значення розрахункових коефіцієнтів КР для
// мереж живлення напругою до 1000 В
type TableEntry struct {
	N            int
	Coefficients map[float64]float64
}

// для таблиці 6.4 Значення розрахункових коефіцієнтів КР на шинах низької
// напруги цехових трансформаторів і магістральних шинопроводів
type TableRange struct {
    Start        int
    End          *int
    Coefficients map[float64]float64
}
// запит для розрахунку
type CalculationRequest struct {
	EquipmentList []ElectricalEquipment `json:"equipmentList"`
}

// функція розрахунку проміжних результатів по кожному ЕП
func calculateEquipmentValues(equipment ElectricalEquipment) ElectricalEquipment {
	// n * Pн
	powerTotal := float64(equipment.Quantity) * float64(equipment.RatedPower)
	// n * Pн * Кв
	utilizationPower := powerTotal * equipment.UtilizationRate
	// n * Pн * Кв * tgφ
	reactivePower := utilizationPower * equipment.ReactivePowerFactor
	// n * Pн^2
	squaredPower := float64(equipment.Quantity) * math.Pow(float64(equipment.RatedPower), 2.0)
	// (n * Pн) / √3 * Uн * cosφ * nн
	current := powerTotal / (math.Sqrt(3.0) * equipment.LoadVoltage *
		equipment.LoadPowerFactor * equipment.EfficiencyFactor)
	currentTruncate := math.Floor(current*10) / 10

	equipment.CalculatedValues = EquipmentCalculatedValues{
		PowerTotal:       powerTotal,
		UtilizationPower: utilizationPower,
		ReactivePower:    reactivePower,
		SquaredPower:     squaredPower,
		Current:          currentTruncate,
	}

	return equipment
}

// функція для пошуку розрахункового коефіцієнту активної потужності по таблиці 6.3
func findCoefficient(n int, coefficient float64, tableData []TableEntry) float64 {
	for _, entry := range tableData {
		if entry.N == n {
			if value, exists := entry.Coefficients[coefficient]; exists {
				return value
			}
		}
	}
	return 1.25
}

// функція для пошуку розрахункового коефіцієнту активної потужності по таблиці 6.4
func findCoefficientRange(n int, coefficient float64, tableDataRange []TableRange) float64 {
	for _, entry := range tableDataRange {
		if n >= entry.Start && (entry.End == nil || n <= *entry.End) {
			if value, exists := entry.Coefficients[coefficient]; exists {
				return value
			}
		}
	}
	return 0.7
}

func intPtr(i int) *int {
    return &i
}

func calculateWorkshopResults(equipmentList []ElectricalEquipment) WorkshopResults {
	// значення з таблиці 6.3
	tableData := []TableEntry{
			{
					N: 1,
					Coefficients: map[float64]float64{
							0.1:  8.00,
							0.15: 5.33,
							0.2:  4.00,
							0.3:  2.67,
							0.4:  2.00,
							0.5:  1.60,
							0.6:  1.33,
							0.7:  1.14,
							0.8:  1.0,
					},
			},
			{
					N: 2,
					Coefficients: map[float64]float64{
							0.1:  6.22,
							0.15: 4.33,
							0.2:  3.39,
							0.3:  2.45,
							0.4:  1.98,
							0.5:  1.60,
							0.6:  1.33,
							0.7:  1.14,
							0.8:  1.0,
					},
			},
			{
					N: 3,
					Coefficients: map[float64]float64{
							0.1:  4.06,
							0.15: 2.89,
							0.2:  2.31,
							0.3:  1.74,
							0.4:  1.45,
							0.5:  1.34,
							0.6:  1.22,
							0.7:  1.14,
							0.8:  1.0,
					},
			},
			{
					N: 4,
					Coefficients: map[float64]float64{
							0.1:  3.23,
							0.15: 2.29,
							0.2:  1.83,
							0.3:  1.39,
							0.4:  1.21,
							0.5:  1.13,
							0.6:  1.08,
							0.7:  1.03,
							0.8:  1.0,
					},
			},
			{
					N: 5,
					Coefficients: map[float64]float64{
							0.1:  2.84,
							0.15: 2.06,
							0.2:  1.65,
							0.3:  1.31,
							0.4:  1.15,
							0.5:  1.10,
							0.6:  1.05,
							0.7:  1.01,
							0.8:  1.0,
					},
			},
			{
					N: 6,
					Coefficients: map[float64]float64{
							0.1:  2.64,
							0.15: 1.96,
							0.2:  1.62,
							0.3:  1.28,
							0.4:  1.14,
							0.5:  1.13,
							0.6:  1.06,
							0.7:  1.01,
							0.8:  1.0,
					},
			},
			{
					N: 7,
					Coefficients: map[float64]float64{
							0.1:  2.49,
							0.15: 1.86,
							0.2:  1.54,
							0.3:  1.23,
							0.4:  1.12,
							0.5:  1.10,
							0.6:  1.04,
							0.7:  1.0,
							0.8:  1.0,
					},
			},
			{
					N: 8,
					Coefficients: map[float64]float64{
							0.1:  2.37,
							0.15: 1.78,
							0.2:  1.48,
							0.3:  1.19,
							0.4:  1.10,
							0.5:  1.08,
							0.6:  1.02,
							0.7:  1.0,
							0.8:  1.0,
					},
			},
			{
					N: 9,
					Coefficients: map[float64]float64{
							0.1:  2.27,
							0.15: 1.71,
							0.2:  1.43,
							0.3:  1.16,
							0.4:  1.09,
							0.5:  1.07,
							0.6:  1.01,
							0.7:  1.0,
							0.8:  1.0,
					},
			},
			{
					N: 10,
					Coefficients: map[float64]float64{
							0.1:  2.18,
							0.15: 1.65,
							0.2:  1.39,
							0.3:  1.13,
							0.4:  1.07,
							0.5:  1.05,
							0.6:  1.0,
							0.7:  1.0,
							0.8:  1.0,
					},
			},
			{
					N: 12,
					Coefficients: map[float64]float64{
							0.1:  2.04,
							0.15: 1.56,
							0.2:  1.32,
							0.3:  1.08,
							0.4:  1.05,
							0.5:  1.03,
							0.6:  1.0,
							0.7:  1.0,
							0.8:  1.0,
					},
			},
			{
					N: 14,
					Coefficients: map[float64]float64{
							0.1:  1.94,
							0.15: 1.49,
							0.2:  1.27,
							0.3:  1.05,
							0.4:  1.02,
							0.5:  1.0,
							0.6:  1.0,
							0.7:  1.0,
							0.8:  1.0,
					},
			},
			{
					N: 16,
					Coefficients: map[float64]float64{
							0.1:  1.85,
							0.15: 1.43,
							0.2:  1.23,
							0.3:  1.02,
							0.4:  1.0,
							0.5:  1.0,
							0.6:  1.0,
							0.7:  1.0,
							0.8:  1.0,
					},
			},
			{
					N: 18,
					Coefficients: map[float64]float64{
							0.1:  1.78,
							0.15: 1.39,
							0.2:  1.19,
							0.3:  1.0,
							0.4:  1.0,
							0.5:  1.0,
							0.6:  1.0,
							0.7:  1.0,
							0.8:  1.0,
					},
			},
			{
					N: 20,
					Coefficients: map[float64]float64{
							0.1:  1.72,
							0.15: 1.35,
							0.2:  1.16,
							0.3:  1.0,
							0.4:  1.0,
							0.5:  1.0,
							0.6:  1.0,
							0.7:  1.0,
							0.8:  1.0,
					},
			},
			{
					N: 25,
					Coefficients: map[float64]float64{
							0.1:  1.60,
							0.15: 1.27,
							0.2:  1.10,
							0.3:  1.0,
							0.4:  1.0,
							0.5:  1.0,
							0.6:  1.0,
							0.7:  1.0,
							0.8:  1.0,
					},
			},
			{
					N: 30,
					Coefficients: map[float64]float64{
							0.1:  1.51,
							0.15: 1.21,
							0.2:  1.05,
							0.3:  1.0,
							0.4:  1.0,
							0.5:  1.0,
							0.6:  1.0,
							0.7:  1.0,
							0.8:  1.0,
					},
			},
			{
					N: 35,
					Coefficients: map[float64]float64{
							0.1:  1.44,
							0.15: 1.16,
							0.2:  1.0,
							0.3:  1.0,
							0.4:  1.0,
							0.5:  1.0,
							0.6:  1.0,
							0.7:  1.0,
							0.8:  1.0,
					},
			},
			{
					N: 40,
					Coefficients: map[float64]float64{
							0.1:  1.40,
							0.15: 1.13,
							0.2:  1.0,
							0.3:  1.0,
							0.4:  1.0,
							0.5:  1.0,
							0.6:  1.0,
							0.7:  1.0,
							0.8:  1.0,
					},
			},
			{
					N: 50,
					Coefficients: map[float64]float64{
							0.1:  1.30,
							0.15: 1.07,
							0.2:  1.0,
							0.3:  1.0,
							0.4:  1.0,
							0.5:  1.0,
							0.6:  1.0,
							0.7:  1.0,
							0.8:  1.0,
					},
			},
			{
					N: 60,
					Coefficients: map[float64]float64{
							0.1:  1.25,
							0.15: 1.03,
							0.2:  1.0,
							0.3:  1.0,
							0.4:  1.0,
							0.5:  1.0,
							0.6:  1.0,
							0.7:  1.0,
							0.8:  1.0,
					},
			},
			{
					N: 80,
					Coefficients: map[float64]float64{
							0.1:  1.16,
							0.15: 1.0,
							0.2:  1.0,
							0.3:  1.0,
							0.4:  1.0,
							0.5:  1.0,
							0.6:  1.0,
							0.7:  1.0,
							0.8:  1.0,
					},
			},
			{
					N: 100,
					Coefficients: map[float64]float64{
							0.1:  1.0,
							0.15: 1.0,
							0.2:  1.0,
							0.3:  1.0,
							0.4:  1.0,
							0.5:  1.0,
							0.6:  1.0,
							0.7:  1.0,
							0.8:  1.0,
					},
			},
	}

	// значення з таблиці 6.4
	tableDataRange := []TableRange{
		{
				Start: 1,
				End:  intPtr(1), // або &end, якщо потрібна конкретна змінна
				Coefficients: map[float64]float64{
						0.1:  8.00,
						0.15: 5.33,
						0.2:  4.00,
						0.3:  2.67,
						0.4:  2.00,
						0.5:  1.60,
						0.6:  1.33,
						0.7:  1.14,
				},
		},
		{
				Start: 2,
				End:   intPtr(2),
				Coefficients: map[float64]float64{
						0.1:  5.01,
						0.15: 3.44,
						0.2:  2.69,
						0.3:  1.90,
						0.4:  1.52,
						0.5:  1.24,
						0.6:  1.11,
						0.7:  1.0,
				},
		},
		{
				Start: 3,
				End:   intPtr(3),
				Coefficients: map[float64]float64{
						0.1:  2.40,
						0.15: 2.17,
						0.2:  1.80,
						0.3:  1.42,
						0.4:  1.23,
						0.5:  1.14,
						0.6:  1.08,
						0.7:  1.0,
				},
		},
		{
				Start: 4,
				End:   intPtr(4),
				Coefficients: map[float64]float64{
						0.1:  2.28,
						0.15: 1.73,
						0.2:  1.46,
						0.3:  1.19,
						0.4:  1.06,
						0.5:  1.04,
						0.6:  1.0,
						0.7:  0.97,
				},
		},
		{
				Start: 5,
				End:   intPtr(5),
				Coefficients: map[float64]float64{
						0.1:  1.31,
						0.15: 1.12,
						0.2:  1.02,
						0.3:  1.0,
						0.4:  0.98,
						0.5:  0.96,
						0.6:  0.94,
						0.7:  0.93,
				},
		},
		{
				Start: 6,
				End:   intPtr(8),
				Coefficients: map[float64]float64{
						0.1:  1.20,
						0.15: 1.0,
						0.2:  0.96,
						0.3:  0.95,
						0.4:  0.94,
						0.5:  0.93,
						0.6:  0.92,
						0.7:  0.91,
				},
		},
		{
				Start: 9,
				End:   intPtr(10),
				Coefficients: map[float64]float64{
						0.1:  1.10,
						0.15: 0.97,
						0.2:  0.91,
						0.3:  0.90,
						0.4:  0.90,
						0.5:  0.90,
						0.6:  0.90,
						0.7:  0.90,
				},
		},
		{
				Start: 10,
				End:   intPtr(25),
				Coefficients: map[float64]float64{
						0.1:  0.80,
						0.15: 0.80,
						0.2:  0.80,
						0.3:  0.85,
						0.4:  0.85,
						0.5:  0.85,
						0.6:  0.90,
						0.7:  0.90,
				},
		},
		{
				Start: 25,
				End:   intPtr(50),
				Coefficients: map[float64]float64{
						0.1:  0.75,
						0.15: 0.75,
						0.2:  0.75,
						0.3:  0.75,
						0.4:  0.75,
						0.5:  0.80,
						0.6:  0.85,
						0.7:  0.85,
				},
		},
		{
				Start: 50,
				End:   nil,
				Coefficients: map[float64]float64{
						0.1:  0.65,
						0.15: 0.65,
						0.2:  0.65,
						0.3:  0.70,
						0.4:  0.70,
						0.5:  0.75,
						0.6:  0.80,
						0.7:  0.80,
				},
		},

		
	}

	// Σ n * Pн
	totalPower := 0.0
	// Σ n * Pн * Кв
	totalUtilizationPower := 0.0
	// Σ n * Pн^2
	totalSquaredPower := 0.0
	// Σ n * Pн * Кв * tgφ
	totalReactivePower := 0.0

	for _, equipment := range equipmentList {
		totalPower += equipment.CalculatedValues.PowerTotal
		totalUtilizationPower += equipment.CalculatedValues.UtilizationPower
		totalSquaredPower += equipment.CalculatedValues.SquaredPower
		totalReactivePower += equipment.CalculatedValues.ReactivePower
	}

	// Груповий коефіцієнт використання
	groupUtilizationRate := 0.0
	if totalPower > 0 {
		groupUtilizationRate = totalUtilizationPower / totalPower
	}

	// Ефективна кількість ЕП
	effectiveEquipmentCount := 0.0
	if totalSquaredPower > 0 {
		effectiveEquipmentCount = math.Ceil(math.Pow(totalPower, 2.0) / totalSquaredPower)
	}

	// Розрахунковий коефіцієнт активної потужності
	roundedEffectiveEquipmentCount := int(math.Ceil(effectiveEquipmentCount))

	// Округлення groupUtilizationRate до першого знаку після коми для пошуку в таблиці
	roundedGroupUtilizationRate := math.Round(groupUtilizationRate*10.0) / 10.0

	// Отримання коефіцієнту з таблиці з перевіркою на null
	estimatedActivePowerFactor := findCoefficient(roundedEffectiveEquipmentCount, roundedGroupUtilizationRate, tableData)

	// Розрахункове активне навантаження
	estimatedActiveLoad := estimatedActivePowerFactor * totalUtilizationPower

	// Розрахункове реактивне навантаження
	estimatedReactiveLoad := totalReactivePower

	// Повна потужність
	fullPower := math.Sqrt(math.Pow(estimatedActiveLoad, 2.0) +
		math.Pow(estimatedReactiveLoad, 2.0))

	// Розрахунковий груповий струм
	estimatedGroupCurrent := 0.0
	if len(equipmentList) > 0 && equipmentList[0].LoadVoltage > 0 {
		estimatedGroupCurrent = estimatedActiveLoad / equipmentList[0].LoadVoltage
	}

	// цех в цілому
	// кількість ЕП
	// equipmentNumber := 81
	// n * Pн
	powerAll := 2330
	// n * Pн * Кв
	utilizationPowerAll := 752
	// n * Pн * Кв * tgφ
	reactivePowerAll := 657
	// n * Pн^2
	squaredPowerAll := 96388

	// Коефіцієнти використання цеху в цілому
	utilizationRateAll := float64(utilizationPowerAll) / float64(powerAll)

	// Ефективна кількість ЕП цеху в цілому
	effectiveEquipmentCountAll := int(math.Pow(float64(powerAll), 2.0) / float64(squaredPowerAll))

	// Округлення utilizationRateAll до першого знаку після коми для пошуку в таблиці
	roundedUtilizationRateAll := math.Round(utilizationRateAll*10.0) / 10.0

	// Розрахунковий коефіцієнт активної потужності цеху в цілому
	estimatedActivePowerFactorAll := findCoefficientRange(effectiveEquipmentCountAll, roundedUtilizationRateAll, tableDataRange)

	// Розрахункове активне навантаження на шинах 0,38 кВ ТП
	estimatedActiveLoadTires := estimatedActivePowerFactorAll * float64(utilizationPowerAll)

	// Розрахункове реактивне навантаження на шинах 0,38 кВ ТП
	estimatedReactiveLoadTires := estimatedActivePowerFactorAll * float64(reactivePowerAll)

	// Повна потужність на шинах 0,38 кВ ТП
	fullPowerTires := math.Sqrt(math.Pow(estimatedActiveLoadTires, 2.0) +
		math.Pow(estimatedReactiveLoadTires, 2.0))

	// Розрахунковий груповий струм на шинах 0,38 кВ ТП
	estimatedGroupCurrentTires := 0.0
	if len(equipmentList) > 0 && equipmentList[0].LoadVoltage > 0 {
		estimatedGroupCurrentTires = estimatedActiveLoadTires / equipmentList[0].LoadVoltage
	}

	return WorkshopResults{
		GroupUtilizationRate:           groupUtilizationRate,
		EffectiveEquipmentCount:        effectiveEquipmentCount,
		EstimatedActivePowerFactor:     estimatedActivePowerFactor,
		EstimatedActiveLoad:            estimatedActiveLoad,
		EstimatedReactiveLoad:          estimatedReactiveLoad,
		FullPower:                      fullPower,
		EstimatedGroupCurrent:          estimatedGroupCurrent,
		UtilizationRateAll:             roundedUtilizationRateAll,
		EffectiveEquipmentCountAll:     effectiveEquipmentCountAll,
		EstimatedActivePowerFactorAll:  estimatedActivePowerFactorAll,
		EstimatedActiveLoadTires:       estimatedActiveLoadTires,
		EstimatedReactiveLoadTires:     estimatedReactiveLoadTires,
		FullPowerTires:                 fullPowerTires,
		EstimatedGroupCurrentTires:     estimatedGroupCurrentTires,
	}
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request CalculationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Розрахунок проміжних значень для кожного обладнання
	calculatedEquipment := make([]ElectricalEquipment, len(request.EquipmentList))
	for i, equipment := range request.EquipmentList {
		calculatedEquipment[i] = calculateEquipmentValues(equipment)
	}

	// Розрахунок загальних результатів
	results := calculateWorkshopResults(calculatedEquipment)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func singleEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var equipment ElectricalEquipment
	err := json.NewDecoder(r.Body).Decode(&equipment)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Розрахунок проміжних значень для обладнання
	calculatedEquipment := calculateEquipmentValues(equipment)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(calculatedEquipment)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func main() {
	// статичні файли
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	// маршрути API
	http.HandleFunc("/api/calculate", calculateHandler)
	http.HandleFunc("/api/equipment", singleEquipmentHandler)
	
	// маршрут для головної сторінки
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	log.Println("Сервер запущено на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}