package upload_session

import (
	p_file_upload_service "docTrack/file_upload_service"
	upload_session_service "docTrack/services/upload_session"
	"encoding/json"
	"net/http"
)

func GetInitUploadSessionHandler(fUploadData *p_file_upload_service.FileUploadServiceInfo) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Close the body when we’re done
		defer r.Body.Close()

		// 1. Decode request
		var req struct {
			Filename string `json:"filename"`
			FileSize int64  `json:"fileSize"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		// 2. Create session
		session, err := upload_session_service.InitUploadSession(req.Filename, 0, 0, req.FileSize)
		if err != nil {
			http.Error(w, "Could not initiate upload", http.StatusInternalServerError)
			return
		}

		// 3. Build response
		resp := struct {
			UploadID    string `json:"uploadID"`
			ChunkSize   int64  `json:"chunkSize"`
			TotalChunks int    `json:"totalChunks"`
		}{
			UploadID:    session.ID,
			ChunkSize:   int64(session.ChunkSize),
			TotalChunks: session.TotalChunks,
		}

		// 4. Send JSON
		w.Header().Set("Content-Type", "application/json")
		// (Optional) w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
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

	err = upload_session_service.WriteChunkAt(uploadID, idx, data)
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

	if err := upload_session_service.UploadSessionFinalConfirmation(uploadID); err != nil {
		http.Error(w, "unprocessable entity ", http.StatusUnprocessableEntity)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&map[string]string{
		"message": "UploadComplete",
	})

}
*/
