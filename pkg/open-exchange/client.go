package open_exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/devdammit/shekel/pkg/currency"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	accessKey string
}

func NewClient(accessKey string) *Client {
	return &Client{accessKey: accessKey}
}

type HistoricalRates struct {
	Date  datetime.Date             `json:"date"`
	Base  currency.Code             `json:"base"`
	Rates map[currency.Code]float64 `json:"rates"`
}

func (c *Client) GetByDate(base currency.Code, codes []currency.Code, date datetime.Date) (*HistoricalRates, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	strCodes := make([]string, len(codes))
	for i, code := range codes {
		strCodes[i] = code.String()
	}

	url := fmt.Sprintf("http://api.exchangeratesapi.io/api/v1/%s?access_key=%s&base=%s&symbols=%s", date.Format("2006-01-02"), c.accessKey, base.String(), strings.Join(strCodes, ","))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var historicalRates HistoricalRates

	err = json.NewDecoder(res.Body).Decode(&historicalRates)
	if err != nil {
		return nil, err
	}

	return &historicalRates, nil
}
