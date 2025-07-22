package upload_session

import (
	pfileUploadService "docTrack/file_upload_service"
	pglobalConfigs "docTrack/global_configs"
	p_uploadSessionService "docTrack/services/upload_session"
	"encoding/json"
	"net/http"
)

func GetInitUploadSessionHandler(fUploadData *pfileUploadService.FileUploadServiceInfo) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Close the body when we’re done
		defer r.Body.Close()

		// 1. Decode request
		var req struct {
			Filename string `json:"filename"`
			FileSize uint   `json:"fileSize"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		// 2. Create session
		response, err := p_uploadSessionService.InitUploadSessionService(fUploadData, 0, 0, req.Filename, req.FileSize, pglobalConfigs.CHUNKSIZE)
		if err != nil {
			http.Error(w, "Could not initiate upload", http.StatusInternalServerError)
			return
		}

		// 4. Send JSON
		w.Header().Set("Content-Type", "application/json")
		// (Optional) w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}

	}
}

/*
!!! Depricated
func UplodaChunk(w http.ResponseWriter, r *http.Request) {

	// important
	defer r.Body.Close()

	// get the upload ID
	vars := mux.Vars(r)
	uploadID := vars["uploadID"]

	if uploadID == "" {
		http.Error(w, "invalid Upload ID ", http.StatusBadRequest)
		return
	}

	// get the chunk id from the querry string
	idxStr := r.URL.Query().Get("index")
	idx, err := strconv.Atoi(idxStr)

	if err != nil {
		http.Error(w, "invalid chunk index ", http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "server could not read data ", http.StatusInternalServerError)
		return
	}

	err = puploadSessionService.WriteChunkAt(uploadID, idx, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 = success, no body

}

// / need to repurpose this function
func CompleteUploadSession(w http.ResponseWriter, r *http.Request) {

	// need to check if the request is from the file upload micro-service
	defer r.Body.Close()

	// need the uploadID
	vars := mux.Vars(r)
	uploadID := vars["uploadID"]

	if uploadID == "" {
		http.Error(w, "invalid upload ID ", http.StatusBadRequest)
	}
	// right now im going to set a single status code for any errors
	// but i should distinguish between server errors- 500 and unprocessableEntity- 422

	if err := puploadSessionService.UploadSessionFinalConfirmation(uploadID); err != nil {
		http.Error(w, "unprocessable entity ", http.StatusUnprocessableEntity)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&map[string]string{
		"message": "UploadComplete",
	})

}
*/
