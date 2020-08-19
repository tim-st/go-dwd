package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/tim-st/go-uniseg"
	"golang.org/x/text/encoding/charmap"
)

type station struct {
	id          int
	name        string
	dateFirst   time.Time
	dateLast    time.Time
	height      int
	latitude    float64
	longitude   float64
	state       string
	resolution  resolution
	zipURLs     map[string][]string
	readClosers map[string]io.ReadCloser
}

type stationlist []station

func sortStationsByID(s stationlist) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].id < s[j].id
	})
}

func sortStationsByDistance(s stationlist, latitude, longitude float64) {
	sort.Slice(s, func(i, j int) bool {
		return distanceInKiloMeters(s[i].latitude, s[i].longitude, latitude, longitude) < distanceInKiloMeters(s[j].latitude, s[j].longitude, latitude, longitude)
	})
}

func sortStationsByStringDistance(s stationlist, search string) {

	searchLower := strings.ToLower(search)
	sort.Slice(s, func(i, j int) bool {

		siNameLower := strings.ToLower(s[i].name)
		sjNameLower := strings.ToLower(s[j].name)

		ci1 := strings.Contains(siNameLower, searchLower)
		cj1 := strings.Contains(sjNameLower, searchLower)
		if ci1 && !cj1 {
			return true
		}
		if cj1 && !ci1 {
			return false
		}

		if ci1 && cj1 {

			di := int(^uint(0) >> 1)
			dj := di

			for _, segment := range uniseg.Segments([]byte(siNameLower)) {
				if segment.Category.IsWord() {
					d := editDistance(string(segment.Segment), searchLower)
					if d < di {
						di = d
					}
				}
			}

			for _, segment := range uniseg.Segments([]byte(sjNameLower)) {
				if segment.Category.IsWord() {
					d := editDistance(string(segment.Segment), searchLower)
					if d < dj {
						dj = d
					}
				}
			}

			if di == dj {
				return len(siNameLower) < len(sjNameLower)
			}

			return di < dj
		}

		siStateLower := strings.ToLower(s[i].state)
		sjStateLower := strings.ToLower(s[j].state)

		ci2 := strings.Contains(siStateLower, searchLower)
		cj2 := strings.Contains(sjStateLower, searchLower)
		if ci2 && !cj2 {
			return true
		}
		if cj2 && !ci2 {
			return false
		}

		if ci2 && cj2 {
			di := int(^uint(0) >> 1)
			dj := di

			for _, segment := range uniseg.Segments([]byte(siStateLower)) {
				if segment.Category.IsWord() {
					d := editDistance(string(segment.Segment), searchLower)
					if d < di {
						di = d
					}
				}
			}

			for _, segment := range uniseg.Segments([]byte(sjStateLower)) {
				if segment.Category.IsWord() {
					d := editDistance(string(segment.Segment), searchLower)
					if d < dj {
						dj = d
					}
				}
			}

			if di == dj {
				return len(siStateLower) < len(sjStateLower)
			}

			return di < dj
		}

		di := int(^uint(0) >> 1)
		dj := di

		for _, segment := range uniseg.Segments([]byte(siNameLower)) {
			if segment.Category.IsWord() {
				d := editDistance(string(segment.Segment), searchLower)
				if d < di {
					di = d
				}
			}
		}

		for _, segment := range uniseg.Segments([]byte(sjNameLower)) {
			if segment.Category.IsWord() {
				d := editDistance(string(segment.Segment), searchLower)
				if d < dj {
					dj = d
				}
			}
		}

		if di == dj {
			return len(siNameLower) < len(sjNameLower)
		}

		return di < dj
	})
}

func sortStationsByHeightAscending(s stationlist) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].height < s[j].height
	})
}

func sortStationsByYearsAscending(s stationlist) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].dateLast.Sub(s[i].dateFirst) < s[j].dateLast.Sub(s[j].dateFirst)
	})
}

func mergeStations(s1, s2 stationlist) stationlist {
	var result = append(stationlist(nil), s1...)
	for _, station2 := range s2 {
		alreadySeen := false
		for idx, station1 := range s1 {
			if station2.id == station1.id {
				if result[idx].dateFirst.After(station2.dateFirst) {
					result[idx].dateFirst = station2.dateFirst
				}
				if result[idx].dateLast.Before(station2.dateLast) {
					result[idx].dateLast = station2.dateLast
				}
				alreadySeen = true
				break
			}
		}
		if !alreadySeen {
			result = append(result, station2)
		}
	}
	sortStationsByID(result)
	return result
}

