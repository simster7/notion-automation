package tasks

import (
	"context"
	log "github.com/sirupsen/logrus"

	"github.com/simster7/notion-automation/client"
	"github.com/simster7/notion-automation/common"
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

	loc, err := getYesterdaysLocation(ctx, notion, logger)
	if err != nil {
		return err
	}

	date := common.GetTime().NotionDate()
	props := client.DatabasePageProperties{
		"Name": client.DatabasePageProperty{
			Type: client.DBPropTypeTitle,
			Title: []client.RichText{
				{
					Text: &client.Text{
						Content: common.GetTime().Format("2 January 2006"),
					},
					Type: client.RichTextTypeText,
				},
			},
		},
		"Date": client.DatabasePageProperty{
			Type: client.DBPropTypeDate,
			Date: &client.Date{
				Start: date,
			},
		},
		"Location": client.DatabasePageProperty{
			Type: client.DBPropTypeSelect,
			Select: &client.SelectOptions{
				Name: loc,
			},
		},
	}

	logger.Infof("creating journal entry for '%s'", date)
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

func getYesterdaysLocation(ctx context.Context, notion *client.Client, logger *log.Entry) (string, error) {
	yesterday := common.GetTime().AddDate(0, 0, -1).NotionDate()
	logger.Infof("making journal database query to get yesterday's ('%s') location", yesterday)
	res, err := notion.QueryDatabase(ctx, common.JournalDbId, &client.DatabaseQuery{
		Filter: client.DatabaseQueryFilter{
			Property: "Date",
			Date: &client.DateDatabaseQueryFilter{
				Equals: yesterday,
			},
		},
	})

	if err != nil {
		return "", common.LogAndError(logger, "failed to get yesterday's location: %w", err)
	}

	if len(res.Results) != 1 {
		return "", common.LogAndError(logger, "invalid state: '%s' did not have a journal entry or more than one existed", yesterday)
	}

	loc, ok := common.GetDataBasePageProperty(res.Results[0], "Location")
	if !ok {
		return "", common.LogAndError(logger, "invalid state: '%s' does not have location", yesterday)
	}

	return loc.Select.Name, nil
}
