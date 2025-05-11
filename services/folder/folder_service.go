package folder

import (
	db "docTrack/config"
	folder_errors "docTrack/errors/folder"
	user_errors "docTrack/errors/user"
	folder_model "docTrack/models/folders"
	user_service "docTrack/services/user"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// collapse any sequence of invalid chars into one hyphen
var invalidRe = regexp.MustCompile(`[\x00-\x1F/\\:*?"<>|]+`)

func CreateFolder(folderName string, ownerID uint, parentID uint) (*folder_model.Folder, error) {

	// need to check the length too

	if folderName == "" {
		return nil, folder_errors.ErrInvalidFolderName
	}
	var count uint8 = 0
	for checkIsFileNameExist(folderName, ownerID, parentID) {
		count++
		folderName = fmt.Sprintf("%s_%d", folderName, count)
	}

	// TODO: remove this later
	// need to check if users exists or not
	// this is temporary code, remove this later
	_, err := user_service.FindUserByID(ownerID)
	if err != nil {
		return nil, user_errors.GetErrInvalidUserID(ownerID)
	}

	// need to get the parent Folder record frm the database
	// if it doesnt exist then we return a error
	var parentFolderPtr *folder_model.Folder
	parentFolderPtr, err = GetFolderByID(parentID)

	if err != nil {
		return nil, err
	}

	// I can now create the folder
	folderNew := folder_model.Folder{
		Name:             folderName,
		Slug:             getSlugFrmFolderName(folderName),
		OwnerID:          ownerID,
		ParentID:         parentID,
		Depth:            parentFolderPtr.Depth + 1,
		MaterializedPath: getParentFolderPath(parentFolderPtr),
	}

	return &folderNew, db.DB.Create(&folderName).Error

}

// need a function for updating the folder
// right now the only thing I can think of updating is the name

// need a function to Get Path

func GetFolderByID(folderID uint) (*folder_model.Folder, error) {
	ptr := new(folder_model.Folder)

	return ptr, db.DB.Where("id = ?", folderID).First(ptr).Error
}

func getParentFolderPath(folder *folder_model.Folder) string {
	return filepath.Join(folder.MaterializedPath, folder.Slug)
}

func getSlugFrmFolderName(folderName string) string {
	// need to get rid of white spaces infront and at the end of the name
	// unicode normalization ?
	// lower case the entire string
	// check for invalid characters
	// collapse multiple hyphens
	// leading and trailing hyphens

	folderName = invalidRe.ReplaceAllString(folderName, "_")

	// get rid of trailing white spaces, hyphens
	folderName = strings.Trim(folderName, " ")
	folderName = strings.Trim(folderName, "_")
	folderName = strings.Trim(folderName, "-")

	// check out unicode normalization

	// lower case the entire string
	folderName = strings.ToLower(folderName)
	return folderName

}

func checkIsFileNameExist(fileName string, ownerID uint, parentID uint) bool {
	if db.DB.Model(&folder_model.Folder{}).Where("name = ? AND owner_id = ? AND parent_id = ?", fileName, ownerID, parentID).Error != nil {
		return false
	}
	return true

}
