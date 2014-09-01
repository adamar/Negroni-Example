
package main


import (
    "net/http"
    "github.com/codegangsta/negroni"
    "github.com/unrolled/render"
    "github.com/goincremental/negroni-sessions"
    "log"
    "github.com/ziutek/mymysql/mysql"
    _ "github.com/ziutek/mymysql/native"
      )



func main() {

  db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "", "negroni")
  err := db.Connect()
     errHandler(err)
  defer db.Close()


  mux := http.NewServeMux()
  n := negroni.Classic()


  store := sessions.NewCookieStore([]byte("ohhhsooosecret")) 
  n.Use(sessions.Sessions("gloabl_session_store", store))

  mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    SimplePage(w, r, "mainpage")
  })


  mux.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
    SimpleAuthenticatedPage(w, r, "mainpage")
  })



  mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        SimplePage(w, r, "login")
    } else if r.Method == "POST" {
        LoginPost(w, r, db)
    }
  })


  mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        SimplePage(w, r, "signup")
    } else if r.Method == "POST" {
        SignupPost(w, r, db)
    }
  })


  mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
    SimplePage(w, r, "logout")
  })




  n.UseHandler(mux)
  n.Run(":3000")

}


func errHandler(err error) {
    if err != nil {
        log.Print(err)
    }
}



func SimplePage(w http.ResponseWriter, req *http.Request, template string) {

    r := render.New(render.Options{})
    r.HTML(w, http.StatusOK, template, nil)

}


func SimpleAuthenticatedPage(w http.ResponseWriter, req *http.Request, template string) {

    session := sessions.GetSession(req)
    session.Set("hello", "world")

    sess := session.Get("hello")

    if sess == nil {
        http.Redirect(w, req, "/notauthenticated", 301)
    }

    r := render.New(render.Options{})
    r.HTML(w, http.StatusOK, template, nil)

}




func LoginPost(w http.ResponseWriter, req *http.Request, db mysql.Conn) {

    session := sessions.GetSession(req)
    session.Set("hello", "world")

    username := req.FormValue("username")
    password := req.FormValue("password")

    rows, _, err := db.Query("SELECT * FROM users WHERE password = '?' AND username = '?'", username, password)


    errHandler(err)

    for _, row := range rows {
        log.Print(row.Str(1))
    }


    r := render.New(render.Options{})
    r.HTML(w, http.StatusOK, "example", nil)

}



func SignupPost(w http.ResponseWriter, req *http.Request, db mysql.Conn) {


    username := req.FormValue("username")
    password := req.FormValue("password")
    email := req.FormValue("email")

    //SQL := "INSERT INTO users (username, password, email) VAUES ('?', '?', '?')", username, email, password
    rows, _, err := db.Query("INSERT INTO users (username, password, email) VAUES ('?', '?', '?')", username, email, password)

    errHandler(err)

    for _, row := range rows {
        log.Print(row.Str(1))
    }


    r := render.New(render.Options{})
    r.HTML(w, http.StatusOK, "example", nil)

}



