package highway

func (hw *Highway) recv() (head, body []byte, err error) {
	p := make([]byte, 1)
	if _, err = hw.conn.Read(p); err != nil {
		return
	}
	p = make([]byte, 4)
	if _, err = hw.conn.Read(p); err != nil {
		return
	}
	l1 := int(p[0])<<24 | int(p[1])<<16 | int(p[2])<<8 | int(p[3])<<0
	p = make([]byte, 4)
	if _, err = hw.conn.Read(p); err != nil {
		return
	}
	l2 := int(p[0])<<24 | int(p[1])<<16 | int(p[2])<<8 | int(p[3])<<0
	head = make([]byte, l1)
	i, n := 0, 0
	for i < l1 {
		n, err = hw.conn.Read(head[i:])
		if err != nil {
			return
		}
		i += n
	}
	body = make([]byte, l2)
	i, n = 0, 0
	for i < l2 {
		n, err = hw.conn.Read(body[i:])
		if err != nil {
			return
		}
		i += n
	}
	p = make([]byte, 1)
	if _, err = hw.conn.Read(p); err != nil {
		return
	}
	return
}
