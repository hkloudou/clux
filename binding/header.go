package binding

import (
	"net/textproto"
	"reflect"
)

type headerBinding struct{}

func (headerBinding) Name() string {
	return "header"
}

func (headerBinding) Bind(req RequestTransportData, obj interface{}) error {

	if err := mapHeader(obj, req.Head()); err != nil {
		return err
	}

	return validate(obj)
}

func mapHeader(ptr interface{}, h map[string][]string) error {
	return mappingByPtr(ptr, headerSource(h), "header")
}

type headerSource map[string][]string

var _ setter = headerSource(nil)

func (hs headerSource) TrySet(value reflect.Value, field reflect.StructField, tagValue string, opt setOptions) (bool, error) {
	return setByForm(value, field, hs, textproto.CanonicalMIMEHeaderKey(tagValue), opt)
}

type headerRawBinding struct{}

func (headerRawBinding) Name() string {
	return "header"
}

func (headerRawBinding) Bind(req RequestTransportData, obj interface{}) error {

	if err := mapHeaderRaw(obj, req.Head()); err != nil {
		return err
	}

	return validate(obj)
}

func mapHeaderRaw(ptr interface{}, h map[string][]string) error {
	return mappingByPtr(ptr, headerRawSource(h), "header")
}

type headerRawSource map[string][]string

var _ setter = headerSource(nil)

func (hs headerRawSource) TrySet(value reflect.Value, field reflect.StructField, tagValue string, opt setOptions) (bool, error) {
	return setByForm(value, field, hs, tagValue, opt)
}
