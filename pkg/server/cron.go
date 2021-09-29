/**
 * @File: cron.go
 * @Author: hsien
 * @Description:
 * @Date: 9/18/21 5:33 PM
 */

package server

import (
	"errors"
	"github.com/robfig/cron/v3"
	"sync"
)

type CronJob interface {
	Run()
	Spec() string
	Name() string
}

type CronSrv struct {
	cron *cron.Cron
	task sync.Map
}

type Entry struct {
	Name string
	cron.Entry
}

type job struct {
	name string
	spec string
	cmd  func()
}

func (j job) Run() { j.cmd() }

func (j job) Name() string { return j.name }

func (j job) Spec() string { return j.spec }

func NewCron() *CronSrv {
	return &CronSrv{
		cron: cron.New(cron.WithSeconds()),
		task: sync.Map{},
	}
}

func (c *CronSrv) AddFunc(name, spec string, cmd func()) error {
	return c.AddJob(&job{
		name: name,
		spec: spec,
		cmd:  cmd,
	})
}

func (c *CronSrv) AddJob(job CronJob) error {
	if c.cron == nil {
		return errors.New("cron not initialize")
	}

	if job.Name() == "" {
		return errors.New("job no name")
	}

	if _, ok := c.task.Load(job.Name()); ok {
		return errors.New("exists the same name job")
	}

	entryId, err := c.cron.AddJob(job.Spec(), job)
	if err != nil {
		return err
	}
	c.task.Store(job.Name(), int(entryId))
	return nil
}

func (c *CronSrv) RemoveById(entryID int) {
	c.cron.Remove(cron.EntryID(entryID))
}

func (c *CronSrv) RemoveByName(entryName string) {
	if entryId, ok := c.task.LoadAndDelete(entryName); ok {
		if ok {
			c.RemoveById(entryId.(int))
		}
	}
}

func (c *CronSrv) Start() {
	c.cron.Start()
}

func (c *CronSrv) Stop() {
	c.cron.Stop()
}

func (c *CronSrv) Entries() []Entry {
	cronEntries := c.cron.Entries()
	var entries []Entry
	for _, cronEntry := range cronEntries {
		entry := Entry{
			"",
			cronEntry,
		}

		c.task.Range(func(k, v interface{}) bool {
			if v.(int) == int(cronEntry.ID) {
				entry.Name = k.(string)
				return false
			}
			return true
		})
		entries = append(entries, entry)
	}
	return entries
}

func (c *CronSrv) Count() int {
	var count int
	c.task.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	return count
}
