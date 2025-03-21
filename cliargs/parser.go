package cliargs

import (
	"flag"
)

type parser struct {
	args         []string
	flags        *flag.FlagSet
	stringValues map[string]*string
}

func New(args []string) *parser {
	flags := flag.NewFlagSet("", flag.ExitOnError)
	return &parser{
		args:         args,
		flags:        flags,
		stringValues: map[string]*string{},
	}
}

func (p *parser) Parse() error {
	err := p.flags.Parse(p.args)
	if err != nil {
		return err
	}

	return nil
}

func (p *parser) Add(name string) {
	p.stringValues[name] = p.flags.String(name, "", "")
}

func (p *parser) Args() []string {
	return p.flags.Args()
}

func (p *parser) Values() map[string]any {
	r := map[string]any{}

	for k, v := range p.stringValues {
		r[k] = *v
	}

	return r
}
