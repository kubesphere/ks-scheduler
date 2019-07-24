package run

import (
	"flag"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/soulseen/ks-schduler/pkg/predicate"
	"github.com/soulseen/ks-schduler/pkg/prioritize"
	"github.com/soulseen/ks-schduler/pkg/routes"

	log "github.com/golang/glog"
)

func Run() {
	//log.Info("Log level was set to ", strings.ToUpper(level.String()))

	flag.Parse()

	router := httprouter.New()
	routes.AddVersion(router)

	predicates := []predicate.Predicate{predicate.TruePredicate}
	for _, p := range predicates {
		routes.AddPredicate(router, p)
	}

	priorities := []prioritize.Prioritize{prioritize.PipelinePriority}
	for _, p := range priorities {
		routes.AddPrioritize(router, p)
	}

	log.Info("info: server starting on the port :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
