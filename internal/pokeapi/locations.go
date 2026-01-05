package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type mapResponse struct {
	Count    int `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (c *Client) ListLocations(pageURL string) (mapResponse, error) {
	url := baseURL + "/location-area"
	if pageURL != "" {
		url = pageURL
	}

	if result, hit := c.cache.Get(url); hit {
		cachedData := mapResponse{}
		if err := json.Unmarshal(result, &cachedData); err != nil {
			return mapResponse{}, fmt.Errorf("Error unmarshalling json body: %s", err)
		}
		c.cache.Extend(url)
		return cachedData, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return mapResponse{}, fmt.Errorf("Error creating request %s", err)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return mapResponse{}, fmt.Errorf("Error executing request %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return mapResponse{}, fmt.Errorf("Response Status Code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return mapResponse{}, fmt.Errorf("Error reading response body: %s", err)
	}

	c.cache.Add(url, body)

	mapRes := mapResponse{}
	if err := json.Unmarshal(body, &mapRes); err != nil {
		return mapResponse{}, fmt.Errorf("Error unmarshalling json body: %s", err)
	}
	return mapRes, nil
}
