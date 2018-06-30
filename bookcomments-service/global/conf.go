package global

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/davecgh/go-spew/spew"
)

type conf struct {
	ServiceName     string    `yaml:"service_name"`
	ProjectRealPath string
	DB_BOOK         mysqlConf `yaml:"db_book"`
	Zipkin          zipkinConf
	Redis           redisConf
	HttpServer struct {
		Addr string
	} `yaml:"http_server"`
	GrpcServer struct {
		Addr string
	} `yaml:"grpc_server"`
	DebugServer struct {
		Addr string
	} `yaml:"debug_server"`
	MetricsServer struct {
		Addr string
	} `yaml:"metrics_server"`
}

type zipkinConf struct {
	Url         string
	ServiceName string `yaml:"service_name"`
	Reporter struct {
		Timeout       int
		BatchSize     int `yaml:"batch_size"`
		BatchInterval int `yaml:"batch_interval"`
		MaxBacklog    int `yaml:"max_backlog"`
	}
}

type mysqlConf struct {
	Username        string
	Password        string
	Host            string
	Port            int
	DBName          string `yaml:"db_name"`
	Driver          string
	Charset         string
	ParseTime       string `yaml:"parse_time"`
	Local           string
	ConnMaxLifeTime int    `yaml:"conn_max_life_time"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
}

type redisConf struct {
	Addr     string
	Password string
	DB       int
}

var Conf conf

const (
	RUN_MODE_LOCAL     = "local"
	RUN_MODE_CONTAINER = "container"
)

var ProjectRealPath = os.Getenv("GOPATH") + "/src/bookinfo/bookcomments-service"
var RuntimeRealPath = ProjectRealPath + "/runtime"
var LogPath = RuntimeRealPath + "/logs"

func loadConf() {

	os.MkdirAll(LogPath, os.ModePerm)

	log.Println(ProjectRealPath)
	if err := godotenv.Load(ProjectRealPath + "/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	runMode := os.Getenv("RUN_MODE")
	log.Println("run mode:", runMode)

	var confFile string
	switch runMode {
	case RUN_MODE_LOCAL:
		confFile = ProjectRealPath + "/conf/local.yaml"
	case RUN_MODE_CONTAINER:
		confFile = ProjectRealPath + "/conf/container.yaml"
	default:
		log.Fatalln("unsuppoer run mode! supports:[local,container]")
	}

	conf, _ := ioutil.ReadFile(confFile)
	if err := yaml.Unmarshal(conf, &Conf); err != nil {
		log.Fatalln("conf load failed", err)
	}

	spew.Dump(Conf)
}
