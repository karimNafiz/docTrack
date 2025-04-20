package upload_session

import (
	db "docTrack/config"
	upload_session_model "docTrack/models/upload_sessions"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

var tempUploadDir string = "temp_uploads"

const defaultChunkSize = 5 << 20 // this is basically short for 5 * (1 << 20) , (1 << 20) is basically 2^20 as we are shifting 20 bits

func InitUploadSession(filename string, fileSize int64) (*upload_session_model.UploadSession, error) {

	var totalChunkSize int = int((fileSize + int64(defaultChunkSize-1)) / int64(defaultChunkSize)) // this is more efficient version of math.ceil which is for float

	uploadSession := upload_session_model.UploadSession{
		ID:          uuid.NewString(),
		Filename:    filename,
		FileSize:    fileSize,
		ChunkSize:   defaultChunkSize,
		TotalChunks: totalChunkSize,
	}

	if err := db.DB.Create(&uploadSession).Error; err != nil {
		return nil, err
	}

	// need to make the temporary Directory
	if err := os.MkdirAll(filepath.Join(tempUploadDir, uploadSession.ID), 0755); err != nil {
		// could not create the temporary directory, need to delete the database entry
		db.DB.Delete(&upload_session_model.UploadSession{}, "id = ?", uploadSession.ID)
		return nil, errors.New("could not create tempory directory for upload")
	}

	return &uploadSession, nil

}

//TODO need to delete the upload session if there are any errors
// need to add a column to the

func WriteChunkAt(uploadID string, chunkNo int, data []byte) error {
	// first i need to check if the uploadID exists or not
	var uploadSession upload_session_model.UploadSession

	err := db.DB.Where("id = ?", uploadID).First(&uploadSession).Error
	if err != nil {
		return err
	}

	if chunkNo >= uploadSession.TotalChunks || chunkNo < 0 {
		return errors.New(" no of chunks have exceeded the total no of chunks ")
	}

	if len(data) > uploadSession.ChunkSize {
		return errors.New(" the chunk size exceeds the pre-defined chunk size ")
	}
	partPath := filepath.Join(filepath.Join(tempUploadDir, uploadSession.ID), fmt.Sprintf("%06d.part", chunkNo))

	err = writeChunkAt(partPath, data, int64(uploadSession.ChunkSize)*int64(chunkNo))
	return err

}

func writeChunkAt(path string, data []byte, offset int64) error {
	f, err := os.OpenFile(path,
		os.O_CREATE|os.O_RDWR,
		0644,
	)
	if err != nil {
		return err
	}
	defer f.Close() // very important

	if _, err := f.WriteAt(data, offset); err != nil {
		return err

	}
	return nil

}
