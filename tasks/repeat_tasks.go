package tasks

import (
	"context"
	"fmt"
	"github.com/simster7/notion-automation/client"
	"github.com/simster7/notion-automation/common"
	log "github.com/sirupsen/logrus"
	"time"
)

type RepeatTasks struct{}

var _ Task = &RepeatTasks{}

func GetRepeatTasks() *RepeatTasks {
	return &RepeatTasks{}
}

func (r *RepeatTasks) Do(ctx context.Context, notion *client.Client) error {
	res, err := notion.QueryDatabase(ctx, common.TaskDbId, &client.DatabaseQuery{
		Filter: client.DatabaseQueryFilter{
			And: []client.DatabaseQueryFilter{
				{
					Property: "Do On",
					Date: &client.DateDatabaseQueryFilter{
						Equals: common.GetTime().AddDate(0, 0, -1).NotionDate(),
					},
				},
				{
					Property: "Done",
					Checkbox: &client.CheckboxDatabaseQueryFilter{
						Equals: client.BoolPtr(true),
					},
				},
				{
					Property: "Repeat",
					Select: &client.SelectDatabaseQueryFilter{
						IsNotEmpty: true,
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to query database: %w", err)
	}

	return common.ExecutePages(res.Results, func(page client.Page) error {
		return createRepeatedTask(ctx, notion, page)
	})
}

func createRepeatedTask(ctx context.Context, notion *client.Client, task client.Page) error {
	var cadence string
	if repeat, ok := task.Properties.(client.DatabasePageProperties)["Repeat"]; !ok {
		log.Warn("task", common.GetDataBasePageTitle(task), "is not slated to be repeated. Skipping...")
		return nil
	} else {
		cadence = repeat.Select.Name
	}

	var anchorDate string
	if dueBy, ok := common.GetDataBasePageProperty(task, "Due By"); ok {
		anchorDate = dueBy.Date.Start
	} else if doOn, ok := common.GetDataBasePageProperty(task, "Do On"); ok {
		anchorDate = doOn.Date.Start
	} else {
		log.Warn("task", common.GetDataBasePageTitle(task), "does not have a date to calculate next instance from. Skipping...")
		return nil
	}

	var newDueBy string
	var err error
	switch cadence {
	case "Daily":
		newDueBy, err = nextInstance(anchorDate, 0, 0, 1)
	case "Weekly":
		newDueBy, err = nextInstance(anchorDate, 0, 0, 7)
	case "Monthly":
		newDueBy, err = nextInstance(anchorDate, 0, 1, 0)
	case "Yearly":
		newDueBy, err = nextInstance(anchorDate, 1, 0, 0)
	default:
		log.Warn("task", common.GetDataBasePageTitle(task), "cadence is unknown. Skipping...")
		return nil
	}
	if err != nil {
		return err
	}

	currentName, ok := common.GetDataBasePageProperty(task, "Name")
	if !ok {
		return fmt.Errorf("task does not have a name")
	}

	props := client.DatabasePageProperties{
		"Name": currentName,
		"Due By": client.DatabasePageProperty{
			Type: client.DBPropTypeDate,
			Date: &client.Date{
				Start: newDueBy,
			},
		},
	}

	_, err = notion.CreatePage(ctx, client.CreatePageParams{
		ParentID:               common.TaskDbId,
		DatabasePageProperties: &props,
	})

	return err
}

func nextInstance(anchorDate string, years, months, days int) (string, error) {
	anchorTime, err := time.Parse("2006-01-02", anchorDate)
	if err != nil {
		return "", fmt.Errorf("failed to parse anchor date: %w", err)
	}
	newDueBy := anchorTime.AddDate(years, months, days).Format("2006-01-02")
	return newDueBy, nil
}
