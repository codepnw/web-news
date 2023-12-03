package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/CloudyKit/jet/v6"
	"github.com/codepnw/web-news/models"
	"github.com/codepnw/web-news/utils"
	"github.com/go-chi/chi/v5"
)

func (a *Application) homeHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		a.serverError(w, err)
		return
	}

	filter := models.Filter{
		Query:    r.URL.Query().Get("q"),
		Page:     a.readIntDefault(r, "page", 1),
		PageSize: a.readIntDefault(r, "page_size", 5),
		OrderBy:  r.URL.Query().Get("order_by"),
	}

	posts, meta, err := a.Models.Posts.GetAll(filter)
	if err != nil {
		a.serverError(w, err)
		return
	}

	queryUrl := fmt.Sprintf("page_size=%d&order_by=%s&q=%s", meta.PageSize, filter.OrderBy, filter.Query)
	nextUrl := fmt.Sprintf("%s&page=%d", queryUrl, meta.NextPage)
	prevUrl := fmt.Sprintf("%s&page=%d", queryUrl, meta.PrevPage)

	vars := make(jet.VarMap)
	vars.Set("posts", posts)
	vars.Set("meta", meta)
	vars.Set("nextUrl", nextUrl)
	vars.Set("prevUrl", prevUrl)
	vars.Set("form", utils.NewForm(r.Form))

	err = a.Render(w, r, "index", vars)

	if err != nil {
		log.Fatal(err)
	}
}

func (a *Application) commentHandler(w http.ResponseWriter, r *http.Request) {
	vars := make(jet.VarMap)

	postId, err := strconv.Atoi(chi.URLParam(r, "postId"))
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	post, err := a.Models.Posts.Get(postId)
	if err != nil {
		a.serverError(w, err)
		return
	}

	comments, err := a.Models.Comments.GetForPost(post.ID)
	if err != nil {
		a.serverError(w, err)
		return
	}

	vars.Set("post", post)
	vars.Set("comments", comments)

	err = a.Render(w, r, "comments", vars)
	if err != nil {
		a.serverError(w, err)
		return
	}
}

func (a *Application) commentPostHandler(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, 1024*2)
	postId, err := strconv.Atoi(chi.URLParam(r, "postId"))
	if err != nil {
		a.serverError(w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		a.serverError(w, err)
		return
	}

	userId := a.Session.GetInt(r.Context(), sessionKeyUserId)

	form := utils.NewForm(r.PostForm)
	form.MinLength("comment", 10).MaxLength("comment", 255)

	if !form.Valid() {
		a.ErrLog.Printf("%+v", form.Errors)
		a.Session.Put(r.Context(), "flash", "Error: your comment is not valid: min: 10, max: 255")
		http.Redirect(w, r, fmt.Sprintf("/comments/%d", postId), http.StatusSeeOther)
		return
	}

	err = a.Models.Comments.Insert(form.Get("comment"), postId, userId)
	if !form.Valid() {
		a.Session.Put(r.Context(), "flash", "Error: "+err.Error())
		http.Redirect(w, r, fmt.Sprintf("/comments/%d", postId), http.StatusSeeOther)
		return
	}

	a.Session.Put(r.Context(), "flash", "comment created")
	http.Redirect(w, r, fmt.Sprintf("/comments/%d", postId), http.StatusSeeOther)
}

func (a *Application) loginHandler(w http.ResponseWriter, r *http.Request) {
	err := a.Render(w, r, "login", nil)
	if err != nil {
		a.serverError(w, err)
		return
	}
}

func (a *Application) loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*2)

	err := r.ParseForm()
	if err != nil {
		a.serverError(w, err)
		return
	}

	form := utils.NewForm(r.PostForm)
	form.Email("email")
	form.MinLength("password", 3)

	if !form.Valid() {
		vars := make(jet.VarMap)
		vars.Set("errors", form.Errors)

		err := a.Render(w, r, "login", vars)
		if err != nil {
			a.serverError(w, err)
			return
		}
	}

	user, err := a.Models.Users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		a.Session.Put(r.Context(), "flash", "login error: "+err.Error())
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	a.Session.RenewToken(r.Context())
	a.Session.Put(r.Context(), sessionKeyUserId, user.ID)
	a.Session.Put(r.Context(), sessionKeyUserName, user.Name)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (a *Application) signUpHandler(w http.ResponseWriter, r *http.Request) {
	vars := make(jet.VarMap)
	vars.Set("form", utils.NewForm(r.PostForm))

	err := a.Render(w, r, "signup", vars)
	if err != nil {
		a.serverError(w, err)
		return
	}
}

func (a *Application) signUpPostHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*2)

	err := r.ParseForm()
	if err != nil {
		a.serverError(w, err)
		return
	}

	vars := make(jet.VarMap)

	form := utils.NewForm(r.PostForm)
	vars.Set("form", form)

	form.Required("name", "email", "password").
		Email("email")

	if !form.Valid() {
		vars.Set("errors", form.Errors)
		err := a.Render(w, r, "signup", vars)
		if err != nil {
			a.serverError(w, err)
		}
		return
	}

	user := models.User{
		Name:      form.Get("name"),
		Email:     form.Get("email"),
		Password:  form.Get("password"),
		Activated: true, 
	}
	err = a.Models.Users.Insert(&user)
	if err != nil {
		form.Fail("signup", "failed to create account: "+err.Error())
		vars.Set("errors", form.Errors)
		err := a.Render(w, r, "signup", vars)
		if err != nil {
			a.serverError(w, err)
		}
		return
	}

	a.Session.Put(r.Context(), "flash", "account created!")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (a *Application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	a.Session.Remove(r.Context(), sessionKeyUserId)
	a.Session.Remove(r.Context(), sessionKeyUserName)
	a.Session.Destroy(r.Context())
	a.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}