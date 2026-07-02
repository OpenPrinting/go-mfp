# MFP - Miulti-Function Printers and scanners toolkit
# Printer and scanner modeling.
#
# Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
# See LICENSE for license terms and conditions
#
# IPP-related definitions

from dataclasses import dataclass
from enum import Enum, IntEnum
from helpers import collection

# IPP operation codes
class OP(IntEnum):
    PRINT_JOB = 0x02
    PRINT_URI = 0x03
    VALIDATE_JOB = 0x04
    CREATE_JOB = 0x05
    SEND_DOCUMENT = 0x06
    SEND_URI = 0x07
    CANCEL_JOB = 0x08
    GET_JOB_ATTRIBUTES = 0x09
    GET_JOBS = 0x0a
    GET_PRINTER_ATTRIBUTES = 0x0b
    HOLD_JOB = 0x0c
    RELEASE_JOB = 0x0d
    RESTART_JOB = 0x0e
    PAUSE_PRINTER = 0x10
    RESUME_PRINTER = 0x11
    PURGE_JOBS = 0x12
    SET_PRINTER_ATTRIBUTES = 0x13
    SET_JOB_ATTRIBUTES = 0x14
    GET_PRINTER_SUPPORTED_VALUES = 0x15
    CREATE_PRINTER_SUBSCRIPTIONS = 0x16
    CREATE_JOB_SUBSCRIPTIONS = 0x17
    GET_SUBSCRIPTION_ATTRIBUTES = 0x18
    GET_SUBSCRIPTIONS = 0x19
    RENEW_SUBSCRIPTION = 0x1a
    CANCEL_SUBSCRIPTION = 0x1b
    GET_NOTIFICATIONS = 0x1c
    SEND_NOTIFICATIONS = 0x1d
    GET_RESOURCE_ATTRIBUTES = 0x1e
    GET_RESOURCE_DATA = 0x1f
    GET_RESOURCES = 0x20
    GET_PRINT_SUPPORT_FILES = 0x21
    ENABLE_PRINTER = 0x22
    DISABLE_PRINTER = 0x23
    PAUSE_PRINTER_AFTER_CURRENT_JOB = 0x24
    HOLD_NEW_JOBS = 0x25
    RELEASE_HELD_NEW_JOBS = 0x26
    DEACTIVATE_PRINTER = 0x27
    ACTIVATE_PRINTER = 0x28
    RESTART_PRINTER = 0x29
    SHUTDOWN_PRINTER = 0x2a
    STARTUP_PRINTER = 0x2b
    REPROCESS_JOB = 0x2c
    CANCEL_CURRENT_JOB = 0x2d
    SUSPEND_CURRENT_JOB = 0x2e
    RESUME_JOB = 0x2f
    PROMOTE_JOB = 0x30
    SCHEDULE_JOB_AFTER = 0x31
    CANCEL_DOCUMENT = 0x33
    GET_DOCUMENT_ATTRIBUTES = 0x34
    GET_DOCUMENTS = 0x35
    DELETE_DOCUMENT = 0x36
    SET_DOCUMENT_ATTRIBUTES = 0x37
    CANCEL_JOBS = 0x38
    CANCEL_MY_JOBS = 0x39
    RESUBMIT_JOB = 0x3a
    CLOSE_JOB = 0x3b
    IDENTIFY_PRINTER = 0x3c
    VALIDATE_DOCUMENT = 0x3d
    ADD_DOCUMENT_IMAGES = 0x3e
    ACKNOWLEDGE_DOCUMENT = 0x3f
    ACKNOWLEDGE_IDENTIFY_PRINTER = 0x40
    ACKNOWLEDGE_JOB = 0x41
    FETCH_DOCUMENT = 0x42
    FETCH_JOB = 0x43
    GET_OUTPUT_DEVICE_ATTRIBUTES = 0x44
    UPDATE_ACTIVE_JOBS = 0x45
    DEREGISTER_OUTPUT_DEVICE = 0x46
    UPDATE_DOCUMENT_STATUS = 0x47
    UPDATE_JOB_STATUS = 0x48
    UPDATE_OUTPUT_DEVICE_ATTRIBUTES = 0x49
    GET_NEXT_DOCUMENT_DATA = 0x4a
    ALLOCATE_PRINTER_RESOURCES = 0x4b
    CREATE_PRINTER = 0x4c
    DEALLOCATE_PRINTER_RESOURCES = 0x4d
    DELETE_PRINTER = 0x4e
    GET_PRINTERS = 0x4f
    SHUTDOWN_ONE_PRINTER = 0x50
    STARTUP_ONE_PRINTER = 0x51
    CANCEL_RESOURCE = 0x52
    CREATE_RESOURCE = 0x53
    INSTALL_RESOURCE = 0x54
    SEND_RESOURCE_DATA = 0x55
    SET_RESOURCE_ATTRIBUTES = 0x56
    CREATE_RESOURCE_SUBSCRIPTIONS = 0x57
    CREATE_SYSTEM_SUBSCRIPTIONS = 0x58
    DISABLE_ALL_PRINTERS = 0x59
    ENABLE_ALL_PRINTERS = 0x5a
    GET_SYSTEM_ATTRIBUTES = 0x5b
    GET_SYSTEM_SUPPORTED_VALUES = 0x5c
    PAUSE_ALL_PRINTERS = 0x5d
    PAUSE_ALL_PRINTERS_AFTER_CURRENT_JOB = 0x5e
    REGISTER_OUTPUT_DEVICE = 0x5f
    RESTART_SYSTEM = 0x60
    RESUME_ALL_PRINTERS = 0x61
    SET_SYSTEM_ATTRIBUTES = 0x62
    SHUTDOWN_ALL_PRINTERS = 0x63
    STARTUP_ALL_PRINTERS = 0x64
    CUPS_GET_DEFAULT = 0x4001
    CUPS_GET_PRINTERS = 0x4002
    CUPS_ADD_MODIFY_PRINTER = 0x4003
    CUPS_DELETE_PRINTER = 0x4004
    CUPS_GET_CLASSES = 0x4005
    CUPS_ADD_MODIFY_CLASS = 0x4006
    CUPS_DELETE_CLASS = 0x4007
    CUPS_ACCEPT_JOBS = 0x4008
    CUPS_REJECT_JOBS = 0x4009
    CUPS_SET_DEFAULT = 0x400a
    CUPS_GET_DEVICES = 0x400b
    CUPS_GET_PPDS = 0x400c
    CUPS_MOVE_JOB = 0x400d
    CUPS_AUTHENTICATE_JOB = 0x400e
    CUPS_GET_PPD = 0x400f
    CUPS_GET_DOCUMENT = 0x4027
    CUPS_CREATE_LOCAL_PRINTER = 0x4028

    # Formatting
    def __repr__ (self):
        return 'ipp.OP.' + self.name

