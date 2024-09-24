package model

type ScheduleEvent struct {
	EventDateEnd       int64    `json:"eventDateEnd"`   // 2021-09-01 00:00:00
	EventDateStart     int64    `json:"eventDateStart"` // 2021-09-01 00:00:00
	EventId            string   `json:"eventId"`        // id
	EventLinkCall      string   `json:"eventLinkCall"`
	EventLinkMaterials string   `json:"eventLinkMaterials"`
	EventName          string   `json:"eventName"`
	EventType          string   `json:"eventType"` // лекция, практика, лабораторная
	GroupId            []string `json:"groupId"`
	RecurrenceDateEnd  int      `json:"recurrenceDateEnd"`
	RoomIds            []string `json:"roomIds"`
	TeacherIds         []string `json:"teacherIds"`
	WeeklyRecurrence   int      `json:"weeklyRecurrence"`
}
