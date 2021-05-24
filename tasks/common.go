package tasks

import (
	"context"
	"github.com/simster7/notion-automation/client"
)

type Task interface {
	Do(context.Context, *client.Client) error
	GetName() string
}
