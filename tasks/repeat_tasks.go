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

func (r *RepeatTasks) GetName() string {
	return "RepeatTasks"
}

func (r *RepeatTasks) Do(ctx context.Context, notion *client.Client) error {
	logger := common.GetLogger().WithField("task", r.GetName())
	logger.Info("starting task")
	defer func() {
		logger.Info("finished task")
	}()

	logger.Info("querying task database for repeated tasks")
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
		return common.LogAndError(logger, "failed to query database: %w", err)
	}
	logger.Info("database query successful")

	return common.ExecutePages(res.Results, func(page client.Page) error {
		return createRepeatedTask(ctx, notion, logger, page)
	})
}

func createRepeatedTask(ctx context.Context, notion *client.Client, logger *log.Entry, task client.Page) error {
	taskLogger := logger.WithField("notion_task", common.GetDataBasePageName(task))
	taskLogger.Infof("creating repeated task")

	var cadence string
	if repeat, ok := task.Properties.(client.DatabasePageProperties)["Repeat"]; !ok {
		taskLogger.Warn("task is not slated to be repeated, skipping...")
		return nil
	} else {
		cadence = repeat.Select.Name
	}

	var anchorDate string
	if dueBy, ok := common.GetDataBasePageProperty(task, "Due By"); ok {
		anchorDate = dueBy.Date.Start
		taskLogger.Infof("basing anchor date on due by '%s'", anchorDate)
	} else if doOn, ok := common.GetDataBasePageProperty(task, "Do On"); ok {
		anchorDate = doOn.Date.Start
		taskLogger.Infof("basing anchor date on do on '%s'", anchorDate)
	} else {
		taskLogger.Warn("task does not have a date to calculate next instance from, skipping...")
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
		taskLogger.Warn("task cadence is unknown, skipping...")
		return nil
	}
	if err != nil {
		return err
	}
	taskLogger.Infof("creating new task for due by '%s' based on cadence '%s' and anchor date '%s'", newDueBy, cadence, anchorDate)

	currentName, ok := common.GetDataBasePageProperty(task, "Name")
	if !ok {
		return fmt.Errorf("task does not have a name")
	}

	exists, err := checkIfTaskExists(ctx, notion, logger, common.GetDataBasePageName(task), newDueBy)
	if err != nil {
		return err
	}

	if exists {
		taskLogger.Warn("task already exists, skipping...")
		return nil
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

func checkIfTaskExists(ctx context.Context, notion *client.Client, logger *log.Entry, name, newDueBy string) (bool, error) {
	logger.Infof("making journal database query to see if task already exists with name '%s' and due by '%s'", name, newDueBy)
	res, err := notion.QueryDatabase(ctx, common.TaskDbId, &client.DatabaseQuery{
		Filter: client.DatabaseQueryFilter{
			And: []client.DatabaseQueryFilter{
				{
					Property: "Name",
					Text: &client.TextDatabaseQueryFilter{
						Equals: name,
					},
				},
				{
					Property: "Due By",
					Date: &client.DateDatabaseQueryFilter{
						Equals: newDueBy,
					},
				},
			},
		},
	})

	if err != nil {
		return false, common.LogAndError(logger, "failed to validate task did not already exist: %w", err)
	}

	return len(res.Results) > 0, nil
}
