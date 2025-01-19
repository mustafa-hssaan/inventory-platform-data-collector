package client_handlers

import (
	"encoding/json"
	"inventory-platform-data-collector/internal/ebay/auth/service"
	clients "inventory-platform-data-collector/internal/ebay/integration/factory"
	"inventory-platform-data-collector/internal/ebay/integration/models/analytics"
	"inventory-platform-data-collector/internal/ebay/integration/models/browseitem"
	"inventory-platform-data-collector/internal/ebay/integration/models/pagination"
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
	params := browseitem.ItemBrowseParams{
		FieldGroups:                 query.Get("fieldGroups"),
		QuantityForShippingEstimate: query.Get("quantity_for_shipping_estimate"),
	}

	findingClient := h.clientFactory.NewItemBrowseClient()
	result, err := findingClient.FindItemDetailsByID(r.Context(), browseitem.GetBrowseItemRequest{
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
func (h *ClientHandler) GetTrafficReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query()
	params := analytics.TrafficReportParam{
		Dimension: query.Get("dimension"),
		Filter:    query.Get("filter"),
		Metric:    query.Get("metric"),
		Sort:      query.Get("sort"),
	}

	analyticsClient := h.clientFactory.NewAnalyticsClient()
	result, err := analyticsClient.GetTrafficReport(r.Context(), params, r.Header)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "No reports found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
func (h *ClientHandler) GetInventoryItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query()
	params := pagination.PaginationParams{}
	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			params.Limit = &limit
		} else {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}
	if offsetStr := query.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			params.Offset = &offset
		} else {
			http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
			return
		}
	}

	inventoryItemsClient := h.clientFactory.NewInvetoryItemsClient()
	result, err := inventoryItemsClient.GetInventoryItemsDetails(r.Context(), params, r.Header)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "No reports found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
