package main

type params struct {
	Address string    `kong:"help='Address of lora server API',default='127.0.0.1:4711'"`
	App     appCmd    `kong:"cmd,help='Application commands',aliases='application,a'"`
	Dev     devCmd    `kong:"cmd,help='Device commands',aliases='device,d'"`
	GW      gwCmds    `kong:"cmd,help='Gateway commands',aliases='gateway,g'"`
	Inbox   inboxCmd  `kong:"cmd,help='Show upstream messages for devices',aliases='in,upstream,data'"`
	Outbox  outboxCmd `kong:"cmd,help='Show downstream messages for devices',aliases='out,downstream'"`
	Send    sendCmd   `kong:"cmd,help='Send message to device',aliase='s,msg'"`
}
