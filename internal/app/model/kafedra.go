package model

type KafedraAnswer struct {
	Data []Kafedra `json:"kafedras"`
}

type Kafedra struct {
	KafId    string `json:"kafId,omitempty"`
	KafTitle string `json:"kafTitle,omitempty"`
	FacId    string `json:"facId,omitempty"`
}

type ExportKafedra struct {
	DepartmentId   string `json:"departmentId,omitempty"`
	DepartmentName string `json:"departmentName,omitempty"`
}
