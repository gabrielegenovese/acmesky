package travelPrefWorker

import (
	acmeskyEntities "acmesky/entities"
	travelPreferenceRepo "acmesky/repository/travel_preference"
	zbSingleton "acmesky/workers"
	chanBPRepo "acmesky/workers/utils/channel_bp_repository"
	"context"
	"fmt"
	"log"
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
		client.
			NewJobWorker().
			JobType("responseTravelPreference").
			Handler(HandleResponseTravelPreference).
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

	flight_subscription := acmeskyEntities.CustomerFlightSubscriptionRequestFromMap(vars)
	newPrefID, insertErr := travelPreferenceRepo.AddCustomerSubscribtionPreference(flight_subscription)

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	if insertErr != nil {
		command, _ := client.NewThrowErrorCommand().
			JobKey(job.Key).
			ErrorCode("DB_ERROR").
			ErrorMessage(insertErr.Error()).
			VariablesFromMap(map[string]interface{}{
				// Note: in bpmn editor need to define variable mapping errorCode->errorCode;
				// otherwise errorCode is not forwarded (zeebe)
				"errorCode": "DB_ERROR",
				"errorMsg":  insertErr.Error(),
			})

		_, err = command.Send(ctx)

		if err != nil {
			log.Println(fmt.Errorf("[BPMNERROR] error on failing job with key [%d]: [%s]", job.Key, err))
		} else {
			log.Println(insertErr)
		}
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
	}
}

func HandleResponseTravelPreference(client worker.JobClient, job entities.Job) {

	vars, _ := job.GetVariablesAsMap()

	var bpk string = vars["bpk"].(string)
	result := chanBPRepo.GetContext(bpk)
	result <- vars

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	_, err := client.NewCompleteJobCommand().
		JobKey(job.Key).
		Send(ctx)

	if err != nil {
		log.Println(fmt.Errorf("[BPMNERROR] failed to create command to complete job [%d] due to [%s]", job.Key, err))
	}
}
