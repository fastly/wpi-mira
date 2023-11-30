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

function addLabels(madOutliers, madOutliersTimes, shakeOutliers, shakeOutliersTimes) {
    //mad outliers
    const madOutliersContainer = document.getElementById('madOutliers');
    const formattedArrayMad = madOutliers.join(', '); // Format the array for display
    madOutliersContainer.innerText = `MAD Outliers: [${formattedArrayMad}]`;

    //mad outlier timestamps
    const madOutliersTimesContainer = document.getElementById('madOutliersTimes');
    const formattedArrayMadTimes = madOutliersTimes.join(', '); // Format the array for display
    madOutliersTimesContainer.innerText = `MAD Outliers: [${formattedArrayMadTimes}]`;

    //add shake alert outlier counts
    const shakeOutliersContainer = document.getElementById('shakeOutliers');
    const formattedArrayShake = shakeOutliers.join(', '); // Format the array for display
    shakeOutliersContainer.innerText = `ShakeAlertOutliers: [${formattedArrayShake}]`;

    //shake times
    const shakeOutliersTimesContainer = document.getElementById('shakeOutliersTimes');
    const formattedArrayShakeTimes = shakeOutliersTimes.join(', '); // Format the array for display
    shakeOutliersTimesContainer.innerText = `ShakeAlertOutliers: [${formattedArrayShakeTimes}]`;
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

            const madOutliers = data.map(result => result.MADOutliers).flat();
            const madOutliersTimes = data.map(result => result.MADTimestamps).flat();

            const shakeAlertOutliers = data.map(result => result.ShakeAlertOutliers).flat();
            const shakeAlertOutliersTimes = data.map(result => result.ShakeAlertTimestamps).flat();

            addData(chart, frequencies);
            addLabels(madOutliers, madOutliersTimes, shakeAlertOutliers, shakeAlertOutliersTimes)
        })
        .catch(error => {
            console.error('Error fetching data:', error);
        });
}, 3000);

