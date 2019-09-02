package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/seamounts/go-mod-tidy/coderepo"
	"github.com/spf13/cobra"
)

var (
	ReplacePkgs = map[string]string{
		"google.golang.org/grpc":      "github.com/grpc/grpc-go",
		"google.golang.org/appengine": "github.com/golang/appengine",
		"google.golang.org/genproto":  "github.com/google/go-genproto",
		"google.golang.org/api":       "github.com/googleapis/google-api-go-client",

		"cloud.google.com/go": "github.com/googleapis/google-cloud-go",
		"golang.org/x/crypto": "github.com/golang/crypto",
		"golang.org/x/net":    "github.com/golang/net",
		"golang.org/x/oauth2": "github.com/golang/oauth2",
		"golang.org/x/sync":   "github.com/golang/sync",
		"golang.org/x/tools":  "github.com/golang/tools",
		"golang.org/x/sys":    "github.com/golang/sys",
		"golang.org/x/text":   "github.com/golang/text",
		"golang.org/x/lint":   "github.com/golang/lint",
		"golang.org/x/term":   "github.com/golang/term",
		"golang.org/x/time":   "github.com/golang/time",
		"golang.org/x/vgo":    "github.com/golang/vgo",
		"golang.org/x/image":  "github.com/golang/image",
		"golang.org/x/exp":    "github.com/golang/exp",
		"golang.org/x/mobile": "github.com/golang/mobile",
	}
)

func init() {
	coderepo.WorkRoot = "/tmp/vcswork"
}

var GoMod = &cobra.Command{
	Use:   "go-mod-tidy",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		Run()
	},
}

func Run() {
	m := make(ModExt, 0)
	for {
		FailedPkgs, err := m.Tidy()
		if err != nil {
			log.Fatalln(err)
		}
		if len(FailedPkgs) == 0 {
			break
		}
		m.replace(FailedPkgs)
	}
}

type ModExt map[string]string

func (m ModExt) Tidy() (map[string]string, error) {
	FailedPkgs := make(map[string]string, 0)
	cmd := "go"
	params := []string{"mod", "tidy", "-v"}

	_, r, err := ExecCmd(cmd, params)
	if err != nil {
		return FailedPkgs, err
	}
	for {
		line, err := r.ReadString('\n')

		if io.EOF == err {
			break
		}
		if err != nil && io.EOF != err {
			log.Fatalln(err)
		}

		fmt.Println(line)
		switch {
		case strings.Contains(line, "unrecognized import path"):
			for originRepo, _ := range ReplacePkgs {
				if !strings.Contains(line, originRepo) {
					continue
				}
				if len(strings.Split(line, " ")) > 2 && strings.Contains(strings.Split(line, " ")[1], "@") {
					versionMod := strings.TrimSuffix(strings.Split(line, " ")[1], ":")
					FailedPkgs[originRepo] = strings.Split(versionMod, "@")[1]
				}
				break
			}

		case strings.HasPrefix(line, "https fetch failed"):
			for originRepo, _ := range ReplacePkgs {
				if !strings.Contains(line, originRepo) {
					continue
				}
				if _, ok := FailedPkgs[originRepo]; !ok {
					FailedPkgs[originRepo] = ""
				}
				break
			}
		}

	}
	return FailedPkgs, nil
}

func (m ModExt) replace(modRepos map[string]string) {
	for originRepo, version := range modRepos {
		if version != "" {
			m[fmt.Sprintf("%s@%s", originRepo, version)] = fmt.Sprintf("%s@%s", ReplacePkgs[originRepo], version)
		} else {
			repo, err := coderepo.NewGitRepo(fmt.Sprintf("https://%s", ReplacePkgs[originRepo]), false)
			if err != nil {
				fmt.Println("%v", err)
				return
			}
			info, err := repo.Stat("latest")
			if err != nil {
				fmt.Println("%v", err)
				return
			}
			codeVersion := getVersion(info)
			m[fmt.Sprintf("%s@%s", originRepo, codeVersion)] = fmt.Sprintf("%s@%s", ReplacePkgs[originRepo], codeVersion)
		}
	}
	if len(m) > 0 {
		m.save()
	}
}

func (m ModExt) save() {
	log.Println("replace mod...")

	cmd := "go"
	params := []string{"mod", "edit"}
	for originRepo, newRepo := range m {
		log.Println(fmt.Sprintf("-replace=%s=%s", originRepo, newRepo))

		params = append(params, fmt.Sprintf("-replace=%s=%s", originRepo, newRepo))
		if _, _, err := ExecCmd(cmd, params); err != nil {
			log.Fatalln(err)
		}
	}
}

func getVersion(repoInfo *coderepo.RevInfo) string {
	return fmt.Sprintf("v0.0.0-%s-%s", repoInfo.Time.UTC().Format("20060102150405"), repoInfo.Short)
}

func ExecCmd(cmd string, params []string) (*bufio.Reader, *bufio.Reader, error) {
	c := exec.Command(cmd, params...)

	stdout, err := c.StdoutPipe()
	errout, err := c.StderrPipe()

	if err != nil {
		return nil, nil, err
	}
	stdr := bufio.NewReader(stdout)
	errr := bufio.NewReader(errout)
	err = c.Start()

	return stdr, errr, err
}
