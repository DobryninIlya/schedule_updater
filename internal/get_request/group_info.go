package get_request

type GroupInfo struct {
	Id    int    `json:"id"`
	Group string `json:"group"`
	Form  string `json:"forma,omitempty"`
}
