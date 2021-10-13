package render

import "net/url"

type ResponseWriter interface {
	Header() url.Values
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}
