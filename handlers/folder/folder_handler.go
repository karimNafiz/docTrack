package folder

import (
	logger "docTrack/logger"
	folder_service "docTrack/services/folder"
	"encoding/json"
	"fmt"
	"net/http"
)

func CreateFolderHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// what do i need
	// folderName
	// parentID
	// ownerID
	var requestBody struct {
		FolderName string `json:"folderName"`
		ParentID   uint   `json:"parentID"`
	}

	fmt.Println("parsed the JSON BODY ")
	// i did not do any checks doesn't matter

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
		return
	}
	// im hard coding the owner ID
	// ill change it later
	_, err := folder_service.CreateFolder(requestBody.FolderName, 01, requestBody.ParentID)
	fmt.Println("Created the folder ")
	if err != nil {
		// need to handler error better
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Success Folder Created",
	})

}

func CopyFolderHandler(w http.ResponseWriter, r *http.Request) {
	// need to close the Body of the request
	// It is an encapsulation for a tcp connection
	// if we do not release the connection
	// go will create another tcp connection with the client
	// this creates a lot of overhead (three way handshake and very wasteful on system resources)
	defer r.Body.Close()

	// we need to parse the json
	// so we create a struct to mimic the json Body
	var requestBody struct {
		FolderID uint `json:"folderID"`
		ParentID uint `json:"parentID"`
	}

	// We are trying to parse the JSON Body
	// if we have an error
	// we use http.Error to send a message stating Invalid JSON Body
	// with the status code http.StatusBadRequest
	// always remember, the header is sent before message
	// the header contains the status code
	// when we do http.Error we are sending an header
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		logger.ErrorLogger.Println(err)
		http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
		return
	}

	// after parsing the JSON body we need to call the folder_service.CopyFolder method that will do all the heavy lifting
	err := folder_service.CopyFolder(requestBody.FolderID, 1, requestBody.ParentID)

	if err != nil {
		http.Error(w, "Internal Server Error ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Success Folder Copied",
	})

}

// func CopyFolderHandler(w http.ResponseWriter , r *http.Request){
// 	defer r.Body.Close()
// 	var requestBody struct{
// 		FolderID uint `json:"folderID"`
// 		DstFolderID uint `json:"dstFolderID"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
// 		http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
// 		return
// 	}

// 	folder_service.CopyFolder(requestBody.FolderID , 01 , requestBody.DstFolderID)

// }
