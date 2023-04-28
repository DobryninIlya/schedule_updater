package service

import (
	"context"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"log"
	getRequest "main/internal/get_request"
	"os"
	"reflect"
	"sort"
	"time"
)

var oldDateUpdate = time.Now().AddDate(-1, 0, 0) //last year
var newDateUpdate = time.Now().AddDate(0, 0, -7) //week date ago

type Updater struct {
	conn          *pgx.Conn
	groupsParsed  []getRequest.GroupInfo
	ScheduleSaved []SavedSchedule
	timeout       time.Duration
	ctx           context.Context
}

func (s *Updater) getGroupsCount() int {
	var count int
	err := s.conn.QueryRow("SELECT COUNT(*) FROM saved_timetable WHERE date_update > $1 and date_update< $2 and groupp>1", oldDateUpdate, newDateUpdate).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count
}

func (s *Updater) getAllGroupsCount() int {
	var count int
	err := s.conn.QueryRow("SELECT COUNT(*) FROM saved_timetable WHERE date_update > $1 and groupp>1", oldDateUpdate).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count
}

func sortSavedSchedule(data []SavedSchedule) []SavedSchedule {
	sort.Slice(data, func(i, j int) bool {
		return data[i].group > data[j].group
	})
	return data
}

func (s *Updater) CollectGroups() []SavedSchedule {
	var (
		group      int
		dateUpdate pgtype.Date
		Schedule   []byte
	)
	count := s.getGroupsCount()
	list := make([]SavedSchedule, count)
	if count == 0 {
		return list
	}
	rows, err := s.conn.Query("select groupp, date_update, shedule from saved_timetable where date_update > $1 and date_update< $2 and groupp>1 ORDER BY date_update", oldDateUpdate, newDateUpdate)
	defer rows.Close()
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	i := 0
	for rows.Next() {
		err := rows.Scan(&group, &dateUpdate, &Schedule)
		if err != nil {
			log.Println(err)
		}
		list[i] = SavedSchedule{
			group:      group,
			dateUpdate: dateUpdate,
			Schedule:   Schedule,
		}
		i++
	}
	sortedList := sortSavedSchedule(list)
	return sortedList
}

func (s *Updater) UpdateSchedule() {
	for _, group := range s.ScheduleSaved { // Итерация по устаревшему расписанию
		newSchedule := getRequest.GetScheduleByGroup(getRequest.GroupInfo{
			Id:    group.group,
			Group: "",
		})

		shedUnmarshaled := getRequest.GetUnmarshaledSchedule(group.Schedule)
		newShedMarshaled := getRequest.GetMarshaledSchedule(newSchedule)
		if !reflect.DeepEqual(shedUnmarshaled, newSchedule) { // Если расписание изменилось, обновляем его в базе данных
			s.conn.QueryRow("UPDATE saved_timetable SET shedule = $1, date_update=Now() WHERE groupp = $2", newShedMarshaled, group.group)
			log.Printf("Обновлено расписание группы %v", group.group)
			time.Sleep(s.timeout)
			// TODO добавить оповещатор для студентов здесь

		}
	}
}

func (s *Updater) reduceNewScheduleData(list []getRequest.GroupInfo) []getRequest.GroupInfo { // Убираем из полного списка то, что уже содержится в БД
	savedDataList := make([]int, s.getAllGroupsCount())
	rows, err := s.conn.Query("SELECT groupp FROM saved_timetable WHERE date_update > $1 and groupp>1", oldDateUpdate)
	if err != nil {
		log.Printf("Ошибка при причесывании списка для внесения в БД: \n %v", err)
		return nil
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		var group int
		rows.Scan(&group)
		savedDataList[i] = group
		i++
	}
	for _, groupId := range savedDataList { // Убираем из полного списка групп те, которые есть в БД (в пределах даты)
		for i := 0; i < len(list); i++ {
			if groupId == list[i].Id {
				list = append(list[:i], list[i+1:]...)
				continue
			}
		}
	}
	return list
}

func (s *Updater) UpdateNewSchedule(data []getRequest.GroupInfo) { // Обновляем то, чего нет в БД (сохраненного ранее, пользователи не использовали больше года)
	for _, group := range data {
		schedule := getRequest.GetScheduleByGroup(group)
		marshaledSchedule := getRequest.GetMarshaledSchedule(schedule) // TODO Начать здесь, сделать вставку расписания в БД + проверка insert/update
		_, err := s.conn.Exec("INSERT INTO saved_timetable (groupp, date_update, shedule) VALUES ($1, Now(), $2)", group.Id, marshaledSchedule)
		if err != nil {
			_, err := s.conn.Exec("UPDATE saved_timetable SET shedule = $1, date_update=Now() WHERE groupp = $2", marshaledSchedule, group.Group)
			if err != nil {
				log.Printf("Ошибка обновления расписания новых групп: %v, %v", group.Group, group.Id)
				log.Printf(err.Error())
			}
		}
		//s.conn.QueryRow("UPDATE saved_timetable SET shedule = $1, date_update=Now() WHERE groupp = $2", marshaledSchedule, group.Group)
		log.Printf("Обновлено расписание группы %v", group.Group)
		time.Sleep(s.timeout)
	}
}

func (s *Updater) Run() {
	defer s.conn.Close()
	parsedRes := getRequest.GetGroupsList()
	s.ScheduleSaved = s.CollectGroups()
	s.UpdateSchedule() // Обновляем группы, которыми пользовались недавно
	reducedGroupsList := s.reduceNewScheduleData(parsedRes)
	s.UpdateNewSchedule(reducedGroupsList) //Добавляем в БД инфо о группах, которые еще не пользовались
	//fmt.Println(parsedRes)

}

func NewUpdater(ctx context.Context, timeout time.Duration, pgConfig pgx.ConnConfig) *Updater {
	conn, err := pgx.Connect(pgConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return &Updater{
		conn:    conn,
		ctx:     ctx,
		timeout: timeout,
	}
}
