package server

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

// temporary directory location
var tmpDir = filepath.FromSlash("/home/chicmic/Desktop/test")

func Server() {

    // default route
    http.HandleFunc( "/", func( res http.ResponseWriter, req *http.Request ) {
        fmt.Fprint( res, "Hello Golang!" )
    } )

    // return a `.html` file for `/index.html` route
    http.HandleFunc( "/maggie", func( res http.ResponseWriter, req *http.Request ) {
        http.ServeFile( res, req, filepath.Join( tmpDir, "/123.mp3" ) );
    } )

    // start HTTP server with `http.DefaultServeMux` handler
    log.Fatal(http.ListenAndServe( ":9000", nil ))

}