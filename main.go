package main

import (
  "fmt"
  "net/http"
  "html/template"
  "forum/golang"
  "regexp"
  "strconv"
  "golang.org/x/crypto/bcrypt"
)

var (
  errConnexion string = ""
  errInscription string = ""
  connexionDisplay ConnexionDisplay
  inscriptionDisplay InscriptionDisplay
)

type IndexDisplay struct {
  ErrIndexMessage string
}
type ConnexionDisplay struct {
  ErrConnexionMessage string
}

type InscriptionDisplay struct {
  ErrInscriptionMessage string
}

type PageUtilisateurDisplay struct {
  Posts []golang.Post
  Username string
}

//* Fonction qui gère la page d'accueil
func indexHandler(w http.ResponseWriter, r *http.Request) {
  indexDisplay := IndexDisplay{}
  tmpl := template.Must(template.ParseFiles("html/index.html"))

  if r.Method == http.MethodGet {
    tmpl.Execute(w, indexDisplay)
    
  } 
  if r.Method == http.MethodPost {
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

    //* Réccupère l'ID de l'utilisateur des cookies
    userCookie, err := r.Cookie("UserID")
    if err == nil {
      userID, err := strconv.Atoi(userCookie.Value)
      if err == nil {
        postSend.UserID = uint(userID)
      }
    } else {
      fmt.Println("Error retrieving user cookie:", err)
    }

    //* Vérifie si l'utilisateur est connecté
    if userCookie != nil && userCookie.Value != "" {
      //* Écris dans la base de données si le titre ou le contenu n'est pas vide
      if title != "" || content != "" {
        fmt.Println("AddPostInDataBase database...")
        golang.AddPostInDataBase(postSend)
        fmt.Println("AddPostInDataBase ended.")
      }
    } else {
      errIndexMessage := "Connectez vous pour poster !"
      indexDisplay.ErrIndexMessage = errIndexMessage
      fmt.Println(errIndexMessage)
    }
    
    tmpl.Execute(w, indexDisplay)
  }
}
  
//* Fonction qui gère la page de connexion
func connexionHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("html/connexion.html"))
  if r.Method == http.MethodGet {
    tmpl.Execute(w, connexionDisplay)
    return
  }
  if r.Method == http.MethodPost {
    err := r.ParseForm()
    if err != nil {
      http.Error(w, "Error parsing form.", http.StatusBadRequest)
      return
    }

    //* Vérifie si l'utilisateur veut se déconnecter
    logout := r.FormValue("logout")
    fmt.Println("Logout:", logout)
    if logout == "true" {
      //* Vide le cookie de l'utilisateur
      userCookie := http.Cookie{
        Name:   "UserID",
        Value:  "",
        Path:   "/",
        MaxAge: -1, //? Durée de vie du cookie
      }
      http.SetCookie(w, &userCookie)
      http.Redirect(w, r, "/", http.StatusSeeOther)
      return
    }

    nameOrMail := r.FormValue("nameOrMail")
    password := r.FormValue("password")
    
    //! Erreur si le nom, le mail ou le mot de passe est vide
    if nameOrMail == "" || password == "" {
      errConnexion = "Nom, Email ou Mot de passe vide.\n"
      connexionDisplay.ErrConnexionMessage += errConnexion
      tmpl.Execute(w, connexionDisplay)
      return
    }

    //* Vérifie si le nom d'utilisateur ou le mail est correct
    user, valid := golang.CheckUserPassword(nameOrMail, password)
    if valid {
      //* Crée un cookie pour l'utilisateur
      userCookie := http.Cookie{
        Name:  "UserID",
        Value: strconv.Itoa(int(user.ID)),
        Path:  "/",
      }
      http.SetCookie(w, &userCookie)
      http.Redirect(w, r, "/", http.StatusSeeOther)
      return
    } else {
      errConnexion = "Nom ou Email incorrect.\n"
      fmt.Println(errConnexion)
      connexionDisplay.ErrConnexionMessage += errConnexion
    }
  }
  tmpl.Execute(w, connexionDisplay)
}

//* Fonction qui gère la page d'inscription
func inscriptionHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("html/inscription.html"))
  if r.Method == http.MethodGet {
    tmpl.Execute(w, inscriptionDisplay)
    return
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
    if golang.CheckUser(username) == false && errInscription == ""{
      var userSend = golang.User{}
      userSend.Username = username
      userSend.Email = email

      //* Hash the password
      hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
      if err != nil {
        http.Error(w, "Error hashing password.", http.StatusInternalServerError)
        return
      }
      userSend.Password = string(hashedPassword)

      //* Écris dans la base de données le User
      fmt.Println("Starting AddUserInDataBase...")
      golang.AddUserInDataBase(userSend)
      fmt.Println("AddUserInDataBase ended.")
    } else if golang.CheckUser(username) && errInscription != "" {
      errInscription = "Nom déjà utilisé.\n"
      fmt.Println(errInscription)
    }
    inscriptionDisplay.ErrInscriptionMessage = errInscription
  }
  tmpl.Execute(w, inscriptionDisplay)   
}

//* Fonction qui gère la page de l'utilisateur
func pageUtilisateurHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("html/pageUtilisateur.html"))
  var pageUtilisateurDisplay PageUtilisateurDisplay
  if r.Method == http.MethodGet {
    userCookie, err := r.Cookie("UserID")
    if err == nil {
      userID, _ := strconv.Atoi(userCookie.Value)
      fmt.Println("User ID:", userID)
  
      //* Récupère les posts de l'utilisateur
      posts := golang.GetPostsByUserID(userID)
      pageUtilisateurDisplay = PageUtilisateurDisplay{Posts: posts}

      user := golang.GetUserByID(userID)
      if user.Username != "" {
        pageUtilisateurDisplay.Username = user.Username
      } else {
        http.Error(w, "Connectez vous pour voir vos Post !", http.StatusUnauthorized)
      }
    }
    tmpl.Execute(w, pageUtilisateurDisplay)
  }
  if r.Method == http.MethodPost {
    r.ParseForm()
    deletePostID, err := strconv.Atoi(r.FormValue("deletePost"))
    if err != nil {
        return
    }
    golang.DeletePost(deletePostID)
    http.Redirect(w, r, "/pageUtilisateur", http.StatusSeeOther)
    tmpl.Execute(w, pageUtilisateurDisplay)
  }
}

func main() {

  golang.CreateAdminUser()
  
  http.HandleFunc("/", indexHandler)
  http.HandleFunc("/connexion", connexionHandler)
  http.HandleFunc("/inscription", inscriptionHandler)
  http.HandleFunc("/pageUtilisateur", pageUtilisateurHandler)
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	fmt.Println("http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}