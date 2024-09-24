package model

type FacultiesAnswer struct {
	Data []Faculty `json:"faculties"`
}

type Faculty struct {
	FacId    string `json:"facId"`
	FacTitle string `json:"facTitle"`
}

type ExportFaculty struct {
	FacultyId   string `json:"facultyId,omitempty"`
	FacultyCode string `json:"facultyCode,omitempty"`
	FacultyName string `json:"facultyName,omitempty"`
}
