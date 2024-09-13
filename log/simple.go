// MFP - Miulti-Function Printers and scanners toolkit
// Logging facilities
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Simple logging API

package log

import (
	"context"
	"encoding"
)

// Trace writes a Trace-level message to the [Logger] associated
// with the Context.
//
// If Logger is not available, [DefaultLogger] will be used.
// The [context.Context] parameter may be safely passed as nil.
func Trace(ctx context.Context, format string, v ...any) {
	CtxLogger(ctx).Trace(CtxPrefix(ctx), format, v...)
}

// Debug writes a Debug-level message to the [Logger] associated
// with the Context.
//
// If Logger is not available, [DefaultLogger] will be used.
// The [context.Context] parameter may be safely passed as nil.
func Debug(ctx context.Context, format string, v ...any) {
	CtxLogger(ctx).Debug(CtxPrefix(ctx), format, v...)
}

// Info writes a Info-level message to the [Logger] associated
// with the Context.
//
// If Logger is not available, [DefaultLogger] will be used.
// The [context.Context] parameter may be safely passed as nil.
func Info(ctx context.Context, format string, v ...any) {
	CtxLogger(ctx).Info(CtxPrefix(ctx), format, v...)
}

// Warning writes a Warning-level message to the [Logger] associated
// with the Context.
//
// If Logger is not available, [DefaultLogger] will be used.
// The [context.Context] parameter may be safely passed as nil.
func Warning(ctx context.Context, format string, v ...any) {
	CtxLogger(ctx).Warning(CtxPrefix(ctx), format, v...)
}

// Error writes a Error-level message to the [Logger] associated
// with the Context.
//
// If Logger is not available, [DefaultLogger] will be used.
// The [context.Context] parameter may be safely passed as nil.
func Error(ctx context.Context, format string, v ...any) {
	CtxLogger(ctx).Error(CtxPrefix(ctx), format, v...)
}

// Fatal writes a Fatal-level message to the [Logger] associated
// with the Context.
//
// It calls os.Exit(1) and never returns.
//
// If Logger is not available, [DefaultLogger] will be used.
// The [context.Context] parameter may be safely passed as nil.
func Fatal(ctx context.Context, format string, v ...any) {
	CtxLogger(ctx).Fatal(CtxPrefix(ctx), format, v...)
}

// Begin initiates creation of a new multi-line log [Record].
//
// See [Logger.Begin] for details.
func Begin(ctx context.Context) *Record {
	return CtxLogger(ctx).Begin(CtxPrefix(ctx))
}

// Object writes any object that implements [encoding.TextMarshaler]
// interface to the Logger associated with the Context.
//
// If Logger is not available, [DefaultLogger] will be used.
// The [context.Context] parameter may be safely passed as nil.
//
// If [encoding.TextMarshaler.MarshalText] returns an error, it
// will be written to log with the [Error] log level, regardless
// of the level specified by the first parameter.
func Object(ctx context.Context, level Level, indent int,
	obj encoding.TextMarshaler) context.Context {
	CtxLogger(ctx).Object(CtxPrefix(ctx), level, indent, obj)
	return ctx
}
