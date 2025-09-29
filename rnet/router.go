package rnet

import "Rinx/riface"

/*
BaseRouter: 实现接口IRouter的基类
具体路由实现根据需要对基类BaseRouter进行重写就可以了
有的路由请求不需要前置处理和后置处理，如果没有BaseRouter，那么就需要对IRouter的所有方法进行重写，有了BaseRouter就可以只重写Handler
*/
type BaseRouter struct{}

func (br *BaseRouter) PreHandler(req riface.IRequest) {}

func (br *BaseRouter) Handler(req riface.IRequest) {}

func (br *BaseRouter) PostHandler(req riface.IRequest) {}
