// MFP - Miulti-Function Printers and scanners toolkit
// XML mini library
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// XML element

package xmldoc

import (
	"os"

	"github.com/OpenPrinting/go-mfp/util/generic"
)

// Element represents a single decoded XML Element
type Element struct {
	Name     string    // Name of this element (ns:name)
	Text     string    // Element body
	Line     int       // Decode: document position, Encode: ignored
	Attrs    []Attr    // Element attributes
	Children []Element // All children
}

// Attr represents an XML element attribute
type Attr struct {
	Name  string // Attribute name
	Value string // Attribute value
}

// IsZero reports if Element equal to zero element
func (root Element) IsZero() bool {
	return root.Name == "" && root.Text == "" &&
		root.Attrs == nil && root.Children == nil
}

// Equal tests that two XML trees are equal, except difference
// in the [Element.Line].
func (root Element) Equal(root2 Element) bool {
	switch {
	case root.Name != root2.Name || root.Text != root2.Text:
		return false
	case len(root.Attrs) != len(root2.Attrs):
		return false
	case len(root.Children) != len(root2.Children):
		return false
	}

	for i := range root.Attrs {
		if !root.Attrs[i].Equal(root2.Attrs[i]) {
			return false
		}
	}

	for i := range root.Children {
		if !root.Children[i].Equal(root2.Children[i]) {
			return false
		}
	}

	return true
}

// Similar tests that two XML trees are similar, which means that
// they are equal, ignoring difference in order of children and
// attributes with the different names.
//
// Note, children and attributes with the same name still needs
// to be ordered equally.
func (root Element) Similar(root2 Element) bool {
	// Perform simple checks
	switch {
	case root.Name != root2.Name || root.Text != root2.Text:
		return false
	case len(root.Attrs) != len(root2.Attrs):
		return false
	case len(root.Children) != len(root2.Children):
		return false
	}

	// Compare attributes and children
	return root.attrsSimilar(root2) && root.childrenSimilar(root2)
}

// attrsSimilar tests that attributes of two XML elements are similar.
func (root Element) attrsSimilar(root2 Element) bool {
	attrs := make(map[string][]Attr, len(root.Attrs))
	attrs2 := make(map[string][]Attr, len(root2.Attrs))

	for i := range root.Attrs {
		attr := root.Attrs[i]
		list := attrs[attr.Name]
		attrs[attr.Name] = append(list, attr)

		attr = root2.Attrs[i]
		list = attrs2[attr.Name]
		attrs2[attr.Name] = append(list, attr)
	}

	for name, list := range attrs {
		list2 := attrs2[name]

		if len(list) != len(list2) {
			return false
		}

		for i, attr := range list {
			if attr.Value != list2[i].Value {
				return false
			}
		}
	}

	return true
}

// childrenSimilar tests that children of two XML elements are similar,
// recursively.
func (root Element) childrenSimilar(root2 Element) bool {
	children := make(map[string][]Element, len(root.Children))
	children2 := make(map[string][]Element, len(root2.Children))

	for i := range root.Children {
		child := root.Children[i]
		list := children[child.Name]
		children[child.Name] = append(list, child)

		child = root2.Children[i]
		list = children2[child.Name]
		children2[child.Name] = append(list, child)
	}

	for name, list := range children {
		list2 := children2[name]

		if len(list) != len(list2) {
			return false
		}

		for i, child := range list {
			if !child.Similar(list2[i]) {
				return false
			}
		}
	}

	return true
}

// Equal tests that two XML attrubutes are equal
func (attr Attr) Equal(attr2 Attr) bool {
	return attr == attr2
}

// Expand replaces ${var} or $var in the XML [Element.Text] and
// [Attr.Value] of the entire XML tree based on the mapping
// function and returns rewritten XML tree.
//
// It uses [os.Expand] function for strings substitution.
func (root Element) Expand(mapping func(string) string) Element {
	root.Text = os.Expand(root.Text, mapping)
	root.Attrs = generic.CopySlice(root.Attrs)
	root.Children = generic.CopySlice(root.Children)

	for i := range root.Attrs {
		attr := &root.Attrs[i]
		attr.Value = os.Expand(attr.Value, mapping)
	}

	for i := range root.Children {
		root.Children[i] = root.Children[i].Expand(mapping)
	}

	return root
}

