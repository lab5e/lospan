package main

import (
	"encoding/base64"
	"fmt"

	"github.com/lab5e/lospan/pkg/pb/lospan"
)

type sendCmd struct {
	DeviceEUI string `kong:"help='Device EUI',required"`
	Payload   string `kong:"help='Payload, base64 encoded',required"`
	Port      uint8  `kong:"help='Port for message',default=1"`
	Ack       bool   `kong:"help='Request message ack',default=false"`
}

func (*sendCmd) Run(args *params) error {
	p := args.Send

	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()

	buf, err := base64.StdEncoding.DecodeString(p.Payload)
	if err != nil {
		return err
	}
	req := &lospan.DownstreamMessage{
		Eui:     p.DeviceEUI,
		Payload: buf,
		Port:    int32(p.Port),
		Ack:     p.Ack,
	}
	_, err = client.SendMessage(ctx, req)
	if err != nil {
		return err
	}
	fmt.Println("Message queued in outbox")
	return nil
}
