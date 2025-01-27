package controlplaneclient

import (
	"context"
	"encoding/json"
	"io"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/cloud"
	cloudtestworkflow "github.com/kubeshop/testkube/pkg/cloud/data/testworkflow"
	"github.com/kubeshop/testkube/pkg/repository/channels"
)

type RunnerRequestsWatcher channels.Watcher[RunnerRequest]
type NotificationWatcher channels.Watcher[*testkube.TestWorkflowExecutionNotification]

type RunnerClient interface {
	GetRunnerOngoingExecutions(ctx context.Context) ([]*cloud.UnfinishedExecution, error)
	WatchRunnerRequests(ctx context.Context) RunnerRequestsWatcher
	ProcessExecutionNotificationRequests(ctx context.Context, process func(ctx context.Context, req *cloud.TestWorkflowNotificationsRequest) NotificationWatcher) error
	ProcessExecutionParallelWorkerNotificationRequests(ctx context.Context, process func(ctx context.Context, req *cloud.TestWorkflowParallelStepNotificationsRequest) NotificationWatcher) error
	ProcessExecutionServiceNotificationRequests(ctx context.Context, process func(ctx context.Context, req *cloud.TestWorkflowServiceNotificationsRequest) NotificationWatcher) error
}

func (c *client) GetRunnerOngoingExecutions(ctx context.Context) ([]*cloud.UnfinishedExecution, error) {
	if c.IsLegacy() {
		return c.legacyGetRunnerOngoingExecutions(ctx)
	}
	res, err := call(ctx, c.metadata().GRPC(), c.client.GetUnfinishedExecutions, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	result := make([]*cloud.UnfinishedExecution, 0)
	for {
		// Take the context error if possible
		if err == nil && ctx.Err() != nil {
			err = ctx.Err()
		}

		// End when it's done
		if errors.Is(err, io.EOF) {
			return result, nil
		}

		// Handle the error
		if err != nil {
			return nil, err
		}

		// Get the next execution to monitor
		var exec *cloud.UnfinishedExecution
		exec, err = res.Recv()
		if err != nil {
			continue
		}

		result = append(result, exec)
	}
}

// Deprecated
func (c *client) legacyGetRunnerOngoingExecutions(ctx context.Context) ([]*cloud.UnfinishedExecution, error) {
	jsonPayload, err := json.Marshal(cloudtestworkflow.ExecutionGetRunningRequest{})
	if err != nil {
		return nil, err
	}
	s := structpb.Struct{}
	if err := s.UnmarshalJSON(jsonPayload); err != nil {
		return nil, err
	}
	req := cloud.CommandRequest{
		Command: string(cloudtestworkflow.CmdTestWorkflowExecutionGetRunning),
		Payload: &s,
	}
	cmdResponse, err := call(ctx, c.metadata().GRPC(), c.client.Call, &req)
	if err != nil {
		return nil, err
	}
	var response cloudtestworkflow.ExecutionGetRunningResponse
	err = json.Unmarshal(cmdResponse.Response, &response)
	if err != nil {
		return nil, err
	}
	result := make([]*cloud.UnfinishedExecution, 0)
	for i := range response.WorkflowExecutions {
		// Ignore if it's not assigned to any runner
		if response.WorkflowExecutions[i].RunnerId == "" && len(response.WorkflowExecutions[i].Signature) == 0 {
			continue
		}

		// Ignore if it's assigned to a different runner
		if response.WorkflowExecutions[i].RunnerId != c.proContext.AgentID {
			continue
		}

		result = append(result, &cloud.UnfinishedExecution{EnvironmentId: c.proContext.EnvID, Id: response.WorkflowExecutions[i].Id})
	}
	return result, err
}

func (c *client) WatchRunnerRequests(ctx context.Context) RunnerRequestsWatcher {
	stream, err := watch(ctx, c.metadata().GRPC(), c.client.GetRunnerRequests)
	if err != nil {
		return channels.NewError[RunnerRequest](err)
	}
	watcher := channels.NewWatcher[RunnerRequest]()
	sendMu := sync.Mutex{}
	send := func(v *cloud.RunnerResponse) error {
		sendMu.Lock()
		defer sendMu.Unlock()
		return stream.Send(v)
	}
	go func() {
		defer watcher.Close(err)
		for {
			// Ignore if it's not implemented in the Control Plane
			if getGrpcErrorCode(err) == codes.Unimplemented {
				return
			}

			// Take the context error if possible
			if errors.Is(err, context.Canceled) && ctx.Err() != nil {
				err = context.Cause(ctx)
			}

			// Handle the error
			if err != nil {
				return
			}

			// Get the next runner request
			var req *cloud.RunnerRequest
			req, err = stream.Recv()
			if err != nil {
				continue
			}

			if req.Type == cloud.RunnerRequestType_PING {
				err = send(&cloud.RunnerResponse{Type: cloud.RunnerRequestType_PING})
				if err != nil {
					return
				}
				continue
			}

			watcher.Send(&runnerRequestData{data: req, send: send})
		}
	}()
	return watcher
}

func (c *client) ProcessExecutionNotificationRequests(ctx context.Context, process func(ctx context.Context, req *cloud.TestWorkflowNotificationsRequest) NotificationWatcher) error {
	return processNotifications(
		ctx,
		c.metadata().GRPC(),
		c.client.GetTestWorkflowNotificationsStream,
		buildPongNotification,
		buildCloudNotification,
		buildCloudError,
		process,
	)
}

func (c *client) ProcessExecutionParallelWorkerNotificationRequests(ctx context.Context, process func(ctx context.Context, req *cloud.TestWorkflowParallelStepNotificationsRequest) NotificationWatcher) error {
	return processNotifications(
		ctx,
		c.metadata().GRPC(),
		c.client.GetTestWorkflowParallelStepNotificationsStream,
		buildParallelStepPongNotification,
		buildParallelStepCloudNotification,
		buildParallelStepCloudError,
		process,
	)
}

func (c *client) ProcessExecutionServiceNotificationRequests(ctx context.Context, process func(ctx context.Context, req *cloud.TestWorkflowServiceNotificationsRequest) NotificationWatcher) error {
	return processNotifications(
		ctx,
		c.metadata().GRPC(),
		c.client.GetTestWorkflowServiceNotificationsStream,
		buildServicePongNotification,
		buildServiceCloudNotification,
		buildServiceCloudError,
		process,
	)
}
