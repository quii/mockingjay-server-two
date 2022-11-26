package stub

import "net/textproto"

type Headers map[string][]string

func (h Headers) compile() {
	for key, value := range h {
		delete(h, key)
		h[textproto.CanonicalMIMEHeaderKey(key)] = value
	}
}
