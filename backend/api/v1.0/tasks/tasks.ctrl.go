package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Ismail14098/agyn_test_rest/database/models"
	"github.com/Ismail14098/agyn_test_rest/lib/status"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func findOwnTaskFromRedisById(ctx *gin.Context,id uint) (models.Task, bool)  {
	user := ctx.MustGet("user").(models.User)
	rdb := ctx.MustGet("rdb").(*redis.Client)
	defCtx := ctx.MustGet("defCtx").(*context.Context)
	key := "tasks_user_"+strconv.Itoa(int(user.ID))

	var task models.Task
	found := false

	value := rdb.LRange(*defCtx, key, 0, -1).Val()
	for _, val := range value{
		var taskRedis models.Task
		json.Unmarshal([]byte(val), &taskRedis)
		if taskRedis.ID == id {
			task = taskRedis
			found = true
			break
		}
	}
	return task,found
}

func deleteOwnTaskFromRedisById(ctx *gin.Context,id uint) bool  {
	user := ctx.MustGet("user").(models.User)
	rdb := ctx.MustGet("rdb").(*redis.Client)
	defCtx := ctx.MustGet("defCtx").(*context.Context)
	key := "tasks_user_"+strconv.Itoa(int(user.ID))

	found := false

	value := rdb.LRange(*defCtx, key, 0, -1).Val()
	var deleteValue []byte
	for _, val := range value{
		var taskRedis models.Task
		json.Unmarshal([]byte(val), &taskRedis)
		if taskRedis.ID == id {
			deleteValue, _ = json.Marshal(&taskRedis)
			found = true
			break
		}
	}

	rdb.LRem(*defCtx,key,0, string(deleteValue))
	return found
}

func deleteTaskFromRedisById(ctx *gin.Context, id uint) (bool, uint){
	rdb := ctx.MustGet("rdb").(*redis.Client)
	defCtx := ctx.MustGet("defCtx").(*context.Context)

	found := false
	var authorId uint

	var keys []string
	keys, _ = rdb.Keys(*defCtx, "tasks_user_*").Result()
	for _, i := range keys{
		value := rdb.LRange(*defCtx, i, 0, -1).Val()
		var deleteValue []byte
		for _, val := range value{
			var taskRedis models.Task
			json.Unmarshal([]byte(val), &taskRedis)
			if taskRedis.ID == id {
				authorId = taskRedis.AuthorId
				deleteValue, _ = json.Marshal(&taskRedis)
				found = true
				break
			}
		}
		rdb.LRem(*defCtx,i,0, string(deleteValue))
	}
	return found, authorId
}


func create(ctx *gin.Context){
	db := ctx.MustGet("db").(*gorm.DB)
	type RequestBody struct {
		Title string `json:"title" binding:"required,correctTitle"`
		Text string `json:"text" binding:"required,correctText"`
		CategoryId uint `json:"category" binding:"required"`
		ExpTime string`json:"expTime" binding:"required,correctDate"`
	}

	var body RequestBody
	if err := ctx.BindJSON(&body); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(400)
		return
	}

	expTimeParsed, _ := time.Parse(time.RFC3339, body.ExpTime)

	user := ctx.MustGet("user").(models.User)

	var category models.Category
	result := db.Find(&category, int(body.CategoryId))
	if result.Error != nil {
		ctx.AbortWithStatus(400)
		return
	}

	task := models.Task{
		Title:      body.Title,
		Text:       body.Text,
		CategoryId: body.CategoryId,
		Category: category,
		AuthorId:   user.ID,
		Author: user,
		Status:     status.Statuses["InProgress"],
		ExpTime:    expTimeParsed,
	}

	rdb := ctx.MustGet("rdb").(*redis.Client)
	defCtx := ctx.MustGet("defCtx").(*context.Context)
	key := "tasks_user_"+strconv.Itoa(int(user.ID))

	db.Create(&task)
	serializedValue, _ := json.Marshal(task)
	taskStr := string(serializedValue)

	rdb.RPush(*defCtx, key, taskStr)
	val, err := rdb.LLen(*defCtx, key).Result()

	if val > 10 && err == nil{
		rdb.LPop(*defCtx, key).Val()
	}

	ctx.JSON(200, gin.H{
		"task": task,
	})
}

