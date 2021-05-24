package common

import (
	"fmt"
	"github.com/simster7/notion-automation/client"
	log "github.com/sirupsen/logrus"
	"strings"
)

func GetDataBasePageName(page client.Page) string {
	if p, ok := GetDataBasePageProperty(page, "Name"); ok {
		return p.Title[0].PlainText
	}
	return ""
}

func GetDataBasePageProperty(page client.Page, property string) (client.DatabasePageProperty, bool) {
	p, ok := page.Properties.(client.DatabasePageProperties)[property]
	return p, ok
}

func LogAndError(logger *log.Entry, format string, args ...interface{}) error {
	// %w is only used in Errorf
	logger.Errorf(strings.ReplaceAll(format, "%w", "%s"), args)
	return fmt.Errorf(format, args)
}
