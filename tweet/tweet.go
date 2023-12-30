package tweet

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"hello/common"
	"hello/user"
	"sort"
	"time"
)

type Tweet struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

var db *gorm.DB

func Setup(router *gin.Engine, database *gorm.DB) {
	db = database

	tweetRoutes := router.Group("/tweet")
	{
		tweetRoutes.GET("/", getFollowedUsersTweets)
		tweetRoutes.POST("/", uploadTweet)
	}
}

func getFollowedUsersTweets(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(400, gin.H{"error": "id is required in request query"})
		return
	}
	
	var currentUser user.User

	if err := db.Where("id = ?", id).First(&currentUser).Error; err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	if len(currentUser.FollowedUsers) <= 0 {
		c.JSON(200, gin.H{"error": "You don't follow any users"})
		return
	}

	var users []user.User

	if err := db.Where("id = ANY(?)", currentUser.FollowedUsers).Find(&users).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var tweetIds pq.Int64Array
	for _, currUser := range users {
		tweetIds = append(tweetIds, currUser.Tweets...)
	}

	var tweets []Tweet

	if err := db.Where("id = ANY(?)", tweetIds).Find(&tweets).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	sort.Slice(tweets, func(i, j int) bool { return tweets[i].CreatedAt.After(tweets[j].CreatedAt) })

	c.JSON(200, tweets)
}

func uploadTweet(c *gin.Context) {
	uploaderId := c.Query("uploaderId")
	if len(uploaderId) == 0 {
		c.JSON(400, gin.H{"error": "uploaderId is required in request query"})
		return
	}

	var currentUser user.User
	var tweet Tweet

	content, exists := c.GetPostForm("content")
	if !exists {
		c.JSON(400, gin.H{"error": "Content is required"})
		return
	}

	tweet.Content = content

	if err := db.Create(&tweet).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := db.Where("id = ?", uploaderId).First(&currentUser).Error; err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	currentUser.Tweets = common.AddUniqueIdToArray(currentUser.Tweets, int64(tweet.ID))

	if err := db.Save(&currentUser).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tweet)
}
