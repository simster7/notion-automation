package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/simster7/notion-automation/client"
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

func MD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
