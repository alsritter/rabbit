package server

import "alsritter.icu/rabbit-template/internal/pkg/cron"

// NewCronServer 用来注册一些定时任务
func NewCronServer() *cron.Server {
	return cron.NewServer()
}
