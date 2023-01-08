package handler

import "net/http"

// 拦截器模式
func AuthInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			username := r.Form.Get("username")
			token := r.Form.Get("token")

			if len(username) < 3 || !IsTokenValid(username, token) {

				w.WriteHeader(http.StatusForbidden)
				return

			}
			h(w, r)
		})
}
