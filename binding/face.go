package binding

import (
	"io"
	"net/url"
)

type RequestTransportData interface {
	Head() url.Values
	Body() io.ReadCloser
}