func downloadStationsMetaData(r resolution, url string) (stationlist, error) {
	const errDwdParsingStationsMetaData = "dwd: parsing stations list failed"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var bufReader = bufio.NewReaderSize(resp.Body, 256*1024)
	var result []station
	dec := charmap.Windows1250.NewDecoder()

	for {
		line, lineErr := bufReader.ReadBytes('\n')
		if lineErr != nil {
			if lineErr == io.EOF {
				break
			}
			return nil, lineErr
		}
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if firstRune, runeLen := utf8.DecodeRune(line); runeLen < 1 {
			return nil, errors.New(errDwdParsingStationsMetaData)
		} else if !unicode.IsDigit(firstRune) {
			continue
		}

		idx := bytes.LastIndex(line, []byte(" "))
		if idx < 5 {
			return nil, errors.New(errDwdParsingStationsMetaData)
		}

		var s station
		s.resolution = r

		dec.Reset()
		if decoded, decErr := dec.Bytes(line[idx+1:]); decErr != nil {
			s.state = string(line[idx+1:])
		} else {
			s.state = string(decoded)
		}
		line = line[:idx]
		split := bytes.Split(line, []byte("   "))
		var parts [][]byte
		for _, s := range split {
			s = bytes.TrimSpace(s)
			if len(s) > 0 {
				parts = append(parts, s)
			}
		}
		if len(parts) != 4 {
			return nil, errors.New(errDwdParsingStationsMetaData)
		}
		height, heightErr := strconv.Atoi(string(parts[1]))
		if heightErr != nil {
			return nil, heightErr
		}
		s.height = height

		lat, latErr := strconv.ParseFloat(string(parts[2]), 32)
		if latErr != nil {
			return nil, latErr
		}
		s.latitude = lat

		parts0 := bytes.Split(parts[0], []byte{' '})
		if len(parts0) != 3 {
			return nil, errors.New(errDwdParsingStationsMetaData)
		}

		sid, sidErr := strconv.Atoi(string(parts0[0]))
		if sidErr != nil {
			return nil, sidErr
		}
		s.id = sid

		timeStart, timeStartErr := time.Parse("20060102", string(parts0[1]))
		if timeStartErr != nil {
			fmt.Println(string(line))
			return nil, timeStartErr
		}
		s.dateFirst = timeStart

		timeEnd, timeEndErr := time.Parse("20060102", string(parts0[2]))
		if timeEndErr != nil {
			return nil, timeEndErr
		}
		s.dateLast = timeEnd

		parts3 := bytes.SplitN(parts[3], []byte{' '}, 2)
		if len(parts3) != 2 {
			return nil, errors.New(errDwdParsingStationsMetaData)
		}

		long, longErr := strconv.ParseFloat(string(parts3[0]), 32)
		if longErr != nil {
			return nil, longErr
		}
		s.longitude = long

		dec.Reset()
		if decoded, decErr := dec.Bytes(parts3[1]); decErr != nil {
			s.name = string(parts3[1])
		} else {
			s.name = string(decoded)
		}

		s.name = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s.name, "  ", " "), " )", ")"), " -", "-"), "( ", "(")

		result = append(result, s)
	}

	resp.Body.Close()
	return result, nil
}

func stations(r resolution) (stationlist, error) {
	var urls []string
	switch r {
	case resolutionDaily:
		urls = []string{
			"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/daily/more_precip/recent/RR_Tageswerte_Beschreibung_Stationen.txt",
			"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/daily/kl/historical/KL_Tageswerte_Beschreibung_Stationen.txt",
			"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/daily/solar/ST_Tageswerte_Beschreibung_Stationen.txt",
			"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/daily/soil_temperature/historical/EB_Tageswerte_Beschreibung_Stationen.txt",
		}
	case resolutionHourly:
		urls = []string{
			"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/precipitation/historical/RR_Stundenwerte_Beschreibung_Stationen.txt",
			"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/wind/historical/FF_Stundenwerte_Beschreibung_Stationen.txt",
			"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/visibility/historical/VV_Stundenwerte_Beschreibung_Stationen.txt",
			"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/sun/historical/SD_Stundenwerte_Beschreibung_Stationen.txt",
			"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/soil_temperature/historical/EB_Stundenwerte_Beschreibung_Stationen.txt",
			"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/pressure/historical/P0_Stundenwerte_Beschreibung_Stationen.txt",
			"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/hourly/dew_point/historical/TD_Stundenwerte_Beschreibung_Stationen.txt",
		}
	case resolution10Minutes:
		urls = []string{
			"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/air_temperature/recent/zehn_min_tu_Beschreibung_Stationen.txt",
			"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/solar/historical/zehn_min_sd_Beschreibung_Stationen.txt",
			"https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/10_minutes/wind/historical/zehn_min_ff_Beschreibung_Stationen.txt",
		}
	default:
		return nil, errors.New("dwd: unsupported resolution")
	}

	var result []station
	for _, url := range urls {
		currentStations, err := downloadStationsMetaData(r, url)
		if err != nil {
			return nil, err
		}
		result = mergeStations(result, currentStations)
	}
	return result, nil
}
