package rnet

import (
	"Rinx/riface"
	"fmt"
	"strconv"
)

// 路由策略实现
type MsgHandler struct {
	// msgId和处理方法的对应关系
	Apis map[uint32]riface.IRouter
}

// 初始化MsgHandler
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		make(map[uint32]riface.IRouter),
	}
}
func (mh *MsgHandler) DoMsgHandler(req riface.IRequest) {
	id := req.GetMessage().GetMsgId()
	handle, ok := mh.Apis[id]
	if !ok {
		fmt.Println("执行失败...")
		return
	}

	handle.PreHandler(req)
	handle.Handler(req)
	handle.PostHandler(req)
	fmt.Println("执行成功...")

}

func (mh *MsgHandler) AddRouter(msgId uint32, router riface.IRouter) {
	if _, ok := mh.Apis[msgId]; ok {
		fmt.Println("msgId: ", msgId, "对应的路由策略已经存在，注册失败...")
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	mh.Apis[msgId] = router
	fmt.Println("msgId", msgId, "注册成功...")

}
