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

type TextMarshaler interface {
	MarshalText() ([]byte, error)
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

func (s *Stream) SetRetry(retry time.Duration) error {
	data := fmt.Sprintf(`retry: %v\n`, retry.Milliseconds())
	return s.write([]byte(data))
}

func (s *Stream) WriteJSON(id, event string, v interface{}) error {
	raw, err := json.Marshal(v)
	if err != nil {
		return err
	}

	data := encode(id, event, raw)
	return s.write(data)
}

func (s *Stream) WriteBinary(id, event string, message BinaryMarshaler) error {
	raw, err := message.MarshalBinary()
	if err != nil {
		return err
	}

	data := encode(id, event, raw)
	return s.write(data)
}

func (s *Stream) WriteText(id, event string, message TextMarshaler) error {
	text, err := message.MarshalText()
	if err != nil {
		return err
	}

	data := encode(id, event, text)
	return s.write(data)
}

func (s *Stream) WriteBytes(id, event string, raw []byte) error {
	data := encode(id, event, raw)
	return s.write(data)
}

func (s *Stream) WriteString(id, event string, data string) error {
	raw := encode(id, event, []byte(data))
	return s.write(raw)
}

func (s *Stream) WriteInt(id, event string, num int64) error {
	str := strconv.FormatInt(num, 10)
	data := encode(id, event, []byte(str))
	return s.write(data)
}

func (s *Stream) WriteFloat(id, event string, num float64) error {
	str := strconv.FormatFloat(num, 'f', 5, 64)
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
