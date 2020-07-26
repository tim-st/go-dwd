package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {

	var targetPath string
	var resolutionString string
	var searchString string
	var latitude float64
	var longitude float64
	var limit int
	var descending bool
	var minYears int
	var minHeight int
	var maxHeight int
	var sortByYears bool
	var sortByHeight bool
	var sortByDistance bool
	var download bool
	var stationID int

	const defaultMinHeight = -100000
	const defaultMaxHeight = +100000

	flag.StringVar(&targetPath, "targetPath", "", "Path where the output files should go")
	flag.StringVar(&resolutionString, "resolution", "daily", `resolution of data updates in {"daily", "hourly", "10minutes"}`)
	flag.IntVar(&stationID, "stationID", -1, "StationID of the weather station to show")
	flag.StringVar(&searchString, "search", "", "Search a weather station by name or place")
	flag.Float64Var(&latitude, "latitude", 0, "Latitude of the target place")
	flag.Float64Var(&longitude, "longitude", 0, "Longitude of the target place")
	flag.IntVar(&limit, "limit", -1, "Limit the number of weather stations")
	flag.IntVar(&minYears, "minYears", 0, "Minimum number of years of data the weather station should have")
	flag.IntVar(&minHeight, "minHeight", defaultMinHeight, "Minimum height of the weather station in meters")
	flag.IntVar(&maxHeight, "maxHeight", defaultMaxHeight, "Maximum height of the weather station in meters")
	flag.BoolVar(&sortByYears, "sortByYears", false, "Sort the weather stations by number of years of data")
	flag.BoolVar(&sortByHeight, "sortByHeight", false, "Sort the weather stations by height")
	flag.BoolVar(&sortByDistance, "sortByDistance", false, "Sort the weather stations by their great-circle distance to given latitude and longitude")
	flag.BoolVar(&descending, "descending", false, "Sets the sort order of the weather stations to descending")
	flag.BoolVar(&download, "download", false, "If true, downloads weather data and creates a CSV file for each station")

	flag.Parse()

	var r resolution
	switch resolutionString {
	case "daily":
		r = resolutionDaily
	case "hourly":
		r = resolutionHourly
	case "10minutes", "10_minutes":
		r = resolution10Minutes
	default:
		log.Fatal("dwd: unknown resolution")
	}

	stations, err := stations(r)
	if err != nil {
		log.Fatal(err)
	}

	if stationID > 0 {
		var s []station
		for _, station := range stations {
			if station.id == stationID {
				s = append(s, station)
				break
			}
		}
		stations = s
	}

	if stationID <= 0 && minYears > 0 {
		var s []station
		for _, station := range stations {
			if int(station.dateLast.Sub(station.dateFirst).Hours())/24/365 >= minYears {
				s = append(s, station)
			}
		}
		stations = s
	}

	if stationID <= 0 && minHeight > defaultMinHeight {
		var s []station
		for _, station := range stations {
			if station.height >= minHeight {
				s = append(s, station)
			}
		}
		stations = s
		if limit > 0 && limit < len(stations) {
			stations = stations[:limit]
		}
	}

	if stationID <= 0 && maxHeight < defaultMaxHeight {
		var s []station
		for _, station := range stations {
			if station.height <= maxHeight {
				s = append(s, station)
			}
		}
		stations = s
		if limit > 0 && limit < len(stations) {
			stations = stations[:limit]
		}
	}

	searchString = strings.TrimSpace(searchString)

	if stationID <= 0 && len(searchString) > 0 {
		sortStationsByStringDistance(stations, searchString)
		if limit > 0 && limit < len(stations) {
			stations = stations[:limit]
		}
	}

	if sortByYears {
		sortStationsByYearsAscending(stations)
	}
	if sortByHeight {
		sortStationsByHeightAscending(stations)
	}
	if sortByDistance && latitude > 0 && longitude > 0 {
		sortStationsByDistance(stations, latitude, longitude)
	}
	if descending {
		for i, j := 0, len(stations)-1; i < j; i, j = i+1, j-1 {
			stations[i], stations[j] = stations[j], stations[i]
		}
	}
	if stationID <= 0 && limit > 0 && limit < len(stations) {
		stations = stations[:limit]
	}

	for _, station := range stations {
		fmt.Printf("% 6d %4d–%4d % 5dm %6.3f %6.3f %s (%s)\n",
			station.id, station.dateFirst.Year(), station.dateLast.Year(), station.height,
			station.latitude, station.longitude, station.name, station.state)
	}

	if download {
		err = downloadStations(stations, r, targetPath)
		if err != nil {
			log.Fatal(err)
		}
	}

}
