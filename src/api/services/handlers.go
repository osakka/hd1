package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"holodeck1/logging"
	"holodeck1/service"
)

type Handler struct {
	serviceRegistry *service.Registry
}

func NewHandler(serviceRegistry *service.Registry) *Handler {
	return &Handler{
		serviceRegistry: serviceRegistry,
	}
}

func (h *Handler) ListServices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Parse query parameters
	filter := &service.ServiceFilter{
		Type:   r.URL.Query().Get("type"),
		Status: r.URL.Query().Get("status"),
	}
	
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	}
	
	if offset := r.URL.Query().Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			filter.Offset = o
		}
	}
	
	// Set default limit if not provided
	if filter.Limit == 0 {
		filter.Limit = 50
	}
	
	services, err := h.serviceRegistry.ListServices(ctx, filter)
	if err != nil {
		logging.Error("failed to list services", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success":  true,
		"services": services,
		"pagination": map[string]interface{}{
			"limit":    filter.Limit,
			"offset":   filter.Offset,
			"has_more": len(services) >= filter.Limit,
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RegisterService(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var req service.RegisterServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate required fields
	if req.Name == "" || req.Type == "" || req.Endpoint == "" {
		http.Error(w, "Name, type, and endpoint are required", http.StatusBadRequest)
		return
	}
	
	service, err := h.serviceRegistry.RegisterService(ctx, &req)
	if err != nil {
		logging.Error("failed to register service", map[string]interface{}{
			"service_name": req.Name,
			"error":        err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success": true,
		"service": service,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetService(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	serviceID, err := uuid.Parse(vars["serviceId"])
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}
	
	service, err := h.serviceRegistry.GetService(ctx, serviceID)
	if err != nil {
		if err.Error() == "service not found" {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		logging.Error("failed to get service", map[string]interface{}{
			"service_id": serviceID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success": true,
		"service": service,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UpdateService(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	serviceID, err := uuid.Parse(vars["serviceId"])
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}
	
	var req service.UpdateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	service, err := h.serviceRegistry.UpdateService(ctx, serviceID, &req)
	if err != nil {
		if err.Error() == "service not found" {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		logging.Error("failed to update service", map[string]interface{}{
			"service_id": serviceID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success": true,
		"service": service,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteService(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	serviceID, err := uuid.Parse(vars["serviceId"])
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}
	
	err = h.serviceRegistry.DeleteService(ctx, serviceID)
	if err != nil {
		if err.Error() == "service not found" {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		logging.Error("failed to delete service", map[string]interface{}{
			"service_id": serviceID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CheckServiceHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	serviceID, err := uuid.Parse(vars["serviceId"])
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}
	
	health, err := h.serviceRegistry.CheckServiceHealth(ctx, serviceID)
	if err != nil {
		logging.Error("failed to check service health", map[string]interface{}{
			"service_id": serviceID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success": true,
		"health":  health,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetServiceHealthHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	serviceID, err := uuid.Parse(vars["serviceId"])
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}
	
	healthHistory, err := h.serviceRegistry.GetServiceHealth(ctx, serviceID)
	if err != nil {
		logging.Error("failed to get service health history", map[string]interface{}{
			"service_id": serviceID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success": true,
		"health":  healthHistory,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}