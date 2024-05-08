package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Group struct {
	Name  string
	X     string
	Gid   string
	Users []string
}

func execCommand(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)

	var buffer bytes.Buffer
	cmd.Stdout = &buffer

	err := cmd.Run()
	if err != err {
		return "", err
	}
	return buffer.String(), nil
}

func buildUser(ctx context.Context) (*Group, error) {

	input, err := execCommand(ctx, "getent", "group", "nixbld")
	if err != err {
		return nil, err
	}
	fmt.Println("getent", input)

	group := Group{}

	// NAME:X:GID:MEMBERS,...

	splits := strings.SplitN(input, ":", 4)

	group.Name = splits[0]
	group.X = splits[1]
	group.Gid = splits[2]
	group.Users = strings.Split(splits[3], ",")

	// fmt.Println(group)

	return &group, nil
}

func pgrep(ctx context.Context, user string) (string, error) {
	input, err := execCommand(ctx, "pgrep", "-fu", user)
	if err != err {
		return "", err
	}
	// fmt.Println(input)
	return input, nil
}

func activeBuildUsers(ctx context.Context) (string, error) {
	group, err := buildUser(ctx)
	if err != nil {
		return "", err
	}

	for _, user := range group.Users {

		fmt.Println(pgrep(ctx, user))

	}

	return "", nil
}
