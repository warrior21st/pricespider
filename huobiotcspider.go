package pricespider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"
)

//根据指定的完整huobi apiurl获取otc价格
func GetHuobiOTCPrice(huobiApiUrl string) (*PriceModel, error) {
	resp, err := http.Get(huobiApiUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("huobi otc api response error:" + string(bodyBytes))
		return nil, err
	}
	if resp.StatusCode != 200 {
		fmt.Println("huobi otc api response error:" + string(bodyBytes))
		return nil, errors.New("get otc price from huobi error,response:" + string(bodyBytes))
	}

	var resultObj interface{}
	err = json.Unmarshal(bodyBytes, &resultObj)
	if err != nil {
		return nil, err
	}
	resultObjMap := resultObj.(map[string]interface{})
	if resultObjMap["code"].(float64) != 200 {
		fmt.Println("huobi otc api response error:" + string(bodyBytes))
		return nil, errors.New("get otc price from huobi error,response:" + string(bodyBytes))
	}

	jsonArr := resultObjMap["data"].([]interface{})
	var tokenSum float64 = 0
	var currencySum float64 = 0
	for i := 0; i < 5 && i < len(jsonArr); i++ {
		m := jsonArr[i].(map[string]interface{})

		tokenAmount, err := strconv.ParseFloat(m["tradeCount"].(string), 64)
		if err != nil {
			return nil, err
		}
		price, err := strconv.ParseFloat(m["price"].(string), 64)
		if err != nil {
			return nil, err
		}
		tokenSum += tokenAmount
		currencySum += tokenAmount * price
	}
	finalPrice := math.Round(currencySum/tokenSum*(math.Pow(10, 6))) / math.Pow(10, 6)
	result := &PriceModel{
		Price:     finalPrice,
		Timestamp: time.Now().Unix(),
	}

	//fmt.Println(time.Now().Add(8*time.Hour).Format("2006-01-02 15:04:05") + "   " + "otc price: " + strconv.FormatFloat(result.Price, 'f', 6, 64) + ",   ts:" + strconv.FormatInt(result.Timestamp, 10) + "  from huobi otc api.")
	return result, nil

}
