/**
 * @File: print.go
 * @Author: hsien
 * @Description:
 * @Date: 9/28/21 6:05 PM
 */

package cronjob

import "custom_server/pkg/log"

type Print struct {
	value string
}

func NewPrint() *Print {
	return &Print{}
}

func (p *Print) Run() {
	log.Debug("cron job demo: print")
}

func (p *Print) Spec() string {
	return "@every 10s"
}

func (p *Print) Name() string {
	return "cron print demo"
}
