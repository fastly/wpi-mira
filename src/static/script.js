function addData(chart, newData) {
    /*const maxDataPoints = 400;

    if (newData.length === 0) {
        return; // No new data to add
    }

    const currentDataLength = chart.data.labels.length;

    if (currentDataLength >= maxDataPoints) {
        const dataToRemove = currentDataLength - maxDataPoints + newData.length;

        for (let i = 0; i < dataToRemove; i++) {
            chart.data.labels.shift();
            chart.data.datasets[0].data.shift();
        }
    }

    for (let i = 0; i < newData.length; i++) {
        chart.data.labels.push(`Point ${currentDataLength + i + 1}`);
        chart.data.datasets[0].data.push(newData[i]);
    }

    chart.update();*/
    const maxDataPoints = 400; //should not need this after the results array gets correctly modified

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
    }

    chart.update();
}


function getFrequencies() {
    return Math.floor(Math.random() * 100);
}

const ctx = document.getElementById('myChart').getContext('2d');
const chart = new Chart(ctx, {
    type: 'line',
    data: {
        labels: ['Point 1', 'Point 2', 'Point 3', 'Point 4', 'Point 5'],
        datasets: [{
            label: 'Sample Data',
            data: [20, 40, 30, 50, 25],
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

/*setInterval(() => {
    /!*const newDataPoints = [];

    // Generate multiple random data points and add them
    for (let i = 0; i < 5; i++) {
        const newDataPoint = generateRandomDataPoint();
        newDataPoints.push(newDataPoint);
    }*!/
    addData(chart, frequencies);
}, 3000)*/
