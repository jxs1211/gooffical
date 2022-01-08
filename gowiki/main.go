package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/golang/glog"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

var validPath = regexp.MustCompile("^/(edit|view|save)/([A-Za-z0-9]+)$")

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title
	return os.WriteFile(filename, p.Body, 0600)
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("not match url")
	}
	return m[2], nil
}

func renderTemplate(w http.ResponseWriter, t_name string, p *Page) {
	// t, err := template.ParseFiles(t_name)
	err := templates.ExecuteTemplate(w, t_name+".html", p)
	if err != nil {
		log.Println("err: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func loadPage(filename string) (*Page, error) {
	if filename == "" {
		return nil, errors.New("filename must not be empty")
	}
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: filename, Body: b}, nil
}

// Data structures
func showData() {
	title := "TestPage"
	p := Page{
		Title: title,
		Body:  []byte("This is a sample Page."),
	}
	_ = p.save()
	res, _ := loadPage(title)
	fmt.Printf("p: %+v, p.Body: %s\n", *res, res.Body)
}

// Introducing the net/http package (an interlude)
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func recoverPanic(w http.ResponseWriter) {
	if p := recover(); p != nil {
		fmt.Fprintf(w, "panic: %s\n", p)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	defer recoverPanic(w)
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		log.Println("err: ", err)
		// fmt.Fprintf(w, "err: %v\n", err)
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	defer recoverPanic(w)
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		log.Println("err: ", err)
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func viewHandler2(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		log.Println("err: ", err)
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler2(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		log.Println("err: ", err)
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler2(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	// http.HandleFunc("/", handler)
	flag.Set("v", "4")
	glog.V(2).Info("Starting http server...")
	http.HandleFunc("/view/", makeHandler(viewHandler2))
	http.HandleFunc("/edit/", makeHandler(editHandler2))
	http.HandleFunc("/save/", makeHandler(saveHandler2))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
