package os

import (
	"fmt"
)

func GetErrFolderAlreadyExits(filename string, path string) error {

	return fmt.Errorf("folder %s at path %s already exists", filename, path)
}
