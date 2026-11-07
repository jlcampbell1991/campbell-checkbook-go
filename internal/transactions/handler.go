package transactions

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// importCSVHandler handles POST /transactions/import multipart form with a CSV file.
func importCSVHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form (max 10MB)
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
			return
		}
		// Expect form field "file"
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "file field is required", http.StatusBadRequest)
			return
		}
		defer file.Close()
		// Use repository ImportCSV
		count, err := repo.ImportCSV(r.Context(), file)
		if err != nil {
			http.Error(w, fmt.Sprintf("import failed: %v", err), http.StatusBadRequest)
			return
		}
		// Return JSON summary
		w.Header().Set("Content-Type", "application/json")
		summary := map[string]int{"imported": count}
		_ = json.NewEncoder(w).Encode(summary)
	}
}

// For completeness, we expose a simple health handler for transactions (optional).
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "transactions ok"})
}
