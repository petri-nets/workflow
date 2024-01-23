package main

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/petri-nets/workflow/dao"
	"github.com/petri-nets/workflow/flow"
	toolstester "github.com/petri-nets/workflow/test/tools"
)

// import seq_test "github.com/petri-nets/workflow/test/seq"

func main() {

	toolstester.InitConfig("/test/conf/config.yaml")

	db := toolstester.GetDBConnection(&toolstester.Conf.Database)

	flowDao := &dao.GormDrive{DB: db}

	// 注册数据查询dao
	flow.RegFlowDao(flowDao)

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
		flow.ContinueCase(cs.WfCase.AppID, w.WfWorkitem.JobID, wfContext)

		return true, nil
	})

	// 清空数据
	toolstester.TruncateData(db)

	// 获取用例配置
	toolstester.FindTestYamlFiles()

	time.Sleep(3 * time.Second)

}
