package sse

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"unsafe"
)

var (
	textHeadBadRequest          = statusText(http.StatusBadRequest)
	textHeadInternalServerError = statusText(http.StatusInternalServerError)
	textHeadUpgradeRequired     = statusText(http.StatusUpgradeRequired)

	textTailErrHandshakeBadProtocol   = "ErrHandshakeBadProtocol"
	textTailErrHandshakeBadMethod     = "ErrHandshakeBadMethod"
	textTailErrHandshakeBadHost       = "ErrHandshakeBadHost"
	textTailErrHandshakeBadUpgrade    = "ErrHandshakeBadUpgrade"
	textTailErrHandshakeBadConnection = "ErrHandshakeBadConnection"
	textTailErrHandshakeBadSecAccept  = "ErrHandshakeBadSecAccept"
	textTailErrHandshakeBadSecKey     = "ErrHandshakeBadSecKey"
	textTailErrHandshakeBadSecVersion = "ErrHandshakeBadSecVersion"
	textTailErrUpgradeRequired        = "ErrHandshakeUpgradeRequired"
)

// pow for integers implementation.
// See Donald Knuth, The Art of Computer Programming, Volume 2, Section 4.6.3
func pow(a, b int) int {
	p := 1
	for b > 0 {
		if b&1 != 0 {
			p *= a
		}
		b >>= 1
		a *= a
	}
	return p
}

func btrim(bts []byte) []byte {
	var i, j int
	for i = 0; i < len(bts) && (bts[i] == ' ' || bts[i] == '\t'); {
		i++
	}
	for j = len(bts); j > i && (bts[j-1] == ' ' || bts[j-1] == '\t'); {
		j--
	}
	return bts[i:j]
}

// asciiToInt converts bytes to int.
func asciiToInt(bts []byte) (ret int, err error) {
	// ASCII numbers all start with the high-order bits 0011.
	// If you see that, and the next bits are 0-9 (0000 - 1001) you can grab those
	// bits and interpret them directly as an integer.
	var n int
	if n = len(bts); n < 1 {
		return 0, fmt.Errorf("converting empty bytes to int")
	}
	for i := 0; i < n; i++ {
		if bts[i]&0xf0 != 0x30 {
			return 0, fmt.Errorf("%s is not a numeric character", string(bts[i]))
		}
		ret += int(bts[i]&0xf) * pow(10, n-i-1)
	}
	return ret, nil
}

func bsplit3(buf []byte, sep byte) (b1, b2, b3 []byte) {
	a := bytes.IndexByte(buf, sep)
	b := bytes.IndexByte(buf[a+1:], sep)
	if a == -1 || b == -1 {
		return buf, nil, nil
	}
	b += a + 1
	return buf[:a], buf[a+1 : b], buf[b+1:]
}

func b2s(bts []byte) (str string) {
	return *(*string)(unsafe.Pointer(&bts))
}

// readLine reads line from br. It reads until '\n' and returns bytes without
// '\n' or '\r\n' at the end.
// It returns err if and only if line does not end in '\n'. Note that read
// bytes returned in any case of error.
//
// It is much like the textproto/Reader.ReadLine() except the thing that it
// returns raw bytes, instead of string. That is, it avoids copying bytes read
// from br.
//
// textproto/Reader.ReadLineBytes() is also makes copy of resulting bytes to be
// safe with future I/O operations on br.
//
// We could control I/O operations on br and do not need to make additional
// copy for safety.
//
// NOTE: it may return copied flag to notify that returned buffer is safe to
// use.
func readLine(br *bufio.Reader) ([]byte, error) {
	var line []byte
	for {
		bts, err := br.ReadSlice('\n')
		if err == bufio.ErrBufferFull {
			// Copy bytes because next read will discard them.
			line = append(line, bts...)
			continue
		}

		// Avoid copy of single read.
		if line == nil {
			line = bts
		} else {
			line = append(line, bts...)
		}
		if err != nil {
			return line, err
		}

		// Size of line is at least 1.
		// In other case bufio.ReadSlice() returns error.
		n := len(line)

		// Cut '\n' or '\r\n'.
		if n > 1 && line[n-2] == '\r' {
			line = line[:n-2]
		} else {
			line = line[:n-1]
		}
		return line, nil
	}
}

// statusText is a non-performant status text generator.
// NOTE: Used only to generate constants.
func statusText(code int) string {
	var buf bytes.Buffer
	bw := bufio.NewWriter(&buf)
	writeStatusText(bw, code)
	bw.Flush()
	return buf.String()
}

