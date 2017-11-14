package eval

import (
	"bytes"
	_ "fmt"
	"io"
	"io/ioutil"
	"monkey/ast"
	"net/http"
	"strings"
	"time"
)

const (
	HTTP_OBJ = "HTTP_OBJ"
	HTTPCLIENT_OBJ = "HTTPCLIENT_OBJ"
	HTTPSERVER_OBJ = "HTTPSERVER_OBJ"

	HTTPRESPONSE_OBJ = "HTTPRESPONSE_OBJ"
	HTTPREQUEST_OBJ  = "HTTPREQUEST_OBJ"

	HTTPRESPONSEWRITER_OBJ = "HTTPRESPONSEWRITER_OBJ"
	HTTPHEADER_OBJ         = "HTTPHEADER_OBJ"
)

//HTTP Object
const http_name = "http"

type HttpObj struct {
}

func NewHTTPObj() Object {
	ret := &HttpObj{}
	SetGlobalObj(http_name, ret)

	SetGlobalObj(http_name+".METHOD_GET", NewString(http.MethodGet))
	SetGlobalObj(http_name+".METHOD_HEAD", NewString(http.MethodHead))
	SetGlobalObj(http_name+".METHOD_POST", NewString(http.MethodPost))
	SetGlobalObj(http_name+".METHOD_PUT", NewString(http.MethodPut))
	SetGlobalObj(http_name+".METHOD_PATCH", NewString(http.MethodPatch))
	SetGlobalObj(http_name+".METHOD_DELETE", NewString(http.MethodDelete))
	SetGlobalObj(http_name+".METHOD_CONNECT", NewString(http.MethodConnect))
	SetGlobalObj(http_name+".METHOD_OPTIONS", NewString(http.MethodOptions))
	SetGlobalObj(http_name+".METHOD_TRACE", NewString(http.MethodTrace))

	SetGlobalObj(http_name+".STATUS_CONTINUE", NewInteger(int64(http.StatusContinue)))
	SetGlobalObj(http_name+".STATUS_SWITCHINGPROTOCOLS", NewInteger(int64(http.StatusSwitchingProtocols)))
	SetGlobalObj(http_name+".STATUS_PROCESSING", NewInteger(int64(http.StatusProcessing)))
	SetGlobalObj(http_name+".STATUS_OK", NewInteger(int64(http.StatusOK)))
	SetGlobalObj(http_name+".STATUS_CREATED", NewInteger(int64(http.StatusCreated)))
	SetGlobalObj(http_name+".STATUS_ACCEPTED", NewInteger(int64(http.StatusAccepted)))
	SetGlobalObj(http_name+".STATUS_NONAUTHORITATIVEINFO", NewInteger(int64(http.StatusNonAuthoritativeInfo)))
	SetGlobalObj(http_name+".STATUS_NOCONTENT", NewInteger(int64(http.StatusNoContent)))
	SetGlobalObj(http_name+".STATUS_RESETCONTENT", NewInteger(int64(http.StatusResetContent)))
	SetGlobalObj(http_name+".STATUS_PARTIALCONTENT", NewInteger(int64(http.StatusPartialContent)))
	SetGlobalObj(http_name+".STATUS_MULTISTATUS", NewInteger(int64(http.StatusMultiStatus)))
	SetGlobalObj(http_name+".STATUS_ALREADYREPORTED", NewInteger(int64(http.StatusAlreadyReported)))
	SetGlobalObj(http_name+".STATUS_IMUSED", NewInteger(int64(http.StatusIMUsed)))
	SetGlobalObj(http_name+".STATUS_MULTIPLECHOICES", NewInteger(int64(http.StatusMultipleChoices)))
	SetGlobalObj(http_name+".STATUS_MOVEDPERMANENTLY", NewInteger(int64(http.StatusMovedPermanently)))
	SetGlobalObj(http_name+".STATUS_FOUND", NewInteger(int64(http.StatusFound)))
	SetGlobalObj(http_name+".STATUS_SEEOTHER", NewInteger(int64(http.StatusSeeOther)))
	SetGlobalObj(http_name+".STATUS_NOTMODIFIED", NewInteger(int64(http.StatusNotModified)))
	SetGlobalObj(http_name+".STATUS_USEPROXY", NewInteger(int64(http.StatusUseProxy)))
	SetGlobalObj(http_name+".STATUS_TEMPORARYREDIRECT", NewInteger(int64(http.StatusTemporaryRedirect)))
	SetGlobalObj(http_name+".STATUS_PERMANENTREDIRECT", NewInteger(int64(http.StatusPermanentRedirect)))
	SetGlobalObj(http_name+".STATUS_BADREQUEST", NewInteger(int64(http.StatusBadRequest)))
	SetGlobalObj(http_name+".STATUS_UNAUTHORIZED", NewInteger(int64(http.StatusUnauthorized)))
	SetGlobalObj(http_name+".STATUS_PAYMENTREQUIRED", NewInteger(int64(http.StatusPaymentRequired)))
	SetGlobalObj(http_name+".STATUS_FORBIDDEN", NewInteger(int64(http.StatusForbidden)))
	SetGlobalObj(http_name+".STATUS_NOTFOUND", NewInteger(int64(http.StatusNotFound)))
	SetGlobalObj(http_name+".STATUS_METHODNOTALLOWED", NewInteger(int64(http.StatusMethodNotAllowed)))
	SetGlobalObj(http_name+".STATUS_NOTACCEPTABLE", NewInteger(int64(http.StatusNotAcceptable)))
	SetGlobalObj(http_name+".STATUS_PROXYAUTHREQUIRED", NewInteger(int64(http.StatusProxyAuthRequired)))
	SetGlobalObj(http_name+".STATUS_REQUESTTIMEOUT", NewInteger(int64(http.StatusRequestTimeout)))
	SetGlobalObj(http_name+".STATUS_CONFLICT", NewInteger(int64(http.StatusConflict)))
	SetGlobalObj(http_name+".STATUS_GONE", NewInteger(int64(http.StatusGone)))
	SetGlobalObj(http_name+".STATUS_LENGTHREQUIRED", NewInteger(int64(http.StatusLengthRequired)))
	SetGlobalObj(http_name+".STATUS_PRECONDITIONFAILED", NewInteger(int64(http.StatusPreconditionFailed)))
	SetGlobalObj(http_name+".STATUS_REQUESTENTITYTOOLARGE", NewInteger(int64(http.StatusRequestEntityTooLarge)))
	SetGlobalObj(http_name+".STATUS_REQUESTURITOOLONG", NewInteger(int64(http.StatusRequestURITooLong)))
	SetGlobalObj(http_name+".STATUS_UNSUPPORTEDMEDIATYPE", NewInteger(int64(http.StatusUnsupportedMediaType)))
	SetGlobalObj(http_name+".STATUS_REQUESTEDRANGENOTSATISFIABLE", NewInteger(int64(http.StatusRequestedRangeNotSatisfiable)))
	SetGlobalObj(http_name+".STATUS_EXPECTATIONFAILED", NewInteger(int64(http.StatusExpectationFailed)))
	SetGlobalObj(http_name+".STATUS_TEAPOT", NewInteger(int64(http.StatusTeapot)))
	SetGlobalObj(http_name+".STATUS_UNPROCESSABLEENTITY", NewInteger(int64(http.StatusUnprocessableEntity)))
	SetGlobalObj(http_name+".STATUS_LOCKED", NewInteger(int64(http.StatusLocked)))
	SetGlobalObj(http_name+".STATUS_FAILEDDEPENDENCY", NewInteger(int64(http.StatusFailedDependency)))
	SetGlobalObj(http_name+".STATUS_UPGRADEREQUIRED", NewInteger(int64(http.StatusUpgradeRequired)))
	SetGlobalObj(http_name+".STATUS_PRECONDITIONREQUIRED", NewInteger(int64(http.StatusPreconditionRequired)))
	SetGlobalObj(http_name+".STATUS_TOOMANYREQUESTS", NewInteger(int64(http.StatusTooManyRequests)))
	SetGlobalObj(http_name+".STATUS_REQUESTHEADERFIELDSTOOLARGE", NewInteger(int64(http.StatusRequestHeaderFieldsTooLarge)))
	SetGlobalObj(http_name+".STATUS_UNAVAILABLEFORLEGALREASONS", NewInteger(int64(http.StatusUnavailableForLegalReasons)))
	SetGlobalObj(http_name+".STATUS_INTERNALSERVERERROR", NewInteger(int64(http.StatusInternalServerError)))
	SetGlobalObj(http_name+".STATUS_NOTIMPLEMENTED", NewInteger(int64(http.StatusNotImplemented)))
	SetGlobalObj(http_name+".STATUS_BADGATEWAY", NewInteger(int64(http.StatusBadGateway)))
	SetGlobalObj(http_name+".STATUS_SERVICEUNAVAILABLE", NewInteger(int64(http.StatusServiceUnavailable)))
	SetGlobalObj(http_name+".STATUS_GATEWAYTIMEOUT", NewInteger(int64(http.StatusGatewayTimeout)))
	SetGlobalObj(http_name+".STATUS_HTTPVERSIONNOTSUPPORTED", NewInteger(int64(http.StatusHTTPVersionNotSupported)))
	SetGlobalObj(http_name+".STATUS_VARIANTALSONEGOTIATES", NewInteger(int64(http.StatusVariantAlsoNegotiates)))
	SetGlobalObj(http_name+".STATUS_INSUFFICIENTSTORAGE", NewInteger(int64(http.StatusInsufficientStorage)))
	SetGlobalObj(http_name+".STATUS_LOOPDETECTED", NewInteger(int64(http.StatusLoopDetected)))
	SetGlobalObj(http_name+".STATUS_NOTEXTENDED", NewInteger(int64(http.StatusNotExtended)))
	SetGlobalObj(http_name+".STATUS_NETWORKAUTHENTICATIONREQUIRED", NewInteger(int64(http.StatusNetworkAuthenticationRequired)))

	return ret
}

