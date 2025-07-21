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
	"sync"

	utils "docTrack/utils"

	"github.com/google/uuid"

	logger "docTrack/logger"
	pdf_service "docTrack/services/pdf"

	folder_errors "docTrack/errors/folder"
	folder_service "docTrack/services/folder"

	os_inhouse "docTrack/os_inhouse"
	os_inhouse_chunkjob "docTrack/os_inhouse/chunk_job"
)

const tempUploadDir = "temp_uploads"
const finalPDFDir = "pdfs"
const defaultChunkSize = 5 << 20 // this is basically short for 5 * (1 << 20) , (1 << 20) is basically 2^20 as we are shifting 20 bits
const uploadFileExt = ".part"

// InitUploadSession initializes a new upload session for a given file.
// It creates a database entry for tracking the upload and also creates
// a temporary directory on disk where file chunks will be stored.
//
// It runs the database creation and file system operations in parallel,
// waits for both to finish, and performs rollback on partial failure.
// this design is depricated im moving everything related to uploading of files to different micro service
// i need to use that micro service
// depricatedFunction
func InitUploadSession(filename string, userID uint, parentID uint, fileSize int64) (*upload_session_model.UploadSession, error) {
	logger.InfoLogger.Printf("started upload session for file %s ", filename)

	// Calculate the number of chunks the file will be divided into.
	// This is equivalent to math.Ceil(fileSize / chunkSize), but done using integers.
	totalChunks := int((fileSize + int64(defaultChunkSize-1)) / int64(defaultChunkSize))

	// Step 1: Validate that the parent folder exists in the database.
	if _, err := folder_service.GetFolderByID(parentID); err != nil {
		logger.ErrorLogger.Println(folder_errors.ErrInvalidFolderID)
		return nil, folder_errors.ErrInvalidFolderID
	}

	// Step 2: Construct the upload session model.
	uploadSession := upload_session_model.UploadSession{
		ID:          uuid.NewString(), // Unique session ID
		UserID:      userID,           // Uploader's user ID
		ParentID:    parentID,         // Parent folder ID
		Filename:    filename,         // Name of the uploaded file
		FileSize:    fileSize,         // Total file size in bytes
		ChunkSize:   defaultChunkSize, // Standard chunk size (constant)
		TotalChunks: totalChunks,      // Number of chunks calculated
	}

	// Declare variables to hold the potential errors from each parallel task.
	var (
		wg    sync.WaitGroup // Used to wait for both goroutines
		dbErr error          // Error from database insertion
		fsErr error          // Error from filesystem folder creation
	)

	// Add two tasks to the WaitGroup
	wg.Add(2)

	// Step 3: Spawn a goroutine to create the upload session in the database.
	go func() {
		defer wg.Done() // Mark this task as done when the goroutine finishes
		dbErr = createUploadSessionRecord(&uploadSession)
	}()

	// Step 4: Spawn a goroutine to create the temporary upload directory on disk.
	go func() {
		defer wg.Done() // Mark this task as done when the goroutine finishes
		fsErr = os_inhouse.CreateFolder(tempUploadDir, uploadSession.ID)
	}()

	// Step 5: Wait for both goroutines to complete
	wg.Wait()

	// Step 6: Handle any errors that occurred during the goroutines
	if dbErr != nil || fsErr != nil {
		// If the database creation succeeded but the folder creation failed
		if dbErr == nil {
			// Roll back the database entry to maintain consistency
			_ = db.DB.Delete(&upload_session_model.UploadSession{}, "id = ?", uploadSession.ID).Error
		}

		// If the folder creation succeeded but the database insert failed
		if fsErr == nil {
			// Remove the temporary folder to avoid orphaned data
			_ = os.RemoveAll(filepath.Join(tempUploadDir, uploadSession.ID))
		}

		// Return the actual error that occurred (prefer DB error if both failed)
		if dbErr != nil {
			logger.ErrorLogger.Println(dbErr)
			return nil, dbErr
		}
		logger.ErrorLogger.Println(fsErr)
		return nil, fsErr
	}

	// Step 7: Return the initialized upload session (both tasks succeeded)
	return &uploadSession, nil
}

// needs to take the port and address of the file uploaded microservice
// for the microservice
// ill use it as a package
// you can download this package
// im not sure about the design right now
// ill jus use the micro service architecture
// need to store every in a config file
// fusInfo - stores meta data for the file upload service
func InitUploadSession_new(fusInfo map[string]string, userID uint, parentID uint, filename string, fileSize int64, chunkSize int) (map[string]string, error) {
	// check if the parenID is valid or not
	// this is prototyping, no need to check if the parent exists or not
	// i need to generate an uploadID
	uploadID := utils.GenerateUploadID()
	// need to send the filename for the filename need to produce a slug not doing it right now
	parentPtr, err := folder_service.GetFolderByID(parentID)
	if err != nil {
		logger.ErrorLogger.Println("parent folder not found, parentID : ", parentID) // temporary code
		return nil, err
	}
	// need the final path
	// TODO: -
	// consider making this function part of the folder struct
	finalPath := folder_service.GetFolderMaterializedPath(parentPtr)
	// need the chunk size
	// need total chunks
	totalChunks := int((fileSize + int64(chunkSize-1)) / int64(chunkSize))
	// need the service

}