# IPP tags
class TAG(Enum):
    # Delimiters
    ZERO = 0x00
    END = 0x03

    # Groups of attributes
    OPERATION = 0x01
    JOB = 0x02
    PRINTER = 0x04
    UNSUPPORTED_GROUP = 0x05
    SUBSCRIPTION = 0x06
    EVENT_NOTIFICATION = 0x07
    RESOURCE = 0x08
    DOCUMENT = 0x09
    SYSTEM = 0x0a

    # Special values
    UNSUPPORTED_VALUE = 0x10
    DEFAULT = 0x11
    UNKNOWN = 0x12
    NOVALUE = 0x13
    NOTSETTABLE = 0x15
    DELETEATTR = 0x16
    ADMINDEFINE = 0x17

    # Values
    INTEGER = 0x21
    BOOLEAN = 0x22
    ENUM = 0x23
    STRING = 0x30
    DATE = 0x31
    RESOLUTION = 0x32
    RANGE = 0x33
    TEXTLANG = 0x35
    NAMELANG = 0x36
    TEXT = 0x41
    NAME = 0x42
    KEYWORD = 0x44
    URI = 0x45
    URISCHEME = 0x46
    CHARSET = 0x47
    LANGUAGE = 0x48
    MIMETYPE = 0x49
    EXTENSION = 0x7f

    # Collections
    BEGIN_COLLECTION = 0x34
    END_COLLECTION = 0x37
    MEMBERNAME = 0x4a

    # Formatting
    def __repr__ (self):
        return 'ipp.TAG.' + self.name

# IPP_TAG_UNSUPPORTED_VALUE
class UNSUPPORTED_VALUE:
    def __repr__ (self):
        return 'ipp.UNSUPPORTED_VALUE()'

