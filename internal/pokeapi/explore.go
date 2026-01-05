package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type exploreResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func (c *Client) ExploreLocation(area string) (exploreResponse, error) {
	url := baseURL + "/location-area/" + area
	
	if result, hit := c.cache.Get(url); hit {
		cachedData := exploreResponse{}
		if err := json.Unmarshal(result, &cachedData); err != nil {
			return exploreResponse{}, fmt.Errorf("Error unmarshalling json body: %s", err)
		}
		c.cache.Extend(url)
		return cachedData, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return exploreResponse{}, fmt.Errorf("Error creating request %s", err)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return exploreResponse{}, fmt.Errorf("Error executing request %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		return exploreResponse{}, fmt.Errorf("Area '%s' does not exist", area)
	}
	if res.StatusCode > 299 {
		return exploreResponse{}, fmt.Errorf("Response Status Code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return exploreResponse{}, fmt.Errorf("Error reading response body: %s", err)
	}

	c.cache.Add(url, body)

	exploreRes := exploreResponse{}
	if err := json.Unmarshal(body, &exploreRes); err != nil {
		return exploreResponse{}, fmt.Errorf("Error unmarshalling json body: %s", err)
	}
	return exploreRes, nil
}
