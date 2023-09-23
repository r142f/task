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
