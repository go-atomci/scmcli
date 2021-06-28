package scmcli

// Provider ..
type Provider interface {
	MergeBranch(soureBranch, TargerBranch string) (bool, string, error)
}
