export let numbersGiven; //this works sometimes?
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
export function openPage() {

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

                const filters = Object.keys(data); //create urls based localhost:8080/filter
                const counts = Object.values(data);
                numbersGiven = [1,2,3,4]
                console.log(numbersGiven)


                const newPage = window.open('mainData.html', '_blank');


                if (newPage) {
                    // Wait for the new page to load
                    newPage.onload = function() {
                        // Access the document within the new window
                        const newDoc = newPage.document;
                        // Find an element in the new window and add data
                        const newData =  counts.join(', '); // Your new data
                        // For example, let's find a div with id 'dataContainer' and add the new data
                        const dataContainer = newDoc.getElementById('outliers');
                        if (dataContainer) {
                            dataContainer.innerHTML = newData;
                        } else {
                            console.error('Element not found in new window');
                        }
                    };

                    //newPage.document.close();
                } else {
                    alert('Pop-up blocked! Please allow pop-ups for this website.');
                }

            } else {
                console.log('No data fetched');
            }
        });
}

document.getElementById('goButton').addEventListener('click', openPage);

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
                /*const results = Object.values(data)

                const firstFilter = filters[0]
                const firstResult = results[0]*/

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