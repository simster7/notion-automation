package tasks

import (
	"context"
	"fmt"

	//"fmt"
	"github.com/simster7/notion-automation/client"
	"github.com/simster7/notion-automation/common"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/calendar/v3"
	"os"
)

var calendarId = os.Getenv("CALENDAR_ID")

type AddCalendarEvents struct{}

var _ Task = &AddCalendarEvents{}

func GetAddCalendarEvents() *AddCalendarEvents {
	return &AddCalendarEvents{}
}

func (a AddCalendarEvents) GetName() string {
	return "AddCalendarEvents"
}

func (a AddCalendarEvents) Do(ctx context.Context, notion *client.Client) error {
	logger := common.GetLogger().WithField("task", a.GetName())
	logger.Info("starting task")
	defer func() {
		logger.Info("finished task")
	}()

	cal, err := calendar.NewService(ctx)
	if err != nil {
		return common.LogAndError(logger, "Unable to retrieve calendar client: %v", err)
	}

	// Don't forget to accept calendar (only necessary once)
	//
	//		srv.CalendarList.Insert(&calendar.CalendarListEntry{Id: "[Can be found in Calendar Sharing Settings]"}).Do()
	//

	logger.Info("querying task database for today's must tasks")
	res, err := notion.QueryDatabase(ctx, common.TaskDbId, &client.DatabaseQuery{
		Filter: client.DatabaseQueryFilter{
			And: []client.DatabaseQueryFilter{
				{
					Property: "Do On",
					Date: &client.DateDatabaseQueryFilter{
						OnOrBefore: common.GetTime().NotionDate(),
					},
				},
				{
					Property: "Done",
					Checkbox: &client.CheckboxDatabaseQueryFilter{
						Equals: client.BoolPtr(false),
					},
				},
				{
					Property: "Must",
					Checkbox: &client.CheckboxDatabaseQueryFilter{
						Equals: client.BoolPtr(true),
					},
				},
			},
		},
	})
	if err != nil {
		return common.LogAndError(logger, "failed to query database: %w", err)
	}
	logger.Info("database query successful")

	return common.ExecutePages(res.Results, func(page client.Page, i int) error {
		return createMustEvent(cal, logger, page, i)
	})
}

func createMustEvent(cal *calendar.Service, logger *log.Entry, task client.Page, index int) error {
	taskLogger := logger.WithField("notion_task", common.GetDataBasePageName(task))
	taskLogger.Infof("creating calendar event for must task")

	start, end := common.GetTime().GetCalendarEventTimes(index)
	fmt.Println(index, start, end)
	taskTitle := task.Properties.(client.DatabasePageProperties)["Name"].Title[0].PlainText
	_, err := cal.Events.Insert(calendarId, &calendar.Event{
		Summary: fmt.Sprintf("[MUST] %s", taskTitle),
		Start:   &calendar.EventDateTime{DateTime: start, TimeZone: common.TimeZone},
		End:     &calendar.EventDateTime{DateTime: end, TimeZone: common.TimeZone},
	}).Do()
	if err != nil {
		return common.LogAndError(taskLogger, "failed to create must event: %s", err)
	}
	taskLogger.Info("must event creation successful")
	return nil
}
