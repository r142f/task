package main

import (
	"backend-trainee-assignment-2023/db"
	"backend-trainee-assignment-2023/segments"
	"backend-trainee-assignment-2023/users_segments"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	db.Init(false)
	defer db.DB.Close()

	http.HandleFunc("/createSegment", segments.CreateSegment)
	http.HandleFunc("/deleteSegment", segments.DeleteSegment)
	http.HandleFunc("/updateUserSegments", users_segments.UpdateUserSegments)
	http.HandleFunc("/userSegments", users_segments.UserSegments)

	http.ListenAndServe(":8080", nil)
}
