package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"opt-switch/app/device/service"
	"opt-switch/app/device/service/dto"
)

// CommandAPI handles command execution HTTP requests
type CommandAPI struct {
	api.Api
}

// ExecuteCommand executes a single command
// @Summary Execute a single command on the device
// @Description Executes a single CLI command on the device and returns the output
// @Tags device
// @Accept json
// @Produce json
// @Param request body dto.CommandExecuteReq true "Command execution request"
// @Success 200 {object} response.Response{data=dto.CommandExecuteResp}
// @Failure 400 {object} response.Response
// @Failure 429 {object} response.Response "Service busy, please try again later"
// @Failure 500 {object} response.Response
// @Router /api/v1/device/command/execute [post]
func (e *CommandAPI) ExecuteCommand(c *gin.Context) {
	req := dto.CommandExecuteReq{}
	s := service.CommandService{}
	err := e.MakeContext(c).
		MakeService(&s.Service).
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	resp, err := s.ExecuteCommand(c, &req)
	if err != nil {
		statusCode, msg := s.MapError(err)
		e.Error(statusCode, err, msg)
		return
	}

	e.OK(resp, "Command executed successfully")
}

// ExecuteBatch executes multiple commands
// @Summary Execute multiple commands on the device
// @Description Executes multiple CLI commands sequentially on the device
// @Tags device
// @Accept json
// @Produce json
// @Param request body dto.BatchCommandReq true "Batch command request"
// @Success 200 {object} response.Response{data=dto.BatchCommandResp}
// @Failure 400 {object} response.Response
// @Failure 429 {object} response.Response "Service busy, please try again later"
// @Failure 500 {object} response.Response
// @Router /api/v1/device/command/batch [post]
func (e *CommandAPI) ExecuteBatch(c *gin.Context) {
	req := dto.BatchCommandReq{}
	s := service.CommandService{}
	err := e.MakeContext(c).
		MakeService(&s.Service).
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	resp, err := s.ExecuteBatch(c, &req)
	if err != nil {
		statusCode, msg := s.MapError(err)
		e.Error(statusCode, err, msg)
		return
	}

	e.OK(resp, "Commands executed successfully")
}

// GetHistory retrieves command execution history
// @Summary Get command execution history
// @Description Retrieves historical command execution records from log file
// @Tags device
// @Accept json
// @Produce json
// @Param limit query int true "Limit number of records" minimum(1) maximum(1000)
// @Param offset query int true "Offset for pagination" minimum(0)
// @Success 200 {object} response.Response{data=dto.CommandHistoryResp}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/device/command/history [get]
func (e *CommandAPI) GetHistory(c *gin.Context) {
	req := dto.CommandHistoryReq{}
	s := service.CommandService{}
	err := e.MakeContext(c).
		MakeService(&s.Service).
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(400, err, err.Error())
		return
	}

	resp, err := s.GetHistory(&req)
	if err != nil {
		e.Error(500, err, "Failed to get history")
		return
	}

	e.OK(resp, "History retrieved successfully")
}

// GetStatus returns the device connection status
// @Summary Get device connection status
// @Description Returns the current status of device connections and queue
// @Tags device
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=dto.DeviceStatusResp}
// @Failure 500 {object} response.Response
// @Router /api/v1/device/status [get]
func (e *CommandAPI) GetStatus(c *gin.Context) {
	s := service.CommandService{}
	err := e.MakeContext(c).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	resp := s.GetStatus()
	e.OK(resp, "Status retrieved successfully")
}

// GetDeviceInfo is a simple health check endpoint
// @Summary Get device information
// @Description Returns basic device information
// @Tags device
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/device [get]
func (e *CommandAPI) GetDeviceInfo(c *gin.Context) {
	e.OK(gin.H{
		"status": "online",
		"type":   "switch",
	}, "Device is online")
}
