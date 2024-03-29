package zeebeUtils

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	zeebeEntities "github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
)

/*
notifies Zeebe of the fail and decrease the retries count and try again after the specified delay
if it is an error due to a rate limit or similar and returns true. false otherwise or null and do nothing.
*/
func Handle_BP_fail_allow_retry(client worker.JobClient, job zeebeEntities.Job, err error, retryBackoff time.Duration) bool {
	var handled bool = false
	if err != nil {
		if strings.Contains(err.Error(), "Error 1040: Too many connections") {
			ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancelFn()

			command, errCmd := client.
				NewFailJobCommand().
				JobKey(job.Key).
				Retries(job.GetRetries() - 1).
				RetryBackoff(retryBackoff).
				ErrorMessage(err.Error()).
				VariablesFromMap(map[string]interface{}{
					"errorCode": "DB_ERROR",
					"errorMsg":  err.Error(),
				})

			if errCmd != nil {
				log.Println(fmt.Errorf("[BPMNERROR] error on failing job with key [%d]: [%s]", job.Key, errCmd))
			}

			command.Send(ctx)
			handled = true
			return handled
		}
	}
	return handled
}
