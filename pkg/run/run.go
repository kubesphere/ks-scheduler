package run

import (
	"flag"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/soulseen/ks-pipeline-schduler/pkg/predicate"
	"github.com/soulseen/ks-pipeline-schduler/pkg/prioritize"
	"github.com/soulseen/ks-pipeline-schduler/pkg/routes"
	_ "github.com/soulseen/ks-pipeline-schduler/pkg/sqlite"

	log "github.com/golang/glog"
)

var (
	PipelinePriority = prioritize.Prioritize{
		Name:      "pipeline",
		Func:      prioritize.Pipeline,
	}

	TruePredicate = predicate.Predicate{
		Name: "alwaystrue",
		Func: predicate.AlwaysTrue,
	}
)

func Run() {
	//log.Info("Log level was set to ", strings.ToUpper(level.String()))

	flag.Parse()


	router := httprouter.New()
	routes.AddVersion(router)

	predicates := []predicate.Predicate{TruePredicate}
	for _, p := range predicates {
		routes.AddPredicate(router, p)
	}

	priorities := []prioritize.Prioritize{PipelinePriority}
	for _, p := range priorities {
		routes.AddPrioritize(router, p)
	}

	log.Info("info: server starting on the port :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
