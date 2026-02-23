package manage

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"gopkg.in/macaron.v1"
)

const (
	defaultStatsDays = 90
	maxStatsDays     = 365
)

type nameValueItem struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func Stats(ctx *macaron.Context) string {
	startTime, endTime, days, err := parseStatsTimeRange(ctx)
	jsonResp := utils.JsonResponse{}
	if err != nil {
		return jsonResp.CommonFailure("时间范围参数错误", err)
	}
	overviewAll, err := queryOverviewAll()
	if err != nil {
		logger.Error(err)
		return jsonResp.CommonFailure("查询总量统计失败", err)
	}

	taskList, err := queryTaskRange(startTime, endTime)
	if err != nil {
		logger.Error(err)
		return jsonResp.CommonFailure("查询任务统计失败", err)
	}
	taskLogList, err := queryTaskLogRange(startTime, endTime)
	if err != nil {
		logger.Error(err)
		return jsonResp.CommonFailure("查询定时日志统计失败", err)
	}
	userList, err := queryUserRange(startTime, endTime)
	if err != nil {
		logger.Error(err)
		return jsonResp.CommonFailure("查询用户统计失败", err)
	}
	loginLogList, err := queryLoginLogRange(startTime, endTime)
	if err != nil {
		logger.Error(err)
		return jsonResp.CommonFailure("查询登录日志统计失败", err)
	}

	tagMap := map[string]int{}
	taskCreatedMap := map[string]int{}
	for _, item := range taskList {
		tag := strings.TrimSpace(item.Tag)
		if tag == "" {
			tag = "未设置标签"
		}
		tagMap[tag]++
		taskCreatedMap[item.Created.Format("2006-01-02")]++
	}

	protocolMap := map[string]int{"http": 0, "shell": 0, "其他": 0}
	durationMap := map[string]int{"0-1s": 0, "1-3s": 0, "3-10s": 0, "10-30s": 0, "30s+": 0}
	statusMap := map[string]int{"失败": 0, "执行中": 0, "成功": 0, "取消": 0}
	for _, item := range taskLogList {
		if item.Protocol == models.TaskHTTP {
			protocolMap["http"]++
		} else if item.Protocol == models.TaskRPC {
			protocolMap["shell"]++
		} else {
			protocolMap["其他"]++
		}

		end := item.EndTime
		if item.Status == models.Running {
			end = time.Now()
		}
		durationSec := int(end.Sub(item.StartTime).Seconds())
		if durationSec <= 1 {
			durationMap["0-1s"]++
		} else if durationSec <= 3 {
			durationMap["1-3s"]++
		} else if durationSec <= 10 {
			durationMap["3-10s"]++
		} else if durationSec <= 30 {
			durationMap["10-30s"]++
		} else {
			durationMap["30s+"]++
		}

		if item.Status == models.Failure {
			statusMap["失败"]++
		} else if item.Status == models.Running {
			statusMap["执行中"]++
		} else if item.Status == models.Finish {
			statusMap["成功"]++
		} else if item.Status == models.Cancel {
			statusMap["取消"]++
		}
	}

	loginUserMap := map[string]int{}
	loginIpMap := map[string]int{}
	for _, item := range loginLogList {
		username := strings.TrimSpace(item.Username)
		if username == "" {
			username = "未知用户"
		}
		loginUserMap[username]++
		ip := strings.TrimSpace(item.Ip)
		if ip == "" {
			ip = "未知IP"
		}
		loginIpMap[ip]++
	}

	result := map[string]interface{}{
		"range": map[string]interface{}{
			"start_time": startTime.Format(models.DefaultTimeFormat),
			"end_time":   endTime.Format(models.DefaultTimeFormat),
			"days":       days,
		},
		"overview": map[string]int{
			"task_total":      len(taskList),
			"task_log_total":  len(taskLogList),
			"user_total":      len(userList),
			"login_log_total": len(loginLogList),
		},
		"overview_all":      overviewAll,
		"task_tag":          toNameValueList(tagMap),
		"task_created":      toNameValueList(taskCreatedMap),
		"task_log_protocol": toNameValueList(protocolMap),
		"task_log_duration": toNameValueList(durationMap),
		"task_log_status":   toNameValueList(statusMap),
		"login_user":        toNameValueList(loginUserMap),
		"login_ip":          toNameValueList(loginIpMap),
	}

	return jsonResp.Success(utils.SuccessContent, result)
}

