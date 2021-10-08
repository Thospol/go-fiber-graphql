package middlewares

import (
	"bytes"
	"encoding/json"
	"fiber-graphql/internal/core/config"
	"fiber-graphql/internal/core/context"
	"fiber-graphql/internal/core/utils"

	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Logger is log request
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		if c.Method() != http.MethodGet && strings.EqualFold(c.Get(fiber.HeaderContentType), "application/json") {
			err := utils.JSONDuplicate(json.NewDecoder(io.NopCloser(bytes.NewBuffer(c.Request().Body()))), nil)
			if err != nil {
				return c.
					Status(config.RR.JSONDuplicateOrInvalidFormat.HTTPStatusCode()).
					JSON(config.RR.JSONDuplicateOrInvalidFormat.WithLocale(c))
			}
		}

		err := c.Next()
		if err != nil {
			return err
		}

		// write response log
		logs := logrus.Fields{
			"host":         c.Hostname(),
			"method":       c.Method(),
			"path":         c.OriginalURL(),
			"language":     c.Locals(context.LangKey),
			"ip":           c.IP(),
			"user_agent":   c.Get("User-Agent"),
			"prefix_path":  c.Get("Prefix-Path"),
			"body_size":    fmt.Sprintf("%.5f MB", float64(bytes.NewReader(c.Request().Body()).Len())/1024.00/1024.00),
			"status_code":  fmt.Sprintf("%d %s", c.Response().StatusCode(), http.StatusText(c.Response().StatusCode())),
			"process_time": int64(time.Since(start)),
			"platform":     c.Get("Platform"),
		}

		userID := context.New(c).GetUserID()
		if userID > 0 {
			logs["user_id"] = userID
		}

		parameters := c.Locals(context.ParametersKey)
		if parameters != nil {
			b, _ := json.Marshal(parameters)
			for _, f := range []string{"password"} {
				if res := gjson.GetBytes(b, f); res.Exists() {
					b, _ = sjson.SetBytes(b, f, "**********")
				}
			}
			logs["parameters"] = string(b)
		} else {
			logs["parameters"] = "{}"
		}
		logrus.WithFields(logs).Infof("[%s][%s] response: %v", c.Method(), c.OriginalURL(), string(c.Response().Body()))

		return nil
	}
}
