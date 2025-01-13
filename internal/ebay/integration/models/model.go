package models

import "time"

type Item struct {
	ItemID          string            `json:"itemId"`
	Title           string            `json:"title"`
	Price           Price             `json:"price"`
	Condition       string            `json:"condition"`
	QuantitySold    int               `json:"quantitySold"`
	WatchCount      int               `json:"watchCount"`
	ViewCount       int               `json:"viewCount"`
	LastSold        time.Time         `json:"lastSold,omitempty"`
	Seller          SellerInfo        `json:"seller"`
	Location        Location          `json:"location"`
	ShippingOptions []Shipping        `json:"shippingOptions"`
	Specifics       map[string]string `json:"itemSpecifics"`
	PriceHistory    []PricePoint      `json:"priceHistory"`
	Reviews         []Review          `json:"reviews"`
	SalesVelocity   float64           `json:"salesVelocity"`
	MarketPosition  MarketPosition    `json:"marketPosition"`
}
type SellerInfo struct {
	Username        string  `json:"username"`
	FeedbackScore   int     `json:"feedbackScore"`
	PositivePercent float64 `json:"positivePercent"`
}
type Location struct {
	Country    string `json:"country"`
	PostalCode string `json:"postalCode"`
	Region     string `json:"region"`
}
type Shipping struct {
	Service      string `json:"service"`
	Cost         Price  `json:"cost"`
	DeliveryDays int    `json:"deliveryDays"`
}
type PricePoint struct {
	Price     Price     `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}
type Review struct {
	Rating  int       `json:"rating"`
	Comment string    `json:"comment"`
	Date    time.Time `json:"date"`
}
type MarketPosition struct {
	CompetitorPriceMin Price   `json:"competitorPriceMin"`
	CompetitorPriceMax Price   `json:"competitorPriceMax"`
	CompetitorPriceAvg Price   `json:"competitorPriceAvg"`
	RelativePosition   float64 `json:"relativePosition"`
}

type Price struct {
	Value    string `json:"value" xml:"Value"`
	Currency string `json:"currency" xml:"Currency"`
}

type Product struct {
	ID          string   `json:"productId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ImageURLs   []string `json:"imageUrls"`
	Brand       string   `json:"brand"`
	Categories  []string `json:"categories"`
}
type GetProductRequest struct {
	ProductID string
}

type ProductResponse struct {
	Product Product `json:"product"`
}
type GetDealsRequest struct {
	CategoryID string
	Limit      int
}

type DealsResponse struct {
	Items []Item `json:"items"`
}

type SearchParams struct {
	Q          string `json:"q"`
	CategoryID string `json:"categoryId,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Offset     int    `json:"offset,omitempty"`
	SortOrder  string `json:"sortOrder,omitempty"`
}
type GetItemRequest struct {
	ItemID string
}

type GetItemResponse struct {
	Item Item `xml:"Item"`
}

type SearchResponse struct {
	Items      []Item     `json:"itemSummaries"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
