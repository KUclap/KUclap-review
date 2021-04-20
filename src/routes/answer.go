package routes 

import (
	"log"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"kuclap-review-api/src/models"
	"kuclap-review-api/src/helper"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// IndexAnswersHandler is index routing for answer usecase
func IndexAnswersHandler(r *mux.Router) {

	var prefixPath = "/answer"

	r.HandleFunc(prefixPath, CreateAnswerEndPoint).Methods("POST")
	r.HandleFunc("/answers/{questionid}", AllAnswersByQuestionIDEndPoint).Methods("GET")
	r.HandleFunc(prefixPath + "/{answerid}", FindAnswerEndpoint).Methods("GET")
	r.HandleFunc(prefixPath + "/{answerid}", DeleteAnswerByIDEndPoint).Methods("DELETE")
	r.HandleFunc("/answers", AllAnswersEndPoint).Methods("GET")
	r.HandleFunc(prefixPath + "/report", CreateAnswerReportEndPoint).Methods("POST")
}

// CreateAnswerEndPoint is POST a new answer
func CreateAnswerEndPoint(w http.ResponseWriter, r *http.Request) {
	
	var	answer		models.Answer

	defer r.Body.Close()

	if err		:=	json.NewDecoder(r.Body).Decode(&answer); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updateAt	:=	time.Now().UTC().Add(7 * time.Hour)
	if err		:=	mgoDAO.UpdateNumberAnswerByQuestionID(answer.QuestionID, 1, updateAt); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	answer.CreatedAt	=	time.Now().UTC().Add(7 * time.Hour)
	answer.UpdateAt		=	answer.CreatedAt
	answer.ID			=	bson.NewObjectId()

	if err		:=	mgoDAO.CreateAnswer(answer); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusCreated, answer)

}

// AllAnswersByClassIDEndPoint is GET list of answers by class_id
func AllAnswersByQuestionIDEndPoint(w http.ResponseWriter, r *http.Request) {

	params			:=	mux.Vars(r)
	
	answers, err	:=	mgoDAO.FindLastAnswersByQuestionID(params["questionid"])

	if err	!=	nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid answer classid")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, answers)

}

// FindAnswerEndpoint is GET a answer by its ID
func FindAnswerEndpoint(w http.ResponseWriter, r *http.Request) {

	params			:= mux.Vars(r)
	answer, err	:= mgoDAO.FindAnswerByID(params["answerid"])

	if err	!=	nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid answer-id or haven't your id in DB")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, answer)

}

// AllAnswersEndPoint is GET list of answers
func AllAnswersEndPoint(w http.ResponseWriter, r *http.Request) {
	
	answers, err	:=	mgoDAO.FindAllAnswers()

	if err	!=	nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, answers)

}

// DeleteAnswerByIDEndPoint is DELETE an existing answer
func DeleteAnswerByIDEndPoint(w http.ResponseWriter, r *http.Request) {

	params			:=	mux.Vars(r) 

	reqToken		:=	r.Header.Get("Authorization")
	splitToken		:=	strings.Split(reqToken, " ")
	reqAuth			:=	splitToken[1]

	answer, err	:=	mgoDAO.FindAnswerAllPropertyByID(params["answerid"])
	if err	!=	nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid answer-id or haven't your id in DB")
		return
	}

	if answer.Auth == reqAuth {
		
		if err	:=	mgoDAO.DeleteAnswerByID(params["answerid"]); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		updateAt		:=	time.Now().UTC().Add(7 * time.Hour)
		if err	=	mgoDAO.UpdateNumberAnswerByQuestionID(answer.QuestionID, -1, updateAt); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		helper.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})

	} else {

		helper.RespondWithError(w, http.StatusInternalServerError, "your auth isn't match.")

	}
}

// CreateAnswerReportEndPoint is POST create report for the review
func CreateAnswerReportEndPoint(w http.ResponseWriter, r *http.Request) { 
	
	var report models.Report
	defer r.Body.Close()
	
	if err	:=	json.NewDecoder(r.Body).Decode(&report); err != nil {
		log.Println("Error in decoding body on request", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updateAt	:=	time.Now().UTC().Add(7 * time.Hour)

	if err	:=	mgoDAO.UpdateAnswerReportByID(report.AnswerID, updateAt); err != nil {
		log.Println("Error in UpdateReportByID DAO", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review id")
		return
	}

	report.CreatedAt	=	time.Now().UTC().Add(7 * time.Hour)

	if err	:=	mgoDAO.InsertReport(report); err != nil {
		log.Println("Error in InsertReport DAO", err.Error())
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusCreated, report)
	
}