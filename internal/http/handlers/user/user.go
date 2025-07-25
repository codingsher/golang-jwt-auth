package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/codingsher/user-jwt-auth/internal/storage"
	"github.com/codingsher/user-jwt-auth/internal/types"
	"github.com/codingsher/user-jwt-auth/internal/utils/jwt"
	"github.com/codingsher/user-jwt-auth/internal/utils/response"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating user")

		var user types.UserRegister
		err := json.NewDecoder(r.Body).Decode(&user)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}

		// request validation
		if err := validator.New().Struct(user); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		fmt.Println(user.PasswordHash)
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)

		lastId, err := storage.NewUserRegister(user.UserName, user.Email, string(hashedPassword))
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
		}
		slog.Info("User Created Successfully", slog.String("user id", fmt.Sprint(lastId)))

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func LoginUser(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user types.UserLogin
		err := json.NewDecoder(r.Body).Decode(&user)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}

		// request validation
		if err := validator.New().Struct(user); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		UserPassword, err := storage.GetUserByEmail(user.Email)
		if err != nil {
			slog.Error("error getting users password")
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(UserPassword), []byte(user.Password))
		if err != nil {
			slog.Error("password incorrect")
			fmt.Println(UserPassword)
			fmt.Println(user.Password)
			fmt.Println(bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// generate jwt tokens now. as the user is in database and the password matches
		accessToken, err := jwt.GenerateJWT(user.Email, "access", 30*time.Minute)
		if err != nil {
			slog.Error("JWT generation failed", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		refreshToken, err := jwt.GenerateJWT(user.Email, "refresh", 7*24*time.Hour)
		if err != nil {
			slog.Error("Refresh token generation failed", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}
