// MFP       - Miulti-Function Printers and scanners toolkit
// TRANSPORT - Transport protocol implementation
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// HTTP server

package transport

import (
	"context"
	syslog "log"
	"net"
	"net/http"
	"sync"

	"github.com/OpenPrinting/go-mfp/log"
)

// Server wraps [http.Server]
type Server struct {
	http.Server                 // Underlying http.Server
	ctx         context.Context // Server context
	handler     http.Handler    // Request handler
}

// NewServer creates a new [Server].
//
// The provided [context.Context] used for logging and becomes
// the base context for incoming requests.
//
// The [http.Server] used only as configuration template. If it
// is nil, the reasonable defaults will be used instead.
func NewServer(ctx context.Context,
	template *http.Server, handler http.Handler) *Server {
	if template == nil {
		template = &http.Server{}
	}

	errorLog := syslog.New(log.Writer(ctx, log.LevelWarning), "", 0)
	srvr := &Server{
		Server: http.Server{
			Addr:                         template.Addr,
			DisableGeneralOptionsHandler: template.DisableGeneralOptionsHandler,
			TLSConfig:                    template.TLSConfig,
			ReadTimeout:                  template.ReadTimeout,
			ReadHeaderTimeout:            template.ReadHeaderTimeout,
			WriteTimeout:                 template.WriteTimeout,
			IdleTimeout:                  template.IdleTimeout,
			MaxHeaderBytes:               template.MaxHeaderBytes,
			TLSNextProto:                 template.TLSNextProto,
			ConnState:                    template.ConnState,
			BaseContext:                  template.BaseContext,
			ConnContext:                  template.ConnContext,
			ErrorLog:                     errorLog,
		},
		ctx:     ctx,
		handler: handler,
	}

	srvr.Server.BaseContext = func(net.Listener) context.Context {
		return srvr.ctx
	}

	srvr.Handler = http.HandlerFunc(srvr.handlerFunc)

	return srvr
}

// handlerFunc wraps the http.Server.Handler.
func (srvr *Server) handlerFunc(w http.ResponseWriter, r *http.Request) {
	// Catch panics to log
	defer func() {
		v := recover()
		if v != nil {
			log.Panic(srvr.ctx, v)
		}
	}()

	// Call the handler
	srvr.handler.ServeHTTP(w, r)
}

// ServeAutoTLS is similar to the [http.Server.Serve] and
// [http.Server.ServeTLS].
//
// It accepts incoming connections on the [http.Listener],
// automatically detect encrypted (TLS) and non-encrypted
// connections and serves them respectively to their type.
//
// This function runs "forewer" and returns only in a case
// of error. Use Server.Shutdown or Server.Close to force
// this function to exit.
func (srvr *Server) ServeAutoTLS(l net.Listener) error {
	plain, encrypted := NewAutoTLSListener(l)

	errchan := make(chan error, 2)
	var done sync.WaitGroup

	// Go's http.Server lazily initializes its TLS and HTTP/2
	// internals at the start of Server.Serve and Server.ServeTLS.
	//
	// If these functions are started concurrently, the server
	// may end up with inconsistently initialized state.
	//
	// I ran into this on Go 1.24.2 and could not reproduce it
	// on other versions, but that doesn't mean the problem
	// doesn't exist elsewhere.
	//
	// Details (in Russian):
	//
	//   https://habr.com/ru/articles/906796/
	//
	// As a workaround, we first call ServeTLS with a Listener
	// that fails immediately, causing the function to exit.
	// Before hitting the Listener and returning, however, it
	// performs the lazy initialization — which is exactly what
	// we want.
	srvr.ServeTLS(ErrListener{}, "", "")

	done.Add(2)

	go func() {
		err := srvr.Serve(plain)
		errchan <- err
		done.Done()
	}()

	go func() {
		err := srvr.ServeTLS(encrypted, "", "")
		errchan <- err
		done.Done()
	}()

	err := <-errchan

	plain.Close()
	encrypted.Close()

	done.Wait()

	return err
}
