package main

import (
	"backend-trainee-assignment-2023/db"
	"backend-trainee-assignment-2023/segments"
	"backend-trainee-assignment-2023/users_segments"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	// "os"
	// "path"
	"reflect"
	"sort"
	"testing"
)

func init() {
	db.Init(true)
}

func clearDB() {
	for _, table := range db.TABLES {
		_, err := db.DB.Exec(fmt.Sprintf("DELETE FROM %v;", table))
		if err != nil {
			log.Printf("Couldn't clear table %v: %v\n", table, err)
		}
	}
}

func createSegment(segmentName string) *http.Response {
	body, _ := json.Marshal(&segments.Segment{Name: segmentName})
	req := httptest.NewRequest(http.MethodPost, "/createSegment", bytes.NewReader(body))
	w := httptest.NewRecorder()

	segments.CreateSegment(w, req)

	return w.Result()
}

func deleteSegment(segmentName string) *http.Response {
	body, _ := json.Marshal(&segments.Segment{Name: segmentName})
	req := httptest.NewRequest(http.MethodDelete, "/deleteSegment", bytes.NewReader(body))
	w := httptest.NewRecorder()

	segments.DeleteSegment(w, req)

	return w.Result()
}

func updateUserSegments(segmentsToAdd, segmentsToDelete []string, userId int) *http.Response {
	userSegments := &struct {
		SegmentsToAdd    []string
		SegmentsToDelete []string
		UserId           int
	}{segmentsToAdd, segmentsToDelete, userId}

	body, _ := json.Marshal(userSegments)
	req := httptest.NewRequest(http.MethodPost, "/updateUserSegments", bytes.NewReader(body))
	w := httptest.NewRecorder()

	users_segments.UpdateUserSegments(w, req)

	return w.Result()
}

func userSegments(userId int) *http.Response {
	uri := fmt.Sprintf("/userSegments?userId=%v", userId)
	body, _ := json.Marshal(&struct{ UserId int }{userId})
	req := httptest.NewRequest(http.MethodGet, uri, bytes.NewReader(body))
	w := httptest.NewRecorder()

	users_segments.UserSegments(w, req)

	return w.Result()
}

func getUserSegmentNames(t *testing.T, userId int) []string {
	res := userSegments(userId)
	if res.StatusCode != http.StatusOK {
		t.Errorf("%v\nexpected status code to be %v, but got %v", res.Status, http.StatusOK, res.StatusCode)
	}

	defer res.Body.Close()
	segmentNames := make([]string, 0)
	json.NewDecoder(res.Body).Decode(&segmentNames)

	return segmentNames
}

func TestCreateSegment(t *testing.T) {
	clearDB()

	segmentName := "AVITO_VOICE_MESSAGES"
	res := createSegment(segmentName)
	if res.StatusCode != http.StatusCreated {
		t.Errorf("%v\nexpected status code to be %v, but got %v", res.Status, http.StatusCreated, res.StatusCode)
	}

	res = createSegment(segmentName)
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("%v\nexpected status code to be %v, but got %v", res.Status, http.StatusCreated, res.StatusCode)
	}
}

func TestDeleteSegment(t *testing.T) {
	clearDB()

	segmentName := "AVITO_VOICE_MESSAGES"
	res := createSegment(segmentName)
	if res.StatusCode != http.StatusCreated {
		t.Errorf("%v\nexpected status code to be %v, but got %v", res.Status, http.StatusCreated, res.StatusCode)
	}

	res = deleteSegment(segmentName)
	if res.StatusCode != http.StatusOK {
		t.Errorf("%v\nexpected status code to be %v, but got %v", res.Status, http.StatusOK, res.StatusCode)
	}
}

