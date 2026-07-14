package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func HandleValidationErrors(err error) gin.H {
	// Type assertion: the validator library returns errors as the interface `error`,
	// but the concrete type underneath is `validator.ValidationErrors` (a slice of
	// per-field errors). `ok` is true only when the assertion matches, so we only
	// enter this block for actual validation failures.
	if validationError, ok := err.(validator.ValidationErrors); ok {
		// Result map: field-path (snake_case) → human-readable Vietnamese message.
		errors := make(map[string]string)

		// Loop over each field error; `_` discards the slice index.
		for _, e := range validationError {
			// e.Namespace() is the full path incl. the struct name,
			// e.g. "CreateNewsRequest.Author.FirstName". [0] → the root struct name.
			root := strings.Split(e.Namespace(), ".")[0]

			// Strip the root prefix (+ its dot) → "Author.FirstName".
			rawPath := strings.TrimPrefix(e.Namespace(), root+".")

			// Split the remaining path into segments → ["Author", "FirstName"].
			parts := strings.Split(rawPath, ".")

			// Convert each segment to snake_case; index `i` is needed to write back.
			for i, part := range parts {
				// A slice element looks like "Images[0]": snake_case only the name
				// and keep the "[0]" suffix intact.
				// NOTE: bug — "part" is quoted, so this checks a literal string and is
				// always false. It should be strings.Contains(part, "[").
				if strings.Contains("part", "[") {
					idx := strings.Index(part, "[")   // position of '['
					base := camelToSnake(part[:idx])  // name before '[' → snake_case
					index := part[idx:]               // "[0]" kept as-is
					parts[i] = base + index           // e.g. "images[0]"
				} else {
					// No brackets → snake_case the whole segment ("FirstName" → "first_name").
					parts[i] = camelToSnake(part)
				}
			}

			// Rejoin segments into a dotted path → "author.first_name".
			fieldPath := strings.Join(parts, ".")

			switch e.Tag() {
			case "gt":
				errors[fieldPath] = fmt.Sprintf("%s phải lớn hơn %s", fieldPath, e.Param())
			case "lt":
				errors[fieldPath] = fmt.Sprintf("%s phải nhỏ hơn %s", fieldPath, e.Param())
			case "gte":
				errors[fieldPath] = fmt.Sprintf("%s phải lớn hơn hoặc bằng %s", fieldPath, e.Param())
			case "lte":
				errors[fieldPath] = fmt.Sprintf("%s phải nhỏ hơn hoặc bằng %s", fieldPath, e.Param())
			case "uuid":
				errors[fieldPath] = fmt.Sprintf("%s phải là UUID hợp lệ", fieldPath)
			case "slug":
				errors[fieldPath] = fmt.Sprintf("%s chỉ được chứa chữ thường, số, dấu gạch ngang hoặc dấu chấm", fieldPath)
			case "min":
				errors[fieldPath] = fmt.Sprintf("%s phải nhiều hơn %s ký tự", fieldPath, e.Param())
			case "max":
				errors[fieldPath] = fmt.Sprintf("%s phải ít hơn %s ký tự", fieldPath, e.Param())
			case "min_int":
				errors[fieldPath] = fmt.Sprintf("%s phải có giá trị lớn hơn %s", fieldPath, e.Param())
			case "max_int":
				errors[fieldPath] = fmt.Sprintf("%s phải có giá trị bé hơn %s", fieldPath, e.Param())
			case "oneof":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[fieldPath] = fmt.Sprintf("%s phải là một trong các giá trị: %s", fieldPath, allowedValues)
			case "required":
				errors[fieldPath] = fmt.Sprintf("%s là bắt buộc", fieldPath)
			case "search":
				errors[fieldPath] = fmt.Sprintf("%s chỉ được chứa chữ thường, in hoa, số và khoảng trắng", fieldPath)
			case "email":
				errors[fieldPath] = fmt.Sprintf("%s phải đúng định dạng là email", fieldPath)
			case "datetime":
				errors[fieldPath] = fmt.Sprintf("%s phải theo đúng định dạng YYYY-MM-DD", fieldPath)
			case "file_ext":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[fieldPath] = fmt.Sprintf("%s chỉ cho phép những file có extension: %s", fieldPath, allowedValues)
			}
		}

		return gin.H{"error": errors}

	}

	return gin.H{"error": "Yêu cầu không hợp lệ" + err.Error()}
}

func RegisterValidators() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("failed to get validator engine")
	}

	var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:[-.][a-z0-9]+)*$`)
	v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		return slugRegex.MatchString(fl.Field().String())
	})

	var searchRegex = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
	v.RegisterValidation("search", func(fl validator.FieldLevel) bool {
		return searchRegex.MatchString(fl.Field().String())
	})

	v.RegisterValidation("min_int", func(fl validator.FieldLevel) bool {
		minStr := fl.Param()
		minVal, err := strconv.ParseInt(minStr, 10, 64)
		if err != nil {
			return false
		}

		return fl.Field().Int() >= minVal
	})

	v.RegisterValidation("max_int", func(fl validator.FieldLevel) bool {
		maxStr := fl.Param()
		// 10 là cơ số, 64 là kích thước bit
		maxVal, err := strconv.ParseInt(maxStr, 10, 64)
		if err != nil {
			return false
		}

		return fl.Field().Int() <= maxVal
	})

	v.RegisterValidation("file_ext", func(fl validator.FieldLevel) bool {
		filename := fl.Field().String()

		allowedStr := fl.Param()
		if allowedStr == "" {
			return false
		}

		allowedExt := strings.Fields(allowedStr) // "jpg png gif"  →  []string{"jpg", "png", "gif"}
		// 1. filepath.Ext(filename) — returns the extension including the dot.
		// "photo.JPG" → ".JPG"
		// 2. strings.ToLower(...) — lowercases it so comparison is case-insensitive.
		// ".JPG" → ".jpg"
		// 3. strings.TrimPrefix(..., ".") — removes the leading "." (only if present).
		// ".jpg" → "jpg"
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")

		//The two forms of range

		// for i, v := range slice { }  // i = index, v = value
		// for _, v := range slice { }  // ignore index, only value  ← this code
		// for i := range slice { }     // only index, no value

		// Equivalent without range

		// for i := 0; i < len(allowedExt); i++ {
		//     if ext == strings.ToLower(allowedExt[i]) {
		//         return true
		//     }
		// }

		for _, allowed := range allowedExt {
			if ext == strings.ToLower(allowed) {
				return true
			}
		}

		return false
	})

	return nil
}

func ValidationRequired(field, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", field)
	}
	return nil
}

func ValidationStringLength(field, value string, min, max int) error {
	length := len(value)
	if length < min || length > max {
		return fmt.Errorf("%s must be between %d and %d characters long", field, min, max)
	}
	return nil
}

func ValidationRegex(field, value string, re *regexp.Regexp, errorMessage string) error {
	if !re.MatchString(value) {
		return fmt.Errorf("%s: %s", field, errorMessage)
	}
	return nil
}

func ValidationPositiveInt(field, value string) (int, error) {
	val, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid integer", field)
	}
	if val <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", field)
	}
	return val, nil
}

func ValidationUUID(field, value string) (uuid.UUID, error) {
	uid, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s must be a valid UUID", field)
	}
	return uid, nil
}

func ValidationInList(field, value string, allowed map[string]bool) error {
	if !allowed[value] {
		return fmt.Errorf("%s must be one of the following values: %s", field, keys(allowed))
	}
	return nil
}

func keys(m map[string]bool) string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}
