package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"net/url"
	"os"
)

// RequestInfo представляет из себя структуру с токеном, url для запроса и объектом http.Client
type RequestInfo struct {
	apiToken   string
	requestUrl string
	client     *http.Client
}

func main() {
	// загружаем в окружение переменные из .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("Ошибка загрузки переменных окружений: ", err)
		return
	}
	apiToken := os.Getenv("API_TOKEN")

	requestUrl := "https://development.kpi-drive.ru/_api/facts/save_fact"

	// данные для отправки
	data := url.Values{
		"period_start":            {"2024-05-01"},
		"period_end":              {"2024-05-31"},
		"period_key":              {"month"},
		"indicator_to_mo_id":      {"227373"},
		"indicator_to_mo_fact_id": {"0"},
		"value":                   {"1"},
		"fact_time":               {"2024-05-31"},
		"is_plan":                 {"0"},
		"auth_user_id":            {"40"},
		"comment":                 {"buffer BendaTestTestTestTestTest50"},
	}

	// создаем объект структуры
	reqInfo := RequestInfo{apiToken: apiToken, requestUrl: requestUrl, client: &http.Client{}}

	for i := 0; i < 10; i++ {
		if err := reqInfo.sendRequest(data); err != nil {
			fmt.Println("Ошибка! Текст ошибки: ", err)
			return
		}
	}
}

func (reqInfo *RequestInfo) sendRequest(data url.Values) error {
	buffer := bytes.NewBufferString(data.Encode())
	req, err := http.NewRequest(http.MethodPost, reqInfo.requestUrl, buffer) // не совершаем запрос, а только создаем его, чтобы можно было поработать с заголовками
	if err != nil {
		return errors.New(fmt.Sprintf("Ошибка при создании запроса: %s", err))
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")         // явно укажем формат отправляемых данных
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", reqInfo.apiToken)) // указываем Bearer токен

	resp, err := reqInfo.client.Do(req)

	defer resp.Body.Close() // лучше закрыть тело запроса, чтобы не допустить утечки памяти

	// сервер сохраняет данные даже при 500 ошибке, поэтому проверяем только на актуальные 400-ые ошибки
	switch resp.StatusCode {
	case http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden:
		return errors.New(fmt.Sprintf("Запрос не выполнен: %d", resp.StatusCode))
	default:
		fmt.Println("Успех: ", resp.Status)
	}

	return nil
}
