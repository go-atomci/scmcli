package gogit

import (
	"fmt"

	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestGoGit(projectPath string) *Client {

	projectPathSplit := strings.Split(strings.Split(projectPath, "://")[1], "/")
	projectName := strings.Join(projectPathSplit[1:], "/")

	return &Client{
		Project:   projectName,
		Workspace: "/tmp/abcd",
		GoGit: GoGit{
			URL:   projectPath,
			User:  "ci-cd-token",
			Token: "EhkBvVvSawydty_4HbPC",
		},
	}
}

func TestGetStatus(t *testing.T) {
	git := newTestGoGit("http://gitlab.example.com/colynn/hello")
	t.Run("Merge Branch", func(t *testing.T) {
		merged, msg, err := git.MergeBranch("devddd", "20191025_162032_28_1")
		assert.Nil(t, err)
		if err != nil {
			t.Fatalf("merge branch failed: %v", err)
		}
		fmt.Printf("merged: %v, msg: %v", merged, msg)
	})
}
