package upload_session

import (
	db "docTrack/config"
	upload_session_model "docTrack/models/upload_sessions"
	"errors"
	"math"
	"os"

	"github.com/google/uuid"
)

var path string = "need to set it later"

func InitUploadSession(filename string, fileSize int64, chunkSize int) (*upload_session_model.UploadSession, error) {
	var uploadSession upload_session_model.UploadSession

	uploadSession = upload_session_model.UploadSession{
		ID:           uuid.NewString(),
		Filename:     filename,
		File_size:    fileSize,
		Chunk_size:   chunkSize,
		Total_chunks: int(math.Ceil(float64(fileSize) / float64(chunkSize))),
	}

	if err := db.DB.Create(&uploadSession).Error; err != nil {
		return nil, err
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

	if chunkNo > uploadSession.Total_chunks {
		return errors.New(" no of chunks have exceeded the total no of chunks ")
	}

	if len(data) > uploadSession.Chunk_size {
		return errors.New(" the chunk size exceeds the pre-defined chunk size ")
	}

	err = writeChunkAt(path, data, int64(uploadSession.Chunk_size*chunkNo))
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
