package model

type RoomsAnswer struct {
	Data []Room `json:"auditories"`
}

type Room struct {
	BuildingId   string `json:"buildingId"`
	AuditoryId   string `json:"auditoryId"`
	AuditoryName string `json:"auditoryName"`
}

type ExportRoom struct {
	BuildingId string `json:"buildingId,omitempty"`
	Floor      int    `json:"floor,omitempty"`
	RoomId     string `json:"roomId,omitempty"`
	RoomName   string `json:"roomName,omitempty"`
}
