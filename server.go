package main

import(
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
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

func formValue(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// リクエストされたデータを対象に、構造体Formの中に該当する値が複数あっても、最初の要素しか取得しない
	fmt.Fprintln(w, r.FormValue("first_name"))
}

func postFormValue(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// フォームのデータのみを対象に、構造体Formの中に該当する値が複数あっても、最初の要素しか取得しない
	fmt.Fprintln(w, r.PostFormValue("first_name"))
}

func multipartForm(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.ParseMultipartForm(1024)
	fileHeader := r.MultipartForm.File["uploaded"][0]
	file, err := fileHeader.Open()
	if err == nil {
		data, err := ioutil.ReadAll(file)
		if err == nil {
			fmt.Fprintln(w, string(data))
		}
	}
}

func formFile(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	file, _, err := r.FormFile("uploaded")
	if err == nil {
		data, err := ioutil.ReadAll(file)
		if err == nil {
			fmt.Fprintln(w, string(data))
		}
	}
}

func main() {
	mux := httprouter.New()
	mux.GET("/hello/:name", hello)
	mux.GET("/header/", header)
	mux.POST("/body", body)
	mux.POST("/form", form)
	mux.POST("/post_form", postForm)
	mux.POST("/form_value", formValue)
	mux.POST("/post_form_value", postFormValue)
	mux.POST("/multipart_form", multipartForm)
	mux.POST("/form_file", formFile)
	server := http.Server{
		Addr: "127.0.0.1:8080",
		Handler: mux,
	}

	server.ListenAndServe()
}