package common

import (
	"github.com/simster7/notion-automation/client"
	log "github.com/sirupsen/logrus"
	"sync"
)

type ExecutePageFunc func(client.Page) error

func ExecutePages(pages []client.Page, fn ExecutePageFunc) error {
	var wg sync.WaitGroup
	wg.Add(len(pages))

	for _, page := range pages {
		page := page
		go func() {
			defer wg.Done()
			err := fn(page)
			if err != nil {
				log.Errorf("error executing on page: %s", err)
			}
		}()
	}

	wg.Wait()
	return nil
}
