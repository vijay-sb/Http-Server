package main

import (
	"fmt"
	"net"
	"os"
  "strings"
)

// ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit
const CRLF = "\r\n"

func main() {
	// you can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("logs from your program will appear here!")

	// uncomment this block to pass the first stage
	//
	 l, err := net.Listen("tcp", "0.0.0.0:4221")
	 if err != nil {
	 	fmt.Println("failed to bind to port 4221")
	 	os.Exit(1)
	 }
	
  conn, err := l.Accept()
	 if err != nil {
	 	fmt.Println("error accepting connection: ", err.Error())
	 	os.Exit(1)
	 }
     
  msg := []byte("HTTP/1.1 200 OK\r\n\r\n")
   conn.Write(msg)
  
  conn.Close()
  


  //server will extract the URL path from an HTTP request, and respond with either a 200 or 404, depending on the path.
  buff := make([]byte,1024)
  conn.Read(buff)

  parts := strings.Split(string(buff),CRLF)

  lineparts := strings.Split(parts[0]," ")
  if lineparts[1] == "/"  {
    conn.Write([]byte("http/1.1 200 OK\r\n\r\n"))
    conn.Close()
  } else if strings.HasPrefix(lineparts[1],"/echo") {
    newParts := strings.Split(lineparts[1],"/")
    if len(newParts) > 3 {
      conn.Write([]byte("http/1.1 404 not found\r\n\r\n"))
    conn.Close()
    }
    text := newParts[2]
    textlength := len(newParts[2])

    conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",textlength,text)))
    conn.Close()

 
  } else {
conn.Write([]byte("http/1.1 404 not found\r\n\r\n"))
    conn.Close()
  }

   



  



}
