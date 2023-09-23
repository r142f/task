package users_segments

import (
	"backend-trainee-assignment-2023/config"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// @Summary Update user segments
// @Tags user_segment
// @Description Method to add / delete user segments
// @ID UpdateSegments
// @Accept json
// @Param input body UserSegments true "Segment names to add/delete, user id"
// @Success 201 "Created"
// @Failure 400 "Bad Request"
// @Failure 405 "Method Not Allowed"
// @Failure 500 "Internal server error"
// @Router /updateUserSegments [post]
func UpdateUserSegments(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	userSegments := &UserSegments{}

	err := json.NewDecoder(req.Body).Decode(&userSegments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userSegments.UserId < 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	tx, err := config.DB.BeginTx(context.Background(), nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for _, segmentToAdd := range userSegments.SegmentsToAdd {
		err := InsertUserSegmentWithSegmentName(tx, userSegments.UserId, segmentToAdd)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			err = tx.Rollback()
			if err != nil {
				log.Printf("tx.Rollback(): %v\n", err)
			}

			return
		}
	}

	for _, segmentToDelete := range userSegments.SegmentsToDelete {
		err := DeleteUserSegmentWithSegmentName(tx, userSegments.UserId, segmentToDelete)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			err = tx.Rollback()
			if err != nil {
				log.Printf("tx.Rollback(): %v\n", err)
			}

			return
		}
	}

	if err = tx.Commit(); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		err = tx.Rollback()
		if err != nil {
			log.Printf("tx.Rollback(): %v\n", err)
		}
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Get user segments
// @Tags user_segment
// @Description Method to get user segments
// @ID UserSegments
// @Param userId query int true "Get segments by userId"
// @Success 200 {array} string
// @Failure 400 "Bad Request"
// @Failure 405 "Method Not Allowed"
// @Failure 500 "Internal server error"
// @Router /userSegments [get]
func SegmentsByUser(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	userId, err := strconv.Atoi(req.FormValue("userId"))
	if err != nil || userId < 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	segmentNames, err := SelectUserSegments(userId)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(segmentNames); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
