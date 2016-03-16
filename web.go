package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"github.com/gorilla/mux"
	"io"
	// "bytes"
	"text/template"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/hello/{name}", handler)
	r.HandleFunc("/", func(w http.ResponseWriter,r *http.Request){
		http.Redirect(w, r, "/hello/anonymous", http.StatusTemporaryRedirect)
	})
	http.Handle("/", r)

	bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	fmt.Printf("listening on %s...", bind)

	// Start: ...
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	renderTemplate(w, "Hello {{.Name}},  runtime:{{ .Version}}", &homeData{name , runtime.Version()})
}

func renderTemplate(w io.Writer, templateString string, data interface{}) {
	tmpl, err := template.New("test").Parse(templateString)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

type homeData struct {
	Name    string
	Version string
}
