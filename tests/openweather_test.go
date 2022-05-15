package tests

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testbackend/internal/weather"
	"testing"
)

func TestOpenWeather(t *testing.T) {
	req, err := weather.NewOpenWeatherMapRequest(context.Background(), "Moscow", "ru")
	assert.NoError(t, err, "request creating failed")

	fmt.Println(req.URL)

	client := http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err, "failed response")

	defer resp.Body.Close()
	parsed, err := weather.ParseResponseData(resp)
	assert.NoError(t, err)

	assert.Equal(t, parsed.Cod, 200)

	fmt.Println(parsed)

	w, err := weather.NewWeather(parsed)
	assert.NoError(t, err)

	fmt.Println(w)
}
