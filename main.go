package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

const (
	urlViaCep     = "https://viacep.com.br/"
	urlWeatherApi = "https://api.weatherapi.com/v1/"
	weatherApiKey = "66f074cb079d40c1bee04932241302"
)

type ViaCepResponse struct {
	Erro        bool   `json:"erro"`
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type WeatherApiResponse struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}

func main() {
	app := fiber.New()
	app.Get("/:cep", handleRequest)
	log.Fatal(app.Listen(":8080"))
}

func handleRequest(c *fiber.Ctx) error {
	cep := c.Params("cep")
	cep = strings.Replace(cep, "-", "", -1)
	if len(cep) != 8 {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": "invalid zipcode"})
	}

	url := urlViaCep + "ws/" + cep + "/json"
	response, err := fetchData(c, url)
	if err != nil || response == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "can not find zipcode"})
	}
	cepResponse := ViaCepResponse{}
	err = json.Unmarshal(response, &cepResponse)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing zipcode data"})
	}

	if cepResponse.Erro {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "can not found zipcode"})
	}

	city := string(cepResponse.Localidade)
	city = removeAccents(city)
	state := string(cepResponse.Uf)

	url = urlWeatherApi + "current.json?key=" + weatherApiKey + "&q=" + city + " - " + state + " - Brazil&aqi=no&tides=no"
	url = strings.Replace(url, " ", "%20", -1)

	response, err = fetchData(c, url)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching weather data"})
	}

	weatherResponse := WeatherApiResponse{}
	err = json.Unmarshal(response, &weatherResponse)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing weather data"})
	}

	tempC := strconv.FormatFloat(weatherResponse.Current.TempC, 'f', -1, 64)
	tempF := strconv.FormatFloat(weatherResponse.Current.TempF, 'f', -1, 64)
	tempK := strconv.FormatFloat(weatherResponse.Current.TempC+273.15, 'f', -1, 64)

	response = []byte(`{ "temp_C": ` + tempC + `, "temp_F": ` + tempF + `, "temp_K": ` + tempK + ` }`)

	return c.Send(response)
}

func fetchData(c *fiber.Ctx, url string) (response []byte, err error) {
	client := fasthttp.Client{}

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(url)

	if err := client.DoTimeout(req, res, 30*time.Second); err != nil {
		return nil, err
	}

	if res.StatusCode() != fiber.StatusOK {
		return nil, errors.New("invalid statuscode")
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	return res.Body(), nil
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}
