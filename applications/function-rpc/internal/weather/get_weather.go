package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mathiasXie/gin-web/config"
	"github.com/mathiasXie/gin-web/pkg/logger"
)

// CitySearchResult 定义城市搜索结果的结构体
type CitySearchResult struct {
	Code     string `json:"code"`
	Location []struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Lat       string `json:"lat"`
		Lon       string `json:"lon"`
		Adm2      string `json:"adm2"`
		Adm1      string `json:"adm1"`
		Country   string `json:"country"`
		Tz        string `json:"tz"`
		UtcOffset string `json:"utcOffset"`
		IsDst     string `json:"isDst"`
		Type      string `json:"type"`
		Rank      string `json:"rank"`
		FxLink    string `json:"fxLink"`
	} `json:"location"`
	Refer struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	} `json:"refer"`
}

// NowWeather 定义实时天气信息的结构体
type NowWeather struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	FxLink     string `json:"fxLink"`
	Now        struct {
		ObsTime   string `json:"obsTime" text:"观测时间"`
		Temp      string `json:"temp" text:"温度"`
		FeelsLike string `json:"feelsLike" text:"体感温度"`
		Icon      string `json:"icon" text:"天气图标"`
		Text      string `json:"text" text:"天气状况"`
		Wind360   string `json:"wind360" text:"风向360角度"`
		WindDir   string `json:"windDir" text:"风向"`
		WindScale string `json:"windScale" text:"风力"`
		WindSpeed string `json:"windSpeed" text:"风速"`
		Humidity  string `json:"humidity" text:"湿度"`
		Precip    string `json:"precip" text:"降水量"`
		Pressure  string `json:"pressure" text:"气压"`
		Vis       string `json:"vis" text:"能见度"`
		Cloud     string `json:"cloud" text:"云量"`
		Dew       string `json:"dew" text:"露点温度"`
	} `json:"now"`
	Refer struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	} `json:"refer"`
}

var weatherCodeMap = map[string]string{
	"100": "晴", "101": "多云", "102": "少云", "103": "晴间多云", "104": "阴",
	"150": "晴", "151": "多云", "152": "少云", "153": "晴间多云",
	"300": "阵雨", "301": "强阵雨", "302": "雷阵雨", "303": "强雷阵雨", "304": "雷阵雨伴有冰雹",
	"305": "小雨", "306": "中雨", "307": "大雨", "308": "极端降雨", "309": "毛毛雨/细雨",
	"310": "暴雨", "311": "大暴雨", "312": "特大暴雨", "313": "冻雨", "314": "小到中雨",
}

// GetLocationID 根据城市名称获取 location id
func GetLocationID(cityName string) (string, string, error) {
	// 构建城市搜索请求 URL

	params := url.Values{}
	params.Add("location", cityName)
	params.Add("key", config.Instance.FunctionRPC.Weather.Qweather.ApiKey)

	fullURL := fmt.Sprintf("https://%s/geo/v2/city/lookup?%s", config.Instance.FunctionRPC.Weather.Qweather.ApiHost, params.Encode())

	// 发送 HTTP 请求
	// 添加请求头
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return "", "", err
	}
	//req.Header.Add("X-QW-Api-Key", config.Instance.FunctionRPC.Weather.Qweather.ApiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// 解析 JSON 数据
	var citySearchResult CitySearchResult
	err = json.Unmarshal(body, &citySearchResult)
	if err != nil {
		return "", "", err
	}

	if len(citySearchResult.Location) > 0 {
		return citySearchResult.Location[0].ID, citySearchResult.Location[0].FxLink, nil
	}

	return "", "", fmt.Errorf("未找到城市: %s", cityName)
}

