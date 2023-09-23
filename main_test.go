package main

import (
	"backend-trainee-assignment-2023/config"
	"backend-trainee-assignment-2023/handlers/records"
	"backend-trainee-assignment-2023/handlers/segments"
	"backend-trainee-assignment-2023/handlers/users_segments"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"
)

func init() {
	config.Init(true)
}

func clearDB() {
	for _, table := range config.TABLES {
		_, err := config.DB.Exec(fmt.Sprintf("DELETE FROM %v;", table))
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
	userSegments := &users_segments.UserSegments{
		SegmentsToAdd:    segmentsToAdd,
		SegmentsToDelete: segmentsToDelete,
		UserId:           userId,
	}

	body, _ := json.Marshal(userSegments)
	req := httptest.NewRequest(http.MethodPost, "/updateUserSegments", bytes.NewReader(body))
	w := httptest.NewRecorder()

	users_segments.UpdateUserSegments(w, req)

	return w.Result()
}

func userSegments(userId int) *http.Response {
	uri := fmt.Sprintf("/userSegments?userId=%v", userId)
	req := httptest.NewRequest(http.MethodGet, uri, nil)
	w := httptest.NewRecorder()

	users_segments.SegmentsByUser(w, req)

	return w.Result()
}

func getUserSegmentNames(t *testing.T, userId int) []string {
	res := userSegments(userId)
	if res.StatusCode != http.StatusOK {
		t.Errorf("%v\nexpected status code to be %v, but got %v", res.Status, http.StatusOK, res.StatusCode)
	}

	segmentNames := make([]string, 0)
	json.NewDecoder(res.Body).Decode(&segmentNames)
	defer res.Body.Close()

	return segmentNames
}

func generateReport(year, month, userId int) *http.Response {
	uri := fmt.Sprintf("/report?year=%v&month=%v&userId=%v", year, month, userId)
	req := httptest.NewRequest(http.MethodGet, uri, nil)
	w := httptest.NewRecorder()

	records.Report(w, req)

	return w.Result()
}

func getReport(t *testing.T, year, month, userId int) []*records.Record {
	res := generateReport(year, month, userId)
	if res.StatusCode != http.StatusOK {
		t.Errorf("%v\nexpected status code to be %v, but got %v", res.Status, http.StatusOK, res.StatusCode)
	}

	resJson := &records.ReportLink{}
	json.NewDecoder(res.Body).Decode(resJson)
	defer res.Body.Close()

	uri := fmt.Sprintf("/%v", resJson.Link)
	req := httptest.NewRequest(http.MethodGet, uri, nil)
	w := httptest.NewRecorder()

	records.Reports(w, req)

	res = w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("%v\nexpected status code to be %v, but got %v", res.Status, http.StatusOK, res.StatusCode)
	}

	report := make([]*records.Record, 0)
	reader := csv.NewReader(res.Body)
	defer res.Body.Close()
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
		}

		record := &records.Record{}
		fmt.Sscan(fields[0], &record.UserId)
		fmt.Sscan(fields[1], &record.SegmentName)
		fmt.Sscan(fields[2], &record.Operation)
		fmt.Sscan(fields[3], &record.Time)

		report = append(report, record)
	}

	return report
}

func updateRecordTime(t *testing.T, year, month, userId, limit int) {
	config.DB.Exec(`
		UPDATE Records as r1 SET Time=$1
		WHERE r1.RecordId=(
			SELECT r2.RecordId
			FROM Records as r2
			WHERE r2.UserId=$2
			LIMIT $3
		);
	`, time.Date(year, time.Month(month), 0, 0, 0, 0, 0, time.UTC), userId, limit)
}

func TestCreateSegment(t *testing.T) {
	clearDB()

	segmentName := "AVITO_VOICE_MESSAGES"
	res := createSegment(segmentName)
	segment := &segments.Segment{}
	json.NewDecoder(res.Body).Decode(&segment)
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("%v\nexpected status code to be %v, but got %v", res.Status, http.StatusCreated, res.StatusCode)
	}
	if segment.Name != segmentName {
		t.Errorf("expected created segment name to be %v, but got %v", segmentName, segment.Name)
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

func TestGenerateReport(t *testing.T) {
	clearDB()

	for _, segmentName := range []string{"s1", "s2", "s3", "s4"} {
		createSegment(segmentName)
	}

	updateUserSegments(
		[]string{"s1", "s2", "s3"},
		[]string{"s2", "s3"},
		1,
	)

	updateUserSegments(
		[]string{"s1", "s2", "s3"},
		[]string{"s2", "s3"},
		2,
	)

	report := getReport(t, time.Now().Year(), int(time.Now().Month()), 1)
	if len(report) != 5 {
		t.Errorf("expected report len to be %v, but got %v", 5, len(report))
	}

	delete, add := 0, 0
	for _, record := range report {
		if record.Operation == "delete" {
			delete++
		} else {
			add++
		}
	}
	if delete != 2 || add != 3 {
		t.Errorf("expected amount of delete/add records to be %v/%v, but got %v/%v", 2, 3, delete, add)
	}

	updateRecordTime(t, time.Now().Year()-1, int(time.Now().Month()), 1, 1)
	if report = getReport(t, time.Now().Year(), int(time.Now().Month()), 1); len(report) != 4 {
		t.Errorf("expected report len to be %v, but got %v", 4, len(report))
	}

	err := os.RemoveAll(filepath.Join(config.GENERATED_DIRNAME, config.REPORTS_DIRNAME))
	if err != nil {
		t.Error(err)
	}
}
