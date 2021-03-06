package helmops

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

type Chart struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	ApiVersion  string   `json:"apiVersion"`
	AppVersion  string   `json:"appVersion"`
	Type        string   `json:"type"`
	Urls        []string `json:"urls"`
	Created     string   `json:"created"`
	Digest      string   `json:"digest"`
	Message     string   `json:"error"`
}

type ErrorResponse struct {
	Message string `json:"error"`
}

var ChartMuseumHost = os.Getenv("CHARTMUSEUM_SERVER_IP")

/*
 * Returns all charts with all versions
 */
func GetLaunches() (map[string][]Chart, error) {

	resp, err := http.Get(ChartMuseumHost + "/api/charts")
	if err != nil {
		return nil, err
	}
	var launches map[string][]Chart
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &launches)
	if err != nil {
		return nil, err
	}

	return launches, nil
}

/*
 * Returns a chart with all versions
 */
func GetLaunch(name string) ([]Chart, error) {

	resp, err := http.Get(ChartMuseumHost + "/api/charts/" + name)
	if err != nil {
		return nil, err
	}
	var launch []Chart
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &launch)
	if err != nil {
		var errResponse ErrorResponse
		err = json.Unmarshal(body, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errResponse.Message)
	}

	return launch, nil
}

/*
 * Returns a chart with all versions
 */
func GetLaunchWithVersion(name string, version string) (Chart, error) {

	resp, err := http.Get(ChartMuseumHost + "/api/charts/" + name + "/" + version)
	if err != nil {
		return Chart{}, err
	}
	var launch Chart
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Chart{}, err
	}

	err = json.Unmarshal(body, &launch)
	if err != nil {
		return Chart{}, err
	}

	if launch.Message != "" {
		return Chart{}, errors.New(launch.Message)
	}

	return launch, nil
}
