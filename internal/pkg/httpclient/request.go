package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ContentType string

const (
	ContentTypeJSON      ContentType = "application/json"
	ContentTypeXML       ContentType = "application/xml"
	ContentTypeForm      ContentType = "application/x-www-form-urlencoded"
	ContentTypeFormData  ContentType = "application/x-www-form-urlencoded"
	ContentTypeURLEncode ContentType = "application/x-www-form-urlencoded"
	ContentTypeHTML      ContentType = "text/html"
	ContentTypeText      ContentType = "text/plain"
	ContentTypeMultipart ContentType = "multipart/form-data"
)

const (
	Post    = http.MethodPost
	Get     = http.MethodGet
	Head    = http.MethodHead
	Put     = http.MethodPut
	Delete  = http.MethodDelete
	Patch   = http.MethodPatch
	Options = http.MethodOptions
)

type Request struct {
	target  string
	method  string
	tracer  trace.Tracer
	header  http.Header
	params  url.Values
	query   url.Values
	body    io.Reader
	Client  *http.Client
	cookies []*http.Cookie
	files   map[string]*file
}

type file struct {
	name     string
	filename string
	filepath string
}

func (r *Request) SetContentType(contentType ContentType) {
	r.SetHeader("Content-Type", string(contentType))
}

func (r *Request) AddHeader(key, value string) {
	r.header.Add(key, value)
}

func (r *Request) DelHeader(key string) {
	r.header.Del(key)
}

func (r *Request) SetHeader(key, value string) {
	r.header.Set(key, value)
}

func (r *Request) SetHeaders(header http.Header) {
	r.header = header
}

func (r *Request) SetBody(body io.Reader) {
	r.body = body
}

func (r *Request) AddParam(key, value string) {
	r.params.Add(key, value)
}

func (r *Request) DelParam(key string) {
	r.params.Del(key)
}

func (r *Request) SetParam(key, value string) {
	r.params.Set(key, value)
}

func (r *Request) SetParams(params url.Values) {
	r.params = params
}

func (r *Request) AddQuery(key, value string) {
	r.query.Add(key, value)
}

func (r *Request) DelQuery(key string) {
	r.query.Del(key)
}

func (r *Request) SetQuery(key, value string) {
	r.query.Set(key, value)
}

func (r *Request) AddFile(name, filename, filepath string) {
	if r.files == nil {
		r.files = make(map[string]*file)
	}
	if filename == "" {
		filename = name
	}
	r.files[name] = &file{name, filename, filepath}
}

func (r *Request) DelFile(name string) {
	if r.files != nil {
		delete(r.files, name)
	}
}

func (r *Request) AddCookie(cookie *http.Cookie) {
	r.cookies = append(r.cookies, cookie)
}

func (r *Request) do(ctx context.Context) (*http.Response, error) {
	var req *http.Request
	var err error
	var body io.Reader
	var transform bool

	if r.method == http.MethodGet ||
		r.method == http.MethodTrace ||
		r.method == http.MethodOptions ||
		r.method == http.MethodHead ||
		r.method == http.MethodDelete {
		transform = true
	}

	if r.body != nil {
		body = r.body
	} else if len(r.files) > 0 {
		var bodyBuffer = &bytes.Buffer{}
		var bodyWriter = multipart.NewWriter(bodyBuffer)

		for _, file := range r.files {
			fileContent, err := os.ReadFile(file.filepath)
			if err != nil {
				return nil, err
			}
			fileWriter, err := bodyWriter.CreateFormFile(file.name, file.filename)
			if err != nil {
				return nil, err
			}
			if _, err = fileWriter.Write(fileContent); err != nil {
				return nil, err
			}
		}
		for key, values := range r.params {
			for _, value := range values {
				bodyWriter.WriteField(key, value)
			}
		}

		if err = bodyWriter.Close(); err != nil {
			return nil, err
		}

		r.SetContentType(ContentType(bodyWriter.FormDataContentType()))
		body = bodyBuffer
	} else if len(r.params) > 0 && !transform {
		body = strings.NewReader(r.params.Encode())
	}

	req, err = http.NewRequestWithContext(ctx, r.method, r.target, body)
	if err != nil {
		return nil, err
	}

	if transform {
		for key, values := range r.params {
			for _, value := range values {
				r.query.Add(key, value)
			}
		}
	}

	req.URL.RawQuery = r.query.Encode()
	req.Header = r.header

	for _, cookie := range r.cookies {
		req.AddCookie(cookie)
	}

	return r.Client.Do(req)
}

func (r *Request) Exec(ctx context.Context) *Response {
	kind := trace.SpanKindClient
	ctx, span := r.tracer.Start(ctx,
		r.target,
		trace.WithAttributes(
			attribute.String("target", r.target),
			attribute.String("method", r.method),
			attribute.String("params", r.params.Encode()),
			attribute.String("type", "req"),
		),
		trace.WithSpanKind(kind),
	)

	rsp, err := r.do(ctx)
	defer func() {
		if err != nil {
			span.RecordError(err)
		}
		span.End()
	}()

	if rsp != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		return &Response{rsp, nil, err}
	}
	data, err := io.ReadAll(rsp.Body)
	return &Response{rsp, data, err}
}

func (r *Request) Download(ctx context.Context, savePath string) *Response {
	kind := trace.SpanKindClient
	ctx, span := r.tracer.Start(ctx,
		r.target,
		trace.WithAttributes(
			attribute.String("target", r.target),
			attribute.String("method", r.method),
			attribute.String("params", r.params.Encode()),
			attribute.String("type", "download"),
		),
		trace.WithSpanKind(kind),
	)

	rsp, err := r.do(ctx)
	defer func() {
		if err != nil {
			span.RecordError(err)
		}
		span.End()
	}()

	if rsp != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		return &Response{rsp, nil, err}
	}

	nFile, err := os.Create(savePath)
	if err != nil {
		return &Response{nil, nil, err}
	}
	defer nFile.Close()

	buf := make([]byte, 1024)
	for {
		size, err := rsp.Body.Read(buf)
		if size == 0 || err != nil {
			break
		}
		nFile.Write(buf[:size])
	}
	data := []byte(savePath)
	return &Response{rsp, data, err}
}

// WriteJSON 将一个对象序列化为 JSON 字符串，并将其作为 http 请求的 body 发送给服务器端。
func (r *Request) WriteJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	r.SetBody(bytes.NewReader(data))
	r.SetContentType(ContentTypeJSON)
	return nil
}
