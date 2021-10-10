package krab

type Argument struct {
	Name string `hcl:"name,label"`
	Type string `hcl:"type,optional"`
}

// Arguments represents command line arguments or params that you can pass to action.
//
type Arguments struct {
	Args []*Argument `hcl:"arg,block"`
}

func (a *Arguments) Validate() error {
	for _, a := range a.Args {
		if err := a.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (a *Arguments) InitDefaults() {
	for _, a := range a.Args {
		a.InitDefaults()
	}
}

func (a *Argument) Validate() error {
	return nil
}

func (a *Argument) InitDefaults() {
	if a.Type == "" {
		a.Type = "string"
	}
}
