
package main


import (
    "net/http"
    "github.com/codegangsta/negroni"
    "github.com/unrolled/render"
    "log"
    //"database/sql"
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

  mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    Splash(w, r)
  })

  mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
    Splash(w, r, db)
  })

  mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
    Splash(w, r, db)
  })

  mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
    Splash(w, r, db)
  })

  mux.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
    Splash(w, r, db)
  })



  n.UseHandler(mux)
  n.Run(":3000")

}


func errHandler(err error) {
    if err != nil {
        log.Print(err)
    }
}





func Splash(w http.ResponseWriter, req *http.Request) {

    r := render.New(render.Options{})
    r.HTML(w, http.StatusOK, "example", nil)

}




func Login(w http.ResponseWriter, req *http.Request, db mysql.Conn) {


    rows, _, err := db.Query("SELECT * FROM users")

    errHandler(err)

    for _, row := range rows {
        log.Print(row.Str(1))
    }


    r := render.New(render.Options{})
    r.HTML(w, http.StatusOK, "example", nil)

}












