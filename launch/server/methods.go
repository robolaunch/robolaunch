package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	launchPb "github.com/robolaunch/robolaunch/api"
	launchflow "github.com/robolaunch/robolaunch/launch/pkg/workflow"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func (s *server) CreateLaunch(ctx context.Context, in *launchPb.CreateRequest) (*launchPb.LaunchState, error) {
	//Getting id token from grpc metadata
	headers, _ := metadata.FromIncomingContext(ctx)
	idToken := headers["authorization"][0]
	searchAttributes := map[string]interface{}{
		"DeploymentName":      in.Name,
		"DeploymentNamespace": in.Namespace,
	}
	log.Printf("---CreateLaunch---")
	//TODO: Run Workflow!
	c, err := client.NewClient(client.Options{
		HostPort: os.Getenv("TEMPORAL_SERVER_IP"),
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(idToken)

	options := client.StartWorkflowOptions{
		ID:               uuid.New().String(),
		TaskQueue:        launchflow.LaunchQueue,
		SearchAttributes: searchAttributes,
	}
	we, err := c.ExecuteWorkflow(ctx, options, launchflow.LaunchWorkflow, launchflow.LaunchRequest{
		Username:   in.GetUsername(),
		Name:       in.GetName(),
		Namespace:  in.GetNamespace(),
		IDToken:    idToken,
		Operation:  in.GetOperation(),
		LaunchType: in.GetLaunchType(),
	})
	if err != nil {
		return nil, err
	}
	//TODO: Query given Workflow

	resp, err := c.QueryWorkflow(context.Background(), we.GetID(), we.GetRunID(), "getLaunch")
	if err != nil {
		return nil, err
	}

	var status launchflow.LaunchStatus
	if err = resp.Get(&status); err != nil {
		return nil, err
	}

	return &launchPb.LaunchState{
		Username:       status.Username,
		Namespace:      status.Namespace,
		Name:           status.Name,
		LaunchType:     status.Namespace,
		WorkloadStatus: status.WorkloadStatus,
		TheiaPort:      status.TheiaPort,
		WebrtcPort:     status.WebRpcPort,
		NodeIp:         status.NodeIp,
	}, nil
}

func (s *server) OperateLaunch(ctx context.Context, in *launchPb.OperateRequest) (*launchPb.LaunchState, error) {
	// log.Printf("---OperateLaunch---")
	// //Getting id token from grpc metadata
	headers, _ := metadata.FromIncomingContext(ctx)

	idToken := headers["authorization"][0]
	c, err := client.NewClient(client.Options{
		HostPort: os.Getenv("TEMPORAL_SERVER_IP"),
	})
	if err != nil {
		return nil, err
	}
	// TODO: Get workflow properties

	// r, err := c.ListWorkflow(context.Background(), &workflowservice.ListWorkflowExecutionsRequest{
	// 	Query: `DeploymentNamespace="testplace" and DeploymentName="bestpiece"`,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// workflowId := r.Executions[0].Execution.GetWorkflowId()
	// runId := r.Executions[0].Execution.GetRunId()
	launchDetail, err := c.DescribeWorkflowExecution(context.Background(), in.GetWorkflowId(), in.GetRunId())
	if err != nil {
		panic(err)
	}
	name := strings.Trim(string(launchDetail.GetWorkflowExecutionInfo().SearchAttributes.IndexedFields["DeploymentName"].Data), `\"`)
	namespace := strings.Trim(string(launchDetail.GetWorkflowExecutionInfo().SearchAttributes.IndexedFields["DeploymentNamespace"].Data), `\"`)

	// From workflow list examine Advanced Query Api
	// TODO: Send Start & Stop signal according to incoming request
	err = c.SignalWorkflow(context.Background(), in.WorkflowId, in.RunId, "CHANGE_LAUNCH", launchflow.LaunchRequest{
		Username:   namespace,
		Name:       name,
		LaunchType: "",
		Namespace:  namespace,
		IDToken:    idToken,
		Operation:  in.Operation,
	})
	if err != nil {
		return nil, err
	}
	// //

	return &launchPb.LaunchState{
		Username:       "",
		Namespace:      "",
		Name:           "",
		LaunchType:     "",
		WorkloadStatus: "",
		TheiaPort:      0,
		WebrtcPort:     0,
		NodeIp:         "",
	}, nil
}

func (s *server) ListLaunch(in *launchPb.Empty, stream launchPb.Launch_ListLaunchServer) error {
	log.Printf("---OperateLaunch---")
	//TODO: Get query from here

	c, err := client.NewClient(client.Options{
		HostPort: os.Getenv("TEMPORAL_SERVER_IP"),
	})

	if err != nil {
		return err
	}
	var oldResult = []*launchPb.LaunchView{}
	for {
		r, err := c.ListWorkflow(context.Background(), &workflowservice.ListWorkflowExecutionsRequest{
			Namespace: "default",
			Query:     `DeploymentNamespace="tester"`,
		})
		if err != nil {
			return err
		}
		result := []*launchPb.LaunchView{}

		for _, workflow := range r.Executions {
			result = append(result, &launchPb.LaunchView{
				Name:       strings.Trim(string(workflow.SearchAttributes.IndexedFields["DeploymentName"].Data), `\"`),
				Namespace:  strings.Trim(string(workflow.SearchAttributes.IndexedFields["DeploymentNamespace"].Data), `\"`),
				WorkflowId: workflow.GetExecution().GetWorkflowId(),
				RunId:      workflow.GetExecution().GetRunId(),
			})

		}
		if !compareFlows(result, oldResult) {
			stream.Send(&launchPb.LaunchList{
				Launches: result,
			})
		}
		oldResult = result
		time.Sleep(time.Second * 2)
		err = stream.Context().Err()
		if err != nil {
			fmt.Println("No client. Terminated")
			return nil
		}

	}

}

//FIXME: Create good logic for it!
func compareFlows(a []*launchPb.LaunchView, b []*launchPb.LaunchView) bool {
	if len(a) != len(b) {
		fmt.Println("SEND!")
		return false
	}
	for i, workflow := range a {
		// if workflow != b[i] {
		// 	return false
		// }

		if !proto.Equal(workflow, b[i]) {
			return false
		}
	}
	return true
}
