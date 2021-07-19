package tasks

import (
	"context"
	"os"
	"testing"

	"github.com/simster7/notion-automation/client"
	"github.com/simster7/notion-automation/common"
)

func TestAddCalendarEvents_Do(t *testing.T) {
	common.SetTime()
	common.InitLogger("test")
	a := GetAddCalendarEvents()
	a.Do(context.Background(), client.NewClient(os.Getenv("NOTION_TOKEN")))
}
