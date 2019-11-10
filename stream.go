package sse

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strconv"
)

type Stream struct {
	conn      net.Conn
	w         io.Writer
	flusher   http.Flusher
	autoFlush bool
}

type BinaryMarshaler interface {
	MarshalBinary() ([]byte, error)
}

func (s *Stream) Flush() {
	s.flusher.Flush()
}

func (s *Stream) WriteJSON(id, event string, v interface{}) error {
	raw, err := json.Marshal(v)
	if err != nil {
		return err
	}

	data := encode(id, event, raw)
	return s.write(data)
}

func (s *Stream) WriteMessage(id, event string, message BinaryMarshaler) error {
	raw, err := message.MarshalBinary()
	if err != nil {
		return err
	}

	data := encode(id, event, raw)
	return s.write(data)
}

func (s *Stream) WriteBytes(id, event string, raw []byte) error {
	data := encode(id, event, raw)
	return s.write(data)
}

func (s *Stream) WriteInt(id, event string, num int64) error {
	str := strconv.FormatInt(num, 10)
	data := encode(id, event, []byte(str))
	return s.write(data)
}

func (s *Stream) WriteRaw(data []byte) error {
	return s.write(data)
}

func (s *Stream) write(data []byte) error {
	_, err := s.w.Write(data)
	if s.autoFlush {
		s.flusher.Flush()
	}
	return err
}

// httpError is like the http.Error with additional headers.
func httpError(w http.ResponseWriter, body string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(code)
	w.Write([]byte(body))
}
