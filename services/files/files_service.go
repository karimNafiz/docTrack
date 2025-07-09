package files_service

// need a function to upload a file
/*
client wants to upload file
client does a post request to request for a upload_session more specifically an upload_id
server sends the port, domain of the file upload_service along with an upload_id
the server sends an upload_session_request to the file_upload_service
before everything the main_service must register with the file upload service

the client sends data to the file upload service
the file upload service receives the chunks
uploads everything sends a confirmation back to the client
then must notify the main service that the upload_session is complete


func UploadFile()

*/
