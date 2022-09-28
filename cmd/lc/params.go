package main

type params struct {
	Address string `kong:"help='Address of lora server API',default='127.0.0.1:4711'"`
	App     appCmd `kong:"cmd,help='Application commands',aliases='application,a'"`
	Dev     devCmd `kong:"cmd,help='Device commands',aliases='device,d'"`
	GW      gwCmds `kong:"cmd,help='Gateway commands',aliases='gateway,g'"`
}
