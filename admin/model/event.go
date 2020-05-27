package model

import (
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
	Speed        int `db:"req_limit"`
	BuyLimit     int `db:"buy_limit"`
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

	sql := "INSERT INTO event(name, product_id, start_time, end_time, total, status, req_limit, buy_limit)VALUES(?,?,?,?,?,?,?,?)"
	_, err = Db.Exec(sql, event.EventName, event.ProductID, event.StartTime, event.EndTime, event.Total, event.Status, event.Speed, event.BuyLimit)
	if err != nil {
		logs.Warn("INSERT INTO event failed, error: %v, sql: %v", err, sql)
		return
	}
	return
}
