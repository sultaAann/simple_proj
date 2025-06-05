package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type ParseError struct {
	Message string
	Err     error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

func parse() (map[string]interface{}, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", os.Getenv("API_URL"), nil)

	if err != nil {
		return nil, fmt.Errorf("Creating new Request: %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Cannot do request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response status code is not 200")
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("Reading error: %v", err)
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("Unmarshalling json error: %v", err)
	}

	return data["data"].(map[string]interface{}), nil
}

func GetData(filter string) (interface{}, error) {
	data, err := parse()
	if err != nil {
		return nil, &ParseError{"Error parsing", err}
	}

	if _, err := strconv.Atoi(filter); err == nil {
		year := filter
		v, ok := data[year]
		if !ok {
			return nil, fmt.Errorf("нет данных за %s год", year)
		}
		return v, nil
	}

	// Если filter — это просто ключ (месяц или год)
	if val, ok := data[filter]; ok {
		return val, nil
	}

	// Если filter — это пара "год месяц"
	fs := strings.Split(filter, " ")
	if len(fs) == 2 {
		year, month := fs[0], fs[1]
		v, ok := data[year]
		if !ok {
			return nil, fmt.Errorf("нет данных за %s год", year)
		}
		res, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("неожиданный формат данных за %s год", year)
		}
		if val, ok := res[month]; ok {
			return val, nil
		}
		return nil, fmt.Errorf("нет данных за месяц %s в %s году", month, year)
	}

	return nil, fmt.Errorf("непонятный формат фильтра: %s", filter)
}
