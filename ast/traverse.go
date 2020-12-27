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
	// TraverseDirection is used to represent the traverse direction
	// It is can be IN OUT BOTH
	TraverseDirection byte

	// TraverseVerb is used to represent a traverse verb
	TraverseVerb struct {
		Direction TraverseDirection
		TableName *TableName
	}

	// TraverseChain ise used to represent a group of action to traverse a graph
	TraverseChain struct {
		node
		Verbs []*TraverseVerb
	}
)

const (
	TraverseDirectionIn   TraverseDirection = 0
	TraverseDirectionOut  TraverseDirection = 1
	TraverseDirectionBoth TraverseDirection = 2
)

// String implements the fmt.Stringer interface
func (d TraverseDirection) String() string {
	switch d {
	case TraverseDirectionIn:
		return "IN"
	case TraverseDirectionOut:
		return "OUT"
	case TraverseDirectionBoth:
		return "BOTH"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", d)
	}
}

// Restore implements Node Accept interface.
func (t *TraverseChain) Restore(ctx *format.RestoreCtx) error {
	ctx.WriteKeyWord("TRAVERSE ")
	for i, v := range t.Verbs {
		ctx.WritePlainf("%s(", v.Direction.String())
		err := v.TableName.Restore(ctx)
		if err != nil {
			return err
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