func (h *HttpObj) Inspect() string  { return http_name }
func (h *HttpObj) Type() ObjectType { return HTTP_OBJ }

func (h *HttpObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "get":
		return h.Get(line, args...)
	case "head":
		return h.Head(line, args...)
	case "post":
		return h.Post(line, args...)
	case "postForm":
		return h.PostForm(line, args...)
	case "newRequest":
		return h.NewRequest(line, args...)
	case "listenAndServe":
		return h.ListenAndServe(line, scope, args...)
	case "handle":
		return h.Handle(line, scope, args...)
	case "handleFunc":
		return h.HandleFunc(line, scope, args...)
	case "newServer":
		return h.NewServer(line, args...)
	case "redirect":
		return h.Redirect(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, h.Type()))
}

func (h *HttpObj) Get(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	url, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "get", "*String", args[0].Type()))
	}

	response, err := http.Get(url.String)
	if err != nil {
		return NewNil(err.Error())
	}
	return &HttpResponse{Response: response}
}

func (h *HttpObj) Head(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	url, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "head", "*String", args[0].Type()))
	}

	response, err := http.Head(url.String)
	if err != nil {
		return NewNil(err.Error())
	}
	return &HttpResponse{Response: response}
}

func (h *HttpObj) Post(line string, args ...Object) Object {
	if len(args) != 2 && len(args) != 3 {
		panic(NewError(line, ARGUMENTERROR, "2|3", len(args)))
	}

	urlStr, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "post", "*String", args[0].Type()))
	}

	contentType, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "post", "*String", args[1].Type()))
	}

	var response *http.Response
	var err error

	if len(args) == 2 {
		response, err = http.Post(urlStr.String, contentType.String, nil)
	} else {
		body, ok := args[2].(*String)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "third", "post", "*String", args[2].Type()))
		}
		response, err = http.Post(urlStr.String, contentType.String, strings.NewReader(body.String))
	}

	if err != nil {
		return NewNil(err.Error())
	}

	return &HttpResponse{Response: response}
}

