package http

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func ListenAndServe() http.Server {
	port := 8080
	srv := http.Server{Addr: ":" + strconv.Itoa(port)}

	go func() {
		if err := srv.ListenAndServe(); nil != err {
			if http.ErrServerClosed != err {
				log.Fatal(fmt.Errorf("listening on %d: %w", port, err))
			}
		}
	}()
	return srv
}
