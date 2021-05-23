package common

import "github.com/simster7/notion-automation/client"

func GetDataBasePageTitle(page client.Page) string {
	if p, ok := GetDataBasePageProperty(page, "title"); ok {
		return p.Title[0].PlainText
	}
	return ""
}

func GetDataBasePageProperty(page client.Page, property string) (client.DatabasePageProperty, bool) {
	p, ok := page.Properties.(client.DatabasePageProperties)[property]
	return p, ok
}
