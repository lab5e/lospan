package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/lab5e/lospan/pkg/pb/lospan"
)

type devCmd struct {
	Add    addDevCmd    `kong:"cmd,help='Add device',aliases='create,a'"`
	Update updateDevCmd `kong:"cmd,help='Update device',aliases='up,u'"`
	Get    getDevCmd    `kong:"cmd,help='Get device',aliases='show,g,i'"`
	Del    delDevCmd    `kong:"cmd,help='Delete device',aliases='rm,delete,r,d'"`
	List   listDevCmd   `kong:"cmd,help='List devices',aliases='ls,l'"`
}

// This is common for both the add and update parameters; reuse
type deviceParameters struct {
	Type              string `kong:"help='Device type',enum='none,otaa,abp,disabled',default='none'"`
	DevAddr           string `kong:"help='Device address (4 byte, for ABP devices in hexadecimal)'"`
	AppSessionKey     string `kong:"help='Application session key (for ABP, 16 bytes hexadecimal)'"`
	NetworkSessionKey string `kong:"help='Network session key (for ABP, 16 bytes hexadecimal)'"`
	AppKey            string `kong:"help='Application key (for OTAA, 16 bytes hexadecimal)'"`
	FrameCountUp      int    `kong:"help='Frame counter up',default=-1"`
	FrameCountDown    int    `kong:"help='Frame counter up',default=-1"`
	RelaxedCounter    bool   `kong:"help='Relaxed counter',default=false"`
}

func printDevice(d *lospan.Device) {
	fmt.Printf("   EUI:              %s\n", d.GetEui())
	fmt.Printf("   AppEUI:           %s\n", d.GetApplicationEui())
	fmt.Printf("   State:            %s\n", d.GetState().String())
	fmt.Printf("   DevAddr:          %08x\n", d.GetDevAddr())
	fmt.Printf("   AppKey:           %s\n", hex.EncodeToString(d.AppKey))
	fmt.Printf("   AppSKey:          %s\n", hex.EncodeToString(d.AppSessionKey))
	fmt.Printf("   NwkSKey:          %s\n", hex.EncodeToString(d.NetworkSessionKey))
	fmt.Printf("   Frame count up:   %d\n", d.GetFrameCountUp())
	fmt.Printf("   Frame count down: %d\n", d.GetFrameCountDown())
	fmt.Printf("   Relaxed counter:  %t\n", d.GetRelaxedCounter())
	fmt.Printf("   Key warning:      %t\n", d.GetKeyWarning())
	fmt.Printf("   Nonce history:\n")
	for i := range d.DevNonces {
		fmt.Printf("        %d: %02x\n", i, d.DevNonces[i])
	}
}

func updateDeviceWithParameters(d *lospan.Device, p deviceParameters) error {
	if p.DevAddr != "" {
		v, err := strconv.ParseInt(p.DevAddr, 16, 32)
		if err != nil {
			return fmt.Errorf("invalid DevAddr value")
		}
		d.DevAddr = newPtr(uint32(v))
	}

	if p.Type == "abp" {
		d.State = newPtr(lospan.DeviceState_ABP)
	}
	if p.Type == "otaa" {
		d.State = newPtr(lospan.DeviceState_OTAA)
	}
	if p.Type == "disabled" {
		d.State = newPtr(lospan.DeviceState_DISABLED)
	}
	if p.AppKey != "" {
		buf, err := hex.DecodeString(p.AppKey)
		if err != nil || len(buf) != 16 {
			return fmt.Errorf("invalid AppKey")
		}
		d.AppKey = buf[:]
	}
	if p.AppSessionKey != "" {
		buf, err := hex.DecodeString(p.AppKey)
		if err != nil || len(buf) != 16 {
			return fmt.Errorf("invalid AppSessionKey")
		}
		d.AppSessionKey = buf[:]
	}
	if p.NetworkSessionKey != "" {
		buf, err := hex.DecodeString(p.AppKey)
		if err != nil || len(buf) != 16 {
			return fmt.Errorf("invalid NetworkSessionKey")
		}
		d.NetworkSessionKey = buf[:]
	}
	if p.FrameCountDown > -1 {
		d.FrameCountDown = newPtr(int32(p.FrameCountDown))
	}
	if p.FrameCountUp > -1 {
		d.FrameCountUp = newPtr(int32(p.FrameCountUp))
	}
	d.RelaxedCounter = newPtr(p.RelaxedCounter)
	return nil
}

