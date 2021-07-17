package jobs

import (
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/simster7/notion-automation/client"
	"github.com/simster7/notion-automation/common"
	"github.com/simster7/notion-automation/tasks"
)

var notion = client.NewClient(os.Getenv("NOTION_TOKEN"))

func Nightly(_ http.ResponseWriter, _ *http.Request) {
	common.SetTime()
	common.InitLogger("nightly")
	log := common.GetLogger()
	log.Infof("starting nightly job...")
	ctx := context.Background()

	taskQueue := []tasks.Task{tasks.GetCreateJournalEntry(), tasks.GetDoOnToday(), tasks.GetRepeatTasks()}
	var wg sync.WaitGroup
	wg.Add(len(taskQueue))

	for _, task := range taskQueue {
		task := task
		log.Infof("starting task '%s'", task.GetName())
		go func() {
			defer wg.Done()
			err := task.Do(ctx, notion)
			if err != nil {
				log.Errorf("error completing task '%s': %s", task.GetName(), err)
			}
		}()
	}

	wg.Wait()

	err := tasks.GetAddCalendarEvents().Do(ctx, notion)
	if err != nil {
		log.Errorf("error completing task 'AddCalendarEvents': %s", err)
	}
}
