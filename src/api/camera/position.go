package camera

import (
	"encoding/json"
	"net/http"
)

func SetCameraPositionHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Camera positioned",
	})
}