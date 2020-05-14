package service

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
)

var (
	secKillConf *SecKillConf
)

func NewSecRequest() (secRequest *SecKillRequest) {
	secRequest = &SecKillRequest{
		ResultChan: make(chan *SecKillResult, 1),
	}
	return
}

func userValidate(secRequest *SecKillRequest) (err error) {

	found := false
	for _, ref := range secKillConf.RefererWhiteList {
		if ref == secRequest.ClientReferer {
			found = true
			break
		}
	}

	if !found {
		err = fmt.Errorf("invalid request")
		logs.Warn("user[%d] is rejected by referer, req[%v]", secRequest.UserID, secRequest)
		return
	}
	authData := fmt.Sprintf("%d:%s", secRequest.UserID, secKillConf.CookieSecretKey)
	authSig := fmt.Sprintf("%x", md5.Sum([]byte(authData)))

	if authSig != secRequest.UserAuthSig {
		err = fmt.Errorf("invalid user cookie auth")
		return
	}
	return
}

func SecKill(secRequest *SecKillRequest) (data map[string]interface{}, code int, err error) {
	secKillConf.RwSecKillProductLock.RLock()
	defer secKillConf.RwSecKillProductLock.RUnlock()

	err = userValidate(secRequest)
	if err != nil {
		code = ErrUserValidationFailed
		logs.Warn("userID[%d] is invalid, check failed, req[%v]", secRequest.UserID, secRequest)
		return
	}

	err = antiSpam(secRequest)
	if err != nil {
		code = ErrUserServiceBusy
		logs.Warn("userID[%d] is invalid, check failed, req[%v]", secRequest.UserID, secRequest)
		return
	}
	data, code, err = SecInfoById(secRequest.ProductID)
	if err != nil {
		logs.Warn("userID[%d] SecInfoById failed, req[%v]", secRequest.UserID, secRequest)
		return
	}
	if code != 0 {
		logs.Warn("userID[%d] SecInfoById failed, code[%d] req[%v]", secRequest.UserID, code, secRequest)
		return
	}

	userKey := fmt.Sprintf("%d_%d", secRequest.UserID, secRequest.ProductID)
	secKillConf.SecKillRequestChan <- secRequest
	ticker := time.NewTicker(time.Second * 10)

	defer func() {
		ticker.Stop()
		secKillConf.UserConnMapLock.Lock()
		delete(secKillConf.UserConnMap, userKey)
		secKillConf.UserConnMapLock.Unlock()
	}()

	select {
	case <-ticker.C:
		code = ErrProcessTimeout
		err = fmt.Errorf("request timeout")
		return
	case <-secRequest.CloseNotify:
		code = ErrClientClose
		err = fmt.Errorf("client closed")
		return
	case result := <-secRequest.ResultChan:
		code = result.Code
		data["product_id"] = result.ProductID
		data["token"] = result.Token
		data["user_id"] = result.UserID
	}
	return
}

func SecInfoList() (data []map[string]interface{}, code int, err error) {
	secKillConf.RwSecKillProductLock.RLock()
	defer secKillConf.RwSecKillProductLock.RUnlock()

	for _, v := range secKillConf.ProductInfoMap {
		item, _, err := SecInfoById(v.ProductID)
		if err != nil {
			logs.Error("get product_id[%d] failed, error: %v", v.ProductID, err)
			continue
		}
		data = append(data, item)
	}
	return
}

func SecInfo(productID int) (data []map[string]interface{}, code int, err error) {
	secKillConf.RwSecKillProductLock.RLock()
	defer secKillConf.RwSecKillProductLock.RUnlock()
	item, code, err := SecInfoById(productID)
	if err != nil {
		return
	}
	data = append(data, item)
	return
}

func SecInfoById(productID int) (data map[string]interface{}, code int, err error) {

	secKillConf.RwSecKillProductLock.RLock()
	defer secKillConf.RwSecKillProductLock.RUnlock()

	v, ok := secKillConf.ProductInfoMap[productID]
	if !ok {
		code = ErrNotFoundProductID
		err = fmt.Errorf("not found product_id: %d", productID)
		return
	}

	start := false
	end := false
	status := "success"

	now := time.Now().Unix()

	if now-v.StartTime < 0 {
		start = false
		end = false
		status = "sec kill event hasn't started"
		code = ErrEventNotStart
	}

	if now-v.StartTime > 0 {
		start = true
	}

	if now-v.EndTime > 0 {
		start = false
		end = true
		status = "sec kill event ended"
		code = ErrEventEnded
	}

	if v.Status == ProductStatusForceSaleOut || v.Status == ProductStatusSaleOut {
		start = false
		end = true
		status = "Product is sold out"
		code = ErrEventSoldOut
	}

	data = make(map[string]interface{})
	data["product_id"] = productID
	data["start"] = start
	data["end"] = end
	data["status"] = status

	return
}
