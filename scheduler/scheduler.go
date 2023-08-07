package scheduler

import (
	"context"
	"time"

	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/mgmtmodel"
	"github.com/go-co-op/gocron"
)

func updateRegisteredHostInfo() {
	ctx := context.WithValue(context.Background(), common.TraceIDKey("TraceID"), common.GenerateTraceID())
	hostListModel := mgmtmodel.HostList{}
	hostListModel.Update(ctx)
}

func StartScheduler() {
	// 创建一个新的计划任务
	s := gocron.NewScheduler(time.UTC)

	// 每隔一段时间执行异步任务

	// 将异步任务添加到计划中，每隔2秒执行一次
	s.Every(60).Seconds().Do(updateRegisteredHostInfo)

	// 开始计划任务的调度
	s.StartAsync()

}
