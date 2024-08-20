package controllers

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"libredrive/models"
	"libredrive/types"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAdmin").(bool) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	userParams := models.CreateUserParams{}
	r.ParseForm()
	if r.Form.Get("Username") == "" || r.Form.Get("Password") == "" ||
		len(r.Form.Get("Password")) > 72 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userParams.Username = r.Form.Get("Username")
	userParams.Isadmin = r.Form.Get("IsAdmin") == "on"
	password, _ := bcrypt.GenerateFromPassword([]byte(r.Form.Get("Password")), 14)
	userParams.Password = string(password)

	if user, err := types.Queries.CreateUser(types.CTX, userParams); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	} else {
		os.MkdirAll(path.Join("user_data", strconv.Itoa(int(user.ID))), 0750)
		w.Header().Set("HX-Redirect", "/")
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Username := r.Form.Get("Username")
	Password := r.Form.Get("Password")

	user, err := types.Queries.GetUser(types.CTX, Username)
	if err != nil ||
		bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Password)) != nil {
		http.Error(w, "Incorrect Username or Password", http.StatusForbidden)
		return
	}
	c := http.Cookie{
		Name:   "auth",
		Value:  fmt.Sprintf("%d&%t&%x", user.ID, user.Isadmin, sha256.Sum256([]byte(Password))),
		MaxAge: 1800,
		Path:   "/",
	}
	http.SetCookie(w, &c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
