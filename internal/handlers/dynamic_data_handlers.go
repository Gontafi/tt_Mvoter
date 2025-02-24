package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tt/internal/services"
	"tt/pkg/utils"
)

type DynamicDataHandler struct {
	service services.DynamicTableServiceInterface
}

func NewDynamicDataHandler(service services.DynamicTableServiceInterface) *DynamicDataHandler {
	return &DynamicDataHandler{service: service}

}

func (h *DynamicDataHandler) CreateTable(w http.ResponseWriter, r *http.Request) {
	var request struct {
		TableName string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := utils.ParseUserIDJWTInHandler(w, r)
	if userID == 0 {
		return
	}

	tableID, err := h.service.CreateTable(r.Context(), userID, request.TableName)
	if err != nil {
		utils.SendError(w, "Failed to create table", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"table_id": strconv.FormatInt(tableID, 10)})
}

func (h *DynamicDataHandler) CreateRow(w http.ResponseWriter, r *http.Request) {
	tableIDString := r.PathValue("tableID")
	if tableIDString == "" {
		utils.SendError(w, "missing table_id", http.StatusBadRequest)
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		utils.SendError(w, "invalid input", http.StatusBadRequest)
		return
	}

	tableID, err := strconv.ParseInt(tableIDString, 10, 64)
	if err != nil {
		utils.SendError(w, "invalid input", http.StatusBadRequest)
		return
	}
	rowID, err := h.service.AddRow(r.Context(), tableID, data)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"row_id": strconv.FormatInt(rowID, 10)})
}

func (h *DynamicDataHandler) GetRows(w http.ResponseWriter, r *http.Request) {
	tableIDString := r.PathValue("tableID")
	if tableIDString == "" {
		utils.SendError(w, "missing table_id", http.StatusBadRequest)
		return
	}

	tableID, err := strconv.ParseInt(tableIDString, 10, 64)
	if err != nil {
		utils.SendError(w, "invalid input", http.StatusBadRequest)
		return
	}

	rows, err := h.service.GetAllRows(r.Context(), tableID)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rows)
}

func (h *DynamicDataHandler) UpdateRow(w http.ResponseWriter, r *http.Request) {
	tableIDString := r.PathValue("tableID")
	if tableIDString == "" {
		utils.SendError(w, "missing table_id", http.StatusBadRequest)
		return
	}

	tableID, err := strconv.ParseInt(tableIDString, 10, 64)
	if err != nil {
		utils.SendError(w, "invalid input", http.StatusBadRequest)
		return
	}

	rowIDString := r.PathValue("rowID")
	if rowIDString == "" {
		utils.SendError(w, "missing row_id", http.StatusBadRequest)
		return
	}

	rowID, err := strconv.ParseInt(rowIDString, 10, 64)
	if err != nil {
		utils.SendError(w, "invalid input", http.StatusBadRequest)
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		utils.SendError(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateTableRow(r.Context(), tableID, rowID, data); err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *DynamicDataHandler) DeleteRow(w http.ResponseWriter, r *http.Request) {
	tableIDString := r.PathValue("tableID")
	if tableIDString == "" {
		utils.SendError(w, "missing table_id", http.StatusBadRequest)
		return
	}

	tableID, err := strconv.ParseInt(tableIDString, 10, 64)
	if err != nil {
		utils.SendError(w, "invalid input", http.StatusBadRequest)
		return
	}

	rowIDString := r.PathValue("rowID")
	if rowIDString == "" {
		utils.SendError(w, "missing row_id", http.StatusBadRequest)
		return
	}

	rowID, err := strconv.ParseInt(rowIDString, 10, 64)
	if err != nil {
		utils.SendError(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.RemoveTableRow(r.Context(), tableID, rowID); err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *DynamicDataHandler) GetTables(w http.ResponseWriter, r *http.Request) {
	userID := utils.ParseUserIDJWTInHandler(w, r)
	if userID == 0 {
		return
	}

	tables, err := h.service.GetTables(r.Context(), userID)
	if err != nil {
		utils.SendError(w, "Failed to retrieve tables", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tables)
	if err != nil {
		utils.SendError(w, "Failed to retrieve tables", http.StatusInternalServerError)
		return
	}

}
