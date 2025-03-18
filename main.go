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
  Posts []golang.Post
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

type PostDisplay struct {
  Post golang.Post
}

type ListTopicDisplay struct {
  Topics []golang.Topic
  ErrListTopicMessage string
  MessageAddTopic string
}

type TopicDisplay struct {
  Topic golang.Topic
}

//!-----------------------------------------------------------------------------------------

//* Fonction qui gère la page d'accueil
func indexHandler(w http.ResponseWriter, r *http.Request) {
  indexDisplay := IndexDisplay{}
  tmpl := template.Must(template.ParseFiles("html/index.html"))

  if r.Method == http.MethodGet {
    //* Récupère les posts de la base de données et Formate les dates pour la lisibilité
    posts := golang.GetAllPosts()
    for i := range posts {
      posts[i].TotalUp, posts[i].TotalDown = golang.Totals(posts[i].ID)
    }
    indexDisplay.Posts = posts

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
    postID := r.FormValue("postId")
    vote := r.FormValue("voteType")


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
      return
    }

    //* Vérifie si l'utilisateur veut voter
    if vote != "" {
      postId, err := strconv.Atoi(postID)
      if err != nil {
        fmt.Println("Error converting postID to int:", err)
        return
      }
      
      userID, err := strconv.Atoi(userCookie.Value)
      if err != nil {
        fmt.Println("Error converting userID to int:", err)
        return
      }
      golang.Votes(uint(postId), uint(userID), vote)
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
    
    //* Récupère les posts pour mettre à jour les votes et formatte les dates
    posts := golang.GetAllPosts()
    for i := range posts {
      posts[i].TotalUp, posts[i].TotalDown = golang.Totals(posts[i].ID)
    }
    indexDisplay.Posts = posts
    
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

//* Fonction qui gère la page d'un post
func postHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("html/post.html"))
  postDisplay := PostDisplay{}

  if r.Method == http.MethodGet {
        //* Récupère le postId dans l'URL
        postIdStr := r.URL.Query().Get("postId")
        if postIdStr == "" {
            http.Error(w, "Missing postId query parameter", http.StatusBadRequest)
            return
        }
    
        postId, err := strconv.Atoi(postIdStr)
        if err != nil {
            http.Error(w, "Invalid postId query parameter", http.StatusBadRequest)
            return
        }
    
        //* Va chercher le post dans la base de données
        post := golang.GetPostByPostID(postId)
        
        post.TotalUp, post.TotalDown = golang.Totals(post.ID)

        postDisplay.Post = post
        tmpl.Execute(w, postDisplay)
  }

  if r.Method == http.MethodPost {
    tmpl.Execute(w, postDisplay)
  }
}

func listTopicsHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("html/listTopics.html"))
  var topics ListTopicDisplay

  if r.Method == http.MethodGet {
    topics.Topics = golang.GetAllTopics()
    tmpl.Execute(w, topics)
  }

  if r.Method == http.MethodPost {
    err := r.ParseForm()
    if err != nil {
      http.Error(w, "Error parsing form.", http.StatusBadRequest)
      return
    }

    nameTopic := r.FormValue("nameTopic")
    description := r.FormValue("description")

    //* Vérifie si le nom du topic ou la description est vide
    if nameTopic == "" || description == "" {
      topics.ErrListTopicMessage = "Nom du topic ou description vide."
      tmpl.Execute(w, topics)
      return
    } else {


      //* Récupère l'utilisateur
      userCookie, err := r.Cookie("UserID")
      if err != nil {
        return
      }
      userID, err := strconv.Atoi(userCookie.Value)
      if err != nil {
        return
      }
      user := golang.GetUserByID(userID)

      //* Ajoute le topic dans la base de données
      result := golang.AddTopic(nameTopic, description, user)
      if result != "Topic created" {
        topics.ErrListTopicMessage = result
        tmpl.Execute(w, topics)
        return
      } else {
        topics.MessageAddTopic = "Topic créé."
      }
    }
    //* Récupère les topics
    topics.Topics = golang.GetAllTopics()

    tmpl.Execute(w, topics)
  }

}

func topicHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("html/topic.html"))
  topic := TopicDisplay{}
  if r.Method == http.MethodGet {

    //* Récupère le topicID dans l'URL
    topicIDStr := r.URL.Query().Get("topicId")
    if topicIDStr == "" {
      return
    }
    topicID, err := strconv.Atoi(topicIDStr)
    if err != nil {
      return
    }
    topic.Topic = golang.GetTopic(topicID)
    
    tmpl.Execute(w, topic)
  }
  if r.Method == http.MethodPost {


    tmpl.Execute(w, nil)
  }
}
func main() {

  golang.CreateAdminUser()
  
  http.HandleFunc("/", indexHandler)
  http.HandleFunc("/connexion", connexionHandler)
  http.HandleFunc("/inscription", inscriptionHandler)
  http.HandleFunc("/pageUtilisateur", pageUtilisateurHandler)
  http.HandleFunc("/post", postHandler)
  http.HandleFunc("/listTopics", listTopicsHandler)
  http.HandleFunc("/topic", topicHandler)
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	fmt.Println("http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}