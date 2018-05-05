package plugin
//
//import (
//	"fmt"
//	"log"
//	"net/http"
//	"github.com/scofier/gogateway/proxy"
//	"github.com/scofier/gogateway/manage"
//	"encoding/json"
//	"github.com/scofier/gogateway/interceptor"
//)
////server port
//var port=":9090"
//
////demo service, http://127.0.0.1:9090/demo/v1/info
//var serviceRegistry = proxy.DefaultRegistry{
//	"demo": {
//		"v1": {
//			"localhost" + port,
//		},
//	},
//}
//
//func main() {
//
//	http.HandleFunc("/", proxy.NewMultipleHostReverseProxy(&serviceRegistry, interceptor.JwtInterceptor))
//	//info
//	http.HandleFunc("/info", func(w http.ResponseWriter, req *http.Request) {
//		value, _:=json.Marshal(serviceRegistry)
//		fmt.Fprintf(w, "%s", value)
//	})
//	//ServiceRegistry handler, you can remove this if you don't need
//	manage.RegistryHandler(&serviceRegistry)
//
//	interceptor.SsoHandler()
//
//	//start server
//	println("Start on", port)
//	log.Fatal(http.ListenAndServe(port, nil))
//}