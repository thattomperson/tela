package eddn

import "time"

type Commodity struct {
	SchemaRef string `json:"$schemaRef"`
	Header    struct {
		UploaderID       string    `json:"uploaderID"`
		SoftwareName     string    `json:"softwareName"`
		SoftwareVersion  string    `json:"softwareVersion"`
		GatewayTimestamp time.Time `json:"gatewayTimestamp"`
	} `json:"header"`
	Message struct {
		SystemName  string    `json:"systemName"`
		StationName string    `json:"stationName"`
		MarketID    int       `json:"marketId"`
		Timestamp   time.Time `json:"timestamp"`
		Commodities []struct {
			Name          string  `json:"name"`
			MeanPrice     float64 `json:"meanPrice"`
			BuyPrice      float64 `json:"buyPrice"`
			Stock         int     `json:"stock"`
			SellPrice     float64 `json:"sellPrice"`
			Demand        int     `json:"demand"`
			StockBracket  int     `json:"stockBracket"`
			DemandBracket int     `json:"demandBracket"`
		} `json:"commodities"`
	} `json:"message"`
}
