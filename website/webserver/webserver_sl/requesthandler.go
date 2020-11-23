package webserver_sl

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type RequestHandler struct {
	requestMap map[string]string
}

func (rhandler *RequestHandler) LoadRequestFileMap(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err.Error())
	}

	fileString := ""
	bytes := make([]byte, 1024)
	for {
		n, readerr := file.Read(bytes)

		if n > 0 {
			fileString = fileString + string(bytes[:n])
		}

		if readerr == io.EOF {
			break
		}
		if readerr != nil {
			panic("Error reading request map")
		}
	}
	elements := strings.Split(fileString, "\"")

	i := 0
	for {
		requestIndex := (i * 4) + 1
		resultpathIndex := (i * 4) + 3
		if resultpathIndex >= len(elements) {
			break
		}
		rhandler.requestMap[elements[requestIndex]] = elements[resultpathIndex]
		i = i + 1
	}
}

func (rhandler *RequestHandler) GetRequestFile(request string) *ResponseInfo {
	filepath := rhandler.requestMap[request]

	extension := strings.LastIndex(filepath, ".")
	contentType := getContentTypeFromExtension(filepath[extension+1:])

	fmt.Println("Opening file: " + filepath)
	file, err := os.Open(filepath)
	if err != nil {
		panic(err.Error())
	}
	responseInfo := &ResponseInfo{file, contentType}
	return responseInfo
}

func getContentTypeFromExtension(extension string) string {
	switch extension {
	case "js":
		return "application/javascript; charset=UTF-8"
	case "html":
		return "text/html; charset=UTF-8"
	case "ico":
		return "image/x-icon"
	default:
		fmt.Println("Could not find extension" + extension)
		return ""
	}
}
