# MIRA: Modelling Internet Routing Anomalies

A project done in collaboration between [Fastly](https://www.fastly.com/) and [WPI](https://www.wpi.edu/).
This project is a WPI [Major Qualifying Project](https://www.wpi.edu/project-based-learning/project-based-education/major-qualifying-project).

This tool monitors public Internet routing data and identifies anomalous routing events based on frequency of updates.

The project was inspired by
[ShakeAlert](https://labs.ripe.net/author/marcel-flores/detecting-waves-with-shakealert/).
It makes use of [RIS Live](https://ris-live.ripe.net/) and [Route Views](https://routeviews.org/) static data for its data sources.


# How to run:
In order to run the main program, cd into the /src folder and run the following commands:

```
go mod tidy
go run main.go -config="path_to_config_json"
```

If no config file, then default uses `default-config.json`, which listens to the Fastly ASN (54113)

# Config file format
Create a configuration json file based on the Configuration struct in config.go. For examples of various use cases, refer to example configurations contained in [example_configs](src/example_configs).

# How to download static file for testing:
To provide a URL containing bz2 files, create a config.json file and insert the URL under the staticFilesLink parameter in the same format as the following: "http://routeviews.org/route-views.ny/bgpdata/2021.11/UPDATES/"

To download the static data from the link, cd into the /src folder and run the following command:
`go run static_data/get_static_data.go`
