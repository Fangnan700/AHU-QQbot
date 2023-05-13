package server

import (
	"Kira-qbot/model"
	"Kira-qbot/server/ahu"
	"Kira-qbot/server/cqhttp"
	"Kira-qbot/server/gpt"
	"Kira-qbot/server/redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
	"time"
)

// MessagePreprocessing 消息预处理
func MessagePreprocessing(context *gin.Context) {

	// 解析事件到结构体
	var event model.Event
	_ = context.ShouldBindJSON(&event)
	fmt.Printf("%+v\n", event)

	// 提取消息事件信息
	rawMessage := event.RawMessage
	if rawMessage != "" {
		// 解析指令
		preStr := strings.Split(rawMessage, " ")[0]
		switch preStr {
		case "测试":
			test(event)
			break
		case "帮助":
			help(event)
			break
		case "聊天":
			chat(event)
			break
		case "登录教务系统":
			go loginJWxt(event)
			break
		case "今天上什么课":
			go getTodayCourses(event)
			break
		case "明天上什么课":
			go getTomorrowCourses(event)
			break
		case "我的课表":
			go getAllCourses(event)
			break
		case "我的考试":
			go getExams(event)
			break
		default:
			break
		}
	}
}

// 测试通信
func test(event model.Event) {
	var user model.User
	user.UserId = event.UserId
	cqhttp.SendMsg(event, "测试通过")
}

// 帮助菜单
func help(event model.Event) {
	var user model.User
	user.UserId = event.UserId

	payload := "Hi～我是Kira[CQ:face,id=21]" +
		"\n以下是常用指令：" +
		"\n\n01-- 帮助 " +
		"\n|---- 获取使用说明" +
		"\n\n02-- 聊天 " +
		"\n|---- 与ChatGPT聊天" +
		"\n\n03-- 登录教务系统" +
		"\n|---- 登录安徽大学教务系统" +
		"\n\n04-- 今天上什么课" +
		"\n|---- 查询今天的课表" +
		"\n\n05-- 明天上什么课" +
		"\n|---- 查询明天的课表" +
		"\n\n06-- 我的课表" +
		"\n|---- 查询所有课程信息" +
		"\n\n07-- 我的考试" +
		"\n|---- 查询考试信息"

	cqhttp.SendMsg(event, payload)
}

// 聊天模式
func chat(event model.Event) {
	var user model.User
	user.UserId = event.UserId
	respContent := gpt.ChatGPT(event.RawMessage)

	cqhttp.SendMsg(event, respContent)
}

// 登录教务系统
func loginJWxt(event model.Event) {
	var user model.User
	user = redis.GetUser(event.UserId)

	// 读取指令参数
	cmd := strings.Split(event.RawMessage, " ")

	// 读取缓存，检查是否存在
	if user.UserId == 0 {
		user = model.User{}
		user.UserId = event.UserId
		user.AhuInfo.AhuStatus = "noLogin"
	}

	if user.AhuInfo.AhuStatus == "noLogin" {
		_ = ahu.Login(user)
		cqhttp.SendMsg(event, fmt.Sprintf("验证码：\n[CQ:image,file=%d-verify.jpg]\n请按以下格式发送登录信息：\n登录教务系统 账号 密码 验证码", user.UserId))
	}

	if user.AhuInfo.AhuStatus == "logging" {
		if len(cmd) < 4 {
			redis.DeleteUser(user)
			user = model.User{}
			user.UserId = event.UserId
			user.AhuInfo.AhuStatus = "noLogin"
			redis.AddUser(user)
			cqhttp.SendMsg(event, "参数错误[CQ:face,id=15][CQ:face,id=15][CQ:face,id=15]")
		} else {
			user.AhuInfo.AhuNumber = cmd[1]
			user.AhuInfo.AhuPasswd = cmd[2]
			user.AhuInfo.AhuVerifyCode = cmd[3]
			loginStatus := ahu.Login(user)
			if loginStatus == "logged" {
				cqhttp.SendMsg(event, "登录成功[CQ:face,id=144][CQ:face,id=144][CQ:face,id=144]")
			} else {
				redis.DeleteUser(user)
				user = model.User{}
				user.UserId = event.UserId
				user.AhuInfo.AhuStatus = "noLogin"
				redis.AddUser(user)
				cqhttp.SendMsg(event, "登录失败了[CQ:face,id=96][CQ:face,id=96][CQ:face,id=96]")
			}
		}
	}

	if user.AhuInfo.AhuStatus == "logged" {
		cqhttp.SendMsg(event, "已经登录过啦[CQ:face,id=63][CQ:face,id=63][CQ:face,id=63]")
	}
}

