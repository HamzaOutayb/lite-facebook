package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"social-network/internal/models"
	utils "social-network/pkg"
)

func (H *Handler) Login(w http.ResponseWriter, r *http.Request) {
	utils.SetSessionCookie(w, "550e8400-e29b-41d4-a716-446655440000")
	utils.WriteJson(w, 200, "nice")
	return
	if r.Method != http.MethodPost {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var user models.User
	if err := utils.ParseBody(r, &user); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, "Bad request")
		return
	}

	err := H.Service.LoginUser(&user)
	if err != nil {
		fmt.Println("err", err.Error())
		utils.WriteJson(w, http.StatusBadRequest, err.Error())
		return
	}

	userinfo := models.UserInfo{
		Nickname:   user.Nickname,
		First_Name: user.First_Name,
		Last_Name:  user.Last_Name,
		Image:      user.Image,
	}

	utils.WriteJson(w, http.StatusOK, userinfo)
}

func (H *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	user := H.Service.Extractuser(r)

	// Parse the multipart form (10MB max file size)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, "file too big")
		return
	}

	// Extract profile picture (optional)
	var filePath string
	file, handler, err := r.FormFile("avatar")
	if err == nil { // No error means a file was uploaded

		// Ensure Profile directory exists
		uploadDir := "../backend/internal/repository/profile"
		defer file.Close()
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.Mkdir(uploadDir, os.ModePerm)
		}

		// Save with a unique filename (e.g., user.UUID + filename)
		filePath = filepath.Join(uploadDir, handler.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			fmt.Println(err)
			utils.WriteJson(w, http.StatusInternalServerError, "Could not save file")
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, file)
		if err != nil {
			utils.WriteJson(w, http.StatusInternalServerError, "Failed to save file")
			return
		}

		// Assign file path to user struct
		user.Image = filePath
	}

	// Proccess Data and Insert it
	err = H.Service.RegisterUser(&user)
	if err != nil {
		fmt.Println("yes",err.Error())
		utils.WriteJson(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.WriteJson(w, http.StatusOK, "You'v loged in succesfuly")
}

func (H *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var Uuid string
	err := json.NewDecoder(r.Body).Decode(&Uuid)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, "bad request")
	}
	err = H.Service.DeleteSessionCookie(w, Uuid)
	if err != nil {
		utils.WriteJson(w, http.StatusOK, err.Error())
		return
	}
	utils.WriteJson(w, http.StatusOK, "You Logged Out Successfuly!")
}
