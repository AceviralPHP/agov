package vhost

import (
	"io/ioutil"
	"os"
	"fmt"
	"strings"
)

type VHost struct {
	// url
	Domain string
	Alias  []string

	// file path
	DirRoot    string
	ErrorFile  string
	AccessFile string
	ConfigRoot string
	WebRoot    string
	LogRoot    string

	// listen
	IP   string
	Port string
}

// create a new instance of the vhost struct
func New(domain string) VHost {
	return VHost{
		Domain: domain,
		Alias: []string{},
		DirRoot: "/var/www/vhosts",
		WebRoot: "htdocs",
		LogRoot: "logs",
		ErrorFile: "error_log",
		AccessFile: "access_log",
		ConfigRoot: "/etc/httpd/conf.d",
		IP: "*",
		Port: "80",
	}
}

func (v *VHost) Exists() bool {
	return v.DirExists() &&
		v.ErrorLogExists() &&
		v.AccessLogExists() &&
		v.WebRootExists() &&
		v.ConfExists()
}

func (v *VHost) DirExists() bool {
	_, err := os.Stat(v.GetDirPath())

	return nil == err
}

func (v *VHost) ErrorLogExists() bool {
	_, err := os.Stat(v.GetErrorLogPath())

	return nil == err
}

func (v *VHost) AccessLogExists() bool {
	_, err := os.Stat(v.GetAccessLogPath())

	return nil == err
}

func (v *VHost) WebRootExists() bool {
	_, err := os.Stat(v.GetWebRootPath())

	return nil == err
}

func (v *VHost) ConfExists() bool {
	_, err := os.Stat(v.GetConfPath())

	return nil == err
}

func (v *VHost) GetDirPath() string {
	return fmt.Sprintf("%s/%s", v.DirRoot, v.Domain)
}

func (v *VHost) GetConfPath() string {
	return fmt.Sprintf("%s/%s.conf", v.ConfigRoot, v.Domain)
}

func (v *VHost) GetLogRootPath() string {
	return fmt.Sprintf("%s/%s/%s", v.DirRoot, v.Domain, v.LogRoot)
}

func (v *VHost) GetErrorLogPath() string {
	return fmt.Sprintf("%s/%s/%s/%s", v.DirRoot, v.Domain, v.LogRoot, v.ErrorFile)
}

func (v *VHost) GetAccessLogPath() string {
	return fmt.Sprintf("%s/%s/%s/%s", v.DirRoot, v.Domain, v.LogRoot, v.AccessFile)
}

func (v *VHost) GetWebRootPath() string {
	return fmt.Sprintf("%s/%s/%s", v.DirRoot, v.Domain, v.WebRoot)
}

func (v *VHost) Create() bool {
	if !v.CreateDirs() {
		return false
	}

	if !v.CreateConf() {
		return false
	}

	return true
}

func (v *VHost) CreateDirs() bool {
	if err := os.MkdirAll(v.GetWebRootPath(), 0644); nil != err {
		return false
	}

	if err := os.Mkdir(v.GetLogRootPath(), 0644); nil != err {
		return false
	}

	return true
}

func (v *VHost) CreateConf() bool {
	if v.ConfExists() {
		return false
	}

	return v.OverWriteConf()
}

func (v *VHost) OverWriteConf() bool {
	return nil == ioutil.WriteFile(v.GetConfPath(), []byte(v.CreateConfString()), 0644)
}

func (v *VHost) CreateConfString() string {
	confString := fmt.Sprintf(`<VirtualHost %s:%s>
  ServerName %s
`, v.IP, v.Port, v.Domain)

	if 0 < len(v.Alias) {
		confString += fmt.Sprintf("  ServerAlias %s\n", strings.Join(v.Alias, " "))
	}

	confString += fmt.Sprintf(`  DocumentRoot %s
  ErrorLog %s
  CustomLog %s combined
</VirtualHost>
`, v.GetWebRootPath(), v.GetErrorLogPath(), v.GetAccessLogPath())

	return confString
}

func (v * VHost) Remove() bool {
	return v.RemoveConfig() &&
		v.RemoveDirs()
}

func (v *VHost) RemoveConfig() bool {
	return nil != os.Remove(v.GetConfPath())
}

func (v *VHost) RemoveDirs() bool {
	return nil != os.RemoveAll(v.GetDirPath())
}

// This is basically the same as Remove but it wont return false
// if one or more of the operations fail
func (v *VHost) RevertChanges() {
	v.RemoveConfig()
	v.RemoveDirs()
}
