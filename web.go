package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"github.com/gorilla/mux"
	"io"
	"text/template"
	"github.com/gorilla/sessions"
	"path/filepath"
	//"log"
)

const templatesDir = "./templates"
var store = sessions.NewCookieStore([]byte("something-very-secret"))


func handler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := vars["name"]
	session, e := store.Get(r, name)
	if e!=nil {
		http.Error(w, e.Error(), 500)
	}
	session.Values["name"] = name
	session.Values["Number"] = 42
	data:= &homeData{
		Name: name ,
		Runtime: runtime.Version(),
		Values : session.Values,
	}
	// Save it before we write to the response/return from the handler.
	session.Save(r, w)
	status, e := renderTemplate(w, data, "hello.html")
	if(e!=nil){
		http.Error(w, e.Error(), status)
	}
}

func renderTemplate(w io.Writer, data interface{}, files ...string) (int, error){
	// fix path
	template_name := files[0] // 1st is template ...rest are used by 1st
	// path .join pwd + templates  + file
	paths, e:= rebase(files, templatesDir)
	// load, parse
	tmpl, e := template.New(template_name).ParseFiles(paths...)
	if e != nil {
		return 500, e
	}
	// render
	e = tmpl.Execute(w, data)
	if e != nil {
		return 500, e
	}
	return 200, nil
}

/*
 path .join pwd + templates  + file,
*/
func rebase(paths []string,base string) ([]string , error){
	var out []string;
	pwd, e:= os.Getwd()
	if e!= nil {
		return nil, e
	}
	for _, path:=  range paths {

		combined := filepath.Join(pwd, base, path)
		out = append(out, combined)
		//log.Printf("rebased template path: %v \n" , combined)
	}
	//log.Println(out)
	return out, nil
}

type homeData struct {
	Name    string
	Runtime string
	Values map[interface{}]interface{}
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/hello/{name}", handler)
	r.HandleFunc("/", func(w http.ResponseWriter,r *http.Request){
		http.Redirect(w, r, "/hello/anonymous", http.StatusTemporaryRedirect)
	})
	http.Handle("/", r)

	bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	fmt.Printf("listening on %s... \n", bind)

	// Start: ...
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		panic(err)
	}
}