type httpRequestLine struct {
	method, uri  []byte
	major, minor int
}

// httpParseRequestLine parses http request line like "GET / HTTP/1.0".
func httpParseRequestLine(line []byte) (httpRequestLine, error) {
	var req httpRequestLine
	var proto []byte
	req.method, req.uri, proto = bsplit3(line, ' ')

	var ok bool
	req.major, req.minor, ok = httpParseVersion(proto)
	if !ok {
		return req, errors.New("ErrMalformedRequest")
	}
	return req, nil
}

var (
	httpVersion1_0    = []byte("HTTP/1.0")
	httpVersion1_1    = []byte("HTTP/1.1")
	httpVersion2      = []byte("HTTP/2")
	httpVersionPrefix = []byte("HTTP/")
)

// httpParseVersion parses major and minor version of HTTP protocol. It returns
// parsed values and true if parse is ok.
func httpParseVersion(buf []byte) (major, minor int, ok bool) {
	switch {
	case bytes.Equal(buf, httpVersion2):
		return 2, 0, true
	case bytes.Equal(buf, httpVersion1_0):
		return 1, 0, true
	case bytes.Equal(buf, httpVersion1_1):
		return 1, 1, true
	case len(buf) < 8:
		return 0, 0, false
	case !bytes.Equal(buf[:5], httpVersionPrefix):
		return 0, 0, false
	}
	buf = buf[5:]

	dot := bytes.IndexByte(buf, '.')
	if dot == -1 {
		return 0, 0, false
	}

	var err error
	major, err = asciiToInt(buf[:dot])
	if err != nil {
		return 0, 0, false
	}
	minor, err = asciiToInt(buf[dot+1:])
	if err != nil {
		return 0, 0, false
	}
	return major, minor, true
}

// httpParseHeaderLine parses HTTP header as key-value pair. It returns parsed
// values and true if parse is ok.
func httpParseHeaderLine(line []byte) (k, v []byte, ok bool) {
	colon := bytes.IndexByte(line, ':')
	if colon == -1 {
		return
	}

	k = btrim(line[:colon])
	// TODO(gobwas): maybe use just lower here?
	canonicalizeHeaderKey(k)

	v = btrim(line[colon+1:])

	return k, v, true
}

// Algorithm below is like standard textproto/CanonicalMIMEHeaderKey, except
// that it operates with slice of bytes and modifies it inplace without copying.
func canonicalizeHeaderKey(k []byte) {
	upper := true
	for i, c := range k {
		if upper && 'a' <= c && c <= 'z' {
			k[i] &= toUpper
		} else if !upper && 'A' <= c && c <= 'Z' {
			k[i] |= toLower
		}
		upper = c == '-'
	}
}

func httpWriteResponseError(bw *bufio.Writer, err error, code int) {
	switch code {
	case http.StatusBadRequest:
		bw.WriteString(textHeadBadRequest)
	case http.StatusInternalServerError:
		bw.WriteString(textHeadInternalServerError)
	case http.StatusUpgradeRequired:
		bw.WriteString(textHeadUpgradeRequired)
	default:
		writeStatusText(bw, code)
	}

	// Write custom headers.
	// if header != nil {
	// 	header(bw)
	// }

	switch err.Error() {
	case "ErrHandshakeBadProtocol":
		bw.WriteString(textTailErrHandshakeBadProtocol)
	case "ErrHandshakeBadMethod":
		bw.WriteString(textTailErrHandshakeBadMethod)
	case "ErrHandshakeBadHost":
		bw.WriteString(textTailErrHandshakeBadHost)
	case "ErrHandshakeBadUpgrade":
		bw.WriteString(textTailErrHandshakeBadUpgrade)
	case "ErrHandshakeBadConnection":
		bw.WriteString(textTailErrHandshakeBadConnection)
	// case "ErrHandshakeBadSecAccept":
	// 	bw.WriteString(textTailErrHandshakeBadSecAccept)
	// case "ErrHandshakeBadSecKey":
	// 	bw.WriteString(textTailErrHandshakeBadSecKey)
	// case "ErrHandshakeBadSecVersion":
	// 	bw.WriteString(textTailErrHandshakeBadSecVersion)
	case "ErrHandshakeUpgradeRequired":
		bw.WriteString(textTailErrUpgradeRequired)
	// case nil:
	// 	bw.WriteString(crlf)
	default:
		writeErrorText(bw, err)
	}
}
