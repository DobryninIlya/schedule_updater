package service

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"log"
	"main/internal/app/connector"
	"main/internal/app/get_request"
	"main/internal/app/model"
	"main/internal/app/store"
	"os"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"time"
)

const deltaDays = 1

var oldDateUpdate = time.Now().AddDate(-1, 0, 0)         //last year
var newDateUpdate = time.Now().AddDate(0, 0, -deltaDays) //week date ago

type Updater struct {
	conn          *pgx.Conn
	groupsParsed  []get_request.GroupInfo
	ScheduleSaved []SavedSchedule
	timeout       time.Duration
	n             *Notifier
	apiKey        string
	connector     connector.ConnectorInterface
	//store         *graph.Store
	ctx context.Context
}

func (s *Updater) getNotifyList(group int, source string) []int64 {
	sql := ""
	if source == "vk" {
		sql = "SELECT id_vk FROM users JOIN notify_clients AS n ON users.id_vk = n.destination_id WHERE n.schedule_change = true AND source=$2 AND groupp=$1" // Получение айди клиентов для рассылки в ВК
	} else if source == "tg" {
		sql = "SELECT tg_users.id FROM tg_users JOIN notify_clients AS n ON tg_users.id = n.destination_id WHERE n.schedule_change = true AND source=$2 AND groupid=$1"
	} else {
		log.Printf("Неизвестный источник рассылки.")
		return nil
	}
	rows, err := s.conn.Query(sql, group, source)
	defer rows.Close()
	if err != nil {
		log.Printf("Ошибка получения списков оповещения: %v", err)
	}
	var vkList []int64
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			log.Printf("Ошибка обработки списков оповещения: %v", err)
		}
		vkList = append(vkList, id)

	}
	return vkList
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
	rows, err := s.conn.Query("select groupp, date_update, shedule from saved_timetable where date_update > $1 and date_update< $2 and groupp>1 and groupp < 15000 ORDER BY date_update", oldDateUpdate, newDateUpdate)
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

func (s *Updater) UpdateSavedGroups(groups []model.Group) error {
	// Приведение слайса из групп к формату json
	groupsJson, err := json.Marshal(groups)
	_, err = s.conn.Exec("UPDATE saved_timetable SET shedule = $1, date_update=Now() WHERE groupp = 1", groupsJson)
	if err != nil {
		return err
	}
	return nil
}

func (s *Updater) UpdateSchedule() {
	wg := sync.WaitGroup{}
	log.Printf("Групп к обновлению: %v", len(s.ScheduleSaved))
	for i, group := range s.ScheduleSaved { //
		if group.group < 1000 {
			continue
		} // Итерация по устаревшему расписанию
		//newSchedule := get_request.GetScheduleByGroup(get_request.GroupInfo{
		//	Id:    group.group,
		//	Group: "",
		//})
		newSchedule, err := s.connector.GetScheduleByGroup(strconv.Itoa(group.group))
		if err != nil {
			log.Printf("Ошибка получения расписания для группы %v: %v", group.group, err)
		}

		shedUnmarshaled := get_request.GetUnmarshaledSchedule(group.Schedule)
		newShedMarshaled := get_request.GetMarshaledSchedule(newSchedule)
		nullSchedule := get_request.Schedule{}
		if reflect.DeepEqual(newSchedule, nullSchedule) {
			log.Printf("Полученное расписание у группы %v оказалось пустым. Обновление не произошло", group.group)
			continue
		}
		if !reflect.DeepEqual(shedUnmarshaled, newSchedule) { // Если расписание изменилось, обновляем его в базе данных
			//s.conn.QueryRow("UPDATE saved_timetable SET shedule = $1, date_update=Now() WHERE groupp = $2", newShedMarshaled, group.group)
			notification := get_request.CompareSchedules(shedUnmarshaled, newSchedule)
			sql := "UPDATE saved_timetable SET shedule = $1, date_update=Now() WHERE groupp = $2"
			_, err := s.conn.Exec(sql, newShedMarshaled, group.group)
			if err != nil {
				log.Printf("Ошибка обновления расписания: %v", err)
			}
			log.Printf("Обновлено расписание группы %v", group.group)
			wg.Add(1)
			go func() { // Блокируемся на заданный интервал, чтобы сервис не заблочил за множественные запросы и оповещаем по спискам
				time.Sleep(s.timeout)
				wg.Done()
			}()
			tgList := s.getNotifyList(group.group, "tg")
			vkList := s.getNotifyList(group.group, "vk")
			s.n.NotifyByList(vkList, tgList, group.group, notification)

			wg.Wait()
		}
		log.Printf("Осталось обновить: %v из %v", i, len(s.ScheduleSaved))
	}
}

