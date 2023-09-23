package segments

import (
	"encoding/json"
	"log"
	"net/http"
)

// @Summary Create segment
// @Tags segment
// @Description Method to create a segment
// @ID CreateSegment
// @Accept json
// @Produce json
// @Param input body Segment true "Segment name"
// @Success 201 {object} Segment
// @Failure 400 "Bad Request"
// @Failure 405 "Method Not Allowed"
// @Failure 500 "Internal server error"
// @Router /createSegment [post]
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(segment); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// @Summary Delete segment
// @Tags segment
// @Description Method to delete a segment
// @ID DeleteSegment
// @Accept json
// @Param input body Segment true "Segment name"
// @Success 200 "OK"
// @Failure 400 "Bad Request"
// @Failure 405 "Method Not Allowed"
// @Failure 500 "Internal server error"
// @Router /deleteSegment [delete]
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
