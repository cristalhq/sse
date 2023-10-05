package sse

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

type Stream struct {
	bw *bufio.ReadWriter
	w  io.Writer
	nc net.Conn
}

type BinaryMarshaler interface {
	MarshalBinary() ([]byte, error)
}

type TextMarshaler interface {
	MarshalText() ([]byte, error)
}

func (s *Stream) Flush() error {
	return s.bw.Flush()
}

// Close sends close event with empty data and closes underlying connection.
func (s *Stream) Close() error {
	defer s.nc.Close()

	_, err := s.bw.Write([]byte("event:close\ndata:\n\n"))
	if err != nil {
		return err
	}
	return s.Flush()
}

func (s *Stream) SetRetry(retry time.Duration) error {
	data := fmt.Sprintf("retry:%v\n", retry.Milliseconds())
	_, err := s.bw.Write([]byte(data))
	return err
}

func (s *Stream) SetID(id int64) error {
	str := strconv.FormatInt(id, 10)
	data := []byte("id:" + str + "\n")
	_, err := s.bw.Write(data)
	return err
}

func (s *Stream) SetEvent(event string) error {
	data := []byte("event:" + event + "\n")
	_, err := s.bw.Write(data)
	return err
}

func (s *Stream) WriteEvent(id int64, event string, data []byte) error {
	if err := s.SetID(id); err != nil {
		return err
	}
	if err := s.SetEvent(event); err != nil {
		return err
	}
	if err := s.WriteBytes(data); err != nil {
		return err
	}
	return s.Flush()
}

func (s *Stream) WriteJSON(v interface{}) error {
	raw, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.writeData(raw)
}

func (s *Stream) WriteBinary(message BinaryMarshaler) error {
	data, err := message.MarshalBinary()
	if err != nil {
		return err
	}
	return s.writeData(data)
}

func (s *Stream) WriteText(message TextMarshaler) error {
	text, err := message.MarshalText()
	if err != nil {
		return err
	}
	return s.writeData(text)
}

func (s *Stream) WriteBytes(data []byte) error {
	return s.writeData(data)
}

func (s *Stream) WriteString(data string) error {
	return s.writeData([]byte(data))
}

func (s *Stream) WriteInt(num int64) error {
	str := strconv.FormatInt(num, 10)
	return s.writeData([]byte(str))
}

func (s *Stream) WriteFloat(num float64) error {
	str := strconv.FormatFloat(num, 'f', 5, 64)
	return s.writeData([]byte(str))
}

func (s *Stream) writeData(data []byte) error {
	size := 6 + 1 + len(data) // size of "data\n\n" + ":{data}"
	buf := make([]byte, 0, size)

	buf = append(buf, "data"...)
	if len(data) > 0 {
		buf = append(buf, ':')
		buf = append(buf, data...)
	}
	buf = append(buf, []byte("\n\n")...)

	_, err := s.bw.Write(buf)
	if err != nil {
		return err
	}
	return s.bw.Flush()
}
