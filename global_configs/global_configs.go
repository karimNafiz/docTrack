package global_configs

// a 2MB chunk size
const CHUNKSIZE = 2 * (1 << 20)
const FILEUPLOADSERVICESCHEME = "https"
const FILEUPLOADSERVICEPORT = ":8443"
const FILEUPLOADSERVICEDOMAIN = "localhost"
const FILEUPLOADSERVICEREGISTERENDPOINT = "/register"
const FILEUPLOADSERVICEINITUPLOADSESSION = "/upload_request"
const FILEUPLOADSERVICECALLBACKURL = "/upload_status"

const MAINSERVICEPORT = ":8080"
const MAINSERVICEDOMAIN = "localhost"
const MAINSERVICESCHEME = "http:"
