package main

type gwCmds struct {
	Add    gwAddCmd    `kong:"cmd,help='Add gateway',aliases='create,a'"`
	Del    gwDelCmd    `kong:"cmd,help='Delete gateway',aliases='rm,delete,d'"`
	Update gwUpdateCmd `kong:"cmd,help='Update gateway',aliases='up'"`
	Get    gwGetCmd    `kong:"cmd,help='Get gateway info'"`
	List   gwListCmd   `kong:"cmd,help='List gateways',aliases='ls'"`
}

type gwListCmd struct {
}

type gwAddCmd struct {
}

type gwUpdateCmd struct {
}

type gwDelCmd struct {
}

type gwGetCmd struct {
}
