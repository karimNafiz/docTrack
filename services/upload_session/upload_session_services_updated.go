package upload_session

import (
	pfileUploadServie "docTrack/file_upload_service"
)

func GetInitUploadSessionUpdated(fData *pfileUploadServie.FileUploadServiceInfo) func(uint, uint, string, uint, uint) error {
	return func(userID uint, parentID uint, filename string, filesize uint, chunkSize uint) error {
		// no need to check for userID if it exists or not it will be handled by middleware
		// no need to check for parentID if it exists or not it will be handler by middleware
		// assume all the data is properly formated
		// they will be taken care by the handler
		// need total chunks
		totalChunks := uint((filesize + (chunkSize - 1)) / chunkSize)

	}

}
