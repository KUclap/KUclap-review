package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"kuclap-review-api/src/helper"
	"kuclap-review-api/src/models"
)

// IndexAdminHandler is index routing for review usecase
func IndexAdminHandler(r *mux.Router) {

	var prefixPath = "/admin"

	r.HandleFunc(prefixPath+"/review/delete", AdminDeleteReviewByIDEndPoint).Methods("PUT")

}

// AdminDeleteReviewByIDEndPoint is PUT an existing review
func AdminDeleteReviewByIDEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var body models.AdminDeleteByID

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("Error in decoding body on request", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	review, err := repository.FindReviewAllPropertyByID(ctx, body.ReviewID)
	if err != nil {
		log.Println("Error in FindReviewAllPropertyByID DAO", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid review-id or haven't your id in DB")
		return
	}

	// Delete the review.
	class, err := repository.FindClassByClassID(ctx, review.ClassID)
	if err != nil {
		log.Println("Error in FindClassByClassID DAO", err.Error())
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var newStats models.StatClass
	var oldStats = class.Stats

	if class.NumberReviewer == 1 {
		// Ignore NaN when we divide with zero
		newStats.How = 0
		newStats.Homework = 0
		newStats.Interest = 0
	} else {
		newStats.How = helper.GetNewStatsByDeleted(class.NumberReviewer, oldStats.How, review.Stats.How)
		newStats.Homework = helper.GetNewStatsByDeleted(class.NumberReviewer, oldStats.Homework, review.Stats.Homework)
		newStats.Interest = helper.GetNewStatsByDeleted(class.NumberReviewer, oldStats.Interest, review.Stats.Interest)
	}

	newStats.UpdateAt = time.Now().UTC().Add(7 * time.Hour)

	updateAt := time.Now().UTC().Add(7 * time.Hour)

	if err = repository.UpdateAdminDeleteByID(ctx, body.ReviewID, body.DeleteReason, updateAt); err != nil {
		log.Println("Error in UpdateAdminDeleteByID DAO", err.Error())
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = repository.UpdateStatsClassByDeleted(ctx, review.ClassID, newStats); err != nil {
		log.Println("Error in UpdateStatsClassByDeleted DAO", err.Error())
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if review.RecapID != "" {
		// Delete recap on the review.
		if err := repository.DeleteRecapByID(ctx, review.RecapID); err != nil {
			log.Println("Error in DeleteRecapByID DAO", err.Error())
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err = repository.UpdateNumberRecapByClassID(ctx, review.ClassID, -1, updateAt); err != nil {
			log.Println("Error in UpdateNumberRecapByClassID DAO", err.Error())
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}
