package main


import (
	
	"io/ioutil"
	"net/http"
	"log"
	"html/template"
)


type Page struct {
	Title string
	Body  []byte
}


func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}


func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil

}

//to view wiki page
func viewHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/view/"):]
	p, err:= loadPage(title)
	//if page doesn't exit
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

//edit handler function
func editHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/edit/"):]
	p, err:= loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	//will read the contents of edit.html 
	renderTemplate(w, "edit", p)
}

//save page handler
func saveHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte (body)}
	err :=p.save()
	//err handling here
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
	
}


//function to handle rendering of template
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    t, err := template.ParseFiles(tmpl + ".html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = t.Execute(w, p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}


func main() {
	//editPage handler
	http.HandleFunc	("/edit/", editHandler)
	//savePage Handler
	http.HandleFunc	("/save/", saveHandler)
	//view page handler
	http.HandleFunc("/view/", viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}


