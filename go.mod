module alb

go 1.13

require (
	bitbucket.org/ww/goautoneg v0.0.0-20120707110453-75cd24fc2f2c // indirect
	github.com/Microsoft/go-winio v0.4.12
	github.com/NYTimes/gziphandler v1.1.1
	github.com/PuerkitoBio/purell v1.1.1
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578
	github.com/beorn7/perks v1.0.0
	github.com/coreos/bbolt v1.3.2
	github.com/coreos/etcd v3.3.10+incompatible
	github.com/coreos/go-semver v0.3.0
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e
	github.com/coreos/pkg v0.0.0-20180108230652-97fdf19511ea
	github.com/davecgh/go-spew v1.1.1
	github.com/davyxu/cellnet v0.0.0-20180423063709-b8dc6aba971b
	github.com/davyxu/golog v0.0.0-20190227055902-72e4f3de2d15
	github.com/davyxu/goobjfmt v0.0.0-20180817064625-baf5de0715b1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v0.0.0-20170524085120-eef6495eddab
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.3.3
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/emicklei/go-restful v2.9.3+incompatible
	github.com/emicklei/go-restful-swagger12 v0.0.0-20170208215640-dcef7f557305
	github.com/evanphx/json-patch v4.1.0+incompatible
	github.com/fsnotify/fsnotify v1.4.7
	github.com/ghodss/yaml v1.0.0
	github.com/gin-contrib/sse v0.0.0-20190301062529-5545eab6dad3
	github.com/gin-gonic/gin v1.3.0
	github.com/go-openapi/jsonpointer v0.19.0
	github.com/go-openapi/jsonreference v0.19.0
	github.com/go-openapi/spec v0.19.0
	github.com/go-openapi/swag v0.19.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gogo/protobuf v1.2.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef
	github.com/golang/protobuf v1.3.2
	github.com/google/btree v1.0.0
	github.com/google/gofuzz v1.0.0
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.2.0
	github.com/gorilla/websocket v1.4.0
	github.com/gregjones/httpcache v0.0.0-20190212212710-3befbb6ad0cc
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.8.5
	github.com/hashicorp/golang-lru v0.5.1
	github.com/hashicorp/hcl v1.0.0
	github.com/imdario/mergo v0.3.7
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/jinzhu/gorm v1.9.4
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a
	github.com/joho/godotenv v1.3.0 // indirect
	github.com/jonboulle/clockwork v0.1.0
	github.com/json-iterator/go v1.1.6
	github.com/konsorten/go-windows-terminal-sequences v1.0.2
	github.com/magiconair/properties v1.8.0
	github.com/mailru/easyjson v0.0.0-20190403194419-1ea4449da983
	github.com/mattn/go-isatty v0.0.7
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/matttproud/golang_protobuf_extensions v1.0.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/pborman/uuid v0.0.0-20180906182336-adf5a7427709
	github.com/pelletier/go-toml v1.3.0
	github.com/petar/GoLLRB v0.0.0-20130427215148-53be0d36a84c
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.8.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90
	github.com/prometheus/common v0.3.0
	github.com/prometheus/procfs v0.0.0-20190412120340-e22ddced7142
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/soheilhy/cmux v0.1.4
	github.com/spf13/afero v1.2.2
	github.com/spf13/cast v1.3.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.4.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5
	github.com/ugorji/go v1.1.4
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2
	go.uber.org/zap v1.14.1 // indirect
	google.golang.org/grpc v1.19.0
	gopkg.in/gcfg.v1 v1.2.3
	gopkg.in/go-playground/validator.v8 v8.18.2
	gopkg.in/inf.v0 v0.9.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0-20170531160350-a96e63847dc3
	gopkg.in/square/go-jose.v2 v2.3.1
	gopkg.in/warnings.v0 v0.1.2
	gopkg.in/yaml.v2 v2.2.2
	k8s.io/api v0.0.0-20181126191744-95336914c664
	k8s.io/apiextensions-apiserver v0.0.0-20181126195113-57b8dbfcc51a
	k8s.io/apimachinery v0.0.0-20181126123124-70adfbae261e
	k8s.io/apiserver v0.0.0-20181126193023-8afb6ff1462d
	k8s.io/client-go v0.0.0-20190117233410-4022682532b3
	k8s.io/kube-openapi v0.0.0-20181018171734-e494cc581111
	k8s.io/kubernetes v1.11.8
)
