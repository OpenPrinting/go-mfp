// MFP - Miulti-Function Printers and scanners toolkit
// IPP - Internet Printing Protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Get-Printer-Attributes request

package ipp

import (
	"errors"

	"github.com/OpenPrinting/go-mfp/proto/ipp/iana"
	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/go-mfp/util/optional"
	"github.com/OpenPrinting/goipp"
)

// Standard attribute groups for Get-Printer-Attributes.
const (
	// GetPrinterAttributesAll requests all printer attributes,
	// except the media-col-database.
	GetPrinterAttributesAll = "all"

	// GetPrinterAttributesJobTemplate requests the Job Template
	// Attributes.
	GetPrinterAttributesJobTemplate = "job-template"

	// GetPrinterAttributesPrinterDescription requests the
	// Printer Description Attributes.
	GetPrinterAttributesPrinterDescription = "printer-description"

	// GetPrinterAttributesMediaColDatabase requests the collection
	// of supported media types.
	//
	// Note, the "media-col-database" is not returned by the
	// printer unless explicitly requested, even if "all" attributes
	// are requested.
	GetPrinterAttributesMediaColDatabase = "media-col-database"
)

// GetPrinterAttributesRequest operation (0x000b) returns
// the requested printer attributes.
type GetPrinterAttributesRequest struct {
	ObjectRawAttrs
	RequestHeader
	OperationGroup

	// Operation attributes
	PrinterURI          string               `ipp:"printer-uri"`
	RequestedAttributes []string             `ipp:"requested-attributes"`
	DocumentFormat      optional.Val[string] `ipp:"document-format"`
}

// GetPrinterAttributesResponse is the CUPS-Get-Default Response.
type GetPrinterAttributesResponse struct {
	ObjectRawAttrs
	ResponseHeader
	OperationGroup

	// Names of unsupported attributes
	UnsupportedAttributes []string

	// Returned printer attributes
	Printer *PrinterAttributes
}

// GetOp returns GetPrinterAttributesRequest IPP Operation code.
func (rq *GetPrinterAttributesRequest) GetOp() goipp.Op {
	return goipp.OpGetPrinterAttributes
}

// Encode encodes GetPrinterAttributesRequest into the goipp.Message.
func (rq *GetPrinterAttributesRequest) Encode() *goipp.Message {
	enc := ippEncoder{}

	groups := goipp.Groups{
		{
			Tag:   goipp.TagOperationGroup,
			Attrs: enc.Encode(rq),
		},
	}

	msg := goipp.NewMessageWithGroups(rq.Version, goipp.Code(rq.GetOp()),
		rq.RequestID, groups)

	return msg
}

// Decode decodes GetPrinterAttributesRequest from goipp.Message.
func (rq *GetPrinterAttributesRequest) Decode(
	msg *goipp.Message, opt *DecoderOptions) error {

	rq.Version = msg.Version
	rq.RequestID = msg.RequestID

	dec := NewDecoder(opt)
	defer dec.Free()

	err := dec.Decode(rq, msg.Operation)
	if err != nil {
		return err
	}

	return nil
}

// Encode encodes GetPrinterAttributesResponse into goipp.Message.
func (rsp *GetPrinterAttributesResponse) Encode() *goipp.Message {
	enc := ippEncoder{}

	var attrs goipp.Attributes
	if rsp.Printer != nil {
		attrs = enc.Encode(rsp.Printer)
	}

	return rsp.EncodeRaw(attrs)
}

// EncodeRaw is like [GetPrinterAttributesResponse.Encode],
// but it accepts printer attributes as parameter and ignores
// the [GetPrinterAttributesResponse.Printer] field.
//
// This function is convenient when there is a need to constrict
// Get-Printer-Attributes response from the raw set of printer
// attributes.
func (rsp *GetPrinterAttributesResponse) EncodeRaw(
	rawPrinterAttrs goipp.Attributes) *goipp.Message {

	enc := ippEncoder{}

	groups := goipp.Groups{
		{
			Tag:   goipp.TagOperationGroup,
			Attrs: enc.Encode(rsp),
		},
	}

	if len(rsp.UnsupportedAttributes) > 0 {
		names := make(goipp.Values, 0, len(rsp.UnsupportedAttributes))
		for _, name := range rsp.UnsupportedAttributes {
			names.Add(goipp.TagKeyword, goipp.String(name))
		}

		attr := goipp.Attribute{
			Name:   "requested-attributes",
			Values: names,
		}

		groups.Add(goipp.Group{
			Tag:   goipp.TagUnsupportedGroup,
			Attrs: goipp.Attributes{attr},
		})
	}

	if rawPrinterAttrs != nil {
		groups.Add(goipp.Group{
			Tag:   goipp.TagPrinterGroup,
			Attrs: rawPrinterAttrs,
		})
	}

	msg := goipp.NewMessageWithGroups(rsp.Version, goipp.Code(rsp.Status),
		rsp.RequestID, groups)

	return msg
}

