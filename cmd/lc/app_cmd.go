package main

import (
	"github.com/lab5e/l5log/pkg/lg"
	"github.com/lab5e/lospan/pkg/pb/lospan"
)

type appCmd struct {
	List listAppCmd   `kong:"cmd,help='List applications',aliases='ls'"`
	Add  addAppCmd    `kong:"cmd,help='Add application',aliases='create,new'"`
	Del  deleteAppCmd `kong:"cmd,help='Remove application',aliases='rm,delete'"`
}

func (c *appCmd) Run(args *params) error {
	return nil
}

type listAppCmd struct {
}

func (c *listAppCmd) Run(args *params) error {

	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()

	lg.Info("Listing applications....")
	res, err := client.ListApplications(ctx, &lospan.ListApplicationsRequest{})
	if err != nil {
		return err
	}
	lg.Info("%d applications found", len(res.Applications))
	for _, app := range res.Applications {
		lg.Info("App: %+v", app)
	}
	return nil
}

type addAppCmd struct {
}

func (c *addAppCmd) Run(args *params) error {
	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()

	lg.Info("Creating new application...")
	res, err := client.CreateApplication(ctx, &lospan.CreateApplicationRequest{})
	if err != nil {
		return err
	}

	lg.Info("Created application with EUI %s", res.Eui)
	return nil
}

type deleteAppCmd struct {
	EUI string `kong:"help='Application EUI to delete'"`
}

func (deleteAppCmd) Run(args *params) error {
	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()

	app, err := client.DeleteApplication(ctx, &lospan.DeleteApplicationRequest{Eui: args.App.Del.EUI})
	if err != nil {
		return err
	}
	lg.Info("Removed application with EUI %s", app.Eui)
	return nil
}
