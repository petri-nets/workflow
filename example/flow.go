package core

import (
	"context"

	"github.com/petri-nets/workflow/dao"
	"github.com/petri-nets/workflow/flow"
	"github.com/petri-nets/workflow/wfmod"
)

var WfSvc *FlowService

// FlowService 工作流主服务程序
type FlowService struct {
	ctx context.Context
}

// NewFlowService 创建一个工作流服务
func NewFlowService(ctx context.Context) (s *FlowService, err error) {
	if WfSvc != nil {
		return WfSvc, nil
	}

	WfSvc = &FlowService{
		ctx: ctx,
	}

	// 注册数据查询dao
	flow.RegFlowDao(getFlowDao())

	// you can register wfTransition some fields validator
	flow.RegCustomerValidator(func(transition *wfmod.WfTransition) error {
		return err
	})

	// register transition action executor, here you can executor self logic
	flow.RegExecutorHandle(func(cs *flow.Case, w *flow.Workitem) (bool, error) {

		wfContext := w.MergeCaseContext(cs)

		// merge system
		systemContext := w.GetSystemContext(cs)

		for k, v := range systemContext {
			wfContext[k] = v
		}

		// push async task

		// continue case
		go flow.ContinueCase(cs.WfCase.AppID, w.WfWorkitem.ID, wfContext)

		return true, err
	})

	return WfSvc, nil
}

func getFlowDao() wfmod.FlowDao {
	return &dao.GormDrive{}
}

func main() {
	wfCtx := wfmod.WfContextType{}
	flow.StartWorkflow(0, 0, wfCtx, "edenzou")
}

// Start 开启工作流服务
func (f *FlowService) Start() {
	// 监听触发事件
	// f.listener()
}

// // 监听事件程序
// func (f *FlowService) listener() {
// 	eventQueueChannel := config.GetWorkflowEventQueueKey()
// 	// 从队列获取数据
// 	for {
// 		select {
// 		case <-f.ctx.Done():
// 			log.WithFields(log.Fields{"err": f.ctx.Err()}).Warn("[workflowListener] ctx done")
// 			return
// 		default:
// 			if reply, err := f.QueueDao.BlPop(eventQueueChannel, 5); err == nil {
// 				if jobMessage, err := config.JobMessageDecode([]byte(reply)); err == nil {
// 					// 通过start job启动工作流新Case实例
// 					go f.startNewCase(jobMessage)

// 					// 中间job执行完成信息 继续之前Case
// 					go f.continueCase(jobMessage)
// 				}
// 			}
// 		}
// 	}
// }
