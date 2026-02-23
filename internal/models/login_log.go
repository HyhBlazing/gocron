package models

import (
	"strings"
	"time"

	"github.com/go-xorm/xorm"
)

// 用户登录日志

type LoginLog struct {
	Id        int       `json:"id" xorm:"pk autoincr notnull "`
	Username  string    `json:"username" xorm:"varchar(32) notnull"`
	Ip        string    `json:"ip" xorm:"varchar(15) not null"`
	Created   time.Time `json:"created" xorm:"datetime notnull created"`
	BaseModel `json:"-" xorm:"-"`
}

func (log *LoginLog) Create() (insertId int, err error) {
	_, err = Db.Insert(log)
	if err == nil {
		insertId = log.Id
	}

	return
}

func (log *LoginLog) List(params CommonMap) ([]LoginLog, error) {
	log.parsePageAndPageSize(params)
	list := make([]LoginLog, 0)
	session := Db.NewSession().Desc("id")
	log.parseWhere(session, params)
	err := session.Limit(log.PageSize, log.pageLimitOffset()).Find(&list)

	return list, err
}

func (log *LoginLog) Total(params CommonMap) (int64, error) {
	session := Db.NewSession()
	log.parseWhere(session, params)
	return session.Count(log)
}

func (log *LoginLog) UsernameList(keyword string) ([]string, error) {
	usernames := make([]string, 0)
	session := Db.NewSession().Distinct("username").Cols("username").Desc("id")
	if strings.TrimSpace(keyword) != "" {
		session.And("username LIKE ?", "%"+strings.TrimSpace(keyword)+"%")
	}
	type usernameItem struct {
		Username string `xorm:"username"`
	}
	list := make([]usernameItem, 0)
	err := session.Limit(50).Find(&list)
	if err != nil {
		return usernames, err
	}
	for _, item := range list {
		if strings.TrimSpace(item.Username) != "" {
			usernames = append(usernames, item.Username)
		}
	}

	return usernames, nil
}

func (log *LoginLog) parseWhere(session *xorm.Session, params CommonMap) {
	if len(params) == 0 {
		return
	}
	username, ok := params["Username"]
	if ok && strings.TrimSpace(username.(string)) != "" {
		session.And("username LIKE ?", "%"+strings.TrimSpace(username.(string))+"%")
	}
	ip, ok := params["Ip"]
	if ok && strings.TrimSpace(ip.(string)) != "" {
		session.And("ip LIKE ?", "%"+strings.TrimSpace(ip.(string))+"%")
	}
	startTime, ok := params["StartTime"]
	if ok && strings.TrimSpace(startTime.(string)) != "" {
		session.And("created >= ?", strings.TrimSpace(startTime.(string)))
	}
	endTime, ok := params["EndTime"]
	if ok && strings.TrimSpace(endTime.(string)) != "" {
		session.And("created <= ?", strings.TrimSpace(endTime.(string)))
	}
}
