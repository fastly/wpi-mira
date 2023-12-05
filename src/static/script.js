function addData(chart, allFreqMap) {
    const newData = Object.values(allFreqMap)
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


function addLabels(allOutliers) {
  const outlierList = allOutliers.AllOutliers

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

const data = {
    "2023-12-04T18:44:00-05:00": 1,
    "2023-12-04T18:45:00-05:00": 8,
    "2023-12-04T18:46:00-05:00": 2,
};
const float64Values = Object.values(data);
console.log(float64Values);
addData(chart, float64Values)
addLabels(float64Values)



setInterval(() => {
    fetch('http://localhost:8080/data')
        .then(response => response.json())
        .then(data => {
            const allFreqMap = data.AllFreq
            const allOutliers = data.AllOutliers
            addData(chart, allFreqMap);
        })
        .catch(error => {
            console.error('Error fetching data:', error);
        });
}, 3000);