// printerAttrGroups maps the standard attribute-group keywords
// ("all", "printer-description", "job-template") to the set of
// individual attribute names that belong to each group, for
// Get-Printer-Attributes requests (used by both Print Services
// (RFC8011) and Scan Services (PWG5100.17)).
//
// Other Get-XXX-Attributes operations (e.g. Get-Job-Attributes)
// will define their own analogous group maps.
var printerAttrGroups = buildPrinterAttrGroups()

// buildPrinterAttrGroups constructs the printer attribute-group
// expansion map from the IANA registration database.
func buildPrinterAttrGroups() map[string]generic.Set[string] {
	all := generic.NewSet[string]()
	for name := range iana.PrinterDescription {
		all.Add(name)
	}
	for name := range iana.PrinterStatus {
		all.Add(name)
	}
	all.Del("media-col-database")

	jobTemplate := generic.NewSet[string]()
	all.ForEach(func(name string) {
		if iana.PrinterDescription[name+"-default"] != nil {
			jobTemplate.Add(name + "-default")
		}
		if iana.PrinterDescription[name+"-supported"] != nil {
			jobTemplate.Add(name + "-supported")
		}
	})

	printerDescription := all.Clone()
	jobTemplate.ForEach(func(name string) {
		printerDescription.Del(name)
	})

	return map[string]generic.Set[string]{
		"all":                 all,
		"printer-description": printerDescription,
		"job-template":        jobTemplate,
	}
}

// filterAttributes filters encoded against the requested groups/names
// defined by attrGroups. Returns the filtered subset (preserving the
// original order of encoded) and the list of requested names that were
// neither a known group nor a supported attribute.
func filterAttributes(
	requestedAttrs []string,
	encoded goipp.Attributes,
	attrGroups map[string]generic.Set[string],
) (filtered goipp.Attributes, unsupported []string) {

	supported := generic.NewSet[string]()
	for _, attr := range encoded {
		supported.Add(attr.Name)
	}

	filter := generic.NewSet[string]()
	seenUnsupported := generic.NewSet[string]()

	for _, name := range requestedAttrs {
		if group, ok := attrGroups[name]; ok {
			filter.Merge(group)
		} else if supported.Contains(name) {
			filter.Add(name)
		} else if seenUnsupported.TestAndAdd(name) {
			unsupported = append(unsupported, name)
		}
	}

	for _, attr := range encoded {
		if filter.Contains(attr.Name) {
			filtered = append(filtered, attr)
		}
	}

	return
}

// Apply applies the request's requested-attributes filter to attrs and
// returns the encoded response message. If useRawAttrs is true, the source
// attrs are taken from attrs.RawAttrs().All() (lossless wire form);
// otherwise they are produced by encoding attrs.
func (rq *GetPrinterAttributesRequest) Apply(
	attrs *PrinterAttributes,
	useRawAttrs bool,
) *goipp.Message {

	var encoded goipp.Attributes
	if useRawAttrs {
		encoded = attrs.RawAttrs().All()
	} else {
		enc := ippEncoder{}
		encoded = enc.Encode(attrs)
	}

	filtered, unsupported := filterAttributes(
		rq.RequestedAttributes, encoded, printerAttrGroups)

	status := goipp.StatusOk
	if len(unsupported) > 0 {
		status = goipp.StatusOkIgnoredOrSubstituted
	}

	rsp := &GetPrinterAttributesResponse{
		ResponseHeader:        rq.ResponseHeader(status),
		UnsupportedAttributes: unsupported,
	}

	return rsp.EncodeRaw(filtered)
}

// Decode decodes GetPrinterAttributesResponse from goipp.Message.
func (rsp *GetPrinterAttributesResponse) Decode(
	msg *goipp.Message, opt *DecoderOptions) error {

	rsp.Version = msg.Version
	rsp.RequestID = msg.RequestID
	rsp.Status = goipp.Status(msg.Code)

	dec := NewDecoder(opt)
	defer dec.Free()

	err := dec.Decode(rsp, msg.Operation)
	if err != nil {
		return err
	}

	if len(msg.Printer) == 0 {
		err = errors.New("missed printer attributes in response")
		return err
	}

	rsp.Printer, err = DecodePrinterAttributes(msg.Printer, opt)
	if err != nil {
		return err
	}

	return nil
}
