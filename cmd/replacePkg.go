package cmd

import (
	"fmt"
	"go-mod-tidy/coderepo"
)

var (
	ReplacePkgs map[string]string = map[string]string{
		"google.golang.org/grpc":      "github.com/grpc/grpc-go",
		"google.golang.org/appengine": "github.com/golang/appengine",
		"google.golang.org/genproto":  "github.com/google/go-genproto",
		"cloud.google.com/go":         "github.com/googleapis/google-cloud-go",
		"golang.org/x/crypto":         "github.com/golang/crypto",
		"golang.org/x/net":            "github.com/golang/net",
		"golang.org/x/oauth2":         "github.com/golang/oauth2",
		"golang.org/x/sync":           "github.com/golang/sync",
		"golang.org/x/tools":          "github.com/golang/tools",
		"golang.org/x/sys":            "github.com/golang/sys",
		"golang.org/x/text":           "github.com/golang/text",
		"golang.org/x/lint":           "github.com/golang/lint",
	}
)

type ModReplace map[string]string

func (mod ModReplace) save() {
	fmt.Println("\n\nSet Follow Replace To go.mod:")
	for k, v := range mod {
		fmt.Printf("%s => %s\n", k, v)
	}
}

func replaceMod(modRepos map[string]string) {
	modReplaces := ModReplace{}

	for originRepo, version := range modRepos {
		if version != "" {
			modReplaces[fmt.Sprintf("%s %s", originRepo, version)] = fmt.Sprintf("%s %s", ReplacePkgs[originRepo], version)
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
			modReplaces[fmt.Sprintf("%s %s", originRepo, codeVersion)] = fmt.Sprintf("%s %s", ReplacePkgs[originRepo], codeVersion)
		}
	}
	modReplaces.save()
}

func getVersion(repoInfo *coderepo.RevInfo) string {
	return fmt.Sprintf("v0.0.0-%s-%s", repoInfo.Time.UTC().Format("20060102150405"), repoInfo.Short)
}
