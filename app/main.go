package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"unicode/utf8"
)


var _ = net.Listen
var _ = os.Exit

const CRLF = "\r\n"

func main() {
	fmt.Println("logs from your program will appear here!")

	
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("error accepting connection: ", err.Error())
			os.Exit(1)

		}
		go Handlereq(conn) //handle multiple threads automatically (concurrently)
	}
}

// msg := []byte("HTTP/1.1 200 OK\r\n\r\n")
// conn.Write(msg)

// conn.Close()

// server will extract the URL path from an HTTP request, and respond with either a 200 or 404, depending on the path.
func Handlereq(conn net.Conn) {
	tempDirectory := "temp";
	buff := make([]byte, 1024)
	conn.Read(buff)

	parts := strings.Split(string(buff), CRLF)
    len_of_parts := len(parts)
	lineparts := strings.Split(parts[0], " ")
	method := lineparts[0]
	body_content := parts[len_of_parts-1]


	header := make(map[string]string)
	type HTTPRequest struct {
		Headers map[string]string
		Url     string
		Method  string
		Body	string
 	}
	request := HTTPRequest{
		Url:     lineparts[1],
		Headers: header,
		Method: method,
		Body: body_content,
	}
	for i := 1; i < len(parts); i++ {
		lines := parts[i]
		if lines == " " {
			break
		}
		headerParts := strings.SplitN(lines, ": ", 2)
		if len(headerParts) == 2 {
			key := headerParts[0]
			value := headerParts[1]
			header[key] = value
		}
	}

	if lineparts[1] == "/" {
		conn.Write([]byte("http/1.1 200 OK\r\n\r\n"))
		fmt.Println(("These is /"))

	} else if strings.HasPrefix(lineparts[1], "/echo") {
		newParts := strings.Split(lineparts[1], "/")
		if len(newParts) > 3 {
			conn.Write([]byte("http/1.1 404 Not Found\r\n\r\n"))
		}
		text := newParts[2]
		textlength := len(newParts[2])

		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", textlength, text)))

	} else if strings.HasPrefix(lineparts[1], "/user-agent") {
		content := header["User-Agent"]
		contentlen := len(content)
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", contentlen, content)))

	}else if strings.HasPrefix(request.Url,"/files"){
		fileparts := strings.Split(request.Url,"/")
		
		filename := fileparts[2]
		filepath := fmt.Sprintf("%s/%s", tempDirectory, filename)


		if len(fileparts) > 3 {
			conn.Write([]byte("http/1.1 404 Not Found\r\n\r\n"))
			return
		}

		if _,err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
			conn.Write([]byte("http/1.1 404 Not Found\r\n\r\n"))
			return
		}
		content,_ := os.ReadFile(filepath)
		contentLength := utf8.RuneCountInString(string(content))
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", contentLength, content)))
	}else if  request.Method == "POST"{
		

	} else {
		conn.Write([]byte("http/1.1 404 not found\r\n\r\n"))

	}
}
