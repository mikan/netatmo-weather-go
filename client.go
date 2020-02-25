package netatmo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/oauth2"
)

// TargetMeasurements defines list of target measurement attributes.
var TargetMeasurements = []string{"Temperature", "CO2", "Humidity", "Pressure", "Noise", "WindStrength", "WindAngle",
	"GustStrength", "GustAngle"}

// Client implements Netatmo API client.
type Client struct {
	oauth  *oauth2.Config
	client *http.Client
}

// Measure defines each measurable series.
type Measure struct {
	DeviceID     string
	ModuleID     string
	Timestamp    int64
	Temperature  *float64 // Nullable
	CO2          *int     // Nullable
	Humidity     *int     // Nullable
	Pressure     *float64 // Nullable
	Noise        *int     // Nullable
	WindStrength *int     // Nullable
	WindAngle    *int     // Nullable
	GustStrength *int     // Nullable
	GustAngle    *int     // Nullable
}

// Place defines place attributes.
type Place struct {
	Altitude int       `json:"altitude"`
	City     string    `json:"city"`     // Name of city (ex. 千代田区)
	Country  string    `json:"country"`  // Country code (ex. JP)
	Timezone string    `json:"timezone"` // TZ Database name (ex. Asia/Tokyo)
	Location []float64 `json:"location"` // Lat, Lon (ex. 139.752778, 35.682500)
}

// Latitude returns latitude value from location data.
func (n *Place) Latitude() float64 {
	if len(n.Location) != 2 {
		return 0
	}
	return n.Location[0]
}

// Longitude returns longitude value location data.
func (n *Place) Longitude() float64 {
	if len(n.Location) != 2 {
		return 0
	}
	return n.Location[1]
}

// DashboardData defines newest measured data gathered by device or module.
type DashboardData struct {
	UTCTime             int64    `json:"time_utc"`
	Temperature         *float64 // Nullable
	MinTemperature      *float64 `json:"min_temp"`      // Nullable
	MaxTemperature      *float64 `json:"max_temp"`      // Nullable
	MinTemperatureTime  *int64   `json:"date_min_temp"` // Nullable
	MaxTemperatureTime  *int64   `json:"date_max_temp"` // Nullable
	TemperatureTrend    *string  `json:"temp_trend"`    // Nullable
	CO2                 *int     // Nullable
	Humidity            *int     // Nullable
	Noise               *int     // Nullable
	Pressure            *float64 // Nullable
	AbsolutePressure    *float64 // Nullable
	PressureTrend       *string  `json:"pressure_trend"` // Nullable
	Rain                *float64 // Nullable
	RainPerHour         *float64 `json:"sum_rain_1"`  // Nullable
	RainPerDay          *float64 `json:"sum_rain_24"` // Nullable
	GustAngle           *int     // Nullable
	GustStrength        *int     // Nullable
	WindAngle           *int     // Nullable
	WindStrength        *int     // Nullable
	MaxWindStrength     *int     `json:"max_wind_str"`      // Nullable
	MaxWindStrengthTime *int64   `json:"date_max_wind_str"` // Nullable
	HealthIndex         *int     `json:"health_idx"`        // Nullable
}

// Module defines netatmo module attributes.
type Module struct {
	ID              string         `json:"_id"`
	Type            string         `json:"type"`
	ModuleName      string         `json:"module_name"`
	DataTypes       []string       `json:"data_type"`
	LastSetupTime   int64          `json:"last_setup"`
	Reachable       bool           `json:"reachable"`
	Firmware        int            `json:"firmware"`
	LastMessageTime int64          `json:"last_message"`
	LastSeenTime    int64          `json:"last_seen"`
	RFStatus        int            `json:"rf_status"`
	BatteryVP       int            `json:"battery_vp"`
	BatteryPercent  int            `json:"battery_percent"`
	DashboardData   *DashboardData `json:"dashboard_data"` // Nullable
}

