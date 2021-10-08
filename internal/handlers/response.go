package handlers

import (
	"fiber-graphql/internal/core/context"
	"fiber-graphql/internal/handlers/render"
	"fiber-graphql/internal/models"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// ResponseObject handle response object
func ResponseObject(c *fiber.Ctx, fn interface{}, request interface{}) error {
	ctx := context.New(c)
	err := ctx.BindValue(request, true)
	if err != nil {
		logrus.Errorf("bind value error: %s", err)
		return render.Error(c, err)
	}

	out := reflect.ValueOf(fn).Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(request),
	})
	errObj := out[1].Interface()
	if errObj != nil {
		logrus.Errorf("call service error: %s", errObj)
		return errObj.(error)
	}

	return render.JSON(c, out[0].Interface())
}

// ResponseObjectWithoutRequest handle response object without request
func ResponseObjectWithoutRequest(c *fiber.Ctx, fn interface{}) error {
	out := reflect.ValueOf(fn).Call([]reflect.Value{
		reflect.ValueOf(context.New(c)),
	})
	errObj := out[1].Interface()
	if errObj != nil {
		logrus.Errorf("call service error: %s", errObj)
		return errObj.(error)
	}

	return render.JSON(c, out[0].Interface())
}

// ResponseSuccess handle response success
func ResponseSuccess(c *fiber.Ctx, fn interface{}, request interface{}) error {
	ctx := context.New(c)
	err := ctx.BindValue(request, true)
	if err != nil {
		logrus.Errorf("bind value error: %s", err)
		return render.Error(c, err)
	}

	out := reflect.ValueOf(fn).Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(request),
	})
	errObj := out[0].Interface()
	if errObj != nil {
		logrus.Errorf("call service error: %s", errObj)
		return errObj.(error)
	}
	return render.JSON(c, models.NewSuccessMessage())
}

// ResponseSuccessWithoutRequest handle response success without request
func ResponseSuccessWithoutRequest(c *fiber.Ctx, fn interface{}) error {
	ctx := context.New(c)
	out := reflect.ValueOf(fn).Call([]reflect.Value{
		reflect.ValueOf(ctx),
	})
	errObj := out[0].Interface()
	if errObj != nil {
		logrus.Errorf("call service error: %s", errObj)
		return errObj.(error)
	}
	return render.JSON(c, models.NewSuccessMessage())
}
