package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
)

const (
	goBinaryPath     string = "https://dl.google.com/go/%s.%s-%s.tar.gz"
	goInstallPath    string = "/usr/local"
	goVersionPattern string = `go version (.+) .+/.+`
)

type executer interface {
	Execute() error
}

type cliCommand struct {
	//interfaces
	executer

	//own
	flags *flag.FlagSet
}

func (command cliCommand) Execute() error {
	return command.flags.Parse(os.Args[2:])
}

type cliVersionCommand struct {
	cliCommand
}

func versionCommand() executer {
	return cliVersionCommand{
		cliCommand: cliCommand{
			flags: flag.NewFlagSet("version", flag.ExitOnError),
		},
	}
}

func (command cliVersionCommand) Execute() error {
	command.cliCommand.Execute()

	cmd := exec.Command("go", "version")

	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(out)
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	fmt.Println(string(regexp.MustCompile(goVersionPattern).FindSubmatch(body)[1]))

	return nil
}

type cliInitCommand struct {
	cliCommand
}

func initCommand() executer {
	return cliInitCommand{
		cliCommand: cliCommand{
			flags: flag.NewFlagSet("init", flag.ExitOnError),
		},
	}
}

func (command cliInitCommand) Execute() error {
	command.cliCommand.Execute()

	exePath, err := executableDir()
	if err != nil {
		return err
	}

	tplPath := path.Join(exePath, "tpl")
	files, err := ioutil.ReadDir(tplPath)
	if err != nil {
		return err
	}

	initPath, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		file, err := ioutil.ReadFile(path.Join(tplPath, fileInfo.Name()))
		if err != nil {
			return err
		}

		outPath := path.Join(initPath, fileInfo.Name())
		err = ioutil.WriteFile(outPath, file, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

type cliListCommand struct {
	cliCommand
}

func listCommand() executer {
	return cliListCommand{
		cliCommand: cliCommand{
			flags: flag.NewFlagSet("list", flag.ExitOnError),
		},
	}
}

func (command cliListCommand) Execute() error {
	command.cliCommand.Execute()

	releases, err := fetchReleases()
	if err != nil {
		return err
	}

	for _, release := range releases.names() {
		fmt.Println(release)
	}

	return nil
}

type cliInstallCommand struct {
	cliCommand

	installSource  *bool
	installBinary  *bool
	installVersion *string
}

func installCommand() executer {
	flagSet := flag.NewFlagSet("install", flag.ExitOnError)

	return cliInstallCommand{
		cliCommand: cliCommand{
			flags: flagSet,
		},
		installSource:  flagSet.Bool("source", false, "install from source"),
		installBinary:  flagSet.Bool("binary", false, "install from binary"),
		installVersion: flagSet.String("version", "", "version to install"),
	}
}

func (command cliInstallCommand) Execute() error {
	command.cliCommand.Execute()

	if !*command.installSource && !*command.installBinary {
		fmt.Println("Must specify install from source or binary")
	}

	if *command.installVersion == "" {
		fmt.Println("Version was not set")
		return nil
	}

	if *command.installSource {
		fmt.Println("Install Source unimplemented")
		return nil
	}

	if *command.installBinary {
		return doInstallBinary(*command.installVersion)
	}

	return nil
}

func doInstallBinary(version string) error {
	resp, err := http.Get(fmt.Sprintf(goBinaryPath, version, runtime.GOOS, runtime.GOARCH))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Status code %d received from source", resp.StatusCode)
	}

	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzipReader)
	if err != nil {
		return err
	}

	err = os.RemoveAll(path.Join(goInstallPath, "go"))
	if err != nil {
		return err
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		filepath := path.Join(goInstallPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(filepath, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
		case tar.TypeReg:
			file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			_, err = io.Copy(file, tarReader)
			if err != nil {
				return err
			}
		}
	}

	fmt.Println("Update path to include /usr/local/go/bin")

	return nil
}
