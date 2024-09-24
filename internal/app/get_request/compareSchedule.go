package get_request

import (
	"fmt"
	"main/internal/app/model"
	"reflect"
	"strings"
)

func CompareSchedules(oldSchedule, newSchedule []model.GroupSchedule) string {
	var changes string
	days := []struct {
		oldDay, newDay []model.GroupSchedule
		dayName        string
	}{
		{oldSchedule, newSchedule, "понедельник"},
		{oldSchedule, newSchedule, "вторник"},
		{oldSchedule, newSchedule, "среду"},
		{oldSchedule, newSchedule, "четверг"},
		{oldSchedule, newSchedule, "пятницу"},
		{oldSchedule, newSchedule, "субботу"},
	}

	for _, day := range days {
		if !reflect.DeepEqual(day.oldDay, day.newDay) {
			changesDay := getChanges(day.oldDay, day.newDay)
			if changesDay != "" {
				changes += fmt.Sprintf("* ### Изменения в расписании на %s:\n*%s", day.dayName, changesDay)
			}
		}
	}

	if changes == "" {
		return "Изменений в расписании не было"
	}
	return changes
}

func getChanges(oldLessons, newLessons []model.GroupSchedule) string {
	var changes string
	for _, oldLesson := range oldLessons {
		removed := true
		for _, newLesson := range newLessons {
			if oldLesson.Daydate == newLesson.Daydate && oldLesson.Daytime == newLesson.Daytime && oldLesson.Auditory == newLesson.Auditory && oldLesson.Disciplname == newLesson.Disciplname {
				removed = false
				if oldLesson.Discipltype != newLesson.Discipltype || oldLesson.Prepodfio != newLesson.Prepodfio || oldLesson.KafTitle != newLesson.KafTitle ||
					oldLesson.Daynum != newLesson.Daynum || oldLesson.Building != newLesson.Building || oldLesson.PrepodID != newLesson.PrepodID ||
					oldLesson.Disciplnum != newLesson.Disciplnum || oldLesson.KafId != newLesson.KafId || oldLesson.FacId != newLesson.FacId ||
					oldLesson.AuditoryId != newLesson.AuditoryId {
					changes += fmt.Sprintf("(x) Занятие \"%s\" - %s - %s - обновлена информация\n", strings.TrimSpace(oldLesson.Disciplname), strings.TrimSpace(oldLesson.Daydate), strings.TrimSpace(oldLesson.Daytime))
				}
				break
			}
		}
		if removed {
			changes += fmt.Sprintf("(-) Занятие \"%s\" - %s - %s было удалено\n", strings.TrimSpace(oldLesson.Disciplname), strings.TrimSpace(oldLesson.Daydate), strings.TrimSpace(oldLesson.Daytime))
		}
	}
	for _, newLesson := range newLessons {
		added := true
		for _, oldLesson := range oldLessons {
			if oldLesson.Daydate == newLesson.Daydate && oldLesson.Daytime == newLesson.Daytime && oldLesson.Auditory == newLesson.Auditory && oldLesson.Disciplname == newLesson.Disciplname {
				added = false
				break
			}
		}
		if added {
			changes += fmt.Sprintf("(+) Занятие \"%s\" - %s - %s было добавлено\n", strings.TrimSpace(newLesson.Disciplname), strings.TrimSpace(newLesson.Daydate), strings.TrimSpace(newLesson.Daytime))
		}
	}
	return changes
}
