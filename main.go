package main

import (
	"fmt"
	"forum/golang"
	"html/template"
	"net/http"
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
  SearchUsers []golang.User
  SearchPosts []golang.Post
  SearchTopics []golang.Topic
  User golang.User
  ErrSearchMessage string
  IsSearch bool
}

type ConnexionDisplay struct {
  ErrConnexionMessage string
}

type InscriptionDisplay struct {
  ErrInscriptionMessage string
}

type ProfilDisplay struct {
  Posts []golang.Post
  User golang.User
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
  Posts []golang.Post
  User golang.User
  ErrTopicMessage string
}

type ActifUserDisplay struct {
  Users []golang.User
}

type AdminDisplay struct {
  Users []golang.User
  Posts []golang.Post
  Topics []golang.Topic
  Comments []golang.Comment
  UserConnected golang.User
  ErrAdminMessage string
}

//!-----------------------------------------------------------------------------------------

//* Fonction qui gère la page d'accueil
func indexHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("html/index.html"))
  var indexDisplay IndexDisplay

  if r.Method == http.MethodGet {
    //* Récupère le cookie de l'utilisateur
    userCookie, err := r.Cookie("UserID")
    if err == nil {
      userID, err := strconv.Atoi(userCookie.Value)
      if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
      }
      indexDisplay.User = golang.GetUserByID(userID)
    } else {
      indexDisplay.User = golang.User{}
    }
  }
  
  if r.Method == http.MethodPost {
    indexDisplay.IsSearch = false
    err := r.ParseForm()
    if err != nil {
      http.Error(w, "Error parsing form.", http.StatusBadRequest)
      return
    }
    
    //* Info de la recherche
    search := r.FormValue("search")
    
    //* Recherche dans la base de données 
    if search != "" {
      //* Récupère les utilisateurs, posts et topics correspondant à la recherche
      searchUsers, searchPosts, searchTopics := golang.SearchUserPostTopic(search)
      
      if len(searchUsers) == 0 && len(searchPosts) == 0 && len(searchTopics) == 0 {
        indexDisplay.ErrSearchMessage = "Aucun résultat trouvé."
        
      } else {
        if len(searchUsers) > 0 {
          for i := range searchUsers {
            //* Fait le total des posts de l'utilisateur
            posts := golang.GetPostsByUserID(int(searchUsers[i].ID))
            if len(posts) > 0 {
              searchUsers[i].TotalPost = uint(len(posts))
            } else {
              searchUsers[i].TotalPost = 0
            }

            //* Fait le total des commentaires de l'utilisateur
            comments := golang.GetCommentsByUserID(int(searchUsers[i].ID))
            if len(comments) > 0 {
              searchUsers[i].TotalComment = uint(len(comments))
            } else {
              searchUsers[i].TotalComment = 0
            }
      
            //* Fait le total des votes de l'utilisateur
            searchUsers[i].TotalVote = golang.TotalVotes(searchUsers[i].ID)
          }
        }
        //* Met dans la template les utilisateurs, posts et topics trouvés
        indexDisplay = IndexDisplay{
          SearchUsers: searchUsers,
          SearchPosts: searchPosts,
          SearchTopics: searchTopics,
          IsSearch: true,
          ErrSearchMessage: "",
        }
      }
    }
  }
  
  //! Exécute le template
  tmpl.Execute(w, indexDisplay)
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

      connexionDisplay.ErrConnexionMessage += errConnexion
      http.Redirect(w, r, "/connexion", http.StatusSeeOther)
      return
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
    }
    
    //* Vérifie si le mail possède un format correct
    matchedEmail, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
    if !matchedEmail {
      errInscription += "Email incorrect."
      fmt.Println(errInscription)
    }
    
    //* Vérifie si le mot de passe possède un format correct
    matchedPassword, _ := regexp.MatchString(`^[a-zA-Z0-9.?!_%+-]{8,}$`, password)
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
func profilHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("html/profil.html"))
  var profilDisplay ProfilDisplay

  if r.Method == http.MethodGet {
    //* Récupère l'ID de l'utilisateur des cookies
    userCookie, err := r.Cookie("UserID")
    if err == nil {
      userID, _ := strconv.Atoi(userCookie.Value)
      // fmt.Println("User ID:", userID)
  
      //* Récupère les posts de l'utilisateur
      posts := golang.GetPostsByUserID(userID)
      profilDisplay = ProfilDisplay{Posts: posts}

      user := golang.GetUserByID(userID)
      if user.Username == "" {
        profilDisplay.User = golang.User{}
      } else {
        profilDisplay.User = user
        if len(posts) > 0 {
          user.TotalPost = uint(len(posts))
        } else {
          user.TotalPost = 0
        }
        //* Fait le total des votes de l'utilisateur
        user.TotalVote = golang.TotalVotes(user.ID)
      }
      profilDisplay.User = user
    }
    tmpl.Execute(w, profilDisplay)
  }

  if r.Method == http.MethodPost {
    r.ParseForm()
    deletePostIDStr := r.FormValue("deletePost")
    deletePostID, err := strconv.Atoi(deletePostIDStr)
    if err != nil {
        return
    }

    // fmt.Println("Delete post ID:", deletePostID)
    golang.DeletePost(deletePostID)
    
    http.Redirect(w, r, "/profil", http.StatusSeeOther)
    tmpl.Execute(w, profilDisplay)
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

//* Fonction qui gère la page des topics
func listTopicsHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("html/listTopics.html"))
  var topicsDisplay ListTopicDisplay

  if r.Method == http.MethodGet {
    topicsDisplay.Topics = golang.GetAllTopics()
    tmpl.Execute(w, topicsDisplay)
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
      topicsDisplay.ErrListTopicMessage = "Nom du topic ou description vide."

      //* Récupère les topics
      topicsDisplay.Topics = golang.GetAllTopics()
      tmpl.Execute(w, topicsDisplay)
      return
    } else {

      //* Récupère l'utilisateur
      userCookie, err := r.Cookie("UserID")
      if err != nil {
        topicsDisplay.ErrListTopicMessage = "Connectez vous pour créer un topic !"
        //* Récupère les topics ARRETE DE l'OUBLIER
        topicsDisplay.Topics = golang.GetAllTopics()
        tmpl.Execute(w, topicsDisplay)
        return
      }
      userID, err := strconv.Atoi(userCookie.Value)
      if err != nil {
        topicsDisplay.ErrListTopicMessage = "Connectez vous pour créer un topic !"
        //* Récupère les topics ARRETE DE l'OUBLIER
        topicsDisplay.Topics = golang.GetAllTopics()
        tmpl.Execute(w, topicsDisplay)
        return
      }
      user := golang.GetUserByID(userID)

      //* Vérifie si l'utilisateur est connecté
      if user.Username == "" {
        topicsDisplay.ErrListTopicMessage = "Connectez vous pour créer un topic !"

        //* Récupère les topics
        topicsDisplay.Topics = golang.GetAllTopics()
        tmpl.Execute(w, topicsDisplay)
        return
      } else {
        //* Ajoute le topic dans la base de données
        result := golang.AddTopic(nameTopic, description, user)
        if result != "Topic created" {
          topicsDisplay.ErrListTopicMessage = result
          tmpl.Execute(w, topicsDisplay)
          return
        } else {
          topicsDisplay.MessageAddTopic = "Topic créé."
        }
      }
    }
    //* Récupère les topics
    topicsDisplay.Topics = golang.GetAllTopics()

    tmpl.Execute(w, topicsDisplay)
  }

}

