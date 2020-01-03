package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"time"
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

func writeHTMLToResponseBody(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	str := `<html>
<head><title>Go Web Programing</title></head>
<body><h1>Hello World</h1></body>
</html>
`
	w.Write([]byte(str))
}

func setStatusCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.WriteHeader(501)
	fmt.Fprintln(w, "そのようなサービスはありません")
}

func writeResponseHeader(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Location", "http://www.google.com")
	w.WriteHeader(302)
}

type Post struct {
	User string
	Threads []string
}

func writeJSONToResponseBody(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	post := &Post{
		User:    "Yamada Tarou",
		Threads: []string{"a", "b", "c"},
	}
	jsonData, _ := json.Marshal(post)
	w.Write(jsonData)
}

func setCookie(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	msg1 := []byte("Hello World 100%")
	c1 := http.Cookie{
		Name:       "first_cookie",
		Value:      base64.URLEncoding.EncodeToString(msg1),
		HttpOnly:   true,
	}

	c2 := http.Cookie{
		Name:       "second_cookie",
		Value:      "Programing Language Go",
		HttpOnly:   true,
	}
	w.Header().Set("Set-Cookie", c1.String())
	w.Header().Add("Set-Cookie", c2.String())
}

func getCookie(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	firstCookie, err := r.Cookie("first_cookie")
	if err != nil {
		fmt.Fprintln(w, "Can not get the first cookie")
	}
	cookies := r.Cookies()
	fmt.Fprintln(w, firstCookie)
	fmt.Fprintln(w, cookies)
}

func deleteAndShowCookie(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c, err := r.Cookie("first_cookie")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Fprintln(w, "Cookieがありません")
		}
	}else{
		c1 := http.Cookie{
			Name:    "first_cookie",
			MaxAge:   -1,
			Expires: time.Unix(1, 0),
		}
		w.Header().Set("Set-Cookie", c1.String())

		val, _ := base64.URLEncoding.DecodeString(c.Value)
		fmt.Fprintln(w, string(val))
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
	mux.GET("/write_html", writeHTMLToResponseBody)
	mux.GET("/write_json", writeJSONToResponseBody)
	mux.GET("/set_status_code", setStatusCode)
	mux.GET("/write_response_header", writeResponseHeader)
	mux.GET("/set_cookie", setCookie)
	mux.GET("/get_cookie", getCookie)
	mux.GET("/delete_and_show_cookie", deleteAndShowCookie)
	server := http.Server{
		Addr: "127.0.0.1:8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
