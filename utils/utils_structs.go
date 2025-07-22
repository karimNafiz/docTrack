package utils

import "encoding/json"

type InitUploadSessionClientResponse struct {
	Status           string `json:"status"`
	FileUploadDomain string `json:"file_upload_domain"`
	FileUploadPort   string `json:"file_upload_port"`
	UploadID         string `json:"uploadID"`
	ChunkSize        uint   `json:"chunk_size"`
	TotalChunks      uint   `json:"total_chunks"`
}

func (u InitUploadSessionClientResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(u)
}

type FUploadServiceHttpRequestBody struct {
	UploadID    string `json:"uploadID"`
	Filename    string `json:"filename"`
	FinalPath   string `json:"final_path"`
	ChunkSize   uint   `json:"chunk_size"`
	TotalChunks uint   `json:"total_chunks"`
	ServiceID   string `json:"serviceID"`
}

func (f FUploadServiceHttpRequestBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(f)
}