// Device defines netatmo device attributes.
type Device struct {
	ID                  string         `json:"_id"`
	CipherID            string         `json:"cipher_id"`
	SetupTime           int64          `json:"date_setup"`
	LastSetupTime       int64          `json:"last_setup"`
	Type                string         `json:"type"`
	LastStatusStoreTime int64          `json:"last_status_store"`
	ModuleName          string         `json:"module_name"`
	Firmware            int            `json:"firmware"`
	LastUpgradeTime     int64          `json:"last_upgrade"`
	WiFiStatus          int            `json:"wifi_status"`
	Reachable           bool           `json:"reachable"`
	CO2Calibrating      bool           `json:"co2_calibrating"`
	StationName         string         `json:"station_name"`
	DataTypes           []string       `json:"data_type"`
	Place               Place          `json:"place"`
	DashboardData       *DashboardData `json:"dashboard_data"` // Nullable
	Modules             []Module       `json:"modules"`
}

// Administrative defines user administrative attributes.
type Administrative struct {
	Language          string `json:"lang"`           // user locale
	DisplayLocale     string `json:"reg_locale"`     // user regional preferences (used for displaying date)
	Country           string `json:"country"`        // user country
	Unit              int    `json:"unit"`           // 0 -> metric system, 1 -> imperial system
	WindUnit          int    `json:"windunit"`       // 0 -> kph, 1 -> mph, 2 -> ms, 3 -> beaufort, 4 -> knot
	PressureUnit      int    `json:"pressureunit"`   // 0 -> mbar, 1 -> inHg, 2 -> mmHg
	FeelLikeAlgorithm int    `json:"feel_like_algo"` // algorithm used to compute feel like temperature, 0 -> humidex, 1 -> heat-index
}

// DescribeUnit describes unit identifier.
func (a *Administrative) DescribeUnit() string {
	switch a.Unit {
	case 0:
		return "metric system"
	case 1:
		return "imperial system"
	default:
		return fmt.Sprintf("unknown unit: %d", a.Unit)
	}
}

// DescribeWindUnit describes wind unit identifier.
func (a *Administrative) DescribeWindUnit() string {
	switch a.WindUnit {
	case 0:
		return "kph"
	case 1:
		return "mph"
	case 2:
		return "ms"
	case 3:
		return "beaufort"
	case 4:
		return "knot"
	default:
		return fmt.Sprintf("unknown wind unit: %d", a.WindUnit)
	}
}

// DescribePressureUnit describes pressure unit identifier.
func (a *Administrative) DescribePressureUnit() string {
	switch a.PressureUnit {
	case 0:
		return "mbar"
	case 1:
		return "inHg"
	case 2:
		return "mmHg"
	default:
		return fmt.Sprintf("unknown pressure unit: %d", a.PressureUnit)
	}
}

// DescribeFeelLikeAlgorithm describes identifier of algorithm used to compute feel like temperature.
func (a *Administrative) DescribeFeelLikeAlgorithm() string {
	switch a.FeelLikeAlgorithm {
	case 0:
		return "humidex"
	case 1:
		return "heat-index"
	default:
		return fmt.Sprintf("unknown feel like algorithm: %d", a.FeelLikeAlgorithm)
	}
}

// User defines user attributes.
type User struct {
	Mail           string         `json:"mail"`
	Administrative Administrative `json:"administrative"`
}

type stationsDataBody struct {
	Devices []Device `json:"devices"`
	User    User     `json:"user"`
}

type getStationsDataResponse struct {
	Body       stationsDataBody `json:"body"`
	Status     string           `json:"status"`
	ExecTime   float64          `json:"time_exec"`
	ServerTime int64            `json:"time_server"`
}

type measureBody struct {
	BeginTime int64        `json:"beg_time"`
	StepTime  int64        `json:"step_time"`
	Value     [][]*float64 `json:"value"`
}

type getMeasureResponse struct {
	Body       []measureBody `json:"body"`
	Status     string        `json:"status"`
	ExecTime   float64       `json:"time_exec"`
	ServerTime int64         `json:"time_server"`
}

