package travel_preference

import (
	acmeskyEntities "acmesky/entities"
	flight_subscription_repository "acmesky/repository"
	zbSingleton "acmesky/workers"
	"context"
	"log"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

func RegisterWorkers() []worker.JobWorker {
	var client zbc.Client
	client = *zbSingleton.GetInstance()
	workers := []worker.JobWorker{
		client.
			NewJobWorker().
			JobType("saveTravelPreference").
			Handler(HandleSaveTravelPreference).
			Open(),
	}
	return workers
}

func UnregisterWorkers(workers []worker.JobWorker) {
	for i := 0; i < len(workers); i++ {
		workers[i].Close()
	}
}

func HandleSaveTravelPreference(client worker.JobClient, job entities.Job) {

	vars, err := job.GetVariablesAsMap()
	if err != nil {
		log.Printf("failed to get variables for job %d: [%s]", job.Key, err)
		return
	}

	flight_subscription := acmeskyEntities.CustomerFlightSubscriptionFromMap(vars)
	flight_subscription_repository.AddCustomerSubscribtionPreference(flight_subscription)

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	_, err = client.NewCompleteJobCommand().JobKey(job.Key).Send(ctx)
	if err != nil {
		log.Printf("failed to complete job with key %d: [%s]", job.Key, err)
	} else {
		log.Printf("completed job %d successfully", job.Key)
	}

}
