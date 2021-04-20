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

// IndexQuestionsHandler is index routing for question usecase
func IndexQuestionsHandler(r *mux.Router) {

	var prefixPath = "/question"

	r.HandleFunc(prefixPath, CreateQuestionEndPoint).Methods("POST")
	r.HandleFunc("/questions/last", LastQuestionsEndPoint).Methods("GET")
	r.HandleFunc("/questions/{classid}", AllQuestionsByClassIDEndPoint).Methods("GET")
	r.HandleFunc(prefixPath + "/{questionid}", FindQuestionEndpoint).Methods("GET")
	r.HandleFunc(prefixPath + "/{questionid}", DeleteQuestionByIDEndPoint).Methods("DELETE")
	r.HandleFunc("/questions", AllQuestionsEndPoint).Methods("GET")
	r.HandleFunc(prefixPath + "/report", CreateQuestionReportEndPoint).Methods("POST")
	// r.HandleFunc("/questions/reported", FindQuestionReportedEndpoint).Methods("GET")
	// r.HandleFunc("/questions/{questionid}", UpdateQuestionEndPoint).Methods("PUT")
}

// CreateQuestionEndPoint is POST a new question
func CreateQuestionEndPoint(w http.ResponseWriter, r *http.Request) {

	var	class		models.Class
	var	question	models.Question

	defer r.Body.Close()

	if err		:=	json.NewDecoder(r.Body).Decode(&question); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	
	class, err	:=	mgoDAO.FindClassByClassID(question.ClassID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	updateAt	:=	time.Now().UTC().Add(7 * time.Hour)
	if err	=	mgoDAO.UpdateNumberQuestionByClassID(question.ClassID, 1, updateAt); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	question.ClassNameTH	=	class.NameTH
	question.ClassNameEN	=	class.NameEN
	question.CreatedAt		=	time.Now().UTC().Add(7 * time.Hour)
	question.UpdateAt		=	question.CreatedAt
	question.ID				=	bson.NewObjectId()
	// question.Answer			=	make([]models.Answer, 0)

	if err	:=	mgoDAO.CreateQuestion(question); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusCreated, question)

}

// LastQuestionsEndPoint is GET list of questions 
// Read param on UrlQuery (eg. /last?offset=5 )
// Paging by query: page={number_page} offset={number_offset}
func LastQuestionsEndPoint(w http.ResponseWriter, r *http.Request) {
	
	page			:=	r.URL.Query().Get("page")
	offset			:=	r.URL.Query().Get("offset")
	questions, err	:=	mgoDAO.LastQuestions(page ,offset)

	if err	!=	nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid question offset")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, questions)

}

// AllQuestionsByClassIDEndPoint is GET list of questions by class_id
func AllQuestionsByClassIDEndPoint(w http.ResponseWriter, r *http.Request) {

	params			:=	mux.Vars(r)
	page			:=	r.URL.Query().Get("page")
	offset			:=	r.URL.Query().Get("offset")
	questions, err	:=	mgoDAO.FindQuestionsByClassID(params["classid"], page, offset)

	if err	!=	nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid question classid")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, questions)

}

// FindQuestionEndpoint is GET a question by its ID
func FindQuestionEndpoint(w http.ResponseWriter, r *http.Request) {

	params			:= mux.Vars(r)
	question, err	:= mgoDAO.FindQuestionByID(params["questionid"])

	if err	!=	nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid question-id or haven't your id in DB")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, question)

}

// AllQuestionsEndPoint is GET list of questions
func AllQuestionsEndPoint(w http.ResponseWriter, r *http.Request) {
	
	questions, err	:=	mgoDAO.FindAllQuestions()

	if err	!=	nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, questions)

}

// DeleteQuestionByIDEndPoint is DELETE an existing question
func DeleteQuestionByIDEndPoint(w http.ResponseWriter, r *http.Request) {

	params			:=	mux.Vars(r) 

	reqToken		:=	r.Header.Get("Authorization")
	splitToken		:=	strings.Split(reqToken, " ")
	reqAuth			:=	splitToken[1]

	question, err	:=	mgoDAO.FindQuestionAllPropertyByID(params["questionid"])

	if err	!=	nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid question-id or haven't your id in DB")
		return
	}

	if question.Auth == reqAuth {
		
		if err	:=	mgoDAO.DeleteQuestionByID(params["questionid"]); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err	=	mgoDAO.DeleteAnswersByQuestionID(params["questionid"]); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		updateAt		:=	time.Now().UTC().Add(7 * time.Hour)
	
		if err	=	mgoDAO.UpdateNumberQuestionByClassID(question.ClassID, -1, updateAt); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		helper.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})

	} else {
		helper.RespondWithError(w, http.StatusInternalServerError, "your auth isn't match.")
	}
}

// CreateQuestionReportEndPoint is POST create report for the review
func CreateQuestionReportEndPoint(w http.ResponseWriter, r *http.Request) { 
	
	var report models.Report
	defer r.Body.Close()
	
	if err	:=	json.NewDecoder(r.Body).Decode(&report); err != nil {
		log.Println("Error in decoding body on request", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updateAt	:=	time.Now().UTC().Add(7 * time.Hour)

	if err	:=	mgoDAO.UpdateQuestionReportByID(report.QuestionID, updateAt); err != nil {
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