func (h *HttpObj) PostForm(line string, args ...Object) Object {
	//	if len(args) != 1 {
	//		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	//	}
	//
	//	url, ok := args[0].(*String)
	//	if !ok {
	//		panic(NewError(line, PARAMTYPEERROR, "first", "postForm", "*String", args[0].Type()))
	//	}
	//
	//	response, err := http.Post(url.String)
	//	if err != nil {
	//		//return FALSE
	//		panic(NewError(line, GENERICERROR, err.Error()))
	//	}
	//	return &HttpResponse{Response:response}
	return NIL
}

func (h *HttpObj) NewRequest(line string, args ...Object) Object {
	if len(args) != 2 && len(args) != 3 {
		panic(NewError(line, ARGUMENTERROR, "2|3", len(args)))
	}

	method, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "newRequest", "*String", args[0].Type()))
	}

	urlStr, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "newRequest", "*String", args[1].Type()))
	}

	var request *http.Request
	var err error

	if len(args) == 2 {
		request, err = http.NewRequest(method.String, urlStr.String, nil)
	} else {
		body, ok := args[2].(*String)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "second", "newRequest", "*String", args[2].Type()))
		}
		request, err = http.NewRequest(method.String, urlStr.String, strings.NewReader(body.String))
	}

	if err != nil {
		return NewNil(err.Error())
	}

	return &HttpRequest{Request: request}
}

