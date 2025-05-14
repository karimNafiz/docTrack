package folder

import (
	db "docTrack/config"
	folder_errors "docTrack/errors/folder"
	user_errors "docTrack/errors/user"
	folder_model "docTrack/models/folders"
	user_service "docTrack/services/user"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gorm.io/gorm"
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

	fmt.Println("checked if the file name existis or not ")

	// TODO: remove this later
	// need to check if users exists or not
	// this is temporary code, remove this later
	_, err := user_service.FindUserByID(ownerID)
	if err != nil {
		return nil, user_errors.GetErrInvalidUserID(ownerID)
	}
	fmt.Println("user validated ")

	// need to get the parent Folder record frm the database
	// if it doesnt exist then we return a error
	var parentFolderPtr *folder_model.Folder
	parentFolderPtr, err = GetFolderByID(parentID)

	if err != nil {
		return nil, err
	}
	fmt.Println(" got the parent folder ////////////////////////////////////////////////// ")

	// I can now create the folder
	folderNew := folder_model.Folder{
		Name:                   folderName,
		Slug:                   getSlugFrmFolderName(folderName),
		OwnerID:                ownerID,
		ParentID:               parentID,
		Depth:                  parentFolderPtr.Depth + 1,
		ParentMaterializedPath: getFolderMaterializedPath(parentFolderPtr),
	}
	err = db.DB.Create(&folderNew).Error
	fmt.Println("error in creating the folder ", err)
	return &folderNew, err

}

func CopyFolder(folderID uint, ownerID uint, dstFolderID uint) {

	// get the original folder struct, from the database
	// so imagine we have the folder structure
	// x
	// --y
	//----z
	// if we copy y
	// we need to get the struct representing y
	// from the database
	originalFolderStructPtr, err := GetFolderByID(folderID)
	if err != nil {
		fmt.Println("the folder ID of the folder to be copied is invalid ")
		return // something sensible
	}

	// using the originalFolderStructPtr we create a duplicate Folder
	// but we pass the parentID the destination ID
	// so if we want to copy y to x'
	// we need to pass the ID of x'
	// so the dstFolderID is the folderID of x'
	duplicateFolderStructPtr, err := CreateFolder(originalFolderStructPtr.Name, ownerID, dstFolderID)

	if err != nil {
		fmt.Println(" error in top level folder duplicate")
		return // something sensible
	}
	// get the path to the Original Folder in the server
	// because we need to duplicate it's content
	originalFolderPath := getFolderMaterializedPath(originalFolderStructPtr)

	// in the copyFolderRecursive function we need to pass the path to the original folder in the server
	// then as parent we need to pass the ptr to duplicateFolderStruct
	copyFolderRecursive(originalFolderPath, duplicateFolderStructPtr, ownerID)

}

func copyFolderRecursive(originalFolderPath string, parentFolderPtr *folder_model.Folder, ownerID uint) {
	// read the entries from the originalFolderPath
	entries, err := os.ReadDir(originalFolderPath)
	if err != nil {
		return // something sensible
	}
	// loop through all the entries
	for _, entry := range entries {
		// if its not a dir
		// we don't have to do anything recursive
		// we just duplicate the File
		// set its parent to the parentFolderPtr passed into this function
		if !entry.IsDir() {
			// start copying the files
			// create file entries on the data base
		}

		// if its a folder
		// get the folder name
		// we need to create a duplicate from it
		childFolderName := entry.Name()
		// so we create a duplicate in the database
		// get a pointer to the struct
		childFolderPtr, err := CreateFolder(childFolderName, parentFolderPtr.ID, ownerID)
		if err != nil {
			return // something sensible
		}
		// continue recursively calling this function
		copyFolderRecursive(filepath.Join(originalFolderPath, entry.Name()), childFolderPtr, ownerID)
	}
}

func copyFile() {}

// need a function for updating the folder
// right now the only thing I can think of updating is the name

// need a function to Get Path

func GetFolderByID(folderID uint) (*folder_model.Folder, error) {
	ptr := new(folder_model.Folder)

	return ptr, db.DB.Where("id = ?", folderID).First(ptr).Error
}

// func getParentFolderMaterializedPath(parentFolderPtr *folder_model.Folder) string {
// 	return filepath.Join(parentFolderPtr.ParentMaterializedPath, parentFolderPtr.Slug)
// }

func getFolderMaterializedPath(folder *folder_model.Folder) string {
	return filepath.Join(folder.ParentMaterializedPath, folder.Slug)
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

func checkIsFileNameExist(name string, ownerID, parentID uint) bool {
	var f folder_model.Folder
	err := db.DB.
		Where("name = ? AND owner_id = ? AND parent_id = ?", name, ownerID, parentID).
		Take(&f).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	// if err is nil, record exists; if err is some other DB‐error, you may want to panic or log
	return err == nil
}
