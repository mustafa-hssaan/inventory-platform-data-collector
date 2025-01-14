package client_handlers

import (
	"encoding/json"
	"inventory-platform-data-collector/internal/ebay/auth/service"
	clients "inventory-platform-data-collector/internal/ebay/integration/factory"
	"inventory-platform-data-collector/internal/ebay/integration/models"
	model_GetBrowseItemRequest "inventory-platform-data-collector/internal/ebay/integration/models/GetBrowseItemRequest"
	model_ItemBrowseParams "inventory-platform-data-collector/internal/ebay/integration/models/ItemBrowseParams"

	"net/http"
	"strconv"
	"strings"
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
	itemID := r.URL.Query().Get("item_id")
	if itemID == "" {
		http.Error(w, "Missing item_id", http.StatusBadRequest)
		return
	}
	query := r.URL.Query()
	params := model_ItemBrowseParams.ItemBrowseParams{
		FieldGroups:                 query.Get("fieldGroups"),
		QuantityForShippingEstimate: query.Get("quantity_for_shipping_estimate"),
	}

	findingClient := h.clientFactory.NewFindingClient()
	result, err := findingClient.FindItemDetailsByID(r.Context(), model_GetBrowseItemRequest.GetBrowseItemRequest{
		ItemID: itemID,
	}, params, r.Header)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Item not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