// GetNowWeather 根据 location id 获取实时天气信息
func GetNowWeather(locationID string) (*NowWeather, error) {
	// 构建实时天气请求 URL
	params := url.Values{}
	params.Add("location", locationID)
	params.Add("key", config.Instance.FunctionRPC.Weather.Qweather.ApiKey)
	fullURL := fmt.Sprintf("https://%s/v7/weather/now?%s", config.Instance.FunctionRPC.Weather.Qweather.ApiHost, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	//req.Header.Add("X-QW-Api-Key", config.Instance.FunctionRPC.Weather.Qweather.ApiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// 解析 JSON 数据
	var nowWeather NowWeather
	err = json.Unmarshal(body, &nowWeather)
	if err != nil {
		return nil, err
	}

	return &nowWeather, nil
}

// GetWeatherByCityName 根据城市名称获取实时天气信息
func GetWeatherByCityName(cityName string) (*NowWeather, error) {
	locationID, _, err := GetLocationID(cityName)
	if err != nil {
		logger.CtxError(context.Background(), "GetWeatherByCityName", "err", err)
		return nil, err
	}

	return GetNowWeather(locationID)
}

func GetWeatherReport(cityName, lang string) (string, error) {
	locationID, fxLink, err := GetLocationID(cityName)
	if err != nil {
		return "", err
	}
	currentAbstract, currentWeather, weekWeather, err := parseWeatherInfo(locationID, fxLink)
	if err != nil {
		return "", err
	}

	var weatherReport []string

	weatherReport = append(weatherReport, fmt.Sprintf("根据下列数据，用%s回应用户的查询天气请求:%s\n", lang, cityName))

	if currentAbstract != "" {
		weatherReport = append(weatherReport, fmt.Sprintf("当前天气: %s\n", currentAbstract))
	}
	if currentWeather != "" {
		weatherReport = append(weatherReport, fmt.Sprintf("当前天气参数: %s\n", currentWeather))
	}
	if weekWeather != "" {
		weatherReport = append(weatherReport, fmt.Sprintf("未来7天天气: %s\n", weekWeather))
	}
	weatherReport = append(weatherReport, "(确保只报告指定单日的天气情况，除非未来会出现异常天气；或者用户明确要求想要了解多日天气，如果未指定，默认报告今天的天气。")
	weatherReport = append(weatherReport, "参数为0的值不需要报告给用户，每次都报告体感温度，根据语境选择合适的参数内容告知用户，并对参数给出相应评价)")

	return strings.Join(weatherReport, ""), nil
}

func parseWeatherInfo(locationID string, fxLink string) (string, string, string, error) {

	var currentAbstract, currentWeather, weekWeather string

	doc, err := fetchWeatherPage(fxLink)
	if err == nil {
		// 当前天气摘要
		doc.Find("div.current-abstract").Each(func(i int, s *goquery.Selection) {
			// 获取元素的文本内容
			currentAbstract = strings.ReplaceAll(strings.TrimSpace(s.Text()), "\n", "")
		})
		// 当前天气参数
		doc.Find(".c-city-weather-current .current-basic .current-basic___item").Each(func(i int, s *goquery.Selection) {
			ps := s.Find("p")
			if ps.Length() == 2 {
				score := ps.Eq(0).Text()
				title := ps.Eq(1).Text()
				currentWeather += fmt.Sprintf("%s: %s\n", title, score)
			}
		})
		// 未来7天天气
		doc.Find(".city-forecast-tabs__row").Each(func(i int, s *goquery.Selection) {
			date := s.Find(".date p").Eq(0).Text()
			low := s.Find(".temp").Eq(1).Text()
			high := s.Find(".temp").Eq(0).Text()
			src, ok := s.Find(".date-bg .icon").Attr("src")
			weather := "未知"
			if ok {
				code := strings.Split(strings.Split(src, "/")[len(strings.Split(src, "/"))-1], ".")[0]
				weather = weatherCodeMap[code]
			}
			weekWeather += fmt.Sprintf("%s:%s到%s,%s\n", date, low, high, weather)
		})
	}
	if len(currentWeather) == 0 {
		//遍历nowWeather的now的属性，用反射获取text的属性值
		nowWeather, err := GetNowWeather(locationID)
		if err != nil {
			return "", "", "", err
		}
		value := reflect.ValueOf(nowWeather.Now)
		typ := reflect.TypeOf(nowWeather.Now)
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			fieldType := typ.Field(i)
			textTag := fieldType.Tag.Get("text")
			if textTag != "" {
				currentWeather += fmt.Sprintf("%s: %v;", textTag, field.Interface())
			}
		}
	}

	return currentAbstract, currentWeather, weekWeather, nil
}

func fetchWeatherPage(url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("请求失败，状态码: %d\n", resp.StatusCode)
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
