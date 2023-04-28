package get_request

type lesson struct {
	PrepodNameEnc  string `json:"prepodNameEnc"`
	DayDate        string `json:"dayDate"`
	AudNum         string `json:"audNum"`
	DisciplName    string `json:"disciplName"`
	BuildNum       string `json:"buildNum"`
	OrgUnitName    string `json:"orgUnitName"`
	DayTime        string `json:"dayTime"`
	DayNum         string `json:"dayNum"`
	Potok          string `json:"potok"`
	PrepodName     string `json:"prepodName"`
	DisciplNum     string `json:"disciplNum"`
	OrgUnitId      string `json:"orgUnitId"`
	PrepodLogin    string `json:"prepodLogin"`
	DisciplType    string `json:"disciplType"`
	DisciplNameEnc string `json:"disciplNameEnc"`
}

type Schedule struct {
	Day3 []lesson `json:"3,omitempty"`
	Day2 []lesson `json:"2,omitempty"`
	Day1 []lesson `json:"1,omitempty"`
	Day6 []lesson `json:"6,omitempty"`
	Day5 []lesson `json:"5,omitempty"`
	Day4 []lesson `json:"4,omitempty"`
}
