package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/lab5e/lospan/pkg/pb/lospan"
)

type outboxCmd struct {
	DeviceEUI string `kong:"help='Device EUI'"`
}

func (*outboxCmd) Run(args *params) error {
	p := args.Outbox

	client, ctx, done, err := createClient(args.Address)
	if err != nil {
		return err
	}
	defer done()

	res, err := client.Outbox(ctx, &lospan.OutboxRequest{Eui: p.DeviceEUI})
	if err != nil {
		return err
	}

	table := tabwriter.NewWriter(os.Stdout, 8, 3, 2, ' ', 0)
	table.Write([]byte("Port\tAck\tCreated\tSent\tAck time\tPayload\n"))
	for _, msg := range res.Messages {
		table.Write([]byte(fmt.Sprintf("%d\t%t\t%d\t%d\t%d\t%s\n",
			msg.Port, msg.Ack, msg.GetCreated(), msg.GetSent(), msg.GetAckTime(), ellipsisString(hex.EncodeToString(msg.Payload), 40))))
	}
	table.Flush()
	return nil
}
