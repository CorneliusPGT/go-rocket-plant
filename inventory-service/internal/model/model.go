package model

import "time"

type Part struct {
	UUID          string            `bson:"uuid"`
	Name          string            `bson:"name"`
	Description   string            `bson:"description"`
	Price         float64           `bson:"price"`
	StockQuantity int64             `bson:"stock_quantity"`
	Category      int32             `bson:"category"`
	Dimensions    *Dimensions       `bson:"dimensions"`
	Manufacter    *Manufacter       `bson:"manufacter"`
	Tags          []string          `bson:"tags"`
	Metadata      map[string]string `bson:"metadata"`
	CreatedAt     time.Time         `bson:"created_at"`
	UpdatedAt     time.Time         `bson:"updated_at"`
}

type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacter struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}
