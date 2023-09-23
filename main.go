package main

import (
	"backend-trainee-assignment-2023/config"
	"backend-trainee-assignment-2023/handlers/records"
	"backend-trainee-assignment-2023/handlers/segments"
	"backend-trainee-assignment-2023/handlers/users_segments"
	"fmt"
	"net/http"

	_ "backend-trainee-assignment-2023/docs"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title TITLE
// @description DESCRIPTION

// @host localhost:8080
func main() {
	config.Init(false)
	defer config.DB.Close()

	http.HandleFunc("/createSegment", segments.CreateSegment)
	http.HandleFunc("/deleteSegment", segments.DeleteSegment)
	http.HandleFunc("/updateUserSegments", users_segments.UpdateUserSegments)
	http.HandleFunc("/userSegments", users_segments.UserSegments)
	http.HandleFunc("/generateReport", records.Report)
	http.HandleFunc(fmt.Sprintf("/%v/", config.REPORTS_DIRNAME), records.Reports)
	http.HandleFunc("/docs/", httpSwagger.WrapHandler)

	http.ListenAndServe(":8080", nil)
}
