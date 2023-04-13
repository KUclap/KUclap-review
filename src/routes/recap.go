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
	"kuclap-review-api/src/storage"
)

func IndexRecapHandler(r *mux.Router) {

	var prefixPath = "/recap"

	// Storage Adpater
	r.HandleFunc(prefixPath+"/presigned/download/{recapid}", PresignedDownloadRecapEndPoint).Methods("GET")
	r.HandleFunc(prefixPath+"/presigned/upload/{classid}", PresignedUploadRecapEndPoint).Methods("GET")
	r.HandleFunc(prefixPath+"/upload", CreateRecapEndPoint).Methods("POST")

	// CRUD on Database
	r.HandleFunc("/recaps/last", LastRecapsEndPoint).Methods("GET")
	r.HandleFunc("/recaps/{classid}", AllRecapsByClassIDEndPoint).Methods("GET")
	r.HandleFunc(prefixPath+"/{recapid}", FindRecapEndpoint).Methods("GET")
	r.HandleFunc(prefixPath+"/{recapid}", DeleteRecapByIDEndPoint).Methods("DELETE")
	r.HandleFunc("/recaps", AllRecapsEndPoint).Methods("GET")
	r.HandleFunc(prefixPath+"/report", CreateRecapReportEndPoint).Methods("POST")

}

func PresignedDownloadRecapEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	params := mux.Vars(r)

	recapID := params["recapid"]
	recap, err := repository.FindRecapByID(ctx, recapID)

	if err != nil {
		log.Println("[ERR] mgo find recap: ", err)
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid recap-id or haven't your id in DB")
		return
	}

	url, err := storage.PresignedURLDownloadRecapS3(recap)

	if err != nil {
		log.Println("[ERR] presigned error: ", err)
		helper.RespondWithError(w, http.StatusBadRequest, "presigned error: your object isn't invalid.")
		return
	}

	updateAt := time.Now().UTC().Add(7 * time.Hour)
	if err := repository.UpdateNumberDownloadedByRecapID(ctx, recapID, updateAt); err != nil {
		log.Println("[ERR] mgo update: ", err)
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

func PresignedUploadRecapEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	params := mux.Vars(r)
	author := r.URL.Query().Get("author")

	recapID := primitive.NewObjectID().Hex()
	classID := params["classid"]

	if classID == "" {
		log.Println("[ERR] no classId from param")
		helper.RespondWithError(w, http.StatusBadRequest, "Your class-id or review-id is invalid.")
		return
	}

	class, err := repository.FindClassByClassID(ctx, classID)
	if err != nil {
		log.Println("[ERR] mgo find class: ", err)
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	payload, err := storage.PresignedURLUploadRecapS3(class, author, recapID)

	log.Println(payload, err)

	if err != nil {
		log.Println("[ERR] presigned error: ", err)
		helper.RespondWithError(w, http.StatusBadRequest, "presigned error: your object isn't invalid.")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, payload)

}

func CreateRecapEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var newRecap models.Recap
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&newRecap); err != nil {
		log.Println("[ERR] decode err: ", err)
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	newRecap.CreatedAt = time.Now().UTC().Add(7 * time.Hour)

	if err := repository.CreateRecap(ctx, newRecap); err != nil {
		log.Println("[ERR] mgo create recap: ", err)
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	updateAt := time.Now().UTC().Add(7 * time.Hour)
	if err := repository.UpdateNumberRecapByClassID(ctx, newRecap.ClassID, 1, updateAt); err != nil {
		log.Println("[ERR] mgo update recap: ", err)
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusCreated, newRecap)

}

func LastRecapsEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	page := r.URL.Query().Get("page")
	offset := r.URL.Query().Get("offset")

	recaps, err := repository.LastRecaps(ctx, page, offset)

	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid recap offset")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, recaps)

}

func AllRecapsByClassIDEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	params := mux.Vars(r)

	page := r.URL.Query().Get("page")
	offset := r.URL.Query().Get("offset")

	recaps, err := repository.FindRecapsByClassID(ctx, params["classid"], page, offset)

	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid recap classid")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, recaps)

}

func FindRecapEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	params := mux.Vars(r)
	recap, err := repository.FindRecapByID(ctx, params["recapid"])

	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid recap-id or haven't your id in DB")
		return
	}

	helper.RespondWithJson(w, http.StatusOK, recap)

}

func DeleteRecapByIDEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	params := mux.Vars(r)

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, " ")
	reqAuth := splitToken[1]

	recap, err := repository.FindRecapAllPropertyByID(ctx, params["recapid"])

	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid recap-id or haven't your id in DB")
		return
	}

	if recap.Auth == reqAuth {

		// Delete the recap
		if err := repository.DeleteRecapByID(ctx, params["recapid"]); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		updateAt := time.Now().UTC().Add(7 * time.Hour)
		if err = repository.UpdateNumberRecapByClassID(ctx, recap.ClassID, -1, updateAt); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if recap.ReviewID != "" {
			// Delete review that storing the recap.
			review, err := repository.FindReviewAllPropertyByID(ctx, recap.ReviewID)
			if err != nil {
				helper.RespondWithError(w, http.StatusBadRequest, "Invalid review-id or haven't your id in DB")
				return
			}

			class, err := repository.FindClassByClassID(ctx, recap.ClassID)
			if err != nil {
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

			if err := repository.DeleteByID(ctx, recap.ReviewID); err != nil {
				helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			if err = repository.UpdateStatsClassByDeleted(ctx, recap.ClassID, newStats); err != nil {
				helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		helper.RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})

	} else {
		helper.RespondWithError(w, http.StatusInternalServerError, "your auth isn't match.")
	}

}

func AllRecapsEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	recaps, err := repository.FindAllRecaps(ctx)

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJson(w, http.StatusOK, recaps)

}

func CreateRecapReportEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var report models.Report
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		log.Println("Error in decoding body on request", err.Error())
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updateAt := time.Now().UTC().Add(7 * time.Hour)
	if err := repository.UpdateRecapReportByID(ctx, report.RecapID, updateAt); err != nil {
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
