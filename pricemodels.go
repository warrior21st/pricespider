package pricespider

type PriceModel struct {
	Price     float64
	Timestamp int64
}

type BinanceAvgPrice struct {
	Price     float64
	Timestamp int64
	Mins      int
}
