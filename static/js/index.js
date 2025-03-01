document.addEventListener('DOMContentLoaded', function () {
  // Початкові налаштування
  let equipmentList = [];
  const equipmentContainer = document.getElementById('equipment-container');
  const addEquipmentBtn = document.getElementById('add-equipment');
  const calculateBtn = document.getElementById('calculate-btn');
  const resultsContainer = document.getElementById('results-container');

  // Таблиця коефіцієнтів 6.3
  const tableData = [
    {
      n: 1,
      coefficients: {
        0.1: 8.0,
        0.15: 5.33,
        0.2: 4.0,
        0.3: 2.67,
        0.4: 2.0,
        0.5: 1.6,
        0.6: 1.33,
        0.7: 1.14,
        0.8: 1.0,
      },
    },
  ];

  // Таблиця коефіцієнтів 6.4
  const tableDataRange = [
    {
      start: 1,
      end: 1,
      coefficients: {
        0.1: 8.0,
        0.15: 5.33,
        0.2: 4.0,
        0.3: 2.67,
        0.4: 2.0,
        0.5: 1.6,
        0.6: 1.33,
        0.7: 1.14,
      },
    },
  ];

  // Класи для представлення даних

  class EquipmentCalculatedValues {
    constructor() {
      this.powerTotal = 0.0; // n × Pн
      this.utilizationPower = 0.0; // n × Pн × Кв
      this.reactivePower = 0.0; // n × Pн × Кв × tgφ
      this.squaredPower = 0.0; // n × Pн²
      this.current = 0.0; // Ip
    }
  }

  class ElectricalEquipment {
    constructor() {
      this.name = '';
      this.efficiencyFactor = 0.0;
      this.loadPowerFactor = 0.0;
      this.loadVoltage = 0.0;
      this.quantity = 0;
      this.ratedPower = 0;
      this.utilizationRate = 0.0;
      this.reactivePowerFactor = 0.0;
      this.calculatedValues = new EquipmentCalculatedValues();
    }
  }

  // Додавання початкових електроприймачів
  for (let i = 0; i < 8; i++) {
    addEquipmentCard();
  }

  // Обробка подій
  addEquipmentBtn.addEventListener('click', addEquipmentCard);
  calculateBtn.addEventListener('click', calculateResults);

  function addEquipmentCard() {
    const template = document.getElementById('equipment-template');
    const clone = template.content.cloneNode(true);
    const equipmentNumber =
      document.querySelectorAll('.equipment-card').length + 1;

    // Встановлення номера електроприймача
    clone.querySelector(
      '.equipment-number',
    ).textContent = `Електроприймач ${equipmentNumber}`;

    // Додавання обробників подій для розрахунку проміжних значень
    setupEquipmentCardEvents(clone);

    // Додавання кнопки видалення
    clone.querySelector('.delete-btn').addEventListener('click', function (e) {
      const card = e.target.closest('.equipment-card');
      if (card) {
        card.remove();
        updateEquipmentNumbers();
      }
    });

    equipmentContainer.appendChild(clone);
  }

  function updateEquipmentNumbers() {
    const cards = document.querySelectorAll('.equipment-card');
    cards.forEach((card, index) => {
      card.querySelector('.equipment-number').textContent = `Електроприймач ${
        index + 1
      }`;
    });
  }

  function setupEquipmentCardEvents(cardElement) {
    const inputs = cardElement.querySelectorAll('input');
    inputs.forEach(input => {
      input.addEventListener('input', function () {
        const card = this.closest('.equipment-card');
        calculateEquipmentValues(card);
      });
    });
  }

  function calculateEquipmentValues(card) {
    // Отримання введених значень
    const name = card.querySelector('.equipment-name').value;
    const efficiencyFactor =
      parseFloat(card.querySelector('.efficiency-factor').value) || 0;
    const loadPowerFactor =
      parseFloat(card.querySelector('.load-power-factor').value) || 0;
    const loadVoltage =
      parseFloat(card.querySelector('.load-voltage').value) || 0;
    const quantity = parseInt(card.querySelector('.quantity').value) || 0;
    const ratedPower = parseInt(card.querySelector('.rated-power').value) || 0;
    const utilizationRate =
      parseFloat(card.querySelector('.utilization-rate').value) || 0;
    const reactivePowerFactor =
      parseFloat(card.querySelector('.reactive-power-factor').value) || 0;

    // Перевірка на наявність всіх потрібних значень
    if (
      efficiencyFactor &&
      loadPowerFactor &&
      loadVoltage &&
      quantity &&
      ratedPower &&
      utilizationRate &&
      reactivePowerFactor
    ) {
      // n * Pн
      const powerTotal = quantity * ratedPower;
      // n * Pн * Кв
      const utilizationPower = powerTotal * utilizationRate;
      // n * Pн * Кв * tgφ
      const reactivePower = utilizationPower * reactivePowerFactor;
      // n * Pн^2
      const squaredPower = quantity * Math.pow(ratedPower, 2.0);
      // (n * Pн) / √3 * Uн * cosφ * nн
      const current =
        powerTotal /
        (Math.sqrt(3.0) * loadVoltage * loadPowerFactor * efficiencyFactor);
      const currentTruncate = Math.floor(current * 10) / 10;

      // Відображення обчислених значень
      card.querySelector('.calculated-values').classList.remove('hidden');
      card.querySelector('.power-total').textContent = powerTotal.toFixed(2);
      card.querySelector('.utilization-power').textContent =
        utilizationPower.toFixed(2);
      card.querySelector('.reactive-power').textContent =
        reactivePower.toFixed(2);
      card.querySelector('.squared-power').textContent =
        squaredPower.toFixed(2);
      card.querySelector('.current').textContent = currentTruncate.toFixed(2);
    }
  }

  function getAllEquipment() {
    const equipmentCards = document.querySelectorAll('.equipment-card');
    const equipmentList = [];

    equipmentCards.forEach(card => {
      const equipment = new ElectricalEquipment();
      equipment.name = card.querySelector('.equipment-name').value || '';
      equipment.efficiencyFactor =
        parseFloat(card.querySelector('.efficiency-factor').value) || 0;
      equipment.loadPowerFactor =
        parseFloat(card.querySelector('.load-power-factor').value) || 0;
      equipment.loadVoltage =
        parseFloat(card.querySelector('.load-voltage').value) || 0;
      equipment.quantity = parseInt(card.querySelector('.quantity').value) || 0;
      equipment.ratedPower =
        parseInt(card.querySelector('.rated-power').value) || 0;
      equipment.utilizationRate =
        parseFloat(card.querySelector('.utilization-rate').value) || 0;
      equipment.reactivePowerFactor =
        parseFloat(card.querySelector('.reactive-power-factor').value) || 0;

      if (
        equipment.efficiencyFactor &&
        equipment.loadPowerFactor &&
        equipment.loadVoltage &&
        equipment.quantity &&
        equipment.ratedPower &&
        equipment.utilizationRate &&
        equipment.reactivePowerFactor
      ) {
        // Розрахунок проміжних значень
        // n * Pн
        equipment.calculatedValues.powerTotal =
          equipment.quantity * equipment.ratedPower;
        // n * Pн * Кв
        equipment.calculatedValues.utilizationPower =
          equipment.calculatedValues.powerTotal * equipment.utilizationRate;
        // n * Pн * Кв * tgφ
        equipment.calculatedValues.reactivePower =
          equipment.calculatedValues.utilizationPower *
          equipment.reactivePowerFactor;
        // n * Pн^2
        equipment.calculatedValues.squaredPower =
          equipment.quantity * Math.pow(equipment.ratedPower, 2.0);
        // (n * Pн) / √3 * Uн * cosφ * nн
        equipment.calculatedValues.current =
          equipment.calculatedValues.powerTotal /
          (Math.sqrt(3.0) *
            equipment.loadVoltage *
            equipment.loadPowerFactor *
            equipment.efficiencyFactor);
        equipment.calculatedValues.current =
          Math.floor(equipment.calculatedValues.current * 10) / 10;

        equipmentList.push(equipment);
      }
    });

    return equipmentList;
  }

  function findCoefficient(n, coefficient) {
    for (const entry of tableData) {
      if (entry.n === n) {
        if (entry.coefficients[coefficient] !== undefined) {
          return entry.coefficients[coefficient];
        }
      }
    }
    return 1.25; // Значення за замовчуванням
  }

  function findCoefficientRange(n, coefficient) {
    for (const entry of tableDataRange) {
      if (n >= entry.start && (entry.end === null || n <= entry.end)) {
        if (entry.coefficients[coefficient] !== undefined) {
          return entry.coefficients[coefficient];
        }
      }
    }
    return 0.7; // Значення за замовчуванням
  }

  function calculateResults() {
    const equipmentList = getAllEquipment();

    if (equipmentList.length === 0) {
      alert('Додайте хоча б один електроприймач з усіма заповненими полями!');
      return;
    }

    // Σ n * Pн
    const totalPower = equipmentList.reduce(
      (sum, equipment) => sum + equipment.calculatedValues.powerTotal,
      0,
    );
    // Σ n * Pн * Кв
    const totalUtilizationPower = equipmentList.reduce(
      (sum, equipment) => sum + equipment.calculatedValues.utilizationPower,
      0,
    );
    // Σ n * Pн^2
    const totalSquaredPower = equipmentList.reduce(
      (sum, equipment) => sum + equipment.calculatedValues.squaredPower,
      0,
    );
    // Σ n * Pн * Кв * tgφ
    const totalReactivePower = equipmentList.reduce(
      (sum, equipment) => sum + equipment.calculatedValues.reactivePower,
      0,
    );

    // Груповий коефіцієнт використання
    let groupUtilizationRate = 0;
    if (totalPower > 0) {
      groupUtilizationRate = totalUtilizationPower / totalPower;
    }

    // Ефективна кількість ЕП
    let effectiveEquipmentCount = 0;
    if (totalSquaredPower > 0) {
      effectiveEquipmentCount = Math.ceil(
        Math.pow(totalPower, 2.0) / totalSquaredPower,
      );
    }

    // Розрахунковий коефіцієнт активної потужності
    const roundedEffectiveEquipmentCount = Math.ceil(effectiveEquipmentCount);

    // Округлення groupUtilizationRate до першого знаку після коми для пошуку в таблиці
    const roundedGroupUtilizationRate =
      Math.round(groupUtilizationRate * 10.0) / 10.0;

    // Отримання коефіцієнту з таблиці
    const estimatedActivePowerFactor = findCoefficient(
      roundedEffectiveEquipmentCount,
      roundedGroupUtilizationRate,
    );

    // Розрахункове активне навантаження
    const estimatedActiveLoad =
      estimatedActivePowerFactor * totalUtilizationPower;

    // Розрахункове реактивне навантаження
    const estimatedReactiveLoad = totalReactivePower;

    // Повна потужність
    const fullPower = Math.sqrt(
      Math.pow(estimatedActiveLoad, 2.0) + Math.pow(estimatedReactiveLoad, 2.0),
    );

    // Розрахунковий груповий струм
    let estimatedGroupCurrent = 0;
    if (equipmentList.length > 0 && equipmentList[0].loadVoltage > 0) {
      estimatedGroupCurrent =
        estimatedActiveLoad / equipmentList[0].loadVoltage;
    }

    // Цех в цілому
    // кількість ЕП
    const equipmentNumber = 81;
    // n * Pн
    const powerAll = 2330;
    // n * Pн * Кв
    const utilizationPowerAll = 752;
    // n * Pн * Кв * tgφ
    const reactivePowerAll = 657;
    // n * Pн^2
    const squaredPowerAll = 96388;

    // Коефіцієнти використання цеху в цілому
    const utilizationRateAll = utilizationPowerAll / powerAll;

    // Ефективна кількість ЕП цеху в цілому
    const effectiveEquipmentCountAll = Math.floor(
      Math.pow(powerAll, 2.0) / squaredPowerAll,
    );

    // Округлення utilizationRateAll до першого знаку після коми для пошуку в таблиці
    const roundedUtilizationRateAll =
      Math.round(utilizationRateAll * 10.0) / 10.0;

    // Розрахунковий коефіцієнт активної потужності цеху в цілому
    const estimatedActivePowerFactorAll = findCoefficientRange(
      effectiveEquipmentCountAll,
      roundedUtilizationRateAll,
    );

    // Розрахункове активне навантаження на шинах 0,38 кВ ТП
    const estimatedActiveLoadTires =
      estimatedActivePowerFactorAll * utilizationPowerAll;

    // Розрахункове реактивне навантаження на шинах 0,38 кВ ТП
    const estimatedReactiveLoadTires =
      estimatedActivePowerFactorAll * reactivePowerAll;

    // Повна потужність на шинах 0,38 кВ ТП
    const fullPowerTires = Math.sqrt(
      Math.pow(estimatedActiveLoadTires, 2.0) +
        Math.pow(estimatedReactiveLoadTires, 2.0),
    );

    // Розрахунковий груповий струм на шинах 0,38 кВ ТП
    let estimatedGroupCurrentTires = 0;
    if (equipmentList.length > 0 && equipmentList[0].loadVoltage > 0) {
      estimatedGroupCurrentTires =
        estimatedActiveLoadTires / equipmentList[0].loadVoltage;
    }

    // Відображення результатів
    displayResults({
      groupUtilizationRate,
      effectiveEquipmentCount,
      estimatedActivePowerFactor,
      estimatedActiveLoad,
      estimatedReactiveLoad,
      fullPower,
      estimatedGroupCurrent,
      utilizationRateAll: roundedUtilizationRateAll,
      effectiveEquipmentCountAll,
      estimatedActivePowerFactorAll,
      estimatedActiveLoadTires,
      estimatedReactiveLoadTires,
      fullPowerTires,
      estimatedGroupCurrentTires,
    });
  }

  function displayResults(results) {
    document.getElementById('group-utilization-rate').textContent =
      results.groupUtilizationRate.toFixed(4);
    document.getElementById('effective-equipment-count').textContent =
      results.effectiveEquipmentCount.toFixed(2);
    document.getElementById('estimated-active-power-factor').textContent =
      results.estimatedActivePowerFactor.toFixed(2);
    document.getElementById('estimated-active-load').textContent =
      results.estimatedActiveLoad.toFixed(2);
    document.getElementById('estimated-reactive-load').textContent =
      results.estimatedReactiveLoad.toFixed(2);
    document.getElementById('full-power').textContent =
      results.fullPower.toFixed(2);
    document.getElementById('estimated-group-current').textContent =
      results.estimatedGroupCurrent.toFixed(2);

    document.getElementById('utilization-rate-all').textContent =
      results.utilizationRateAll.toFixed(2);
    document.getElementById('effective-equipment-count-all').textContent =
      results.effectiveEquipmentCountAll;
    document.getElementById('estimated-active-power-factor-all').textContent =
      results.estimatedActivePowerFactorAll.toFixed(2);
    document.getElementById('estimated-active-load-tires').textContent =
      results.estimatedActiveLoadTires.toFixed(2);
    document.getElementById('estimated-reactive-load-tires').textContent =
      results.estimatedReactiveLoadTires.toFixed(2);
    document.getElementById('full-power-tires').textContent =
      results.fullPowerTires.toFixed(2);
    document.getElementById('estimated-group-current-tires').textContent =
      results.estimatedGroupCurrentTires.toFixed(2);

    resultsContainer.classList.remove('hidden');
  }

  // Заповнення початковими даними (опціонально)
  function fillWithSampleData() {
    const sampleEquipment = [
      {
        name: 'Шліфувальний верстат',
        efficiencyFactor: 0.92,
        loadPowerFactor: 0.9,
        loadVoltage: 0.38,
        quantity: 4,
        ratedPower: 20,
        utilizationRate: 0.15,
        reactivePowerFactor: 1.33,
      },
      {
        name: 'Свердлильний верстат',
        efficiencyFactor: 0.92,
        loadPowerFactor: 0.9,
        loadVoltage: 0.38,
        quantity: 2,
        ratedPower: 14,
        utilizationRate: 0.12,
        reactivePowerFactor: 1.0,
      },
      {
        name: 'Фугувальний верстат',
        efficiencyFactor: 0.92,
        loadPowerFactor: 0.9,
        loadVoltage: 0.38,
        quantity: 4,
        ratedPower: 42,
        utilizationRate: 0.15,
        reactivePowerFactor: 1.33,
      },
      {
        name: 'Циркулярна пила',
        efficiencyFactor: 0.92,
        loadPowerFactor: 0.9,
        loadVoltage: 0.38,
        quantity: 1,
        ratedPower: 20,
        utilizationRate: 0.5,
        reactivePowerFactor: 0.75,
      },
      {
        name: 'Прес',
        efficiencyFactor: 0.92,
        loadPowerFactor: 0.9,
        loadVoltage: 0.38,
        quantity: 1,
        ratedPower: 20,
        utilizationRate: 0.5,
        reactivePowerFactor: 0.75,
      },
      {
        name: 'Полірувальний верстат',
        efficiencyFactor: 0.92,
        loadPowerFactor: 0.9,
        loadVoltage: 0.38,
        quantity: 1,
        ratedPower: 40,
        utilizationRate: 0.2,
        reactivePowerFactor: 1.0,
      },
      {
        name: 'Фрезерний верстат',
        efficiencyFactor: 0.92,
        loadPowerFactor: 0.9,
        loadVoltage: 0.38,
        quantity: 2,
        ratedPower: 32,
        utilizationRate: 0.2,
        reactivePowerFactor: 1.0,
      },
      {
        name: 'Вентилятор',
        efficiencyFactor: 0.92,
        loadPowerFactor: 0.9,
        loadVoltage: 0.38,
        quantity: 1,
        ratedPower: 20,
        utilizationRate: 0.65,
        reactivePowerFactor: 0.75,
      },
    ];

    const equipmentCards = document.querySelectorAll('.equipment-card');

    for (
      let i = 0;
      i < Math.min(sampleEquipment.length, equipmentCards.length);
      i++
    ) {
      const card = equipmentCards[i];
      const data = sampleEquipment[i];

      card.querySelector('.equipment-name').value = data.name;
      card.querySelector('.efficiency-factor').value = data.efficiencyFactor;
      card.querySelector('.load-power-factor').value = data.loadPowerFactor;
      card.querySelector('.load-voltage').value = data.loadVoltage;
      card.querySelector('.quantity').value = data.quantity;
      card.querySelector('.rated-power').value = data.ratedPower;
      card.querySelector('.utilization-rate').value = data.utilizationRate;
      card.querySelector('.reactive-power-factor').value =
        data.reactivePowerFactor;

      calculateEquipmentValues(card);
    }
  }

  // Опційно викликати fillWithSampleData() для заповнення тестовими даними
  fillWithSampleData();
});
