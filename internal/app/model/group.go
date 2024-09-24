package model

type GroupAnswer struct {
	Data []Group `json:"groups"`
}

type Group struct {
	Id       string `json:"id"`
	GroupNum string `json:"groupNum"`
	SpecNum  string `json:"specNum"`
}

type ExportGroup struct {
	CourseNumber   int    `json:"courseNumber,omitempty"`
	EducationForm  string `json:"educationForm,omitempty"`
	EducationLevel string `json:"educationLevel,omitempty"`
	FacultyId      string `json:"facultyId,omitempty"`
	GroupCode      string `json:"groupCode,omitempty"`
	GroupId        string `json:"groupId,omitempty"`
	GroupName      string `json:"groupName,omitempty"`
	SpecialityCode string `json:"specialityCode,omitempty"`
}
