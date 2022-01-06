package routers

import (
	"hellclientswitch/modules/systems/gateway"

	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/herb/middleware/router"
	"github.com/herb-go/herb/middleware/router/httprouter"
)

//WsMiddlewares middlewares which should be used on router.
var WsMiddlewares = func() middleware.Middlewares {
	return middleware.Middlewares{}
}

//RouterWsFactory ws router factory.
var RouterWsFactory = router.NewFactory(func() router.Router {
	var Router = httprouter.New()
	//Put your router configure code here
	Router.ALL("/ws").HandleFunc(gateway.Action)
	return Router
})
