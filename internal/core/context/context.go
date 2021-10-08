package context

import (
	"fiber-graphql/internal/core/config"
	"fiber-graphql/internal/core/sql"
	"fiber-graphql/internal/core/utils"
	"net/http"
	"reflect"
	"strconv"

	"github.com/valyala/fasthttp"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

const (
	pathKey            = "path"
	compositeFormDepth = 3
	// UserKey user key
	UserKey = "user"
	// LangKey lang key
	LangKey = "lang"
	// PostgreDatabaseKey database `postgre` key
	PostgreDatabaseKey = "postgre_database"
	// ParametersKey parameters key
	ParametersKey = "parameters"
	// UsernameKey username key
	UsernameKey = "username"
)

// Context context
type Context struct {
	*fiber.Ctx
}

// New new custom fiber context
func New(c *fiber.Ctx) *Context {
	return &Context{Ctx: c}
}

// BindValue bind value
func (c *Context) BindValue(i interface{}, validate bool) error {
	switch c.Method() {
	case http.MethodGet:
		_ = c.QueryParser(i)

	default:
		_ = c.BodyParser(i)
	}

	c.PathParser(i, 1)
	c.Locals(ParametersKey, i)
	utils.TrimSpace(i, 1)

	if validate {
		err := c.Validate(i)
		if err != nil {
			return err
		}
	}
	return nil
}

// PathParser parse path param
func (c *Context) PathParser(i interface{}, depth int) {
	formValue := reflect.ValueOf(i)
	if formValue.Kind() == reflect.Ptr {
		formValue = formValue.Elem()
	}
	t := reflect.TypeOf(formValue.Interface())
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		paramValue := formValue.FieldByName(fieldName)
		if paramValue.IsValid() {
			if depth < compositeFormDepth && paramValue.Kind() == reflect.Struct {
				depth++
				c.PathParser(paramValue.Addr().Interface(), depth)
			}
			tag := t.Field(i).Tag.Get(pathKey)
			if tag != "" {
				setValue(paramValue, c.Params(tag))
			}
		}
	}
}

func setValue(paramValue reflect.Value, value string) {
	if paramValue.IsValid() && value != "" {
		switch paramValue.Kind() {
		case reflect.Uint:
			number, _ := strconv.ParseUint(value, 10, 32)
			paramValue.SetUint(number)

		case reflect.String:
			paramValue.SetString(value)

		default:
			number, err := strconv.Atoi(value)
			if err != nil {
				paramValue.SetString(value)
			} else {
				paramValue.SetInt(int64(number))
			}
		}
	}
}

// Validate validate
func (c *Context) Validate(i interface{}) error {
	if err := config.CF.Validator.Struct(i); err != nil {
		return config.RR.CustomMessage(err.Error(), err.Error()).WithLocale(c.Ctx)
	}

	return nil
}

// GetDatabase get connection database `postgresql`
func (c *Context) GetDatabase() *gorm.DB {
	val := c.Locals(PostgreDatabaseKey)
	if val == nil {
		return sql.Database
	}

	return val.(*gorm.DB)
}

// Claims jwt claims
type Claims struct {
	jwt.StandardClaims
	Role int
}

// GetClaims get user claims
func (c *Context) GetClaims() *Claims {
	user := c.Locals(UserKey).(*jwt.Token)
	return user.Claims.(*Claims)
}

// GetUserID get user claims
func (c *Context) GetUserID() uint {
	token, ok := c.fiberCtx().Locals(UserKey).(*jwt.Token)
	if ok {
		if cl := token.Claims.(*Claims); cl != nil {
			subject := c.GetClaims().Subject
			i, _ := strconv.Atoi(subject)
			return uint(i)
		}
	}

	return 0
}

// GetRole get role
func (c *Context) GetRole() int {
	token, ok := c.fiberCtx().Locals(UserKey).(*jwt.Token)
	if ok {
		cl := token.Claims.(*Claims)
		if cl != nil {
			return c.GetClaims().Role
		}
	}

	return 0
}

// GetLanguage get language
func (c *Context) GetLanguage() string {
	return c.fiberCtx().Locals(LangKey).(string)
}

func (c *Context) fiberCtx() *Context {
	if c.Ctx == nil {
		c.Ctx = fiber.New().AcquireCtx(&fasthttp.RequestCtx{})
	}

	return c
}

// GetSource get source
func (c *Context) GetSource() string {
	return c.Get("Source")
}

// Platform get source
func (c *Context) Platform() string {
	return c.Get("Platform")
}

// Username username (basic auth)
func (c *Context) Username() string {
	v, ok := c.Locals(UsernameKey).(string)
	if ok {
		return v
	}

	return ""
}
