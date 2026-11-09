# LPD protocol implementation

```
import "github.com/OpenPrinting/go-mfp/proto/lpd"
```

This package implements the server side of the Line Printer Daemon
(LPD) protocol, as defined by [RFC 1179](https://www.rfc-editor.org/rfc/rfc1179).

The [cups-lpd](https://github.com/OpenPrinting/cups/blob/master/scheduler/cups-lpd.c)
mini-daemon from CUPS was used as the practical reference.
