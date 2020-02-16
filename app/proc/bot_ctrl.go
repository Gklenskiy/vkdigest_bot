package proc

import (
	"fmt"
	"strings"
	"time"

	"github.com/Gklenskiy/vkdigest_bot/app/models"
	log "github.com/go-pkgz/lgr"
	tb "gopkg.in/tucnak/telebot.v2"
)

// BotCtrl bot controller
type BotCtrl struct {
	BotCtrlSettings
}

// BotCtrlSettings settings for bot controller
type BotCtrlSettings struct {
	VkBaseURL    string
	VkAPIVersion string
	VkAppID      string
	AuthURL      string
}

// NewBotCtrl creates Bot
func NewBotCtrl(settings BotCtrlSettings) *BotCtrl {
	return &BotCtrl{
		BotCtrlSettings: settings,
	}
}

// PingCtrl handle ping command
func (ctrl *BotCtrl) PingCtrl(b *tb.Bot, m *tb.Message) {
	_, err := b.Send(m.Sender, "pong 4")
	if err != nil {
		log.Printf("[Error] while send message to user %+v", err)
	}
}

// StartCtrl handler /start command
func (ctrl *BotCtrl) StartCtrl(b *tb.Bot, m *tb.Message) {
	sendStartMsg(ctrl.VkAppID, ctrl.AuthURL, ctrl.VkAPIVersion, b, m)
}

// TrendsCtrl for handle trend command
func (ctrl *BotCtrl) TrendsCtrl(b *tb.Bot, m *tb.Message) {
	sortType := 0
	daysDeadline := 7
	deadline := time.Now().AddDate(0, 0, -daysDeadline).Unix()
	log.Printf("[DEBUG] Deadline date %v \n", time.Unix(deadline, 0))

	exist, token := ctrl.tryGetUserVkToken(b, m)
	if !exist {
		return
	}

	log.Printf("[DEBUG] Start process Trends Request")
	startProcess := time.Now()
	vkClient := NewVkClient(ctrl.VkBaseURL, token, ctrl.VkAPIVersion, time.Millisecond)

	var allPosts Posts
	log.Printf("[DEBUG] Start get Posts")
	start := time.Now()
	ch := make(chan Posts)

	userID := m.Sender.ID
	sources, err := models.GetSources(userID)
	if err != nil {
		log.Printf("[ERROR] while get sources %+v", err)
		return
	}

	if len(sources) == 0 {
		_, err := b.Send(m.Sender, "В списке источников пока ничего нет\n"+
			"Добавить страницу можно командой /add\n")
		if err != nil {
			log.Printf("[Error] while send message to user %+v", err)
		}
		return
	}

	for _, domain := range sources {
		go processDomain(vkClient, deadline, domain, ch)
	}

	for range sources {
		posts := <-ch
		allPosts = append(allPosts, posts...)
	}
	elapsed := time.Since(start)
	log.Printf("[DEBUG] End get Posts took %s", elapsed)

	allPosts.Filter(sortType)

	res := make([]string, 0)
	for _, v := range allPosts[:7] {
		postInfo := fmt.Sprintf(`https://vk.com/wall%d_%d 
		Likes: %d 
		Reposts: %d 
		Views: %d 
		Comments: %d 
		LikesOnView: %.2f 
		RepostsOnView: %.2f 
		CommentsOnView: %.2f 
		Rating: %.2f  
		Avg: %.2f
		Div: %.2f`, v.OwnerID, v.ID, v.Likes, v.Reposts, v.Views, v.Comments, v.LikesOnView(), v.RepostsOnView(), v.CommentsOnView(), v.Rating(), v.AverageRating, v.DeviationFromAverage())
		res = append(res, postInfo)
	}

	elapsedProcess := time.Since(startProcess)
	log.Printf("[DEBUG] End process Trends Request %s", elapsedProcess)
	for _, v := range res {
		_, err := b.Send(m.Sender, v)
		if err != nil {
			log.Printf("[Error] while send message to user %+v", err)
		}
	}
}

// AddCtrl handles command for add source
func (ctrl *BotCtrl) AddCtrl(b *tb.Bot, m *tb.Message) {
	exist, token := ctrl.tryGetUserVkToken(b, m)
	if !exist {
		return
	}

	commandParam := getCommandParams(m.Text)
	log.Printf("[DEBUG] Add: %s", commandParam)
	if commandParam == "" {
		_, err := b.Send(m.Sender,
			"Введи ссылку на страницу вк, например https://vk.com/newalbums \n"+
				"Либо короткое имя страницы, например: newalbums")
		if err != nil {
			log.Printf("[Error] while send message to user %+v", err)
		}
		return
	}

	domain := getDomain(commandParam)
	vkClient := NewVkClient(ctrl.VkBaseURL, token, ctrl.VkAPIVersion, time.Millisecond)
	if !vkClient.IsValidDomain(domain) {
		_, err := b.Send(m.Sender, "Ссылка недействительна \n"+
			"Введи ссылку на страницу вк, например https://vk.com/newalbums \n"+
			"Либо короткое имя страницы, например: newalbums")
		if err != nil {
			log.Printf("[Error] while send message to user %+v", err)
		}
		return
	}

	log.Printf("[DEBUG] Add: Success! Save %s", domain)
	userID := m.Sender.ID
	err := models.CreateSource(userID, domain)
	if err != nil {
		log.Printf("[ERROR] failed to save sourse %s for user %d, %+v", domain, userID, err)
		_, err := b.Send(m.Sender, "Хмммммм... Какие-то проблемы, попробуй ещё раз")
		if err != nil {
			log.Printf("[Error] while send message to user %+v", err)
		}

		return
	}

	_, err = b.Send(m.Sender, "Успех!")
	if err != nil {
		log.Printf("[Error] while send message to user %+v", err)
	}
}

