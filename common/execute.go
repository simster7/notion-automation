package common

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/simster7/notion-automation/client"
)

type ExecutePageFunc func(client.Page, int) error

func ExecutePages(pages []client.Page, fn ExecutePageFunc) error {
	var wg sync.WaitGroup
	wg.Add(len(pages))

	for i, page := range pages {
		i := i
		page := page
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
