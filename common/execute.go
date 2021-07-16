package common

import (
	"github.com/simster7/notion-automation/client"
	log "github.com/sirupsen/logrus"
	"sync"
)

type ExecutePageFunc func(client.Page, int) error

func ExecutePages(pages []client.Page, fn ExecutePageFunc) error {
	var wg sync.WaitGroup
	wg.Add(len(pages))

	for i, page := range pages {
		page := page
		i := i
		go func() {
			defer wg.Done()
			err := fn(page, i)
			if err != nil {
				log.Errorf("error executing on page: %s", err)
			}
		}()
	}

	wg.Wait()
	return nil
}
