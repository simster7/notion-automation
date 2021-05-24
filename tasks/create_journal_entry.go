package tasks

import (
	"context"
	"github.com/simster7/notion-automation/client"
	"github.com/simster7/notion-automation/common"
	log "github.com/sirupsen/logrus"
)

type CreateJournalEntry struct{}

var _ Task = &CreateJournalEntry{}

func GetCreateJournalEntry() *CreateJournalEntry {
	return &CreateJournalEntry{}
}

func (c *CreateJournalEntry) GetName() string {
	return "CreateJournalEntry"
}

func (c *CreateJournalEntry) Do(ctx context.Context, notion *client.Client) error {
	logger := common.GetLogger().WithField("task", c.GetName())
	logger.Info("starting task")
	defer func() {
		logger.Info("finished task")
	}()

	exists, err := checkIfEntryExists(ctx, logger, notion)
	if err != nil {
		return err
	}

	if exists {
		logger.Warn("journal entry already exists, skipping...")
		return nil
	}

	date := common.GetTime().NotionDate()
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
				Start: date,
			},
		},
	}

	logger.Info("creating journal entry for '%s'", date)
	_, err = notion.CreatePage(ctx, client.CreatePageParams{
		ParentID:               common.JournalDbId,
		DatabasePageProperties: &props,
	})

	if err != nil {
		return common.LogAndError(logger, "failed to create journal entry: %w", err)
	}

	logger.Info("journal entry created")

	return nil
}

func checkIfEntryExists(ctx context.Context, logger *log.Entry, notion *client.Client) (bool, error) {
	date := common.GetTime().NotionDate()
	logger.Infof("making journal database query to see if entry exists for '%s'", date)
	res, err := notion.QueryDatabase(ctx, common.JournalDbId, &client.DatabaseQuery{
		Filter: client.DatabaseQueryFilter{
			Property: "Date",
			Date: &client.DateDatabaseQueryFilter{
				Equals: date,
			},
		},
	})

	if err != nil {
		return false, common.LogAndError(logger, "failed to validate journal entry did not already exist: %w", err)
	}

	return len(res.Results) > 0, nil
}
