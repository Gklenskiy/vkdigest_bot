package proc

import (
	"fmt"
	"time"

	"github.com/Gklenskiy/vkdigest_bot/app/models"
	log "github.com/go-pkgz/lgr"
	tb "gopkg.in/tucnak/telebot.v2"
)

// BotCtrl
type BotCtrl struct {
	BotCtrlSettings
}

type BotCtrlSettings struct {
	VkBaseURL    string
	VkApiVersion string
	Domains      []VkDomain
	VkAppId      string
	AuthUrl      string
}

// NewBotCtrl
func NewBotCtrl(settings BotCtrlSettings) *BotCtrl {
	return &BotCtrl{
		BotCtrlSettings: settings,
	}
}

// PingCtrl
func (ctrl *BotCtrl) PingCtrl(b *tb.Bot, m *tb.Message) {
	b.Send(m.Sender, "pong 4")
}

// StartCtrl handler /start command
func (ctrl *BotCtrl) StartCtrl(b *tb.Bot, m *tb.Message) {
	SendStartMsg(ctrl.VkAppId, ctrl.AuthUrl, ctrl.VkApiVersion, b, m)
}

// TrendsCtrl for handle trend command
func (ctrl *BotCtrl) TrendsCtrl(b *tb.Bot, m *tb.Message) {
	sortType := 0
	daysDeadline := 7
	deadline := time.Now().AddDate(0, 0, -daysDeadline).Unix()
	log.Printf("Deadline date %v \n", time.Unix(deadline, 0))

	userID := m.Sender.ID
	token, err := models.GetToken(userID)
	if err != nil {
		log.Printf("[ERROR] failed to get token for userId %d, %+v", userID, err)
	}

	if token == "" {
		log.Printf("[DEBUG] User with ID %d don't have Token", userID)
		SendStartMsg(ctrl.VkAppId, ctrl.AuthUrl, ctrl.VkApiVersion, b, m)
		return

	}

	log.Printf("[DEBUG] Start process Trends Request")
	startProcess := time.Now()
	vkClient := NewVkClient(ctrl.VkBaseURL, token, ctrl.VkApiVersion, 500*time.Millisecond)

	var allPosts Posts
	log.Printf("[DEBUG] Start get Posts")
	start := time.Now()
	ch := make(chan Posts)
	for _, domain := range ctrl.Domains {
		go ProcessDomain(vkClient, deadline, domain, ch)
	}
	for range ctrl.Domains {
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
		b.Send(m.Sender, v)
	}
}

func ProcessDomain(vkClient *VkClient, deadline int64, domain VkDomain, ch chan<- Posts) {
	posts, err := vkClient.GetPosts(0, deadline, domain.Name)
	if err != nil {
		log.Printf("[ERROR] failed to get posts from Vk, %+v", err)
	}

	posts.SetAverageRating()
	ch <- posts
}

func SendStartMsg(vkApi string, authURL string, vkApiVersion string, b *tb.Bot, m *tb.Message) {
	url := fmt.Sprintf(`https://oauth.vk.com/authorize?client_id=%s&display=page&redirect_uri=%s&scope=wall,offline&response_type=code&v=%s&state=%d`,
		vkApi, authURL, vkApiVersion, m.Sender.ID)

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

	b.Send(m.Sender, "Салют!"+
		"Я использую функции Вконтакте, поэтому мне необходимо"+
		"чуть больше возможностей."+
		"Обещаю использовать их по назначению."+
		"Переходи по ссылке ниже и жми Разрешить", &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}
