package tasks

import (
	"context"
	"fmt"
	"github.com/simster7/notion-automation/client"
	"github.com/simster7/notion-automation/common"
)

type DoOnToday struct{}

var _ Task = &DoOnToday{}

func GetDoOnToday() *DoOnToday {
	return &DoOnToday{}
}

func (d *DoOnToday) Do(ctx context.Context, notion *client.Client) error {
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
		return fmt.Errorf("failed to query database: %w", err)
	}

	return common.ExecutePages(res.Results, func(page client.Page) error {
		return setTaskDoOnToToday(ctx, notion, page)
	})
}

func setTaskDoOnToToday(ctx context.Context, notion *client.Client, task client.Page) error {
	_, err := notion.UpdatePageProps(ctx, task.ID, client.UpdatePageParams{
		DatabasePageProperties: &client.DatabasePageProperties{
			"Do On": client.DatabasePageProperty{
				Type: client.DBPropTypeDate,
				Date: &client.Date{
					Start: common.GetTime().NotionDate(),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update task '%s': %w", common.GetDataBasePageTitle(task), err)
	}
	return nil
}
