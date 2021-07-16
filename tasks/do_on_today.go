package tasks

import (
	"context"
	"github.com/simster7/notion-automation/client"
	"github.com/simster7/notion-automation/common"
	log "github.com/sirupsen/logrus"
)

type DoOnToday struct{}

var _ Task = &DoOnToday{}

func GetDoOnToday() *DoOnToday {
	return &DoOnToday{}
}

func (d *DoOnToday) GetName() string {
	return "DoOnToday"
}

func (d *DoOnToday) Do(ctx context.Context, notion *client.Client) error {
	logger := common.GetLogger().WithField("task", d.GetName())
	logger.Info("starting task")
	defer func() {
		logger.Info("finished task")
	}()

	logger.Info("querying task database for old due on tasks")
	res, err := notion.QueryDatabase(ctx, common.TaskDbId, &client.DatabaseQuery{
		Filter: client.DatabaseQueryFilter{
			And: []client.DatabaseQueryFilter{
				{
					Property: "Do On",
					Date: &client.DateDatabaseQueryFilter{
						OnOrBefore: common.GetTime().AddDate(0, 0, -1).NotionDate(),
					},
				},
				{
					Property: "Done",
					Checkbox: &client.CheckboxDatabaseQueryFilter{
						Equals: client.BoolPtr(false),
					},
				},
			},
		},
	})
	if err != nil {
		return common.LogAndError(logger, "failed to query database: %w", err)
	}
	logger.Info("database query successful")

	return common.ExecutePages(res.Results, func(page client.Page, _ int) error {
		return setTaskDoOnToToday(ctx, notion, logger, page)
	})
}

func setTaskDoOnToToday(ctx context.Context, notion *client.Client, logger *log.Entry, task client.Page) error {
	dueOn := common.GetTime().NotionDate()
	taskLogger := logger.WithField("notion_task", common.GetDataBasePageName(task))
	taskLogger.Infof("updating task due on date to '%s'", dueOn)

	_, err := notion.UpdatePageProps(ctx, task.ID, client.UpdatePageParams{
		DatabasePageProperties: &client.DatabasePageProperties{
			"Do On": client.DatabasePageProperty{
				Type: client.DBPropTypeDate,
				Date: &client.Date{
					Start: dueOn,
				},
			},
			"Must": client.DatabasePageProperty{
				Type:     client.DBPropTypeCheckbox,
				Checkbox: client.BoolPtr(false),
			},
		},
	})
	if err != nil {
		return common.LogAndError(taskLogger, "failed to update task: %s", err)
	}
	taskLogger.Info("update query successful")
	return nil
}
