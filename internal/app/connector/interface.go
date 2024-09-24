package connector

import (
	"main/internal/app/model"
)

type ConnectorInterface interface {
	GetGroups() ([]model.Group, error)
	GetFaculties() ([]model.Faculty, error)
	GetBuildings() ([]model.Building, error)
	GetAuditoriums() ([]model.Room, error)
	GetScheduleByGroup(groupnum string) ([]model.GroupSchedule, error)
	GetKafedras() ([]model.Kafedra, error)
}
