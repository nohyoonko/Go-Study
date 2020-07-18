package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/urfave/negroni"

	"github.com/gorilla/pat"
	"github.com/unrolled/render"
)

type User struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

var rd *render.Render

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := User{Name: "tucker", Email: "tucker@naver.com"}

	rd.JSON(w, http.StatusOK, user)
	/*
		rd.JSON(w, http.StatusOK, user)
					==
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data, _ := json.Marshal(user)
		fmt.Fprint(w, string(data))
	*/
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		rd.Text(w, http.StatusBadRequest, err.Error())
		/*
			rd.Text(w, http.StatusBadRequest, err.Error())
						==
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
		*/
		return
	}
	user.CreatedAt = time.Now()
	rd.JSON(w, http.StatusOK, user)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	user := User{Name: "Tucker", Email: "Tucker@naver.com"}
	rd.HTML(w, http.StatusOK, "body", user)
	/*
		rd.HTML(w, http.StatusOk, "hello.tmpl", "Tucker")
					==
		tmpl, err := template.New("Hello").ParseFiles("templates/hello.tmpl")
		if err != nil {
			rd.Text(w, http.StatusBadRequest, err.Error())
			return
		}
		tmpl.ExecuteTemplate(w, "hello.tmpl", "Tucker")
	*/
}

func main() {
	rd = render.New(render.Options{
		Directory:  "template",
		Extensions: []string{".html", ".tmpl"},
		Layout:     "hello",
	})
	mux := pat.New()

	mux.Get("/users", getUserInfoHandler)
	mux.Post("/users", addUserHandler)
	mux.Get("/hello", helloHandler)

	n := negroni.Classic()
	n.UseHandler(mux)
	http.ListenAndServe(":3000", n)
}
