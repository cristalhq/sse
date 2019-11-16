package sse

import (
	"bufio"
)

const (
	textStatusLine = "HTTP/1.1 200\r\n"

	crlf          = "\r\n"
	colonAndSpace = ": "
	commaAndSpace = ", "
)

func httpWriteResponseUpgrade(bw *bufio.Writer, nonce []byte, protocol string) {
	bw.WriteString(textStatusLine)

	httpWriteHeader(bw, "Cache-Control", "no-cache")
	httpWriteHeader(bw, "Connection", "keep-alive")
	httpWriteHeader(bw, "Content-Type", "text/event-stream")
	bw.WriteString(crlf)
}

func httpWriteHeader(bw *bufio.Writer, key, value string) {
	bw.WriteString(key)
	bw.WriteString(colonAndSpace)
	bw.WriteString(value)
	bw.WriteString(crlf)
}
