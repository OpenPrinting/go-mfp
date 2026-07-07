# MFP - Miulti-Function Printers and scanners toolkit
# Printer and scanner modeling.
#
# Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
# See LICENSE for license terms and conditions
#
# Python helper classes

from types import SimpleNamespace
import reprlib
import sys

# collection is the base class for classes that are essentially
# a collection of attributes.
class collection(SimpleNamespace):
    def __repr__(self):
        # formatter adds indentation to the standard reprlib.Repr.
        #
        # We have  to implement it for compatibility with Python versions
        # before 3.12 (where indentation is supported by stdlib)
        class formatter(reprlib.Repr):
            def __init__(self):
                super().__init__()

                # Disable all length truncation limits for lossless formatting
                self.maxdict = sys.maxsize
                self.maxlist = sys.maxsize
                self.maxtuple = sys.maxsize
                self.maxset = sys.maxsize
                self.maxfrozenset = sys.maxsize
                self.maxdeque = sys.maxsize
                self.maxarray = sys.maxsize
                self.maxstring = sys.maxsize
                self.maxlong = sys.maxsize
                self.maxother = sys.maxsize

                # Configure identation
                self.indent_str = "    "

            # _format_sequence is the helper function that formats sequences
            def _format_sequence(self, items, left_bracket, right_bracket, level):
                if not items:
                    return f"{left_bracket}{right_bracket}"
                if level <= 0:
                    return f"{left_bracket}...{right_bracket}"

                pieces = []
                newline_indent = "\n" + self.indent_str

                for item in items:
                    item_repr = self.repr(item)
                    if "\n" in item_repr:
                        item_repr = item_repr.replace("\n", newline_indent)
                    pieces.append(f"{self.indent_str}{item_repr}")

                # Fix: Extract join away from f-string expression
                joined_elements = ",\n".join(pieces)
                return f"{left_bracket}\n{joined_elements},\n{right_bracket}"

            # repr_dict formats nested dictionaries
            def repr_dict(self, x, level):
                if not x:
                    return "{}"
                if level <= 0:
                    return "{...}"

                pieces = []
                newline_indent = "\n" + self.indent_str

                for key, value in x.items():
                    val_repr = self.repr(value)
                    if "\n" in val_repr:
                        val_repr = val_repr.replace("\n", newline_indent)
                    pieces.append(f"{self.indent_str}{self.repr(key)}: {val_repr}")

                # Fix: Extract join away from f-string expression
                joined_dict = ",\n".join(pieces)
                return "{\n" + joined_dict + ",\n}"

            # repr_list formats nested lists
            def repr_list(self, x, level):
                return self._format_sequence(x, "[", "]", level)

            # repr_list formats nested tuples
            def repr_tuple(self, x, level):
                if len(x) == 1:
                    newline_indent = "\n" + self.indent_str
                    single_repr = self.repr(x[0]).replace("\n", newline_indent)
                    return f"(\n{self.indent_str}{single_repr},\n)"
                return self._format_sequence(x, "(", ")", level)

            # repr_list formats nested sets
            def repr_set(self, x, level):
                return self._format_sequence(x, "{", "}", level)

            # repr_list formats nested frozen sets
            def repr_frozenset(self, x, level):
                inner = self._format_sequence(x, "{", "}", level)
                if "\n" in inner:
                    newline_indent = "\n" + self.indent_str
                    inner = inner.replace("\n", newline_indent)
                return f"frozenset(\n{self.indent_str}{inner},\n)"

        # Format attributes, one by one
        fmt = formatter()
        pieces = []
        indent = "    "
        newline_indent = "\n" + indent

        for key, value in vars(self).items():
            val_repr = fmt.repr(value)
            if "\n" in val_repr:
                val_repr = val_repr.replace("\n", newline_indent)

            pieces.append(f"{indent}{key} = {val_repr}")

        # Join formatted attributes
        joined_pieces = ""
        if len(pieces) > 0:
            joined_pieces = ",\n".join(pieces)
            joined_pieces = f"\n{joined_pieces},\n"

        class_name = self.__class__.__module__ + "." + self.__class__.__name__

        return f"{class_name}({joined_pieces})"

# meta_keyword is a hepler metaclass for keyword implementation
class meta_keyword (type):
    def __str__ (self):
        return self.__name__

    def __repr__ (self):
        return self.__module__ + "." + self.__name__

# keyword is a base class of classes that can be used
# as keywords.
#
# for classes, inherited from keyword:
#   - str(cls) returns the keyword name (i.e., cls.__name__)
#   - repr(cls) returns keyword representation in the Python
#     syntax (i.e. cls.__module__ + "." + cls.__name__)
class keyword (metaclass = meta_keyword):
    pass

# iskeyword reports if cls is a keyword
def iskeyword(cls):
    # issubclass will trow an exception, if cls is not
    # a type object. Just catch it and return false, if
    # it happens.
    try:
        return issubclass(cls, keyword)
    except:
        pass

    return False

# Purge unexported symbols from the module namespace.
del meta_keyword
