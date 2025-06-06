/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package utils provides utility functions for flow processing.
package utils

import (
	"fmt"

	"github.com/asgardeo/thunder/internal/executor/authassert"
	"github.com/asgardeo/thunder/internal/executor/basicauth"
	"github.com/asgardeo/thunder/internal/executor/githubauth"
	"github.com/asgardeo/thunder/internal/executor/googleauth"
	"github.com/asgardeo/thunder/internal/flow/constants"
	"github.com/asgardeo/thunder/internal/flow/jsonmodel"
	"github.com/asgardeo/thunder/internal/flow/model"
	idpmodel "github.com/asgardeo/thunder/internal/idp/model"
	idpservice "github.com/asgardeo/thunder/internal/idp/service"
)

// BuildGraphFromDefinition builds a graph from a graph definition json.
func BuildGraphFromDefinition(definition *jsonmodel.GraphDefinition) (model.GraphInterface, error) {
	if definition == nil || len(definition.Nodes) == 0 {
		return nil, fmt.Errorf("graph definition is nil or has no nodes")
	}

	// Create a graph
	g := model.NewGraph(definition.ID)

	// Map to track which nodes have incoming edges
	hasIncomingEdge := make(map[string]bool)

	// First, mark all nodes that have incoming edges
	for _, targetIDs := range definition.Edges {
		for _, targetID := range targetIDs {
			hasIncomingEdge[targetID] = true
		}
	}

	// Find the start node (node without incoming edges)
	startNodeID := ""
	for _, node := range definition.Nodes {
		if !hasIncomingEdge[node.ID] {
			startNodeID = node.ID
			break
		}
	}

	// If no start node found, fallback to the first node
	if startNodeID == "" && len(definition.Nodes) > 0 {
		startNodeID = definition.Nodes[0].ID
	}

	// Validate that we have a valid start node
	if startNodeID == "" {
		return nil, fmt.Errorf("no valid start node found in the graph definition")
	}

	// Add all nodes to the graph
	for _, nodeDef := range definition.Nodes {
		isStartNode := (nodeDef.ID == startNodeID)
		isFinalNode := (nodeDef.Type == string(constants.NodeTypeAuthSuccess))

		// Construct a new node
		node := model.NewNode(nodeDef.ID, nodeDef.Type, isStartNode, isFinalNode)

		// Convert and set input data from definition
		inputData := make([]model.InputData, len(nodeDef.InputData))
		for i, input := range nodeDef.InputData {
			inputData[i] = model.InputData{
				Name:     input.Name,
				Type:     input.Type,
				Required: input.Required,
			}
		}
		node.SetInputData(inputData)

		// Set the executor config if defined
		if nodeDef.Executor.Name != "" {
			executor, err := getExecutorConfigByName(nodeDef.Executor)
			if err != nil {
				return nil, fmt.Errorf("error while getting executor %s: %w", nodeDef.Executor, err)
			}
			node.SetExecutorConfig(executor)
		} else if nodeDef.Type == string(constants.NodeTypeAuthSuccess) {
			executor, err := getExecutorConfigByName(jsonmodel.ExecutorDefinition{
				Name: "AuthAssertExecutor",
			})
			if err != nil {
				return nil, fmt.Errorf("error while getting default AuthAssertExecutor: %w", err)
			}
			node.SetExecutorConfig(executor)
		}

		err := g.AddNode(node)
		if err != nil {
			return nil, fmt.Errorf("failed to add node %s to the graph: %w", nodeDef.ID, err)
		}
	}

	err := g.SetStartNodeID(startNodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to set start node ID: %w", err)
	}

	// Add all edges to the graph
	for sourceID, targetIDs := range definition.Edges {
		for _, targetID := range targetIDs {
			err := g.AddEdge(sourceID, targetID)
			if err != nil {
				return nil, fmt.Errorf("failed to add edge from %s to %s: %w", sourceID, targetID, err)
			}
		}
	}

	// Set PreviousNodeID and NextNodeID based on edges
	for fromNodeID, toNodeIDs := range g.GetEdges() {
		if len(toNodeIDs) > 0 {
			// Set the NextNodeID for the source node
			if sourceNode, exists := g.GetNode(fromNodeID); exists {
				sourceNode.SetNextNodeID(toNodeIDs[0])
				// Update the source node in the graph
				err := g.AddNode(sourceNode)
				if err != nil {
					return nil, fmt.Errorf("failed to update source node %s in the graph: %w", fromNodeID, err)
				}
			}

			// Set the PreviousNodeID for each target node
			for _, toNodeID := range toNodeIDs {
				if targetNode, exists := g.GetNode(toNodeID); exists {
					targetNode.SetPreviousNodeID(fromNodeID)
					// Update the target node in the graph
					err := g.AddNode(targetNode)
					if err != nil {
						return nil, fmt.Errorf("failed to update target node %s in the graph: %w", toNodeID, err)
					}
				}
			}
		}
	}

	return g, nil
}

