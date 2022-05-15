package tests

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testbackend/internal/weather"
	"testing"
	"time"
)

type T struct {
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	Timezone       string  `json:"timezone"`
	TimezoneOffset int     `json:"timezone_offset"`
	Daily          []struct {
		Dt        int     `json:"dt"`
		Sunrise   int     `json:"sunrise"`
		Sunset    int     `json:"sunset"`
		Moonrise  int     `json:"moonrise"`
		Moonset   int     `json:"moonset"`
		MoonPhase float64 `json:"moon_phase"`
		Temp      struct {
			Day   float64 `json:"day"`
			Min   float64 `json:"min"`
			Max   float64 `json:"max"`
			Night float64 `json:"night"`
			Eve   float64 `json:"eve"`
			Morn  float64 `json:"morn"`
		} `json:"temp"`
		FeelsLike struct {
			Day   float64 `json:"day"`
			Night float64 `json:"night"`
			Eve   float64 `json:"eve"`
			Morn  float64 `json:"morn"`
		} `json:"feels_like"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		DewPoint  float64 `json:"dew_point"`
		WindSpeed float64 `json:"wind_speed"`
		WindDeg   int     `json:"wind_deg"`
		WindGust  float64 `json:"wind_gust"`
		Weather   []struct {
			Id          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds int     `json:"clouds"`
		Pop    float64 `json:"pop"`
		Rain   float64 `json:"rain,omitempty"`
		Uvi    float64 `json:"uvi"`
	} `json:"daily"`
}

func TestMoonPhase(t *testing.T) {
	u, err := url.Parse("https://api.openweathermap.org/data/2.5/onecall?lat=37.6156&lon=55.7522&exclude=current,minutely,hourly,alerts&appid=" + weather.apiKey)
	assert.NoError(t, err)
	r, err := (&http.Client{}).Do(&http.Request{URL: u})
	assert.NoError(t, err)
	var tmp T
	err = json.NewDecoder(r.Body).Decode(&tmp)
	assert.NoError(t, err)

	for _, day := range tmp.Daily {
		timeCalc := time.Unix(int64(day.Dt), 0)
		assert.True(t, weather.CalcMoonPhase(timeCalc)-day.MoonPhase <= 0.3)
		fmt.Println(weather.MoonPhase(timeCalc))
	}
}
