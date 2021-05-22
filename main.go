package main

import (
	"context"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/simster7/notion-automation/client"
	"github.com/simster7/notion-automation/common"
	"github.com/simster7/notion-automation/tasks"
)

var notion = client.NewClient(os.Getenv("NOTION_TOKEN"))

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