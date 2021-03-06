package cmd

import (
	"fmt"
	"io/ioutil"
	"time"

	log "github.com/go-pkgz/lgr"
	"gopkg.in/yaml.v2"

	"github.com/Gklenskiy/vkdigest_bot/app/proc"
)

// ConsoleCommand with params
type ConsoleCommand struct {
	//Domains      []string `long:"domains" default:"newalbums" description:"port for listen" env-delim:","`	DaysDeadline int    `long:"deadline" default:"7" description:"service url"`
	DaysDeadline int    `long:"deadline" default:"7" description:"service url"`
	SortType     int    `long:"sort_type" default:"1" description:"token for telegram bot"`
	VkToken      string `long:"vk_token" env:"VK_TOKEN" required:"true" description:"Vk Token"`
	Conf         string `short:"c" long:"conf" env:"MR_CONF" default:"conf.yml" description:"config file (yml)"`

	CommonOpts
}

// Execute is the entry point for "console" command, called by flag parser
func (cmd *ConsoleCommand) Execute(args []string) error {
	conf, err := loadConfig(cmd.Conf)
	if err != nil {
		log.Fatalf("[ERROR] can't load config %s, %v", cmd.Conf, err)
	}
	vkClient := proc.NewVkClient(conf.Sources["vk"].BaseURL, cmd.VkToken, conf.Sources["vk"].APIVersion, time.Millisecond)

	deadline := time.Now().AddDate(0, 0, -cmd.DaysDeadline).Unix()
	log.Printf("Deadline date %v \n", time.Unix(deadline, 0))

	var allPosts proc.Posts
	for _, domain := range conf.Sources["vk"].Domains {
		posts, err := vkClient.GetPosts(0, deadline, domain.Name)
		if err != nil {
			log.Printf("[ERROR] failed to get posts from Vk, %+v", err)
			return err
		}

		posts.SetAverageRating()
		allPosts = append(allPosts, posts...)
	}

	allPosts.Filter(cmd.SortType)

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

	for _, v := range res {
		log.Printf(v)
	}

	return nil
}

func loadConfig(fname string) (res *proc.Conf, err error) {
	log.Printf(fname)
	res = &proc.Conf{}
	data, err := ioutil.ReadFile(fname) // nolint
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, res); err != nil {
		return nil, err
	}

	return res, nil
}
