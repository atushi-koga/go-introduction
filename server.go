package main

import(
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(w, "World, %s\n", p.ByName("name"))
}

func header(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(w, r.Header)
	fmt.Fprintln(w, r.Header["Accept-Encoding"])
	fmt.Fprintln(w, r.Header.Get("Accept-Encoding"))
}

func body(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	length := r.ContentLength
	body := make([]byte, length)
	r.Body.Read(body)
	fmt.Fprintln(w, string(body))
}

func form(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.ParseForm()
	fmt.Fprintln(w, r.Form)
	fmt.Fprintln(w, r.Form["first_name"])
}

func postForm(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.ParseForm()
	fmt.Fprintln(w, r.PostForm)
	fmt.Fprintln(w, r.PostForm["first_name"])
}

func main() {
	mux := httprouter.New()
	mux.GET("/hello/:name", hello)
	mux.GET("/header/", header)
	mux.POST("/body", body)
	mux.POST("/form", form)
	mux.POST("/post_form", postForm)
	server := http.Server{
		Addr: "127.0.0.1:8080",
		Handler: mux,
	}

	server.ListenAndServe()
}