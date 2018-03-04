package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"

	"github.com/go-spirit/go-spirit/builder/fetcher"
	"github.com/go-spirit/go-spirit/utils"
	"github.com/gogap/config"
)

type Project struct {
	Name       string
	conf       config.Configuration
	fetchers   map[string]fetcher.Fetcher
	needUpdate bool
}

type Builder struct {
	conf         config.Configuration
	projects     map[string]*Project
	projectsKeys []string
}

type Option func(*Options)

type Options struct {
	ConfigFile string
	NeedUpdate bool
}

type fetchRepo struct {
	repoConf   config.Configuration
	Fetcher    fetcher.Fetcher
	Url        string
	Revision   string
	NeedUpdate bool
}

func ConfigFile(file string) Option {
	return func(o *Options) {
		o.ConfigFile = file
	}
}

func NeedUpdate(needUpdate bool) Option {
	return func(o *Options) {
		o.NeedUpdate = needUpdate
	}
}

func (p *fetchRepo) Pull() (err error) {
	err = p.Fetcher.Fetch(p.Url, p.Revision, p.NeedUpdate, p.repoConf)
	return
}

func NewBuilder(opts ...Option) (builder *Builder, err error) {
	builderOpts := Options{}

	for _, o := range opts {
		o(&builderOpts)
	}

	conf := config.NewConfig(
		config.ConfigFile(builderOpts.ConfigFile),
	)

	var projs = make(map[string]*Project)
	var projKeys []string

	for _, projName := range conf.Keys() {
		var proj *Project
		proj, err = NewProject(projName, builderOpts.NeedUpdate, conf.GetConfig(projName))
		if err != nil {
			return
		}

		if _, exist := projs[projName]; exist {
			if exist {
				err = fmt.Errorf("project: %s already exist", projName)
				return
			}
		}

		projs[projName] = proj
		projKeys = append(projKeys, projName)
	}

	builder = &Builder{
		conf:         conf,
		projectsKeys: projKeys,
		projects:     projs,
	}

	return
}

func NewProject(projName string, needUpdate bool, conf config.Configuration) (proj *Project, err error) {

	if conf == nil {
		err = fmt.Errorf("could not inital project of %s config because of config is nil", projName)
		return
	}

	fetchers := make(map[string]fetcher.Fetcher)
	fetchersConf := conf.GetConfig("fetchers")

	if fetchersConf == nil {
		err = fmt.Errorf("could not inital project of %s config because of fetchers config is not set", projName)
		return
	}

	for _, fetcherName := range fetchersConf.Keys() {
		var f fetcher.Fetcher
		f, err = fetcher.NewFetcher(
			fetcherName,
			fetchersConf.GetConfig(fetcherName),
		)

		if err != nil {
			return
		}

		fetchers[fetcherName] = f
	}

	proj = &Project{
		Name:       projName,
		conf:       conf,
		fetchers:   fetchers,
		needUpdate: needUpdate,
	}

	return
}

func (p *Project) getFetchRepos() (repos []*fetchRepo, err error) {
	reposConf := p.conf.GetConfig("repos")

	if reposConf == nil {
		return
	}

	var fetchRepos []*fetchRepo

	for _, repoName := range reposConf.Keys() {
		repoConf := reposConf.GetConfig(repoName)
		if repoConf == nil {
			err = fmt.Errorf("repo's config is nil, project: %s, repo: %s", p.Name, repoName)
			return
		}

		url := repoConf.GetString("url")

		if len(url) == 0 {
			err = fmt.Errorf("repo's url is empty, project: %s, repo: %s", p.Name, repoName)
			return
		}

		f, exist := p.fetchers[repoConf.GetString("fetcher", "goget")]
		if !exist {
			err = fmt.Errorf("fetcher %s not exist, project: %s, repo: %s", p.fetchers[repoConf.GetString("fetcher", "goget")], p.Name, repoName)
			return
		}

		revision := repoConf.GetString("revision", "master")

		r := &fetchRepo{
			repoConf:   repoConf,
			Url:        url,
			Fetcher:    f,
			Revision:   revision,
			NeedUpdate: p.needUpdate,
		}

		fetchRepos = append(fetchRepos, r)
	}

	repos = fetchRepos

	return
}

