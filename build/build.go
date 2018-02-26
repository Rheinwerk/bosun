// Simple script to build bosun and scollector. This is not required, but it will properly insert version date and commit
// metadata into the resulting binaries, which `go build` will not do by default.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	shaFlag         = flag.String("sha", "", "SHA to embed.")
	buildBosun      = flag.Bool("bosun", false, "Only build Bosun.")
	buildTsdb       = flag.Bool("tsdbrelay", false, "Only build tsdbrelay")
	buildScollector = flag.Bool("scollector", false, "Only build scollector.")
	output          = flag.String("output", "", "Output directory; defaults to $GOPATH/bin.")
	esv5            = flag.Bool("esv5", false, "Build with esv5 support instead of v3")
	targetos        = flag.String("targetos", "", "Specify OS to compile for. Passed to the compiler os GOOS. Defaults to current OS.")

	allProgs = []string{"bosun", "scollector", "tsdbrelay"}
)

func main() {
	flag.Parse()
	// Get current commit SHA
	sha := *shaFlag
	if sha == "" {
		cmd := exec.Command("git", "rev-parse", "HEAD")
		cmd.Stderr = os.Stderr
		output, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		sha = strings.TrimSpace(string(output))
	}

	timeStr := time.Now().UTC().Format("20060102150405")
	ldFlags := fmt.Sprintf("-X bosun.org/_version.VersionSHA=%s -X bosun.org/_version.VersionDate=%s", sha, timeStr)

	progs := allProgs
	if *buildBosun {
		progs = []string{"bosun"}
	} else if *buildScollector {
		progs = []string{"scollector"}
	} else if *buildTsdb {
		progs = []string{"tsdbrelay"}
	}
	for _, app := range progs {
		fmt.Println("building", app)
		var args []string
		if *output != "" {
			suffix := ""
			if *targetos != "" {
				suffix = fmt.Sprintf("-%s", *targetos)
			}
			args = append(args, "build", "-o", fmt.Sprintf("%s%s", filepath.Join(*output, app), suffix))
		} else {
			args = append(args, "install")
		}
		if *esv5 {
			args = append(args, "-tags", "esv5")
		}
		args = append(args, "-ldflags", ldFlags, fmt.Sprintf("bosun.org/cmd/%s", app))
		fmt.Println("go", strings.Join(args, " "))
		cmd := exec.Command("go", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if *targetos != "" {
			env:=os.Environ()
			env = append(env, fmt.Sprintf("GOOS=%s", *targetos))
			cmd.Env = env
		}
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
