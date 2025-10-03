package rnet

import (
	"Rinx/riface"
	"Rinx/utils"
	"fmt"
	"strconv"
)

// 路由策略实现
type MsgHandler struct {
	// msgId和处理方法的对应关系
	Apis map[uint32]riface.IRouter
	// 协程池的大小
	WorkerPoolSize uint32
	// 任务队列,个数暂定和协程池数量一致
	TaskQueue []chan riface.IRequest
}

// 初始化MsgHandler
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]riface.IRouter),
		WorkerPoolSize: utils.Config.WorkerPoolSize,
		TaskQueue:      make([]chan riface.IRequest, utils.Config.WorkerPoolSize),
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

// 启动协程池
func (mh *MsgHandler) StartWorkerPool() {
	fmt.Println("[协程池开启...]")
	for i := uint32(0); i < mh.WorkerPoolSize; i++ {
		// 开辟任务队列的空间容量
		mh.TaskQueue[i] = make(chan riface.IRequest, utils.Config.TaskQueueCapacity)
		// 启动任务
		go mh.StartOneWorker(i, mh.TaskQueue[i])

	}
}

// 取出任务队列中的request进行业务处理
func (mh *MsgHandler) StartOneWorker(workerId uint32, request chan riface.IRequest) {
	// TODO 退出机制
	fmt.Println("[Worker Id = ", workerId, " is started...]")
	// 等待任务队列的request到来
	for {
		select {
		// 有消息到来，取出消息进行处理
		case req := <-request:
			mh.DoMsgHandler(req)
		}
	}
}

// 将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandler) SendReqToTaskQueue(request riface.IRequest) {
	// 根据connId来决定request放入哪个任务队列中
	// connId对线程池大小进行模运算确定
	workerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " request msgID=", request.GetMessage().GetMsgId(), "to workerID=", workerId)
	// 消息放入任务队列
	mh.TaskQueue[workerId] <- request
}