func view(ctx *gin.Context){
	db := ctx.MustGet("db").(*gorm.DB)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatus(400)
		return
	}

	var task models.Task
	task, found := findOwnTaskFromRedisById(ctx, uint(id))

	if !found {
		db.Preload("Category").Preload("Author").First(&task, id)
	}

	ctx.JSON(200, gin.H{
		"task": task,
	})
}

func all(ctx *gin.Context){
	db := ctx.MustGet("db").(*gorm.DB)
	user:= ctx.MustGet("user").(models.User)

	var tasks []models.Task
	page := ctx.DefaultQuery("page", "1")
	pageInt, err := strconv.Atoi(page)

	if err != nil || pageInt < 1 {
		ctx.AbortWithStatus(400)
		return
	}

	rdb := ctx.MustGet("rdb").(*redis.Client)
	defCtx := ctx.MustGet("defCtx").(*context.Context)
	key := "tasks_user_" + strconv.Itoa(int(user.ID))

	len := rdb.LLen(*defCtx, key).Val()
	// Fill redis by last 10 tasks
	var counter int
	var filledRedis bool
	if len == 0 {
		db.Limit(10).Last(&tasks)
		for task := range tasks {
			serializedValue, _ := json.Marshal(task)
			rdb.Expire(*defCtx, key, time.Hour)
			rdb.RPush(*defCtx, key, string(serializedValue))
		}
		filledRedis = true
	}
	//If first page and we filled redis we can return
	if pageInt == 1 && filledRedis {
		ctx.JSON(200, gin.H{
			"tasks": tasks,
		})
		return
	}

	if pageInt == 1{
		//get last tasks from redis
		value := rdb.LRange(*defCtx, key, 0, -1).Val()
		for _, val := range value {
			counter++
			var task models.Task
			json.Unmarshal([]byte(val), &task)
			tasks = append(tasks, task)
		}
	}

	if pageInt > 1 && counter == 10 {
		offset := 10 * (pageInt-1)
		db.Where("author_id = ", user.ID).Offset(offset).Limit(10).Preload("Category").Preload("Author").Find(&tasks)
	}

	ctx.JSON(200, gin.H{
		"tasks": tasks,
	})
}

func show(ctx *gin.Context){
	db := ctx.MustGet("db").(*gorm.DB)

	var tasks []models.Task
	page := ctx.DefaultQuery("page", "1")
	pageInt, err := strconv.Atoi(page)

	if err != nil || pageInt < 1 {
		ctx.JSON(400, gin.H{
			"error": "check id",
		})
		return
	}

	offset := 10 * (pageInt-1)
	db.Offset(offset).Limit(10).Preload("Category").Preload("Author").Find(&tasks)

	if len(tasks) == 0 {
		ctx.AbortWithStatus(400)
		return
	}

	ctx.JSON(200, gin.H{
		"tasks": tasks,
	})
}

