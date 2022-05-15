package weather

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"testbackend/internal/cache"
	"time"
)

var apiKey = os.Getenv("OPENWEATHER_API_KEY")

const serviceURL = "https://api.openweathermap.org/data/2.5/weather"

type OpenWeatherMapService struct {
	requestTime time.Time
	client      *http.Client
	cache       cache.Cache
}

func (o *OpenWeatherMapService) Init(context.Context) error {
	o.client = http.DefaultClient
	o.cache = cache.NewMapCache()
	_, _ = o.GetWeather()
	return nil
}

func (o *OpenWeatherMapService) Ping(context.Context) error {
	return nil
}

func (o *OpenWeatherMapService) Close() error {
	return nil
}

func (o *OpenWeatherMapService) GetWeather() (Weather, error) {
	strTime := time.Now().Format("01-01-2002")
	value, err := o.cache.Get(strTime)
	if err == nil {
		return value.(Weather), nil
	}
	req, err := NewOpenWeatherMapRequest(context.Background(), "Moscow", "ru")
	if err != nil {
		return Weather{}, err
	}
	resp, err := o.client.Do(req)
	if err != nil {
		return Weather{}, err
	}
	defer resp.Body.Close()
	parsed, err := ParseResponseData(resp)
	if err != nil {
		return Weather{}, err
	}
	w, err := NewWeather(parsed)
	if err = o.cache.Set(strTime, w); err != nil {
		return Weather{}, err
	}
	return w, nil
}

func NewOpenWeatherMapRequest(ctx context.Context, city, country string) (*http.Request, error) {
	URL, err := url.ParseRequestURI(serviceURL)
	if err != nil {
		return nil, err
	}
	query := url.Values{
		"q":     []string{city + "," + country},
		"appid": []string{apiKey},
	}
	URL.RawQuery = query.Encode()
	return http.NewRequestWithContext(ctx, http.MethodGet, URL.String(), http.NoBody)
}

func ParseResponseData(resp *http.Response) (result OpenWeatherResponse, err error) {
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}
	if result.Cod != 200 {
		err = errors.New("wrong response status")
	}
	return
}

func NewWeather(data OpenWeatherResponse) (Weather, error) {
	const AbsoluteZero = 273.15
	if len(data.Weather) == 0 {
		return Weather{}, errors.New("invalid OpenWeather object")
	}
	weather := Weather{
		Description: DescriptionModel{
			Main:        data.Weather[0].Main,
			Description: data.Weather[0].Description,
		},
		Metrics: MetricsModel{
			Temp:      int(data.Main.Temp - AbsoluteZero),
			TempMin:   int(data.Main.TempMin - AbsoluteZero),
			TempMax:   int(data.Main.TempMax - AbsoluteZero),
			Pressure:  data.Main.Pressure,
			Humidity:  data.Main.Humidity,
			MoonPhase: MoonPhase(time.Unix(int64(data.Dt), 0)),
		},
		Name: data.Name,
		Icon: data.Weather[0].Icon,
	}
	return weather, nil
}

func (w *Weather) WriteTo(writer io.Writer) (int64, error) {
	rawData, err := json.Marshal(w)
	if err != nil {
		return 0, err
	}
	n, err := writer.Write(rawData)
	return int64(n), err
}
