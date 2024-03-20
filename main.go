package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/vistaarjuneja/dockertasks/tasks"
)

func main() {
	fmt.Println("Hello, World!")
	stageID := uuid.New().String()
	// setup
	req := tasks.SetupRequest(stageID)
	err := tasks.HandleSetup(req)
	if err != nil {
		panic(err)
	}

	// execute two steps
	step1ID := uuid.New().String()
	q := tasks.ExecRequest(step1ID, stageID, []string{"echo", "hello"})
	resp, err := tasks.HandleExec(q)
	if err != nil {
		panic(err)
	}
	fmt.Printf("poll response: %+v", resp)

	step2ID := uuid.New().String()
	r := tasks.ExecRequest(step2ID, stageID, []string{"sleep", "10"})
	resp, err = tasks.HandleExec(r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("poll response: %+v", resp)

	// cleanup
	d := tasks.DestroyRequest(stageID)
	err = tasks.HandleDestroy(d)
	if err != nil {
		panic(err)
	}

	fmt.Println("successfully completed!")

}
