package controllers

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"encoding/json"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

type Scanner struct {
	Hostname string `form:"hostname"`
	TCP      string `form:"tcp"`
	UDP      string `form:"udp"`
}
type Response struct {
	Hostname  string
	IPAddress string
	Protocol  string
	Port      string
	Status    string
}

func (c *MainController) Get() {

	c.TplName = "index.html"
}
func (c *MainController) Post() {
	var cmd *exec.Cmd
	scannerInfo := Scanner{}
	if err := c.ParseForm(&scannerInfo); err != nil {
		//error
	}
	if scannerInfo.TCP == "" {
		cmd = exec.Command("static/bin/scanport", "-host", scannerInfo.Hostname, "-udp", scannerInfo.UDP)
	} else if scannerInfo.UDP == "" {
		cmd = exec.Command("static/bin/scanport", "-host", scannerInfo.Hostname, "-tcp", scannerInfo.TCP)
	} else {
		cmd = exec.Command("static/bin/scanport", "-host", scannerInfo.Hostname, "-tcp", scannerInfo.TCP, "-udp", scannerInfo.UDP)
	}
	stdout, _ := cmd.Output()
	formatOut := formatPostData(&stdout)
	jsonOut,_ := json.Marshal(formatOut)
	c.Ctx.ResponseWriter.Write([]byte(jsonOut))
}

//Format and massage data
func formatPostData(stdout *[]byte) []Response {
	//   0              1            3   4       5
	//172.217.0.46 (google.com) ==> TCP/443 is open
	output := strings.Split(string(*stdout), "\n")
	digitMatch := regexp.MustCompile(`^\d.+`)
	parenMatch := regexp.MustCompile(`(\(|\))`)
	var responses []Response

	for _, i := range output {
		if digitMatch.MatchString(i) {
			var r Response
			j := strings.Split(i, " ")
			fqdn := parenMatch.ReplaceAllString(j[1], "")
			protoport := strings.Split(j[3], "/")
			r.Hostname = fqdn
			r.IPAddress = j[0]
			r.Protocol = protoport[0]
			r.Port = protoport[1]
			r.Status = j[5]
			responses = append(responses, r)
		}
	}
	fmt.Printf("%#v\n", responses)
	return responses
}