// NewClient will creates Netatmo client object.
func NewClient(ctx context.Context, clientID, clientSecret, username, password string) (*Client, error) {
	oauth := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"read_station"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.netatmo.net/",
			TokenURL: "https://api.netatmo.net/oauth2/token",
		},
	}
	token, err := oauth.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		return nil, err
	}
	return &Client{
		oauth:  oauth,
		client: oauth.Client(ctx, token),
	}, err
}

// GetStationsData gathers station data from Netatmo API.
// Reference: https://dev.netatmo.com/apidocumentation/weather#getstationsdata
func (c *Client) GetStationsData() ([]Device, *User, error) {
	resp, err := c.client.Get("https://api.netatmo.com/api/getstationsdata")
	if err != nil {
		return nil, nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	var respData getStationsDataResponse
	if err := json.Unmarshal(data, &respData); err != nil {
		return nil, nil, err
	}
	return respData.Body.Devices, &respData.Body.User, nil
}

// GetMeasureByTimeRange gathers measure data by specified time window.
// Reference: https://dev.netatmo.com/apidocumentation/weather#getmeasure
func (c *Client) GetMeasureByTimeRange(deviceID, moduleID string, begin, end int64) ([]Measure, error) {
	resp, err := c.client.Get("https://api.netatmo.com/api/getmeasure" +
		"?device_id=" + deviceID +
		"&module_id=" + moduleID +
		"&scale=max" + // {max, 30min, 1hour, 3hours, 1day, 1week, 1month}
		"&type=" + strings.Join(TargetMeasurements, ",") +
		"&real_time=true" + // default: false
		"&date_begin=" + strconv.FormatInt(begin, 10) +
		"&date_end=" + strconv.FormatInt(end, 10))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buildGetMeasureResponse(deviceID, moduleID, data)
}

// GetMeasureByNewest gathers newest measure data.
// Reference: https://dev.netatmo.com/apidocumentation/weather#getmeasure
func (c *Client) GetMeasureByNewest(deviceID, moduleID string) (*Measure, error) {
	resp, err := c.client.Get("https://api.netatmo.com/api/getmeasure" +
		"?device_id=" + deviceID +
		"&module_id=" + moduleID +
		"&scale=max" + // {max, 30min, 1hour, 3hours, 1day, 1week, 1month}
		"&type=" + strings.Join(TargetMeasurements, ",") +
		"&date_end=last")
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	measures, err := buildGetMeasureResponse(deviceID, moduleID, data)
	if err != nil {
		return nil, err
	}
	if measures == nil {
		return nil, nil // No Data
	}
	return &measures[len(measures)-1], nil
}

func buildGetMeasureResponse(deviceID, moduleID string, data []byte) ([]Measure, error) {
	var response getMeasureResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	var measures []Measure
	for _, v := range response.Body {
		for i, m := range v.Value {
			measure := Measure{
				DeviceID:     deviceID,
				ModuleID:     moduleID,
				Timestamp:    v.BeginTime + (v.StepTime * int64(i)),
				Temperature:  handleFloat(m[0]),
				CO2:          handleInt(m[1]),
				Humidity:     handleInt(m[2]),
				Pressure:     handleFloat(m[3]),
				Noise:        handleInt(m[4]),
				WindStrength: handleInt(m[5]),
				WindAngle:    handleInt(m[6]),
				GustStrength: handleInt(m[7]),
				GustAngle:    handleInt(m[8]),
			}
			measures = append(measures, measure)
		}
	}
	if len(measures) == 0 {
		return nil, nil
	}
	return measures, nil
}

func handleFloat(v *float64) *float64 {
	if v == nil {
		return nil
	}
	if *v == 0.0 { // If the value exactly matches 0.0, treat it as null value
		return nil
	}
	return v
}

func handleInt(v *float64) *int {
	if v == nil {
		return nil
	}
	if *v == 0.0 { // If the value exactly matches 0.0, treat it as null value
		return nil
	}
	iv := int(*v)
	return &iv
}
