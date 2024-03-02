package zbSingleton

import (
	"log"
	"os"
	"sync"

	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

var lockZeebeInstance = &sync.Mutex{}
var zeebeInstance *zbc.Client

func createClient() zbc.Client {

	config := zbc.ClientConfig{
		GatewayAddress:         os.Getenv("ZEEBE_HOST"),
		UsePlaintextConnection: true,
	}

	client, err := zbc.NewClient(&config)
	if err != nil {
		panic(err)
	}

	return client
}

func closeClient(client zbc.Client) {
	log.Println("closing client")
	_ = client.Close()
}

func GetInstance() *zbc.Client {
	var client zbc.Client

	if zeebeInstance == nil {
		lockZeebeInstance.Lock()
		defer lockZeebeInstance.Unlock()
		if zeebeInstance == nil {
			client = createClient()
			zeebeInstance = &client
		}
	}

	return zeebeInstance
}
