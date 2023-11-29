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

function printNumber() {
    var number = 42; // You can set any number you want to display
    var div = document.createElement('div');
    div.style.position = 'fixed';
    div.style.bottom = '10px';
    div.style.left = '10px';
    div.style.backgroundColor = 'rgba(0, 0, 0, 0.5)';
    div.style.color = '#fff';
    div.style.padding = '5px';
    div.style.borderRadius = '5px';
    div.textContent = 'Number: ' + number;
    document.body.appendChild(div);
}



setInterval(() => {
    fetch('http://localhost:8080/data')
        .then(response => response.json())
        .then(data => {
            const frequencies = data.map(result => result.Frequencies).flat();
            addData(chart, frequencies);
            printNumber()
        })
        .catch(error => {
            console.error('Error fetching data:', error);
        });
}, 3000);

