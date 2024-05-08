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

func buildUser(ctx context.Context) (*Group, error) {

	input, err := execCommand(ctx, "getent", "group", "nixbld")
	if err != err {
		return nil, err
	}

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

func pgrep(ctx context.Context, user string) ([]int32, error) {
	input, err := execCommand(ctx, "pgrep", "-fu", user)
	if err != err {
		return []int32{}, err
	}

	pids := lo.FlatMap(strings.Split(input, "\n"), func(pid string, index int) []int32 {
		if pid == "" {
			return []int32{}
		}
		pidInt, _ := strconv.ParseInt(pid, 10, 32)
		return []int32{int32(pidInt)}
	})

	return pids, nil
}

type ActiveUser struct {
	User      string
	UserObj   *user.User
	Processes []*process.Process
	Dir       fs.DirEntry
}

func (a ActiveUser) DirName() string {
	return strings.ReplaceAll(strings.ReplaceAll(a.Dir.Name(), "nix-build-", ""), ".drv-0", "")
}

func activeBuildUsers(ctx context.Context, users []string) ([]ActiveUser, error) {

	dirs, _ := os.ReadDir("/tmp")

	activeUsers := []ActiveUser{}
	for _, userName := range users {

		userObj, _ := user.Lookup(userName)

		uid, _ := strconv.ParseUint(userObj.Uid, 10, 32)

		pids, _ := pgrep(ctx, userName)

		if len(pids) > 0 {

			processes := lo.Map(pids, func(pid int32, index int) *process.Process {
				processInfo, _ := process.NewProcess(pid)
				return processInfo
			})

			dir, _ := lo.Find(dirs, func(dir fs.DirEntry) bool {
				info, _ := dir.Info()
				obj := info.Sys()
				// fmt.Println(obj)

				if obj, ok := obj.(*syscall.Stat_t); ok {
					return obj.Uid == uint32(uid)
				}

				return false
			})

			activeUsers = append(activeUsers, ActiveUser{
				User:      userName,
				UserObj:   userObj,
				Processes: processes,
				Dir:       dir,
			})
		}

	}

	return activeUsers, nil
}

type ActiveBuildersResponse = *[]ActiveUser

func GetActiveBuilders(ctx context.Context) (ActiveBuildersResponse, error) {
	group, err := buildUser(ctx)
	if err != nil {
		return nil, err
	}

	// fmt.Println(group)

	active, err := activeBuildUsers(ctx, group.Users)
	if err != nil {
		return nil, err
	}

	return &active, nil
}
