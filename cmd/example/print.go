package main

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mikan/netatmo-weather-go"
)

func printStationsData(devices []netatmo.Device, user netatmo.User, w io.Writer) error {
	tw := new(tabwriter.Writer).Init(w, 0, 8, 1, '\t', 0)
	must(fmt.Fprintln(tw, "User information:"))
	must(fmt.Fprintf(tw, "\tMail:\t%s\n", user.Mail))
	must(fmt.Fprintf(tw, "\tLanguage:\t%s\n", user.Administrative.Language))
	must(fmt.Fprintf(tw, "\tDisplay locale:\t%s\n", user.Administrative.DisplayLocale))
	must(fmt.Fprintf(tw, "\tCountry:\t%s\n", user.Administrative.Country))
	must(fmt.Fprintf(tw, "\tUnit:\t%s\n", user.Administrative.DescribeUnit()))
	must(fmt.Fprintf(tw, "\tWind unit:\t%s\n", user.Administrative.DescribeWindUnit()))
	must(fmt.Fprintf(tw, "\tPressure unit:\t%s\n", user.Administrative.DescribePressureUnit()))
	must(fmt.Fprintf(tw, "\tFeel like algorithm:\t%s\n", user.Administrative.DescribeFeelLikeAlgorithm()))
	for i := 0; i < len(devices); i++ {
		d := devices[i]
		must(fmt.Fprintln(tw))
		must(fmt.Fprintf(tw, "Device %d of %d:\n", i+1, len(devices)))
		must(fmt.Fprintf(tw, "\tDevice ID:\t%s\n", d.ID))
		must(fmt.Fprintf(tw, "\tModule name:\t%s\n", d.ModuleName))
		must(fmt.Fprintf(tw, "\tStation name:\t%s\n", d.StationName))
		must(fmt.Fprintf(tw, "\tType:\t%s\n", d.Type))
		must(fmt.Fprintf(tw, "\tData types:\t%s\n", strings.Join(d.DataTypes, ", ")))
		must(fmt.Fprintf(tw, "\tCipher ID:\t%s\n", d.CipherID))
		must(fmt.Fprintf(tw, "\tFirmware:\t%d\n", d.Firmware))
		must(fmt.Fprintf(tw, "\tWi-Fi status:\t%d\n", d.WiFiStatus))
		must(fmt.Fprintf(tw, "\tReachable:\t%t\n", d.Reachable))
		must(fmt.Fprintf(tw, "\tCO2 calibrating:\t%t\n", d.CO2Calibrating))
		must(fmt.Fprintf(tw, "\tCountry:\t%s\n", d.Place.Country))
		must(fmt.Fprintf(tw, "\tCity:\t%s\n", d.Place.City))
		must(fmt.Fprintf(tw, "\tTime zone:\t%s\n", d.Place.Timezone))
		must(fmt.Fprintf(tw, "\tAltitude:\t%d\n", d.Place.Altitude))
		must(fmt.Fprintf(tw, "\tLocation:\t%f, %f\n", d.Place.Latitude(), d.Place.Longitude()))
		must(fmt.Fprintf(tw, "\tSetup time:\t%s\n", formatTimestamp(d.SetupTime)))
		must(fmt.Fprintf(tw, "\tLast setup time:\t%s\n", formatTimestamp(d.LastSetupTime)))
		must(fmt.Fprintf(tw, "\tLast upgrade time:\t%s\n", formatTimestamp(d.LastUpgradeTime)))
		must(fmt.Fprintf(tw, "\tLast status store time:\t%s\n", formatTimestamp(d.LastStatusStoreTime)))
		printDashboardData("", tw, d.DashboardData, d.DataTypes)
		for j := 0; j < len(d.Modules); j++ {
			m := d.Modules[j]
			must(fmt.Fprintln(tw))
			must(fmt.Fprintf(tw, "\tModule %d of %d:\n", j+1, len(d.Modules)))
			must(fmt.Fprintf(tw, "\t\tModule ID:\t%s\n", m.ID))
			must(fmt.Fprintf(tw, "\t\tModule name:\t%s\n", m.ModuleName))
			must(fmt.Fprintf(tw, "\t\tData types:\t%s\n", strings.Join(m.DataTypes, ", ")))
			must(fmt.Fprintf(tw, "\t\tFirmware:\t%d\n", m.Firmware))
			must(fmt.Fprintf(tw, "\t\tRF status:\t%d\n", m.RFStatus))
			must(fmt.Fprintf(tw, "\t\tBattery:\t%d %% (vp: %d)\n", m.BatteryPercent, m.BatteryVP))
			must(fmt.Fprintf(tw, "\t\tReachable:\t%t\n", m.Reachable))
			must(fmt.Fprintf(tw, "\t\tLast setup time:\t%s\n", formatTimestamp(m.LastSetupTime)))
			must(fmt.Fprintf(tw, "\t\tLast message time:\t%s\n", formatTimestamp(m.LastMessageTime)))
			must(fmt.Fprintf(tw, "\t\tLast seen time:\t%s\n", formatTimestamp(m.LastSeenTime)))
			printDashboardData("\t", tw, m.DashboardData, m.DataTypes)
		}
	}
	return tw.Flush()
}

