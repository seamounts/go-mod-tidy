package cmd

import (
	"bytes"
	"fmt"
	"github.com/seamounts/go-mod-tidy/coderepo"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

func init() {
	coderepo.WorkRoot = "/tmp/vcswork"
}

var goModTidy = &cobra.Command{
	Use:   "go-mod-tidy",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		runModTidy()
	},
}

func Execute() {
	if err := goModTidy.Execute(); err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}

func runModTidy() {
	var stderr bytes.Buffer
	var stdout bytes.Buffer

	gocmd := exec.Command("go", "mod", "tidy", "-v")
	gocmd.Stdout = &stdout
	gocmd.Stderr = &stderr
	gocmd.Run()
	if stderr.String() == "" {
		fmt.Println("all modules verified")
		return
	}
	parseStderr(stderr.String())
}

func parseStderr(stderr string) {
	outs := strings.Split(stderr, "\n")
	fetchFailedMods := getFailedPkgs(outs)
	replaceMod(fetchFailedMods)
}

func getFailedPkgs(outs []string) map[string]string {
	FailedPkgs := make(map[string]string, 0)
	notFoundPkgs := []string{}

	for _, line := range outs {
		switch {
		case strings.Contains(line, "unrecognized import path"):
			notFound := true
			for originRepo, _ := range ReplacePkgs {
				if !strings.Contains(line, originRepo) {
					continue
				}
				notFound = false
				if len(strings.Split(line, " ")) > 2 && strings.Contains(strings.Split(line, " ")[1], "@") {
					versionMod := strings.TrimSuffix(strings.Split(line, " ")[1], ":")
					FailedPkgs[originRepo] = strings.Split(versionMod, "@")[1]
				}
				break
			}
			if notFound {
				notFoundPkgs = append(notFoundPkgs, line)
			}

		case strings.HasPrefix(line, "https fetch failed"):
			notFound := true
			for originRepo, _ := range ReplacePkgs {
				if !strings.Contains(line, originRepo) {
					continue
				}
				notFound = false
				if _, ok := FailedPkgs[originRepo]; !ok {
					FailedPkgs[originRepo] = ""
				}
				break
			}
			if notFound {
				notFoundPkgs = append(notFoundPkgs, line)
			}
		}
	}
	if len(notFoundPkgs) > 0 {
		fmt.Println("not found pkgs:")
		for _, notFoundPkg := range notFoundPkgs {
			fmt.Println(notFoundPkg)
		}
	}

	return FailedPkgs
}
