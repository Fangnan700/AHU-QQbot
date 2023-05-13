package ahu

import (
	"Kira-qbot/model"
	"Kira-qbot/server/redis"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var CurWeek = GetCurrentWeek()
var CqHttpPath string

func init() {
	var config model.Config

	configFileName := "config/Config.yml"
	configFile, _ := os.Open(configFileName)
	decoder := yaml.NewDecoder(configFile)
	_ = decoder.Decode(&config)
	CqHttpPath = config.CqHttpPath
}

func Login(user model.User) string {

	var err error
	var request *http.Request
	var response *http.Response
	var cookie http.Cookie

	JwxtUrl := "https://jwxt3.ahu.edu.cn"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// 登录预备状态
	if user.AhuInfo.AhuStatus == "noLogin" {
		request, err = http.NewRequest("GET", JwxtUrl, nil)
		if err != nil {
			LogError(err)
		}

		// 请求教务系统登录页
		request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
		response, err = client.Do(request)
		if err != nil {
			LogError(err)
		}

		// 提取cookie
		for _, _cookie := range response.Cookies() {
			if _cookie.Name == "ASP.NET_SessionId" {
				cookie.Name = "ASP.NET_SessionId"
				cookie.Value = _cookie.Value
				cookie.Path = _cookie.Path
			}
		}

		// 解析页面元素
		var doc *goquery.Document
		respBody := response.Body
		doc, err = goquery.NewDocumentFromReader(respBody)
		if err != nil {
			LogError(err)
		}

		// 收集加密参数
		et, _ := doc.Find("#txtKeyExponent").Attr("value")
		nt, _ := doc.Find("#txtKeyModulus").Attr("value")

		// 收集表单参数
		LASTFOCUS, _ := doc.Find("#__LASTFOCUS").Attr("value")
		VIEWSTATE, _ := doc.Find("#__VIEWSTATE").Attr("value")
		VIEWSTATEGENERATOR, _ := doc.Find("#__VIEWSTATEGENERATOR").Attr("value")
		EVENTTARGET, _ := doc.Find("#__EVENTTARGET").Attr("value")
		EVENTARGUMENT, _ := doc.Find("#__EVENTARGUMENT").Attr("value")

		verifySrc, _ := doc.Find("#icode").Attr("src")
		verifySrc = fmt.Sprintf("%s%s", JwxtUrl, verifySrc)

		// 请求验证码
		request, err = http.NewRequest("GET", verifySrc, nil)
		request.AddCookie(&cookie)
		response, err = client.Do(request)
		if err != nil {
			LogError(err)
		}

		imgFile, _ := os.Create(fmt.Sprintf("%s/data/images/%d-verify.jpg", CqHttpPath, user.UserId))
		_, err = io.Copy(imgFile, response.Body)
		if err != nil {
			LogError(err)
		}

		// 写入用户信息
		user.AhuInfo.AhuCookie = cookie
		user.AhuInfo.ET = et
		user.AhuInfo.NT = nt
		user.AhuInfo.LASTFOCUS = LASTFOCUS
		user.AhuInfo.VIEWSTATE = VIEWSTATE
		user.AhuInfo.VIEWSTATEGENERATOR = VIEWSTATEGENERATOR
		user.AhuInfo.EVENTTARGET = EVENTTARGET
		user.AhuInfo.EVENTARGUMENT = EVENTARGUMENT
		user.AhuInfo.AhuStatus = "logging"

		// 写入缓存
		redis.AddUser(user)
		LogInfo(fmt.Sprintf("新增用户：%d", user.UserId))

		return "logging"
	}

	// 登录就绪状态
	if user.AhuInfo.AhuStatus == "logging" {

		user.AhuInfo.AhuPasswdRsa = AhuPasswdEncode(user.AhuInfo.ET, user.AhuInfo.NT, user.AhuInfo.AhuPasswd)

		// 创建请求
		payload := url.Values{}
		payload.Add("__LASTFOCUS", user.AhuInfo.LASTFOCUS)
		payload.Add("__VIEWSTATE", user.AhuInfo.VIEWSTATE)
		payload.Add("__VIEWSTATEGENERATOR", user.AhuInfo.VIEWSTATEGENERATOR)
		payload.Add("__EVENTTARGET", user.AhuInfo.EVENTTARGET)
		payload.Add("__EVENTARGUMENT", user.AhuInfo.EVENTARGUMENT)
		payload.Add("txtUserName", user.AhuInfo.AhuNumber)
		payload.Add("TextBox2", user.AhuInfo.AhuPasswdRsa)
		payload.Add("txtSecretCode", user.AhuInfo.AhuVerifyCode)
		payload.Add("txtKeyExponent", user.AhuInfo.ET)
		payload.Add("txtKeyModulus", user.AhuInfo.NT)
		payload.Add("RadioButtonList1", "学生")
		payload.Add("Button1", "登录")

		request, err = http.NewRequest("POST", JwxtUrl, strings.NewReader(payload.Encode()))
		request.Header.Add("Cookie", fmt.Sprintf("ASP.NET_SessionId=%s", user.AhuInfo.AhuCookie.Value))
		request.Header.Add("Accept", "*/*")
		request.Header.Add("Host", "jwxt3.ahu.edu.cn")
		request.Header.Add("Connection", "keep-alive")
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")

		response, err = client.Do(request)
		if err != nil {
			LogError(err)
		}

		// 登录成功
		if response.StatusCode == 302 {
			user.AhuInfo.AhuStatus = "logged"
			redis.AddUser(user)
			LogInfo(fmt.Sprintf("用户登录成功：%d", user.UserId))
			return "logged"
		}
	}

	// 已登录状态
	if user.AhuInfo.AhuStatus == "logged" {
		return "logged"
	}

	// 登录失败
	LogInfo(fmt.Sprintf("用户登录失败：%d", user.UserId))
	return "fail"
}

func GetCourses(user model.User, weekday int) []model.Course {
	var courses []model.Course

	// 检查是否已登录
	if user.AhuInfo.AhuStatus == "logged" {
		CoursesUrl := fmt.Sprintf("https://jwxt3.ahu.edu.cn/xsxkqk.aspx?xh=%s", user.AhuInfo.AhuNumber)

		client := &http.Client{}
		request, _ := http.NewRequest("POST", CoursesUrl, nil)
		request.Header.Add("Cookie", fmt.Sprintf("ASP.NET_SessionId=%s", user.AhuInfo.AhuCookie.Value))
		request.Header.Add("Referer", fmt.Sprintf("https://jwxt3.ahu.edu.cn/xs_main.aspx?xh=%s", user.AhuInfo.AhuNumber))
		request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
		request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		request.Header.Add("Host", "jwxt3.ahu.edu.cn")
		request.Header.Add("Connection", "keep-alive")

		response, err := client.Do(request)
		if err != nil {
			LogError(err)
		}

		// 解析页面元素
		body := response.Body

		var doc *goquery.Document
		doc, err = goquery.NewDocumentFromReader(body)
		if err != nil {
			LogError(err)
		}

		trIndex := 0
		doc.Find("#DBGrid tr").Each(func(i int, s_tr *goquery.Selection) {
			if trIndex >= 1 {
				var item []string
				s_tr.Find("td").Each(func(j int, s_td *goquery.Selection) {
					item = append(item, strings.TrimSpace(s_td.Text()))
				})

				// 提取课程名称
				cName := item[2]

				// 提取教师名称
				tName := item[5]

				// 提取星期几
				cDay := regexp.MustCompile(`周.{1}`).FindString(item[8])
				cDayNum := 0
				switch cDay {
				case "周日":
					cDayNum = 0
					break
				case "周一":
					cDayNum = 1
					break
				case "周二":
					cDayNum = 2
					break
				case "周三":
					cDayNum = 3
					break
				case "周四":
					cDayNum = 4
					break
				case "周五":
					cDayNum = 5
					break
				case "周六":
					cDayNum = 6
					break
				default:
					break
				}

				// 提取第几周
				weeks := regexp.MustCompile(`\d+-\d+`).FindString(item[8])
				startWeek, _ := strconv.Atoi(strings.Split(weeks, "-")[0])
				endWeek, _ := strconv.Atoi(strings.Split(weeks, "-")[1])

				// 提取第几节
				class := regexp.MustCompile(`第\d+(,\d+)*节`).FindAllString(item[8], 3)

				// 提取上课时间
				var max, min, cIndex int
				var startTime, endTime string
				for _, str := range class {
					// 查找匹配的字符串
					matches := regexp.MustCompile(`\d+`).FindAllString(str, -1)

					// 将字符串转换为数字，并更新最大值和最小值
					for _, match := range matches {
						num, _ := strconv.Atoi(match)
						if num > max {
							max = num
						}
						if num < min || min == 0 {
							min = num
						}
					}
				}

				switch min {
				case 1:
					startTime = "8:20"
					break
				case 3:
					startTime = "10:20"
					break
				case 5:
					startTime = "14:00"
					break
				case 7:
					startTime = "15:50"
				case 9:
					startTime = "19:00"
				default:
					break
				}

				switch max {
				case 2:
					endTime = "10:00"
				case 3:
					endTime = "11:05"
					break
				case 4:
					endTime = "12:00"
					break
				case 6:
					endTime = "15:40"
					break
				case 7:
					endTime = "15:40"
				case 8:
					endTime = "17:35"
					break
				case 11:
					endTime = "21:30"
					break
				default:
					break
				}
				cIndex = min
				cTime := fmt.Sprintf("%s-%s", startTime, endTime)

				// 提取上课地点
				address := strings.Split(item[9], ";")[0]

				// 创建课程对象
				course := model.Course{
					cIndex,
					class,
					cName,
					tName,
					cDayNum,
					cDay,
					cTime,
					address,
					startWeek,
					endWeek,
					CurWeek,
				}
				courses = append(courses, course)
			}
			trIndex += 1
		})
	}

	// weekday = -1 表示查询所有课程
	if weekday != -1 {

		var data []model.Course
		switch weekday {
		case 0:
			for _, course := range courses {
				if course.CDay == "周日" && course.CurrentWeek >= course.StartWeek && course.CurrentWeek <= course.EndWeek {
					data = append(data, course)
				}
			}
			break
		case 1:
			for _, course := range courses {
				if course.CDay == "周一" && course.CurrentWeek >= course.StartWeek && course.CurrentWeek <= course.EndWeek {
					data = append(data, course)
				}
			}
			break
		case 2:
			for _, course := range courses {
				if course.CDay == "周二" && course.CurrentWeek >= course.StartWeek && course.CurrentWeek <= course.EndWeek {
					data = append(data, course)
				}
			}
			break
		case 3:
			for _, course := range courses {
				if course.CDay == "周三" && course.CurrentWeek >= course.StartWeek && course.CurrentWeek <= course.EndWeek {
					data = append(data, course)
				}
			}
			break
		case 4:
			for _, course := range courses {
				if course.CDay == "周四" && course.CurrentWeek >= course.StartWeek && course.CurrentWeek <= course.EndWeek {
					data = append(data, course)
				}
			}
			break
		case 5:
			for _, course := range courses {
				if course.CDay == "周五" && course.CurrentWeek >= course.StartWeek && course.CurrentWeek <= course.EndWeek {
					data = append(data, course)
				}
			}
			break
		case 6:
			for _, course := range courses {
				if course.CDay == "周六" && course.CurrentWeek >= course.StartWeek && course.CurrentWeek <= course.EndWeek {
					data = append(data, course)
				}
			}
			break
		}
		sort.Slice(data, func(i, j int) bool {
			return data[i].Index <= data[j].Index
		})
		return data
	}

	sort.Slice(courses, func(i, j int) bool {
		return courses[i].CDayNum <= courses[j].CDayNum
	})
	return courses
}

// GetExams 考试查询
func GetExams(user model.User) string {
	client := http.Client{}
	examUrl := fmt.Sprintf("http://kskw.ahu.edu.cn/bkcx.asp?xh=%s", user.AhuInfo.AhuNumber)

	request, _ := http.NewRequest("GET", examUrl, nil)
	response, _ := client.Do(request)

	htmlBytes, _ := io.ReadAll(response.Body)
	html := string(htmlBytes)
	body := strings.Split(html, "信息如下：")[1]
	body = strings.ReplaceAll(body, "<br>", "\n")
	body = strings.TrimSpace(body)
	body = fmt.Sprintf("学号：%s", body)

	if strings.Contains(body, "你不是学生") {
		return ""
	}
	return body
}
