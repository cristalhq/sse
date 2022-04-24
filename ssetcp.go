package sse

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (u Upgrader) Upgrade(conn io.ReadWriter) (*Stream, error) {
	// return nil, nil

	br := bufio.NewReaderSize(conn, 1000)
	// br := pbufio.GetReader(conn,
	// nonZero(u.ReadBufferSize, DefaultServerReadBufferSize),
	// nonZero(10, 100), //DefaultServerReadBufferSize),
	// )
	bw := bufio.NewWriterSize(conn, 1000)
	brw := bufio.NewReadWriter(br, bw)
	// bw := pbufio.GetWriter(conn,
	// 	// nonZero(u.WriteBufferSize, DefaultServerWriteBufferSize),
	// 	nonZero(10, 100), //DefaultServerWriteBufferSize),
	// )

	// defer func() {
	// pbufio.PutReader(br)
	// pbufio.PutWriter(bw)
	// }()

	// Read HTTP request line like "GET /sse HTTP/2".
	rl, err := readLine(br)
	if err != nil {
		return nil, err
	}
	req, err := httpParseRequestLine(rl)
	if err != nil {
		return nil, err
	}

	// Parse and check HTTP request.
	// As RFC6455 says:
	//   The client's opening handshake consists of the following parts. If the
	//   server, while reading the handshake, finds that the client did not
	//   send a handshake that matches the description below (note that as per
	//   [RFC2616], the order of the header fields is not important), including
	//   but not limited to any violations of the ABNF grammar specified for
	//   the components of the handshake, the server MUST stop processing the
	//   client's handshake and return an HTTP response with an appropriate
	//   error code (such as 400 Bad Request).
	//
	// See https://tools.ietf.org/html/rfc6455#section-4.2.1

	// An HTTP/1.1 or higher GET request, including a "Request-URI".
	//
	// Even if RFC says "1.1 or higher" without mentioning the part of the
	// version, we apply it only to minor part.
	switch {
	case req.major != 1 || req.minor < 1:
		// Abort processing the whole request because we do not even know how to actually parse it.
		err = errors.New("ErrHandshakeBadProtocol")

	case b2s(req.method) != http.MethodGet:
		err = errors.New("ErrHandshakeBadMethod")

	default:
		// if onRequest := u.OnRequest; onRequest != nil {
		// 	err = onRequest(req.uri)
		// }
	}

	// Start headers read/parse loop.
	var (
	// headerSeen reports which header was seen by setting corresponding bit on.
	// headerSeen byte

	// nonceSize = 24 // base64.StdEncoding.EncodedLen(nonceKeySize)
	// nonce = make([]byte, nonceSize)
	)
	for err == nil {
		line, e := readLine(br)
		if e != nil {
			return nil, e
		}
		if len(line) == 0 {
			// Blank line, no more lines to read.
			break
		}

		k, v, ok := httpParseHeaderLine(line)
		if !ok {
			err = errors.New("ErrMalformedRequest")
			break
		}
		_ = v

		switch b2s(k) {
		// case headerHostCanonical:
		// 	headerSeen |= headerSeenHost
		// 	if onHost := u.OnHost; onHost != nil {
		// 		err = onHost(v)
		// 	}

		// case headerUpgradeCanonical:
		// 	headerSeen |= headerSeenUpgrade
		// 	if !bytes.Equal(v, specHeaderValueUpgrade) && !btsEqualFold(v, specHeaderValueUpgrade) {
		// 		err = ErrHandshakeBadUpgrade
		// 	}

		// case headerConnectionCanonical:
		// 	headerSeen |= headerSeenConnection
		// 	if !bytes.Equal(v, specHeaderValueConnection) && !btsHasToken(v, specHeaderValueConnectionLower) {
		// 		err = ErrHandshakeBadConnection
		// 	}

		// case headerSecVersionCanonical:
		// 	headerSeen |= headerSeenSecVersion
		// 	if !bytes.Equal(v, specHeaderValueSecVersion) {
		// 		err = ErrHandshakeUpgradeRequired
		// 	}

		// case headerSecKeyCanonical:
		// 	headerSeen |= headerSeenSecKey
		// 	if len(v) != nonceSize {
		// 		err = ErrHandshakeBadSecKey
		// 	} else {
		// 		copy(nonce[:], v)
		// 	}

		// case headerSecProtocolCanonical:
		// 	if custom, check := u.ProtocolCustom, u.Protocol; hs.Protocol == "" && (custom != nil || check != nil) {
		// 		var ok bool
		// 		if custom != nil {
		// 			hs.Protocol, ok = custom(v)
		// 		} else {
		// 			hs.Protocol, ok = btsSelectProtocol(v, check)
		// 		}
		// 		if !ok {
		// 			err = ErrMalformedRequest
		// 		}
		// 	}

		// case headerSecExtensionsCanonical:
		// 	if custom, check := u.ExtensionCustom, u.Extension; custom != nil || check != nil {
		// 		var ok bool
		// 		if custom != nil {
		// 			hs.Extensions, ok = custom(v, hs.Extensions)
		// 		} else {
		// 			hs.Extensions, ok = btsSelectExtensions(v, hs.Extensions, check)
		// 		}
		// 		if !ok {
		// 			err = ErrMalformedRequest
		// 		}
		// 	}

		default:
			// if onHeader := u.OnHeader; onHeader != nil {
			// 	err = onHeader(k, v)
			// }
		}
	}
	switch {
	case err == nil: //s&& headerSeen != headerSeenAll:
		// switch {
		// case (headerSeen & headerSeenHost) == 0:
		// 	// As RFC2616 says:
		// 	//   A client MUST include a Host header field in all HTTP/1.1
		// 	//   request messages. If the requested URI does not include an
		// 	//   Internet host name for the service being requested, then the
		// 	//   Host header field MUST be given with an empty value. An
		// 	//   HTTP/1.1 proxy MUST ensure that any request message it
		// 	//   forwards does contain an appropriate Host header field that
		// 	//   identifies the service being requested by the proxy. All
		// 	//   Internet-based HTTP/1.1 servers MUST respond with a 400 (Bad
		// 	//   Request) status code to any HTTP/1.1 request message which
		// 	//   lacks a Host header field.
		// 	err = errors.New("ErrHandshakeBadHost")
		// case (headerSeen & headerSeenUpgrade) == 0:
		// 	err = errors.New("ErrHandshakeBadUpgrade")
		// case (headerSeen & headerSeenConnection) == 0:
		// 	err = errors.New("ErrHandshakeBadConnection")
		// case (headerSeen & headerSeenSecVersion) == 0:
		// 	// In case of empty or not present version we do not send 426 status,
		// 	// because it does not meet the ABNF rules of RFC6455:
		// 	//
		// 	// version = DIGIT | (NZDIGIT DIGIT) |
		// 	// ("1" DIGIT DIGIT) | ("2" DIGIT DIGIT)
		// 	// ; Limited to 0-255 range, with no leading zeros
		// 	//
		// 	// That is, if version is really invalid – we sent 426 status as above, if it
		// 	// not present – it is 400.
		// 	err = errors.New("ErrHandshakeBadSecVersion")
		// case (headerSeen & headerSeenSecKey) == 0:
		// 	err = errors.New("ErrHandshakeBadSecKey")
		// default:
		// 	panic("unknown headers state")
		// }

		// case err == nil && u.OnBeforeUpgrade != nil:
		// 	header[1], err = u.OnBeforeUpgrade()
	}
	if err != nil {
		fmt.Printf("omg %#v\n", err)
		var code int
		// if rej, ok := err.(*rejectConnectionError); ok {
		// 	code = rej.code
		// 	header[1] = rej.header
		// }
		if code == 0 {
			code = http.StatusInternalServerError
		}
		httpWriteResponseError(bw, err, code)
		// Do not store Flush() error to not override already existing one.
		bw.Flush()
		return nil, err
	}

	httpWriteResponseUpgrade(bw)
	if err := bw.Flush(); err != nil {
		return nil, err
	}

	s := &Stream{
		bw: brw,
		w:  bw,
	}
	return s, nil
}