func edit(ctx *gin.Context){
	db := ctx.MustGet("db").(*gorm.DB)
	type RequestBody struct {
		Id uint `json:"id" binding:"required"`
		Title string `json:"title" binding:"required,correctTitle"`
		Text string `json:"text" binding:"required,correctText"`
		Status string `json:"status" binding:"-"`
		CategoryId uint `json:"category" binding:"required"`
		ExpTime string`json:"expTime" binding:"-"`
	}

	var body RequestBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{
			"error": "check input data",
		})
		return
	}

	var task models.Task
	task, found := findOwnTaskFromRedisById(ctx, body.Id)
	if !found {
		result := db.Find(&task, body.Id)
		if result.Error != nil {
			ctx.AbortWithStatus(400)
			return
		}
	}

	user := ctx.MustGet("user").(models.User)
	if task.AuthorId != user.ID {
		ctx.AbortWithStatus(401)
		return
	}


	if body.ExpTime != "" {
		expTimeParsed, err := time.Parse(time.RFC3339, body.ExpTime)
		if err == nil && time.Now().Before(expTimeParsed) {
			task.ExpTime = expTimeParsed
		}
	}


	var category models.Category
	result := db.Find(&category, int(body.CategoryId))
	if result.Error != nil {
		ctx.AbortWithStatus(400)
		return
	}


	task.Title = body.Title
	task.CategoryId = body.CategoryId
	task.Text = body.Text

	_, ok := status.Statuses[body.Status]
	if body.Status != "" && ok {
		task.Status = status.Statuses[body.Status]
	} else {
		task.Status = status.Statuses["InProgress"]
	}

	if found {
		deleteOwnTaskFromRedisById(ctx,task.ID)

		rdb := ctx.MustGet("rdb").(*redis.Client)
		defCtx := ctx.MustGet("defCtx").(*context.Context)
		key := "tasks_user_"+strconv.Itoa(int(user.ID))

		serializedValue, _ := json.Marshal(task)
		taskStr := string(serializedValue)

		rdb.RPush(*defCtx,key,taskStr)
	} else {
		db.Save(&task)
	}

	ctx.JSON(200, gin.H{
		"task": task,
	})
}

func editSudo(ctx *gin.Context){
	db := ctx.MustGet("db").(*gorm.DB)
	type RequestBody struct {
		Id string `json:"id" binding:"required"`
		Title string `json:"title" binding:"required,correctTitle"`
		Text string `json:"text" binding:"required,correctText"`
		Status string `json:"status" binding:"-"`
		CategoryId string `json:"category" binding:"required"`
		ExpTime string`json:"expTime" binding:"-"`
	}

	var body RequestBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{
			"error": "check input data",
		})
		return
	}

	id, errId := strconv.Atoi(body.Id)
	categoryId, errCatId := strconv.Atoi(body.CategoryId)
	if errId != nil || errCatId != nil {
		ctx.AbortWithStatus(400)
		return
	}

	var task models.Task
	db.Find(&task,id)
	task.Title = body.Title
	task.Text = body.Text
	task.CategoryId = uint(categoryId)

	statusStr, ok := status.Statuses[body.Status]
	if body.Status != "" && ok {
		task.Status = statusStr
	} else {
		task.Status = status.Statuses["InProgress"]
	}

	//testTimeParsed, err := time.Parse(time.RFC3339,"2020-11-30T13:42:00Z")
	//fmt.Println(testTimeParsed, err)
	if body.ExpTime != "" {
		expTimeParsed, err := time.Parse(time.RFC3339, body.ExpTime)
		if err == nil && time.Now().Before(expTimeParsed) {
			task.ExpTime = expTimeParsed
		}
	}

	db.Save(&task)
	ctx.JSON(200, gin.H{
		"task": task,
	})
}

func drop(ctx *gin.Context){
	db := ctx.MustGet("db").(*gorm.DB)
	user := ctx.MustGet("user").(models.User)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "check id",
		})
		return
	}

	found := deleteOwnTaskFromRedisById(ctx, uint(id))
	if !found {
		var task models.Task
		count := db.Where("id = ? AND author_id = ?", id, user.ID).Delete(&task).RowsAffected
		if count == 0 {
			ctx.JSON(400, gin.H{
				"status": "error",
			})
			return
		}
	}

	ctx.JSON(200, gin.H{
		"status": "successful",
	})
}

func dropSudo(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "check id",
		})
		return
	}

	var task models.Task
	var authorId uint
	_, authorId = deleteTaskFromRedisById(ctx, uint(id))

	result := db.Where("id = ?", id).Find(&task)
	authorId = task.AuthorId
	if result.RowsAffected == 0 {
		ctx.AbortWithStatus(400)
		return
	}
	result = db.Delete(&task)
	num := result.RowsAffected
	if num > 0 {
		notification := models.Notification{
			Text:   "Your task with id="+strconv.Itoa(id)+" was deleted !",
			UserId: authorId,
		}
		db.Create(&notification)

		ctx.JSON(200, gin.H{
			"status": "successful",
		})
	} else {
		ctx.JSON(200, gin.H{
			"status": "error",
		})
	}
}