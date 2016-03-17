package shared

import (
	"io"
	"html/template"
	"os"
	"path/filepath"
)

const templatesDir = "./templates"

func RenderTemplate(w io.Writer, data interface{}, files ...string) (int, error) {
	// fix path
	template_name := files[0] // 1st is template ...rest are used by 1st
	// path .join pwd + templates  + file
	paths, e := rebase(files, templatesDir)
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
func rebase(paths []string, base string) ([]string, error) {
	var out []string;
	pwd, e := os.Getwd()
	if e != nil {
		return nil, e
	}
	for _, path := range paths {

		combined := filepath.Join(pwd, base, path)
		out = append(out, combined)
		//log.Printf("rebased template path: %v \n" , combined)
	}
	//log.Println(out)
	return out, nil
}
