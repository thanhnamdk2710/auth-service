package validation

var tagMessages = map[string]string{
	"required": "is required",
	"gte":      "must be at least %s characters",
	"lte":      "must be at most %s characters",
	"email":    "must be a valid email address",
	"eqfield":  "must match the %s field",
}
