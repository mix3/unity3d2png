package unity3d2png

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

type Extract struct {
	Unity3d string
	TGA     string
	PNG     string
}

type Service struct {
	Java     string
	Disunity string
	Convert  string
}

func (s Service) java() string {
	if s.Java == "" {
		return "java"
	}
	return s.Java
}

func (s Service) disunity() string {
	if s.Disunity == "" {
		return "./disunity.jar"
	}
	return s.Disunity
}

func (s Service) convert() string {
	if s.Convert == "" {
		return "convert"
	}
	return s.Convert
}

func (s Service) exists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func (s Service) Extract(file string) ([]Extract, error) {
	stdout, stderr, err := execCmd(s.java(), "-jar", s.disunity(), "extract", file)
	if 0 < len(stdout) {
		logrus.Info(string(stdout))
	}
	if 0 < len(stderr) {
		logrus.Warn(string(stderr))
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed disunity")
	}

	if !s.exists(strings.TrimSuffix(file, ".unity3d")) {
		return []Extract{}, nil
	}

	var tgas []string
	err = filepath.Walk(strings.TrimSuffix(file, ".unity3d"), func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".tga") {
			tgas = append(tgas, path)
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	es := []Extract{}
	for _, tga := range tgas {
		png := strings.TrimSuffix(tga, ".tga") + ".png"
		stdout, stderr, err := execCmd(s.convert(), tga, png)
		if 0 < len(stdout) {
			logrus.Info(string(stdout))
		}
		if 0 < len(stderr) {
			logrus.Warn(string(stderr))
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed convert")
		}
		es = append(es, Extract{
			Unity3d: file,
			TGA:     tga,
			PNG:     png,
		})
	}

	return es, nil
}

func execCmd(arg string, more ...string) ([]byte, []byte, error) {
	args := append([]string{arg}, more...)
	logrus.Info(args)
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}
