package config

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// RR -> for use to return result model
var (
	RR = &ReturnResult{}
)

// Result result
type Result struct {
	Code        int               `json:"code" mapstructure:"code"`
	Description LocaleDescription `json:"message" mapstructure:"localization"`
}

// SwaggerInfoResult swagger info result
type SwaggerInfoResult struct {
	Code        int    `json:"code"`
	Description string `json:"message"`
}

// WithLocale with locale
func (rs Result) WithLocale(c *fiber.Ctx) Result {
	lacale, ok := c.Locals("lang").(string)
	if !ok {
		rs.Description.Locale = "th"
	}
	rs.Description.Locale = lacale
	return rs
}

// Error error description
func (rs Result) Error() string {
	if rs.Description.Locale == "th" {
		return rs.Description.TH
	}
	return rs.Description.EN
}

// ErrorCode error code
func (rs Result) ErrorCode() int {
	return rs.Code
}

// HTTPStatusCode http status code
func (rs Result) HTTPStatusCode() int {
	switch rs.Code {
	case 0, 200: // success
		return http.StatusOK
	case 404: // not found
		return http.StatusNotFound
	case 401: // unauthorized
		return http.StatusUnauthorized
	}

	return http.StatusBadRequest
}

// ReturnResult return result model
type ReturnResult struct {
	JSONDuplicateOrInvalidFormat    Result `mapstructure:"json_duplicate_or_invalid_format"`
	InvalidUsernameOrPassword       Result `mapstructure:"invalid_username_or_password"`
	InvalidCitizenID                Result `mapstructure:"invalid_citizen_id"`
	InvalidTaxID                    Result `mapstructure:"invalid_tax_id"`
	InvalidPassword                 Result `mapstructure:"invalid_password"`
	InvalidPhoneNumber              Result `mapstructure:"invalid_phone_number"`
	InvalidEmail                    Result `mapstructure:"invalid_email"`
	InvalidReservedEmail            Result `mapstructure:"invalid_reserved_email"`
	InvalidToken                    Result `mapstructure:"invalid_token"`
	InvalidPermissionRole           Result `mapstructure:"invalid_permission_role"`
	CitizenIDAlreadyExists          Result `mapstructure:"citizen_id_already_exists"`
	TaxIDAlreadyExists              Result `mapstructure:"tax_id_already_exists"`
	PhoneNumberAlreadyExists        Result `mapstructure:"phone_number_already_exists"`
	EmailAlreadyExists              Result `mapstructure:"email_already_exists"`
	ReservedEmailAlreadyExists      Result `mapstructure:"reserved_email_already_exists"`
	MissingOccupation               Result `mapstructure:"missing_occupation"`
	MissingCompanyName              Result `mapstructure:"missing_company_name"`
	MissingEmail                    Result `mapstructure:"missing_email"`
	InvalidPrefixUpload             Result `mapstructure:"invalid_prefix_path_upload"`
	InvalidMaximumSize              Result `mapstructure:"invalid_maximum_size"`
	InvalidTypeImageFile            Result `mapstructure:"invalid_type_image_file"`
	UploadFileFail                  Result `mapstructure:"upload_file_fail"`
	QueueFull                       Result `mapstructure:"queue_full"`
	InvalidUserID                   Result `mapstructure:"invalid_user_id"`
	InvalidCodeOrExpired            Result `mapstructure:"invalid_code_or_expired"`
	AlreadyVerified                 Result `mapstructure:"already_verified"`
	InvalidAddressType              Result `mapstructure:"invalid_address_type"`
	UsernameNotFound                Result `mapstructure:"username_not_found"`
	InvalidExpressType              Result `mapstructure:"invalid_express_type"`
	InvalidSecretKey                Result `mapstructure:"invalid_secret_key"`
	InvalidIdentificationNumber     Result `mapstructure:"invalid_identification_number"`
	AlreadyApproved                 Result `mapstructure:"already_approved"`
	InvalidApprovalType             Result `mapstructure:"invalid_approval_type"`
	NoProductPurchase               Result `mapstructure:"no_product_purchase"`
	AddressTypeNotDeletable         Result `mapstructure:"address_type_not_deletable"`
	AlreadySubmittedBooking         Result `mapstructure:"already_submitted_booking"`
	InvalidUsername                 Result `mapstructure:"invalid_username"`
	AddressTypeNotCreatable         Result `mapstructure:"address_type_not_creatable"`
	InvalidTotalProduct             Result `mapstructure:"invalid_total_product"`
	InvalidPermissionAccess         Result `mapstructure:"invalid_permission_access"`
	InvalidBookingStatus            Result `mapstructure:"invalid_booking_status"`
	UnableToProcessMultipleCarQueue Result `mapstructure:"unable_to_process_multiple_car_queue"`
	InvalidCarQueueStatus           Result `mapstructure:"invalid_car_queue_status"`
	InvalidRating                   Result `mapstructure:"invalid_rating"`
	AlreadyRating                   Result `mapstructure:"already_rating"`
	Internal                        struct {
		Success          Result `mapstructure:"success" json:"success"`
		General          Result `mapstructure:"general" json:"general"`
		BadRequest       Result `mapstructure:"bad_request" json:"bad_request"`
		ConnectionError  Result `mapstructure:"connection_error" json:"connection_error"`
		DatabaseNotFound Result `mapstructure:"database_not_found" json:"database_not_found"`
		Unauthorized     Result `mapstructure:"unauthorized" json:"unauthorized"`
	} `mapstructure:"internal" json:"internal"`
}

// LocaleDescription locale description
type LocaleDescription struct {
	EN     string `mapstructure:"en"`
	TH     string `mapstructure:"th"`
	Locale string `mapstructure:"success"`
}

// MarshalJSON marshall json
func (ld LocaleDescription) MarshalJSON() ([]byte, error) {
	if strings.ToLower(ld.Locale) == "th" {
		return json.Marshal(ld.TH)
	}
	return json.Marshal(ld.EN)
}

// UnmarshalJSON unmarshal json
func (ld *LocaleDescription) UnmarshalJSON(data []byte) error {
	var res string
	err := json.Unmarshal(data, &res)
	if err != nil {
		return err
	}
	ld.EN = res
	ld.Locale = "en"
	return nil
}

// InitReturnResult init return result
func InitReturnResult(configPath string) error {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName("return_result")

	if err := v.ReadInConfig(); err != nil {
		logrus.Error("read config file error:", err)
		return err
	}

	if err := bindingReturnResult(v, RR); err != nil {
		logrus.Error("binding config error:", err)
		return err
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		logrus.Info("config file changed:", e.Name)
		if err := bindingReturnResult(v, RR); err != nil {
			logrus.Error("binding error:", err)
		}
		logrus.Infof("Initial 'Return Result'. %+v", RR)
	})
	return nil
}

// bindingReturnResult binding return result
func bindingReturnResult(vp *viper.Viper, rr *ReturnResult) error {
	if err := vp.Unmarshal(&rr); err != nil {
		logrus.Error("unmarshal config error:", err)
		return err
	}
	return nil
}

// CustomMessage custom message
func (rr *ReturnResult) CustomMessage(messageEN, messageTH string) Result {
	return Result{
		Code: 999,
		Description: LocaleDescription{
			EN: messageEN,
			TH: messageTH,
		},
	}
}
