package internal

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/leonardom/go-jar-launcher/configs"
)

type Command struct {
	Name string
	Args []string
}

type launcher struct {
	config *configs.Config
}

func NewLauncher(config *configs.Config) *launcher {
	return &launcher{
		config: config,
	}
}

func (l *launcher) Execute() error {
	command := l.getCommand()
	cmd := exec.Command(command.Name, command.Args...)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return errors.New("error executing [" + cmd.String() + "'] Error: " + err.Error())
	}
	return nil
}

func (l *launcher) getCommand() Command {
	var args []string
	args = append(args, l.config.JVMOptions...)
	args = append(args, "-jar")
	args = append(args, l.config.JARFile)
	args = append(args, l.config.Args...)
	return Command{
		Name: getJava(l.config.JavaHome),
		Args: args,
	}
}

func getJava(jre string) string {
	if len(strings.TrimSpace(jre)) > 0 {
		return strings.Join([]string{jre, "bin", "java"}, string(os.PathSeparator))
	}
	return "java"
}
