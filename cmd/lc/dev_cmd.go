package main

type devCmd struct {
	Add    addDevCmd    `kong:"cmd,help='Add device',aliases='create,a'"`
	Update updateDevCmd `kong:"cmd,help='Update device',aliases='up,u'"`
	Get    getDevCmd    `kong:"cmd,help='Get device',aliases='show,g,i'"`
	Del    delDevCmd    `kong:"cmd,help='Delete device',aliases='rm,delete,r,d'"`
	List   listDevCmd   `kong:"cmd,help='List devices',aliases='ls,l'"`
}

type addDevCmd struct {
	AppEUI string `kong:"help='Application EUI',required"`
	EUI    string `kong:"help='Device EUI'"`
}

func (*addDevCmd) Run(args *params) error {
	return nil
}

type updateDevCmd struct {
	EUI string `kong:"help='Device EUI',required"`
}

func (*updateDevCmd) Run(args *params) error {
	return nil
}

type getDevCmd struct {
	EUI string `kong:"help='Device EUI',required"`
}

func (*getDevCmd) Run(args *params) error {
	return nil
}

type delDevCmd struct {
	EUI string `kong:"help='Device EUI',required"`
}

func (*delDevCmd) Run(args *params) error {
	return nil
}

type listDevCmd struct {
}

func (*listDevCmd) Run(args *params) error {
	return nil
}
