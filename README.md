# go-dwd

Inofficial client for downloading weather data from [DWD](https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/) into single CSV files.

## Terms of use

<https://opendata.dwd.de/climate_environment/CDC/Terms_of_use.txt>

## Installation

>go get -u github.com/tim-st/go-dwd/...

The compiled binary `go-dwd` will be located at `$GOPATH\bin\`

## Usage

>go-dwd -help
```
Usage of go-dwd:
  -descending
        Sets the sort order of the weather stations to descending
  -download
        If true, downloads weather data and creates a CSV file for each station
  -latitude float
        Latitude of the target place
  -limit int
        Limit the number of weather stations (default -1)
  -longitude float
        Longitude of the target place
  -maxHeight int
        Maximum height of the weather station in meters (default 100000)
  -minHeight int
        Minimum height of the weather station in meters (default -100000)
  -minYears int
        Minimum number of years of data the weather station should have
  -resolution string
        resolution of data updates in {"daily", "hourly", "10minutes"} (default "daily")
  -search string
        Search a weather station by name or place
  -sortByDistance
        Sort the weather stations by their great-circle distance to given latitude and longitude
  -sortByHeight
        Sort the weather stations by height
  -sortByYears
        Sort the weather stations by number of years of data
  -stationID int
        StationID of the weather station to show (default -1)
  -targetPath string
        Path where the output files should go
```

### List weather stations

>go-dwd -search=Place

### List 10 weather stations

>go-dwd -search=Place -limit=10

### Download weather data

The following command will download weather data for every station in the results with `daily` resolution
>go-dwd -search=Place -limit=10 -resolution=daily -download

## Example

An example file showing how to import the created CSV file using `pandas` is [here](https://github.com/tim-st/go-dwd/blob/master/beispiel_tagesdaten.ipynb).
