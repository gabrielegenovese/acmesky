package travelPrefWorker

import (
	acmeskyEntities "acmesky/entities"
	travelPreferenceRepo "acmesky/repository/travel_preference"
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
	newPrefID, insertErr := travelPreferenceRepo.AddCustomerSubscribtionPreference(flight_subscription)

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	if insertErr == nil {
		command, err := client.NewCompleteJobCommand().
			JobKey(job.Key).
			VariablesFromMap(map[string]interface{}{
				"travel_preference_id": newPrefID,
			})
		_, err = command.Send(ctx)
		if err != nil {

		}
	} else {
		_, err = client.NewFailJobCommand().
			JobKey(job.Key).
			Retries(1).
			ErrorMessage(err.Error()).
			Send(ctx)

	}

	if err != nil {
		log.Printf("failed to complete job with key %d: [%s]", job.Key, err)
	} else {
		log.Printf("completed job %d successfully", job.Key)
	}

}