// AttrByName searches element attributes by name.
//
// It returns (attr, true) if element was found or (Attr{}, false) if not.
func (root Element) AttrByName(name string) (Attr, bool) {
	for _, attr := range root.Attrs {
		if attr.Name == name {
			return attr, true
		}
	}

	return Attr{}, false
}

// AttrsMap returns map of the Element's attributes, indexed by name.
//
// If there are multiple attributes with the same name,
// only the first is returned.
func (root Element) AttrsMap() map[string]Attr {
	attrs := make(map[string]Attr)

	for _, attr := range root.Attrs {
		if _, found := attrs[attr.Name]; !found {
			attrs[attr.Name] = attr
		}
	}

	return attrs
}

// ChildByName searches child element by name.
//
// It returns (child, true) if element was found or (Element{}, false) if not.
func (root Element) ChildByName(name string) (Element, bool) {
	for _, elm := range root.Children {
		if elm.Name == name {
			return elm, true
		}
	}

	return Element{}, false
}

// ChildrenMap returns map of the Element's children, indexed by name.
//
// If there are multiple children with the same name,
// only the first is returned.
func (root Element) ChildrenMap() map[string]Element {
	children := make(map[string]Element)

	for _, elm := range root.Children {
		if _, found := children[elm.Name]; !found {
			children[elm.Name] = elm
		}
	}

	return children
}

// Lookup searches for children elements, multiple elements a time.
//
//	Use pattern:
//
//	  element1 = Lookup(Name: "Element1"}
//	  element2 = Lookup(Name: "Element2"}
//	  element3 = Lookup(Name: "Element3"}
//
//	  root.Lookup(&element1,&element2,&element3)
//	  if element1.Found{
//	      . . .
//	  }
//	  if element2.Found{
//	      . . .
//	  }
//	  if element2.Found{
//	      . . .
//	  }
//
// If root contains multiple children with the same name, the
// first child always returned.
//
// It returns the first not found Lookup element with the
// [Lookup.Required] flag set, or nil of all required elements
// are found.
//
// Note, if there are missed required elements, it is not
// guaranteed that all lookup requests will be processed.
func (root Element) Lookup(lookups ...*Lookup) *Lookup {
	if len(lookups)*len(root.Children) <= 16 {
		// If we have a small amount of Lookup/Children combination,
		// just go straightforward, it will be less expensive, that
		// doing via maps of children
		for _, l := range lookups {
			l.Elem, l.Found = root.ChildByName(l.Name)
			if !l.Found && l.Required {
				return l
			}
		}
	} else {
		// Otherwise, obtain map of children and then fill the answer.
		children := root.ChildrenMap()
		for _, l := range lookups {
			l.Elem, l.Found = children[l.Name]
			if !l.Found && l.Required {
				return l
			}
		}
	}

	return nil
}

// LookupAttrs is like [Element.Lookup], but for attributes, not for
// children elements.
func (root Element) LookupAttrs(lookups ...*LookupAttr) *LookupAttr {
	if len(lookups)*len(root.Attrs) <= 16 {
		// If we have a small amount of Lookup/Attrs combination,
		// just go straightforward, it will be less expensive, that
		// doing via maps of attributes
		for _, l := range lookups {
			l.Attr, l.Found = root.AttrByName(l.Name)
			if !l.Found && l.Required {
				return l
			}
		}
	} else {
		// Otherwise, obtain map of attributes and then fill the answer.
		attrs := root.AttrsMap()
		for _, l := range lookups {
			l.Attr, l.Found = attrs[l.Name]
			if !l.Found && l.Required {
				return l
			}
		}
	}

	return nil
}
