package pdf

import (
	db "docTrack/config"
	pdf_model "docTrack/models/pdfs"
	user_service "docTrack/services/user"
	"fmt"
	"os"
	"path/filepath"
)

func CreatePDF(userID uint, filename string, filesaveLoation string, sizeBytes int64, pageNo uint) (*pdf_model.PDF, error) {

	_, err := user_service.FindUserByID(userID) // right now i dont need the user

	// if the error is not nil that means the user does not exists
	if err != nil {
		fmt.Println(fmt.Errorf("no user with this userID exists %w ", err))
		return nil, err
	}
	// the path variable is temporary for testing
	_, err = os.Stat(filepath.Join(filesaveLoation, filename))
	fmt.Println("filesaveLocation " + filesaveLoation)
	fmt.Println("filename " + filename)
	if err != nil {
		fmt.Println(fmt.Errorf(" no file of the filename %s exits %w ", filename, err))
		return nil, err
	}

	pdf := pdf_model.PDF{
		UserID:           userID,
		OriginalFilename: filename,
		FileSaveLocation: filesaveLoation,
		SizeBytes:        sizeBytes,
		PageNo:           pageNo,
	}

	err = db.DB.Create(&pdf).Error

	return &pdf, err

}
