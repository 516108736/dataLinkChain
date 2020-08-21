package routers

import (
	"github.com/516108736/dataLinkChain/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
