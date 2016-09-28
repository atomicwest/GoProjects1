// create wiki pages

package main

import (
  // "fmt"
  "io/ioutil"
  "net/http"
  "html/template"
  "regexp"
  "errors"
)

// global variable to store template names, .Must will abort the application if there are loading errors
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// global var for regex validations
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")  //aborts if expression compilation fails

//use validPath to validate path, retrieve page title
func getTitle(w http.ResponseWriter, r *http.Request) (string, error){
  m := validPath.FindStringSubmatch(r.URL.Path)
  if m == nil {
    http.NotFound(w,r)
    return "", errors.New("Not a Valid Page Title")
  }
  return m[2], nil    //m[2] is the title, returned with nil errors if valid title
}


// create a struct and describe how page data will be stored in memory
type Page struct {
  Title string
  Body []byte     //byte slice represents the page content
}

// create save method on the Page struct
// saves body to a text file with the same name as page's title
func (p *Page) save() error {
  filename := p.Title + ".txt"
  return ioutil.WriteFile(filename, p.Body, 0600)
}

// create load method for pages
func loadPage(title string) (*Page, error) {
  filename := title + ".txt"
  // body, _ := ioutil.ReadFile(filename) //blank identifier (underscore) throws out error
  body, err := ioutil.ReadFile(filename)
  // if file does not exist
  if err != nil {
    return nil, err
  }
  // return *Page and error, successful page load should reach this
  // point of the code and return nil for error
  return &Page{Title: title, Body: body}, nil
}

func main() {
  http.HandleFunc("/view/", makeHandler(viewHandler))    //handle requests under /view/ path
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))
  http.ListenAndServe(":8080", nil)
  // p1 := &Page{Title: "TestPage", Body: []byte("Welcome to The Page")}
  // p1.save()
  // p2, _ := loadPage("TestPage")
  // fmt.Println(string(p2.Body))
}


//make wrapper for handlers, return function of http.HandlerFunc type
//there will be no need to call getTitle from the handler functions once we call makeHandler in main
func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request){
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil{
      http.NotFound(w,r)
      return
    }
    fn(w,r,m[2])
  }
}


//create handler to allow users to view wiki page
func viewHandler(w http.ResponseWriter, r *http.Request, title string){
  // title := r.URL.Path[len("/view/"):]
  // title, err := getTitle(w,r)
  // if err != nil {
  //   return
  // }

  p, err := loadPage(title)               //currently dropping error from loadpage()
  // fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)  //use without templates
  if err != nil {
    http.Redirect(w,r, "/edit/"+title, http.StatusFound)      //redirect to edit page if no view page exists
    return
  }
  renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
  // title := r.URL.Path[len("/edit/"):]
  // title, err := getTitle(w,r)
  // if err != nil {
  //   return
  // }

  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }

  //create form with hard-coded HTML
  // fmt.Fprint(w, "<h1>Editing %s<h1>" +
  //   "<form action=\"/save/%s\" method=\"POST\">"+
  //   "<textarea name=\"body\">%s</textarea><br>"+
  //   "<input type=\"submit\" value=\"Save\">"+
  //   "</form>",
  //   p.Title, p.Title, p.Body)

  //create form relying on a separate html template
  renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string){
  // title := r.URL.Path[len("/save/"):]   //parses out the "/save/" from the pathname

  // title, err := getTitle(w,r)
  // if err != nil {
  //   return
  // }

  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}    //convert body string to []byte
  err := p.save()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w,r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){

  //use without templates global var
  // t, err := template.ParseFiles(tmpl + ".html")
  // if err != nil {
  //   http.Error(w, err.Error(), http.StatusInternalServerError)
  //   return
  // }
  // err = t.Execute(w,p)
  // if err != nil {
  //   http.Error(w, err.Error(), http.StatusInternalServerError)
  // }

  err := templates.ExecuteTemplate(w, tmpl+".html", p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

// run with:
// $ go build wiki.go
// $ ./wiki
