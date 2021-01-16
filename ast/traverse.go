// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package ast

import (
	"fmt"

	"github.com/pingcap/parser/format"
)

var _ Node = &TraverseChain{}

type (
	// TraverseAction is used to represent the traverse direction
	// It is can be IN OUT BOTH
	TraverseAction byte

	TraverseTarget struct {
		Name  *TableName
		Where ExprNode
	}

	// TraverseVerb is used to represent a traverse verb
	TraverseVerb struct {
		Action  TraverseAction
		Targets []*TraverseTarget
	}

	// TraverseChain ise used to represent a group of action to traverse a graph
	TraverseChain struct {
		node
		Verbs []*TraverseVerb
	}
)

const (
	TraverseActionIn   TraverseAction = 0
	TraverseActionOut  TraverseAction = 1
	TraverseActionBoth TraverseAction = 2
	TraverseActionTags TraverseAction = 3
)

// String implements the fmt.Stringer interface
func (d TraverseAction) String() string {
	switch d {
	case TraverseActionIn:
		return "IN"
	case TraverseActionOut:
		return "OUT"
	case TraverseActionBoth:
		return "BOTH"
	case TraverseActionTags:
		return "TAGS"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", d)
	}
}

// Restore implements Node Accept interface.
func (t *TraverseChain) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("TRAVERSE ")
	for i, v := range t.Verbs {
		ctx.WritePlainf("%s(", v.Action.String())
		for j, n := range v.Targets {
			if err := n.Name.Restore(ctx); err != nil {
				return err
			}
			if err := n.Where.Restore(ctx); err != nil {
				return err
			}
			if j != len(v.Targets)-1 {
				ctx.WritePlain(",")
			}
		}
		ctx.WritePlain(")")
		if i != len(t.Verbs)-1 {
			ctx.WritePlain(".")
		}
	}
	return nil
}

// Accept implements Node Accept interface.
func (t *TraverseChain) Accept(v Visitor) (node Node, ok bool) {
	newNode, _ := v.Enter(t)
	// TODO: visit all children expression if we support filter expression in Verbs in the future
	return v.Leave(newNode.(*TraverseChain))
}
