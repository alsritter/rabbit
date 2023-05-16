package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Response struct {
	*http.Response
	data  []byte
	error error
}

func (r *Response) Status() string {
	if r.Response != nil {
		return r.Response.Status
	}
	return fmt.Sprintf("%d ServiceUnavailable", http.StatusServiceUnavailable)
}

func (r *Response) StatusCode() int {
	if r.Response != nil {
		return r.Response.StatusCode
	}
	return http.StatusServiceUnavailable
}

func (r *Response) Proto() string {
	if r.Response != nil {
		return r.Response.Proto
	}
	return ""
}

func (r *Response) ProtoMajor() int {
	if r.Response != nil {
		return r.Response.ProtoMajor
	}
	return 1
}

func (r *Response) ProtoMinor() int {
	if r.Response != nil {
		return r.Response.ProtoMinor
	}
	return 0
}

func (r *Response) Header() http.Header {
	if r.Response != nil {
		return r.Response.Header
	}
	return http.Header{}
}

func (r *Response) ContentLength() int64 {
	if r.Response != nil {
		return r.Response.ContentLength
	}
	return 0
}

func (r *Response) TransferEncoding() []string {
	if r.Response != nil {
		return r.Response.TransferEncoding
	}
	return nil
}

func (r *Response) Close() bool {
	if r.Response != nil {
		return r.Response.Close
	}
	return true
}

func (r *Response) Uncompressed() bool {
	if r.Response != nil {
		return r.Response.Uncompressed
	}
	return true
}

func (r *Response) Trailer() http.Header {
	if r.Response != nil {
		return r.Response.Trailer
	}
	return http.Header{}
}

func (r *Response) Request() *http.Request {
	if r.Response != nil {
		return r.Response.Request
	}
	return nil
}

func (r *Response) TLS() *tls.ConnectionState {
	if r.Response != nil {
		return r.Response.TLS
	}
	return nil
}

func (r *Response) Cookies() []*http.Cookie {
	if r.Response != nil {
		return r.Response.Cookies()
	}
	return nil
}

func (r *Response) Location() (*url.URL, error) {
	if r.Response != nil {
		return r.Response.Location()
	}
	return nil, nil
}

func (r *Response) ProtoAtLeast(major, minor int) bool {
	if r.Response != nil {
		return r.Response.ProtoAtLeast(major, minor)
	}
	return false
}

func (r *Response) Write(w io.Writer) error {
	if r.Response != nil {
		return r.Response.Write(w)
	}
	return nil
}

func (r *Response) Error() error {
	return r.error
}

func (r *Response) Reader() io.Reader {
	return bytes.NewReader(r.data)
}

func (r *Response) Bytes() ([]byte, error) {
	return r.data, r.error
}

func (r *Response) MustBytes() []byte {
	return r.data
}

func (r *Response) String() (string, error) {
	return string(r.data), r.error
}

func (r *Response) MustString() string {
	return string(r.data)
}

func (r *Response) Unmarshal(v any) error {
	if r.error != nil {
		return r.error
	}
	return json.Unmarshal(r.data, &v)
}
