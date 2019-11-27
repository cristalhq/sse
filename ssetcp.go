package sse

import (
	"io"
)

func (u Upgrader) Upgrade(conn io.ReadWriter) (*Stream, error) {
	return nil, nil
}
