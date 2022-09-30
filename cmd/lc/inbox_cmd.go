package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/lab5e/lospan/pkg/pb/lospan"
)

type inboxCmd struct {
	DeviceEUI string `kong:"help='Device EUI'"`
}

func (*inboxCmd) Run(args *params) error {
	p := args.Inbox

	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()

	res, err := client.Inbox(ctx, &lospan.InboxRequest{Eui: p.DeviceEUI})
	if err != nil {
		return err
	}

	table := tabwriter.NewWriter(os.Stdout, 8, 3, 2, ' ', 0)
	table.Write([]byte("DevAddr\tGateway\tData rate\tRSSI\tSNR\tFrequency\tPayload\n"))
	for _, msg := range res.Messages {
		table.Write([]byte(fmt.Sprintf("%08x\t%s\t%s\t%d\t%3.2f\t%3.2f\t%s\n",
			msg.DevAddr, msg.GatewayEui, msg.DataRate, msg.Rssi, msg.Snr, msg.Frequency, ellipsisString(hex.EncodeToString(msg.Payload), 40))))
	}
	table.Flush()
	return nil
}
