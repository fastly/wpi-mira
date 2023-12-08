/*
TODO: potentially link the url and the data for the prefix in one struct
 */
// Function to populate the dropdown with values from a list of strings
function populateDropdown(listOfStrings) {
    const select = document.getElementById("subscriptionSelect");
    const existingOptions = Array.from(select.options).map(option => option.textContent.toLowerCase());

    listOfStrings.forEach(function(string) {
        const lowercaseString = string.toLowerCase();
        if (!existingOptions.includes(lowercaseString)) {
            const option = document.createElement("option");
            option.value = lowercaseString.replace(/\s+/g, '-'); // Example conversion to lowercase and replacing spaces with hyphens
            option.textContent = string;
            select.appendChild(option);
        }
    });
}

function getEndpointForSub (query) { //take in a subscription string and turn it into a link for the end point
    const baseUrl = 'http://localhost:8080/frequencies';
    // Encode the query object as a URI component
    const encodedQuery = encodeURIComponent(JSON.stringify(query));

    // Construct the final URL with the encoded query as a parameter
    const finalUrl = `${baseUrl}?subscription=${query}`;
    console.log(finalUrl)

    return finalUrl;
}


//opens up a page based on the prefix selected
function openPage() {

    //the items are the filters
    const select = document.getElementById("subscriptionSelect");
    const selectedOption = select.options[select.selectedIndex];
    const selectedUrl = selectedOption.value;

    //for every item that is selected on the page; pop up the data for frequencies and counts from the given subscription endpoint

    const url = getEndpointForSub(selectedUrl);
    fetchByUrl(url)
        .then(data => {
            if (data) {
                console.log('Fetched data:', data);
                //add subscriptions to the dropdown as they populate in the result

                const filters = Object.keys(data) //create urls based localhost:8080/filter
                const counts = Object.values(data)

                const newPage = window.open('mainData.html', '_blank');
                if (newPage) {
                    newPage.document.close();
                } else {
                    alert('Pop-up blocked! Please allow pop-ups for this website.');
                }

            } else {
                console.log('No data fetched');
            }
        });
    /*    const select = document.getElementById("subscriptionSelect");
        const selectedOption = select.options[select.selectedIndex];
        const selectedUrl = selectedOption.value;
        window.open(selectedUrl, '_blank');
        // window.location.href = selectedUrl;

     if (selectedUrl) {
            fetchByUrl("http://localhost:8080/data")
                .then(data => {
                    if (data) {
                        console.log('Fetched data:', data);
                        // Process fetched data as needed (e.g., displaying in UI)
                        // Example: addData(chart, data);
                        // Example: addLabels(data);
                    } else {
                        console.log('No data fetched');
                    }
                });
        }*/
}

/*
function openPage() {
    const select = document.getElementById("subscriptionSelect");
    const selectedUrl = select.value;
    window.location.href = selectedUrl;
}*/

function createChart(){
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
            maintainAspectRatio: true,
            aspectRatio: 0.5,
            scales: {
                y: {
                    beginAtZero: true
                }
            }
        }
    });


}


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

//works with a box right now; need to make this more readable
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

                populateDropdown(filters)


                //populateDropdown(filters);

                //get the results by keys
                //addData(chart, firstResult); //getting data
                //addLabels(firstResult);
            } else {
                console.log('No data fetched');
            }
        });
}, 3000); //updates every 3 seconds

