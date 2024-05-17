package user

import (
	"fmt"
	"net/http"

	"github.com/xhoang0509/ecom-api/config"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/xhoang0509/ecom-api/services/auth"
	"github.com/xhoang0509/ecom-api/types"
	"github.com/xhoang0509/ecom-api/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", h.handleRegister).Methods(http.MethodPost)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user types.LoginUserPayload

	if err := utils.ParseJson(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// check user exist with email

	foundUser, _ := h.store.GetUserByEmail(user.Email)
	if foundUser == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user with email %s found", user.Email))
		return
	}

	err := auth.CompareHashPassword(user.Password, foundUser.Password)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, foundUser.ID)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var user types.RegisterUserPayload
	if err := utils.ParseJson(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// check if the user exists
	foundUser, err := h.store.GetUserByEmail(user.Email)
	if foundUser != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", user.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// if it doesn't, we create the new user
	err = h.store.CreateUser(types.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, nil)
}
