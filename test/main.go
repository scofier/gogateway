package main

import "fmt"
//
//import (
//	"fmt"
//	"log"
//	"net/http"
//
//	"github.com/creack/goproxy"
//	"github.com/creack/goproxy/registry"
//)
//
//// ServiceRegistry is a local registry of services/versions
//var ServiceRegistry = registry.DefaultRegistry{
//	"demo": {
//		"v1": {
//			"localhost:12341",
//			"localhost:12342",
//		},
//	},
//}
//
//func main() {
//	http.HandleFunc("/", goproxy.NewMultipleHostReverseProxy(ServiceRegistry))
//	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
//		fmt.Fprintf(w, "%v\n", ServiceRegistry)
//	})
//	println("ready")
//	log.Fatal(http.ListenAndServe(":9090", nil))
//}
type path []byte

func (p path) ToUpper() {
	for i, b := range p {
		if 'a' <= b && b <= 'z' {
			p[i] = b + 'A' - 'a'
		}
	}
}

func main() {
	pathName := path("/usr/bin/tso")
	pathName.ToUpper()
	fmt.Printf("%s\n", pathName)
}