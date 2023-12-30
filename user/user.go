package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"hello/common"
	"strconv"
)

type User struct {
	ID            uint          `json:"id"`
	Tweets        pq.Int64Array `json:"tweets" gorm:"type:integer[]"`
	Username      string        `json:"username"`
	FollowedUsers pq.Int64Array `json:"followed_users" gorm:"type:integer[]"`
}

var db *gorm.DB
var err error

func Setup(router *gin.Engine, database *gorm.DB) {
	db = database

	userRoutes := router.Group("/user")
	{
		userRoutes.GET("/", getUser)
		userRoutes.PUT("/", editUser)
		userRoutes.POST("/", createUser)
		userRoutes.DELETE("/", deleteUser)
		userRoutes.PUT("/follow", followUser)
		userRoutes.PUT("/unfollow", unfollowUser)
	}
}

func getUser(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(400, gin.H{"error": "id is required in request query"})
		return
	}

	var user User

	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(200, user)
	}
}

func createUser(c *gin.Context) {
	var user User

	username, _ := c.GetPostForm("username")

	user.Username = username

	if err := db.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(200, user)
	}
}

func editUser(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(400, gin.H{"error": "id is required in request query"})
		return
	}

	var user User

	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	username, exists := c.GetPostForm("username")
	if !exists {
		c.JSON(400, gin.H{"error": "Username is required"})
		return
	}

	user.Username = username

	if err := db.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

func deleteUser(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(400, gin.H{"error": "id is required in request query"})
		return
	}

	var user User

	result := db.Where("id = ?", id).Delete(&user)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": err.Error})
	} else if result.RowsAffected > 0 {
		c.JSON(200, gin.H{"id #" + id: "deleted"})
	} else {
		c.JSON(404, gin.H{"error": "User not found"})
	}
}

func initFollowActions(c *gin.Context) (user User, followedId int64) {
	userId := c.Query("id")
	if len(userId) == 0 {
		c.JSON(400, gin.H{"error": "id is required in request query"})
		return
	}

	if err := db.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	strFollowedId, exists := c.GetPostForm("followedId")
	if !exists {
		c.JSON(400, gin.H{"error": "followedId is required"})
		return
	}

	followedId, _ = strconv.ParseInt(strFollowedId, 10, 64)

	return
}

func removeFollow(followedUsers pq.Int64Array, followedId int64) pq.Int64Array {
	var result pq.Int64Array

	for _, id := range followedUsers {
		if id != int64(followedId) {
			result = append(result, id)
		}
	}

	return result
}

func followUser(c *gin.Context) {
	user, followedId := initFollowActions(c)
	var followedUser User

	if err := db.Where("id = ?", followedId).First(&followedUser).Error; err != nil {
		c.JSON(404, gin.H{"error": "Did not find user #" + strconv.FormatInt(followedId, 10)})
		return
	}

	user.FollowedUsers = common.AddUniqueIdToArray(user.FollowedUsers, followedId)

	if err := db.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

func unfollowUser(c *gin.Context) {
	user, followedId := initFollowActions(c)

	user.FollowedUsers = removeFollow(user.FollowedUsers, followedId)

	if err := db.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}
