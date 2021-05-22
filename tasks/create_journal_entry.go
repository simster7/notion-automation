package tasks

import (
	"context"
	"fmt"
	"github.com/simster7/notion-automation/client"
	"github.com/simster7/notion-automation/common"
)

type CreateJournalEntry struct{}

var _ Task = &CreateJournalEntry{}

func GetCreateJournalEntry() *CreateJournalEntry {
	return &CreateJournalEntry{}
}

func (c *CreateJournalEntry) Do(ctx context.Context, notion *client.Client) error {
	props := client.DatabasePageProperties{
		"Name": client.DatabasePageProperty{
			Type: client.DBPropTypeTitle,
			Title: []client.RichText{
				{
					Text: &client.Text{
						Content: common.GetTime().Format("2 Jan 2006"),
					},
				},
			},
		},
		"Date": client.DatabasePageProperty{
			Type: client.DBPropTypeDate,
			Date: &client.Date{
				Start: common.GetTime().NotionDate(),
			},
		},
	}

	_, err := notion.CreatePage(ctx, client.CreatePageParams{
		ParentID:   common.JournalDbId,
		DatabasePageProperties: &props,
	})

	if err != nil {
		return fmt.Errorf("failed to create journal entry: %w", err)
	}

	return nil
}
