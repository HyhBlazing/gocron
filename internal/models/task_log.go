package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
)

type TaskType int8

// 任务执行日志
type TaskLog struct {
	Id         int64        `json:"id" xorm:"bigint pk autoincr"`
	TaskId     int          `json:"task_id" xorm:"int notnull index default 0"`       // 定时ID
	Name       string       `json:"name" xorm:"varchar(32) notnull"`                  // 定时名称
	Spec       string       `json:"spec" xorm:"varchar(64) notnull"`                  // crontab
	Protocol   TaskProtocol `json:"protocol" xorm:"tinyint notnull index"`            // 协议 1:http 2:RPC
	Command    string       `json:"command" xorm:"varchar(256) notnull"`              // URL地址或shell命令
	Timeout    int          `json:"timeout" xorm:"mediumint notnull default 0"`       // 任务执行超时时间(单位秒),0不限制
	RetryTimes int8         `json:"retry_times" xorm:"tinyint notnull default 0"`     // 任务重试次数
	Hostname   string       `json:"hostname" xorm:"varchar(128) notnull default '' "` // RPC主机名，逗号分隔
	StartTime  time.Time    `json:"start_time" xorm:"datetime created"`               // 开始执行时间
	EndTime    time.Time    `json:"end_time" xorm:"datetime updated"`                 // 执行完成（失败）时间
	Status     Status       `json:"status" xorm:"tinyint notnull index default 1"`    // 状态 0:执行失败 1:执行中  2:执行完毕 3:任务取消(上次任务未执行完成) 4:异步执行
	Result     string       `json:"result" xorm:"mediumtext notnull "`                // 执行结果
	TotalTime  int          `json:"total_time" xorm:"-"`                              // 执行总时长
	Tag        string       `json:"tag" xorm:"-"`
	BaseModel  `json:"-" xorm:"-"`
}

func (taskLog *TaskLog) Create() (insertId int64, err error) {
	_, err = Db.Insert(taskLog)
	if err == nil {
		insertId = taskLog.Id
	}

	return
}

// 更新
func (taskLog *TaskLog) Update(id int64, data CommonMap) (int64, error) {
	return Db.Table(taskLog).ID(id).Update(data)
}

func (taskLog *TaskLog) List(params CommonMap) ([]TaskLog, error) {
	taskLog.parsePageAndPageSize(params)
	list := make([]TaskLog, 0)
	session := Db.Alias("tl").Join("LEFT", []string{TablePrefix + "task", "t"}, "tl.task_id = t.id")
	taskLog.parseWhere(session, params)
	err := session.Desc("tl.id").Cols("tl.*").Limit(taskLog.PageSize, taskLog.pageLimitOffset()).Find(&list)
	if len(list) > 0 {
		tagMap := taskLog.queryTaskTagMap(list)
		for i, item := range list {
			endTime := item.EndTime
			if item.Status == Running {
				endTime = time.Now()
			}
			execSeconds := endTime.Sub(item.StartTime).Seconds()
			list[i].TotalTime = int(execSeconds)
			if tag, ok := tagMap[item.TaskId]; ok {
				list[i].Tag = tag
			}
		}
	}

	return list, err
}

func (taskLog *TaskLog) queryTaskTagMap(logs []TaskLog) map[int]string {
	taskIds := make([]interface{}, 0, len(logs))
	seen := make(map[int]struct{}, len(logs))
	for _, item := range logs {
		if item.TaskId <= 0 {
			continue
		}
		if _, ok := seen[item.TaskId]; ok {
			continue
		}
		seen[item.TaskId] = struct{}{}
		taskIds = append(taskIds, item.TaskId)
	}
	if len(taskIds) == 0 {
		return map[int]string{}
	}
	taskList := make([]Task, 0)
	err := Db.Table(new(Task)).In("id", taskIds...).Cols("id,tag").Find(&taskList)
	if err != nil {
		return map[int]string{}
	}
	tagMap := make(map[int]string, len(taskList))
	for _, item := range taskList {
		tagMap[item.Id] = item.Tag
	}
	return tagMap
}

// 清空表
func (taskLog *TaskLog) Clear() (int64, error) {
	return Db.Where("1=1").Delete(taskLog)
}

// 删除N个月前的日志
func (taskLog *TaskLog) Remove(id int) (int64, error) {
	t := time.Now().AddDate(0, -id, 0)
	return Db.Where("start_time <= ?", t.Format(DefaultTimeFormat)).Delete(taskLog)
}

func (taskLog *TaskLog) Total(params CommonMap) (int64, error) {
	session := Db.NewSession().Alias("tl").Join("LEFT", []string{TablePrefix + "task", "t"}, "tl.task_id = t.id")
	taskLog.parseWhere(session, params)
	list := make([]TaskLog, 0)
	err := session.GroupBy("tl.id").Cols("tl.id").Find(&list)
	return int64(len(list)), err
}

// 解析where
func (taskLog *TaskLog) parseWhere(session *xorm.Session, params CommonMap) {
	if len(params) == 0 {
		return
	}
	taskId, ok := params["TaskId"]
	if ok && taskId.(int) > 0 {
		session.And("tl.task_id = ?", taskId)
	}
	keyword, ok := params["Keyword"]
	if ok && strings.TrimSpace(keyword.(string)) != "" {
		k := strings.TrimSpace(keyword.(string))
		if id, err := strconv.Atoi(k); err == nil && id > 0 {
			session.And("(tl.task_id = ? OR tl.name LIKE ?)", id, "%"+k+"%")
		} else {
			session.And("tl.name LIKE ?", "%"+k+"%")
		}
	}
	protocol, ok := params["Protocol"]
	if ok && protocol.(int) > 0 {
		session.And("tl.protocol = ?", protocol)
	}
	status, ok := params["Status"]
	if ok && status.(int) > -1 {
		session.And("tl.status = ?", status)
	}
	tag, ok := params["Tag"]
	if ok && strings.TrimSpace(tag.(string)) != "" {
		session.And("t.tag LIKE ?", "%"+strings.TrimSpace(tag.(string))+"%")
	}
}
