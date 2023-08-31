package users_segments

import (
	"backend-trainee-assignment-2023/db"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func UpdateUserSegments(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	userSegments := &struct {
		SegmentsToAdd    []string
		SegmentsToDelete []string
		UserId           int
	}{}

	err := json.NewDecoder(req.Body).Decode(&userSegments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userSegments.UserId < 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	tx, err := db.DB.BeginTx(context.Background(), nil)
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

func UserSegments(w http.ResponseWriter, req *http.Request) {
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
