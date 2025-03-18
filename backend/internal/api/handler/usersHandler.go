package handler

import (
	"log"
	"net/http"
	"social-network/internal/models"
	utils "social-network/pkg"
)

func (Handler *Handler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user, ok := r.Context().Value(utils.UserIDKey).(models.UserInfo)
	if !ok {
		utils.WriteJson(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var profile models.UserInfo
	var err error

	err = utils.ParseBody(r, &profile)
	if err != nil {
		log.Println(err)
		utils.WriteJson(w, http.StatusBadRequest, "Bad request")
		return
	}
	if profile.ID == 0 {
		profile.ID = user.ID
	}

	profile, err = Handler.Service.GetUserByID(profile.ID)
	if err != nil {
		log.Println(err)
		utils.WriteJson(w, http.StatusBadRequest, "Bad request")
		return
	}
	var fullProfile models.Profile
	fullProfile.UserInfo = profile
	Handler.Service.GetFollowCounts(&fullProfile)
	Handler.Service.SetAction(&fullProfile, user.ID)


	utils.WriteJson(w, http.StatusOK, fullProfile)
}

func (Handler *Handler) HandleGetProfileAbout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user, ok := r.Context().Value(utils.UserIDKey).(models.UserInfo)
	if !ok {
		utils.WriteJson(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var profile models.User
	var err error

	err = utils.ParseBody(r, &profile)
	if err != nil {
		log.Println(err)
		utils.WriteJson(w, http.StatusBadRequest, "Bad request")
		return
	}
	if profile.ID == 0 {
		profile.ID = user.ID
	}

	err = Handler.Service.GetProfile(&profile, user.ID)
	if err != nil {
		log.Println(err)
		if err.Error() == "profile is private, follow to see"{
			utils.WriteJson(w, http.StatusForbidden, http.StatusText( http.StatusForbidden))
			return
		}
		utils.WriteJson(w, http.StatusBadRequest, "Bad request")
		return
	}

	utils.WriteJson(w, http.StatusOK, profile)
}

func (Handler *Handler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
}
