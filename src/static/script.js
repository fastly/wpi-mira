function addData(chart, allFreqMap) {
    const newData = Object.values(allFreqMap) //frequencies at each individual time stamp
    const timeStamps = Object.keys(allFreqMap) //timestamps for labeling points
    if (newData.length === 0) {
        return; // No new data to add
    }

    // Clear existing data
    chart.data.labels = [];
    chart.data.datasets[0].data = [];

    // Add new data
    for (let i = 0; i < newData.length; i++) {
        chart.data.labels.push(timeStamps[i]);
        chart.data.datasets[0].data.push(newData[i]);
    }
    chart.update();
}


function addLabels(allOutliers) {
    const outlierList = Object.values(allOutliers)

    const timestampList = [];
    const valsList = [];

    for (let i = 0; i < outlierList.length; i++) {
        timestamp = new Date(outlierList[i].Timestamp);
        count = outlierList[i].Count;
        timestampList.push(timestamp)
        valsList.push(count)
    }

    //outliers
    const outliersContainer = document.getElementById('outliers');
    const formattedOutliers = valsList.join(', '); // Format the array for display
    outliersContainer.innerText = `Outliers: [${formattedOutliers}]`;

    //timestamps
    const outliersTimesContainer = document.getElementById('outliersTimes');
    const formattedOutliersTimes = timestampList.join(', '); // Format the array for display
    outliersTimesContainer.innerText = `Outliers Times: [${formattedOutliersTimes}]`;



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
            const allFreqMap = data.AllFreq
            const allOutliers = data.AllOutliers
            addData(chart, allFreqMap);
            addLabels(allOutliers);

        })
        .catch(error => {
            console.error('Error fetching data:', error);
        });
}, 3000);

