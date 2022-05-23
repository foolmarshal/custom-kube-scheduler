package extender

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/comail/colog"
	"github.com/julienschmidt/httprouter"

	v1 "k8s.io/api/core/v1"
)

const (
	versionPath      = "/version"
	apiPrefix        = "/scheduler"
	bindPath         = apiPrefix + "/bind"
	preemptionPath   = apiPrefix + "/preemption"
	predicatesPrefix = apiPrefix + "/predicates"
	prioritiesPrefix = apiPrefix + "/priorities"
)

var (
	version string // injected via ldflags at build time

	Graviton2Filter = Predicate{
		Name: "filter_graviton2",
		Func: func(pod v1.Pod, node v1.Node) (bool, error) {
			// 			log.Print("info: ", "filter_graviton2_allocatable_ephemeral_storage",
			// 				"\n\n ephemeral-storage = ", node.Status.Allocatable["ephemeral-storage"])
			// 			log.Print("info: ", "filter_graviton2_pod", "\n\n pod = ", pod)
			// 			log.Print("info: ", "filter_graviton2_requested_pod_spec", "\n\n pod-spec = ", pod.Spec)
			return true, nil

		},
	}
)

func StringToLevel(levelStr string) colog.Level {
	switch level := strings.ToUpper(levelStr); level {
	case "TRACE":
		return colog.LTrace
	case "DEBUG":
		return colog.LDebug
	case "INFO":
		return colog.LInfo
	case "WARNING":
		return colog.LWarning
	case "ERROR":
		return colog.LError
	case "ALERT":
		return colog.LAlert
	default:
		log.Printf("warning: LOG_LEVEL=\"%s\" is empty or invalid, falling back to \"INFO\".\n", level)
		return colog.LInfo
	}
}

func Run() {
	colog.SetDefaultLevel(colog.LInfo)
	colog.SetMinLevel(colog.LInfo)
	colog.SetFormatter(&colog.StdFormatter{
		Colors: true,
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
	})
	colog.Register()
	level := StringToLevel(os.Getenv("LOG_LEVEL"))
	log.Print("Log level was set to ", strings.ToUpper(level.String()))
	colog.SetMinLevel(level)

	router := httprouter.New()
	AddVersion(router)

	predicates := []Predicate{Graviton2Filter}
	for _, p := range predicates {
		AddPredicate(router, p)
	}
	log.Print("info: server starting on the port :80")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}
