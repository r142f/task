package users_segments

import (
	"backend-trainee-assignment-2023/config"
	"backend-trainee-assignment-2023/handlers/segments"
	"database/sql"
	"errors"
	"net/http"
)

type UserSegments struct {
	SegmentsToAdd    []string `json:"segmentsToAdd" example:"DISCOUNT_30"`
	SegmentsToDelete []string `json:"segmentsToDelete" example:"VOICE_MESSAGES"`
	UserId           int	  `json:"userId" example:"1"`
}

func InsertUserSegmentWithSegmentName(tx *sql.Tx, userId int, segmentName string) error {
	segmentId, err := segments.SelectSegmentIdBySegmentName(segmentName)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO UsersSegments (UserId, SegmentId) VALUES ($1, $2);", userId, segmentId)

	return err
}

func DeleteUserSegmentWithSegmentName(tx *sql.Tx, userId int, segmentName string) error {
	segmentId, err := segments.SelectSegmentIdBySegmentName(segmentName)
	if err != nil {
		return err
	}

	res, err := tx.Exec("DELETE FROM UsersSegments WHERE UserId=$1 AND SegmentId=$2;", userId, segmentId)
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New(http.StatusText(http.StatusBadRequest))
	}

	return err
}

func SelectUserSegments(userId int) (segmentNames []string, err error) {
	segmentNames = make([]string, 0)
	rows, err := config.DB.Query("SELECT SegmentName FROM UsersSegments NATURAL JOIN Segments WHERE UserId=$1", userId)
	if err != nil {
		return
	}

	for rows.Next() {
		var segmentName string

		if err = rows.Scan(&segmentName); err != nil {
			return
		}

		segmentNames = append(segmentNames, segmentName)
	}

	err = rows.Err()

	return
}
