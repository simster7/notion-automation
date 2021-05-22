package main

import (
	"context"
	"github.com/simster7/notion-automation/client"
	"github.com/simster7/notion-automation/common"
	"github.com/simster7/notion-automation/tasks"
	log "github.com/sirupsen/logrus"
	"sync"
)

var notion = client.NewClient("")

func main() {
	common.SetTime()
	ctx := context.Background()

	taskQueue := []tasks.Task{tasks.GetCreateJournalEntry(), tasks.GetDoOnToday()}
	var wg sync.WaitGroup
	wg.Add(len(taskQueue))

	for _, task := range taskQueue {
		task := task
		go func() {
			defer wg.Done()
			err := task.Do(ctx, notion)
			if err != nil {
				log.Errorf("cannot complete task: %s", err)
			}
		}()
	}

	wg.Wait()
}