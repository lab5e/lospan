package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/lab5e/lospan/pkg/pb/lospan"
)

func printGateway(gw *lospan.Gateway) {
	fmt.Printf("    EUI:             %s\n", gw.Eui)
	fmt.Printf("    IP:              %s\n", gw.GetIp())
	fmt.Printf("    Strict IP check: %t\n", gw.GetStrictIp())
	fmt.Printf("    Latitude:        %2.2f\n", gw.GetLatitude())
	fmt.Printf("    Longitude:       %3.2f\n", gw.GetLongitude())
	fmt.Printf("    Altitude:        %2.2f\n", gw.GetAltitude())
}

type gwCmds struct {
	Add    gwAddCmd    `kong:"cmd,help='Add gateway',aliases='create,a'"`
	Del    gwDelCmd    `kong:"cmd,help='Delete gateway',aliases='rm,delete,d'"`
	Update gwUpdateCmd `kong:"cmd,help='Update gateway',aliases='up'"`
	Get    gwGetCmd    `kong:"cmd,help='Get gateway info',aliases='g'"`
	List   gwListCmd   `kong:"cmd,help='List gateways',aliases='ls'"`
}

type gwAddCmd struct {
	EUI       string  `kong:"help='Gateway EUI',required"`
	IP        string  `kong:"help='Gateway IP address',required"`
	Altitude  float32 `kong:"help='Altitude for gateway'"`
	Longitude float32 `kong:"help='Longitude for gateway (-360...360)'"`
	Latitude  float32 `kong:"help='Latitude for gateway (-90...90)'"`
	StrictIP  bool    `kong:"help='Strict IP check',default=true"`
}

func (*gwAddCmd) Run(args *params) error {
	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()

	gw, err := client.CreateGateway(ctx, &lospan.Gateway{
		Eui:       args.GW.Add.EUI,
		Ip:        newPtr(args.GW.Add.IP),
		StrictIp:  newPtr(args.GW.Add.StrictIP),
		Altitude:  newPtr(args.GW.Add.Altitude),
		Longitude: newPtr(args.GW.Add.Longitude),
		Latitude:  newPtr(args.GW.Add.Latitude),
	})
	if err != nil {
		return err
	}
	fmt.Println("Added gateway")
	printGateway(gw)
	return nil
}

type gwUpdateCmd struct {
	EUI       string  `kong:"help='Gateway EUI',required"`
	IP        string  `kong:"help='Gateway IP address',optional"`
	Altitude  float32 `kong:"help='Altitude for gateway',default=-99999"`
	Longitude float32 `kong:"help='Longitude for gateway (-360...360)',default=-999"`
	Latitude  float32 `kong:"help='Latitude for gateway (-90...90)',default=-999"`
	StrictIP  bool    `kong:"help='Strict IP check',default=true,optional"`
}

func (*gwUpdateCmd) Run(args *params) error {
	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()

	gw := &lospan.Gateway{
		Eui: args.GW.Update.EUI,
	}
	if args.GW.Update.IP != "" {
		gw.Ip = newPtr(args.GW.Update.IP)
	}
	if args.GW.Update.Altitude != -99999 {
		gw.Altitude = newPtr(args.GW.Update.Altitude)
	}
	if args.GW.Update.Longitude != -999 {
		gw.Longitude = newPtr(args.GW.Update.Longitude)
	}
	if args.GW.Update.Latitude != -999 {
		gw.Latitude = newPtr(args.GW.Update.Latitude)
	}
	gw.StrictIp = newPtr(args.GW.Update.StrictIP)

	newGW, err := client.UpdateGateway(ctx, gw)
	if err != nil {
		return err
	}
	fmt.Println("Updated gateway")
	printGateway(newGW)
	return nil
}

type gwDelCmd struct {
	EUI string `kong:"help='Gateway EUI',required"`
}

func (*gwDelCmd) Run(args *params) error {
	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()

	gw, err := client.DeleteGateway(ctx, &lospan.DeleteGatewayRequest{
		Eui: args.GW.Del.EUI,
	})
	if err != nil {
		return err
	}
	fmt.Println("Deleted gateway")
	printGateway(gw)
	return nil
}

type gwGetCmd struct {
	EUI string `kong:"help='Gateway EUI',required"`
}

func (*gwGetCmd) Run(args *params) error {
	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()

	gw, err := client.GetGateway(ctx, &lospan.GetGatewayRequest{
		Eui: args.GW.Get.EUI,
	})
	if err != nil {
		return err
	}
	fmt.Println("Found gateway")
	printGateway(gw)
	return nil
}

type gwListCmd struct {
}

func (*gwListCmd) Run(args *params) error {
	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()

	gws, err := client.ListGateways(ctx, &lospan.ListGatewaysRequest{})
	if err != nil {
		return err
	}

	writer := tabwriter.NewWriter(os.Stdout, 3, 4, 2, ' ', 0)
	writer.Write([]byte("EUI\tIP\tStrict\tLat\tLon\tAlt\n"))
	for _, gw := range gws.Gateways {
		writer.Write([]byte(fmt.Sprintf("%s\t%s\t%t\t%3.2f\t%3.2f\t%3.2f\n",
			gw.Eui,
			gw.GetIp(),
			gw.GetStrictIp(),
			gw.GetLatitude(),
			gw.GetLongitude(),
			gw.GetAltitude())))
	}
	writer.Flush()
	return nil
}