func (h *HttpObj) ListenAndServe(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 && len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	addr, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "listenAndServe", "*String", args[0].Type()))
	}

	if len(args) == 1 {
		err := http.ListenAndServe(addr.String, nil)
		if err != nil {
			return NewFalseObj(err.Error())
		}
	} else {
		block, ok := args[1].(*Function)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "second", "listenAndServe", "*Function", args[1].Type()))
		}

		paramCount := len(block.Literal.Parameters)
		if paramCount != 2 {
			panic(NewError(line, FUNCCALLBACKERROR, 2, paramCount))
		}

		err := http.ListenAndServe(addr.String, &customHTTPHandler{Scope: scope, F: block})
		if err != nil {
			return NewFalseObj(err.Error())
		}
	}

	return TRUE
}

func (h *HttpObj) Handle(line string, scope *Scope, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	pattern, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "handle", "*String", args[0].Type()))
	}

	block, ok := args[1].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "handle", "*Function", args[1].Type()))
	}

	paramCount := len(block.Literal.Parameters)
	if paramCount != 2 {
		panic(NewError(line, FUNCCALLBACKERROR, 2, paramCount))
	}

	http.Handle(pattern.String, &customHTTPHandler{Scope: scope, F: block})

	return NIL
}

func (h *HttpObj) HandleFunc(line string, scope *Scope, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	pattern, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "handleFunc", "*String", args[0].Type()))
	}

	block, ok := args[1].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "handleFunc", "*Function", args[1].Type()))
	}

	paramCount := len(block.Literal.Parameters)
	if paramCount != 2 {
		panic(NewError(line, FUNCCALLBACKERROR, 2, paramCount))
	}

	http.HandleFunc(pattern.String, func(w http.ResponseWriter, r *http.Request) {
		ServeHTTP(scope, block, w, r)
	})

	return NIL
}

func ServeHTTP(scope *Scope, f *Function, w http.ResponseWriter, r *http.Request) {
	s := NewScope(scope)

	//Save the two variables to `Scope`, so `Eval` can use them
	s.Set(f.Literal.Parameters[0].(*ast.Identifier).Value, &HttpResponseWriter{Writer: w})
	s.Set(f.Literal.Parameters[1].(*ast.Identifier).Value, &HttpRequest{Request: r})
	Eval(f.Literal.Body, s)
}

func (h *HttpObj) NewServer(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	addr, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "newServer", "*String", args[0].Type()))
	}

	return &HttpServer{Server: &http.Server{Addr: addr.String}}
}

func (h *HttpObj) Redirect(line string, args ...Object) Object {
	if len(args) != 4 {
		panic(NewError(line, ARGUMENTERROR, "4", len(args)))
	}

	rw, ok := args[0].(*HttpResponseWriter)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "redirect", "*HttpResponseWriter", args[0].Type()))
	}

	req, ok := args[1].(*HttpRequest)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "redirect", "*HttpRequest", args[1].Type()))
	}

	urlStr, ok := args[2].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "third", "redirect", "*String", args[2].Type()))
	}

	code, ok := args[3].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "fourth", "redirect", "*Integer", args[3].Type()))
	}

	http.Redirect(rw.Writer, req.Request, urlStr.String, int(code.Int64))
	return NIL
}

type customHTTPHandler struct {
	Scope *Scope
	F     *Function
}

func (c *customHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := NewScope(c.Scope)

	//Save the two variables to `Scope`, so `Eval` can use them
	s.Set(c.F.Literal.Parameters[0].(*ast.Identifier).Value, &HttpResponseWriter{Writer: w})
	s.Set(c.F.Literal.Parameters[1].(*ast.Identifier).Value, &HttpRequest{Request: r})
	Eval(c.F.Literal.Body, s)
}

//HTTP Client object
type HttpClient struct {
	Client *http.Client
}

func (h *HttpClient) Inspect() string  { return "httpclient" }
func (h *HttpClient) Type() ObjectType { return HTTPCLIENT_OBJ }
func (h *HttpClient) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "do":
		return h.Do(line, args...)
	case "get":
		return h.Get(line, args...)
	case "head":
		return h.Head(line, args...)
	case "post":
		return h.Post(line, args...)
	case "postForm":
		return h.PostForm(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, h.Type()))
	}
}

