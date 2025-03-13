package main

import (
  "fmt"
  "net/http"
  "html/template"
  "forum/golang"
  "regexp"
)

var (
  errConnexion string = ""
  errInscription string = ""
  connexionDisplay ConnexionDisplay
  inscriptionDisplay InscriptionDisplay
)

type ConnexionDisplay struct {
  ErrConnexionMessage string
}

type InscriptionDisplay struct {
  ErrInscriptionMessage string
}

//* Fonction qui gère la page d'accueil
func indexHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("html/index.html"))
  if r.Method == http.MethodGet {
  }else{
    err := r.ParseForm()
    if err != nil {
      http.Error(w, "Error parsing form.", http.StatusBadRequest)
      return
    }

    title := r.FormValue("title")
    content := r.FormValue("content")

    fmt.Println(title, content)

    var postSend = golang.Post{}
    postSend.Title = title
    postSend.Text = content

    //* Écris dans la base de données si le titre ou le contenu n'est pas vide
    if title != "" || content != "" {
    fmt.Println("AddPostInDataBase database...")
    golang.AddPostInDataBase(postSend)
    fmt.Println("AddPostInDataBase ended.")
    }

  }
  tmpl.Execute(w, nil)
}

//* Fonction qui gère la page de connexion
func connexionHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("html/connexion.html"))
  if r.Method == http.MethodGet {
  }
  if r.Method == http.MethodPost {
    err := r.ParseForm()
    if err != nil {
      http.Error(w, "Error parsing form.", http.StatusBadRequest)
      return
    }
    nameOrMail := r.FormValue("nameOrMail")
    password := r.FormValue("password")
    
    //! Erreur si le nom ou le mail ou le mot de passe est vide
    if nameOrMail == "" || password == "" {
      errConnexion += "Nom ou Email vide.\n"
      return
    }

    //* Vérifie si le nom ou le mail et le mot de passe sont corrects
    if golang.CheckUser(nameOrMail, password) {
      user := http.Cookie{
        Name: "User", Value: nameOrMail,
      }
      http.SetCookie(w, &user)
    } else {
      errConnexion = "Nom ou Email incorrect.\n"
      fmt.Println(errConnexion)
    }
    connexionDisplay.ErrConnexionMessage += errConnexion
  }
  tmpl.Execute(w, connexionDisplay)
}

func inscriptionHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("html/inscription.html"))
  if r.Method == http.MethodGet {
  }
  if r.Method == http.MethodPost {
    err := r.ParseForm()
    if err != nil {
      http.Error(w, "Error parsing form.", http.StatusBadRequest)
      return
    }
    username := r.FormValue("username")
    email := r.FormValue("email")
    password := r.FormValue("password")
    fmt.Println(username, email, password)
    
    //* Vérifie si le nom d'utilisateur possède un format correct
    matchUsername, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]{5,}$`, username)
    if !matchUsername {
      errInscription += "Nom incorrect."
      fmt.Println(errInscription)
      return
    }
    
    //* Vérifie si le mail possède un format correct
    matchedEmail, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
    if !matchedEmail {
      errInscription += "Email incorrect."
      fmt.Println(errInscription)
    }
    
    //* Vérifie si le mot de passe possède un format correct
    matchedPassword, _ := regexp.MatchString(`^[a-zA-Z0-9.?_%+-]{8,}$`, password)
    if !matchedPassword {
      errInscription = "Mot de passe incorrect."
      fmt.Println(errInscription)
    }
    
    //* Vérifie si le nom d'utilisateur ou le mail est déjà utilisé
    if golang.CheckUser(username, password) == false && errInscription == ""{
      var userSend = golang.User{}
      userSend.Username = username
      userSend.Email = email
      userSend.Password = password

      //* Écris dans la base de données le User
      fmt.Println("Starting AddUserInDataBase...")
      golang.AddUserInDataBase(userSend)
      fmt.Println("AddUserInDataBase ended.")
    } else if golang.CheckUser(username, password) && errInscription != "" {
      errInscription = "Nom ou Email déjà utilisé.\n"
      fmt.Println(errInscription)
    }
    inscriptionDisplay.ErrInscriptionMessage = errInscription
    
  }
  tmpl.Execute(w, inscriptionDisplay)   
}

func main() {
  http.HandleFunc("/", indexHandler)
  http.HandleFunc("/connexion", connexionHandler)
  http.HandleFunc("/inscription", inscriptionHandler)
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	fmt.Println("http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}