// MFP - Multi-Function Printers and scanners toolkit
// WS-Scan core protocol
//
// Copyright (C) 2024 and up by Yogesh Singla (yogeshsingla481@gmail.com)
// See LICENSE for license terms and conditions
//
// ValWithOptions: reusable type for elements with
// a value and optional boolean attributes (MustHonor, Override, UsedDefault)

package wsscan

import (
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/go-mfp/util/xmldoc"
)

// ValWithOptions holds a value and optional boolean attributes.
// This is a generic element for patterns like:
// <wscn:Element
//
//	wscn:MustHonor="true"
//	wscn:Override="false"
//	wscn:UsedDefault="true">
//	    value
//
// </wscn:Element>
// The type parameter T allows the value to be any type (string, int, etc.)
type ValWithOptions[T any] struct {
	Val         T
	MustHonor   optional.Val[BooleanElement]
	Override    optional.Val[BooleanElement]
	UsedDefault optional.Val[BooleanElement]
}

// WithOptions is a combination of [WithOptionsGetter]
// and [WithOptionsSetter] interfaces.
type WithOptions interface {
	WithOptionsGetter
	WithOptionsSetter
}

// WithOptionsGetter is the common interface, implemented by
// all variants of ValWithOptions[T], regardless of T.
//
// It provides an uniform get access to the underlying value
// and options.
type WithOptionsGetter interface {
	GetValue() any
	GetMustHonor() optional.Val[BooleanElement]
	GetOverride() optional.Val[BooleanElement]
	GetUsedDefault() optional.Val[BooleanElement]
}

// WithOptionsSetter is the common interface, implemented by
// all variants of ValWithOptions[T], regardless of T.
//
// It provides an uniform set access to the underlying value
// and options.
//
// It requires pointer receiver.
type WithOptionsSetter interface {
	SetValue(any) bool
	SetMustHonor(optional.Val[BooleanElement])
	SetOverride(optional.Val[BooleanElement])
	SetUsedDefault(optional.Val[BooleanElement])
}

// GetValue implements [WithOptions] interface for getting the
// underlying value without options.
func (t ValWithOptions[T]) GetValue() any {
	return t.Val
}

// SetValue implements [WithOptions] interface for setting the
// underlying value.
func (t *ValWithOptions[T]) SetValue(v any) bool {
	var ok bool
	t.Val, ok = v.(T)
	return ok
}

// GetMustHonor implements [WithOptions] interface for getting the
// MustHonor option.
func (t ValWithOptions[T]) GetMustHonor() optional.Val[BooleanElement] {
	return t.MustHonor
}

// SetMustHonor implements [WithOptions] interface for setting the
// MustHonor option.
func (t *ValWithOptions[T]) SetMustHonor(opt optional.Val[BooleanElement]) {
	t.MustHonor = opt
}

// GetOverride implements [WithOptions] interface for getting the
// Override option.
func (t ValWithOptions[T]) GetOverride() optional.Val[BooleanElement] {
	return t.Override
}

// SetOverride implements [WithOptions] interface for setting the
// Override option.
func (t *ValWithOptions[T]) SetOverride(opt optional.Val[BooleanElement]) {
	t.Override = opt
}

// GetUsedDefault implements [WithOptions] interface for getting the
// UsedDefault option.
func (t ValWithOptions[T]) GetUsedDefault() optional.Val[BooleanElement] {
	return t.UsedDefault
}

// SetUsedDefault implements [WithOptions] interface for setting the
// UsedDefault option.
func (t *ValWithOptions[T]) SetUsedDefault(opt optional.Val[BooleanElement]) {
	t.UsedDefault = opt
}

// HasOptions reports if value really has any options set.
// It implements the [Wrapper] interface.
func (t ValWithOptions[T]) HasOptions() bool {
	return t.MustHonor != nil || t.Override != nil || t.UsedDefault != nil
}

// Unwrap returns the underlying value, if t has no options, or the
// t's value itself otherwise.
//
// It implements the [Wrapper] interface.
func (t ValWithOptions[T]) Unwrap() any {
	if !t.HasOptions() {
		return t.Val
	}
	return t
}

// Wrap wraps the simple value into the Wrapper
// type and returns the new wrapped value.
//
// In case the provided value cannot be converted
// into the Wrapper's underlying type, this function
// returns nil.
func (t ValWithOptions[T]) Wrap(v any) any {
	val, ok := v.(T)
	if ok {
		return ValWithOptions[T]{Val: val}
	}
	return nil
}

// decodeValWithOptions fills the struct from an XML element.
// The decoder function converts the XML text to the desired type T.
func (t *ValWithOptions[T]) decodeValWithOptions(
	root xmldoc.Element,
	decoder func(string) (T, error),
) (ValWithOptions[T], error) {
	// Decode the text value using the provided decoder
	val, err := decoder(root.Text)
	if err != nil {
		return *t, err
	}
	t.Val = val

	// Decode MustHonor attribute
	if attr, found := root.AttrByName(NsWSCN + ":MustHonor"); found {
		boolVal := BooleanElement(attr.Value)
		if err := boolVal.Validate(); err != nil {
			return *t, err
		}
		t.MustHonor = optional.New(boolVal)
	}

	// Decode Override attribute
	if attr, found := root.AttrByName(NsWSCN + ":Override"); found {
		boolVal := BooleanElement(attr.Value)
		if err := boolVal.Validate(); err != nil {
			return *t, err
		}
		t.Override = optional.New(boolVal)
	}

	// Decode UsedDefault attribute
	if attr, found := root.AttrByName(NsWSCN + ":UsedDefault"); found {
		boolVal := BooleanElement(attr.Value)
		if err := boolVal.Validate(); err != nil {
			return *t, err
		}
		t.UsedDefault = optional.New(boolVal)
	}

	return *t, nil
}

// toXML creates an XML element from the struct.
// The encoder function converts the value of type T to a string.
func (t ValWithOptions[T]) toXML(
	name string,
	encoder func(T) string,
) xmldoc.Element {
	elm := xmldoc.Element{Name: name, Text: encoder(t.Val)}
	var attrs []xmldoc.Attr

	// Add MustHonor attribute if present
	if t.MustHonor != nil {
		attrs = append(attrs, xmldoc.Attr{
			Name:  NsWSCN + ":MustHonor",
			Value: string(optional.Get(t.MustHonor)),
		})
	}

	// Add Override attribute if present
	if t.Override != nil {
		attrs = append(attrs, xmldoc.Attr{
			Name:  NsWSCN + ":Override",
			Value: string(optional.Get(t.Override)),
		})
	}

	// Add UsedDefault attribute if present
	if t.UsedDefault != nil {
		attrs = append(attrs, xmldoc.Attr{
			Name:  NsWSCN + ":UsedDefault",
			Value: string(optional.Get(t.UsedDefault)),
		})
	}

	if len(attrs) > 0 {
		elm.Attrs = attrs
	}

	return elm
}