func TestUpdateUserSegments(t *testing.T) {
	clearDB()

	segmentNames := []string{"s1", "s2", "s3", "s4"}
	for _, segmentName := range segmentNames {
		createSegment(segmentName)
	}

	var tests = []struct {
		segmentsToAdd    []string
		segmentsToDelete []string
		userId           int
		want             int
	}{
		{nil, nil, 1, http.StatusCreated},
		{[]string{"s1", "s2", "s3"}, nil, 1, http.StatusCreated},
		{[]string{"s1", "s2", "s3"}, nil, 2, http.StatusCreated},
		{[]string{"s1", "s2", "s3"}, nil, 1, http.StatusInternalServerError},
		{[]string{"s1", "s2", "s3"}, nil, 2, http.StatusInternalServerError},
		{nil, []string{"s2"}, 1, http.StatusCreated},
		{nil, []string{"s2"}, 2, http.StatusCreated},
		{[]string{"s4"}, []string{"s1", "s4"}, 1, http.StatusCreated},
		{[]string{"s4"}, []string{"s1", "s4"}, 2, http.StatusCreated},
		{[]string{"s1"}, []string{"s2", "s4"}, 1, http.StatusInternalServerError},
		{[]string{"s1"}, []string{"s2", "s4"}, 2, http.StatusInternalServerError},
	}

	for _, test := range tests {
		if res := updateUserSegments(
			test.segmentsToAdd,
			test.segmentsToDelete,
			test.userId,
		); res.StatusCode != test.want {
			t.Errorf("%v\nexpected status code to be %v, but got %v", res.Status, http.StatusOK, res.StatusCode)
		}
	}
}

func TestUserSegments(t *testing.T) {
	clearDB()

	for _, segmentName := range []string{"s1", "s2", "s3", "s4"} {
		createSegment(segmentName)
	}

	var tests = []struct {
		segmentsToAdd    []string
		segmentsToDelete []string
		userId           int
		want             []string
	}{
		{nil, nil, 1, nil},
		{[]string{"s1", "s2", "s3"}, nil, 1, []string{"s1", "s2", "s3"}},
		{nil, []string{"s2"}, 1, []string{"s1", "s3"}},
		{[]string{"s4"}, []string{"s1", "s4"}, 1, []string{"s3"}},
		{nil, []string{"s3"}, 1, nil},
	}

	for _, test := range tests {
		updateUserSegments(
			test.segmentsToAdd,
			test.segmentsToDelete,
			test.userId,
		)

		userSegments := getUserSegmentNames(t, test.userId)
		sort.Strings(userSegments)
		if !reflect.DeepEqual(userSegments, test.want) {
			t.Errorf("expected userSegments to be %v, but got %v", test.want, userSegments)
		}
	}
}

func TestUserSegmentsWithDeleteSegment(t *testing.T) {
	clearDB()

	for _, segmentName := range []string{"s1", "s2", "s3", "s4"} {
		createSegment(segmentName)
	}

	updateUserSegments(
		[]string{"s1", "s2", "s3"},
		nil,
		1,
	)
	updateUserSegments(
		[]string{"s1", "s2", "s3", "s4"},
		nil,
		2,
	)

	var tests = []struct {
		segmentToDelete string
		userIds         []int
		wants           [][]string
	}{
		{"s3", []int{1, 2}, [][]string{{"s1", "s2"}, {"s1", "s2", "s4"}}},
		{"s4", []int{1, 2}, [][]string{{"s1", "s2"}, {"s1", "s2"}}},
	}

	for _, test := range tests {
		if res := deleteSegment(test.segmentToDelete); res.StatusCode != http.StatusOK {
			t.Errorf("%v\nexpected status code to be %v, but got %v", res.Status, http.StatusOK, res.StatusCode)
		}
		for i := range test.userIds {
			userSegments := getUserSegmentNames(t, test.userIds[i])
			sort.Strings(userSegments)
			if !reflect.DeepEqual(userSegments, test.wants[i]) {
				t.Errorf("expected userSegments to be %v, but got %v", test.wants[i], userSegments)
			}
		}
	}
}
