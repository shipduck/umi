package main

import "testing"

func TestContextClone(t *testing.T) {
	ctx := Context{Title: "foo"}
	clone := ctx.Clone()
	clone.Title = "bar"

	if ctx.Title == clone.Title {
		t.Errorf("ContextClone - expected %q, got %q", ctx, clone)
	}
}
