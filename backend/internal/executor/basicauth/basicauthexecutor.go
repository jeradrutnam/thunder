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

// Package basicauth provides the basic authentication executor for handling username and password authentication.
package basicauth

import (
	"encoding/json"

	authnmodel "github.com/asgardeo/thunder/internal/authn/model"
	flowconst "github.com/asgardeo/thunder/internal/flow/constants"
	flowmodel "github.com/asgardeo/thunder/internal/flow/model"
	"github.com/asgardeo/thunder/internal/system/log"
	userprovider "github.com/asgardeo/thunder/internal/user/provider"
)

const loggerComponentName = "BasicAuthExecutor"

// BasicAuthExecutor implements the ExecutorInterface for basic authentication.
type BasicAuthExecutor struct {
	internal flowmodel.Executor
}

// NewBasicAuthExecutor creates a new instance of BasicAuthExecutor.
func NewBasicAuthExecutor(id, name string) flowmodel.ExecutorInterface {
	return &BasicAuthExecutor{
		internal: flowmodel.Executor{
			Properties: flowmodel.ExecutorProperties{
				ID:   id,
				Name: name,
			},
		},
	}
}

// GetID returns the ID of the BasicAuthExecutor.
func (b *BasicAuthExecutor) GetID() string {
	return b.internal.GetID()
}

// GetName returns the name of the BasicAuthExecutor.
func (b *BasicAuthExecutor) GetName() string {
	return b.internal.GetName()
}

// GetProperties returns the properties of the BasicAuthExecutor.
func (b *BasicAuthExecutor) GetProperties() flowmodel.ExecutorProperties {
	return b.internal.Properties
}

// Execute executes the basic authentication logic.
func (b *BasicAuthExecutor) Execute(ctx *flowmodel.NodeContext) (*flowmodel.ExecutorResponse, error) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName),
		log.String(log.LoggerKeyExecutorID, b.GetID()),
		log.String(log.LoggerKeyFlowID, ctx.FlowID))
	logger.Debug("Executing basic authentication executor")

	execResp := &flowmodel.ExecutorResponse{}

	// Validate for the required input data.
	if b.requiredInputData(ctx, execResp) {
		// If required input data is not provided, return incomplete status.
		logger.Debug("Required input data for basic authentication executor is not provided")
		execResp.Status = flowconst.ExecUserInputRequired
		return execResp, nil
	}

	username := ctx.UserInputData["username"]
	// TODO: Should handle client errors here. Service should return a ServiceError and
	//  client errors should be appended as a failure.
	//  For the moment handling returned error as a authentication failure.
	authenticatedUser, err := getAuthenticatedUser(username, ctx.UserInputData["password"], logger)
	if err != nil {
		execResp.Status = flowconst.ExecFailure
		execResp.FailureReason = "Failed to authenticate user: " + err.Error()
		return execResp, nil
	}
	if authenticatedUser == nil {
		execResp.Status = flowconst.ExecFailure
		execResp.FailureReason = "Authenticated user not found."
		return execResp, nil
	}
	if !authenticatedUser.IsAuthenticated {
		execResp.Status = flowconst.ExecFailure
		execResp.FailureReason = "User authentication failed."
		return execResp, nil
	}

	ctx.AuthenticatedUser = *authenticatedUser
	execResp.Status = flowconst.ExecComplete

	logger.Debug("Basic authentication executor execution completed",
		log.String("status", string(execResp.Status)),
		log.Bool("isAuthenticated", ctx.AuthenticatedUser.IsAuthenticated))

	return execResp, nil
}

// requiredInputData checks and adds the required input data for basic authentication.
// Returns true if needed to request user input data.
func (b *BasicAuthExecutor) requiredInputData(ctx *flowmodel.NodeContext, execResp *flowmodel.ExecutorResponse) bool {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName),
		log.String(log.LoggerKeyExecutorID, b.GetID()),
		log.String(log.LoggerKeyFlowID, ctx.FlowID))

	// TODO: Convert password to a secure type (i.e. byte_array)
	basicReqData := []flowmodel.InputData{
		{
			Name:     "username",
			Type:     "string",
			Required: true,
		},
		{
			Name:     "password",
			Type:     "string",
			Required: true,
		},
	}

	// Check for the required input data. Also appends the authenticator specific input data.
	// TODO: This validation should be moved to the flow composer. Ideally the validation and appending
	//  should happen during the flow definition creation.
	requiredData := ctx.NodeInputData
	if len(requiredData) == 0 {
		logger.Debug("No required input data defined for basic authentication executor")
		requiredData = basicReqData
	} else {
		// Append the default required data if not already present.
		for _, inputData := range basicReqData {
			exists := false
			for _, existingInputData := range requiredData {
				if existingInputData.Name == inputData.Name {
					exists = true
					break
				}
			}
			// If the input data already exists, skip adding it again.
			if !exists {
				requiredData = append(requiredData, inputData)
			}
		}
	}

	requireData := false

	if execResp.RequiredData == nil {
		execResp.RequiredData = make([]flowmodel.InputData, 0)
	}

	if len(ctx.UserInputData) == 0 {
		execResp.RequiredData = append(execResp.RequiredData, requiredData...)
		return true
	}

	// Check if the required input data is provided by the user.
	for _, inputData := range requiredData {
		if _, ok := ctx.UserInputData[inputData.Name]; !ok {
			if !inputData.Required {
				logger.Debug("Skipping optional input data that is not provided by user",
					log.String("inputDataName", inputData.Name))
				continue
			}
			execResp.RequiredData = append(execResp.RequiredData, inputData)
			requireData = true
			logger.Debug("Required input data not provided by user",
				log.String("inputDataName", inputData.Name))
		}
	}

	return requireData
}

// getAuthenticatedUser perform authentication based on the provided username and password and return
// authenticated user details.
func getAuthenticatedUser(username, password string, logger *log.Logger) (*authnmodel.AuthenticatedUser, error) {
	userProvider := userprovider.NewUserProvider()
	userService := userProvider.GetUserService()

	userID, err := userService.IdentityUser("username", username)
	if err != nil {
		logger.Error("Failed to identify user by username",
			log.String("username", log.MaskString(username)),
			log.Error(err))
		return nil, err
	}
	if *userID == "" {
		logger.Error("User not found for the provided username",
			log.String("username", log.MaskString(username)))
		return nil, err
	}

	user, err := userService.VerifyUser(*userID, "password", password)
	if err != nil {
		logger.Error("Failed to verify user credentials", log.String("userID", *userID), log.Error(err))
		return nil, err
	}

	var authenticatedUser authnmodel.AuthenticatedUser
	if user == nil {
		authenticatedUser = authnmodel.AuthenticatedUser{
			IsAuthenticated: false,
		}
	} else {
		var attrs map[string]interface{}
		if err := json.Unmarshal(user.Attributes, &attrs); err != nil {
			logger.Error("Failed to unmarshal user attributes", log.Error(err))
			return nil, err
		}
		authenticatedUser = authnmodel.AuthenticatedUser{
			IsAuthenticated:        true,
			UserID:                 user.ID,
			Username:               attrs["username"].(string),
			AuthenticatedSubjectID: attrs["email"].(string),
			Attributes: map[string]string{
				"email":     attrs["email"].(string),
				"firstName": attrs["firstName"].(string),
				"lastName":  attrs["lastName"].(string),
			},
		}
	}
	return &authenticatedUser, nil
}
