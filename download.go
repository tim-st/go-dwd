package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func readZipURLs(url string, folderName string, stations []station) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bufReader := bufio.NewReader(resp.Body)

	formattedStationIDs := make([]string, len(stations))
	for idx, station := range stations {
		formattedStationIDs[idx] = fmt.Sprintf("_%05d_", station.id)
		if stations[idx].zipURLs == nil {
			stations[idx].zipURLs = make(map[string][]string, 16)
		}
	}

	for {
		line, lineErr := bufReader.ReadString('\n')
		if lineErr != nil {
			if lineErr == io.EOF {
				break
			}
			return lineErr
		}
		parts := strings.Split(line, `"`)
		for _, part := range parts {
			if strings.HasSuffix(part, ".zip") {

				for idx, formattedStationID := range formattedStationIDs {
					if strings.Contains(part, formattedStationID) {
						stations[idx].zipURLs[folderName] = append(stations[idx].zipURLs[folderName], url+part)
					}
				}

				break
			}
		}
	}

	return nil
}

func produktTxtReadCloser(url string) (io.ReadCloser, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 32*1024*1024))
	if err != nil {
		return nil, err
	}

	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, err
	}

	for _, zipFile := range zipReader.File {
		if strings.HasPrefix(zipFile.Name, "produkt_") && strings.HasSuffix(zipFile.Name, ".txt") {
			zf, err := zipFile.Open()
			if err != nil {
				return nil, err
			}
			return zf, nil
		}
	}

	return nil, errors.New("dwd: produkt_*.txt not found in ZIP file")
}

func parseLine(line []byte, s station) (stationID int, t time.Time, features [][]byte, err error) {
	parts := bytes.Split(line, []byte{';'})
	for idx, part := range parts {
		parts[idx] = bytes.TrimSpace(part)
		if bytes.HasSuffix(parts[idx], []byte("-999")) {
			parts[idx] = nil
		}
	}
	if bytes.Equal(parts[len(parts)-1], []byte(colNameEndOfRow)) {
		parts = parts[:len(parts)-1]
	}
	switch s.resolution {
	case resolutionDaily:
		t, err = time.Parse("20060102", string(parts[1][:8]))
	case resolutionHourly:
		t, err = time.Parse("2006010215", string(parts[1][:10]))
	case resolution10Minutes:
		t, err = time.Parse("200601021504", string(parts[1][:12]))
	default:
		err = errors.New("dwd: unsupported resolution")
	}
	if err != nil {
		return
	}
	stationID, err = strconv.Atoi(string(parts[0]))
	features = parts[2:]
	return
}

