package model

type GroupScheduleAnswer struct {
	Data []GroupSchedule `json:"schedule"`
}

type GroupSchedule struct {
	Daynum      string `json:"daynum,omitempty"`
	Dayname     string `json:"dayname,omitempty"`
	Daytime     string `json:"daytime,omitempty"`
	Daydate     string `json:"daydate,omitempty"`
	Disciplname string `json:"disciplname,omitempty"`
	Disciplnum  string `json:"disciplnum,omitempty"`
	Discipltype string `json:"discipltype,omitempty"`
	Auditory    string `json:"auditory,omitempty"`
	Building    string `json:"building,omitempty"`
	Prepodfio   string `json:"prepodfio,omitempty"`
	PrepodID    string `json:"prepodID,omitempty"`
	KafTitle    string `json:"kafTitle,omitempty"`
	KafId       string `json:"kafId,omitempty"`
	FacId       string `json:"facId,omitempty"`
	AuditoryId  string `json:"auditoryId,omitempty"`
}
