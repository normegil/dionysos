package main

import "net/http"

func main() {
	if err := http.ListenAndServe(":8080", nil); nil != err {
		panic(err)
	}
}
