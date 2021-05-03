package pricespider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
)

const (
	_huobiApiPrefix string = "https://api.huobi.pro/market/trade?symbol="
)

//获取火币最后一次成交价
func GetHuobiTradePrice(symbol string) (*PriceModel, error) {
	resp, err := http.Get(_huobiApiPrefix + strings.ToLower(symbol))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("huobiapi response error:" + string(bodyBytes))
		return nil, err
	}
	if resp.StatusCode != 200 {
		fmt.Println("huobiapi response error:" + string(bodyBytes))
		return nil, errors.New("get price from huobiapi error,response:" + string(bodyBytes))
	}

	var resultObj interface{}
	err = json.Unmarshal(bodyBytes, &resultObj)
	if err != nil {
		return nil, err
	}
	resultObjMap := resultObj.(map[string]interface{})
	if resultObjMap["status"].(string) != "ok" {
		fmt.Println("huobiapi response error:" + string(bodyBytes))
		return nil, errors.New("get price from huobiapi error,response:" + string(bodyBytes))
	}

	jsonArr := resultObjMap["tick"].(map[string]interface{})["data"].([]interface{})
	var ts int64 = 0
	var token0Sum float64 = 0
	var token1Sum float64 = 0
	for i := 0; i < len(jsonArr); i++ {
		m := jsonArr[i].(map[string]interface{})
		tempTs := int64(m["ts"].(float64))
		if tempTs > ts {
			ts = tempTs
		}
		amount := m["amount"].(float64)
		token0Sum += amount
		token1Sum += amount * m["price"].(float64)
	}
	finalPrice := math.Round(token1Sum/token0Sum*(math.Pow(10, 6))) / math.Pow(10, 6)
	ts = ts / 1000
	result := &PriceModel{
		Price:     finalPrice,
		Timestamp: ts,
	}

	fmt.Println(symbol + " last trade price: " + strconv.FormatFloat(result.Price, 'f', 6, 64) + ",   ts:" + strconv.FormatInt(result.Timestamp, 10) + "  from huobi api.")
	return result, nil
}
