// cd into this script's directory and type: go run simpleserver.go
// go to localhost:8080/somestring; i.e. http://localhost:8080/donuts

package main

import (
  "fmt"
  "net/http"
)

//this function is of type http.HandlerFunc, takes http.ResponseWriter and http.Request as args
func handler(w http.ResponseWriter, r *http.Request){
  fmt.Fprint(w, "I give you the power of ", r.URL.Path[1:])
}

//r.IR:PATH is a http.Request data structure that is the path component of the request URL
//handler function slices the path to omit the leading "/" from path name

func main(){
  http.HandleFunc("/", handler)   //tells http package to handle all requests to web root("/") with handler
  http.ListenAndServe(":8080", nil) //specify port, blocks until program is terminated
}
