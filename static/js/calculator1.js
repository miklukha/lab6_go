document.addEventListener('DOMContentLoaded', function () {
  const form = document.getElementById('calculatorForm');
  const resultsDiv = document.getElementById('results');

  form.addEventListener('submit', function (e) {
    e.preventDefault();

    const composition = {
      hp: parseFloat(document.getElementById('hp').value) || 0,
      cp: parseFloat(document.getElementById('cp').value) || 0,
      sp: parseFloat(document.getElementById('sp').value) || 0,
      np: parseFloat(document.getElementById('np').value) || 0,
      op: parseFloat(document.getElementById('op').value) || 0,
      wp: parseFloat(document.getElementById('wp').value) || 0,
      ap: parseFloat(document.getElementById('ap').value) || 0,
    };

    fetch('/calculator1', {
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
    // коефіцієнти переходу
    document.getElementById('dryMassCoefficient').textContent =
      results.dryMassCoefficient.toFixed(2);
    document.getElementById('combustibleMassCoefficient').textContent =
      results.combustibleMassCoefficient.toFixed(2);

    // суха маса
    const dryCompositionDiv = document.getElementById('dryComposition');
    dryCompositionDiv.innerHTML = '';
    Object.entries(results.dryComposition).forEach(([key, value]) => {
      dryCompositionDiv.innerHTML += `
                <div class="row">
                    <div class="col-6">${key}</div>
                    <div class="col-6 text-end">${value.toFixed(2)}</div>
                </div>
            `;
    });

    // горюча маса
    const combustibleCompositionDiv = document.getElementById(
      'combustibleComposition',
    );
    combustibleCompositionDiv.innerHTML = '';
    Object.entries(results.combustibleComposition).forEach(([key, value]) => {
      combustibleCompositionDiv.innerHTML += `
                <div class="row">
                    <div class="col-6">${key}</div>
                    <div class="col-6 text-end">${value.toFixed(2)}</div>
                </div>
            `;
    });

    // нижча теплота згоряння
    document.getElementById('lowerHeatingValue').textContent =
      results.lowerHeatingValue.toFixed(2);
    document.getElementById('lowerDryHeatingValue').textContent =
      results.lowerDryHeatingValue.toFixed(2);
    document.getElementById('lowerCombustibleHeatingValue').textContent =
      results.lowerCombustibleHeatingValue.toFixed(2);

    resultsDiv.style.display = 'block';
  }
});