func createFile(s station, bufReaders map[string]*bufio.Reader, targetFolder string) error {
	f, fErr := os.Create(filepath.Join(targetFolder, fmt.Sprintf("dwd_%s_%05d.csv.gz", s.resolution.String(), s.id)))
	if fErr != nil {
		return fErr
	}

	bufWriter := bufio.NewWriter(f)
	gzWriter, _ := gzip.NewWriterLevel(bufWriter, gzip.BestCompression)

	var folderNamesResolution = folderNames[s.resolution]

	var stationIDBytes = []byte(fmt.Sprintf("%d", s.id))

	var weatherData = make([]struct {
		readerFinished bool
		hasReadLine    bool
		t              time.Time
		features       [][]byte
	}, len(folderNamesResolution))

	gzWriter.Write([]byte("station_id,messdatum,"))

	for idx1, folder := range folderNamesResolution {
		columns := columns[s.resolution][folder]
		for idx2, column := range columns {
			if idx2 > 1 && idx2 < len(columns)-1 {
				gzWriter.Write([]byte(folder + "__" + column))
				if idx1 == len(folderNamesResolution)-1 && idx2 == len(columns)-2 {
					continue
				}
				gzWriter.Write([]byte{','})
			}
		}
	}

	gzWriter.Write([]byte{'\n'})

	const defaultYear = 2999

	var lastDate = time.Date(defaultYear, 1, 1, 0, 0, 0, 0, time.UTC)

	for {

		var dateFurthestPast = time.Date(defaultYear, 1, 1, 0, 0, 0, 0, time.UTC)
		var numberReady = 0

		for idx, folderName := range folderNamesResolution {

			if !weatherData[idx].readerFinished && !weatherData[idx].hasReadLine {
				bufReader, folderAvailable := bufReaders[folderName]
				if !folderAvailable {
					weatherData[idx].readerFinished = true
					continue
				}
				line, lineErr := bufReader.ReadSlice('\n')
				if lineErr != nil {
					if lineErr == io.EOF {
						delete(bufReaders, folderName)
						weatherData[idx].readerFinished = true
						continue
					}
					return lineErr
				}
				line = bytes.TrimSpace(line)
				if len(line) == 0 {
					continue
				}
				if bytes.HasPrefix(line, []byte(colNameStationID)) {
					continue
				}
				stationID, t, features, err := parseLine(line, s)
				if err != nil {
					return err
				}
				if stationID != s.id {
					return errors.New("dwd: stationID mismatch")
				}

				if lastDate.Year() != defaultYear && !t.After(lastDate) {
					continue
				}

				weatherData[idx].hasReadLine = true
				weatherData[idx].t = t
				weatherData[idx].features = features
			} else {
				numberReady++
			}

			if weatherData[idx].hasReadLine {
				if weatherData[idx].t.Before(dateFurthestPast) {
					dateFurthestPast = weatherData[idx].t
				}
			}

		}

		if numberReady != len(folderNamesResolution) {
			continue
		}

		if len(bufReaders) == 0 {
			break
		}

		lastDate = dateFurthestPast

		gzWriter.Write(stationIDBytes)
		gzWriter.Write([]byte{','})
		gzWriter.Write([]byte(dateFurthestPast.Format("2006-01-02 15:04:05")))
		gzWriter.Write([]byte{','})

		for idx1, t := range weatherData {
			if t.t.Equal(dateFurthestPast) {
				for idx2, feature := range t.features {
					gzWriter.Write(feature)
					if !(idx1 == len(weatherData)-1 && idx2 == len(t.features)-1) {
						gzWriter.Write([]byte{','})
					}
				}
				weatherData[idx1].hasReadLine = false
			} else {
				numberCols := len(columns[s.resolution][folderNamesResolution[idx1]]) - 3
				for idx2 := 0; idx2 < numberCols; idx2++ {
					if !(idx1 == len(weatherData)-1 && idx2 == numberCols-1) {
						gzWriter.Write([]byte{','})
					}
				}
			}
		}

		gzWriter.Write([]byte{'\n'})
	}

	if err := gzWriter.Close(); err != nil {
		return err
	}

	if err := bufWriter.Flush(); err != nil {
		return err
	}

	return f.Close()
}

func downloadStation(s station, targetFolder string) error {
	var readClosers = map[string][]io.ReadCloser{}
	var bufReaders = map[string]*bufio.Reader{}

	for folderName, urls := range s.zipURLs {
		var readers = []io.Reader{}
		for _, url := range urls {
			readCloser, err := produktTxtReadCloser(url)
			if err != nil {
				return err
			}
			readClosers[folderName] = append(readClosers[folderName], readCloser)
			readers = append(readers, readCloser)
		}
		bufReaders[folderName] = bufio.NewReader(io.MultiReader(readers...))
	}

	if err := createFile(s, bufReaders, targetFolder); err != nil {
		return err
	}

	for _, rcs := range readClosers {
		for _, rc := range rcs {
			if err := rc.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func downloadStations(stations []station, r resolution, targetFolder string) error {
	urlMap, containsResolution := urls[r]

	if !containsResolution {
		return errors.New("dwd: unknown resolution")
	}

	for folder, urls := range urlMap {
		for _, url := range urls {
			if err := readZipURLs(url, folder, stations); err != nil {
				return err
			}
		}
	}

	for _, station := range stations {
		if err := downloadStation(station, targetFolder); err != nil {
			return err
		}
	}

	return nil
}
