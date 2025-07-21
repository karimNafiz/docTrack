package folder

import (
	db "docTrack/config"
	folder_errors "docTrack/errors/folder"
	user_errors "docTrack/errors/user"
	logger "docTrack/logger"
	folder_model "docTrack/models/folders"
	os_inhouse "docTrack/os_inhouse"
	user_service "docTrack/services/user"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

// collapse any sequence of invalid chars into one hyphen
var invalidRe = regexp.MustCompile(`[\x00-\x1F/\\:*?"<>|]+`)

func CreateFolder(flderName string, ownerID uint, parentID uint) (*folder_model.Folder, error) {

	// making sure the foler name is valid
	if flderName == "" {
		return nil, folder_errors.ErrInvalidFolderName
	}
	// if the folder name already exists
	// we just add a number after the folder name
	// imagine we have a folder called x
	// if we create another folder named x, this will make sure the folder name is x_1
	// imagine we also have x_1
	// then we will have x_2
	var count uint8 = 0
	var folderName string = flderName
	for checkIsFileNameExist(folderName, ownerID, parentID) {
		count++
		folderName = fmt.Sprintf("%s_%d", flderName, count)
	}

	// TODO the use checking will be done long before
	// the code even reaches this point
	_, err := user_service.FindUserByID(ownerID)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, user_errors.GetErrInvalidUserID(ownerID)
	}

	// we also need to check if the parent folder also exists in the database
	var parentFolderPtr *folder_model.Folder
	parentFolderPtr, err = GetFolderByID(parentID)

	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}
	slug := getSlugFrmFolderName(folderName)
	parentFolderPath := GetFolderMaterializedPath(parentFolderPtr)
	err = os_inhouse.CreateFolder(parentFolderPath, slug)

	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	// I can now create the folder
	folderNew := folder_model.Folder{
		Name:                   folderName,
		Slug:                   slug,
		OwnerID:                ownerID,
		ParentID:               parentID,
		Depth:                  parentFolderPtr.Depth + 1,
		ParentMaterializedPath: parentFolderPath,
	}
	err = db.DB.Create(&folderNew).Error

	if err != nil {
		logger.ErrorLogger.Println(err)
	} else {
		logger.DebugLogger.Printf("created record id %d in table %s ", folderNew.ID, "folders")
	}
	return &folderNew, err

}

func CopyFolder(folderID uint, ownerID uint, dstFolderID uint) error {

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
		// TODO make an error about invalid Folder ID
		logger.ErrorLogger.Println(err)
		return err
	}

	// using the originalFolderStructPtr we create a duplicate Folder
	// but we pass the parentID the destination ID
	// so if we want to copy y to x'
	// we need to pass the ID of x'
	// so the dstFolderID is the folderID of x'
	duplicateFolderStructPtr, err := CreateFolder(originalFolderStructPtr.Name, ownerID, dstFolderID)

	if err != nil {
		logger.ErrorLogger.Println(err)
		return err // something sensible
	}
	// get the path to the Original Folder in the server
	// because we need to duplicate it's content
	originalFolderPath := GetFolderMaterializedPath(originalFolderStructPtr)

	// in the copyFolderRecursive function we need to pass the path to the original folder in the server
	// then as parent we need to pass the ptr to duplicateFolderStruct
	return copyFolderRecursive(originalFolderPath, duplicateFolderStructPtr, ownerID)

}

func copyFolderRecursive(originalFolderPath string, parentFolderPtr *folder_model.Folder, ownerID uint) error {
	// read the entries from the originalFolderPath
	entries, err := os.ReadDir(originalFolderPath)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err // something sensible
	}
	// loop through all the entries
	for _, entry := range entries {
		// if its not a dir
		// we don't have to do anything recursive
		// we just duplicate the File
		// set its parent to the parentFolderPtr passed into this function
		if !entry.IsDir() {
			srcFilename := entry.Name()
			srcFilePath := filepath.Join(originalFolderPath, srcFilename)
			dstFolderPath := GetFolderMaterializedPath(parentFolderPtr)
			copyFile(srcFilePath, dstFolderPath, srcFilename)
		}

		// if its a folder
		// get the folder name
		// we need to create a duplicate from it
		childFolderName := entry.Name()
		// so we create a duplicate in the database
		// get a pointer to the struct
		childFolderPtr, err := CreateFolder(childFolderName, ownerID, parentFolderPtr.ID)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return err // something sensible
		}
		// continue recursively calling this function
		err = copyFolderRecursive(filepath.Join(originalFolderPath, entry.Name()), childFolderPtr, ownerID)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return err
		}
	}
	return nil
}

// how should the copyFile function work
// get the path to the original file
// get the path to the destination
// open connections to both the paths
// try to check if already a file with the same name exists or not
// if so do the technique I have done before
// just copy
func copyFile(originalFilePath string, destinationPath string, filename string) error {
	// try to open a connection to the originalFilePath
	srcFile, err := os.Open(originalFilePath)

	// if there is an error
	// we need to return the error
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	// at this point we have opened a connection to our file
	// make sure we have this
	// or else memory will leaked
	defer srcFile.Close()

	dstFilePath := filepath.Join(destinationPath, filename)
	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	// at this point we have an open connection to the destination file path
	defer dstFile.Close()

	// copy everything
	_, err = io.Copy(dstFile, srcFile)

	// TODO: make sure dat all the bytes are copied
	// open up os.Stats for the srcFile
	if err != nil {
		logger.ErrorLogger.Println(err)
		// TODO need to delete the destination file created
		os.Remove(dstFilePath)
		return err
	}
	return nil

}

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

func GetFolderMaterializedPath(folder *folder_model.Folder) string {
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
