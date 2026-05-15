package conv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/providers/db/conv"
	"github.com/kust1q/Zapp/main/internal/providers/db/models"
)

func TestDBConv_User(t *testing.T) {
	u := &entity.User{ID: 1, Credential: &entity.Credential{Email: "e"}}
	m := conv.FromDomainToUserModel(u)
	assert.Equal(t, 1, m.ID)
	assert.Equal(t, "e", m.Email)
	assert.Nil(t, conv.FromDomainToUserModel(nil))

	dm := conv.FromUserModelToDomain(&models.User{ID: 2, Email: "e2"})
	assert.Equal(t, 2, dm.ID)
	assert.Equal(t, "e2", dm.Credential.Email)
	assert.Nil(t, conv.FromUserModelToDomain(nil))

	f := &entity.Follow{FollowerID: 1, FollowingID: 2}
	fm := conv.FromDomainToFollowModel(f)
	assert.Equal(t, 1, fm.FollowerID)
	assert.Nil(t, conv.FromDomainToFollowModel(nil))

	df := conv.FromFollowModelToDomain(&models.Follow{FollowerID: 3})
	assert.Equal(t, 3, df.FollowerID)
	assert.Nil(t, conv.FromFollowModelToDomain(nil))
}

func TestDBConv_Tweet(t *testing.T) {
	tw := &entity.Tweet{ID: 1, Author: &entity.SmallUser{ID: 1}}
	tm := conv.FromDomainToTweetModel(tw)
	assert.Equal(t, 1, tm.ID)
	assert.Nil(t, conv.FromDomainToTweetModel(nil))

	dtm := conv.FromTweetModelToDomain(&models.Tweet{ID: 2, UserID: 2})
	assert.Equal(t, 2, dtm.ID)
	assert.Nil(t, conv.FromTweetModelToDomain(nil))
	
	list := conv.FromTweetModelToDomainList([]models.Tweet{{ID: 1}, {ID: 2}})
	assert.Len(t, list, 2)
}

func TestDBConv_RetweetLike(t *testing.T) {
	rt := &entity.Retweet{ID: 1}
	rm := conv.FromDomainToRetweetModel(rt)
	assert.Equal(t, 1, rm.ID)
	assert.Nil(t, conv.FromDomainToRetweetModel(nil))

	drm := conv.FromRetweetModelToDomain(&models.Retweet{ID: 2})
	assert.Equal(t, 2, drm.ID)
	assert.Nil(t, conv.FromRetweetModelToDomain(nil))

	lk := &entity.Like{UserID: 1, TweetID: 10}
	lm := conv.FromDomainToLikeModel(lk)
	assert.Equal(t, 1, lm.UserID)
	assert.Nil(t, conv.FromDomainToLikeModel(nil))

	dlm := conv.FromLikeModelToDomain(&models.Like{UserID: 2})
	assert.Equal(t, 2, dlm.UserID)
	assert.Nil(t, conv.FromLikeModelToDomain(nil))

	llist := conv.FromLikeModelToDomainList([]models.Like{{UserID: 1}, {UserID: 2}})
	assert.Len(t, llist, 2)

	c := &models.Counters{LikeCount: 5}
	dc := conv.FromCountersModelToDomain(c)
	assert.Equal(t, 5, dc.LikeCount)
}

func TestDBConv_Media(t *testing.T) {
	me := &entity.TweetMedia{ID: 1}
	mm := conv.FromDomainToTweetMediaModel(me)
	assert.Equal(t, 1, mm.ID)
	assert.Nil(t, conv.FromDomainToTweetMediaModel(nil))

	dme := conv.FromTweetMediaModelToDomain(&models.TweetMedia{ID: 2})
	assert.Equal(t, 2, dme.ID)
	assert.Nil(t, conv.FromTweetMediaModelToDomain(nil))

	avm := conv.FromDomainToAvatarModel(&entity.Avatar{ID: 1})
	assert.Equal(t, 1, avm.ID)
	assert.Nil(t, conv.FromDomainToAvatarModel(nil))
	
	av := &models.Avatar{ID: 3}
	dav := conv.FromAvatarModelToDomain(av)
	assert.Equal(t, 3, dav.ID)
	assert.Nil(t, conv.FromAvatarModelToDomain(nil))
}
