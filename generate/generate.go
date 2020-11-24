package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/proemergotech/errors"
	"gopkg.in/yaml.v2"
)

type SkeletonMeta struct {
	Version string
	Config  Config
}

type Config struct {
	ProjectName   string `yaml:"project_name"`
	SchemaPackage string `yaml:"schema_package"`
	RedisCache    bool   `yaml:"redis_cache"`
	RedisStore    bool   `yaml:"redis_store"`
	RedisNotice   bool   `yaml:"redis_notice"`
	Elastic       bool   `yaml:"elastic"`
	PublicRest    bool   `yaml:"public_rest"`
	Examples      bool   `yaml:"examples"`
	ConfigFile    bool   `yaml:"config_file"`
	Bootstrap     bool   `yaml:"bootstrap"`
	Changelog     bool   `yaml:"changelog"`
	Test          bool   `yaml:"test"`
}

var tmpl *template.Template
var config *Config
var placeHolderRegex = regexp.MustCompile(`[^\s{}]*%:[^\s{}]*`)
var importCleanupRegex = regexp.MustCompile(`(?m)^.*DELETE.*$\n`)
var schemaDirSuffix = filepath.Join("schema", "skeleton")

func main() {
	var (
		metaFile []byte
		err      error
	)
	metaFile, config, err = processMeta()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	tmpl = template.New("skeleton").Funcs(sprig.TxtFuncMap())

	_ = os.RemoveAll("../output")
	err = processDir("../template", "../output")
	if err != nil {
		_ = os.RemoveAll("../output")
		log.Fatalf("%+v", err)
	}

	err = ioutil.WriteFile("../output/meta.yml", metaFile, 0644)
	if err != nil {
		log.Fatalf("failed writing meta.yml: %+v", err)
	}
}

func processFile(src, dst string) (err error) {
	si, err := os.Stat(src)
	if err != nil {
		return errors.Wrapf(err, "%s Stat failed", src)
	}

	inB, err := ioutil.ReadFile(src)
	if err != nil {
		return errors.Wrapf(err, "%s ReadFile failed", src)
	}
	in := string(inB)

	in = placeHolderRegex.ReplaceAllString(in, "")

	tmpl, err := tmpl.Clone()
	if err != nil {
		return errors.Wrap(err, "template Clone failed")
	}

	_, err = tmpl.Parse(in)
	if err != nil {
		return errors.Wrapf(err, "template Parse failed for %s", src)
	}

	out := bytes.NewBuffer(nil)
	err = tmpl.Execute(out, config)
	if err != nil {
		return errors.Wrapf(err, "template Execute failed for %s", src)
	}

	if strings.TrimSpace(out.String()) == "" {
		// do not copy empty files
		return nil
	}

	outB := out.Bytes()
	if strings.HasSuffix(dst, ".go") {
		outBFormatted, err := goimport(outB, config.ProjectName)
		if err != nil {
			err = errors.Wrapf(err, "%s goimport failed", dst)
			log.Printf("%+v", err.Error())
		} else {
			outB = outBFormatted
		}
	}

	dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, si.Mode())
	if err != nil {
		return errors.Wrapf(err, "%s OpenFile failed", dst)
	}
	defer func() {
		_ = dstFile.Close()
	}()

	_, err = dstFile.Write(outB)
	if err != nil {
		return errors.Wrapf(err, "%s Copy failed", dst)
	}

	_ = dstFile.Sync()

	return
}

func processDir(src string, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcInfo.IsDir() {
		return errors.Errorf("%s is not a directory", src)
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrapf(err, "%s Stat failed", dst)
	}
	if err == nil {
		return errors.Errorf("%s already exists", dst)
	}

	if strings.HasSuffix(dst, schemaDirSuffix) {
		dst = strings.ReplaceAll(dst, schemaDirSuffix, filepath.Join("schema", config.SchemaPackage))
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return errors.Wrapf(err, "%s MkdirAll failed", dst)
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return errors.Wrapf(err, "%s ReadDir failed", src)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = processDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = processFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	dstEntries, err := ioutil.ReadDir(dst)
	if err != nil {
		return errors.Wrapf(err, "%s ReadDir failed", dst)
	}

	// remove empty directory
	if len(dstEntries) == 0 {
		err = os.Remove(dst)
		if err != nil {
			return errors.Wrapf(err, "%s Remove failed", dst)
		}
	} else {
		log.Printf("%s done\n", dst)
	}

	return nil
}

func goimport(b []byte, projectName string) ([]byte, error) {
	cmd := exec.Command("goimports", "-local", "gitlab.com/proemergotech/"+projectName)
	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, errors.Wrap(err, "goimports command StdinPipe failed")
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "goimports command StdoutPipe failed")
	}
	err = cmd.Start()
	if err != nil {
		return nil, errors.Wrap(err, "goimports command Start failed")
	}

	_, err = in.Write(b)
	if err != nil {
		return nil, errors.Wrap(err, "goimports Write failed")
	}
	_ = in.Close()

	result, err := ioutil.ReadAll(out)
	if err != nil {
		return nil, errors.Wrap(err, "goimports ioutil.ReadAll failed")
	}

	err = cmd.Wait()
	if err != nil {
		return nil, errors.Wrap(err, "goimports command Wait failed")
	}

	return result, nil
}

func processMeta() ([]byte, *Config, error) {
	changeLog, err := ioutil.ReadFile("../CHANGELOG.md")
	if err != nil {
		return nil, nil, errors.Errorf("failed reading CHANGELOG.md: %+v", err)
	}
	changelogRegex := regexp.MustCompile(`## (v[0-9.]+)`)
	versionB := changelogRegex.Find(changeLog)
	if len(versionB) == 0 {
		return nil, nil, errors.Errorf("invalid CHANGELOG.md: version not found")
	}

	metaYml, err := ioutil.ReadFile("meta.yml")
	if err != nil {
		return nil, nil, errors.Errorf("failed reading meta.yml: %+v", err)
	}
	meta := make(map[string]interface{})
	err = yaml.Unmarshal(metaYml, &meta)
	if err != nil {
		return nil, nil, errors.Errorf("invalid meta.yml: %+v", err)
	}
	skeletonMetaI, ok := meta["skeleton"]
	if !ok {
		return nil, nil, errors.Errorf("invalid meta.yml: top level key skeleton is missing")
	}
	skeletonMetaB, err := yaml.Marshal(skeletonMetaI)
	if !ok {
		return nil, nil, errors.Errorf("invalid meta.yml: failed re-marshalling skeleton meta")
	}
	skeletonMeta := &SkeletonMeta{}
	err = yaml.Unmarshal(skeletonMetaB, skeletonMeta)
	if err != nil {
		return nil, nil, errors.Errorf("failed marshalling meta.yml skeleton: %+v", err)
	}

	skeletonMeta.Version = string(versionB[3:])

	meta["skeleton"] = skeletonMeta
	metaB, err := yaml.Marshal(meta)
	if err != nil {
		return nil, nil, errors.Errorf("failed marshalling meta.yml: %+v", err)
	}

	return metaB, &skeletonMeta.Config, nil
}
