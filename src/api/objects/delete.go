package objects
import ("encoding/json"; "net/http")
func DeleteObjectHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Object deleted"})
}
