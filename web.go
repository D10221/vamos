package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"path"
	"github.com/gorilla/schema"
	"log"
	"github.com/D10221/vamos/shared"
	db "github.com/D10221/vamos/shared/db"
)


var store = sessions.NewCookieStore([]byte("something-very-secret"))

func signinHandler(w http.ResponseWriter, r *http.Request){

	err := r.ParseForm()

	if err != nil {
		log.Printf("ParseForm: e: %v \n", err.Error())
	}

	decoder := schema.NewDecoder()
	// r.PostForm is a map of our POST form values
	user:= &db.User{}
	err = decoder.Decode(user, r.PostForm)
	if err != nil {
		log.Printf("ParseForm.User: e: %v \n", err.Error())
	}
	log.Print(user)

	data:= struct { Title string} {"SignIn"}
	shared.RenderTemplate(w, "sigin.html", data)
}

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
	status, e := shared.RenderTemplate(w, data, "hello.html")
	if e!=nil {
		http.Error(w, e.Error(), status)
	}
}



type homeData struct {
	Name    string
	Runtime string
	Values map[interface{}]interface{}
}

func main() {

	r := mux.NewRouter()

	// a route.
	r.HandleFunc("/hello/{name}", handler)
	r.HandleFunc("/signin", signInHandler)

	// root
	r.HandleFunc("/", func(w http.ResponseWriter,r *http.Request){
		http.Redirect(w, r, "/hello/anonymous", http.StatusTemporaryRedirect)
	})

	// Static
	//if a path not found until now, e.g. "/image/tiny.png"
	//this will look at "./public/image/tiny.png" at filesystem
	r.PathPrefix("/favicon.png").HandlerFunc(ServeFileHandler)
	http.Handle("/", r)

	// Start: ...
	bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	fmt.Printf("listening on %s... \n", bind)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		panic(err)
	}
}

func ServeFileHandler(res http.ResponseWriter, req *http.Request) {
	fname := path.Base(req.URL.Path)
	http.ServeFile(res, req, "./static/images/"+fname)
}
