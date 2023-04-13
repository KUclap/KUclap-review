package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"kuclap-review-api/src/helper"
	"kuclap-review-api/src/models"
)

// IndexQuestionsHandler is index routing for question usecase
func IndexQuestionsHandler(r *mux.Router) {

	var prefixPath = "/question"

	r.HandleFunc(prefixPath, CreateQuestionEndPoint).Methods("POST")
	r.HandleFunc("/questions/last", LastQuestionsEndPoint).Methods("GET")
	r.HandleFunc("/questions/{classid}", AllQuestionsByClassIDEndPoint).Methods("GET")
	r.HandleFunc(prefixPath+"/{questionid}", FindQuestionEndpoint).Methods("GET")
	r.HandleFunc(prefixPath+"/{questionid}", DeleteQuestionByIDEndPoint).Methods("DELETE")
	r.HandleFunc("/questions", AllQuestionsEndPoint).Methods("GET")
	r.HandleFunc(prefixPath+"/report", CreateQuestionReportEndPoint).Methods("POST")
	// r.HandleFunc("/questions/reported", FindQuestionReportedEndpoint).Methods("GET")
	// r.HandleFunc("/questions/{questionid}", UpdateQuestionEndPoint).Methods("PUT")
}

// CreateQuestionEndPoint is POST a new question
func CreateQuestionEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var question models.Question

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	class, err := repository.FindClassByClassID(ctx, question.ClassID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	updateAt := time.Now().UTC().Add(7 * time.Hour)
	if err = repository.UpdateNumberQuestionByClassID(ctx, question.ClassID, 1, updateAt); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	question.ClassNameTH = class.NameTH
	question.ClassNameEN = class.NameEN
	question.CreatedAt = time.Now().UTC().Add(7 * time.Hour)
	question.UpdateAt = question.CreatedAt
	question.ID = primitive.NewObjectID()
	// question.Answer			=	make([]models.Answer, 0)

	if err := repository.CreateQuestion(ctx, question); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusCreated, question)

}

// LastQuestionsEndPoint is GET list of questions
// Read param on UrlQuery (eg. /last?offset=5 )
// Paging by query: page={number_page} offset={number_offset}
func LastQuestionsEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	page := r.URL.Query().Get("page")
	offset := r.URL.Query().Get("offset")
	questions, err := repository.LastQuestions(ctx, page, offset)

	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid question offset")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, questions)

}

// AllQuestionsByClassIDEndPoint is GET list of questions by class_id
func AllQuestionsByClassIDEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	params := mux.Vars(r)
	page := r.URL.Query().Get("page")
	offset := r.URL.Query().Get("offset")
	questions, err := repository.FindQuestionsByClassID(ctx, params["classid"], page, offset)

	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid question classid")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, questions)

}

// FindQuestionEndpoint is GET a question by its ID
func FindQuestionEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	params := mux.Vars(r)
	question, err := repository.FindQuestionByID(ctx, params["questionid"])

	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid question-id or haven't your id in DB")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, question)

}

// AllQuestionsEndPoint is GET list of questions
func AllQuestionsEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	questions, err := repository.FindAllQuestions(ctx)

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, questions)

}

// DeleteQuestionByIDEndPoint is DELETE an existing question
func DeleteQuestionByIDEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	params := mux.Vars(r)

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, " ")
	reqAuth := splitToken[1]

	question, err := repository.FindQuestionAllPropertyByID(ctx, params["questionid"])

	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid question-id or haven't your id in DB")
		return
	}

	if question.Auth == reqAuth {

		if err := repository.DeleteQuestionByID(ctx, params["questionid"]); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err = repository.DeleteAnswersByQuestionID(ctx, params["questionid"]); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		updateAt := time.Now().UTC().Add(7 * time.Hour)

		if err = repository.UpdateNumberQuestionByClassID(ctx, question.ClassID, -1, updateAt); err != nil {
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
	ctx := context.Background()

	var report models.Report
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		log.Println("Error in decoding body on request", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updateAt := time.Now().UTC().Add(7 * time.Hour)

	if err := repository.UpdateQuestionReportByID(ctx, report.QuestionID, updateAt); err != nil {
		log.Println("Error in UpdateReportByID DAO", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review id")
		return
	}

	report.CreatedAt = time.Now().UTC().Add(7 * time.Hour)

	if err := repository.InsertReport(ctx, report); err != nil {
		log.Println("Error in InsertReport DAO", err.Error())
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusCreated, report)

}
