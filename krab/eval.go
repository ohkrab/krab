package krab

type EvalNode interface {
	// Eval evaluates this node with the given context. The second parameter
	// are the argument values. These will match in order and 1-1 with the
	// results of the Args() return value.
	Eval(EvalContext) (interface{}, error)
}

type GraphNodeEvalable interface {
	EvalTree() EvalNode
}
