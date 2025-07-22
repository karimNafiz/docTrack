package upload_session

import (
	"bytes"
	"context"
	pfileUploadServie "docTrack/file_upload_service"
	pGlobalConfigs "docTrack/global_configs"
	pFolderService "docTrack/services/folder"
	pUtils "docTrack/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

func GetInitUploadSessionUpdated(fData *pfileUploadServie.FileUploadServiceInfo) func(uint, uint, string, uint, uint) (json.Marshaler, error) {
	return func(userID uint, parentID uint, filename string, filesize uint, chunkSize uint) (json.Marshaler, error) {
		// no need to check for userID if it exists or not it will be handled by middleware
		// no need to check for parentID if it exists or not it will be handler by middleware
		// assume all the data is properly formated
		// they will be taken care by the handler
		// need total chunks
		// TODO: remove the code for getting the parent record from the database, getting the parent will be handled by a middleware
		folderPtr, err := pFolderService.GetFolderByID(parentID)
		if err != nil {
			fmt.Println("GetFolderByID err:", err)
			return nil, err
		}
		totalChunks := uint((filesize + (chunkSize - 1)) / chunkSize)
		uploadID := pUtils.GenerateUploadID()
		finalPath := pFolderService.GetFolderMaterializedPath(folderPtr)
		reqBody := pUtils.FUploadServiceHttpRequestBody{
			UploadID:    uploadID,
			Filename:    filename,
			FinalPath:   finalPath,
			ChunkSize:   chunkSize,
			TotalChunks: totalChunks,
			ServiceID:   fData.ServiceID,
		}
		payload, err := reqBody.MarshalJSON()
		url := filepath.Join(pGlobalConfigs.GetFolderServiceRootUrlHttp(), pGlobalConfigs.FILEUPLOADSERVICEINITUPLOADSESSION)
		headers := map[string]string{
			"Content-Type": "application/json",
		}
		response, err := pUtils.SendHttpRequest(context.Background(), http.MethodPost, url, headers, bytes.NewBuffer(payload))
		if err != nil {
			log.Println("error with sending http request err -> ", err.Error())
		}
		// looks fancy nothing else
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(response.Body)
		log.Println("file upload service response status code ", response.StatusCode)
		respBytes, err := io.ReadAll(response.Body)
		respBody := pUtils.InitUploadSessionClientResponse{
			UploadID:    uploadID,
			ChunkSize:   chunkSize,
			TotalChunks: totalChunks,
		}
		err = json.Unmarshal(respBytes, &respBody)
		if err != nil {
			log.Println("error reading response body -> ", err.Error())
			return nil, err
		}
		log.Println("file upload server responded with status ", respBody.Status)
		// TODO: create global structs like these so I don't have to create them on the fly
		return respBody, nil

	}
}
