package gocom

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"testing"
)

func TestFiberContext_InvokeNativeCtx(t *testing.T) {
	type fields struct {
		ctx *fiber.Ctx
	}
	type args struct {
		handlerFunc interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "INITIAL",
			fields: fields{
				ctx: &fiber.Ctx{},
			},
			args: args{
				handlerFunc: func(inner *fiber.Ctx) error {
					fmt.Println("HELLO WORLD")
					return nil
				},
			},
			wantErr: false,
		}, {
			name: "INVALID",
			fields: fields{
				ctx: &fiber.Ctx{},
			},
			args: args{
				handlerFunc: func(inner Context) error {
					return nil
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &FiberContext{
				ctx: tt.fields.ctx,
			}
			if err := o.InvokeNativeCtx(tt.args.handlerFunc); (err != nil) != tt.wantErr {
				t.Errorf("InvokeNativeCtx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