func initUploadSession(domain string, port string)

//TODO need to delete the upload session if there are any errors
// need to add a column to the

func WriteChunkAt(uploadID string, chunkNo int, data []byte) error {
	// first i need to check if the uploadID exists or not
	var uploadSession upload_session_model.UploadSession

	err := db.DB.Where("id = ?", uploadID).First(&uploadSession).Error
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	// redundant check
	// if chunkNo >= uploadSession.TotalChunks || chunkNo < 0 {
	// 	err := errors.New(" no of chunks have exceeded the total no of chunks ")
	// 	logger.ErrorLogger.Println(err)
	// 	return err
	// }

	// check if the data recieved is the expected size
	if len(data) != uploadSession.ChunkSize {
		err := errors.New(" the chunk size exceeds the pre-defined chunk size ")
		logger.ErrorLogger.Println(err)
		return err
	}

	// if everything is fine we try to add a chunk job to our go-routine pool
	err = os_inhouse_chunkjob.AddChunkJob(os_inhouse_chunkjob.CreateChunkJob(uploadID, uint(chunkNo), tempUploadDir, data))

	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	// if there is no error
	// i should put the upload session to PROGRESS if its it not

}

// redundant function

// func writeChunkAt(path string, data []byte, _ int64) error {
// 	// Create or truncate the part file so it starts empty
// 	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
// 	// Simply write the chunk; its file size will == len(data)
// 	_, err = f.Write(data)
// 	return err
// }

func UploadSessionFinalConfirmation(uploadID string) error {

	tempDownloadPath := filepath.Join(tempUploadDir, uploadID)
	defer func() {
		clearDirectory(tempDownloadPath)
		DelUploadSessionFrmUploadID(uploadID)
	}()
	err := uploadSessionFinalConfirmation(uploadID, tempDownloadPath)

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
	fileNumberSet := make(map[int]string)

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
			fileNumberSet[fileNumber] = fileName
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

	return mergeChunksIntoPDF(&fileNumberSet, &uploadSession, tempDownloadPath, finalPDFDir)

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

func mergeChunksIntoPDF(chunkDict *map[int]string, uploadSession *upload_session_model.UploadSession, tempDir string, finalPDFDir string) error {
	// Grab the map and length
	fmt.Println("starting to merge part files into pdf ")
	dict := *chunkDict
	total := len(dict)
	if total == 0 {
		fmt.Println("chunkDict dictionary is empty")
		return errors.New("chunkDict dictionary is empty")
	}

	// 2) Create (or truncate) the final PDF and ensure it’s closed
	finalPath := filepath.Join(finalPDFDir, uploadSession.Filename)
	out, err := os.Create(finalPath)
	if err != nil {
		fmt.Println(fmt.Errorf("could not create final PDF %q: %w", finalPath, err))
		return fmt.Errorf("could not create final PDF %q: %w", finalPath, err)
	}
	defer out.Close()

	// 3) Iterate from chunk 0 up to TotalChunks-1
	for i := 0; i < uploadSession.TotalChunks; i++ {
		// ensure we have a filename for this index
		name, ok := dict[i]
		if !ok {
			fmt.Println(fmt.Errorf("missing chunk %d for upload %s", i, uploadSession.ID))
			return fmt.Errorf("missing chunk %d for upload %s", i, uploadSession.ID)
		}

		// 4) Open the part file and make sure it’s closed promptly
		partPath := filepath.Join(tempDir, name)
		in, err := os.Open(partPath)
		if err != nil {
			fmt.Println(fmt.Errorf("could not open chunk %d at %q: %w", i, partPath, err))
			return fmt.Errorf("could not open chunk %d at %q: %w", i, partPath, err)
		}

		// 5) Copy its bytes into the final PDF using your Copy()
		if _, err := utils.Copy(in, out); err != nil {
			in.Close()
			fmt.Println(fmt.Errorf("failed to copy chunk %d: %w", i, err))
			return fmt.Errorf("failed to copy chunk %d: %w", i, err)
		}
		in.Close()
	}

	// 6) All done
	// write now im jus using a fixed, hard coded, temporary, userID, ill link it later on
	_, err = pdf_service.CreatePDF(1, uploadSession.Filename, finalPDFDir, uploadSession.FileSize, 0)
	return err
}

func createUploadSessionRecord(session *upload_session_model.UploadSession) error {
	if err := db.DB.Create(session).Error; err != nil {
		logger.ErrorLogger.Printf("failed to create upload session record: %v", err)
		return err
	}
	return nil
}
