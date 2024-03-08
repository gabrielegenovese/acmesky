package travelPrefWorker

import (
	acmeskyEntities "acmesky/entities"
	travelPreferenceRepo "acmesky/repository/travel_preference"
	zbSingleton "acmesky/workers"
	ginContextRepo "acmesky/workers/utils/gin_context_repository"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
)

func RegisterWorkers() []worker.JobWorker {
	client := *zbSingleton.GetInstance()
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
		log.Println(fmt.Errorf("[BPMNERROR] failed to get variables for job %d: [%s]", job.Key, err))
	}
	var bpk string = vars["bpk"].(string)
	ginCtx := ginContextRepo.GetContext(bpk)

	flight_subscription := acmeskyEntities.CustomerFlightSubscriptionFromMap(vars)
	newPrefID, insertErr := travelPreferenceRepo.AddCustomerSubscribtionPreference(flight_subscription)

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	if insertErr != nil {
		err = insertErr
		_, err = client.NewFailJobCommand().
			JobKey(job.Key).
			Retries(0).
			ErrorMessage(err.Error()).
			Send(ctx)

		if err != nil {
			log.Println(fmt.Errorf("[BPMNERROR] error on failing job with key [%d]: [%s]", job.Key, err))
		} else {
			log.Println(insertErr)
		}
		ginCtx.Status(http.StatusInternalServerError)
		ginContextRepo.UnsetContext(bpk)
	} else {

		commandComplete, err := client.NewCompleteJobCommand().
			JobKey(job.Key).
			VariablesFromMap(map[string]interface{}{
				"travel_preference_id": newPrefID,
			})

		if err != nil {
			log.Println(fmt.Errorf("[BPMNERROR] failed to create command to complete job [%d] due to [%s]", job.Key, err))
		} else {
			_, err = commandComplete.Send(ctx)

			if err != nil {
				log.Panicf("[BPMNERROR] failed to complete job with key %d: [%s]", job.Key, err)
			} else {
				log.Printf("[BPMN] completed job %d successfully %d", job.Key, newPrefID)
			}
		}

		if err != nil {
			ginCtx.Status(http.StatusInternalServerError)
		} else {
			ginCtx.Status(http.StatusOK)
		}
		ginContextRepo.UnsetContext(bpk)
	}
}
