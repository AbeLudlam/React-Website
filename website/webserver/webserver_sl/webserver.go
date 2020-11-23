package webserver_sl

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

func Start(address string, requestFileMapPath string) {
	fmt.Println("Starting server")

	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Listenting on: " + ln.Addr().String())

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err.Error())
		}

		go func(conn net.Conn) {
			fmt.Println("Connecting to: " + conn.RemoteAddr().String())
			chunk := make([]byte, 1024)
			requestHandler := &RequestHandler{make(map[string]string)}
			requestHandler.LoadRequestFileMap(requestFileMapPath)
			fmt.Println("Reading request file map path")

			for {
				input := make([]byte, 0)

				for {
					fmt.Println("reading...")
					conn.SetReadDeadline(time.Now().Add(time.Millisecond * 3000))
					n, err := conn.Read(chunk)
					fmt.Println(string(n))

					if err, ok := err.(net.Error); ok && err.Timeout() {
						fmt.Println("Timeout reached, closing conn")
						conn.Close()
						return
					}

					if n > 0 {
						input = append(input, chunk...)
					}

					if err == io.EOF {
						break
					}

					if err != nil {
						break
					}
					break
				}

				requestHeader := string(input)
				fmt.Println(string(requestHeader))

				requestHeaderLines := strings.Split(requestHeader, "\n")
				request := strings.Split(requestHeaderLines[0], " ")[1]
				requestFileInfo := requestHandler.GetRequestFile(request)

				fileStats, errStats := requestFileInfo.file.Stat()
				if errStats != nil {
					panic(errStats.Error())
				}

				file := requestFileInfo.file
				header := []byte(generateHTMLHeader(requestFileInfo.contentType, fileStats.Size()))
				wchunk := make([]byte, 1024)

				conn.Write(header)
				for {
					n, err := file.Read(wchunk)

					if n > 0 {
						conn.Write(wchunk[0:n])
					}

					if err == io.EOF {
						break
					}

					if err != nil {
						break
					}
				}
				file.Close()
			}
		}(conn)

	}
}

func generateHTMLHeader(contentType string, contentLength int64) string {
	header := "\r\nHTTP/1.1 200 OK\r\nContent-Type: " + contentType + "\r\n"
	if contentLength > 0 {
		header = header + "Content-Length: " + strconv.FormatInt(contentLength, 10) + ";\r\n"
	}
	header = header + "\r\n"
	return header
}
