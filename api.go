package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type ApiServer struct {
	listenAddr string
	store      Storage
}

func NewApiServer(l string, s Storage) *ApiServer {
	return &ApiServer{
		listenAddr: l,
		store:      s,
	}
}

func (s *ApiServer) Run() {
	r := mux.NewRouter()

	r.HandleFunc("/", s.handleHomePage)

	r.HandleFunc("/api/products", s.handleProducts)
	r.HandleFunc("/api/products/tracksuits", s.handleProductsTracksuits)
	r.HandleFunc("/api/products/windbreakers", s.handleProductsWindbreakers)
	r.HandleFunc("/api/products/{productTitle}", s.handleProduct)

	r.HandleFunc("/api/users", s.protectMiddleware(s.handleUsers))
	r.HandleFunc("/api/users/{userId}", s.protectMiddleware(s.handleUser))

	r.HandleFunc("/api/login", s.handleLogin)

	http.ListenAndServe(s.listenAddr, r)
}

func (s *ApiServer) handleHomePage(w http.ResponseWriter, r *http.Request) {

}

func (s *ApiServer) handleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleGetAllProducts(w, r)
	case "POST":
		s.handleCreateProduct(w, r)
	case "PUT":
		s.handleUpdateProduct(w, r)
	default:
		err := fmt.Errorf("method %s not allowed", r.Method)
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
	}
}

func (s *ApiServer) handleProduct(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleGetProduct(w, r)
	case "DELETE":
		s.handleDeleteProduct(w, r)
	default:
		err := fmt.Errorf("method %s not allowed", r.Method)
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
	}
}

func (s *ApiServer) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleGetAllUsers(w, r)
	case "POST":
		s.handleCreateUser(w, r)
	case "PUT":
		s.handleUpdateUser(w, r)
	default:
		err := fmt.Errorf("method %s not allowed", r.Method)
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
	}
}

func (s *ApiServer) handleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleGetUser(w, r)
	case "DELETE":
		s.handleDeleteUser(w, r)
	default:
		err := fmt.Errorf("method %s not allowed", r.Method)
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
	}
}

func (s *ApiServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	loginReq := new(LoginRequest)
	err := json.NewDecoder(r.Body).Decode(loginReq)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	user, err := s.store.GetUserByUsername(loginReq.Username)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if !user.VerifyPassword(loginReq.Password) {
		http.Error(w, "wrong password", http.StatusBadRequest)
		return
	}

	token, err := CreateJWT(user.Id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	cookie := CreateCookie("token", token)
	http.SetCookie(w, cookie)
	WriteJSON(w, "login successful", http.StatusOK)
}

func (s *ApiServer) handleGetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := s.store.GetAllProducts()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	WriteJSON(w, products, http.StatusOK)
}

func (s *ApiServer) handleProductsTracksuits(w http.ResponseWriter, r *http.Request) {
	products, err := s.store.GetProductsByType("tracksuit")
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	WriteJSON(w, products, http.StatusOK)
}

func (s *ApiServer) handleProductsWindbreakers(w http.ResponseWriter, r *http.Request) {
	products, err := s.store.GetProductsByType("windbreaker")
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	WriteJSON(w, products, http.StatusOK)
}

func (s *ApiServer) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	product := new(Product)
	err := json.NewDecoder(r.Body).Decode(product)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	product, err = s.store.CreateProduct(*product)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	WriteJSON(w, product.Id, http.StatusCreated)
}

func (s *ApiServer) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	product := new(Product)
	err := json.NewDecoder(r.Body).Decode(product)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = s.store.UpdateProduct(*product)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	WriteJSON(w, "product updated", http.StatusCreated)
}

func (s *ApiServer) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["productTitle"]
	title = strings.Replace(title, "-", " ", -1)

	product, err := s.store.GetProductByTitle(title)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	WriteJSON(w, *product, http.StatusOK)
}

func (s *ApiServer) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["productTitle"]
	title = strings.Replace(title, "-", " ", -1)

	product, err := s.store.GetProductByTitle(title)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = s.store.DeleteProductById(product.Id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	WriteJSON(w, "product deleted", http.StatusOK)
}

func (s *ApiServer) handleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.store.GetAllUsers()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	WriteJSON(w, users, http.StatusOK)
}

func (s *ApiServer) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	user, err = s.store.CreateUser(*user)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	WriteJSON(w, user.Id, http.StatusOK)
}

func (s *ApiServer) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = s.store.UpdateUser(*user)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	WriteJSON(w, "user updated", http.StatusOK)
}

func (s *ApiServer) handleGetUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["userId"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	user, err := s.store.GetUserById(id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	WriteJSON(w, *user, http.StatusOK)
}

func (s *ApiServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["userId"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = s.store.DeleteUserById(id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	WriteJSON(w, "user deleted", http.StatusOK)
}

func (s *ApiServer) protectMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "error: not authorized", http.StatusUnauthorized)
			return
		}

		// get token
		token, err := ValidateJWT(cookie.Value)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "error: not authorized", http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			slog.Error(err.Error())
			http.Error(w, "error: not authorized", http.StatusUnauthorized)
			return
		}

		// get userId
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			slog.Error(err.Error())
			http.Error(w, "error: not authorized", http.StatusUnauthorized)
			return
		}
		userId, ok := claims["userId"].(float64)
		if !ok {
			slog.Error(err.Error())
			http.Error(w, "error: not authorized", http.StatusUnauthorized)
			return
		}

		// check for user
		user, err := s.store.GetUserById(int(userId))
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "error: user not found", http.StatusNotFound)
			return
		}

		// check user level
		reqPath := r.URL.Path
		endpoint := strings.Split(reqPath, "/")[1]

		switch endpoint {
		case "products":
			if user.Level > 1 {
				err := fmt.Errorf("user %s not authorized to access %s: userLevel(%d) > 1", user.Username, endpoint, user.Level)
				slog.Error(err.Error())
				http.Error(w, "not authorized", http.StatusUnauthorized)
				return
			}
		case "users":
			if user.Level > 0 {
				err := fmt.Errorf("user %s not authorized to access %s: userLevel(%d) > 1", user.Username, endpoint, user.Level)
				slog.Error(err.Error())
				http.Error(w, "not authorized", http.StatusUnauthorized)
				return
			}
		}

		// call the next func with user in context
		ctx := context.WithValue(r.Context(), ContextKeyUser, user)
		next(w, r.WithContext(ctx))
	}
}
