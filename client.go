package sse

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

var delim = []byte{':', ' '}

// Event ...
type Event struct {
	// URI  string
	Type string
	Data []byte
}

// Client ...
type Client struct {
	url    string
	client *http.Client
	resp   *bufio.Reader
}

// NewClient ...
func NewClient(url string, client *http.Client) (*Client, error) {

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got response status code %d\n", resp.StatusCode)
	}

	c := &Client{
		url:    url,
		client: client,
		resp:   bufio.NewReader(resp.Body),
	}
	return c, nil
}

func (c *Client) ReadNext() ([]byte, error) {
	bs, err := c.resp.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}

	if len(bs) < 2 {
		return nil, nil
	}

	spl := bytes.Split(bs, delim)
	if len(spl) < 2 {
		return nil, nil
	}

	currEvent := &Event{
		// URI: c.url
	}

	switch string(spl[0]) {
	case "event":
		currEvent.Type = string(bytes.TrimSpace(spl[1]))
	case "data":
		currEvent.Data = bytes.TrimSpace(spl[1])
		// evCh <- currEvent
	}
	if err == io.EOF {
		return nil, err
	}
	return nil, nil
}
