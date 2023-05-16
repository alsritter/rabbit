package httpclient

import (
	"net/http"
	"net/url"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

type Client struct {
	tracer trace.Tracer
}

func New(tracer trace.Tracer) *Client {
	return &Client{
		tracer: tracer,
	}
}

func (c *Client) NewRequest(method, target string) *Request {
	var nURL, _ = url.Parse(target)
	var req = &Request{tracer: c.tracer}

	req.method = strings.ToUpper(method)
	req.target = target
	req.params = url.Values{}
	if nURL != nil {
		req.query = nURL.Query()
	} else {
		req.query = url.Values{}
	}
	req.header = http.Header{}
	req.Client = http.DefaultClient
	req.SetContentType(ContentTypeURLEncode)
	return req
}

func (c *Client) NewJSONRequest(method, target string, param interface{}) *Request {
	var r = c.NewRequest(method, target)
	r.WriteJSON(param)
	return r
}
