package main

import (
	"context"
)

type HugoAcend struct{}

func (m *HugoAcend) MyFunction(ctx context.Context, stringArg string) (*Container, error) {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg}).Sync(ctx)
}

func (m *HugoAcend) Lint(ctx context.Context) (*Container, error) {
	//return nil, nil
	return dag.Container().From("node:21").WithExec([]string{"npm", "ci"}).WithExec([]string{"npm", "run", "mdlint"}).Sync(ctx)
}
