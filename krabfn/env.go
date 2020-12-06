package krabfn

func env() {

	// param := function.Parameter{Name: "name", Type: cty.String}
	// fn := function.New(
	// 	&function.Spec{
	// 		Params: []function.Parameter{param},
	// 		Type:   function.StaticReturnType(cty.String),
	// 		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
	// 			return cty.StringVal(os.Getenv(args[0].AsString())), nil
	// 		},
	// 	},
	// )

	// evalContext := hcl.EvalContext{
	// 	Variables: map[string]cty.Value{
	// 		"local": cty.MapVal(map[string]cty.Value{"uri": cty.StringVal("postgres://uri")}),
	// 	},
	// 	Functions: map[string]function.Function{"env": fn},
	// }
	// tt := c.Uri.Variables()
	// for _, t := range tt {
	// 	fmt.Println("t", t)
	// 	a, b := t.TraverseAbs(&evalContext)
	// 	fmt.Println("T", t, a, b)
	// }
	// // if c.Name == "from_env" {
	// val, _ := c.Uri.Value(&evalContext)
	// c.UriVal = val.AsString()
	// }
}
