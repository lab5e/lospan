package main

import (
	"fmt"

	"github.com/lab5e/lospan/pkg/pb/lospan"
)

type appCmd struct {
	List listAppCmd   `kong:"cmd,help='List applications',aliases='ls'"`
	Add  addAppCmd    `kong:"cmd,help='Add application',aliases='create,new,n'"`
	Del  deleteAppCmd `kong:"cmd,help='Remove application',aliases='rm,delete,d'"`
	Get  getAppCmd    `kong:"cmd,help='Retrieve application',aliases='show,retrieve,g'"`
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

	res, err := client.ListApplications(ctx, &lospan.ListApplicationsRequest{})
	if err != nil {
		return err
	}
	fmt.Printf("%d applications found\n", len(res.Applications))
	fmt.Println()
	fmt.Println("EUI")
	for _, app := range res.Applications {
		fmt.Println(app.Eui)
	}
	return nil
}

type addAppCmd struct {
}

func (*addAppCmd) Run(args *params) error {
	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()

	res, err := client.CreateApplication(ctx, &lospan.CreateApplicationRequest{})
	if err != nil {
		return err
	}

	fmt.Printf("Created application with EUI %s\n", res.Eui)
	return nil
}

type deleteAppCmd struct {
	EUI string `kong:"help='Application EUI to delete'"`
}

func (*deleteAppCmd) Run(args *params) error {
	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()

	app, err := client.DeleteApplication(ctx, &lospan.DeleteApplicationRequest{Eui: args.App.Del.EUI})
	if err != nil {
		return err
	}
	fmt.Printf("Removed application with EUI %s\n", app.Eui)
	return nil
}

type getAppCmd struct {
	EUI string `kong:"help='Application EUI to retrieve'"`
}

func (*getAppCmd) Run(args *params) error {
	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()

	app, err := client.GetApplication(ctx, &lospan.GetApplicationRequest{Eui: args.App.Get.EUI})
	if err != nil {
		return err
	}
	fmt.Printf("Application EUI %s\n", app.Eui)
	return nil
}
