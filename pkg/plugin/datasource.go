package plugin

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "io"
    "net/http"

    "github.com/grafana/grafana-plugin-sdk-go/backend"
    "github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
    "github.com/grafana/grafana-plugin-sdk-go/data"
    "github.com/grafana/weather/pkg/models"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces - only those which are required for a particular task.
var (
    _ backend.QueryDataHandler      = (*Datasource)(nil)
    _ backend.CheckHealthHandler    = (*Datasource)(nil)
    _ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

const (
    baseURL = "http://api.openweathermap.org/data/2.5/forecast"
)

// NewDatasource creates a new datasource instance.
func NewDatasource(_ context.Context, _ backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
    return &Datasource{}, nil
}

// Datasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type Datasource struct{}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
    // Clean up datasource instance resources.
}

type WeatherResponse struct {
    List []struct {
       Dt   int64 `json:"dt"`
       Main struct {
          Temp float64 `json:"temp"`
       } `json:"main"`
       Weather []struct {
          Description string `json:"description"`
       } `json:"weather"`
    } `json:"list"`
}

type WeatherData struct {
    Time        time.Time
    Temperature float64
    Description string
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
    // create response struct
    response := backend.NewQueryDataResponse()

    // loop over queries and execute them individually.
    for _, q := range req.Queries {
       res := d.query(ctx, req.PluginContext, q)

       // save the response in a hashmap
       // based on with RefID as identifier
       response.Responses[q.RefID] = res
    }

    return response, nil
}

type queryModel struct {
    City string `json:"city"`
}

func (d *Datasource) GetHistoricalWeather(city string, apiKey string) ([]WeatherData, error) {
    endDate := time.Now()
    startDate := endDate.AddDate(0, 0, -5)

    url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric&start=%d&end=%d", baseURL, city, apiKey, startDate.Unix(), endDate.Unix())

    resp, err := http.Get(url)
    if err != nil {
       return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
       return nil, err
    }

    if resp.StatusCode != http.StatusOK {
       return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
    }

    var weatherResponse WeatherResponse
    err = json.Unmarshal(body, &weatherResponse)
    if err != nil {
       return nil, err
    }

    weatherData := make([]WeatherData, len(weatherResponse.List))
    for i, forecast := range weatherResponse.List {
       weatherData[i] = WeatherData{
          Time:        time.Unix(forecast.Dt, 0),
          Temperature: forecast.Main.Temp,
          Description: forecast.Weather[0].Description,
       }
    }

    return weatherData, nil
}

func (d *Datasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
    var response backend.DataResponse
    // get our secrets which we saved in config editor via plugin overview.
    config, err := models.LoadPluginSettings(*pCtx.DataSourceInstanceSettings)
    if err != nil {
       response.Error = err
       return response
    }

    // Unmarshal the JSON into our queryModel.
    var qm queryModel

    err = json.Unmarshal(query.JSON, &qm)
    if err != nil {
       return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
    }

    // create data frame response.
    // For an overview on data frames and how grafana handles them:
    // https://grafana.com/developers/plugin-tools/introduction/data-frames
    frame := data.NewFrame("response")

    weatherData, err := d.GetHistoricalWeather(qm.City, config.Secrets.ApiKey)
    if err != nil {
       response.Error = err
       return response
    }

    times := make([]time.Time, len(weatherData[:10]))
    temperatures := make([]float64, len(weatherData[:10]))
    descriptions := make([]string, len(weatherData[:10]))

    for i, d := range weatherData[:10] {
       times[i] = d.Time
       temperatures[i] = d.Temperature
       descriptions[i] = d.Description
    }

    frame.Fields = append(frame.Fields,
       data.NewField("time", nil, times),
       data.NewField("description", nil, descriptions),
       data.NewField("temperature", nil, temperatures),
    )
    response.Frames = append(response.Frames, frame)

    return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
    res := &backend.CheckHealthResult{}
    config, err := models.LoadPluginSettings(*req.PluginContext.DataSourceInstanceSettings)

    if err != nil {
       res.Status = backend.HealthStatusError
       res.Message = "Unable to load settings"
       return res, nil
    }

    if config.Secrets.ApiKey == "" {
       res.Status = backend.HealthStatusError
       res.Message = "API key is missing"
       return res, nil
    }

    return &backend.CheckHealthResult{
       Status:  backend.HealthStatusOk,
       Message: "Data source is working",
    }, nil
}