func (h *HttpClient) Do(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	req, ok := args[0].(*HttpRequest)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "do", "*HttpRequest", args[0].Type()))
	}

	response, err := h.Client.Do(req.Request)
	if err != nil {
		return NewNil(err.Error())
	}
	return &HttpResponse{Response: response}
}

func (h *HttpClient) Get(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	url, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "get", "*String", args[0].Type()))
	}

	response, err := h.Client.Get(url.String)
	if err != nil {
		return NewNil(err.Error())
	}
	return &HttpResponse{Response: response}
}

func (h *HttpClient) Head(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	url, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "url", "*String", args[0].Type()))
	}

	response, err := h.Client.Head(url.String)
	if err != nil {
		return NewNil(err.Error())
	}
	return &HttpResponse{Response: response}
}

func (h *HttpClient) Post(line string, args ...Object) Object {
	if len(args) != 2 && len(args) != 3 {
		panic(NewError(line, ARGUMENTERROR, "2|3", len(args)))
	}

	urlStr, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "post", "*String", args[0].Type()))
	}

	contentType, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "post", "*String", args[1].Type()))
	}

	var response *http.Response
	var err error

	if len(args) == 2 {
		response, err = h.Client.Post(urlStr.String, contentType.String, nil)
	} else {
		body, ok := args[2].(*String)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "third", "post", "*String", args[2].Type()))
		}
		response, err = h.Client.Post(urlStr.String, contentType.String, strings.NewReader(body.String))
	}

	if err != nil {
		return NewNil(err.Error())
	}

	return &HttpResponse{Response: response}
}

func (h *HttpClient) PostForm(line string, args ...Object) Object {
	//	if len(args) != 1 {
	//		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	//	}
	//
	//	url, ok := args[0].(*String)
	//	if !ok {
	//		panic(NewError(line, INPUTERROR, args[0].Type(), "postForm"))
	//	}
	//
	//	response, err := h.Client.Post(url.String)
	//	if err != nil {
	//		//return FALSE
	//		panic(NewError(line, GENERICERROR, err.Error()))
	//	}
	//	return &HttpResponse{Response:response}
	return NIL
}

//HTTP Response object
type HttpResponse struct {
	Response *http.Response
}

func (h *HttpResponse) Inspect() string  { return "httpresponse" }
func (h *HttpResponse) Type() ObjectType { return HTTPRESPONSE_OBJ }
func (h *HttpResponse) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "closeBody":
		return h.CloseBody(line, args...)
	case "readAll":
		return h.ReadAll(line, args...)
	case "header":
		return h.Header(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, h.Type()))
	}
}

func (h *HttpResponse) CloseBody(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	h.Response.Body.Close()
	return NIL
}

func (h *HttpResponse) ReadAll(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	b, err := ioutil.ReadAll(h.Response.Body)
	if err != nil {
		return NewNil(err.Error())
	}
	return NewString(string(b))
}

func (h *HttpResponse) Header(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return &HttpHeader{Header: h.Response.Header}
}

//HTTP Request object
type HttpRequest struct {
	Request *http.Request
}

func (h *HttpRequest) Inspect() string  { return "httprequest" }
func (h *HttpRequest) Type() ObjectType { return HTTPREQUEST_OBJ }
func (h *HttpRequest) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "header":
		return h.Header(line, args...)
	case "write":
		return h.Write(line, args...)
	case "formValue":
		return h.FormValue(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, h.Type()))
	}
	return NIL
}

func (h *HttpRequest) Header(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return &HttpHeader{Header: h.Request.Header}
}

func (h *HttpRequest) Write(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	content, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "write", "*String", args[0].Type()))
	}

	buf := bytes.NewBuffer([]byte(content.String))
	err := h.Request.Write(buf)
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

func (h *HttpRequest) FormValue(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	key, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "formValue", "*String", args[0].Type()))
	}

	return NewString(h.Request.FormValue(key.String))
}

//HTTP ResponseWriter object
type HttpResponseWriter struct {
	Writer http.ResponseWriter
}

func (h *HttpResponseWriter) IOWriter() io.Writer { return h.Writer }
func (h *HttpResponseWriter) Inspect() string     { return "http responsewriter" }
func (h *HttpResponseWriter) Type() ObjectType    { return HTTPRESPONSEWRITER_OBJ }
func (h *HttpResponseWriter) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "write":
		return h.Write(line, args...)
	case "writeHeader":
		return h.WriteHeader(line, args...)
	case "header":
		return h.Header(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, h.Type()))
	}
}

