package inventory

type InventoryResponse struct {
	Total          int             `json:"total"`
	Size           int             `json:"size"`
	Href           string          `json:"href"`
	Limit          int             `json:"limit"`
	InventoryItems []InventoryItem `json:"inventoryItems"`
}

type InventoryItem struct {
	Sku          string       `json:"sku"`
	Locale       string       `json:"locale"`
	Product      Product      `json:"product"`
	Condition    string       `json:"condition"`
	Availability Availability `json:"availability"`
}

type Product struct {
	Title       string              `json:"title"`
	Subtitle    string              `json:"subtitle,omitempty"`
	Description string              `json:"description"`
	Brand       string              `json:"brand"`
	Mpn         string              `json:"mpn,omitempty"`
	Epid        string              `json:"epid,omitempty"`
	Upc         []string            `json:"upc,omitempty"`
	Ean         []string            `json:"ean,omitempty"`
	Aspects     map[string][]string `json:"aspects,omitempty"`
	ImageUrls   []string            `json:"imageUrls"`
}

type Availability struct {
	ShipToLocationAvailability   ShipToLocationAvailability     `json:"shipToLocationAvailability"`
	PickupAtLocationAvailability []PickupAtLocationAvailability `json:"pickupAtLocationAvailability,omitempty"`
}

type ShipToLocationAvailability struct {
	Quantity int `json:"quantity"`
}

type PickupAtLocationAvailability struct {
	Quantity            int              `json:"quantity,omitempty"`
	MerchantLocationKey string           `json:"merchantLocationKey"`
	AvailabilityType    string           `json:"availabilityType,omitempty"`
	FulfillmentTime     *FulfillmentTime `json:"fulfillmentTime,omitempty"`
}

type FulfillmentTime struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}
