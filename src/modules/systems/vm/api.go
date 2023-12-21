package vm

type VmAPI interface {
	APISendMessage(id string, msg string)
	APIBroadcast(msg string)
}
