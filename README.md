# Read Me

# How to run:
To provide a URL containing bz2 files, go into config.json and insert the url in the same format as the default under staticFilesLink parameter
To download static data for testing, first cd to src and run:
go run static_data/get_static_data.go

In order to run the main program cd into the /src folder and run the following commands:
go mod tidy

go run main.go -config="path_to_config_json"
(if no config file, then default uses config.json)