package csvgen

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

type AuditGetter interface {
	GetAudit(year, month int) ([][]string, error)
}

func GenReport(month, year int, ag AuditGetter) (string, error) {
	const op = "csv.create.GenPeport"

	report, err := ag.GetAudit(year, month)
	if err != nil {
		return "", fmt.Errorf("%s: %v", op, err)
	}
	filename := strings.Builder{}
	filename.WriteString(`internal\CSV\history\report_`)
	filename.WriteString(takeTime())
	filename.WriteString(".csv")

	if err := saveReport(report, filename.String()); err != nil {
		return "", fmt.Errorf("%s: %v", op, err)
	}

	return filename.String(), nil
}

func saveReport(report [][]string, filename string) error {
	const op = "csv.create.saveReport"

	file, err := os.Create(filename)
	
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	w := csv.NewWriter(file)
	defer w.Flush()

	err = w.WriteAll(report)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}


func takeTime() string {
	date := strings.Builder{}
	now := time.Now()
	date.WriteString(now.Format(time.DateOnly))
	date.WriteByte('_')
	date.WriteString(now.Format(time.TimeOnly))
	tmp := strings.ReplaceAll(date.String(), "-", "_")
	tmp = strings.ReplaceAll(tmp, ":", "_")
	return tmp
}