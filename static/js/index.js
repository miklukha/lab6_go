document.addEventListener('DOMContentLoaded', function () {
  const form = document.getElementById('calculatorForm');
  const resultsDiv = document.getElementById('results');

  form.addEventListener('submit', function (e) {
    e.preventDefault();

    const minerals = {
      coal: parseFloat(document.getElementById('coal').value) || 0,
      mazut: parseFloat(document.getElementById('mazut').value) || 0,
      gas: parseFloat(document.getElementById('gas').value) || 0,
    };

    fetch('/calculator', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(minerals),
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
    // спалювання вугілля
    document.getElementById('coalEmissionFactor').textContent =
      results.coalEmissionFactor.toFixed(2) + ' г/ГДж';
    document.getElementById('coalEmissionValue').textContent =
      results.coalEmissionValue.toFixed(2) + ' т';

    // спалювання мазуту
    document.getElementById('mazutEmissionFactor').textContent =
      results.mazutEmissionFactor.toFixed(2) + ' г/ГДж';
    document.getElementById('mazutEmissionValue').textContent =
      results.mazutEmissionValue.toFixed(2) + ' т';

    // спалювання газу
    document.getElementById('gasEmissionFactor').textContent =
      results.gasEmissionFactor.toFixed(2) + ' г/ГДж';
    document.getElementById('gasEmissionValue').textContent =
      results.gasEmissionValue.toFixed(2) + ' т';

    resultsDiv.style.display = 'block';
  }
});
