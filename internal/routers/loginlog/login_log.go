package loginlog

import (
	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"github.com/ouqiang/gocron/internal/routers/base"
	macaron "gopkg.in/macaron.v1"
)

func Index(ctx *macaron.Context) string {
	loginLogModel := new(models.LoginLog)
	params := parseQueryParams(ctx)
	total, err := loginLogModel.Total(params)
	if err != nil {
		logger.Error(err)
	}
	loginLogs, err := loginLogModel.List(params)
	if err != nil {
		logger.Error(err)
	}

	jsonResp := utils.JsonResponse{}

	return jsonResp.Success(utils.SuccessContent, map[string]interface{}{
		"total": total,
		"data":  loginLogs,
	})
}

func Usernames(ctx *macaron.Context) string {
	userModel := new(models.User)
	list, err := userModel.NameList(ctx.QueryTrim("keyword"))
	jsonResp := utils.JsonResponse{}
	if err != nil {
		logger.Error(err)
		return jsonResp.CommonFailure("查询用户名失败", err)
	}
	if list == nil {
		list = make([]string, 0)
	}
	return jsonResp.Success(utils.SuccessContent, list)
}

func parseQueryParams(ctx *macaron.Context) models.CommonMap {
	params := models.CommonMap{}
	params["Username"] = ctx.QueryTrim("username")
	params["Ip"] = ctx.QueryTrim("ip")
	params["StartTime"] = ctx.QueryTrim("start_time")
	params["EndTime"] = ctx.QueryTrim("end_time")
	base.ParsePageAndPageSize(ctx, params)
	return params
}
