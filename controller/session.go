package controller

import (
	"fp-rpl/common"
	"fp-rpl/dto"
	"fp-rpl/entity"
	"fp-rpl/service"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type sessionController struct {
	sessionService service.SessionService
	areaService service.AreaService
	// filmService service.FilmService
}

type SessionController interface {
	CreateSession(ctx *gin.Context)
	GetAllSessions(ctx *gin.Context)
}

func NewSessionController(sessionS service.SessionService, areaS service.AreaService) SessionController {
	return &sessionController{
		sessionService: sessionS,
		areaService: areaS,
	}
}

func (sessionC *sessionController) CreateSession(ctx *gin.Context) {
	var sessionDTO dto.SessionCreateRequest
	err := ctx.ShouldBind(&sessionDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to process session create request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	// Check for duplicate Sessionname or Email
	sessionCheck, err := sessionC.sessionService.GetSessionByTimeAndPlace(ctx, sessionDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to process session create request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	// Check if duplicate is found
	if !(reflect.DeepEqual(sessionCheck, entity.Session{})) {
		resp := common.CreateFailResponse("session with the exact same attributes already exists", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	// Check Film by ID

	// Check Area by ID
	area, err := sessionC.areaService.GetAreaByID(ctx, sessionDTO.AreaID)
	if err != nil {
		resp := common.CreateFailResponse(err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if reflect.DeepEqual(area, entity.Area{}) {
		resp := common.CreateFailResponse("area not found", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	newSession, err := sessionC.sessionService.CreateNewSession(ctx, sessionDTO, area.SpotCount, area.SpotPerRow)
	if err != nil {
		resp := common.CreateFailResponse("failed to process session create request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	resp := common.CreateSuccessResponse("successfully created session", http.StatusCreated, newSession)
	ctx.JSON(http.StatusCreated, resp)
}

func (sessionC *sessionController) GetAllSessions(ctx *gin.Context) {
	sessions, err := sessionC.sessionService.GetAllSessions(ctx)
	if err != nil {
		resp := common.CreateFailResponse("failed to fetch all sessions", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp common.Response
	if len(sessions) == 0 {
		resp = common.CreateSuccessResponse("no session found", http.StatusOK, sessions)
	} else {
		resp = common.CreateSuccessResponse("successfully fetched all sessions", http.StatusOK, sessions)
	}
	ctx.JSON(http.StatusOK, resp)
}