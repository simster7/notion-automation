package tasks

import (
	"context"
	"fmt"
	"github.com/simster7/notion-automation/client"
	"github.com/simster7/notion-automation/common"
	log "github.com/sirupsen/logrus"
)

type CreateJournalEntry struct{}

var _ Task = &CreateJournalEntry{}

func GetCreateJournalEntry() *CreateJournalEntry {
	return &CreateJournalEntry{}
}

func (c *CreateJournalEntry) Do(ctx context.Context, notion *client.Client) error {
	exists, err := checkIfEntryExists(ctx, notion)
	if err != nil {
		return err
	}

	if exists {
		log.Warn("journal entry already exists, skipping...")
		return nil
	}

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

	_, err = notion.CreatePage(ctx, client.CreatePageParams{
		ParentID:               common.JournalDbId,
		DatabasePageProperties: &props,
	})

	if err != nil {
		return fmt.Errorf("failed to create journal entry: %w", err)
	}

	return nil
}

func checkIfEntryExists(ctx context.Context, notion *client.Client) (bool, error) {
	res, err := notion.QueryDatabase(ctx, common.JournalDbId, &client.DatabaseQuery{
		Filter: client.DatabaseQueryFilter{
			Property: "Date",
			Date: &client.DateDatabaseQueryFilter{
				Equals: common.GetTime().NotionDate(),
			},
		},
	})

	if err != nil {
		return false, fmt.Errorf("failed to validate journal entry did not alread exist: %w", err)
	}

	return len(res.Results) > 0, nil
}
