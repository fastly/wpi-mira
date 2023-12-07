# Read Me

# How to download static file for testing:
To provide a URL containing bz2 files, create a config.json file and insert the URL under the staticFilesLink parameter in the same format as the following: "http://routeviews.org/route-views.ny/bgpdata/2021.11/UPDATES/"

To download the static data from the link, cd into the /src folder and run the following command:
go run static_data/get_static_data.go

# How to run:
In order to run the main program, cd into the /src folder and run the following commands:
go mod tidy

go run main.go -config="path_to_config_json"

If no config file, then default uses default-config.json, which listens to the Fastly ASN (54113)

# Config file format
Create a configuration json file based on the Configuration struct in config.go. For examples of various use cases, refer to the example_configs folder. 