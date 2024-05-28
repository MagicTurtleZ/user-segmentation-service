package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	db *pgx.Conn
}

func New(storageUrl string) (*Storage, error) {
	const op = "storage.postgre.New"
	db, err := pgx.Connect(context.Background(), storageUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close(ctx context.Context) {
	s.db.Close(ctx)
}

func (s *Storage) CreateSegment(segName string) error {
	const op = "storage.postgre.CreateSegment"

	_, err := s.db.Exec(context.Background(), "INSERT INTO segments(segment_name) VALUES($1)", segName)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (s *Storage) CreateUser() error {
	const op = "storage.postgre.CreateUser"

	_, err := s.db.Exec(context.Background(), "INSERT INTO users(user_id) VALUES()")
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (s *Storage) DeleteSegment(segName string) error {
	const op = "storage.postgre.DeleteSegment"

	_, err := s.db.Exec(context.Background(), "DELETE FROM segments WHERE segment_name = $1", segName)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (s *Storage) AddSegment(userId int, segName string, ttl *time.Time) error {
	const op = "storage.postgre.AddSegment"
	
	_, err := s.db.Exec(context.Background(), "INSERT INTO user_segments (user_id, segment_id, time_to_limit) SELECT $1, s.segment_id, $3 FROM segments s WHERE s.segment_name = $2;", userId, segName, ttl)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (s *Storage) RemoveSegment(userId int, segName string) error {
	const op = "storage.postgre.RemoveSegment"
	_, err := s.db.Exec(context.Background(), "DELETE FROM user_segments USING segments WHERE user_segments.segment_id = segments.segment_id AND user_segments.user_id = $1 AND segments.segment_name = $2;", userId, segName)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (s *Storage) GetAllSegments(userId int) ([]string, error) {
	const op = "storage.postgre.GetAllSegments"

	rows, err := s.db.Query(context.Background(), "SELECT s.segment_name, us.time_to_limit FROM user_segments us JOIN segments AS s ON s.segment_id = us.segment_id WHERE us.user_id = $1", userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()
	res := make([]string, 0, 10)
	var (
		val string
		ddl sql.NullTime
	)
	for rows.Next() {
		err := rows.Scan(&val, &ddl)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}
		
		if !ddl.Valid || ddl.Time.After(time.Now()) {
			res = append(res, val)
		} 
	}

	return res, nil
}

func (s *Storage) GetAudit(year, month int) ([][]string, error) {
	const op = "storage.postgre.GetAudit"

	rows, err := s.db.Query(context.Background(), "SELECT * FROM segment_audit WHERE EXTRACT(YEAR FROM created) = $1 AND EXTRACT(MONTH FROM created) = $2", year, month)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()

	res := make([][]string, 0, 10)
	val3 := time.Time{}

	for rows.Next() {
		fields := make([]string, 4)
		err := rows.Scan(&fields[0], &fields[1], &fields[2], &val3)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}
		fields[3] = val3.Format(time.DateTime)
		res = append(res, fields)
	}
	
	return res, nil
}

func (s *Storage) Clean() error {
	const op = "storage.postgre.Clean"

	_, err := s.db.Exec(context.Background(), "DELETE FROM user_segments WHERE time_to_limit IS NOT NULL AND time_to_limit < CURRENT_TIMESTAMP;")
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}