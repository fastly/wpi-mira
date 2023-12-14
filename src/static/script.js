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

function getOutliersEndpoint (query) {
    const baseUrl = 'http://localhost:8080/outliers';
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
    let url = new URL(window.location.href);
    url.searchParams.set("p", selectedUrl);
    window.open(url, '_blank');
}

//adding outlier messages
function addMsgs (outlierUrl) {
    //getting outliers from the outliers endpoint
    fetchByUrl(outlierUrl)
        .then(data => {
            let textboxContent = '';
            //go through all the timestamps and only get the message contents
            // Iterate through each timestamp (assuming each timestamp is a key in the JSON object)
            console.log(data);
            for (const timestamp in data) {
                if (data.hasOwnProperty(timestamp)) {
                    const dataForTimestamp = data[timestamp];
                    const counts = dataForTimestamp.Count;
                    const messages = dataForTimestamp.Msgs.map(msg => {
                        return `Timestamp: ${msg.Timestamp}, Type: ${msg.BGPMessageType}, PeerIP: ${msg.PeerIP}, ASN: ${msg.PeerASN}, Prefix: ${msg.Prefix}`;
                    }).join('\n');

                    // Append the information to the textbox content
                    textboxContent += `Number of Messages in the Outlier Container: ${counts}\n Outlier Messages:\n${messages}\n\n`;
                    const msgsContainer = document.getElementById('msgs');
                    // Set the value of the textbox with the extracted data
                    msgsContainer.innerText = textboxContent;
                } else {
                    console.log('No data fetched');
            }
            }
        });

}

//adding data to the graph
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
function getData(chart, selectedUrl) {
    //for every item that is selected on the page; pop up the data for frequencies and counts from the given subscription endpoint


    //selectedUrl is the specific filter/subscription
    const url = getEndpointForSub(selectedUrl);
    const outlierUrl = getOutliersEndpoint(selectedUrl);
    //console.log(url);

    fetchByUrl(url)
        .then(data => {
            if (data) {
                //console.log(filters);
                const counts = Object.values(data);
                //console.log(selectedUrl);
                addData(chart, counts);
                addMsgs(outlierUrl);

            } else {
                console.log('No data fetched');
            }
        });
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

function onLoad() {
    setInterval(() => {
        const url = 'http://localhost:8080/data';
        fetchByUrl(url)
            .then(data => {
                if (data) {
                    //console.log('Fetched data:', data);
                    //add subscriptions to the dropdown as they populate in the results
                    const filters = Object.keys(data) //create urls based localhost:8080/filter
                    populateDropdown(filters);

                } else {
                    console.log('No data fetched');
                }
            });
    }, 3000); //updates every 3 seconds

    //find query parameter
    const urlSearchParams = new URLSearchParams(window.location.search);
    const subs = urlSearchParams.get("p");
    if (typeof subs === "string" && subs.length !== 0) {
        const ctx = document.getElementById('myChart').getContext('2d');
        ctx.width = window.innerWidth * 0.5; // 50% of the viewport width
        ctx.height = window.innerHeight; // 100% of the viewport height
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
                aspectRatio: 1,
                scales: {
                    y: {
                        beginAtZero: true
                    }
                }
            }
        });

        setInterval(() => {
            getData(chart, subs);
            }, 3000);

    }


}


