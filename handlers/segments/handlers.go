package segments

import (
	"encoding/json"
	"log"
	"net/http"
)

func CreateSegment(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	segment := &Segment{}
	err := json.NewDecoder(req.Body).Decode(&segment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if segment.Name == "" || len(segment.Name) > 255 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	if err = InsertSegment(segment); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteSegment(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	segment := &Segment{}
	err := json.NewDecoder(req.Body).Decode(&segment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if segment.Name == "" || len(segment.Name) > 255 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err = DeleteSegmentByName(segment.Name); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