//* Fonction qui gère la page d'un topic
func topicHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("html/topic.html"))
  topicDisplay := TopicDisplay{}
  
  if r.Method == http.MethodGet {
    //* Récupère le topicID dans l'URL
    topicIDStr := r.URL.Query().Get("topicId")
    if topicIDStr == "" {
      http.Error(w, "Missing topicId query parameter", http.StatusBadRequest)
      return
    }
    topicID, err := strconv.Atoi(topicIDStr)
    if err != nil {
      http.Error(w, "Invalid topicId query parameter", http.StatusBadRequest)
      return
    }
    topicDisplay.Topic = golang.GetTopic(topicID)

    //* Récupère les posts pour mettre à jour les votes et formatte les dates
    posts := golang.GetPostsByTopicID(topicID)
    for i := range posts {
      posts[i].TotalUp, posts[i].TotalDown = golang.Totals(posts[i].ID)
    }
    topicDisplay.Posts = posts

    //* Récupère l'ID de l'utilisateur des cookies
    userCookie, err := r.Cookie("UserID")

    if err != nil {
      topicDisplay.ErrTopicMessage = "Connectez vous pour poster (PAS CONNECTÉ)!"
      for i := range topicDisplay.Posts {
      topicDisplay.Posts[i].IsLoggedIn = false
      topicDisplay.Posts[i].UserConnectedID = 0
      }
    } else {
      userID, err := strconv.Atoi(userCookie.Value)
      if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
      }
      topicDisplay.User = golang.GetUserByID(userID)

      //* Sert à afficher les boutons
      for i := range topicDisplay.Posts {
        topicDisplay.Posts[i].IsLoggedIn = true
        topicDisplay.Posts[i].UserConnectedID = uint(userID)
        for j := range topicDisplay.Posts[i].Comments {
          topicDisplay.Posts[i].Comments[j].UserConnectedID = uint(userID)
          topicDisplay.Posts[i].Comments[j].IsLoggedIn = true
        }
      }
    }

    //! Exécute le template
    err = tmpl.Execute(w, topicDisplay)
    if err != nil {
      http.Error(w, "Error executing template", http.StatusInternalServerError)
    }
    return
  }

  if r.Method == http.MethodPost {
    err := r.ParseForm()
    if err != nil {
      http.Error(w, "Error parsing form.", http.StatusBadRequest)
      return
    }

    //* Info du post
    title := r.FormValue("title")
    content := r.FormValue("content")
    postID := r.FormValue("postId")
    vote := r.FormValue("voteType")
    
    
    //* Info de la modification
    modifyPostIDStr := r.FormValue("modifyPostId")
    modifyTitle := r.FormValue("modifyTitle")
    modifyContent := r.FormValue("modifyContent")
    
    //* Info de la suppression
    deletePostIDStr := r.FormValue("deletePostId")
    deleteCommentIDStr := r.FormValue("deleteCommentId")

    //* Info du commentaire
    commentPostIDStr := r.FormValue("commentPostId")
    commentContent := r.FormValue("commentContent")
    
    //* Récupère le topicID dans l'URL
    topicIDStr := r.FormValue("topicId")
    topicID, err := strconv.Atoi(topicIDStr)
    if err != nil {
      http.Error(w, "Invalid topicId", http.StatusBadRequest)
      return
    }

    //* Récupère les posts pour mettre à jour les votes et formatte les dates
    posts := golang.GetPostsByTopicID(topicID)
    for i := range posts {
      posts[i].TotalUp, posts[i].TotalDown = golang.Totals(posts[i].ID)
    }
    topicDisplay.Posts = posts

    //* Vérifie que l'utilisateur est connecté
    userCookie, err := r.Cookie("UserID")

    if err != nil {
      topicDisplay.ErrTopicMessage = "Connectez vous pour poster (PAS CONNECTÉ)!"
      //* Fait en sorte que les boutons ne s'affichent pas
      for i := range topicDisplay.Posts {
        topicDisplay.Posts[i].IsLoggedIn = false
        topicDisplay.Posts[i].UserConnectedID = 0
      }

    } else {
    
      //* Création de l'ID de l'utilisateur
      userID, err := strconv.Atoi(userCookie.Value)
      if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
      }
      topicDisplay.User = golang.GetUserByID(userID)

      //* Sert à afficher les boutons
      for i := range topicDisplay.Posts {
        topicDisplay.Posts[i].IsLoggedIn = true
        topicDisplay.Posts[i].UserConnectedID = uint(userID)
        for j := range topicDisplay.Posts[i].Comments {
          topicDisplay.Posts[i].Comments[j].UserConnectedID = uint(userID)
          topicDisplay.Posts[i].Comments[j].IsLoggedIn = true
        }
      }
    
      //* Vérifie si l'utilisateur veux poster
      if vote == "" && title != "" && content != "" {

        //* Crée un post
        var postSend = golang.Post{}
        postSend.TopicID = uint(topicID)
        postSend.Title = title
        postSend.Text = content
        postSend.UserID = uint(userID)

        golang.AddPostInDataBase(postSend)
        http.Redirect(w, r, "/topic?topicId="+topicIDStr, http.StatusSeeOther)
        return
      } else {
        topicDisplay.ErrTopicMessage = "Connectez vous pour poster ! (VEUX POSTER)"
      }

      //* Vérifie si l'utilisateur veut voter
      if (postID != "" && vote != "") && (modifyContent == content && modifyTitle == title) {
        postId, err := strconv.Atoi(postID)
        if err != nil {
          http.Error(w, "Invalid post ID", http.StatusBadRequest)
          return
        }
        //* Récupère l'ID de l'utilisateur des cookies
        userID, err := strconv.Atoi(userCookie.Value)
        if err != nil {
          fmt.Println("Error getting user ID")
          http.Error(w, "Invalid user ID", http.StatusBadRequest)
          return
        }
        //* Ajoute le vote dans la base de données
        golang.Votes(uint(postId), uint(userID), vote)

        //* Redirige l'utilisateur après avoir traité le vote
        http.Redirect(w, r, "/topic?topicId="+topicIDStr, http.StatusSeeOther)
        return
      } else {
        topicDisplay.ErrTopicMessage = "Connectez vous pour voter ! (VEUX VOTER)"
      }

      //* Modifier le Post
      if modifyPostIDStr != "" && modifyTitle != title && modifyContent != content && vote == "" {
        modifyPostID, err := strconv.Atoi(modifyPostIDStr)
        if err != nil {
          http.Error(w, "Invalid post ID", http.StatusBadRequest)
          return
        }

        //* Crée un post
        var postSend = golang.Post{}
        postSend.ID = uint(modifyPostID)
        postSend.Title = modifyTitle
        postSend.Text = modifyContent

        golang.UpdatePost(postSend)

        http.Redirect(w, r, "/topic?topicId="+topicIDStr, http.StatusSeeOther)
        return
      } else {
        topicDisplay.ErrTopicMessage = "Connectez le titre ou le contenue pour modifier ! (VEUX MODIFIER)"
      }

      //* Supprimer le Post
      if deletePostIDStr != "" {
        deletePostID, err := strconv.Atoi(deletePostIDStr)
        if err != nil {
          http.Error(w, "Invalid post ID", http.StatusBadRequest)
          return
        }

        // fmt.Println("Delete post ID:", deletePostID)
        golang.DeletePost(deletePostID)

        http.Redirect(w, r, "/topic?topicId="+topicIDStr, http.StatusSeeOther)
        return
      }

      //* Vérifie si l'utilisateur veut commenter
      if commentContent != "" && commentPostIDStr != "" {
        //* Récupère l'ID du post
        commentPostID, err := strconv.Atoi(commentPostIDStr)
        if err != nil {
          http.Error(w, "Invalid post ID", http.StatusBadRequest)
          return
        }

        //* Crée un commentaire
        var commentSend = golang.Comment{}
        commentSend.PostID = uint(commentPostID)
        commentSend.Text = commentContent
        commentSend.UserID = uint(userID)

        golang.AddComment(commentSend)

        http.Redirect(w, r, "/topic?topicId="+topicIDStr, http.StatusSeeOther)
        return
      } else {
        topicDisplay.ErrTopicMessage = "Connectez vous pour commenter ! (VEUX COMMENTER)"
      }

      //* Supprimer le commentaire
      if deleteCommentIDStr != "" {
        deleteCommentID, err := strconv.Atoi(deleteCommentIDStr)
        if err != nil {
          http.Error(w, "Invalid post ID", http.StatusBadRequest)
          return
        }

        // fmt.Println("Delete comment ID:", deleteCommentID)
        golang.DeleteComment(uint(deleteCommentID))

        http.Redirect(w, r, "/topic?topicId="+topicIDStr, http.StatusSeeOther)
        return
      }
    }

    //! Exécute le template
    err = tmpl.Execute(w, topicDisplay)
    if err != nil {
      http.Error(w, "Error executing template", http.StatusInternalServerError)
    }
  }
}

