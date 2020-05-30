package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
)

const (
	EventStatusNormal  = 0
	EventStatusDisable = 1
	EventStatusExpire  = 2
)

type Event struct {
	EventID   int    `db:"id"`
	EventName string `db:"name"`
	ProductID int    `db:"product_id"`
	StartTime int64  `db:"start_time"`
	EndTime   int64  `db:"end_time"`
	Total     int    `db:"total"`
	Status    int    `db:"status"`

	StartTimeStr string
	EndTimeStr   string
	StatusStr    string
	Speed        int     `db:"req_limit"`
	BuyLimit     int     `db:"buy_limit"`
	BuyRate      float64 `db:"buy_rate"`
}

type ProductInfoConf struct {
	ProductID         int
	StartTime         int64
	EndTime           int64
	Status            int
	Total             int
	Left              int
	OnePersonBuyLimit int
	BuyRate           float64
	SoldMaxLimit      int
}

type EventModel struct {
}

func NewEventModel() *EventModel {
	return &EventModel{}
}

func (p *EventModel) GetEventList() (eventList []*Event, err error) {
	sql := "SELECT id, name, product_id, start_time, end_time, total, status, req_limit, buy_limit FROM event order by id"
	err = Db.Select(&eventList, sql)
	if err != nil {
		logs.Error("SELECT event from mysql failed, error: %v", err)
		return
	}

	for _, v := range eventList {
		t := time.Unix(v.StartTime, 0)
		v.StartTimeStr = t.Format("2006-01-02 15:04:05")
		t = time.Unix(v.EndTime, 0)
		v.EndTimeStr = t.Format("2006-01-02 15:04:05")

		now := time.Now().Unix()
		if now > int64(v.EndTime) {
			v.StatusStr = "ended"
			continue
		}
		if v.Status == EventStatusNormal {
			v.StatusStr = "normal"
		} else if v.Status == EventStatusDisable {
			v.StatusStr = "disable"
		}
	}
	return
}

func (p *EventModel) ProductValid(productId, total int) (valid bool, err error) {
	sql := "SELECT id, name, total, status FROM product where id=?"
	var productList []*Product
	err = Db.Select(&productList, sql, productId)
	if err != nil {
		logs.Warn("SELECT product failed, error: %v", err)
		return
	}

	if len(productList) == 0 {
		err = fmt.Errorf("product[%v] does not exist", productId)
		return
	}

	if total > productList[0].Total {
		err = fmt.Errorf("product[%v] total is not valid", productId)
		return
	}

	valid = true
	return
}

func (p *EventModel) CreateEvent(event *Event) (err error) {
	valid, err := p.ProductValid(event.ProductID, event.Total)
	if err != nil {
		logs.Error("product validation failed, error: %v", err)
		return
	}

	if !valid {
		err = fmt.Errorf("product[%d] validation failed, err: %v", event.ProductID, err)
		logs.Error(err)
		return
	}

	if event.StartTime <= 0 || event.EndTime <= 0 {
		err = fmt.Errorf("invalid start[%v] end[%v] time", event.StartTime, event.EndTime)
		logs.Error(err)
		return
	}

	if event.EndTime <= event.StartTime {
		err = fmt.Errorf("start[%v] is greater than end[%v] time", event.StartTime, event.EndTime)
		logs.Error(err)
		return
	}

	now := time.Now().Unix()
	if event.EndTime <= now || event.StartTime <= now {
		err = fmt.Errorf("start[%v] end[%v] time is less thant now[%v]", event.StartTime, event.EndTime, now)
		logs.Error(err)
		return
	}

	sql := "INSERT INTO event(name, product_id, start_time, end_time, total, status, req_limit, buy_limit, buy_rate)VALUES(?,?,?,?,?,?,?,?,?)"
	_, err = Db.Exec(sql, event.EventName, event.ProductID, event.StartTime, event.EndTime, event.Total, event.Status, event.Speed, event.BuyLimit, event.BuyRate)
	if err != nil {
		logs.Warn("INSERT INTO event failed, error: %v, sql: %v", err, sql)
		return
	}
	logs.Debug("insert into database succ")
	err = p.SyncToEtcd(event)
	if err != nil {
		logs.Warn("sync to etcd failed, error: %v data:%v", err, event)
		return
	}
	return
}

func (p *EventModel) SyncToEtcd(event *Event) (err error) {
	// todo pass in AppConf
	productInfoList, err := loadProductFromEtcd(EtcdKey)

	var productInfo ProductInfoConf
	productInfo.StartTime = event.StartTime
	productInfo.EndTime = event.EndTime
	productInfo.OnePersonBuyLimit = event.BuyLimit
	productInfo.ProductID = event.ProductID
	productInfo.SoldMaxLimit = event.Speed
	productInfo.StartTime = event.StartTime
	productInfo.Status = event.Status
	productInfo.Total = event.Total
	productInfo.BuyRate = event.BuyRate

	productInfoList = append(productInfoList, productInfo)

	data, err := json.Marshal(productInfoList)
	if err != nil {
		logs.Error("json marshal failed, error: %v", err)
		return
	}

	_, err = EtcdClient.Put(context.Background(), EtcdKey, string(data))
	if err != nil {
		logs.Error("put to etcd failed, error: %v, data[%v]", err, string(data))
		return
	}
	return
}

func loadProductFromEtcd(key string) (productInfo []ProductInfoConf, err error) {
	logs.Debug("get from etcd start")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := EtcdClient.Get(ctx, key)
	if err != nil {
		logs.Error("get [%s] from etcd failed, err:%v", key, err)
		return
	}

	logs.Debug("get from etcd succ resp:%v", resp)
	for k, v := range resp.Kvs {
		logs.Debug("key[%v] valud[%v]", k, v)
		err = json.Unmarshal(v.Value, &productInfo)
		if err != nil {
			logs.Error("Unmarshal sec product info failed, err:%v", err)
			return
		}
		logs.Debug("sec info conf is [%v]", productInfo)
	}

	// updateSecProductInfo(appConf, produtInfo)
	// logs.Debug("update product info succ, data:%v", produtInfo)
	// initSecProductWatcher(appConf)
	// logs.Debug("initSecProductWatcher succ")
	return
}
