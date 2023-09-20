package main

import (
	"backend-trainee-assignment-2023/config"
	"backend-trainee-assignment-2023/handlers/records"
	"backend-trainee-assignment-2023/handlers/segments"
	"backend-trainee-assignment-2023/handlers/users_segments"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	config.Init(false)
	defer config.DB.Close()

	http.HandleFunc("/createSegment", segments.CreateSegment)
	http.HandleFunc("/deleteSegment", segments.DeleteSegment)
	http.HandleFunc("/updateUserSegments", users_segments.UpdateUserSegments)
	http.HandleFunc("/userSegments", users_segments.UserSegments)
	http.HandleFunc("/generateReport", records.Report)
	http.HandleFunc(fmt.Sprintf("/%v/", config.REPORTS_DIRNAME), records.Reports)

	http.ListenAndServe(":8080", nil)
}
