# Libmason
[![GoDoc](https://godoc.org/github.com/vdemeester/libmason?status.png)](https://godoc.org/github.com/vdemeester/libmason)
[![Build Status](https://travis-ci.org/vdemeester/libmason.svg?branch=master)](https://travis-ci.org/vdemeester/libmason)
[![Go Report Card](https://goreportcard.com/badge/github.com/vdemeester/libmason)](https://goreportcard.com/report/github.com/vdemeester/libmason)
[![License](https://img.shields.io/github/license/vdemeester/libmason.svg)]()

Libmason an helper library to build client-driven docker container image
builder. *It is still very experimental*.

The goal of `libmason` is to provide few helpers to ease the pain of
creating client-side docker image builder for those who find the
`Dockerfile` and `docker build` a little bit too limited.

It uses [engine-api](https://github.com/docker/engine-api) and is
pretty tied to it (some structs of `engine-api` are popping up for now).

## Helpers & Builders

As previously said, `libmason` provides some helpers to create
client-side builders, from the most low-level (almost `API` level) to
some higher level (with concept of Steps, commit/non-commit step,
etc…). Those *helpers* are designed to be composable.

### Base

The base Helper is located in the main package (`libmason`). It's a low level
interface (and implementation) of commands that might be needed for a
builder (get the image, create a container, commit a container to an
image, etc.).

```go
import (
    "github.com/vdemeester/libmason"
    "github.com/docker/engine-api/types"
    "github.com/docker/engine-api/types/container"
)
// […]

ctx := context.Background()

helper := libmason.Newhelper(client)
// […]

image, err := helper.GetImage(ctx, "busybox", types.ImagePullOptions{})
// […]

resp, err := helper.ContainerCreate(ctx, types.ContainerCreateConfig{
    Config: &container.Config{
        Image: image.ID,
    }
}
// […]

imageID, err := helper.ContainerCommit(ctx, resp.ID, types.ContainerCommitOptions{})
// […]
```

### Builder

The `builder` package currently holds a `StepBuilder` which consists
of a composition of Step executed in order.

```go
import (
    "github.com/vdemeester/libmason/builder"
)
// […]

steps := []Step{
    &MyStep{},
    // A step with that needs to create a container
    builder.WithDefaultCreate(&AnotherStep{}),
    // A step that will commit the container
    builder.WithCommit(&AThirdStep{}),
    // Or remove the container
    builder.WithRemove(&AFourthStep{}),
    // Or all of them ?
    builder.WithCreate(build.WithCommitAndRemove(&MyStep{})),
}

builder := builder.WithSteps(builder.DefaultBuilder(client))
image, err := builder.Run()
// […]
```

See the [godoc](https://godoc.org/github.com/vdemeester/libmason) on how to create steps.

