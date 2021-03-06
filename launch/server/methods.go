package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
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

type UserInfo struct {
	Username string `json:"preferred_username"`
}

func (s *server) CreateLaunch(ctx context.Context, in *launchPb.CreateRequest) (*launchPb.LaunchState, error) {
	//Getting id token from grpc metadata
	headers, _ := metadata.FromIncomingContext(ctx)
	idToken := strings.TrimPrefix(headers["authorization"][0], "Bearer ")
	searchAttributes := map[string]interface{}{
		"DeploymentName":      in.Name,
		"DeploymentNamespace": in.Namespace,
		"DeploymentStatus":    "CREATING",
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

	idToken := strings.TrimPrefix(headers["authorization"][0], "Bearer ")
	c, err := client.NewClient(client.Options{
		HostPort: os.Getenv("TEMPORAL_SERVER_IP"),
	})
	if err != nil {
		return nil, err
	}
	if in.Operation == "GET" {
		resp, err := c.QueryWorkflow(context.Background(), in.GetWorkflowId(), in.GetRunId(), "getLaunch")
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
	headers, _ := metadata.FromIncomingContext(stream.Context())
	fmt.Println(headers)
	if headers["x-jwt"][0] == "" {
		return errors.New("parsed token is empty")
	}

	decoded, err := base64.StdEncoding.DecodeString(headers["x-jwt"][0])
	if err != nil {
		return err
	}

	var user = UserInfo{}
	err = json.Unmarshal(decoded, &user)
	if err != nil {
		return err
	}
	fmt.Println(user.Username)
	fmt.Println(`DeploymentNamespace="` + user.Username + `" and ExecutionStatus="Running"`)
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
			Query:     `DeploymentNamespace="` + user.Username + `" and ExecutionStatus="Running"`,
		})
		if err != nil {
			return err
		}
		result := []*launchPb.LaunchView{}

		for _, workflow := range r.Executions {
			result = append(result, &launchPb.LaunchView{
				Name:           strings.Trim(string(workflow.SearchAttributes.IndexedFields["DeploymentName"].Data), `\"`),
				Namespace:      strings.Trim(string(workflow.SearchAttributes.IndexedFields["DeploymentNamespace"].Data), `\"`),
				WorkflowId:     workflow.GetExecution().GetWorkflowId(),
				RunId:          workflow.GetExecution().GetRunId(),
				WorkloadStatus: strings.Trim(string(workflow.SearchAttributes.IndexedFields["DeploymentStatus"].Data), `\"`),
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
