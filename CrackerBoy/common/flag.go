package common

import (
	"CrackerBoy/util"
	"github.com/urfave/cli"
)

var Commands = []cli.Command{
	{
		Name:  "run",
		Usage: "begin the crack",
		Action: util.Begin,

		Flags: []cli.Flag{
			IntFlag("thread, t",10,"set the thread num"),
		},
	},
}

func IntFlag(name string,value int,usage string)cli.IntFlag{
	return cli.IntFlag{
		Name: name,
		Value: value,
		Usage: usage,
	}

}