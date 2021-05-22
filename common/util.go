package common

import "github.com/simster7/notion-automation/client"

func GetDataBasePageTitle(page client.Page) string {
	return page.Properties.(client.DatabasePageProperties)["Name"].Title[0].PlainText
}
