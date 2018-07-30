// Copyright 2017 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package common

import (
	"fmt"
	"strings"
	"time"

	"go.chromium.org/luci/common/api/swarming/swarming/v1"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/isolatedclient"
	"go.chromium.org/luci/common/logging"
	"golang.org/x/net/context"

	admin "infra/tricium/api/admin/v1"
)

const (
	swarmingBasePath      = "/_ah/api/swarming/v1/"
	swarmingDevServerURL  = "https://chromium-swarm-dev.appspot.com"
	swarmingProdServerURL = "https://chromium-swarm.appspot.com"
)

// SwarmingServer implements the TaskServerAPI for the swarming service.
var SwarmingServer swarmingServer

type swarmingServer struct {
}

// Trigger implements the TaskServerAPI.
func (s swarmingServer) Trigger(c context.Context, serverURL, isolateServerURL string, worker *admin.Worker, workerIsolate, pubsubUserdata string, tags []string) (*TriggerResult, error) {
	pubsubTopic := topic(c)
	// Prepare task dimensions.
	dims := []*swarming.SwarmingRpcsStringPair{}
	for _, d := range worker.Dimensions {
		// Extracting dimension key and value.
		// Note that ':' may appear in the value but not the key.
		dim := strings.SplitN(d, ":", 2)
		if len(dim) != 2 {
			return nil, errors.Reason("failed to split dimension: %q", d).Err()
		}
		dims = append(dims, &swarming.SwarmingRpcsStringPair{Key: dim[0], Value: dim[1]})
	}
	// Prepare CIPD input packages.
	cipd := &swarming.SwarmingRpcsCipdInput{}
	for _, p := range worker.CipdPackages {
		cipd.Packages = append(cipd.Packages, &swarming.SwarmingRpcsCipdPackage{
			PackageName: p.PackageName,
			Path:        p.Path,
			Version:     p.Version,
		})
	}
	// Need to increase the timeout to get a response from the Swarming service.
	c, _ = context.WithTimeout(c, 60*time.Second)
	oauthClient, err := getOAuthClient(c)
	if err != nil {
		return nil, errors.Annotate(err, "failed to create oauth client").Err()
	}
	swarmingService, err := swarming.New(oauthClient)
	if err != nil {
		return nil, errors.Annotate(err, "failed to create swarming client").Err()
	}
	// TODO(emso): Read timeouts from the analyzer config.
	// Prepare properties.
	props := &swarming.SwarmingRpcsTaskProperties{
		Dimensions:           dims,
		ExecutionTimeoutSecs: 600,
		IoTimeoutSecs:        600,
		InputsRef: &swarming.SwarmingRpcsFilesRef{
			Isolated:       workerIsolate,
			Isolatedserver: isolateServerURL,
			Namespace:      isolatedclient.DefaultNamespace,
		},
	}
	// Only include CIPD input if there are packages.
	if len(cipd.Packages) > 0 {
		props.CipdInput = cipd
	}
	swarmingService.BasePath = fmt.Sprintf("%s%s", serverURL, swarmingBasePath)
	res, err := swarmingService.Tasks.New(&swarming.SwarmingRpcsNewTaskRequest{
		Name:           "tricium:" + worker.Name,
		Priority:       100,
		ExpirationSecs: 21600,
		Properties:     props,
		PubsubTopic:    pubsubTopic,
		PubsubUserdata: pubsubUserdata,
		Tags:           tags,
	}).Do()
	if err != nil {
		return nil, errors.Annotate(err, "failed to trigger swarming task").Err()
	}
	logging.Fields{
		"task ID":       res.TaskId,
		"worker":        worker.Name,
		"dimensions":    dims,
		"pubsub topic":  pubsubTopic,
		"input isolate": workerIsolate,
	}.Infof(c, "Worker triggered")
	return &TriggerResult{TaskID: res.TaskId}, nil
}

// Collect implements the TaskServerAPI.
func (s swarmingServer) Collect(c context.Context, serverURL, taskID string, buildID int64) (*CollectResult, error) {
	// Need to increase the timeout to get a response from the Swarming service.
	c, _ = context.WithTimeout(c, 60*time.Second)
	oauthClient, err := getOAuthClient(c)
	if err != nil {
		return nil, errors.Annotate(err, "failed to create oauth client").Err()
	}
	swarmingService, err := swarming.New(oauthClient)
	if err != nil {
		return nil, errors.Annotate(err, "failed to create swarming client").Err()
	}
	swarmingService.BasePath = fmt.Sprintf("%s%s", serverURL, swarmingBasePath)
	taskResult, err := swarmingService.Task.Result(taskID).Do()
	if err != nil {
		return nil, errors.Annotate(err, "failed to retrieve task result from swarming").Err()
	}

	var result *CollectResult
	if taskResult.OutputsRef == nil || taskResult.OutputsRef.Isolated == "" {
		logging.Fields{
			"task ID":    taskID,
			"task state": result.State,
		}.Warningf(c, "Task had no output.")
	} else {
		result.IsolatedOutputHash = taskResult.OutputsRef.Isolated
	}

	if taskResult.State == "COMPLETED" {
		if taskResult.ExitCode != 0 {
			result.State = Failure
		} else {
			result.State = Success
		}
	} else if taskResult.State == "PENDING" || taskResult.State == "RUNNING" {
		result.State = Pending
	}

	return result, nil
}
