package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Server implementation
func startServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("Server listening on %s", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Read the raw HTTP request
	rawRequest, isChunked, err := readHttpMessage(reader)
	if err != nil {
		log.Printf("Error reading request: %v", err)
		return
	}

	// Display the raw request
	fmt.Println("\n=== RAW HTTP REQUEST ===")
	fmt.Println(rawRequest)

	if isChunked {
		fmt.Println("(Using chunked encoding)")
	}

	// Parse and display headers and body
	headers, body, err := parseHttpMessage(rawRequest)
	if err != nil {
		log.Printf("Error parsing request: %v", err)
		return
	}

	fmt.Println("\n=== PARSED REQUEST ===")
	fmt.Println("Headers:")
	for k, v := range headers {
		fmt.Printf("%s: %s\n", k, v)
	}
	fmt.Println("\nBody:")
	fmt.Println(body)

	// Send a chunked response
	responseHeader := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Date: " + time.Now().Format(time.RFC1123) + "\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"\r\n"

	conn.Write([]byte(responseHeader))

	time.Sleep(5 * time.Second)

	// Send first chunk
	chunk1 := "Hello,"
	writeChunk(conn, chunk1)

	// Small delay to demonstrate chunking
	time.Sleep(500 * time.Millisecond)

	// Send second chunk
	chunk2 := " World!"
	writeChunk(conn, chunk2)

	// Send final empty chunk with trailers
	conn.Write([]byte("0\r\nX-Trailer: Example-Trailer\r\n\r\n"))
}

// Client implementation
func sendRequest(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Prepare a raw HTTP request with chunked encoding
	requestHeader := "GET / HTTP/1.1\r\n" +
		"Host: localhost\r\n" +
		"User-Agent: RawHTTPClient\r\n" +
		"Accept: */*\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"\r\n"

	// Send the request header
	fmt.Println("\n=== SENDING CHUNKED REQUEST ===")
	fmt.Println(requestHeader)
	conn.Write([]byte(requestHeader))

	// Send first chunk
	chunk1 := "Hello"
	fmt.Printf("Sending chunk: %s\n", chunk1)
	writeChunk(conn, chunk1)

	// Send second chunk
	chunk2 := " there"
	fmt.Printf("Sending chunk: %s\n", chunk2)
	writeChunk(conn, chunk2)

	// Send final empty chunk
	fmt.Println("Sending end chunk")
	conn.Write([]byte("0\r\n\r\n"))

	// Read the response
	reader := bufio.NewReader(conn)
	rawResponse, isChunked, err := readHttpMessage(reader)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// Display the raw response
	fmt.Println("\n=== RAW HTTP RESPONSE ===")
	fmt.Println(rawResponse)

	if isChunked {
		fmt.Println("(Using chunked encoding)")
	}

	// Parse and display headers and body
	headers, body, err := parseHttpMessage(rawResponse)
	if err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}

	fmt.Println("\n=== PARSED RESPONSE ===")
	fmt.Println("Headers:")
	for k, v := range headers {
		fmt.Printf("%s: %s\n", k, v)
	}
	fmt.Println("\nBody:")
	fmt.Println(body)
}

// Write a chunk in the proper format: size in hex + CRLF + data + CRLF
func writeChunk(writer io.Writer, data string) {
	chunkHeader := fmt.Sprintf("%x\r\n", len(data))
	writer.Write([]byte(chunkHeader))
	writer.Write([]byte(data))
	writer.Write([]byte("\r\n"))
}

// Read HTTP message handling both regular and chunked encoding
func readHttpMessage(reader *bufio.Reader) (string, bool, error) {
	var message strings.Builder
	var isChunked bool

	// Read headers
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", false, err
		}

		message.WriteString(line)

		// Check if chunked encoding
		if strings.HasPrefix(strings.ToLower(line), "transfer-encoding: chunked") {
			isChunked = true
		}

		// Empty line indicates end of headers
		if line == "\r\n" || line == "\n" {
			break
		}
	}

	// If chunked, read chunks
	if isChunked {
		for {
			// Read chunk size line
			sizeLine, err := reader.ReadString('\n')
			if err != nil {
				return "", true, err
			}

			message.WriteString(sizeLine)

			// Parse chunk size
			size, err := strconv.ParseInt(strings.TrimSpace(sizeLine), 16, 64)
			if err != nil {
				return "", true, fmt.Errorf("invalid chunk size: %v", err)
			}

			// Zero size means end of body
			if size == 0 {
				// Read trailers or final CRLF
				for {
					trailerLine, err := reader.ReadString('\n')
					if err != nil {
						return "", true, err
					}

					message.WriteString(trailerLine)

					if trailerLine == "\r\n" || trailerLine == "\n" {
						break
					}
				}
				break
			}

			// Read chunk data
			chunk := make([]byte, size)
			_, err = io.ReadFull(reader, chunk)
			if err != nil {
				return "", true, err
			}
			message.Write(chunk)

			// Read the trailing CRLF
			_, err = reader.ReadString('\n')
			if err != nil {
				return "", true, err
			}
			message.WriteString("\r\n")
		}
	} else {
		// Read regular body until connection closes or content-length is reached
		body, err := io.ReadAll(reader)
		if err != nil && err != io.EOF {
			return "", false, err
		}
		message.Write(body)
	}

	return message.String(), isChunked, nil
}

func parseHttpMessage(raw string) (map[string]string, string, error) {
	headers := make(map[string]string)

	// Split headers from body
	parts := strings.SplitN(raw, "\r\n\r\n", 2)
	if len(parts) != 2 {
		return headers, "", fmt.Errorf("invalid HTTP message format")
	}

	headerLines := strings.Split(parts[0], "\r\n")
	// Skip the first line (request/status line)
	for i := 1; i < len(headerLines); i++ {
		if headerLines[i] == "" {
			continue
		}
		headerParts := strings.SplitN(headerLines[i], ":", 2)
		if len(headerParts) == 2 {
			headers[strings.TrimSpace(headerParts[0])] = strings.TrimSpace(headerParts[1])
		}
	}

	// Handle body
	body := parts[1]

	// If chunked encoding, extract the actual content
	if strings.ToLower(headers["Transfer-Encoding"]) == "chunked" {
		body = extractChunkedBody(body)
	}

	return headers, body, nil
}

func extractChunkedBody(chunkedBody string) string {
	var body strings.Builder
	lines := strings.Split(chunkedBody, "\r\n")

	i := 0
	for i < len(lines) {
		// Parse chunk size
		sizeStr := strings.TrimSpace(lines[i])
		if sizeStr == "" {
			i++
			continue
		}

		// Strip any chunk extension
		if idx := strings.IndexByte(sizeStr, ';'); idx >= 0 {
			sizeStr = sizeStr[:idx]
		}

		size, err := strconv.ParseInt(sizeStr, 16, 64)
		if err != nil || size < 0 {
			return "[Error parsing chunked body]"
		}

		i++

		// Zero chunk size marks the end
		if size == 0 {
			break
		}

		// Append chunk data if available
		if i < len(lines) {
			body.WriteString(lines[i])
			i++
		}
	}

	return body.String()
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: program [server|client] [address]")
	}

	mode := os.Args[1]
	address := "localhost:8080"
	if len(os.Args) > 2 {
		address = os.Args[2]
	}

	switch mode {
	case "server":
		startServer(address)
	case "client":
		sendRequest(address)
	default:
		log.Fatalf("Unknown mode: %s. Use 'server' or 'client'", mode)
	}
}
