package client_handlers

import (
	"encoding/json"
	"inventory-platform-data-collector/internal/ebay/auth/service"
	clients "inventory-platform-data-collector/internal/ebay/integration/factory"
	"inventory-platform-data-collector/internal/ebay/integration/models"
	"net/http"
	"strconv"
)

type ClientHandler struct {
	clientFactory *clients.ClientFactory
}

func NewClientHandler(authService *service.Service, environment string) *ClientHandler {
	return &ClientHandler{
		clientFactory: clients.NewClientFactory(authService, environment),
	}
}

func (h *ClientHandler) FindingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	params := models.SearchParams{
		Keywords:   query.Get("keywords"),
		CategoryID: query.Get("category_id"),
		SortOrder:  query.Get("sort_order"),
	}

	if limit := query.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			params.Limit = l
		}
	}

	if offset := query.Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			params.Offset = o
		}
	}

	findingClient := h.clientFactory.NewFindingClient()
	result, err := findingClient.Search(r.Context(), params, r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *ClientHandler) MerchandisingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	request := models.GetDealsRequest{
		CategoryID: query.Get("category_id"),
	}

	if limit := query.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			request.Limit = l
		}
	}

	merchandisingClient := h.clientFactory.NewMerchandisingClient()
	result, err := merchandisingClient.GetDeals(r.Context(), request, r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *ClientHandler) ProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		http.Error(w, "Missing product_id", http.StatusBadRequest)
		return
	}

	productClient := h.clientFactory.NewProductClient()
	result, err := productClient.GetProduct(r.Context(), models.GetProductRequest{
		ProductID: productID,
	}, r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *ClientHandler) TradingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	itemID := r.URL.Query().Get("item_id")
	if itemID == "" {
		http.Error(w, "Missing item_id", http.StatusBadRequest)
		return
	}

	tradingClient := h.clientFactory.NewTradingClient()
	result, err := tradingClient.GetItem(r.Context(), models.GetItemRequest{
		ItemID: itemID,
	}, r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
