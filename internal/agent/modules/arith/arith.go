package arith

import (
	"fmt"
	"golang.org/x/net/context"
)

type ArithCtrl struct {
}

func NewArithCtrl() *ArithCtrl {
	return &ArithCtrl{}
}

func (c *ArithCtrl) Add(ctx context.Context, args *Request, reply *Response) error {
	reply.C = args.A + args.B
	fmt.Printf("call: %d + %d = %d\n", args.A, args.B, reply.C)
	return nil
}
