package gogit

import (
	"atomci/services/scmcli"
	"fmt"
	"os"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	"github.com/astaxie/beego/logs"
)

// GoGit ..
type GoGit struct {
	URL   string
	User  string
	Token string
}

// Client ..
type Client struct {
	GoGit
	Workspace string
	Project   string
}

// NewClient ..
func NewClient(url, user, token, workspace string) (scmcli.Provider, error) {
	client := &Client{}

	urlPathSplit := strings.Split(strings.Split(url, "://")[1], "/")
	projectName := strings.Join(urlPathSplit[1:], "/")
	logs.Debug("git projectpathsplit: %s,\tprojectName: %s", urlPathSplit, projectName)
	client.URL = url
	client.User = user
	client.Token = token
	client.Workspace = workspace
	client.Project = projectName
	return client, nil
}

// MergeBranch merge sourcebranch into targetBranch
// return bool mergeBranch status true mean success, false mean failure
// string mergebranch detail info;
func (ggc *Client) MergeBranch(sourceBranch, targetBranch string) (bool, string, error) {
	// prepare_workspace
	if err := ggc.prepareWorkspace(); err != nil {
		return false, "准备工作空间失败", fmt.Errorf("prepare workspace: %v", err)
	}
	// checkout target branch
	repo, err := ggc.checkoutBasedOnBranchParam(targetBranch)
	if err != nil {
		return false, fmt.Sprintf("拉取分支: %s 失败", targetBranch), fmt.Errorf("checkout: %v", err)
	}

	// merge branch into target branch
	err = ggc.mergeBranchIntoTargetBranch(repo, sourceBranch)
	if err != nil {
		return false, fmt.Sprintf("合并源分支：%s 至目标分支失败", sourceBranch), fmt.Errorf("merge into error: %v", err)
	}

	// push target branch
	err = ggc.pushBranchToRemote(targetBranch)
	if err != nil {
		return false, fmt.Sprintf("推送至远程分支失败，%v", targetBranch), fmt.Errorf("push branch: %v", err)
	}
	return true, fmt.Sprintf("应用：%s 分支%s合并至%s成功", ggc.Project, sourceBranch, targetBranch), nil
}

// prepareWorkspace
func (ggc *Client) prepareWorkspace() error {
	logs.Info("start prepare workspace")
	_, statErr := os.Stat(ggc.Workspace)

	if !os.IsNotExist(statErr) {
		logs.Info("clean workspace: %s", ggc.Workspace)
		if err := os.RemoveAll(ggc.Workspace); err != nil {
			logs.Error("when prepareWorkspace, os.RemoveAll occur error: %s", err.Error())
			return fmt.Errorf("removeAll occur error: %s", err.Error())
		}
	}
	if err := os.MkdirAll(ggc.Workspace, os.ModePerm); err != nil {
		logs.Error("when prepareWorkspace, makedirAll occur error:%s ", err.Error())
		return fmt.Errorf("mkdirAll occur error: %s", err.Error())
	}
	logs.Debug("create workspace: %s,", ggc.Workspace)
	return nil
}

// checkout code based on branch params
func (ggc *Client) checkoutBasedOnBranchParam(targetBranch string) (*git.Repository, error) {
	logs.Info("git clone %s %s", ggc.URL, ggc.Workspace)
	repo, err := git.PlainClone(ggc.Workspace, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: ggc.User,
			Password: ggc.Token,
		},
		URL:           ggc.URL,
		ReferenceName: plumbing.NewBranchReferenceName(targetBranch),
	})
	if err != nil {
		logs.Error("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
		return nil, fmt.Errorf("check branch: %v", err)
	}
	return repo, nil
}

// Merge branch into target branch
func (ggc *Client) mergeBranchIntoTargetBranch(repo *git.Repository, sourceBranch string) error {
	w, err := repo.Worktree()
	if err != nil {
		logs.Error("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
		return fmt.Errorf("w worktree: %v", err)
	}

	// merge remote source branch into target branch
	// update source branch
	err = w.Pull(&git.PullOptions{
		Auth: &http.BasicAuth{
			Username: ggc.User,
			Password: ggc.Token,
		},
		ReferenceName: plumbing.NewBranchReferenceName(sourceBranch),
	})

	return err
}

// push
func (ggc *Client) pushBranchToRemote(targetBranch string) error {
	r, err := git.PlainOpen(ggc.Workspace)
	if err != nil {
		logs.Error("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
		return err
	}
	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: ggc.User,
			Password: ggc.Token,
		},
	})
	if err != nil {
		logs.Error("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
		return err
	}
	return nil
}
