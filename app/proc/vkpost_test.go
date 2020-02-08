package proc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVkPost_CalcLikeOnView(t *testing.T) {
	tbl := []struct {
		in   Post
		out  float64
		name string
	}{
		{Post{
			ID:       1,
			Date:     time.Now().Unix(),
			Reposts:  10,
			Likes:    10,
			Views:    100,
			OwnerID:  123,
			Comments: 3,
			IsPinned: 0,
		}, 10, "simple"},
		{Post{
			ID:       1,
			Date:     time.Now().Unix(),
			Reposts:  10,
			Likes:    10,
			Views:    0,
			OwnerID:  123,
			Comments: 3,
			IsPinned: 0,
		}, 0, "zero views return 0"},
		{Post{
			ID:       1,
			Date:     time.Now().Unix(),
			Reposts:  10,
			Likes:    0,
			Views:    100,
			OwnerID:  123,
			Comments: 3,
			IsPinned: 0,
		}, 0, "zero likes return 0"},
	}

	for _, tt := range tbl {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, tt.in.LikesOnView())
		})
	}
}

func TestVkPost_CalcRepostOnView(t *testing.T) {
	tbl := []struct {
		in   Post
		out  float64
		name string
	}{
		{Post{
			ID:       1,
			Date:     time.Now().Unix(),
			Reposts:  10,
			Likes:    10,
			Views:    100,
			OwnerID:  123,
			Comments: 3,
			IsPinned: 0,
		}, 10, "simple"},
		{Post{
			ID:       1,
			Date:     time.Now().Unix(),
			Reposts:  10,
			Likes:    10,
			Views:    0,
			OwnerID:  123,
			Comments: 3,
			IsPinned: 0,
		}, 0, "zero views return 0"},
		{Post{
			ID:       1,
			Date:     time.Now().Unix(),
			Reposts:  0,
			Likes:    10,
			Views:    100,
			OwnerID:  123,
			Comments: 3,
			IsPinned: 0,
		}, 0, "zero reposts return 0"},
	}

	for _, tt := range tbl {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, tt.in.RepostsOnView())
		})
	}
}
func TestVkPost_Rating(t *testing.T) {
	tbl := []struct {
		in   Post
		out  float64
		name string
	}{
		{Post{
			ID:       1,
			Date:     time.Now().Unix(),
			Reposts:  5,
			Likes:    10,
			Views:    100,
			OwnerID:  123,
			Comments: 3,
			IsPinned: 0,
		}, 25, "simple"},
	}

	for _, tt := range tbl {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, tt.in.Rating())
		})
	}
}

func TestVkPost_SetAverageRating(t *testing.T) {
	tbl := []struct {
		in   Posts
		out  float64
		name string
	}{
		{
			Posts{
				Post{
					ID:       1,
					Date:     time.Now().Unix(),
					Reposts:  10,
					Likes:    10,
					Views:    100,
					OwnerID:  123,
					Comments: 3,
					IsPinned: 0,
				},
				Post{
					ID:       1,
					Date:     time.Now().Unix(),
					Reposts:  10,
					Likes:    0,
					Views:    100,
					OwnerID:  123,
					Comments: 3,
					IsPinned: 0,
				},
			},
			35, "should return AverageRating"},
	}

	for _, tt := range tbl {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, tt.in.SetAverageRating())
		})
	}
}

func TestVkPost_DeviationFromAverage(t *testing.T) {
	tbl := []struct {
		in   Posts
		out  float64
		name string
	}{
		{
			Posts{
				Post{
					ID:       1,
					Date:     time.Now().Unix(),
					Reposts:  10,
					Likes:    10,
					Views:    100,
					OwnerID:  123,
					Comments: 3,
					IsPinned: 0,
				},
				Post{
					ID:       1,
					Date:     time.Now().Unix(),
					Reposts:  10,
					Likes:    0,
					Views:    100,
					OwnerID:  123,
					Comments: 3,
					IsPinned: 0,
				},
			},
			5, "should return DeviationFromAverage"},
	}

	for _, tt := range tbl {
		t.Run(tt.name, func(t *testing.T) {
			tt.in.SetAverageRating()
			assert.Equal(t, tt.out, tt.in[0].DeviationFromAverage())
		})
	}
}

func TestVkPost_Filter(t *testing.T) {
	tbl := []struct {
		in   Posts
		out  Posts
		name string
	}{
		{
			Posts{
				Post{
					ID:       1,
					Date:     time.Now().Unix(),
					Reposts:  10,
					Likes:    0,
					Views:    100,
					OwnerID:  123,
					Comments: 3,
					IsPinned: 0,
				},
				Post{
					ID:       2,
					Date:     time.Now().Unix(),
					Reposts:  10,
					Likes:    10,
					Views:    100,
					OwnerID:  123,
					Comments: 3,
					IsPinned: 0,
				},
			},
			Posts{
				Post{
					ID:            2,
					Date:          time.Now().Unix(),
					Reposts:       10,
					Likes:         10,
					Views:         100,
					OwnerID:       123,
					Comments:      3,
					IsPinned:      0,
					likesOnView:   10,
					repostsOnView: 10,
					AverageRating: 35,
				},
				Post{
					ID:            1,
					Date:          time.Now().Unix(),
					Reposts:       10,
					Likes:         0,
					Views:         100,
					OwnerID:       123,
					Comments:      3,
					IsPinned:      0,
					repostsOnView: 10,
					AverageRating: 35,
				},
			}, "should return right order"},
	}

	for _, tt := range tbl {
		t.Run(tt.name, func(t *testing.T) {
			tt.in.SetAverageRating()
			tt.in.Filter(0)
			assert.Equal(t, tt.out, tt.in)
		})
	}
}
