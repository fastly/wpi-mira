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
    const finalUrl = `${baseUrl}?subscription=${query}`;
    console.log(finalUrl)

    return finalUrl;
}


//opens a page with all the data on it
function openPage() {
    //the items are the filters
    const select = document.getElementById("subscriptionSelect");
    const selectedOption = select.options[select.selectedIndex];
    const selectedUrl = selectedOption.value;
    getData(selectedUrl);
    const newPage = window.open('index.html', '_blank');
}

function addData(chart, newData) {
    if (newData.length === 0) {
        return; // No new data to add
    }

    // Clear existing data
    chart.data.labels = [];
    chart.data.datasets[0].data = [];

    // Add new data
    for (let i = 0; i < newData.length; i++) {
        chart.data.labels.push(i);
        chart.data.datasets[0].data.push(newData[i]);
    }
    chart.update();
}

//opens up a page based on the prefix selected
function getData(selectedUrl) {
    //for every item that is selected on the page; pop up the data for frequencies and counts from the given subscription endpoint

    const url = getEndpointForSub(selectedUrl);
    fetchByUrl(url)
        .then(data => {
            if (data) {
                console.log('Fetched data:', data);
                //add subscriptions to the dropdown as they populate in the result

                const filters = Object.keys(data); //create urls based localhost:8080/filter
                const counts = Object.values(data);
                let numbersGiven = [1,2,3,4];
                addData(chart, numbersGiven)

            } else {
                console.log('No data fetched');
            }
        });
}



//document.getElementById('goButton').addEventListener('click', openPage);

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
    const url = 'http://localhost:8080/data';


    fetchByUrl(url)
        .then(data => {
            if (data) {
                console.log('Fetched data:', data);
                //add subscriptions to the dropdown as they populate in the results
                const filters = Object.keys(data) //create urls based localhost:8080/filter
                populateDropdown(filters);
            } else {
                console.log('No data fetched');
            }
        });
}, 3000); //updates every 3 seconds