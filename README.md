# pricespider
    import "github.com/warrior21st/pricespider"
## sample code
    price,err := pricespider.GetHuobiTradePrice("btcusdt")
	price,err = pricespider.GetBinanceTradePrice("ethusdt")
	avgPrice,err := pricespider.GetBinanceAvgPrice("ethusdt")
