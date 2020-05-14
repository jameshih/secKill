package router

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jameshih/test/SecProxy/controller"
)

func init() {
	logs.Debug("enter router init")
	beego.Router("/seckill", &controller.SkillController{}, "*:SecKill")
	beego.Router("/secinfo", &controller.SkillController{}, "*:SecInfo")
}
