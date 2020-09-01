package agent

import (
	"io"
	"text/template"

	"code.htres.cn/casicloud/alb/pkg/model"
)

const hacfgTmpl = `
global
  daemon
  maxconn 2048

defaults
  mode tcp
  balance roundrobin
  timeout connect 5000ms
  timeout client 50000ms
  timeout server 50000ms

{{if eq .Record.Type 0}}
listen listener-{{.Record.IP.NullString.String}}_{{.Record.Port}}
  stick-table type ip size 200k expire 30m
  stick on src
  bind *:{{.Record.Port}}
  {{range .Endpoints}}server {{.Name}} {{.IP}}:{{.Port}} check inter 10s 
  {{end}}
{{else}}
listen listener-{{.Record.Domain.NullString.String}}
  mode http
  stick-table type ip size 200k expire 30m
  stick on src
  bind {{.Record.Domain.NullString.String}}
  {{range .Endpoints}}server {{.Name}} {{.IP}}:{{.Port}} check inter 10s 
  {{end}}
{{end}}
`

//WriteHaproxyCfg output haproxy.cfg
func WriteHaproxyCfg(wr io.Writer, policy model.LBPolicy) error {
	tmpl := template.Must(template.New("t1").Parse(hacfgTmpl))
	err := tmpl.Execute(wr, policy)
	if err != nil {
		sysLogger.WithField("reason", err).Error("Create haproxy configfile failed")
		return err
	}
	return nil
}

// keepalived configuration template
const keepAlivedCtgTmpl = `
vrrp_instance ALB_KEEPALIVE {
  interface {{ .INet }}
  virtual_router_id {{ .VirutalRouterID }}

  state {{ .State }}
  priority {{ .Priority }}

  unicast_src_ip {{ .UnicastSrcIP }}
  unicast_peer {
      {{- range .UnicastPeer}} 
      {{.}}
      {{- end}}
  }
  virtual_ipaddress {
      {{ .VirtualIP }}
  }
}
`

// WriteKeepalivedCfg  create the keepalived.conf
func WriteKeepalivedCfg(wr io.Writer, c KeepalivedConfig) error {
	tmpl := template.Must(template.New("t1").Parse(keepAlivedCtgTmpl))
	err := tmpl.Execute(wr, c)
	if err != nil {
		sysLogger.WithField("reason", err).Error("Create keepalived configfile failed")
		return err
	}
	return nil
}
