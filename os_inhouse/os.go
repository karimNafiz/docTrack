package os_inhouse

import (
	os_errors "docTrack/errors/os"
	logger "docTrack/logger"
	"os"
	"path/filepath"
)

// need a function to create a folder
func CreateFolder(path string, filename string) error {

	// target path
	targetPath := filepath.Join(path, filename)

	// First I need to check if the folder already exists or not
	flag, err := isFolderExists(filepath.Join(path, filename))

	// This error could signal permission denied
	// not sure what to do at this point
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	// if the flag is true
	// that means the Folder exists
	if flag {
		err = os_errors.GetErrFolderAlreadyExits(filename, path)
		logger.ErrorLogger.Println(err)
		return err
	}

	// if the folder doesnt exists already we need to create one
	err = os.Mkdir(targetPath, 0755)
	if err != nil {
		logger.ErrorLogger.Println(err)

	}

	logger.DebugLogger.Printf("folder %s created at path %s ", filename, path)
	return nil

}

func isFolderExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		// path exists; now check it’s a directory
		return info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		// path doesn’t exist
		return false, nil
	}
	// some other error (e.g. permission denied)
	return false, err
}
