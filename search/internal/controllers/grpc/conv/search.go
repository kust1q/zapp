package conv

import searchproto "github.com/kust1q/Zapp/search/pkg/gen/proto/search"

func ToSearchUserProtoResponse(ids []int) *searchproto.SearchUsersResponse {
	res := make([]int64, 0, len(ids))
	for _, n := range ids {
		res = append(res, int64(n))
	}
	return &searchproto.SearchUsersResponse{
		UserIds: res,
	}
}

func ToSearchTweetProtoResponse(ids []int) *searchproto.SearchTweetsResponse {
	res := make([]int64, 0, len(ids))
	for _, n := range ids {
		res = append(res, int64(n))
	}
	return &searchproto.SearchTweetsResponse{
		TweetIds: res,
	}
}
