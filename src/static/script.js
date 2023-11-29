function addData(chart, newData) {
    if (newData.length === 0) {
        return; // No new data to add
    }

    // Clear existing data
    chart.data.labels = [];
    chart.data.datasets[0].data = [];

    // Add new data
    for (let i = 0; i < newData.length; i++) {
        chart.data.labels.push(`Point ${i + 1}`);
        chart.data.datasets[0].data.push(newData[i]);

        const numberElement = document.createElement('p');
        const numberText = document.createTextNode(`Number ${i + 1}: ${newData[i]}`);
        numberElement.appendChild(numberText);
    }

    chart.update();
}


const ctx = document.getElementById('myChart').getContext('2d');
const chart = new Chart(ctx, {
    type: 'line',
    data: {
        labels: [],
        datasets: [{
            label: 'Message Counts Per Minute',
            data: [],
            backgroundColor: 'rgba(54, 162, 235, 0.2)',
            borderColor: 'rgba(54, 162, 235, 1)',
            borderWidth: 1
        }]
    },
    options: {
        responsive: true,
        maintainAspectRatio: false,
        aspectRatio: 0.5,
        scales: {
            y: {
                beginAtZero: true
            }
        }
    }
});



setInterval(() => {
    fetch('http://localhost:8080/data')
        .then(response => response.json())
        .then(data => {
            const frequencies = data.map(result => result.Frequencies).flat();
            addData(chart, frequencies);
        })
        .catch(error => {
            console.error('Error fetching data:', error);
        });
}, 3000);

