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

// IndexAnswersHandler is index routing for answer usecase
func IndexAnswersHandler(r *mux.Router) {

	var prefixPath = "/answer"

	r.HandleFunc(prefixPath, CreateAnswerEndPoint).Methods("POST")
	r.HandleFunc("/answers/{questionid}", AllAnswersByQuestionIDEndPoint).Methods("GET")
	r.HandleFunc(prefixPath+"/{answerid}", FindAnswerEndpoint).Methods("GET")
	r.HandleFunc(prefixPath+"/{answerid}", DeleteAnswerByIDEndPoint).Methods("DELETE")
	r.HandleFunc("/answers", AllAnswersEndPoint).Methods("GET")
	r.HandleFunc(prefixPath+"/report", CreateAnswerReportEndPoint).Methods("POST")
}

// CreateAnswerEndPoint is POST a new answer
func CreateAnswerEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var answer models.Answer

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updateAt := time.Now().UTC().Add(7 * time.Hour)
	if err := repository.UpdateNumberAnswerByQuestionID(ctx, answer.QuestionID, 1, updateAt); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	answer.CreatedAt = time.Now().UTC().Add(7 * time.Hour)
	answer.UpdateAt = answer.CreatedAt
	answer.ID = primitive.NewObjectID()

	if err := repository.CreateAnswer(ctx, answer); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusCreated, answer)

}

// AllAnswersByClassIDEndPoint is GET list of answers by class_id
func AllAnswersByQuestionIDEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	params := mux.Vars(r)

	answers, err := repository.FindLastAnswersByQuestionID(ctx, params["questionid"])

	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid answer classid")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, answers)

}

// FindAnswerEndpoint is GET a answer by its ID
func FindAnswerEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	params := mux.Vars(r)
	answer, err := repository.FindAnswerByID(ctx, params["answerid"])

	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid answer-id or haven't your id in DB")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, answer)

}

// AllAnswersEndPoint is GET list of answers
func AllAnswersEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	answers, err := repository.FindAllAnswers(ctx)

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, answers)

}

// DeleteAnswerByIDEndPoint is DELETE an existing answer
func DeleteAnswerByIDEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	params := mux.Vars(r)

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, " ")
	reqAuth := splitToken[1]

	answer, err := repository.FindAnswerAllPropertyByID(ctx, params["answerid"])
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid answer-id or haven't your id in DB")
		return
	}

	if answer.Auth == reqAuth {
		if err := repository.DeleteAnswerByID(ctx, params["answerid"]); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		updateAt := time.Now().UTC().Add(7 * time.Hour)
		if err = repository.UpdateNumberAnswerByQuestionID(ctx, answer.QuestionID, -1, updateAt); err != nil {
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
	ctx := context.Background()

	var report models.Report
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		log.Println("Error in decoding body on request", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updateAt := time.Now().UTC().Add(7 * time.Hour)

	if err := repository.UpdateAnswerReportByID(ctx, report.AnswerID, updateAt); err != nil {
		log.Println("Error in UpdateReportByID DAO", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid answer id")
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
