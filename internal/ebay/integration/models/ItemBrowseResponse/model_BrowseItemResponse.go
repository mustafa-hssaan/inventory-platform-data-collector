package model_BrowseItemResponse

type BrowseItemResponse struct {
	ItemID string `json:"itemId"`
	Title  string `json:"title"`
	Price  struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"price"`
	Condition    string `json:"condition"`
	QuantitySold int    `json:"quantitySold"`
	Seller       struct {
		Username string `json:"username"`
	} `json:"seller"`
	ShippingOptions []struct {
		ShippingServiceCode string `json:"shippingServiceCode"`
		ShippingCost        struct {
			Value    float64 `json:"value"`
			Currency string  `json:"currency"`
		} `json:"shippingCost"`
		EstimatedDeliveryDays int `json:"estimatedDeliveryDays"`
	} `json:"shippingOptions"`
	ItemLocation struct {
		StateOrProvince string `json:"stateOrProvince"`
		Country         string `json:"country"`
		PostalCode      string `json:"postalCode"`
		Region          string `json:"region"`
	} `json:"itemLocation"`
	ItemSpecifics map[string][]string `json:"itemSpecifics"`
}
