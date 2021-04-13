package routes 

import (

	"encoding/json"
	"strings"
	"net/http"
	"log"
	"time"
	// "fmt"

	"kuclap-review-api/src/helper"
	"kuclap-review-api/src/models"
	"kuclap-review-api/src/storage"

	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"

)

func IndexRecapHandler(r *mux.Router) {

	var prefixPath = "/recap"

	// Storage Adpater
	r.HandleFunc(prefixPath + "/presigned/download/{recapid}", PresignedDownloadRecapEndPoint).Methods("GET")
	r.HandleFunc(prefixPath + "/presigned/upload/{classid}", PresignedUploadRecapEndPoint).Methods("GET")
	r.HandleFunc(prefixPath + "/upload", CreateRecapEndPoint).Methods("POST")

	// CRUD on Database
	r.HandleFunc("/recaps/last", LastRecapsEndPoint).Methods("GET")
	r.HandleFunc("/recaps/{classid}", AllRecapsByClassIDEndPoint).Methods("GET")
	r.HandleFunc(prefixPath + "/{recapid}", FindRecapEndpoint).Methods("GET")
	r.HandleFunc(prefixPath + "/{recapid}", DeleteRecapByIDEndPoint).Methods("DELETE")
	r.HandleFunc("/recaps", AllRecapsEndPoint).Methods("GET")

}

func PresignedDownloadRecapEndPoint(w http.ResponseWriter, r *http.Request) {
	
	params			:=	mux.Vars(r)
	
	recapID			:=	params["recapid"]
	recap, err		:=	mgoDAO.FindRecapByID(recapID)

	if err	!=	nil {
		log.Println("[ERR] mgo find recap: ", err)
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid recap-id or haven't your id in DB")
		return
	}

	url, err	:=	storage.PresignedURLDownloadRecapS3(&recap)

	if err != nil {
		log.Println("[ERR] presigned error: ", err)
		helper.RespondWithError(w, http.StatusBadRequest, "presigned error: your object isn't invalid.")
		return
	}

	if err	:=	mgoDAO.UpdateNumberDonwloadedByRecapID(recapID); err != nil {
		log.Println("[ERR] mgo update: ", err)
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

func PresignedUploadRecapEndPoint(w http.ResponseWriter, r *http.Request){

	// var	payload	storage.PresignedResponse

	params		:=	mux.Vars(r)
	author			:=	r.URL.Query().Get("author")

	recapID			:=	bson.NewObjectId().Hex()
	classID			:=	params["classid"]

	if classID == "" {
		log.Println("[ERR] no classId from param")
		helper.RespondWithError(w, http.StatusBadRequest, "Your class-id or review-id is invalid.")
		return
	}

	class, err	:=	mgoDAO.FindClassByClassID(classID)
	if err	!=	nil {
		log.Println("[ERR] mgo find class: ", err)
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	payload, err	:=	storage.PresignedURLUploadRecapS3(&class, author, recapID)

	log.Println(payload, err)

	if err != nil {
		log.Println("[ERR] presigned error: ", err)
		helper.RespondWithError(w, http.StatusBadRequest, "presigned error: your object isn't invalid.")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, payload)

}

func CreateRecapEndPoint(w http.ResponseWriter, r *http.Request){

	var newRecap	models.Recap
	defer r.Body.Close()

	if err	:=	json.NewDecoder(r.Body).Decode(&newRecap); err != nil {
		log.Println("[ERR] decode err: ", err)
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	newRecap.CreatedAt	=	time.Now().UTC().Add(7 * time.Hour)

	if err	:=	mgoDAO.CreateRecap(newRecap); err != nil {
		log.Println("[ERR] mgo create recap: ", err)
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err	:=	mgoDAO.UpdateNumberRecapByClassID(newRecap.ClassID, 1); err != nil {
		log.Println("[ERR] mgo update recap: ", err)
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusCreated, newRecap)

}

func LastRecapsEndPoint(w http.ResponseWriter, r *http.Request){
		
	page			:=	r.URL.Query().Get("page")
	offset			:=	r.URL.Query().Get("offset")

	recaps, err	:=	mgoDAO.LastRecaps(page ,offset)

	if err	!=	nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid recap offset")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, recaps)

}

func AllRecapsByClassIDEndPoint(w http.ResponseWriter, r *http.Request){

	params			:=	mux.Vars(r)

	page			:=	r.URL.Query().Get("page")
	offset			:=	r.URL.Query().Get("offset")

	recaps, err	:=	mgoDAO.FindRecapsByClassID(params["classid"], page, offset)

	if err	!=	nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid recap classid")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, recaps)
	
}

func FindRecapEndpoint(w http.ResponseWriter, r *http.Request){

	params			:= mux.Vars(r)
	recap, err	:= mgoDAO.FindRecapByID(params["recapid"])

	if err	!=	nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid recap-id or haven't your id in DB")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, recap)

}

func DeleteRecapByIDEndPoint(w http.ResponseWriter, r *http.Request){

	params			:=	mux.Vars(r) 
	
	reqToken		:=	r.Header.Get("Authorization")
	splitToken		:=	strings.Split(reqToken, " ")
	reqAuth			:=	splitToken[1]

	recap, err		:=	mgoDAO.FindRecapAllPropertyByID(params["recapid"])

	if err	!=	nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid recap-id or haven't your id in DB")
		return
	}

	if recap.Auth == reqAuth {
		
		// Delete the recap
		if err	:=	mgoDAO.DeleteRecapByID(params["recapid"]); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err	=	mgoDAO.UpdateNumberRecapByClassID(recap.ClassID, -1); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Delete review that storing the recap.
		review, err	:=	mgoDAO.FindReviewAllPropertyByID(recap.ReviewID)
		if err	!=	nil {
			helper.RespondWithError(w, http.StatusBadRequest, "Invalid review-id or haven't your id in DB")
			return
		}
		
		class, err	:=	mgoDAO.FindClassByClassID(recap.ClassID)
		if err	!=	nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		var newStats	models.StatClass
		var oldStats =	class.Stats

		if class.NumberReviewer ==  1 {
			// Ignore NaN when we divide with zero
			newStats.How		=	0
			newStats.Homework	=	0
			newStats.Interest	=	0
		} else {
			newStats.How		=	helper.GetNewStatsByDeleted(class.NumberReviewer, oldStats.How, review.Stats.How)
			newStats.Homework	=	helper.GetNewStatsByDeleted(class.NumberReviewer, oldStats.Homework,  review.Stats.Homework)
			newStats.Interest	=	helper.GetNewStatsByDeleted(class.NumberReviewer, oldStats.Interest,  review.Stats.Interest)
		}
		
		if err	:=	mgoDAO.DeleteByID(recap.ReviewID); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err	=	mgoDAO.UpdateStatsClassByDeleted(recap.ClassID, newStats); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		helper.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})

	} else {
		helper.RespondWithError(w, http.StatusInternalServerError, "your auth isn't match.")
	}

}

func AllRecapsEndPoint(w http.ResponseWriter, r *http.Request){

	recaps, err	:=	mgoDAO.FindAllRecaps()

	if err	!=	nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, recaps)

}