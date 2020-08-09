package krab

type EvalContext interface {
	// EvaluateBlock takes the given raw configuration block and associated
	// schema and evaluates it to produce a value of an object type that
	// conforms to the implied type of the schema.
	//
	// The "self" argument is optional. If given, it is the referenceable
	// address that the name "self" should behave as an alias for when
	// evaluating. Set this to nil if the "self" object should not be available.
	//
	// The "key" argument is also optional. If given, it is the instance key
	// of the current object within the multi-instance container it belongs
	// to. For example, on a resource block with "count" set this should be
	// set to a different addrs.IntKey for each instance created from that
	// block. Set this to addrs.NoKey if not appropriate.
	//
	// The returned body is an expanded version of the given body, with any
	// "dynamic" blocks replaced with zero or more static blocks. This can be
	// used to extract correct source location information about attributes of
	// the returned object value.
	// EvaluateBlock(body hcl.Body, schema *configschema.Block, self addrs.Referenceable, keyData InstanceKeyEvalData) (cty.Value, hcl.Body, tfdiags.Diagnostics)

	// EvaluateExpr takes the given HCL expression and evaluates it to produce
	// a value.
	//
	// The "self" argument is optional. If given, it is the referenceable
	// address that the name "self" should behave as an alias for when
	// evaluating. Set this to nil if the "self" object should not be available.
	// EvaluateExpr(expr hcl.Expression, wantType cty.Type, self addrs.Referenceable) (cty.Value, tfdiags.Diagnostics)

	// EvaluationScope returns a scope that can be used to evaluate reference
	// addresses in this context.
	// EvaluationScope(self addrs.Referenceable, keyData InstanceKeyEvalData) *lang.Scope

	// SetModuleCallArguments defines values for the variables of a particular
	// child module call.
	//
	// Calling this function multiple times has merging behavior, keeping any
	// previously-set keys that are not present in the new map.
	// SetModuleCallArguments(addrs.ModuleCallInstance, map[string]cty.Value)

	// GetVariableValue returns the value provided for the input variable with
	// the given address, or cty.DynamicVal if the variable hasn't been assigned
	// a value yet.
	//
	// Most callers should deal with variable values only indirectly via
	// EvaluationScope and the other expression evaluation functions, but
	// this is provided because variables tend to be evaluated outside of
	// the context of the module they belong to and so we sometimes need to
	// override the normal expression evaluation behavior.
	// GetVariableValue(addr addrs.AbsInputVariableInstance) cty.Value
}

type DefaultEvalContext struct {
	Evaluator *Evaluator
}

// func (ctx *BuiltinEvalContext) EvaluateBlock(body hcl.Body, schema *configschema.Block, self addrs.Referenceable, keyData InstanceKeyEvalData) (cty.Value, hcl.Body, tfdiags.Diagnostics) {
// 	var diags tfdiags.Diagnostics
// 	scope := ctx.EvaluationScope(self, keyData)
// 	body, evalDiags := scope.ExpandBlock(body, schema)
// 	diags = diags.Append(evalDiags)
// 	val, evalDiags := scope.EvalBlock(body, schema)
// 	diags = diags.Append(evalDiags)
// 	return val, body, diags
// }

// func (ctx *BuiltinEvalContext) EvaluateExpr(expr hcl.Expression, wantType cty.Type, self addrs.Referenceable) (cty.Value, tfdiags.Diagnostics) {
// 	scope := ctx.EvaluationScope(self, EvalDataForNoInstanceKey)
// 	return scope.EvalExpr(expr, wantType)
// }

// func (ctx *BuiltinEvalContext) EvaluationScope(self addrs.Referenceable, keyData InstanceKeyEvalData) *lang.Scope {
// 	if !ctx.pathSet {
// 		panic("context path not set")
// 	}
// 	data := &evaluationStateData{
// 		Evaluator:       ctx.Evaluator,
// 		ModulePath:      ctx.PathValue,
// 		InstanceKeyData: keyData,
// 		Operation:       ctx.Evaluator.Operation,
// 	}
// 	return ctx.Evaluator.Scope(data, self)
// }

// func (ctx *BuiltinEvalContext) Path() addrs.ModuleInstance {
// 	if !ctx.pathSet {
// 		panic("context path not set")
// 	}
// 	return ctx.PathValue
// }

// func (ctx *BuiltinEvalContext) SetModuleCallArguments(n addrs.ModuleCallInstance, vals map[string]cty.Value) {
// 	ctx.VariableValuesLock.Lock()
// 	defer ctx.VariableValuesLock.Unlock()

// 	if !ctx.pathSet {
// 		panic("context path not set")
// 	}

// 	childPath := n.ModuleInstance(ctx.PathValue)
// 	key := childPath.String()

// 	args := ctx.VariableValues[key]
// 	if args == nil {
// 		args = make(map[string]cty.Value)
// 		ctx.VariableValues[key] = vals
// 		return
// 	}

// 	for k, v := range vals {
// 		args[k] = v
// 	}
// }

// func (ctx *BuiltinEvalContext) GetVariableValue(addr addrs.AbsInputVariableInstance) cty.Value {
// 	ctx.VariableValuesLock.Lock()
// 	defer ctx.VariableValuesLock.Unlock()

// 	modKey := addr.Module.String()
// 	modVars := ctx.VariableValues[modKey]
// 	val, ok := modVars[addr.Variable.Name]
// 	if !ok {
// 		return cty.DynamicVal
// 	}
// 	return val
// }
