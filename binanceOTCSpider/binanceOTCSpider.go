package binanceOTCSpider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

const (
	apiUrl     = "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"
	parasConst = `{"proMerchantAds":false,"page":1,"rows":10,"payTypes":[],"countries":[],"publisherType":"merchant","tradeType":"{tradeType}","asset":"{asset}","fiat":"{fiat}"}`
)

func GetOTCPrice(token string, fiat string, avgCount int) (price decimal.Decimal, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(genParams(token, fiat)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("binance otc api response error:" + string(bodyBytes))
		return nil, err
	}
	if resp.StatusCode != 200 {
		fmt.Println("binance otc api response error:" + string(bodyBytes))
		return nil, errors.New("get price from binance otc api error,response:" + string(bodyBytes))
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

	//fmt.Println(time.Now().Add(8*time.Hour).Format("2006-01-02 15:04:05") + "   " + symbol + " last trade price: " + strconv.FormatFloat(result.Price, 'f', 6, 64) + ",   ts:" + strconv.FormatInt(result.Timestamp, 10) + "  from binance api.")
	return result, nil

}

func genParams(token string, fiat string) string {
	paras := strings.ReplaceAll(parasConst, "{tradeType}", "BUY")
	paras = strings.ReplaceAll(paras, "{asset}", strings.ToUpper(token))
	paras = strings.ReplaceAll(paras, "{fiat}", strings.ToUpper(fiat))

	return paras
}
