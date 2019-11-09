package sse

func encode(id, event string, data []byte) []byte {
	size := 6 // size of "data\n\n"
	if len(event) > 0 {
		size += 6 + len(event) + 1 // size of "event:{event}\n"
	}
	if len(data) > 0 {
		size += 1 + len(data) // size of ":{data}"
	}

	buf := make([]byte, 0, size)
	if len(event) > 0 {
		buf = append(buf, "event:"...)
		buf = append(buf, event...)
		buf = append(buf, '\n')
	}

	buf = append(buf, "data"...)
	if len(data) > 0 {
		buf = append(buf, ':')
		buf = append(buf, data...)
	}
	buf = append(buf, "\n\n"...)
	return buf

}