//* Fonction qui gère la page des utilisateurs actifs
func actifUser(w http.ResponseWriter, r *http.Request) {
  tmp := template.Must(template.ParseFiles("html/actifUser.html"))
  if r.Method == http.MethodGet {
    users := ActifUserDisplay{Users: golang.GetAllUsers()}
    for i := range users.Users {
      //* Fait le total des posts de l'utilisateur
      posts := golang.GetPostsByUserID(int(users.Users[i].ID))
      if len(posts) > 0 {
        users.Users[i].TotalPost = uint(len(posts))
      } else {
        users.Users[i].TotalPost = 0
      }
      //* Fait le total des commentaires de l'utilisateur
      comments := golang.GetCommentsByUserID(int(users.Users[i].ID))
      if len(comments) > 0 {
        users.Users[i].TotalComment = uint(len(comments))
      } else {
        users.Users[i].TotalComment = 0
      }
      
      //* Fait le total des votes de l'utilisateur
      users.Users[i].TotalVote = golang.TotalVotes(users.Users[i].ID)
    }
    err := tmp.Execute(w, users)
    if err != nil {
      http.Error(w, "Error executing template", http.StatusInternalServerError)
      return
    }
  }

  if r.Method == http.MethodPost {
    tmp.Execute(w, nil)
  }

}

