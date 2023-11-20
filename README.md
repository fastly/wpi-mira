# Read Me

# How to run:
If you want to download static data for testing, first cd to src/static_data and run:
go run get_static_data.go

In order to download bz2 files from a URL modify the link containing bz2 files to download in config.json following the format of the default link in the file. Then, change the directory into /src and type:
go run static_data/get_static_fromURL.go

In order to run the main program cd into the /src folder and run the following commands:
go mod tidy
go run main.go
