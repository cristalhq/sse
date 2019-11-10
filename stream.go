package sse

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type Stream struct {
	w       io.Writer
	flusher http.Flusher
}

type BinaryMarshaler interface {
	MarshalBinary() ([]byte, error)
}

func (s *Stream) Flush() {
	s.flusher.Flush()
}

// Close sends close event with empth data.
func (s *Stream) Close() error {
	_, err := s.w.Write([]byte("event:close\ndata:\n\n"))
	s.Flush()
	if err != nil {
		return err
	}
	return nil
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
	s.flusher.Flush()
	return err
}
