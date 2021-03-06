package load

import (
	"fmt"
	"net/http"
	"time"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"github.com/mounlion/markdownwatcher/bot"
	"github.com/mounlion/markdownwatcher/model"
	"log"
	"github.com/mounlion/markdownwatcher/config"
)

var (
	headers = map[string]string{
		"X-Requested-With": "XMLHttpRequest",
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
	}
)

func Catalog(lastProductIndex int, cityKey string) string {
	var html string

	countError := 0;
	countErrorForFailed := 20;

	headers["Cookie"] = getCookei(cityKey);

	for {
		if *config.Config.Logger {log.Printf("Fetch offset %d", lastProductIndex)}
		result, statusCode := fetchCatalog(lastProductIndex)
		if statusCode != 200 {
			countError++
			isFailed := countError > countErrorForFailed;
			if isFailed {
				if *config.Config.Logger {log.Printf("Fetch failed. Status code: %d", statusCode)}
				message := fmt.Sprintf("<b>Обнаружена проблема</b>\n\nСтатус ответа сервера: <code>%d</code>", statusCode)
				bot.SendServiceMessage(message)
				break
			} else {
				var timeOfSleep time.Duration = 15;
				if *config.Config.Logger {log.Printf("Fetch failed. Status code: %d. Sleep: %d", statusCode, timeOfSleep)}
				time.Sleep(time.Second*timeOfSleep)
			}
		}
		html += result.HTML
		if result.IsNextLoadAvailable {
			lastProductIndex = result.LastProductIndex
		} else {
			if *config.Config.Logger {log.Printf("All fetch end")}
			break
		}
	}

	return html
}

func fetchCatalog (offset int)  (model.JSONObject, int)  {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", "https://www.dns-shop.ru/catalogMarkdown/markdown/update/?offset="+strconv.Itoa(offset), nil)
	if err != nil {log.Printf("NewRequest error")}
	for key, val := range headers {
		req.Header.Add(key, val)
	}
	resp, err := netClient.Do(req)
	defer resp.Body.Close()
	if err != nil {log.Printf("Http request error")}
	buf, _ := ioutil.ReadAll(resp.Body)
	jsonObj := model.JSONObject{}
	if resp.StatusCode == 200 {
		json.Unmarshal(buf, &jsonObj)
		return jsonObj, resp.StatusCode
	}

	fmt.Println(string(buf))
	return jsonObj, resp.StatusCode
}

func getCookei(cityName string) string {
    return "city_path="+cityName+"; __COOKIE__"
}