func (h *HttpResponseWriter) Write(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	content, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "write", "*String", args[0].Type()))
	}

	i, err := h.Writer.Write([]byte(content.String))
	if err != nil {
		return NewNil(err.Error())
	}
	return NewInteger(int64(i))
}

func (h *HttpResponseWriter) WriteHeader(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	statusCode, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "writeHeader", "*Integer", args[0].Type()))
	}

	h.Writer.WriteHeader(int(statusCode.Int64))

	return NIL
}

func (h *HttpResponseWriter) Header(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return &HttpHeader{Header: h.Writer.Header()}
}

//HTTP Server object
type HttpServer struct {
	Server *http.Server
}

func (h *HttpServer) Inspect() string  { return "httpserver" }
func (h *HttpServer) Type() ObjectType { return HTTPSERVER_OBJ }
func (h *HttpServer) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "setReadTimeout":
		return h.SetReadTimeout(line, args...)
	case "setWriteTimeout":
		return h.SetWriteTimeout(line, args...)
	case "setMaxHeaderBytes":
		return h.SetMaxHeaderBytes(line, args...)
	case "listenAndServe":
		return h.ListenAndServe(line, args...)
		//	case "close":
		//		return h.Close(line, args...)
	case "setKeepAlivesEnabled":
		return h.SetKeepAlivesEnabled(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, h.Type()))
	}
}

func (h *HttpServer) SetReadTimeout(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	duration, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setReadTimeout", "*Integer", args[0].Type()))
	}

	h.Server.ReadTimeout = time.Duration(duration.Int64) * time.Second

	return NIL
}

func (h *HttpServer) SetWriteTimeout(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	duration, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setWriteTimeout", "*Integer", args[0].Type()))
	}

	h.Server.WriteTimeout = time.Duration(duration.Int64) * time.Second

	return NIL
}

func (h *HttpServer) SetMaxHeaderBytes(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	headerBytes, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setMaxHeaderBytes", "*Integer", args[0].Type()))
	}

	h.Server.MaxHeaderBytes = int(headerBytes.Int64)

	return NIL
}

func (h *HttpServer) ListenAndServe(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	err := h.Server.ListenAndServe()
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

//func (h *HttpServer) Close(line string, args ...Object) Object {
//	if len(args) != 0 {
//		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
//	}
//
//	err := h.Server.Close()
//	if err != nil {
//		//return FALSE
//		panic(NewError(line, GENERICERROR, err.Error()))
//	}
//
//	return NIL
//}

func (h *HttpServer) SetKeepAlivesEnabled(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	b, ok := args[0].(*Boolean)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setKeepAlivesEnabled", "*Boolean", args[0].Type()))
	}

	h.Server.SetKeepAlivesEnabled(b.Bool)
	return NIL
}

//HTTP Header object
type HttpHeader struct {
	Header http.Header
}

func (h *HttpHeader) Inspect() string  { return "httpheader" }
func (h *HttpHeader) Type() ObjectType { return HTTPHEADER_OBJ }
func (h *HttpHeader) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "add":
		return h.Add(line, args...)
	case "del":
		return h.Del(line, args...)
	case "get":
		return h.Get(line, args...)
	case "setHeader":
		return h.Set(line, args...)
	case "write":
		return h.Write(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, h.Type()))
	}
	return NIL
}

func (h *HttpHeader) Add(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	key, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "add", "*String", args[0].Type()))
	}

	value, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "add", "*String", args[1].Type()))
	}

	h.Header.Add(key.String, value.String)
	return NIL
}

func (h *HttpHeader) Del(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	key, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "del", "*String", args[0].Type()))
	}

	h.Header.Del(key.String)
	return NIL
}

func (h *HttpHeader) Get(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	key, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "get", "*String", args[0].Type()))
	}

	str := h.Header.Get(key.String)
	return NewString(str)
}

func (h *HttpHeader) Set(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	key, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "set", "*String", args[0].Type()))
	}

	value, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "set", "*String", args[1].Type()))
	}

	h.Header.Set(key.String, value.String)
	return NIL
}

func (h *HttpHeader) Write(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	content, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "write", "*String", args[0].Type()))
	}

	buf := bytes.NewBuffer([]byte(content.String))
	err := h.Header.Write(buf)
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}