//DeleteCtrl handles command for delete source
func (ctrl *BotCtrl) DeleteCtrl(b *tb.Bot, m *tb.Message) {
	exist, _ := ctrl.tryGetUserVkToken(b, m)
	if !exist {
		return
	}

	commandParam := getCommandParams(m.Text)
	log.Printf("[DEBUG] Delete: %s", commandParam)
	if commandParam == "" {
		_, err := b.Send(m.Sender, "Введите название страницы")
		if err != nil {
			log.Printf("[Error] while send message to user %+v", err)
		}
		return
	}

	domain := getDomain(commandParam)
	userID := m.Sender.ID
	rowCount, err := models.DeleteSource(userID, domain)
	if err != nil {
		log.Printf("[ERROR] failed to delete sourse %s for user %d, %+v", domain, userID, err)
		_, err := b.Send(m.Sender, "Хмммммм... Какие-то проблемы, попробуй ещё раз")
		if err != nil {
			log.Printf("[Error] while send message to user %+v", err)
		}

		return
	}
	if rowCount == 0 {
		_, err := b.Send(m.Sender, "Такой страницы нет в списке")
		if err != nil {
			log.Printf("[Error] while send message to user %+v", err)
		}

		return
	}

	_, err = b.Send(m.Sender, "Готово!")
	if err != nil {
		log.Printf("[Error] while send message to user %+v", err)
	}
}

func getCommandParams(userInput string) string {
	idx := strings.Index(userInput, " ")
	if idx > -1 {
		return userInput[idx+1:]
	} else {
		return ""
	}
}

func getDomain(commandParam string) string {
	slahIdx := strings.LastIndex(commandParam, "/")
	var domain string
	if slahIdx > -1 {
		domain = commandParam[slahIdx+1:]
	} else {
		domain = commandParam
	}
	log.Printf("[DEBUG] Domain specified %s", domain)

	return domain
}

// ListCtrl handle command for list of sources
func (ctrl *BotCtrl) ListCtrl(b *tb.Bot, m *tb.Message) {
	exist, _ := ctrl.tryGetUserVkToken(b, m)
	if !exist {
		return
	}

	userID := m.Sender.ID
	sources, err := models.GetSources(userID)
	if err != nil {
		log.Printf("[ERROR] while get sources %+v", err)
		return
	}

	if len(sources) == 0 {
		_, err := b.Send(m.Sender, "В списке источников пока ничего нет\n"+
			"Добавить страницу можно командой /add\n")
		if err != nil {
			log.Printf("[Error] while send message to user %+v", err)
		}
		return
	}

	var msg string
	for _, v := range sources {
		msg += v + "\n"
	}
	_, err = b.Send(m.Sender, msg)
	if err != nil {
		log.Printf("[Error] while send message to user %+v", err)
	}

}

func processDomain(vkClient *VkClient, deadline int64, domain string, ch chan<- Posts) {
	posts, err := vkClient.GetPosts(0, deadline, domain)
	if err != nil {
		log.Printf("[ERROR] failed to get posts from Vk, %+v", err)
	}

	posts.SetAverageRating()
	ch <- posts
}

func sendStartMsg(vkAPI string, authURL string, vkAPIVersion string, b *tb.Bot, m *tb.Message) {
	url := fmt.Sprintf(`https://oauth.vk.com/authorize?client_id=%s&display=page&redirect_uri=%s&scope=wall,offline&response_type=code&v=%s&state=%d`,
		vkAPI, authURL, vkAPIVersion, m.Sender.ID)

	inlineBtn := tb.InlineButton{
		Unique: "auth",
		URL:    url,
		Text:   "Дать больше власти(мухаха!)",
	}

	inlineKeys := [][]tb.InlineButton{
		[]tb.InlineButton{inlineBtn},
	}

	if !m.Private() {
		return
	}

	_, err := b.Send(m.Sender, "Салют!"+
		"Я использую функции Вконтакте, поэтому мне необходимо"+
		"чуть больше возможностей."+
		"Обещаю использовать их по назначению."+
		"Переходи по ссылке ниже и жми Разрешить", &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
	if err != nil {
		log.Printf("[Error] while send message to user %+v", err)
	}
}

func (ctrl *BotCtrl) tryGetUserVkToken(b *tb.Bot, m *tb.Message) (bool, string) {
	userID := m.Sender.ID
	token, err := models.GetToken(userID)
	if err != nil {
		log.Printf("[ERROR] failed to get token for userId %d, %+v", userID, err)
		return false, ""
	}

	if token == "" {
		log.Printf("[DEBUG] User with ID %d don't have Token", userID)
		sendStartMsg(ctrl.VkAppID, ctrl.AuthURL, ctrl.VkAPIVersion, b, m)
		return false, ""
	}

	return true, token
}
