package awaitility

import (
	"context"
	"fmt"
	"time"
)

func ExampleAwait() {
	checkCondition := func() bool {
		// check for message in queue (like kafka/nats)
		time.Sleep(1 * time.Second)
		return true
	}

	err := Await(1*time.Second, 5*time.Second, checkCondition)
	if err != nil {
		fmt.Printf("error occurred, %s\n", err)
	}

	fmt.Println("Everything is ok")

	// Output:
	// Everything is ok
}

func ExampleAwaitCtx() {
	checkCondition := func() bool {
		// check for message in queue (like kafka/nats)
		time.Sleep(1 * time.Second)
		return true
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := AwaitCtx(ctx, 1*time.Second, 5*time.Second, checkCondition)
	if err != nil {
		fmt.Printf("error occurred, %s\n", err)
	}

	fmt.Println("Everything is ok")

	// Output:
	// Everything is ok
}
