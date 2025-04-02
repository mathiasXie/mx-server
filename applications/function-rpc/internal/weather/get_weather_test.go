package weather

import (
	"testing"

	"github.com/mathiasXie/gin-web/config"
)

func TestGetWeather(t *testing.T) {
	config.Instance.LoadConfig("../../conf/function-rpc.yaml")
	weather, err := GetWeatherByCityName("北京")
	if err != nil {
		t.Errorf("GetWeather failed: %v", err)
	}
	t.Logf("weather: %v", weather)
}

func TestGetWeatherReport(t *testing.T) {
	config.Instance.LoadConfig("../../conf/function-rpc.yaml")
	report, err := GetWeatherReport("新加坡", "zh")
	if err != nil {
		t.Errorf("GetWeatherReport failed: %v", err)
	}
	t.Logf("report: %v", report)
}

/*
 curl -H "X-QW-Api-Key: 20ffef72ccd7432cac1e8c6838d63283" --compressed \
 'https://mn7p3ybh3m.re.qweatherapi.com/geo/v2/city/lookup?location=beij'
*/
