package abstractions

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"strings"
	"time"

	u "net/url"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
	stduritemplate "github.com/std-uritemplate/std-uritemplate/go/v2"
)

var _ RequestBuilder = (*NativeRequestBuilder)(nil)

type Builder[T any] interface {
	Build() (*T, error)
}

type RequestBuilder interface {
	Builder[Request]
}

// NativeRequestBuilder is the base class for all request builders.
type NativeRequestBuilder struct {
	// pathParameters The path parameters for the request
	pathParameters map[string]any
	// requestAdapter The request adapter to use to execute the requests.
	requestAdapter RequestAdapter
	// urlTemplate The url template to use to build the URL for the current request builder
	urlTemplate string
	// method The HTTP method of the request.
	method HttpMethod
	// headers The Request Headers.
	headers *RequestHeaders
	// queryParameters The Query Parameters of the request.
	queryParameters map[string]any
	// content The Request content.
	content []byte
	// options The Request options.
	options map[string]RequestOption
	// context The context of the request.
	context context.Context
}

// NewBaseRequestBuilder creates a new BaseRequestBuilder instance.
func NewBaseRequestBuilder(requestAdapter RequestAdapter, urlTemplate string, pathParameters map[string]any) *NativeRequestBuilder {
	return &NativeRequestBuilder{
		requestAdapter: requestAdapter,
		urlTemplate:    urlTemplate,
		pathParameters: maps.Clone(pathParameters),
	}
}

func (rB *NativeRequestBuilder) WithHeaders(headers *RequestHeaders) *NativeRequestBuilder {
	rB.headers = headers

	return rB
}

func (rB *NativeRequestBuilder) WithOptions(options ...RequestOption) *NativeRequestBuilder {
	if len(options) == 0 {
		return rB
	}
	if len(rB.options) == 0 {
		rB.options = make(map[string]RequestOption, len(options))
	}
	for _, option := range options {
		rB.options[option.GetKey().Key] = option
	}
	return rB
}

func (rB *NativeRequestBuilder) WithQueryParameters(parameters any) *NativeRequestBuilder {
	if parameters == nil || rB == nil {
		return rB
	}
	valOfP := reflect.ValueOf(parameters)
	fields := reflect.TypeOf(parameters)
	numOfFields := fields.NumField()
	for i := 0; i < numOfFields; i++ {
		field := fields.Field(i)
		fieldName := field.Name
		fieldValue := valOfP.Field(i)
		tagValue := field.Tag.Get("uriparametername")
		if tagValue != "" {
			fieldName = tagValue
		}
		value := rB.sanitizeValue(fieldValue.Interface())
		valueOfValue := reflect.ValueOf(value)
		if valueOfValue.IsNil() {
			continue
		}
		strArr, ok := value.([]string)
		if ok && len(strArr) > 0 {
			tmp := make([]any, len(strArr))
			for i, v := range strArr {
				tmp[i] = v
			}
			rB.queryParameters[fieldName] = tmp
		}
		if arr, ok := value.([]any); ok && len(arr) > 0 {
			rB.queryParameters[fieldName] = arr
		}
		normalizedValue := rB.normalizeParameters(valueOfValue, value, true)
		if normalizedValue != nil {
			rB.queryParameters[fieldName] = normalizedValue
		}
	}
	return rB
}

func (rB *NativeRequestBuilder) WithMethod(method HttpMethod) *NativeRequestBuilder {

	rB.method = method

	return rB
}

func (rB *NativeRequestBuilder) WithContent(content []byte, contentType string) *NativeRequestBuilder {
	rB.content = content
	if rB.headers == nil {
		rB.headers = NewRequestHeaders()
	}
	rB.headers.Add(contentTypeHeader, contentType)

	return rB
}

func (rB *NativeRequestBuilder) WithContext(ctx context.Context) *NativeRequestBuilder {
	rB.context = ctx

	return rB
}

