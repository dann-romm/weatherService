# weatherService

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com) [![forthebadge](http://forthebadge.com/images/badges/built-with-love.svg)](http://forthebadge.com)

Simple backend server, providing access to weather data with helpful application lifecycle control manager
> Currently running in docker-container on my vps server

# Usage

Service provides an entry point `98.142.251.23:8900/current-weather`
### The output is a json data like this:
```json
{
    "weather": {
        "main": "Rain",
        "description": "light rain"
    },
    "main": {
        "temp": 19,
        "temp_min": 17,
        "temp_max": 20,
        "pressure": 1017,
        "humidity": 62,
        "moon_phase": "waxing crescent"
    },
    "name": "Moscow",
    "icon": "10n"
}
```

# How it works

There is an `Application` instance that controlls app lifecycle
It has own `Resources`, that we can initialize, check it's health, gracefull stop `Resources` and release it.
In according with [low coupling](https://en.wikipedia.org/wiki/Coupling_(computer_programming)) pattern, `Resources` is an interface that `ServiceKeeper` implements.
`ServiceKeeper` in turn, has a couple of services, that we can initialize, and in case of error, gracefull stop
___
Application uses [openWeatherMap](https://openweathermap.org/`) API to collect weather data, processes it and cache