# IPP_TAG_DEFAULT
class DEFAULT:
    def __repr__ (self):
        return 'ipp.DEFAULT()'

# IPP_TAG_UNKNOWN
class UNKNOWN:
    def __repr__ (self):
        return 'ipp.UNKNOWN()'

# IPP_TAG_NOVALUE
class NOVALUE:
    def __repr__ (self):
        return 'ipp.NOVALUE()'

# IPP_TAG_NOTSETTABLE
class NOTSETTABLE:
    def __repr__ (self):
        return 'ipp.NOTSETTABLE()'

# IPP_TAG_DELETEATTR
class DELETEATTR:
    def __repr__ (self):
        return 'ipp.DELETEATTR()'

# IPP_TAG_ADMINDEFINE
class ADMINDEFINE:
    def __repr__ (self):
        return 'ipp.ADMINDEFINE()'

# IPP_TAG_INTEGER
class INTEGER(int):
    def __repr__ (self):
        return 'ipp.INTEGER(' + repr(int(self)) + ')'

# IPP_TAG_BOOLEAN
@dataclass
class BOOLEAN:
    Val : bool

    def __bool__ (self):
        return self.Val

    def __repr__ (self):
        return 'ipp.BOOLEAN(' + repr(bool(self.Val)) + ')'

# IPP_TAG_ENUM
class ENUM(int):
    def __repr__ (self):
        return 'ipp.ENUM(' + repr(int(self)) + ')'

# IPP_TAG_STRING
class STRING(str):
    def __repr__ (self):
        return 'ipp.STRING(' + repr(str(self)) + ')'

# IPP_TAG_DATE
class DATE(str):
    def __repr__ (self):
        return 'ipp.DATE(' + repr(str(self)) + ')'

# IPP_TAG_RESOLUTION
@dataclass
class RESOLUTION:
    X: int
    Y: int
    Units: str

    def __repr__ (self):
        return 'ipp.RESOLUTION(' + repr(self.X) + ', ' + repr(self.Y) + ', ' + repr(self.Units) + ')'

# IPP_TAG_RANGE
@dataclass
class RANGE:
    Lower: int
    Upper: int

    def __repr__ (self):
        return 'ipp.RANGE(' + repr(self.Lower) + ', ' + repr(self.Upper) + ')'

# IPP_TAG_TEXTLANG
@dataclass
class TEXTLANG:
    Text: str
    Lang: str

    def __repr__ (self):
        return 'ipp.TEXTLANG(' + repr(self.Text) + ', ' + repr(self.Lang) + ')'

# IPP_TAG_NAMELANG
@dataclass
class NAMELANG:
    Name: str
    Lang: str

    def __repr__ (self):
        return 'ipp.NAMELANG(' + repr(self.Name) + ', ' + repr(self.Lang) + ')'

# IPP_TAG_TEXT
class TEXT(str):
    def __repr__ (self):
        return 'ipp.TEXT(' + repr(str(self)) + ')'

# IPP_TAG_NAME
class NAME(str):
    def __repr__ (self):
        return 'ipp.NAME(' + repr(str(self)) + ')'

# IPP_TAG_KEYWORD
class KEYWORD(str):
    def __repr__ (self):
        return 'ipp.KEYWORD(' + repr(str(self)) + ')'

# IPP_TAG_URI
class URI(str):
    def __repr__ (self):
        return 'ipp.URI(' + repr(str(self)) + ')'

# IPP_TAG_URISCHEME
class URISCHEME(str):
    def __repr__ (self):
        return 'ipp.URISCHEME(' + repr(str(self)) + ')'

# IPP_TAG_CHARSET
class CHARSET(str):
    def __repr__ (self):
        return 'ipp.CHARSET(' + repr(str(self)) + ')'

# IPP_TAG_LANGUAGE
class LANGUAGE(str):
    def __repr__ (self):
        return 'ipp.LANGUAGE(' + repr(str(self)) + ')'

# IPP_TAG_MIMETYPE
class MIMETYPE(str):
    def __repr__ (self):
        return 'ipp.MIMETYPE(' + repr(str(self)) + ')'

# IPP_TAG_BEGIN_COLLECTION/IPP_TAG_END_COLLECTION
class COLLECTION(collection): pass

# attrs is the model-settable variable that defines the
# IPP printer attributes
attrs = None