// TODO: shouldn't be a method
func (rB *NativeRequestBuilder) sanitizeValue(value any) any {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case *time.Time:
		return v.Format(time.RFC3339)
	case time.Time:
		return v.Format(time.RFC3339)
	case []*time.Time:
		return castItem(v, func(t *time.Time) string {
			return t.Format(time.RFC3339)
		})
	case []time.Time:
		return castItem(v, func(t time.Time) string {
			return t.Format(time.RFC3339)
		})
	case *s.ISODuration:
		return v.String()
	case s.ISODuration:
		return v.String()
	case []*s.ISODuration:
		return castItem(v, func(v *s.ISODuration) string {
			return v.String()
		})
	case []s.ISODuration:
		return castItem(v, func(v s.ISODuration) string {
			return v.String()
		})
	case *s.TimeOnly:
		return v.String()
	case s.TimeOnly:
		return v.String()
	case []*s.TimeOnly:
		return castItem(v, func(v *s.TimeOnly) string {
			return v.String()
		})
	case []s.TimeOnly:
		return castItem(v, func(v s.TimeOnly) string {
			return v.String()
		})
	case *s.DateOnly:
		return v.String()
	case s.DateOnly:
		return v.String()
	case []*s.DateOnly:
		return castItem(v, func(v *s.DateOnly) string {
			return v.String()
		})
	case []s.DateOnly:
		return castItem(v, func(v s.DateOnly) string {
			return v.String()
		})
	}

	return value
}

// TODO: shouldn't be a method
// Normalize different types to values that can be rendered in an URL:
// enum -> string (name)
// []enum -> []string (containing names)
// []non_interface -> []any (like []int64 -> []any)
func (rB *NativeRequestBuilder) normalizeParameters(valueOfValue reflect.Value, value any, returnNilIfNotNormalizable bool) any {
	if valueOfValue.Kind() == reflect.Slice && valueOfValue.Len() > 0 {
		//type assertions to "enums" don't work if you don't know the enum type in advance, we need to use reflection
		enumArr := valueOfValue.Slice(0, valueOfValue.Len())
		if _, ok := enumArr.Index(0).Interface().(kiotaEnum); ok {
			// testing the first value is an enum to avoid iterating over the whole array if it's not
			strRepresentations := make([]string, valueOfValue.Len())
			for i := range strRepresentations {
				strRepresentations[i] = enumArr.Index(i).Interface().(kiotaEnum).String()
			}
			return strRepresentations
		} else {
			anySlice := make([]any, valueOfValue.Len())
			for i := range anySlice {
				anySlice[i] = enumArr.Index(i).Interface()
			}
			return anySlice
		}
	} else if enum, ok := value.(kiotaEnum); ok {
		return enum.String()
	}

	if returnNilIfNotNormalizable {
		return nil
	} else {
		return value
	}
}

// getURI returns the URI of the request.
func (rB *NativeRequestBuilder) getURI() (*u.URL, error) {
	if rB.urlTemplate == "" {
		return nil, errors.New("uri cannot be empty")
	} else if rB.pathParameters == nil {
		return nil, errors.New("uri template parameters cannot be nil")
	} else if rB.queryParameters == nil {
		return nil, errors.New("uri query parameters cannot be nil")
	} else if rawURL := rB.pathParameters[raw_url_key]; rawURL != nil {
		rawURLString, ok := rawURL.(string)
		if !ok {
			return nil, fmt.Errorf("pathParameter %s is not %T", raw_url_key, "")
		}
		uri, err := u.Parse(rawURLString)
		if err != nil {
			return nil, err
		}
		return uri, nil
	}
	_, baseurlExists := rB.pathParameters["baseurl"]
	if !baseurlExists && strings.Contains(strings.ToLower(rB.urlTemplate), "{+baseurl}") {
		return nil, errors.New("pathParameters must contain a value for \"baseurl\" for the url to be built")
	}

	substitutions := make(map[string]any)
	for key, value := range rB.pathParameters {
		substitutions[key] = rB.normalizeParameters(reflect.ValueOf(value), rB.sanitizeValue(value), false)
	}
	for key, value := range rB.queryParameters {
		substitutions[key] = rB.sanitizeValue(value)
	}
	url, err := stduritemplate.Expand(rB.urlTemplate, substitutions)
	if err != nil {
		return nil, err
	}
	uri, err := u.Parse(url)
	return uri, err
}

func (rB *NativeRequestBuilder) Build() (*Request, error) {
	uri, err := rB.getURI()
	if err != nil {
		return nil, err
	}

	request := &Request{
		Method:          rB.method,
		uri:             uri,
		Headers:         rB.headers,
		QueryParameters: rB.queryParameters,
		Content:         rB.content,
		PathParameters:  rB.pathParameters,
		UrlTemplate:     rB.urlTemplate,
		options:         rB.options,
	}

	return request, nil
}

func ConfigureRequest[T any](request *NativeRequestBuilder, config *RequestConfiguration[T]) {
	if request == nil {
		return
	}
	if config == nil {
		return
	}
	if config.QueryParameters != nil {
		request.WithQueryParameters(*(config.QueryParameters))
	}
	request.WithHeaders(config.Headers)
	request.WithOptions(config.Options...)
}
