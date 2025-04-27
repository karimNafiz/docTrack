package upload_session

import (
	db "docTrack/config"
	upload_session_model "docTrack/models/upload_sessions"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	utils "docTrack/utils"

	"github.com/google/uuid"
)

var tempUploadDir string = "temp_uploads"

const defaultChunkSize = 5 << 20 // this is basically short for 5 * (1 << 20) , (1 << 20) is basically 2^20 as we are shifting 20 bits
const uploadFileExt = ".part"

func InitUploadSession(filename string, fileSize int64) (*upload_session_model.UploadSession, error) {

	fmt.Println("inside intit uploadSession method ")
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
	fmt.Println("writing chunk at path " + partPath + " \n from function WriteChunkAt")
	err = writeChunkAt(partPath, data, int64(uploadSession.ChunkSize)*int64(chunkNo))
	return err

}

func writeChunkAt(path string, data []byte, _ int64) error {
	// Create or truncate the part file so it starts empty
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	// Simply write the chunk; its file size will == len(data)
	_, err = f.Write(data)
	return err
}

func UploadSessionFinalConfirmation(uploadID string) error {
	tempDownloadPath := filepath.Join(tempUploadDir, uploadID)
	err := uploadSessionFinalConfirmation(uploadID, tempDownloadPath)
	if err != nil {
		DelUploadSessionFrmUploadID(uploadID)
		clearDirectory(tempDownloadPath)
		fmt.Println("error : " + err.Error())
	}
	return err
}

func uploadSessionFinalConfirmation(uploadID string, tempDownloadPath string) error {
	var uploadSession upload_session_model.UploadSession

	err := db.DB.Where("id = ?", uploadID).First(&uploadSession).Error
	if err != nil {

		return fmt.Errorf("the uploadID was not found, the upload file may have been corrupted and deleted %w ", err)
	}

	// after this i know the uploadSession record with this uploadID exists
	// need to find the path to the temporary downloaded parts

	// need to read the entries from this path
	entries, err := os.ReadDir(tempDownloadPath)
	if err != nil {

		// that specific directory doesnt exists
		return fmt.Errorf("the directory for the temporary uploaded chunks do not exist ")
	}

	totalChunks := 0
	var totalFileSize int64 = 0
	fileNumberSet := make(map[int]struct{})

	for _, e := range entries {
		// if its a directory then we skip
		// jus being safe
		if e.IsDir() {
			continue
		}

		fileName := e.Name()
		if strings.HasSuffix(fileName, uploadFileExt) {
			fileNumber, err := strconv.Atoi(strings.TrimSuffix(fileName, uploadFileExt))
			if err != nil {
				return fmt.Errorf(" file is corrupted \n actual error %w ", err)
			}
			fileNumberSet[fileNumber] = struct{}{}
			totalChunks++
			info, err := e.Info()
			if err != nil {
				return fmt.Errorf("internal server issue \n actual error %w ", err)
			}
			totalFileSize += info.Size()

		}
	}
	// after looping through all the entries and extracting meta data we will now check
	// if the data is corrupted or not
	// checking the number of chunks
	if totalChunks != uploadSession.TotalChunks {
		return fmt.Errorf("the no of chunks uploaded dont match, data was corrupted while transferring")
	}

	// need to check file size
	if totalFileSize != uploadSession.FileSize {
		fmt.Println("uploaded file size ", totalFileSize)
		fmt.Println("file size sent in header ", uploadSession.FileSize)
		return fmt.Errorf("the uploaded fileSize and actual fileSize do not match, data was corrupted while transferring")
	}

	// need to check if all the chunks are present
	for i := range uploadSession.TotalChunks {
		_, ok := fileNumberSet[i]
		if !ok {
			return fmt.Errorf("not all chunks are present ")
		}
	}

	return nil

}

func DelUploadSessionFrmUploadID(uploadID string) error {
	err := db.DB.Delete(&upload_session_model.UploadSession{}, "id = ?", uploadID).Error
	return err

}

func clearDirectory(directory string) error {
	entries, err := os.ReadDir(directory)
	if err != nil { // the directory is invalid
		return fmt.Errorf("the directory path is invalid  actual error \n %w", err)
	}
	for _, e := range entries {
		err := os.RemoveAll(filepath.Join(directory, e.Name()))
		if err != nil {
			return fmt.Errorf("there was an issue deleting the file %s \n actual error %w", e.Name(), err)
		}
	}
	return nil

}

// need to create pdf out of all the uploaded parts

// i have a map of chunk no, and the file name and the suffix
// use that to

func mergeChunksIntoPDF(chunkDict *map[int]string, upload_session *upload_session_model.UploadSession, path string) error {

	dictLen := len(*chunkDict)

	if dictLen <= 0 {
		return errors.New("chunkDict dictionary is empty ")
	}

	// loop through the entries
	// from 0 to len(chunkDict - 1)
	for i := range dictLen {
		fileName := (*chunkDict)[i]
		in, err := os.Open(filepath.Join(path, fileName))
		if err != nil {
			// need to decide what to do
			return err
		}
		// need to close the stream
		defer in.Close()

		utils.ReadBytes(in, "4KB")
	}

}