// getExecutorConfigByName constructs an executor configuration by its definition if it exists.
func getExecutorConfigByName(execDef jsonmodel.ExecutorDefinition) (*model.ExecutorConfig, error) {
	if execDef.Name == "" {
		return nil, fmt.Errorf("executor name cannot be empty")
	}

	// At this point, we assume executors and attached IDPs are already registered in the system.
	// Hence validations will not be done at this point.
	var executor model.ExecutorConfig
	switch execDef.Name {
	case "BasicAuthExecutor":
		executor = model.ExecutorConfig{
			Name:    "BasicAuthExecutor",
			IdpName: "Local",
		}
	case "GithubOAuthExecutor":
		executor = model.ExecutorConfig{
			Name:    "GithubOAuthExecutor",
			IdpName: execDef.IdpName,
		}
	case "GoogleOIDCAuthExecutor":
		executor = model.ExecutorConfig{
			Name:    "GoogleOIDCAuthExecutor",
			IdpName: execDef.IdpName,
		}
	case "AuthAssertExecutor":
		executor = model.ExecutorConfig{
			Name: "AuthAssertExecutor",
		}
	default:
		return nil, fmt.Errorf("executor with name %s not found", execDef.Name)
	}

	if executor.Name == "" {
		return nil, fmt.Errorf("executor with name %s could not be created", execDef.Name)
	}

	return &executor, nil
}

// GetExecutorByName constructs an executor by its definition.
func GetExecutorByName(execConfig *model.ExecutorConfig) (model.ExecutorInterface, error) {
	if execConfig == nil {
		return nil, fmt.Errorf("executor configuration cannot be nil")
	}
	if execConfig.Name == "" {
		return nil, fmt.Errorf("executor name cannot be empty")
	}

	var executor model.ExecutorInterface
	switch execConfig.Name {
	case "BasicAuthExecutor":
		idp, err := getIDP("Local")
		if err != nil {
			return nil, fmt.Errorf("error while getting IDP for BasicAuthExecutor: %w", err)
		}
		executor = basicauth.NewBasicAuthExecutor(idp.ID, idp.Name)
	case "GithubOAuthExecutor":
		idp, err := getIDP(execConfig.IdpName)
		if err != nil {
			return nil, fmt.Errorf("error while getting IDP for GithubOAuthExecutor: %w", err)
		}
		executor = githubauth.NewGithubOAuthExecutor(idp.ID, idp.Name, idp.ClientID, idp.ClientSecret,
			idp.RedirectURI, idp.Scopes, map[string]string{})
	case "GoogleOIDCAuthExecutor":
		idp, err := getIDP(execConfig.IdpName)
		if err != nil {
			return nil, fmt.Errorf("error while getting IDP for GoogleOIDCAuthExecutor: %w", err)
		}
		executor = googleauth.NewGoogleOIDCAuthExecutor(idp.ID, idp.Name, idp.ClientID, idp.ClientSecret,
			idp.RedirectURI, idp.Scopes, map[string]string{})
	case "AuthAssertExecutor":
		executor = authassert.NewAuthAssertExecutor("auth-assert-executor", "AuthAssertExecutor")
	default:
		return nil, fmt.Errorf("executor with name %s not found", execConfig.Name)
	}

	if executor == nil {
		return nil, fmt.Errorf("executor with name %s could not be created", execConfig.Name)
	}
	return executor, nil
}

// getIDP retrieves the IDP by its name. Returns an error if the IDP does not exist or if the name is empty.
func getIDP(idpName string) (*idpmodel.IDP, error) {
	if idpName == "" {
		return nil, fmt.Errorf("IDP name cannot be empty")
	}

	idpSvc := idpservice.GetIDPService()
	idp, err := idpSvc.GetIdentityProviderByName(idpName)
	if err != nil {
		return nil, fmt.Errorf("error while getting IDP with the name %s: %w", idpName, err)
	}
	if idp == nil {
		return nil, fmt.Errorf("IDP with name %s does not exist", idpName)
	}

	return idp, nil
}
