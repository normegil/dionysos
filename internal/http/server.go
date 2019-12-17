package http

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
)

func ListenAndServe(addr net.TCPAddr, handler http.Handler) func(ctx context.Context) error {
	httpAddr := addr.IP.String() + ":" + strconv.Itoa(addr.Port)
	srv := http.Server{Addr: httpAddr, Handler: handler}

	go func() {
		log.Printf("listening on %s", httpAddr)
		if err := srv.ListenAndServe(); nil != err {
			if http.ErrServerClosed != err {
				log.Fatal(fmt.Errorf("listening on '%s': %w", httpAddr, err))
			}
		}
	}()
	return func(ctx context.Context) error {
		if err := srv.Shutdown(ctx); nil != err {
			return fmt.Errorf("shutdown http server: %w", err)
		}
		return nil
	}
}
