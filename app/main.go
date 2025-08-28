package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
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
	tempDirectory := "temp"
	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)

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
		Body    string
	}
	request := HTTPRequest{
		Url:     lineparts[1],
		Headers: header,
		Method:  method,
		Body:    body_content,
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
		body := "You have hit /" + "\n"
		response := "HTTP/1.1 200 OK\r\n" +
			"Content-Type: text/plain\r\n" +
			"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
			"\r\n" +
			body

		conn.Write([]byte(response))
		fmt.Println("These is /")

	} else if strings.HasPrefix(lineparts[1], "/echo") {
		newParts := strings.Split(lineparts[1], "/")
		if len(newParts) > 3 {
			send404(conn)
			return
		}
		text := newParts[2]
		textlength := len(newParts[2])

		conn.Write([]byte(
			fmt.Sprintf(
				"HTTP/1.1 200 OK\r\n"+
					"Content-Type: text/plain\r\n"+
					"Content-Length: %d\r\n\r\n%s",
				textlength,
				text,
			),
		))

	} else if strings.HasPrefix(lineparts[1], "/user-agent") {
		content := header["User-Agent"]
		contentlen := len(content)
		conn.Write([]byte(
			fmt.Sprintf(
				"HTTP/1.1 200 OK\r\n"+
					"Content-Type: text/plain\r\n"+
					"Content-Length: %d\r\n\r\n%s",
				contentlen,
				content,
			),
		))

	} else if request.Method == "GET" && strings.HasPrefix(request.Url, "/files") {
		fileparts := strings.Split(request.Url, "/")

		filename := fileparts[2]
		filepath := fmt.Sprintf("%s/%s", tempDirectory, filename)

		if len(fileparts) > 3 {
			send404(conn)
			return
		}

		if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
			send404(conn)
			return
		}
		content, _ := os.ReadFile(filepath)
		contentLength := utf8.RuneCountInString(string(content))
		conn.Write([]byte(
			fmt.Sprintf(
				"HTTP/1.1 200 OK\r\n"+
					"Content-Type: application/octet-stream\r\n"+
					"Content-Length: %d\r\n\r\n%s", contentLength, content)))

	} else if request.Method == "POST" && strings.HasPrefix(request.Url, "/files/") {
		fileparts := strings.Split(request.Url, "/")
		filname := fileparts[2]
		filepath := fmt.Sprintf("%s/%s", tempDirectory, filname)

		if len(fileparts) > 3 {
			send404(conn)
			return
		}

		if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
			send404(conn)
			return
		}

		headerEnd := bytes.Index(buff, []byte("\r\n\r\n"))
		headerLength := headerEnd + 4
		if err := os.WriteFile(filepath, buff[headerLength:n], 0644); err == nil {
			send404(conn)
			return
		}
		conn.Write([]byte("http/1.1 201 Created \r\n\\r\n"))
	} else {
		send404(conn)
	}
}

func send404(conn net.Conn) error {
	body := "404 not found\n"
	response := "HTTP/1.1 404 Not Found\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
		"\r\n" +
		body

	_, err := conn.Write([]byte(response))
	fmt.Println("[ERROR]: request did not match any method or route")
	return err
}
