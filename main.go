
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
    SimplePage(w, r, "example")
  })



  mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        SimplePage(w, r, "example")
    } else if r.Method == "POST" {
        LoginPost(w, r, db)
    }
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



func Login(w http.ResponseWriter, req *http.Request, template string) {

    r := render.New(render.Options{})
    r.HTML(w, http.StatusOK, template, nil)

}




func LoginPost(w http.ResponseWriter, req *http.Request, db mysql.Conn) {


    username := req.FormValue("username")
    SQL := "SELECT username, password FROM users WHERE username = " + username
    rows, _, err := db.Query(SQL)

    errHandler(err)

    for _, row := range rows {
        log.Print(row.Str(1))
    }


    r := render.New(render.Options{})
    r.HTML(w, http.StatusOK, "example", nil)

}












