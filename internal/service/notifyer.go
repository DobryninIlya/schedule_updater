package service

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

const vkTemplate = "https://api.vk.com/method/%s?v=5.131&%s"
const vkSendMethod = "messages.send"
const tgTemplate = "https://api.telegram.org/bot%s/%s%s"
const tgSendMethod = "sendMessage?"
const notifyMessage = "[Оповещение] Произошло обновление вашего учебного расписания."
const tgButtons = `{
"inline_keyboard": [
[{
"text": "Оповещения",
"callback_data": "{\"button\":\"notifier_change_menu\",\"data\":\"notifier_change_menu\"}"
}]
]
}`

const vkButtons = `{
"inline": true,
"buttons": [
[
{
"action": {
"type": "text",
"label": "Управление оповещениями"
},
"color": "default"
}
]
]
}`

type Notifier struct {
	vkToken    string
	tgToken    string
	vkTemplate string
	tgTemplate string
	wg         sync.WaitGroup
}

func NewNotifier() *Notifier {
	rand.Seed(time.Now().UnixNano())
	return &Notifier{
		vkToken:    os.Getenv("VK_TOKEN"),
		tgToken:    os.Getenv("TG_TOKEN"),
		vkTemplate: vkTemplate,
		tgTemplate: tgTemplate,
	}
}

func (s Notifier) SendMessageVKids(uId []int64, message string, buttons string) bool {
	if len(uId) == 0 {
		log.Printf("Попытка отправить сообщение пустому списку оповещения")
		return false
	}
	ids := ""
	for _, v := range uId {
		ids += fmt.Sprintf("%v,", v)
	}
	ids = ids[:len(ids)-1]
	randomInt := rand.Int31()
	params := fmt.Sprintf("random_id=%v&peer_ids=%s&access_token=%s&disable_mentions=0&message=%s&keyboard=%s",
		randomInt,
		ids,
		s.vkToken,
		url.QueryEscape(message),
		url.QueryEscape(buttons),
	)
	url := fmt.Sprintf(s.vkTemplate, vkSendMethod, params)
	log.Printf("URL %v", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Ошибка API. Отправка сообщений: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	log.Printf("Результатат отправки в ВК: %v, %v", body, resp.Header)
	if err != nil {
		log.Printf("При этом возникла ошибка: %v", err)
	}
	defer resp.Body.Close()
	return true
}

func (s Notifier) SendMessageTG(uId int64, message string, buttons string) bool {
	if uId == 0 {
		log.Printf("Попытка отправить сообщение некорректному айди")
		return false
	}
	params := fmt.Sprintf("chat_id=%v&text=%v&reply_markup=%s",
		uId,
		url.QueryEscape(message),
		url.QueryEscape(buttons),
	)
	url := fmt.Sprintf(s.tgTemplate, s.tgToken, tgSendMethod, params)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Ошибка API. Отправка сообщений: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	log.Printf("Результатат отправки в ВК: %v", string(body))
	if err != nil {
		log.Printf("При этом возникла ошибка: %v", err)
	}
	defer resp.Body.Close()
	return true
}

func (s Notifier) NotifyByList(listVK, listTG []int64, group int) {
	s.wg.Add(2)
	go func() {
		if len(listVK) >= 100 {
			log.Println("Список оповещения слишком большой для разовой отправки в ВК")
			s.wg.Done()
			return
		}
		result := s.SendMessageVKids(listVK, notifyMessage, vkButtons)
		if !result {
			log.Printf("Группа %v не оповещена в ВК", group)
		} else {
			log.Printf("Оповещена сообщением VK группа: %v", group)
		}
		s.wg.Done()
	}()

	go func() {
		result := false
		for _, uId := range listTG {
			result = s.SendMessageTG(uId, notifyMessage, tgButtons)
		}
		if !result {
			log.Printf("Группа %v не оповещена в ТГ", group)
		} else {
			log.Printf("Оповещена сообщением TG группа: %v", group)
		}
		s.wg.Done()
	}()
	s.wg.Wait()
}
