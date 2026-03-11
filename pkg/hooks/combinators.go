package hooks

// OnBlock returns a ChainLink that only runs when the previous result has Block=true.
func OnBlock(fn ChainFunc) ChainLink {
	return ChainLink{
		Func:      fn,
		Condition: func(r *Result) bool { return r != nil && r.Block },
	}
}

// OnError returns a ChainLink that only runs when the previous result has a non-empty Error.
func OnError(fn ChainFunc) ChainLink {
	return ChainLink{
		Func:      fn,
		Condition: func(r *Result) bool { return r != nil && r.Error != "" },
	}
}

// Transform returns a ChainLink that always runs and can modify the result.
func Transform(fn ChainFunc) ChainLink {
	return ChainLink{
		Func: fn,
	}
}

// HandlerLink wraps an existing Handler as a ChainLink.
func HandlerLink(h Handler) ChainLink {
	return ChainLink{Handler: h}
}
