package service

import "github.com/jackc/pgx/pgtype"

type SavedSchedule struct {
	group      int
	dateUpdate pgtype.Date
	Schedule   []byte
}
