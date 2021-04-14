// Package main 测试用例package main
package main

import (
	"log"
	"net/http"
	"runtime"
)

const (
	httpPort = "8088"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	hello := func(w http.ResponseWriter, req *http.Request) {
		data := "Hello, go-stress-testing! \n"
		w.Header().Add("Server", "golang")
		_, _ = w.Write([]byte(data))
		return
	}

	http.HandleFunc("/", hello)
	err := http.ListenAndServe(":"+httpPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