func queryOverviewAll() (map[string]int, error) {
	taskTotal, err := models.Db.Count(new(models.Task))
	if err != nil {
		return nil, err
	}
	taskLogTotal, err := models.Db.Count(new(models.TaskLog))
	if err != nil {
		return nil, err
	}
	userTotal, err := models.Db.Count(new(models.User))
	if err != nil {
		return nil, err
	}
	loginLogTotal, err := models.Db.Count(new(models.LoginLog))
	if err != nil {
		return nil, err
	}
	return map[string]int{
		"task_total":      int(taskTotal),
		"task_log_total":  int(taskLogTotal),
		"user_total":      int(userTotal),
		"login_log_total": int(loginLogTotal),
	}, nil
}

func parseStatsTimeRange(ctx *macaron.Context) (time.Time, time.Time, int, error) {
	now := time.Now()
	endTime := now
	startTime := now.AddDate(0, 0, -(defaultStatsDays - 1))
	days := defaultStatsDays

	startStr := strings.TrimSpace(ctx.QueryTrim("start_time"))
	endStr := strings.TrimSpace(ctx.QueryTrim("end_time"))
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
	if startStr != "" && endStr != "" {
		start, err := time.ParseInLocation("2006-01-02", startStr, time.Local)
		if err != nil {
			return startTime, endTime, days, err
		}
		end, err := time.ParseInLocation("2006-01-02", endStr, time.Local)
		if err != nil {
			return startTime, endTime, days, err
		}
		if end.After(todayEnd) {
			return startTime, endTime, days, errors.New("结束时间不能超过当天")
		}
		startTime = start
		endTime = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	} else {
		daysQuery := ctx.QueryTrim("days")
		if daysQuery != "" {
			parsedDays, err := strconv.Atoi(daysQuery)
			if err == nil && parsedDays > 0 {
				days = parsedDays
			}
		}
	}

	if days > maxStatsDays {
		days = maxStatsDays
	}
	if startStr == "" || endStr == "" {
		startTime = now.AddDate(0, 0, -(days - 1))
		endTime = now
	}
	if startTime.After(endTime) {
		startTime, endTime = endTime, startTime
	}
	rangeDays := int(endTime.Sub(startTime).Hours()/24) + 1
	if rangeDays > maxStatsDays {
		startTime = endTime.AddDate(0, 0, -(maxStatsDays - 1))
		rangeDays = maxStatsDays
	}

	return startTime, endTime, rangeDays, nil
}

func queryTaskRange(startTime, endTime time.Time) ([]models.Task, error) {
	list := make([]models.Task, 0)
	err := models.Db.Where("created >= ? AND created <= ?", startTime.Format(models.DefaultTimeFormat), endTime.Format(models.DefaultTimeFormat)).
		Cols("id,name,tag,protocol,created").
		Find(&list)
	return list, err
}

func queryTaskLogRange(startTime, endTime time.Time) ([]models.TaskLog, error) {
	list := make([]models.TaskLog, 0)
	err := models.Db.Where("start_time >= ? AND start_time <= ?", startTime.Format(models.DefaultTimeFormat), endTime.Format(models.DefaultTimeFormat)).
		Cols("id,protocol,start_time,end_time,status").
		Find(&list)
	return list, err
}

func queryUserRange(startTime, endTime time.Time) ([]models.User, error) {
	list := make([]models.User, 0)
	err := models.Db.Where("created >= ? AND created <= ?", startTime.Format(models.DefaultTimeFormat), endTime.Format(models.DefaultTimeFormat)).
		Cols("id,is_admin,created").
		Find(&list)
	return list, err
}

func queryLoginLogRange(startTime, endTime time.Time) ([]models.LoginLog, error) {
	list := make([]models.LoginLog, 0)
	err := models.Db.Where("created >= ? AND created <= ?", startTime.Format(models.DefaultTimeFormat), endTime.Format(models.DefaultTimeFormat)).
		Cols("id,username,ip,created").
		Find(&list)
	return list, err
}

func toNameValueList(data map[string]int) []nameValueItem {
	list := make([]nameValueItem, 0, len(data))
	for name, value := range data {
		if value <= 0 {
			continue
		}
		list = append(list, nameValueItem{
			Name:  name,
			Value: value,
		})
	}
	// 简单稳定排序：值降序，值相同按名称升序
	for i := 0; i < len(list)-1; i++ {
		for j := i + 1; j < len(list); j++ {
			if list[j].Value > list[i].Value || (list[j].Value == list[i].Value && list[j].Name < list[i].Name) {
				list[i], list[j] = list[j], list[i]
			}
		}
	}
	return list
}
