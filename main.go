package main

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type HugoAcend struct{}

func (d *Directory) Lint(ctx context.Context) error {

	cache := dag.CacheVolume("node-cache")

	_, err := dag.Container().From("node:21").WithDirectory("/src", d).WithMountedCache("/src/node_modules", cache).WithExec([]string{"npm", "ci"}).WithExec([]string{"npm", "run", "mdlint"}).Sync(ctx)
	return err
}

func (m *HugoAcend) BuildAndPushVariants(ctx context.Context, d *Directory, variants []string) error {
	eg, ctx := errgroup.WithContext(ctx)

	for _, v := range variants {
		eg.Go(func() error {
			return m.BuildAndPush(ctx, d, v)
		})
	}

	return eg.Wait()

}

func (m *HugoAcend) BuildAndPush(ctx context.Context, d *Directory, variant string) error {
	dag.Container().Build(d, ContainerBuildOpts{
		BuildArgs: []BuildArg{BuildArg{Name: "TRAINING_HUGO_ENV", Value: variant}},
	}).WithRegistryAuth("ttl.sh", "foo", dag.SetSecret("a", "bar")).Publish(ctx, "ttl.sh/csaucb9ub")
	return nil
}

func (d *Directory) DeployHelm(ctx context.Context, releaseName, namespace string) error {
	_, err := dag.Container().From("registry.puzzle.ch/cicd/alpine-base:latest").
		WithExec(
			[]string{"helm", "upgrade", releaseName, "acend-training-chart",
				"--install", "--wait", "--kubeconfig=$HOME/.kube/config", "--namespace", namespace,
				"--repo=https://acend.github.io/helm-charts/",
				"--values=helm-chart/values.yaml --atomic"}).Sync(ctx)
	return err
}
