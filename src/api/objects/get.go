package objects
import ("encoding/json"; "net/http")
func GetObjectHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"name": "test", "type": "cube"})
}
