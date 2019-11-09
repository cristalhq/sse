package sse

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNotFlusher  = Error("sse: not an HTTP-flusher")
	ErrNotHijacker = Error("sse: not an HTTP-hijacker")
)
