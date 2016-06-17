package handler

import (
	"html/template"
	"net/http"
)

//Home 主页面处理函数
func Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("home").ParseFiles("views/home.html")
	if err != nil {
		return
	}
	t.ExecuteTemplate(w, "home.html", nil)
}

//Join 新增成员处理函数
func Join(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	t, err := template.New("home").ParseFiles("views/room.html")
	if err != nil {
		w.Write([]byte("Failed to parse room.html"))
		return
	}
	user := map[string]string{"Email": email, "WebSocketHost": r.Host}
	t.ExecuteTemplate(w, "room.html", user)
}
