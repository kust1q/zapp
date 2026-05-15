package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/kust1q/Zapp/search/internal/domain/entity"
	"github.com/sirupsen/logrus"
)

const (
	IndexTweets = "tweets"
	IndexUsers  = "users"
)

type elasticRepository struct {
	client *elasticsearch.Client
}

func NewElasticRepository(client *elasticsearch.Client) *elasticRepository {
	return &elasticRepository{client: client}
}

func (r *elasticRepository) IndexTweet(ctx context.Context, tweet *entity.Tweet) error {
	doc := tweetDoc{
		Content:  tweet.Content,
		Username: tweet.Author.Username,
		UserID:   tweet.Author.ID,
	}
	return r.indexDocument(ctx, IndexTweets, tweet.ID, doc)
}

func (r *elasticRepository) IndexUser(ctx context.Context, user *entity.User) error {
	doc := userDoc{
		Username: user.Username,
		Bio:      user.Bio,
	}
	return r.indexDocument(ctx, IndexUsers, user.ID, doc)
}

func (r *elasticRepository) DeleteTweet(ctx context.Context, tweetID int) error {
	req := esapi.DeleteRequest{
		Index:      IndexTweets,
		DocumentID: strconv.Itoa(tweetID),
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("elastic request error: %w", err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode == 404 {
		logrus.WithField("tweet_id", tweetID).Info("tweet not found in elastic during deletion")
		return nil
	}

	if res.IsError() {
		return fmt.Errorf("elastic delete error: %s", res.String())
	}

	return nil
}

func (r *elasticRepository) DeleteUser(ctx context.Context, userID int) error {
	req := esapi.DeleteRequest{
		Index:      IndexUsers,
		DocumentID: strconv.Itoa(userID),
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("elastic request error: %w", err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode == 404 {
		logrus.WithField("user_id", userID).Info("user not found in elastic during deletion")
		return nil
	}

	if res.IsError() {
		return fmt.Errorf("elastic delete error: %s", res.String())
	}

	return nil
}

func (r *elasticRepository) DeleteTweetsByUserID(ctx context.Context, userID int) error {
	query := map[string]any{
		"query": map[string]any{
			"term": map[string]any{
				"user_id": userID,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return err
	}

	req := esapi.DeleteByQueryRequest{
		Index: []string{IndexTweets},
		Body:  &buf,
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("delete by query req error: %w", err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.IsError() {
		return fmt.Errorf("delete by query error: %s", res.String())
	}

	return nil
}

func (r *elasticRepository) indexDocument(ctx context.Context, index string, id int, body any) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: strconv.Itoa(id),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("elastic request error: %w", err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.IsError() {
		return fmt.Errorf("elastic error indexing %s: %s", index, res.String())
	}
	return nil
}

func (r *elasticRepository) SearchTweets(ctx context.Context, query string) ([]int, error) {
	queryMap := map[string]any{
		"query": map[string]any{
			"multi_match": map[string]any{
				"query":     query,
				"fields":    []string{"content", "username"},
				"fuzziness": "AUTO",
			},
		},
		"_source": false,
	}
	return r.performSearch(ctx, IndexTweets, queryMap)
}

func (r *elasticRepository) SearchUsers(ctx context.Context, query string) ([]int, error) {
	queryMap := map[string]any{
		"query": map[string]any{
			"multi_match": map[string]any{
				"query":     query,
				"fields":    []string{"username^3", "bio"},
				"fuzziness": "AUTO",
			},
		},
		"_source": false,
	}
	return r.performSearch(ctx, IndexUsers, queryMap)
}

func (r *elasticRepository) performSearch(ctx context.Context, index string, queryMap map[string]any) ([]int, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(queryMap); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(index),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hitsObj, ok := result["hits"].(map[string]any)
	if !ok {
		return []int{}, nil
	}

	hitsList, ok := hitsObj["hits"].([]any)
	if !ok {
		return []int{}, nil
	}

	var ids []int
	for _, hit := range hitsList {
		hitMap := hit.(map[string]any)
		idStr := hitMap["_id"].(string)

		if id, err := strconv.Atoi(idStr); err == nil {
			ids = append(ids, id)
		}
	}

	return ids, nil
}
