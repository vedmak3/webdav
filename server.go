package main // import "webdav-server"

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/webdav"
)

func init() {
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}
func main() {
	arg := os.Args[1:]
	rand.Seed(time.Now().UnixNano())
	if len(arg) == 2 {
		salt := RandStringRunes(8)
		fmt.Println(string(salt))
		os.Mkdir("./storage", 0777)
		storagePath := "./storage"

		srv := &webdav.Handler{
			Prefix:     "/" + salt,
			FileSystem: webdav.Dir(storagePath),
			LockSystem: webdav.NewMemLS(),
			Logger: func(r *http.Request, err error) {
				if err != nil {
					fmt.Printf("WebDAV %s: %s, ERROR: %s\n", r.Method, r.URL, err)
				}
			},
		}

		mux := http.NewServeMux()
		// Trailing slash must be inputed to end of path in http.HandleFunc
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println(r)
		})
		mux.HandleFunc("/"+salt+"/", func(w http.ResponseWriter, r *http.Request) {
			username, password, _ := r.BasicAuth()

			// Check credential
			if username == arg[0] && password == arg[1] {
				// User control at here, if required.

				srv.ServeHTTP(w, r)
				return
			}

			w.Header().Set("WWW-Authenticate", `Basic realm="BASIC WebDAV REALM"`)
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
		})

		tlsServer(":443", mux)
	}
}
