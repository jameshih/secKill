package controller

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jameshih/secKill/proxy/service"
)

type SkillController struct {
	beego.Controller
}

func (p *SkillController) SecKill() {
	productID, err := p.GetInt("product_id")
	result := make(map[string]interface{})
	result["code"] = 0
	result["message"] = "success"
	defer func() {
		p.Data["json"] = result
		p.ServeJSON()
	}()

	if err != nil {
		result["code"] = 1001
		result["message"] = "invalid product_id"
		return
	}

	source := p.GetString("src")
	authCode := p.GetString("authcode")
	secTime := p.GetString("time")
	nounce := p.GetString("nounce")

	secRequest := service.NewSecRequest()
	secRequest.ProductID = productID
	secRequest.Source = source
	secRequest.AuthCode = authCode
	secRequest.SecTime = secTime
	secRequest.Nounce = nounce
	secRequest.UserAuthSig = p.Ctx.GetCookie("userAuthSig")
	secRequest.UserID, err = strconv.Atoi(p.Ctx.GetCookie("userID"))
	if err != nil {
		result["code"] = service.ErrInvalidRequest
		result["message"] = fmt.Sprintf("invalid cookie")
		return
	}
	secRequest.AccessTime = time.Now()
	if len(p.Ctx.Request.RemoteAddr) > 0 {
		secRequest.ClientAddr = strings.Split(p.Ctx.Request.RemoteAddr, ":")[0]
	}

	secRequest.ClientReferer = p.Ctx.Request.Referer()
	secRequest.CloseNotify = p.Ctx.ResponseWriter.CloseNotify()

	logs.Debug("client request:[%v]", secRequest)
	if err != nil {
		result["code"] = service.ErrInvalidRequest
		result["message"] = fmt.Sprintf("invalid cookie:userID")
		return
	}

	data, code, err := service.SecKill(secRequest)
	if err != nil {
		result["code"] = code
		result["message"] = err.Error()
		return
	}

	result["data"] = data
	result["code"] = code
	return
}

func (p *SkillController) SecInfo() {
	productID, err := p.GetInt("product_id")
	result := make(map[string]interface{})
	result["code"] = 0
	result["message"] = "success"
	defer func() {
		p.Data["json"] = result
		p.ServeJSON()
	}()

	if err != nil {
		data, code, err := service.SecInfoList()
		if err != nil {
			logs.Error("invalid request, get product_id failed, error: %v", err)
			return
		}
		result["code"] = code
		result["data"] = data

	} else {
		data, code, err := service.SecInfo(productID)
		if err != nil {
			result["code"] = code
			result["message"] = err.Error()
			logs.Error("invalid request, get product_id failed, error: %v", err)
			return
		}
		result["code"] = code
		result["data"] = data
	}
}
