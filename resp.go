package sse

import (
	"bufio"
	"net/http"
	"strconv"
)

const (
	textStatusLine = "HTTP/1.1 200\r\n"

	crlf          = "\r\n"
	colonAndSpace = ": "
	commaAndSpace = ", "
)

const (
	toLower = 'a' - 'A'      // for use with OR.
	toUpper = ^byte(toLower) // for use with AND.
	// toLower8 = uint64(toLower) | uint64(toLower)<<8 |
	// 	uint64(toLower)<<16 | uint64(toLower)<<24 |
	// 	uint64(toLower)<<32 | uint64(toLower)<<40 |
	// 	uint64(toLower)<<48 | uint64(toLower)<<56
)

func httpWriteResponseUpgrade(bw *bufio.Writer) {
	bw.WriteString(textStatusLine)

	httpWriteHeader(bw, "Cache-Control", "no-cache")
	httpWriteHeader(bw, "Connection", "keep-alive")
	httpWriteHeader(bw, "Content-Type", "text/event-stream")
	httpWriteHeader(bw, "Access-Control-Allow-Origin", "*")
	bw.WriteString(crlf)
}

func httpWriteHeader(bw *bufio.Writer, key, value string) {
	bw.WriteString(key)
	bw.WriteString(colonAndSpace)
	bw.WriteString(value)
	bw.WriteString(crlf)
}

func httpWriteHeaderKey(bw *bufio.Writer, key string) {
	bw.WriteString(key)
	bw.WriteString(colonAndSpace)
}

func writeStatusText(bw *bufio.Writer, code int) {
	bw.WriteString("HTTP/1.1 ")
	bw.WriteString(strconv.Itoa(code))
	bw.WriteByte(' ')
	bw.WriteString(http.StatusText(code))
	bw.WriteString(crlf)
	bw.WriteString("Content-Type: text/plain; charset=utf-8")
	bw.WriteString(crlf)
}

func writeErrorText(bw *bufio.Writer, err error) {
	body := err.Error()
	bw.WriteString("Content-Length: ")
	bw.WriteString(strconv.Itoa(len(body)))
	bw.WriteString(crlf)
	bw.WriteString(crlf)
	bw.WriteString(body)
}

// httpError is like the http.Error with additional headers.
func httpError(w http.ResponseWriter, body string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(code)
	w.Write([]byte(body))
}
