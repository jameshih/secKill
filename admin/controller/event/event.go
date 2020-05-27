package event

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jameshih/secKill/admin/model"
)

type EventController struct {
	beego.Controller
}

func (p *EventController) ListEvent() {
	eventModel := model.NewEventModel()
	eventList, err := eventModel.GetEventList()
	if err != nil {
		logs.Warn("get event list failed, error: %v", err)
		return
	}

	p.Data["eventList"] = eventList
	p.Layout = "layout/layout.html"
	p.TplName = "event/list.html"
}

func (p *EventController) NewEvent() {

	p.Layout = "layout/layout.html"
	p.TplName = "event/form.html"
}

func (p *EventController) CreateEvent() {
	var err error
	defer func() {
		if err != nil {
			p.Data["Error"] = err.Error()
			p.Layout = "layout/layout.html"
			p.TplName = "layout/error.html"
		}
	}()

	eventName := p.GetString("event_name")
	if len(eventName) == 0 {
		err = fmt.Errorf("invalid event name")
		logs.Warn(err)
		return
	}

	productId, err := p.GetInt("product_id")
	if err != nil {
		err = fmt.Errorf("invalid product id, error: %v", err)
		logs.Warn(err)
		return
	}

	startTime, err := p.GetInt64("start_time")
	if err != nil {
		err = fmt.Errorf("invalid start time, error: %v", err)
		logs.Warn(err)
		return
	}

	endTime, err := p.GetInt64("end_time")
	if err != nil {
		err = fmt.Errorf("invalid end time, error: %v", err)
		logs.Warn(err)
		return
	}

	total, err := p.GetInt("total")
	if err != nil {
		err = fmt.Errorf("invalid event total, error: %v", err)
		logs.Warn(err)
		return
	}

	speed, err := p.GetInt("req_limit")
	if err != nil {
		err = fmt.Errorf("invalid req limit, error: %v", err)
		logs.Warn(err)
		return
	}

	buyLimit, err := p.GetInt("buy_limit")
	if err != nil {
		err = fmt.Errorf("invalid buy limit, error: %v", err)
		logs.Warn(err)
		return
	}

	eventModel := model.NewEventModel()
	event := model.Event{
		EventName: eventName,
		ProductID: productId,
		StartTime: startTime,
		EndTime:   endTime,
		Total:     total,
		Speed:     speed,
		BuyLimit:  buyLimit,
	}

	err = eventModel.CreateEvent(&event)
	if err != nil {
		err = fmt.Errorf("failed to submit, error: %v", err)
		logs.Warn(err)
		return
	}

	logs.Debug("event name[%s], product id[%d], start time[%d], end time[%d], total[%d], status[%d]", event.EventName, event.ProductID, event.StartTime, event.EventID, event.Total, event.Status)

	p.Redirect("/event/list", http.StatusMovedPermanently)
	return
}
