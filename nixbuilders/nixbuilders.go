package nixbuilders

import (
	"bytes"
	"context"
	"io/fs"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"

	"github.com/samber/lo"
	"github.com/shirou/gopsutil/process"
	"github.com/wfd3/go-groups/src/group"
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

	return strings.Trim(buffer.String(), " \n"), nil
}

func pgrep(ctx context.Context, users []string) (map[string][]*process.Process, error) {

	pr, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return nil, err
	}

	prs := lo.GroupBy(pr, func(p *process.Process) string {

		userName, _ := p.UsernameWithContext(ctx)

		return userName

	})

	return prs, nil
}

type ActiveUser struct {
	User        string
	UserObj     *user.User
	Processes   []*process.Process
	RootProcess *process.Process
	Dir         fs.DirEntry
}

func (a ActiveUser) DirName() string {
	return strings.ReplaceAll(a.Dir.Name(), "nix-build-", "")
}

func activeBuildUsers(ctx context.Context, users []string) ([]ActiveUser, error) {

	processesByUser, err := pgrep(ctx, users)
	if err != nil {
		return nil, err
	}

	dirs, _ := os.ReadDir("/tmp")

	activeUsers := []ActiveUser{}
	for _, userName := range users {

		processes := processesByUser[userName]
		userObj, _ := user.Lookup(userName)

		uid, err := strconv.ParseUint(userObj.Uid, 10, 32)
		if err != nil {
			break
		}

		if len(processes) > 0 {

			dir, _ := lo.Find(dirs, func(dir fs.DirEntry) bool {
				info, err := dir.Info()
				if err != nil {
					return false
				}

				obj := info.Sys()
				// fmt.Println(obj)

				if obj, ok := obj.(*syscall.Stat_t); ok {
					return obj.Uid == uint32(uid)
				}

				return false
			})

			rootProcess, _ := lo.Find(processes, func(pr *process.Process) bool {
				return true
			})

			activeUsers = append(activeUsers, ActiveUser{
				User:        userName,
				UserObj:     userObj,
				Processes:   processes,
				RootProcess: rootProcess,
				Dir:         dir,
			})
		}

	}

	return activeUsers, nil
}

type ActiveBuildersResponse = *[]ActiveUser

func GetActiveBuilders(ctx context.Context) (ActiveBuildersResponse, error) {

	gr, err := group.Lookup("nixbld")
	if err != nil {
		return nil, err
	}

	active, err := activeBuildUsers(ctx, gr.Members)
	if err != nil {
		return nil, err
	}

	return &active, nil
}
