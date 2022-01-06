package gateway

import (
	"bytes"
	"net/http"

	"github.com/herb-go/connections"
	"github.com/herb-go/connections/contexts"
	"github.com/herb-go/connections/websocket"
	"github.com/herb-go/uniqueid"
	"github.com/herb-go/util"
)

//ModuleName module name
const ModuleName = "900systems.gateway"

type Gateway struct {
	Gateway *connections.Gateway
	*contexts.Contexts
}

func (g *Gateway) DoBroadcast(m *connections.Message) {
	go func() {
		list := g.Gateway.ListConn()
		for _, v := range list {
			if v.ID() != m.Conn.ID() {
				v.Send(m.Message)
			}
		}
	}()
}

//OnMessage called when connection message received.
func (g *Gateway) OnMessage(m *connections.Message) {
	if bytes.HasPrefix(m.Message, CommandBroadcast) {
		g.DoBroadcast(m)
	}
}

//OnError called when onconnection error raised.
func (g *Gateway) OnError(e *connections.Error) {
	util.LogError(e.Error)
}

func New() *Gateway {
	g := &Gateway{
		Gateway:  connections.NewGateway(),
		Contexts: contexts.New(),
	}
	g.Gateway.IDGenerator = uniqueid.DefaultGenerator.GenerateID
	return g
}

var DefaultGateway = New()

func Start() {
	go connections.Consume(DefaultGateway.Gateway, DefaultGateway)
}

func Stop() {
	DefaultGateway.Gateway.Stop()
}

func Action(w http.ResponseWriter, r *http.Request) {
	wc, err := websocket.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	_, err = DefaultGateway.Gateway.Register(wc)
	if err != nil {
		panic(err)
	}

}
func init() {
	util.RegisterModule(ModuleName, func() {
		//Init registered initator which registered by RegisterInitiator
		//util.RegisterInitiator(ModuleName, "func", func(){})
		util.InitOrderByName(ModuleName)
		Start()
		util.OnQuit(Stop)
	})
}
