package global_configs

import "fmt"

// a 2MB chunk size
const CHUNKSIZE = 2 * (1 << 20)
const FILEUPLOADSERVICESCHEME = "https"
const FILEUPLOADSERVICEPORT = "8443"
const FILEUPLOADSERVICEDOMAIN = "localhost"
const FILEUPLOADSERVICEREGISTERENDPOINT = "/register"
const FILEUPLOADSERVICEINITUPLOADSESSION = "/upload/init"
const FILEUPLOADSERVICECALLBACKURL = "/upload_status"

const MAINSERVICEPORT = "8080"
const MAINSERVICEDOMAIN = "localhost"
const MAINSERVICESCHEME = "http"

func GetFolderServiceRootUrlHttp() string {
	return fmt.Sprintf("%s://%s:%s", "http", FILEUPLOADSERVICEDOMAIN, "9000")
}
func GetFolderServiceRootUrlHttps() string {
	return fmt.Sprintf("%s://%s:%s", "https", FILEUPLOADSERVICEDOMAIN, "8443")
}