//* Fonction qui gère la page de l'admin
func adminHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("html/admin.html"))
  adminDisplay := AdminDisplay{}
  if r.Method == http.MethodGet {
     //* Récupère le cookie de l'utilisateur
     userCookie, err := r.Cookie("UserID")
     if err == nil {
       userID, err := strconv.Atoi(userCookie.Value)
       if err != nil {
         http.Error(w, "Invalid user ID", http.StatusBadRequest)
         return
       }
       adminDisplay.UserConnected = golang.GetUserByID(userID)
     } else {
       adminDisplay.UserConnected = golang.User{}
     }

    //* Récupère les utilisateurs, posts, topics et commentaires
      adminDisplay.Users = golang.GetAllUsers()
      adminDisplay.Posts = golang.GetAllPosts()
      adminDisplay.Topics = golang.GetAllTopics()
      adminDisplay.Comments = golang.GetAllComments()
  }

  if r.Method == http.MethodPost {
    err := r.ParseForm()
    if err != nil {
      http.Error(w, "Error parsing form.", http.StatusBadRequest)
      return
    }

    //* Récupère le cookie de l'utilisateur
    userCookie, err := r.Cookie("UserID")
    if err == nil {
      userID, err := strconv.Atoi(userCookie.Value)
      if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
      }
      adminDisplay.UserConnected = golang.GetUserByID(userID)
    } else {
      adminDisplay.UserConnected = golang.User{}
    }

    //* Données de la recherche
    searchUser := r.FormValue("searchUser")
    searchPost := r.FormValue("searchPost")
    searchTopic := r.FormValue("searchTopic")
    searchComment := r.FormValue("searchComment")
    searchForum := r.FormValue("searchForum")

    //* Recherche dans la base de données

    //! varibale pour savoir si l'utilisateur a fait plusieurs recherches
    searchUserBool := 0
    searchPostBool := 0
    searchTopicBool := 0
    searchCommentBool := 0
    searchForumBool := 0

    if searchUser != "" {
      searchUserBool = 1
    }
    if searchPost != "" {
      searchPostBool = 1
    }
    if searchTopic != "" {
      searchTopicBool = 1
    }
    if searchComment != "" {
      searchCommentBool = 1
    }
    if searchForum != "" {
      searchForumBool = 1
    }

    if searchUserBool + searchPostBool + searchTopicBool + searchCommentBool + searchForumBool > 1 {
      adminDisplay.ErrAdminMessage = "Vous ne pouvez pas faire plusieurs recherches en même temps."
      adminDisplay.Users = golang.GetAllUsers()
      adminDisplay.Posts = golang.GetAllPosts()
      adminDisplay.Topics = golang.GetAllTopics()
      adminDisplay.Comments = golang.GetAllComments()
    } else {
      //* Recherche les utilisateurs
      if searchUser != "" {
        adminDisplay.Users = golang.SearchUsersByUsername(searchUser)
        adminDisplay.Posts = []golang.Post{}
        adminDisplay.Topics = []golang.Topic{}
        adminDisplay.Comments = []golang.Comment{}

        if len(adminDisplay.Users) == 0 {
          adminDisplay.ErrAdminMessage = "Aucun résultat trouvé."
        } else {
          adminDisplay.ErrAdminMessage = ""
        }     
      }
      
      //* Recherche les posts
      if searchPost != "" {
        adminDisplay.Posts = golang.SearchPostsByTitle(searchPost)
        adminDisplay.Users = []golang.User{}
        adminDisplay.Topics = []golang.Topic{}
        adminDisplay.Comments = []golang.Comment{}
        
        if len(adminDisplay.Posts) == 0 {
          adminDisplay.ErrAdminMessage = "Aucun résultat trouvé."
        } else {
          adminDisplay.ErrAdminMessage = ""
        }
      }
      
      //* Recherche les topics
      if searchTopic != "" {
        adminDisplay.Topics = golang.SearchTopicsByName(searchTopic)
        adminDisplay.Users = []golang.User{}
        adminDisplay.Posts = []golang.Post{}
        adminDisplay.Comments = []golang.Comment{}
        
        if len(adminDisplay.Topics) == 0 {
          adminDisplay.ErrAdminMessage = "Aucun résultat trouvé."
        } else {
          adminDisplay.ErrAdminMessage = ""
        }
      }
      
      //* Recherche les commentaires
      if searchComment != "" {
        adminDisplay.Comments = golang.SearchCommentsByText(searchComment)
        adminDisplay.Users = []golang.User{}
        adminDisplay.Posts = []golang.Post{}
        adminDisplay.Topics = []golang.Topic{}

        if len(adminDisplay.Comments) == 0 {
          adminDisplay.ErrAdminMessage = "Aucun résultat trouvé."
        } else {
          adminDisplay.ErrAdminMessage = ""
        }
      }

      //* Recherche dans le forum
      if searchForum != "" {
        adminDisplay.Users = golang.SearchUsersByUsername(searchForum)
        adminDisplay.Posts = golang.SearchPostsByTitle(searchForum)
        adminDisplay.Topics = golang.SearchTopicsByName(searchForum)
        adminDisplay.Comments = golang.SearchCommentsByText(searchForum)

        if len(adminDisplay.Users) == 0 && len(adminDisplay.Posts) == 0 && len(adminDisplay.Topics) == 0 && len(adminDisplay.Comments) == 0 {
          adminDisplay.ErrAdminMessage = "Aucun résultat trouvé."
        } else {
          adminDisplay.ErrAdminMessage = ""
        }
      }
    }
  }

  //! Exécute le template
  err := tmpl.Execute(w, adminDisplay)
  if err != nil {
    http.Error(w, "Error executing template", http.StatusInternalServerError)
  }
}

func main() {

  golang.CreateAdminUser()
  
  http.HandleFunc("/", indexHandler)
  http.HandleFunc("/connexion", connexionHandler)
  http.HandleFunc("/inscription", inscriptionHandler)
  http.HandleFunc("/profil", profilHandler)
  http.HandleFunc("/post", postHandler)
  http.HandleFunc("/listTopics", listTopicsHandler)
  http.HandleFunc("/topic", topicHandler)
  http.HandleFunc("/actifUser", actifUser)
  http.HandleFunc("/admin", adminHandler)

  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	fmt.Println("http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}