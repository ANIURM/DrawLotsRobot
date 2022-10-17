package controller

func Init()

func initTimer()

func initEvent() {
	// register your handlers here
	// example
	dispatcher.RegisterListener(receiveMessage.Receive, "im.message.receive_v1")
}
