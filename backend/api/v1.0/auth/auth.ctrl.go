package auth

import (
	"github.com/Ismail14098/agyn_test_rest/database/models"
	"github.com/Ismail14098/agyn_test_rest/lib/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"time"
)

type User = models.User

func Hash(password string) (string, error){
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func checkHash(password string, hash string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err != nil
}

func generateToken(data common.JSON) (string, error) {

	//  token is valid for 7days
	date := time.Now().Add(time.Second * 60 * 50)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": data,
		"exp": date,
	})

	pwd, _ := os.Getwd()
	keyPath := pwd+"/jwtsecret.key"

	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(key) //getting signed token
	return tokenString, err
}

func register(ctx *gin.Context){
	db := ctx.MustGet("db").(*gorm.DB)

	type RequestBody struct {
		Username string `json:"username" binding:"required,correctLogin"`
		Email string `json:"email" binding:"required"`
		Firstname string `json:"firstname" binding:"required,correctName"`
		Lastname string `json:"lastname" binding:"required,correctName"`
		Password string `json:"password" binding:"required,correctPassword"`
		RoleId string `json:"roleId" binding:"required"`
	}

	var body RequestBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{
			"error": "check input data",
		})
		return
	}

	var exist User
	if err := db.Where("username = ?", body.Username).First(&exist).Error; err == nil {
		ctx.JSON(409, gin.H{
			"error": "use another username",
		})
		return
	}

	if err := db.Where("email = ?", body.Email).First(&exist).Error; err == nil {
		ctx.JSON(409, gin.H{
			"error": "use another email",
		})
		return
	}

	passwordhash, err := Hash(body.Password)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	user := User{
		Username: body.Username,
		Email: body.Email,
		Firstname: body.Firstname,
		Lastname: body.Lastname,
		PasswordHash: passwordhash,
	}

	var role models.Role

	result := db.First(&role, body.RoleId)
	if result.Error != nil {
		ctx.AbortWithStatus(409)
		return
	}

	db.Create(&user)

	userRole := models.UserRole{
		UserId: user.ID,
		RoleId: role.ID,
	}
	db.Create(&userRole)

	ctx.JSON(200, common.JSON{
		"user":  user.Serialize(),
	})

}

func login(ctx *gin.Context){
	db := ctx.MustGet("db").(*gorm.DB)

	type RequestBody struct {
		Username string `json:"username" binding:"required,correctLogin"`
		Password string `json:"password" binding:"required,correctPassword"`
	}

	var body RequestBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{
			"error": "check input data",
		})
		return
	}

	var user User
	found := false
	if err := db.Where("username = ?", body.Username).First(&user).Error; err == nil {
		found = true
	}

	if !found {
		if err := db.Where("email = ?", body.Username).First(&user).Error; err == nil{
			found = true
		}
	}

	if !found {
		ctx.JSON(400, gin.H{
			"error": "check input data",
		})
		return
	}

	if checkHash(body.Password, user.PasswordHash) {
		ctx.AbortWithStatus(401)
		return
	}

	serialized := user.Serialize()
	token, _ := generateToken(serialized)

	var notifications []models.Notification
	db.Where("user_id = ?", user.ID).Find(&notifications)

	ctx.SetCookie("token", token, 60*5*10, "/", "", false, true)
	ctx.JSON(200, common.JSON{
		"user":  user.Serialize(),
		"token": token,
		"notifications": notifications,
	})
	db.Delete(&notifications)
}

func check(ctx *gin.Context){
	userRaw, ok := ctx.Get("user")
	if !ok {
		ctx.AbortWithStatus(401)
		return
	}

	user := userRaw.(User)
	tokenExpire := ctx.MustGet("token_expire").(string)
	tokenTime, err := time.Parse(time.RFC3339, tokenExpire)
	if err != nil {
		ctx.AbortWithStatus(401)
		return
	}
	now := time.Now().Unix()
	diff := tokenTime.Unix() - now

	if diff < 5 {
		token, _ := generateToken(user.Serialize())
		ctx.SetCookie("token", token, 60*50, "/","", false, true)
		ctx.JSON(200, common.JSON{
			"user":  user.Serialize(),
			"token": token,
		})
		return
	}

	ctx.JSON(200, common.JSON{
		"user":  user.Serialize(),
		"token": nil,
	})
}

func logout(ctx *gin.Context){
	tokenString, err := ctx.Cookie("token")
	if err != nil {
		ctx.AbortWithStatus(401)
		return
	}
	ctx.SetCookie("token", tokenString, -1, "/", "", false, true)
	ctx.JSON(200, common.JSON{
		"status": "success",
	})
}