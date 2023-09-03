package get_request

import (
	"fmt"
	"reflect"
	"strings"
)

func CompareSchedules(oldSchedule, newSchedule Schedule) string {
	changes := ""
	if !reflect.DeepEqual(oldSchedule.Day1, newSchedule.Day1) {
		changesDay := getChanges(oldSchedule.Day1, newSchedule.Day1)
		if changesDay != "" {
			changes += "* ### Изменения в расписании на понедельник:\n*"
			changes += changesDay
		}
	}
	if !reflect.DeepEqual(oldSchedule.Day2, newSchedule.Day2) {
		changesDay := getChanges(oldSchedule.Day2, newSchedule.Day2)
		if changesDay != "" {
			changes += "* ### Изменения в расписании на вторник:\n*"
			changes += changesDay
		}
	}
	if !reflect.DeepEqual(oldSchedule.Day3, newSchedule.Day3) {
		changesDay := getChanges(oldSchedule.Day3, newSchedule.Day3)
		if changesDay != "" {
			changes += "* ### Изменения в расписании на среду:\n*"
			changes += changesDay
		}
	}
	if !reflect.DeepEqual(oldSchedule.Day4, newSchedule.Day4) {
		changesDay := getChanges(oldSchedule.Day4, newSchedule.Day4)
		if changesDay != "" {
			changes += "* ### Изменения в расписании на четверг:\n*"
			changes += changesDay
		}
	}
	if !reflect.DeepEqual(oldSchedule.Day5, newSchedule.Day5) {
		changesDay := getChanges(oldSchedule.Day5, newSchedule.Day5)
		if changesDay != "" {
			changes += "* ### Изменения в расписании на пятницу:\n*"
			changes += changesDay
		}
	}
	if !reflect.DeepEqual(oldSchedule.Day6, newSchedule.Day6) {
		changesDay := getChanges(oldSchedule.Day6, newSchedule.Day6)
		if changesDay != "" {
			changes += "* ### Изменения в расписании на субботу:\n*"
			changes += changesDay
		}
	}
	if changes == "" {
		return "Изменений в расписании не было"
	}
	return changes
}

func getChanges(oldLessons, newLessons []lesson) string {
	changes := ""
	for _, oldLesson := range oldLessons {
		removed := true
		for _, newLesson := range newLessons {
			if oldLesson.DayDate == newLesson.DayDate && oldLesson.DayTime == newLesson.DayTime && oldLesson.AudNum == newLesson.AudNum && oldLesson.DisciplName == newLesson.DisciplName {
				removed = false
				if oldLesson.DisciplType != newLesson.DisciplType || oldLesson.PrepodNameEnc != newLesson.PrepodNameEnc || oldLesson.OrgUnitName != newLesson.OrgUnitName ||
					oldLesson.DayNum != newLesson.DayNum || oldLesson.Potok != newLesson.Potok || oldLesson.PrepodName != newLesson.PrepodName ||
					oldLesson.DisciplNum != newLesson.DisciplNum || oldLesson.OrgUnitId != newLesson.OrgUnitId || oldLesson.DisciplType != newLesson.DisciplType ||
					oldLesson.DisciplNameEnc != newLesson.DisciplNameEnc { // Да, ужасно, но что поделать :(
					changes += fmt.Sprintf("(x) Занятие \"%s\" - %s - %s - обновлена информация\n", strings.TrimSpace(oldLesson.DisciplName), strings.TrimSpace(oldLesson.DayDate), strings.TrimSpace(oldLesson.DayTime))
				}
				break
			}
		}
		if removed {
			changes += fmt.Sprintf("(-) Занятие \"%s\" - %s - %s было удалено\n", strings.TrimSpace(oldLesson.DisciplName), strings.TrimSpace(oldLesson.DayDate), strings.TrimSpace(oldLesson.DayTime))
		}
	}
	for _, newLesson := range newLessons {
		added := true
		for _, oldLesson := range oldLessons {
			if oldLesson.DayDate == newLesson.DayDate && oldLesson.DayTime == newLesson.DayTime && oldLesson.AudNum == newLesson.AudNum && oldLesson.DisciplName == newLesson.DisciplName {
				added = false
				break
			}
		}
		if added {
			changes += fmt.Sprintf("(+) Занятие \"%s\" - %s - %s было добавлено\n", strings.TrimSpace(newLesson.DisciplName), strings.TrimSpace(newLesson.DayDate), strings.TrimSpace(newLesson.DayTime))
		}
	}
	return changes
}
