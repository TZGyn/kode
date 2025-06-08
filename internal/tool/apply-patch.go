package tool

import "os/exec"

func ApplyPatch(patch string) (string, error) {
	checkGitCmd := exec.Command("apply_patch", "--help")

	stdout, err := checkGitCmd.Output()
	if err != nil {
		return string(stdout), err
	}

	return string(stdout), nil
}
