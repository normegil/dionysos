package http

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
)

func ListenAndServe(addr net.TCPAddr, handler http.Handler) http.Server {
	httpAddr := addr.IP.String() + ":" + strconv.Itoa(addr.Port)
	srv := http.Server{Addr: httpAddr, Handler: handler}

	go func() {
		if err := srv.ListenAndServe(); nil != err {
			if http.ErrServerClosed != err {
				log.Fatal(fmt.Errorf("listening on '%s': %w", httpAddr, err))
			}
		}
	}()
	return srv
}
