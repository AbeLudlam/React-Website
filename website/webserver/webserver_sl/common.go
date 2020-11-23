package webserver_sl

import "os"

type ResponseInfo struct {
	file        *os.File
	contentType string
}
