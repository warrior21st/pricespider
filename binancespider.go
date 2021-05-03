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

const (
	_binanceTradePriceApiPrefix = "https://api.binance.com/api/v3/ticker/price?symbol="
	_binanceAvgPriceApiPrefix   = "https://api.binance.com/api/v3/avgPrice?symbol="
)

//获取币安最后一次成交价
func GetBinanceTradePrice(symbol string) (*PriceModel, error) {
	resp, err := http.Get(_binanceTradePriceApiPrefix + strings.ToUpper(symbol))
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
		Price:     price,
		Timestamp: time.Now().Unix(),
	}

	fmt.Println(symbol + " last trade price: " + strconv.FormatFloat(result.Price, 'f', 6, 64) + ",   ts:" + strconv.FormatInt(result.Timestamp, 10) + "  from binance api.")
	return result, nil
}

//获取币安均价
func GetBinanceAvgPrice(symbol string) (*BinanceAvgPrice, error) {
	resp, err := http.Get(_binanceAvgPriceApiPrefix + strings.ToUpper(symbol))
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
	mins := int(resultObjMap["mins"].(float64))
	if err != nil {
		return nil, err
	}
	result := &BinanceAvgPrice{
		Price:     price,
		Timestamp: time.Now().Unix(),
		Mins:      mins,
	}

	fmt.Println(symbol + " avg price: " + strconv.FormatFloat(result.Price, 'f', 6, 64) + ",   ts:" + strconv.FormatInt(result.Timestamp, 10) + ",   mins:" + strconv.Itoa(result.Mins) + "from binance api.")
	return result, nil
}
