document.addEventListener('DOMContentLoaded', function () {
  const form = document.getElementById('calculatorForm');
  const resultsDiv = document.getElementById('results');

  form.addEventListener('submit', function (e) {
    e.preventDefault();

    const composition = {
      carbonCombustible:
        parseFloat(document.getElementById('carbonCombustible').value) || 0,
      hydrogenCombustible:
        parseFloat(document.getElementById('hydrogenCombustible').value) || 0,
      oxygenCombustible:
        parseFloat(document.getElementById('oxygenCombustible').value) || 0,
      sulfurCombustible:
        parseFloat(document.getElementById('sulfurCombustible').value) || 0,
      vanadiumCombustible:
        parseFloat(document.getElementById('vanadiumCombustible').value) || 0,
      moistureContent:
        parseFloat(document.getElementById('moistureContent').value) || 0,
      ashDry: parseFloat(document.getElementById('ashDry').value) || 0,
      heatingValueCombustible:
        parseFloat(document.getElementById('heatingValueCombustible').value) ||
        0,
    };

    fetch('/calculator2', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(composition),
    })
      .then(response => {
        if (!response.ok) {
          throw new Error('Помилка сервера');
        }
        return response.json();
      })
      .then(data => {
        displayResults(data);
      })
      .catch(error => {
        alert('Помилка: ' + error.message);
      });
  });

  function displayResults(results) {
    // склад робочої маси
    const workingCompositionDiv = document.getElementById('workingComposition');
    workingCompositionDiv.innerHTML = '';

    Object.entries(results.workingComposition).forEach(([component, value]) => {
      const row = document.createElement('div');
      row.className = 'row mb-2';

      const labelCol = document.createElement('div');
      labelCol.className = 'col-8';
      labelCol.textContent = component;

      const valueCol = document.createElement('div');
      valueCol.className = 'col-4 text-end';

      // для ванадію використовуємо мг/кг, для інших - відсотки
      if (component.includes('V')) {
        valueCol.textContent = `${value.toFixed(2)} мг/кг`;
      } else {
        valueCol.textContent = `${value.toFixed(2)}%`;
      }

      row.appendChild(labelCol);
      row.appendChild(valueCol);
      workingCompositionDiv.appendChild(row);
    });

    // нижча теплота згоряння
    const heatingValueDiv = document.getElementById('workingHeatingValue');
    heatingValueDiv.textContent = `${results.workingHeatingValue.toFixed(
      2,
    )} МДж/кг`;

    resultsDiv.style.display = 'block';
  }
});
