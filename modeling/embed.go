// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Embedded files

package modeling

import _ "embed" // For go:embed to work

//go:embed model.py
var embedPyModel string

//go:embed legacy_model.py
var embedPyLegacyModel string

//go:embed init.py
var embedPyInit string

//go:embed helpers.py
var embedPyHelpers string

//go:embed query.py
var embedPyQuery string

//go:embed ipp.py
var embedPyIPP string

//go:embed escl.py
var embedPyEscl string

//go:embed wsd.py
var embedPyWSD string
