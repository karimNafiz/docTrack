package folder

import (
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
