# MFP - Miulti-Function Printers and scanners toolkit
# Printer and scanner modeling.
#
# Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
# See LICENSE for license terms and conditions
#
# Python initialization

from dataclasses import dataclass
from uuid import UUID
from typing import TypedDict
from http.client import HTTPMessage

@dataclass
class HTTPQuery:
    Request: HTTPMessage
    Response: HTTPMessage
