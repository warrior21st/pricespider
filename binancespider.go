package pricespider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	_binanceApiPrefix = "https://api.binance.com/api/v3/ticker/price?symbol="
)

func GetTradePriceFromBinance(symbol string) (*PriceModel, error) {
	resp, err := http.Get(_binanceApiPrefix + strings.ToUpper(symbol))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("binanceapi response error:" + string(bodyBytes))
		return nil, err
	}
	if resp.StatusCode != 200 {
		fmt.Println("binanceapi response error:" + string(bodyBytes))
		return nil, errors.New("get price from binanceapi error,response:" + string(bodyBytes))
	}

	var resultObj interface{}
	err = json.Unmarshal(bodyBytes, &resultObj)
	if err != nil {
		return nil, err
	}
	resultObjMap := resultObj.(map[string]interface{})
	price, err := strconv.ParseFloat(resultObjMap["price"].(string), 64)
	if err != nil {
		return nil, err
	}

	result := &PriceModel{
		price:     price,
		timestamp: time.Now().Unix(),
	}

	fmt.Println(symbol + " last trade price: " + strconv.FormatFloat(result.price, 'f', 6, 64) + ",   ts:" + strconv.FormatInt(result.timestamp, 10) + "  from binance api.")
	return result, nil
}
