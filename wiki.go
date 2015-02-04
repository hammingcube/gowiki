package main
import (
	"os"
	"path/filepath"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Page struct {
	Title string
	Body []byte
}

func savePage(p *Page, dir string) error {
	filename := filepath.Join(dir, p.Title + ".txt")
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string, dir string) (*Page, error) {
	filename := filepath.Join(dir, title + ".txt")
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	tmplFile := filepath.Join(TemplatePath, tmpl + ".html")
	fmt.Println(tmplFile)
	t, err := template.ParseFiles(tmplFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title, ContentPath)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title, ContentPath)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	pages := []*Page{}
	files, err := ioutil.ReadDir(ContentPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".txt" {
			continue
		}
		title := f.Name()[:len(f.Name()) - len(".txt")]
		fmt.Println(title)
		p, err := loadPage(title, ContentPath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		pages = append(pages, p)
	}
	renderTemplate(w, "home", pages)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := savePage(p, ContentPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}



var (
	GoPath = os.Getenv("GOPATH")
	ContentPath = GoPath + "/src/github.com/maddyonline/gowiki/.data/"
	TemplatePath = GoPath + "/src/github.com/maddyonline/gowiki/static/templates" 
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8999", nil)
}