type addDevCmd struct {
	EUI    string `kong:"help='Device EUI'"`
	AppEUI string `kong:"help='Application EUI',required"`
	deviceParameters
}

func (*addDevCmd) Run(args *params) error {
	p := args.Dev.Add
	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()
	if p.AppEUI == "" {
		return errors.New("missing application EUI")
	}
	newDevice := &lospan.Device{
		ApplicationEui: newPtr(p.AppEUI),
	}
	if p.EUI != "" {
		newDevice.Eui = newPtr(p.EUI)
	}
	if err := updateDeviceWithParameters(newDevice, p.deviceParameters); err != nil {
		return err
	}

	d, err := client.CreateDevice(ctx, newDevice)
	if err != nil {
		return err
	}
	fmt.Println("Created device")
	printDevice(d)
	return nil
}

type updateDevCmd struct {
	EUI    string `kong:"help='Device EUI',required"`
	AppEUI string `kong:"help='Application EUI'"`
	deviceParameters
}

func (*updateDevCmd) Run(args *params) error {
	p := args.Dev.Update
	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()
	newDevice := &lospan.Device{
		Eui: newPtr(p.EUI),
	}
	if p.AppEUI != "" {
		newDevice.ApplicationEui = newPtr(p.AppEUI)
	}
	if err := updateDeviceWithParameters(newDevice, p.deviceParameters); err != nil {
		return err
	}

	d, err := client.UpdateDevice(ctx, newDevice)
	if err != nil {
		return err
	}
	fmt.Println("Updated device")
	printDevice(d)
	return nil
}

type getDevCmd struct {
	EUI string `kong:"help='Device EUI',required"`
}

func (*getDevCmd) Run(args *params) error {
	p := args.Dev.Get
	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()
	d, err := client.GetDevice(ctx, &lospan.GetDeviceRequest{
		Eui: p.EUI,
	})
	if err != nil {
		return err
	}
	fmt.Println("Device")
	printDevice(d)
	return nil
}

type delDevCmd struct {
	EUI string `kong:"help='Device EUI',required"`
}

func (*delDevCmd) Run(args *params) error {
	p := args.Dev.Del
	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()
	d, err := client.DeleteDevice(ctx, &lospan.DeleteDeviceRequest{
		Eui: p.EUI,
	})
	if err != nil {
		return err
	}
	fmt.Println("Deleted device")
	printDevice(d)
	return nil
}

type listDevCmd struct {
	AppEUI string `kong:"help='Application EUI',required"`
}

func (*listDevCmd) Run(args *params) error {
	p := args.Dev.List
	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()
	devs, err := client.ListDevices(ctx, &lospan.ListDeviceRequest{
		ApplicationEui: p.AppEUI,
	})
	if err != nil {
		return err
	}

	table := tabwriter.NewWriter(os.Stdout, 8, 4, 2, ' ', 0)
	table.Write([]byte("EUI\tState\tDevAddr\tFCntUp\tFCNtDn\tRelaxed\tKey warning\n"))
	for _, dev := range devs.Devices {
		table.Write([]byte(fmt.Sprintf("%s\t%s\t%08x\t%d\t%d\t%t\t%t\n",
			dev.GetEui(),
			dev.GetState().String(),
			dev.GetDevAddr(),
			dev.GetFrameCountUp(),
			dev.GetFrameCountUp(),
			dev.GetRelaxedCounter(),
			dev.GetKeyWarning())))
	}
	table.Flush()
	return nil
}
