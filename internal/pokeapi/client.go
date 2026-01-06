package pokeapi

import (
	"net/http"
	"time"

	"github.com/thomas-reed/go-pokedex/internal/pokecache"
)

const (
	baseURL = "https://pokeapi.co/api/v2"
)

type Client struct {
	httpClient http.Client
	cache *pokecache.Cache
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokecache.NewCache(cacheInterval),
	}
}

type Config struct {
	PokeapiClient Client
	NextLocationURL string
	PreviousLocationURL string
	Pokedex map[string]Pokemon
}
