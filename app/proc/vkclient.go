package proc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"

	log "github.com/go-pkgz/lgr"
)

const getPosts = "execute.getPostsMax"

// VkClient client for work with VK
type VkClient struct {
	AccessToken string
	Version     string
	BaseURL     string
	Timeout     time.Duration
}

type getPostsResponse struct {
	Response      []item                   `json:"response"`
	ExecuteErrors []map[string]interface{} `json:"execute_errors"`
}

type item struct {
	ID       []int   `json:"ids"`
	Dates    []int64 `json:"dates"`
	Reposts  []int   `json:"reposts"`
	Likes    []int   `json:"likes"`
	Views    []int   `json:"views"`
	OwnerID  []int   `json:"ownerId"`
	Comments []int   `json:"comments"`
	IsPinned []int   `json:"isPinned"`
	Last     bool    `json:"stop"`
}

// NewVkClient init vk client
func NewVkClient(vkApiURL, accessToken, version string, timeout time.Duration) *VkClient {

	client := VkClient{
		AccessToken: accessToken,
		Version:     version,
		BaseURL:     vkApiURL,
		Timeout:     timeout,
	}

	return &client
}

// GetPosts get posts from vk wall
func (client VkClient) GetPosts(offset int, deadline int64, domain string) (Posts, error) {
	if domain == "" {
		return nil, nil
	}

	httpClient := &http.Client{Timeout: client.Timeout * time.Second}

	var postsResponses []getPostsResponse
	log.Printf("[DEBUG] Start get Posts for: %s", domain)
	startMain := time.Now()
	for i := 0; ; i++ {
		request, err := client.getPostsURL(i, deadline, domain)
		if err != nil {
			log.Printf("[ERROR] failed on create Get Posts request, %s", err)
			return nil, err
		}

		log.Printf("[DEBUG] Send request: %s", request.URL.RequestURI())
		start := time.Now()
		resp, err := httpClient.Do(request)
		if err != nil {
			log.Printf("[ERROR] failed on Get Posts, %s", err)
			return nil, err
		}
		elapsed := time.Since(start)
		log.Printf("[DEBUG] Get Request took %s", elapsed)

		log.Printf("[DEBUG] Parse response: %s", request.URL.RequestURI())
		start = time.Now()
		parsedResponse, err := parse(resp)
		if err != nil {
			log.Printf("[ERROR] failed on parse response body, %s", err)
			return nil, err
		}
		elapsed = time.Since(start)
		log.Printf("[DEBUG] Parse response took %s", elapsed)

		postsResponses = append(postsResponses, parsedResponse)
		if isLast(parsedResponse) {
			break
		}
	}
	elapsed := time.Since(startMain)
	log.Printf("[DEBUG] End get Posts for %s, took %s", domain, elapsed)

	posts := mapPosts(postsResponses)
	posts = cropByDeadline(posts, deadline)

	log.Printf("[DEBUG] Posts count from: %s is %d", domain, len(posts))
	return posts, nil
}

func isLast(response getPostsResponse) bool {
	len := len(response.Response)
	if len == 0 {
		return true
	}

	lastIdx := sort.Search(len, func(i int) bool {
		return response.Response[i].Last == true
	})

	if lastIdx < len {
		return true
	}

	return false
}

func cropByDeadline(posts []Post, deadline int64) []Post {
	len := len(posts)
	if len == 0 {
		return nil
	}

	deadlineIdx := sort.Search(len, func(i int) bool {
		return posts[i].Date < deadline && posts[i].IsPinned != 1
	})

	if deadlineIdx < len {
		return posts[:deadlineIdx]
	}

	return posts
}

func parse(httpResponse *http.Response) (getPostsResponse, error) {
	var result getPostsResponse
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		log.Printf("[ERROR] failed on read response body, %s", err)
		return result, err
	}

	log.Printf("[DEBUG] Unmurshal: %s", body)
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("[ERROR] failed to unmarshal response, %s", err)
		return result, err
	}

	return result, nil
}

func mapPosts(postsResponses []getPostsResponse) []Post {
	var posts []Post
	for _, postsResponse := range postsResponses {
		mappedPosts := mapPost(postsResponse)
		posts = append(posts, mappedPosts...)
	}

	return posts
}

func mapPost(postsResponse getPostsResponse) []Post {
	var posts []Post
	for i := range postsResponse.Response {
		for idx, elem := range postsResponse.Response[i].ID {
			post := Post{
				ID:       elem,
				Date:     postsResponse.Response[i].Dates[idx],
				Reposts:  postsResponse.Response[i].Reposts[idx],
				Likes:    postsResponse.Response[i].Likes[idx],
				Views:    postsResponse.Response[i].Views[idx],
				OwnerID:  postsResponse.Response[i].OwnerID[idx],
				Comments: postsResponse.Response[i].Comments[idx],
				IsPinned: postsResponse.Response[i].IsPinned[idx],
			}

			posts = append(posts, post)
		}
	}

	return posts
}

func (client VkClient) getPostsURL(offset int, deadline int64, domain string) (*http.Request, error) {
	req, err := http.NewRequest("GET", client.BaseURL+"/"+getPosts, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("offset", strconv.Itoa(offset))
	q.Add("deadline", strconv.FormatInt(deadline, 10))
	q.Add("domain", domain)
	q.Add("v", client.Version)
	q.Add("access_token", client.AccessToken)

	req.URL.RawQuery = q.Encode()

	return req, err
}
