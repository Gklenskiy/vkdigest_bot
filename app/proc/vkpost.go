package proc

import (
	"sort"

	"github.com/Gklenskiy/vkdigest_bot/app/util"
)

// Post model from VK wall
type Post struct {
	ID       int
	Date     int64
	Reposts  int
	Likes    int
	Views    int
	OwnerID  int
	Comments int
	IsPinned int

	likesOnView    float64
	repostsOnView  float64
	commentsOnView float64
	AverageRating  float64
}

// Posts slice
type Posts []Post

// LikesOnView return count of likes on views
func (p *Post) LikesOnView() float64 {
	if p.likesOnView == 0 {
		p.likesOnView = util.RatioInPercent(float64(p.Likes), float64(p.Views))
	}

	return p.likesOnView
}

// RepostsOnView return count of reposts on views
func (p *Post) RepostsOnView() float64 {
	if p.repostsOnView == 0 {
		p.repostsOnView = util.RatioInPercent(float64(p.Reposts), float64(p.Views))
	}

	return p.repostsOnView
}

// CommentsOnView return count of comments on views
func (p *Post) CommentsOnView() float64 {
	if p.commentsOnView == 0 {
		p.commentsOnView = util.RatioInPercent(float64(p.Comments), float64(p.Views))
	}

	return p.commentsOnView
}

// Rating returns cumulative rating by some magic formula
func (p *Post) Rating() float64 {
	return p.LikesOnView() + p.RepostsOnView()*3
}

// Deviation from value
func (p *Post) Deviation(value float64) float64 {
	return p.Rating() - value
}

// DeviationFromAverage from AverageRating
func (p *Post) DeviationFromAverage() float64 {
	return p.Rating() - p.AverageRating
}

// SetAverageRating of posts
func (posts Posts) SetAverageRating() float64 {
	var sum float64
	for _, p := range posts {
		sum += p.Rating()
	}

	avg := sum / float64(len(posts))
	for i := range posts {
		posts[i].AverageRating = avg
	}

	return avg
}

// Filter posts by sort type
func (posts Posts) Filter(sortType int) {
	sort.Slice(posts, func(i, j int) bool {
		switch sortType {
		// 0 - отклонение от среднего
		case 0:
			return posts[i].DeviationFromAverage() > posts[j].DeviationFromAverage()
		// 1 - кол-во лайков
		case 1:
			return posts[i].Likes > posts[j].Likes
		// 2 - кол-во репостов
		case 2:
			return posts[i].Reposts > posts[j].Reposts
		// 3 - кол-во лайков на 1 просмотр
		case 3:
			return posts[i].LikesOnView() > posts[j].LikesOnView()
		// 4 - кол-во репостов на 1 просмотр
		case 4:
			return posts[i].RepostsOnView() > posts[j].RepostsOnView()
		// 5 - кол-во комментариев на 1 просмотр
		case 5:
			return posts[i].CommentsOnView() > posts[j].CommentsOnView()
		// 6 - кол-во комментариев
		case 6:
			return posts[i].Comments > posts[j].Comments
		// 7 - кол-во комментариев
		case 7:
			return posts[i].Rating() > posts[j].Rating()
		default:
			return posts[i].Rating() > posts[j].Rating()
		}
	})
}