func (s *Updater) reduceNewScheduleData(list []model.Group) []model.Group { // Убираем из полного списка то, что уже содержится в БД
	savedDataList := make([]int, s.getAllGroupsCount())
	rows, err := s.conn.Query("SELECT groupp FROM saved_timetable WHERE date_update > $1 and groupp>1 and groupp<20000", oldDateUpdate)
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
			listGroupIdInt, _ := strconv.Atoi(list[i].Id)
			if groupId == listGroupIdInt {
				list = append(list[:i], list[i+1:]...)
				continue
			}
		}
	}
	return list
}

func (s *Updater) UpdateNewSchedule(data []model.Group) { // Обновляем то, чего нет в БД (сохраненного ранее, пользователи не использовали больше года)
	for _, group := range data {
		schedule, err := s.connector.GetScheduleByGroup(group.GroupNum)
		if err != nil {
			log.Printf("Ошибка получения расписания для группы %v: %v", group.GroupNum, err)
		}
		marshaledSchedule := get_request.GetMarshaledSchedule(schedule) // TODO Начать здесь, сделать вставку расписания в БД + проверка insert/update
		_, err = s.conn.Exec("INSERT INTO saved_timetable (groupp, date_update, shedule) VALUES ($1, Now(), $2)", group.Id, marshaledSchedule)
		if err != nil {
			_, err := s.conn.Exec("UPDATE saved_timetable SET shedule = $1, date_update=Now() WHERE groupp = $2", marshaledSchedule, group.GroupNum)
			//var isUpdated bool
			//err := s.conn.QueryRow("SELECT update_saved_timetable($2, 24293);", marshaledSchedule, group.Group).Scan(&isUpdated)
			if err != nil {
				log.Printf("Ошибка обновления расписания новых групп: %v, %v", group.GroupNum, group.Id)
				log.Printf(err.Error())
			}
		}
		//s.conn.QueryRow("UPDATE saved_timetable SET shedule = $1, date_update=Now() WHERE groupp = $2", marshaledSchedule, group.Group)
		log.Printf("Обновлено расписание группы %v", group.GroupNum)
		//time.Sleep(s.timeout)
	}
}

func (s *Updater) Run() {
	defer s.conn.Close()
	//parsedRes := get_request.GetGroupsList()
	parsedRes, err := s.connector.GetGroups()
	if err != nil {
		log.Printf("Ошибка получения списка групп: %v", err)
	}
	s.ScheduleSaved = s.CollectGroups()
	err = s.UpdateSavedGroups(parsedRes) // Обновляем группы
	if err != nil {
		log.Printf("Ошибка обновления групп: %v", err)
	}
	s.UpdateSchedule() // Обновляем группы, которыми пользовались недавно
	reducedGroupsList := s.reduceNewScheduleData(parsedRes)
	s.UpdateNewSchedule(reducedGroupsList) //Добавляем в БД инфо о группах, которые еще не пользовались
	//fmt.Println(parsedRes)

}

func NewUpdater(ctx context.Context, timeout time.Duration, pgConfig pgx.ConnConfig, config *store.Config) *Updater {
	conn, err := pgx.Connect(pgConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return &Updater{
		conn:      conn,
		ctx:       ctx,
		timeout:   timeout,
		n:         NewNotifier(),
		apiKey:    config.ServiceApiKey,
		connector: connector.NewConnector(config.DataGateway, config.Debug, config.ServiceApiKey),
		//store:   graph.NewGraphStore(db),
	}
}
