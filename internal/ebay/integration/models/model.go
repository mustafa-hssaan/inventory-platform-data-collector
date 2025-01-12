package models

type Item struct {
	ID              string   `json:"itemId" xml:"ItemID"`
	Title           string   `json:"title" xml:"Title"`
	Price           Price    `json:"price" xml:"Price"`
	Categories      []string `json:"categoryIds" xml:"CategoryIDs"`
	Condition       string   `json:"condition" xml:"Condition"`
	Description     string   `json:"description" xml:"Description"`
	ImageURLs       []string `json:"imageUrls" xml:"ImageURLs"`
	Location        string   `json:"location" xml:"Location"`
	SellerUsername  string   `json:"sellerUsername" xml:"SellerUsername"`
	ShippingOptions []string `json:"shippingOptions" xml:"ShippingOptions"`
}

type Price struct {
	Value    float64 `json:"value" xml:"Value"`
	Currency string  `json:"currency" xml:"Currency"`
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
	Keywords   string `json:"keywords"`
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
	Items      []Item     `json:"items"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
