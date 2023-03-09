package main

import (
	"log"

	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos/gorilla"

	// 这里由于命名冲突，进行重命名导入
	gorillaWs "github.com/gorilla/websocket"
)

// port 端口
const port string = ":6868"

func main() {
	// iris的实例
	app := iris.New()

	// iris的子包提供了
	// websocket.DefaultGorillaUpgrader
	// 这个默认的Upgrader提供快速上手
	// 但是这个默认的Upgrader是gorilla实现的
	// 里面默认进行了Origin的检查，不允许跨域
	// 因此需要手动建立一个Upgrader，并覆盖默认
	// 的Origin检查方法，允许所有Origin
	myUpgrader := gorilla.Upgrader(gorillaWs.Upgrader{CheckOrigin: func(*http.Request) bool {
		return true
	}})

	// wsChat 用于聊天室的websocket
	// websocket.Events 实际上一个string -> func的map
	// websocket.OnNativeMessage是一个字符串
	wsChat := websocket.New(myUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			log.Printf("Chat server got: %s from [%s]", msg.Body, nsConn.Conn.ID())

			// 广播到所有本服务的连接
			nsConn.Conn.Server().Broadcast(nsConn, msg)
			return nil
		},
	})
	// Connect事件触发的方法
	wsChat.OnConnect = func(c *websocket.Conn) error {
		log.Printf("[%s] Connected to chat server!", c.ID())
		return nil
	}
	// Disconnect事件触发的方法
	wsChat.OnDisconnect = func(c *websocket.Conn) {
		log.Printf("[%s] Disconnected from chat server", c.ID())
	}

	// wsEcho 用于复读机的websocket
	// websocket.Events 实际上一个string -> func的map
	// websocket.OnNativeMessage是一个字符串
	wsEcho := websocket.New(myUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			log.Printf("Echo server got : %s from [%s]", msg.Body, nsConn.Conn.ID())

			nsConn.Conn.Write(msg)
			return nil
		},
	})
	wsEcho.OnConnect = func(c *websocket.Conn) error {
		log.Printf("[%s] Connected to echo server!", c.ID())
		return nil
	}
	wsEcho.OnDisconnect = func(c *websocket.Conn) {
		log.Printf("[%s] Disconnected from echo server!", c.ID())
	}

	// 根路由
	app.Get("/", func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"msg": "HelloWebSocket",
		})
	})

	// websocket 的路由
	app.Get("/chat", websocket.Handler(wsChat))
	app.Get("/echo", websocket.Handler(wsEcho))

	// 启动项目
	app.Run(iris.Addr(port))
}
