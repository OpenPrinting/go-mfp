# MFP - Miulti-Function Printers and scanners toolkit
# Printer and scanner modeling.
#
# Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
# See LICENSE for license terms and conditions
#
# Template for generated models

#-
# This is the generated MFP model file.
# You probably need to edit it appropriately before use.

#-ipp
# IPP printer attributes:
ipp.attrs = $IPP

#-escl
# eSCL scanner parameters:
escl.caps = $ESCL

# ----- PUT YOUR ESCL HOOKS HERE -----

# Called on request:  POST /{root}/ScanJobs
#
# def escl_onScanJobsRequest (q: query.Query, rq: escl.ScanSettings):

# Called on response: GET /{JobUri}/NextDocument
#
# def escl_onNextDocumentResponse (q: query.Query, flt: escl.ImageFilter):

#-wsd
# WS-Scan (WSD) scanner parameters:
wsd.caps = $WSD
