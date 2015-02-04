package main
import (
	"os"
	"path/filepath"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"flag"
)


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

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
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
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)

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
	ContentPath = GoPath + '.data/'
	SrcPath = "github.com/maddyonline/gowiki"
	TemplatePath = GoPath + "/src/" + SrcPath + "/static/templates" 
)

func init() {
	flag.StringVar(&dataDir, "datadir", filepath.Join(os.Getenv("HOME"), "gowikidata"), "Directory to store/retrieve blog posts")
}


func main() {
	flag.Parse()
	fmt.Println(filepath.Join(dataDir, "abc.txt"))
	return
	switchToDataDir()
	fmt.Println(TemplatePath)
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8999", nil)
}


