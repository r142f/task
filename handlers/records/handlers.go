package records

import (
	"backend-trainee-assignment-2023/config"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// @Summary Form a report and get a link
// @Tags report
// @Description A method for generating a report with the history of a user entering/exiting a segment from specified month and year until now. Returns link to the report
// @ID GenerateReport
// @Produce json
// @Param userId query int true "Id of the user"
// @Param year query int true "Year from"
// @Param month query int true "Month from"
// @Success 200 {object} ReportLink
// @Failure 400 "Bad Request"
// @Failure 405 "Method Not Allowed"
// @Failure 500 "Internal server error"
// @Router /generateReport [get]
func Report(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	year, err := strconv.Atoi(req.FormValue("year"))
	if err != nil || year <= 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	month, err := strconv.Atoi(req.FormValue("month"))
	if err != nil || month <= 0 || month > 12 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(req.FormValue("userId"))
	if err != nil || userId < 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	t := time.Date(year, time.Month(month), 0, 0, 0, 0, 0, time.UTC)

	link, err := GenerateReport(userId, &t)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(&ReportLink{link}); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

type customFS struct {
	http.FileSystem
}

// @Summary Get a report
// @Tags report
// @Description A method for getting a report
// @ID Report
// @Produce text/csv
// @Param uuid path string true "uuid of the report"
// @Success 200 "OK"
// @Router /reports/{uuid} [get]
func (fs *customFS) Open(name string) (http.File, error) {
	f, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, os.ErrNotExist
	}

	return f, nil
}

func Reports(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/csv")
	http.FileServer(
		&customFS{http.Dir(config.GENERATED_DIRNAME)},
	).ServeHTTP(w, req)
}