func (p *Project) Pull() (err error) {
	repos, err := p.getFetchRepos()
	if err != nil {
		return
	}

	for _, repo := range repos {
		err = repo.Pull()
		if err != nil {
			return
		}
	}

	return
}

func (p *Project) Build() (err error) {
	pkgs := p.conf.GetStringList("packages")
	if len(pkgs) == 0 {
		return
	}

	buf := bytes.NewBuffer(nil)

	for _, pkg := range pkgs {
		buf.WriteString(fmt.Sprintf("import _ \"%s\"\n", pkg))
	}

	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	mainName := fmt.Sprintf("main_spirit_%s.go", p.Name)
	workdir := fmt.Sprintf("%s/%s", os.TempDir(), uuid.New())

	err = os.MkdirAll(workdir, 0755)
	if err != nil {
		return
	}

	mainPath := filepath.Join(workdir, mainName)

	mainSrc := strings.Replace(mainTmpl, "##imports##", buf.String(), 1)
	mainSrc = strings.Replace(mainSrc, "##config##", "`"+p.conf.String()+"`", 1)

	err = ioutil.WriteFile(mainPath, []byte(mainSrc), 0644)
	if err != nil {
		err = fmt.Errorf("write %s failure to temp dir: %s", mainName, err)
		return
	}

	mainSrc = strings.Replace(mainSrc, "\"##revision##\"", "`"+p.revisions(workdir)+"`", 1)
	err = ioutil.WriteFile(mainPath, []byte(mainSrc), 0644)
	if err != nil {
		err = fmt.Errorf("write %s failure to temp dir: %s", mainName, err)
		return
	}

	defer os.Remove(mainPath)

	// go get before build
	appendGetArgs := p.conf.GetStringList("build-args.go-get")
	gogetArgs := []string{"get", "-d"}
	gogetArgs = append(gogetArgs, appendGetArgs...)

	utils.ExecCommandSTDWD("go", workdir, gogetArgs...)

	// go build
	appendBuildArgs := p.conf.GetStringList("build-args.go-build")
	buildArgs := []string{"build"}
	buildArgs = append(buildArgs, appendBuildArgs...)
	buildArgs = append(buildArgs, "-o", filepath.Join(cwd, p.Name), mainPath)

	err = utils.ExecCommandSTD("go", buildArgs...)
	if err != nil {
		return
	}

	return
}

type packageRevision struct {
	Package  string `json:"package"`
	Branch   string `json:"branch"`
	Revision string `json:"revision"`
}

func (p *Project) revisions(wkdir string) string {

	pkgs, _ := utils.GoDeps(wkdir)

	if len(pkgs) == 0 {
		return ""
	}

	var pkgsRevision []packageRevision

	gopath := os.Getenv("GOPATH")

	for _, pkg := range pkgs {
		pkgPath := filepath.Join(gopath, "src", pkg)
		pkgHash, err := utils.GetCommitSHA(pkgPath)
		if err != nil {
			continue
		}
		branchName, err := utils.GetBranchOrTagName(pkgPath)
		if err != nil {
			continue
		}

		pkgsRevision = append(pkgsRevision, packageRevision{Package: pkg, Revision: pkgHash, Branch: branchName})
	}

	data, _ := json.MarshalIndent(pkgsRevision, "", "    ")

	return string(data)
}

func (p *Builder) ListProject() []string {
	var porj []string
	for _, c := range p.projectsKeys {
		porj = append(porj, c)
	}
	return porj
}

func (p *Builder) Build(porj ...string) (err error) {
	for _, c := range porj {
		logrus.WithField("project", c).Infoln("building")
		err = p.projects[c].Build()
		if err != nil {
			return
		}
	}
	return
}

func (p *Builder) Pull(porj ...string) (err error) {
	for _, c := range porj {
		err = p.projects[c].Pull()
		if err != nil {
			return
		}
	}
	return
}
