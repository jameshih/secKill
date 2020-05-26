package event

import (
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

}

func (p *EventController) CreateEvent() {

}
