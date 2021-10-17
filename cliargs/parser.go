package cliargs

import (
	"errors"
	"flag"
)

type parser struct {
	args                []string
	flags               *flag.FlagSet
	requiredNonFlagArgs int
}

func New(args []string) *parser {
	flags := flag.NewFlagSet("", flag.ExitOnError)
	return &parser{args: args, flags: flags, requiredNonFlagArgs: -1}
}

func (p *parser) RequireNonFlagArgs(n int) {
	p.requiredNonFlagArgs = n
}

func (p *parser) Parse() error {
	err := p.flags.Parse(p.args)
	if err != nil {
		return err
	}

	if p.requiredNonFlagArgs != -1 {
		if len(p.flags.Args()) != p.requiredNonFlagArgs {
			return errors.New("Invalid number of arguments")
		}
	}

	return nil
}

func (p *parser) Add(name string) {
	p.flags.String(name, "", "")
}

func (p *parser) Args() []string {
	return p.flags.Args()
}