func printMeasures(values []netatmo.Measure, w io.Writer) error {
	tw := new(tabwriter.Writer).Init(w, 0, 8, 1, '\t', 0)
	must(fmt.Fprintln(tw, "Timestamp\t"+strings.Join(netatmo.TargetMeasurements, "\t")))
	for _, m := range values {
		must(fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			time.Unix(m.Timestamp, 0).Format("2006/01/02 15:04:05"),
			f64OrNull(m.Temperature),
			intOrNull(m.CO2),
			intOrNull(m.Humidity),
			f64OrNull(m.Pressure),
			intOrNull(m.Noise),
			intOrNull(m.WindStrength),
			intOrNull(m.WindAngle),
			intOrNull(m.GustStrength),
			intOrNull(m.GustAngle)))
	}
	return tw.Flush()
}

func printDashboardData(prefix string, w io.Writer, data *netatmo.DashboardData, types []string) {
	if data == nil {
		must(fmt.Fprintln(w, prefix+"\tDashboard data:\t(no data)"))
		return
	}
	must(fmt.Fprintln(w, prefix+"\tDashboard data:"))
	must(fmt.Fprintf(w, prefix+"\t\tTime (UTC):\t%s\n", formatTimestamp(data.UTCTime)))
	if sliceContains(types, "Temperature") {
		must(fmt.Fprintf(w, prefix+"\t\tTemperature:\t%.1f °C (trend: %s)\n", *data.Temperature, *data.TemperatureTrend))
		must(fmt.Fprintf(w, prefix+"\t\tMinimum temperature:\t%.1f °C (at %s)\n", *data.MinTemperature,
			formatTimestamp(*data.MinTemperatureTime)))
		must(fmt.Fprintf(w, prefix+"\t\tMaximum temperature:\t%.1f °C (at %s)\n", *data.MaxTemperature,
			formatTimestamp(*data.MaxTemperatureTime)))
	}
	if sliceContains(types, "CO2") {
		must(fmt.Fprintf(w, prefix+"\t\tCO2:\t%d ppm\n", *data.CO2))
	}
	if sliceContains(types, "Humidity") {
		must(fmt.Fprintf(w, prefix+"\t\tHumidity:\t%d %%\n", *data.Humidity))
	}
	if sliceContains(types, "Noise") {
		must(fmt.Fprintf(w, prefix+"\t\tNoise:\t%d db\n", *data.Noise))
	}
	if sliceContains(types, "Pressure") {
		must(fmt.Fprintf(w, prefix+"\t\tPressure:\t%.1f mb (trend: %s)\n", *data.Pressure, *data.PressureTrend))
		must(fmt.Fprintf(w, prefix+"\t\tAbsolute pressure:\t%.1f mb\n", *data.AbsolutePressure))
	}
	if sliceContains(types, "Rain") {
		must(fmt.Fprintf(w, prefix+"\t\tRain:\t%.1f mm\n", *data.Rain))
		must(fmt.Fprintf(w, prefix+"\t\tRain per hour:\t%.1f mm\n", *data.RainPerHour))
		must(fmt.Fprintf(w, prefix+"\t\tRain per day:\t%.1f mm\n", *data.RainPerDay))
	}
	if sliceContains(types, "Wind") {
		must(fmt.Fprintf(w, prefix+"\t\tWind:\t%d km/h (angle: %d °)\n", *data.WindStrength, *data.WindAngle))
		must(fmt.Fprintf(w, prefix+"\t\tGust:\t%d km/h (angle: %d °)\n", *data.GustStrength, *data.GustAngle))
	}
}

func must(_ int, err error) {
	if err != nil {
		panic(err)
	}
}

func formatTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

func sliceContains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
func f64OrNull(v *float64) string {
	if v == nil {
		return "null"
	}
	return fmt.Sprintf("%v", *v)
}

func intOrNull(v *int) string {
	if v == nil {
		return "null"
	}
	return fmt.Sprintf("%v", *v)
}
