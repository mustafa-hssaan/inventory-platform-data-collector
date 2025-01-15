package browseitem

type BrowseItemResponse struct {
	ItemID                    string                    `json:"itemId"`
	Title                     string                    `json:"title"`
	CategoryPath              string                    `json:"categoryPath"`
	Price                     Price                     `json:"price"`
	Condition                 string                    `json:"condition"`
	Seller                    Seller                    `json:"seller"`
	ItemLocation              ItemLocation              `json:"itemLocation"`
	ItemSpecifics             map[string][]string       `json:"itemSpecifics"`
	EstimatedAvailabilities   []EstimatedAvailabilities `json:"estimatedAvailabilities"`
	BuyingOptions             []string                  `json:"buyingOptions"`
	TopRatedBuyingExperience  bool                      `json:"topRatedBuyingExperience"`
	ImmediatePay              bool                      `json:"immediatePay"`
	EnabledForGuestCheckout   bool                      `json:"enabledForGuestCheckout"`
	EligibleForInlineCheckout bool                      `json:"eligibleForInlineCheckout"`
	AdultOnly                 bool                      `json:"adultOnly"`
	CategoryId                string                    `json:"categoryId"`
	ListingMarketplaceId      string                    `json:"listingMarketplaceId"`
	ReturnTerms               ReturnTerms               `json:"returnTerms"`
}

type Price struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Seller struct {
	Username           string `json:"username"`
	FeedbackPercentage string `json:"feedbackPercentage"`
	FeedbackScore      int    `json:"feedbackScore"`
}

type ShippingOption struct {
	ShippingServiceCode   string `json:"shippingServiceCode"`
	ShippingCost          Price  `json:"shippingCost"`
	EstimatedDeliveryDays int    `json:"estimatedDeliveryDays"`
}

type ItemLocation struct {
	StateOrProvince string `json:"stateOrProvince"`
	Country         string `json:"country"`
	PostalCode      string `json:"postalCode"`
	City            string `json:"city"`
}
type EstimatedAvailabilities struct {
	DeliveryOption              []string `json:"deliveryOptions"`
	EstimatedAvailabilityStatus string   `json:"estimatedAvailabilityStatus"`
	EstimatedAvailableQuantity  int      `json:"estimatedAvailableQuantity"`
	EstimatedSoldQuantity       int      `json:"estimatedSoldQuantity"`
	EstimatedRemainingQuantity  int      `json:"estimatedRemainingQuantity"`
}
type PaymentMethods struct {
	PaymentMethodType  string                `json:"paymentMethodType"`
	PaymentMethodBrand []PaymentMethodBrands `json:"paymentMethodBrands"`
}
type PaymentMethodBrands struct {
	PaymentMethodBrandType string `json:"paymentMethodBrandType"`
}
type ReturnTerms struct {
	ReturnsAccepted bool `json:"returnsAccepted"`
}
