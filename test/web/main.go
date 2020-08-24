package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type MyMux struct {}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/registry/v3/microservices" {
		sayhelloName(w, r)
		return
	} else {
		log.Println(r.URL.Path)
	}

	http.NotFound(w, r)
	return
}


func sayhelloName(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		return
	}
	r.Header.Set("content-type","application/json")
	r.Header.Set("x-domain-name","default")
	fmt.Println("path:", r.URL.Path)
//	fmt.Fprintf(w, "hello go")



	body, _ := ioutil.ReadAll(r.Body)
	//    r.Body.Close()
	body_str := string(body)
	fmt.Println(body_str)






}

func main() {
	mux := &MyMux{}
	err := http.ListenAndServe(":30100", mux)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}