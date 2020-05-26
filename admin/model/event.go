package model

import (
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
)

type Event struct {
	EventID   int    `db:"id"`
	EventName string `db:"name"`
	ProductID int    `db:"product_id"`
	StartTime int    `db:"start_time"`
	EndTime   int    `db:"end_time"`
	Total     int    `db:"total"`
	Status    int    `db:"status"`
}

type EventModel struct {
}

func NewEventModel() *EventModel {
	return &EventModel{}
}

func (p *EventModel) GetEventList() (eventList []*Event, err error) {
	sql := "SELECT id, name, product_id, start_time, end_time, total, status FROM event order by id"
	err = Db.Select(&eventList, sql)
	if err != nil {
		logs.Error("SELECT event from mysql failed, error: %v", err)
		return
	}
	return
}

func (p *EventModel) CreateEvent(event *Event) (err error) {
	sql := "INSERT INTO event(name, product_id, start_time, end_time, total, status)VALUES(?,?,?,?,?,?)"
	_, err = Db.Exec(sql, event.EventName, event.ProductID, event.StartTime, event.EndTime, event.Total, event.Status)
	if err != nil {
		logs.Warn("INSERT INTO event failed, error: %v, sql: %v", err, sql)
		return
	}
	return
}
