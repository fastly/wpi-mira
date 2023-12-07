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


function addData(chart, result) {
    const allFreqMap = result.AllFreq

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


function addLabels(result) {
    const allOutliers = result.AllOutliers
    const outlierList = Object.values(allOutliers)

    const timestampList = [];
    const valsList = [];

    for (let i = 0; i < outlierList.length; i++) {
        timestamp = new Date(outlierList[i].Timestamp);
        count =outlierList[i].Count;
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


async function fetchByUrl(url) {
    try {
        const response = await fetch(url);

        if (!response.ok) {
            throw new Error('Network response was not ok.');
        }
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching data:', error);
        return null;
    }
}


/*//open a new page for every subscription
function openPage() {
    var sub = document.getElementById("subscriptionSelect").value;
    var url = "index.html";
    var newWindow = window.open(url, "_blank");
}



function populateDropdown(itemsToAdd) {
    const dropdown = document.getElementById("subscriptionSelect");

    itemsToAdd.forEach(item => {
        // Check if the item already exists in the dropdown
        const exists = [...dropdown.options].some(option => option.text === item);
        if (!exists) {
            // Create a new option element
            const newOption = document.createElement("option");
            newOption.value = item; // Set value (you can change this if needed)
            newOption.text = item; // Set text

            // Append the new option to the dropdown
            dropdown.appendChild(newOption);
        }
    });
}*/



// Function to populate the table with data
function populateTable(data) {
    const tableBody = document.getElementById('tableBody');

    // Clear any existing rows
    tableBody.innerHTML = '';

    // Loop through the data and create table rows
    data.forEach(item => {
        const row = document.createElement('tr');


        const timeCell = document.createElement('td');
        timeCell.textContent = item.time;
        row.appendChild(timeCell);

        const countsCell = document.createElement('td');
        countsCell.textContent = item.counts;
        row.appendChild(countsCell);

        tableBody.appendChild(row);
    });
}

setInterval(() => {
    const url = 'http://localhost:8080/data';
    fetchByUrl(url)
        .then(data => {
            if (data) {
                console.log('Fetched data:', data);
                //add subscriptions to the dropdown as they populate in the results
                const filters = Object.keys(data) //create urls based localhost:8080/filter
                const results = Object.values(data)

                const firstFilter = filters[0]
                const firstResult = results[0]

                //populateDropdown(filters);

                //get the results by keys
                addData(chart, firstResult); //getting data
                addLabels(firstResult);
            } else {
                console.log('No data fetched');
            }
        });



   // data = fetchData('http://localhost:8080/data')
   /* fetch('http://localhost:8080/data')
        .then(response => response.json())
        .then(data => {
            //add subscriptions to the dropdown as they populate in the results
            const filters = Object.keys(data) //create urls based localhost:8080/filter
            const results = Object.values(data)

            const firstFilter = filters[0]
            const firstResult = results[0]

            //populateDropdown(filters);

            //get the results by keys
            addData(chart, firstResult); //getting data
            addLabels(firstResult);

        })
        .catch(error => {
            console.error('Error fetching data:', error);
        });*/
}, 3000); //updates every 3 seconds

