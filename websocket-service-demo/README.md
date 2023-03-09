# websocket-service-demo

golang 在 iris 中使用 websocket 的极简示例，搭建了两个 websocket 服务，分别是 **echo 复读机服务**与**实时公共聊天室**。

## 启动

首先需要配置go开发环境，下载 [go](https://go.dev/dl/) 并安装。

```cmd
git clone https://github.com/zhlu32/range-for-jmeter.git
cd range-for-jmeter/websocket-service-demo
go mod tidy
go run main.go
```

如果出现端口冲突，可以在 main.go 中修改 port 的值

```go
const port string = ":6868"
```

启动后可访问[http://localhost:6868](http://localhost:6868)，应该返回

```json
{ "msg": "HelloWebSocket" }
```

通过[websocket 测试网站](http://www.easyswoole.com/wstool.html)测试 echo 复读机服务:

1. 在服务地址输入[ws://localhost:6868/echo](ws://localhost:6868/echo)
2. 点击连接
3. 调试信息显示 OPENED
4. 在输入框输入内容
5. 点击发送到服务端
6. 右侧消息记录应该可以看到服务器回复了一模一样的内容

通过[websocket 测试网站](http://www.easyswoole.com/wstool.html)测试 chat 聊天服务:

1. 在服务地址输入[ws://localhost:6868/chat](ws://localhost:6868/chat)
2. 点击连接
3. 新开一个页面重复 1、2 模拟多人使用
4. 调试信息显示 OPENED
5. 在输入框输入内容
6. 点击发送到服务端
7. 两个页面应该都可以互相看到对方发送的消息

## 说明

本项目使用 golang 作为开发语言，在 iris 这个后端框架上搭建 websocket 服务，所以并不能算是在 golang 基础上的极简示例，而是在 iris 这个框架基础上的极简示例。至于为什么要在 iris 的基础上使用，原因有两个，一是网络上已经有了使用原生 http 库配合一个 websocket 库实现的文章，而缺少在 iris 上搭建 websocket 的文章，二是个人使用 golang 搭建后端的时候都会基于 iris，因此在 iris 上搭建 websocket 更符合我的个人需求。

### 为什么需要 websocket

> 不讲细节，只给大家一个粗略的认知，如果希望了解更多的，可以自行百度，有许多讲的很详细的文章。

每一项技术的出现通常都是为了解决某些痛点、需求的，websocket 也不例外。

websocket 解决了什么痛点？基本上接触过 web 开发的人都知道，http api 有一个致命的痛点——服务端没有主动推送能力，一切的数据都是由客户端（前端）主动发起请求（ajax，http request），服务器被动响应返回的。

这样的模式导致了一个后果，一切的操作都是我们主动提出，服务器被动执行的，如查看天气信息，浏览视频，发布消息（朋友圈）等，但是我们熟悉的微信聊天，qq 聊天，我们不仅需要主动提出发送信息，还需要被动的接收消息。

之所以说我们是被动接收消息，是因为对方何时发送信息是永远都无法确定下来的，因此你无法做到每次对方发送消息的时候，你都能刚刚好向服务器发出获取对方消息的请求。

因此基于这个模式，对方的消息永远无法完美的被我们按时获取，因此要放弃完美接收的思路，使用另一种比较笨而暴力的方法——轮询，轮询是因为无法预测对方发送消息的时机，转而使用高频率的定时询问去弥补。

举个例子，就像你没有声音的情况下不知道什么时候厨房里的水壶会烧开，你就每隔 2 分钟去厨房看看一样，尽管可能无法完美的卡在水刚刚烧开的那个时刻，但效果也差不太多。

说回聊天，使用轮询的方式就是每隔半秒甚至更短的时间（时间取决于需求），就请求一次服务器，获取消息，尽管对方可能没有发送任何消息。这样可以一定程度上达到实时接收的效果，但付出了太多的代价，不间断的高频率请求会造成极大的资源消耗，而这一切的原因都来自于服务器没有主动推送的能力，如果服务器有主动推送能力，就可以在 b 的新消息发出的时候，主动推送给你，这样消息不仅是完美实时的，而且不需要轮询。

websocket 最大的特点，就是使得客户端和服务端之间建立起一个连接，这个连接会一直保持直到一方中断，而在连接上双方都可以互相发送、接收消息——全双工通信

### websocket 用途

前面已经对 websocket 的特点说得比较详细了，websocket 提供的全双工通信对于实时性要求高的场景非常适用，最核心的使用场景就是服务器有主动推送需求。比如游戏服务器、聊天服务器、以及前端程序员应该都使用过的 vscode 插件 live-server。

### 项目使用技术说明

语言方面使用 golang，框架选择了 iris，在 websocket 实现上，使用了 iris 的 websocket 子包，iris 的 websocket 子包需要配合 iris 作者的 neffos 包，而实际上，neffos 包中依赖 gorilla 和 gobwas 两位作者实现的 websocket，使用时可选使用其中的一个。

kataras.iris <- kataras.iris.websocket <- kataras.neffos <- gorilla.websocket/gobwas.websocket

本项目只使用了 neffos 中的较为基本的功能，为了兼容浏览器中的协议[MDN WebSocket API](https://developer.mozilla.org/zh-CN/docs/Web/API/WebSocket)