// 查询今日课表
func getTodayCourses(event model.Event) {
	now := time.Now()
	weekday := int(now.Weekday())

	var user model.User
	user = redis.GetUser(event.UserId)

	courses := ahu.GetCourses(user, weekday)
	if courses != nil {
		payload := "查询成功[CQ:face,id=144][CQ:face,id=144][CQ:face,id=144]\n课表信息："
		for _, c := range courses {
			payload = payload +
				"\n--" + strings.Split(c.CName, "（")[0] +
				"\n|--教师：" + c.TName +
				"\n|--时间：" + c.CTime +
				"\n|--地点：" + regexp.MustCompile(`\p{Han}+\w+\d+`).FindString(c.Address)
		}

		cqhttp.SendMsg(event, payload)
	} else {
		cqhttp.SendMsg(event, "查询失败了[CQ:face,id=96]")
		redis.DeleteUser(user)
		user = model.User{}
		user.UserId = event.UserId
		user.AhuInfo.AhuStatus = "noLogin"
		redis.AddUser(user)
	}
}

// 查询明日课表
func getTomorrowCourses(event model.Event) {
	now := time.Now()
	weekday := int(now.Weekday()) + 1

	var user model.User
	user = redis.GetUser(event.UserId)

	courses := ahu.GetCourses(user, weekday)
	if courses != nil {

		payload := "查询成功[CQ:face,id=144][CQ:face,id=144][CQ:face,id=144]\n课表信息："

		for _, c := range courses {
			payload = payload +
				"\n--" + strings.Split(c.CName, "（")[0] +
				"\n|--教师：" + c.TName +
				"\n|--时间：" + c.CTime +
				"\n|--地点：" + regexp.MustCompile(`\p{Han}+\w+\d+`).FindString(c.Address)
		}

		cqhttp.SendMsg(event, payload)
	} else {
		cqhttp.SendMsg(event, "查询失败了[CQ:face,id=96]")
		redis.DeleteUser(user)
		user = model.User{}
		user.UserId = event.UserId
		user.AhuInfo.AhuStatus = "noLogin"
		redis.AddUser(user)
	}
}

// 查询所有课表
func getAllCourses(event model.Event) {

	var user model.User
	user = redis.GetUser(event.UserId)
	courses := ahu.GetCourses(user, -1)

	if courses != nil {
		payload := "查询成功[CQ:face,id=144][CQ:face,id=144][CQ:face,id=144]\n课表信息："
		for _, c := range courses {
			payload = payload +
				"\n--" + strings.Split(c.CName, "（")[0] +
				"\n|--教师：" + c.TName +
				"\n|--时间：" + c.CDay + " " + c.CTime +
				"\n|--地点：" + regexp.MustCompile(`\p{Han}+\w+\d+`).FindString(c.Address) +
				"\n|--周次：" + fmt.Sprintf("第%d-%d周", c.StartWeek, c.EndWeek) +
				"\n"
		}
		payload = payload + fmt.Sprintf("\n当前周：第%d周", ahu.GetCurrentWeek())

		cqhttp.SendMsg(event, payload)
	} else {
		cqhttp.SendMsg(event, "查询失败了[CQ:face,id=96]")
		redis.DeleteUser(user)
		user = model.User{}
		user.UserId = event.UserId
		user.AhuInfo.AhuStatus = "noLogin"
		redis.AddUser(user)
	}
}

// 查询考试信息
func getExams(event model.Event) {
	var user model.User
	user = redis.GetUser(event.UserId)
	if user.UserId != 0 {

		exams := ahu.GetExams(user)

		if exams != "" {
			payload := fmt.Sprintf("查询成功[CQ:face,id=144][CQ:face,id=144][CQ:face,id=144]\n考试信息：\n%s", exams)
			cqhttp.SendMsg(event, payload)
		} else {
			cqhttp.SendMsg(event, "查询失败了[CQ:face,id=96]")
		}
	} else {
		user.UserId = event.UserId
		cqhttp.SendMsg(event, "查询失败了[CQ:face,id=96]")
	}
}
