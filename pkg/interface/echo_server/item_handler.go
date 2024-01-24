package echoserver

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/taroxii/vote-api/pkg/constants"
	"github.com/taroxii/vote-api/pkg/entity"
)

type ItemHandler struct {
	itemUsecase entity.ItemUseCase
}

func NewItemHandler(e *echo.Echo, iuc entity.ItemUseCase) {
	handler := &ItemHandler{
		itemUsecase: iuc,
	}
	mdw := InitMiddleware()
	itemsRoute := e.Group("/items")
	itemsRoute.Use(mdw.AuthMiddleware)
	itemsRoute.GET("", handler.FetchItems)
	itemsRoute.POST("", handler.CreateItem)
	itemsRoute.PATCH("/:id", handler.Update)
	itemsRoute.DELETE("/:id", handler.Delete)
	itemsRoute.PATCH("/:id/vote", handler.VoteItem)
	itemsRoute.PATCH("/:id/clear", handler.ClearVote)
}

func (i *ItemHandler) FetchItems(c echo.Context) error {
	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()
	listAr, nextCursor, err := i.itemUsecase.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listAr)
}

func (i *ItemHandler) Delete(c echo.Context) error {
	idParams := c.Param("id")
	id, err := strconv.Atoi(idParams)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	err = i.itemUsecase.Delete(ctx, int64(id))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, ResponseError{Message: "the record has been deleted"})
}

func (i *ItemHandler) Update(c echo.Context) error {
	idParams := c.Param("id")
	id, err := strconv.Atoi(idParams)
	itemReq := new(entity.Item)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()

	if err = c.Bind(itemReq); err != nil {
		return c.JSON(http.StatusBadGateway, ResponseError{Message: err.Error()})
	}
	itemReq.ID = int64(id)

	err = i.itemUsecase.Update(ctx, itemReq)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, itemReq)
}

func (i *ItemHandler) VoteItem(c echo.Context) error {
	idParams := c.Param("id")
	id, err := strconv.Atoi(idParams)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	usr := c.Get(constants.USER_CONTEXT_KEY).(*entity.JWTClaims)
	ite, err := i.itemUsecase.Vote(ctx, int64(id), usr.ID)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, ite)
}

func (i *ItemHandler) ClearVote(c echo.Context) error {
	idParams := c.Param("id")
	id, err := strconv.Atoi(idParams)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()

	ite, err := i.itemUsecase.ClearVote(ctx, int64(id))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, ite)
}

func (i *ItemHandler) CreateItem(c echo.Context) error {
	req := new(entity.Item)
	err := c.Bind(req)
	ctx := c.Request().Context()
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	err = i.itemUsecase.Insert(ctx, req)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, req)
}
