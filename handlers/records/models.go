package records

import (
	"backend-trainee-assignment-2023/config"
	"database/sql"
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type ReportLink struct {
	Link string `json:"link"`
}

type Record struct {
	RecordId    int
	UserId      int
	SegmentName string
	Operation   string
	Time        time.Time
}

func writeRecordRowsToCSV(rows *sql.Rows, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	for rows.Next() {
		record := &Record{}

		if err := rows.Scan(
			&record.RecordId,
			&record.UserId,
			&record.SegmentName,
			&record.Operation,
			&record.Time,
		); err != nil {
			return err
		}

		if err := csvWriter.Write([]string{
			strconv.Itoa(record.UserId),
			record.SegmentName,
			record.Operation,
			record.Time.Format("2006-01-02 15:04:05"),
		}); err != nil {
			return err
		}
	}

	return nil
}

func GenerateReport(userId int, t *time.Time) (string, error) {
	rows, err := config.DB.Query("SELECT * FROM Records WHERE Time >= $1 AND UserId=$2", t, userId)
	if err != nil {
		return "", err
	}

	link := filepath.Join(config.REPORTS_DIRNAME, uuid.New().String()+".csv")
	err = writeRecordRowsToCSV(rows, filepath.Join(config.GENERATED_DIRNAME, link))

	return link, err
}
