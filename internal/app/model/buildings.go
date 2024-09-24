package model

type BuildingsAnswer struct {
	Data []Building `json:"buildings"`
}

type Building struct {
	Id    string `json:"id"`
	Title string `json:"Title"`
}

type ExportBuilding struct {
	BuildingAddressText string `json:"buildingAddressText,omitempty"`
	BuildingId          string `json:"buildingId,omitempty"`
	City                string `json:"cit,omitemptyy"`
	Corps               string `json:"corps,omitempty"`
	Entrance            string `json:"entrance,omitempty"`
	House               string `json:"house,omitempty"`
	Latitude            string `json:"latitude,omitempty"`
	Locality            string `json:"locality,omitempty"`
	Longitude           string `json:"longitude,omitempty"`
	Microdistrict       string `json:"microdistrict,omitempty"`
	Region              string `json:"region,omitempty"`
	Street              string `json:"street,omitempty"`
	Structure           string `json:"structure,omitempty"`
}
