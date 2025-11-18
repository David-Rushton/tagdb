package cli

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func unmarshalArgs(args []string, target any) error {
	// Validate input.
	ptrValue := reflect.ValueOf(target)
	if ptrValue.Kind() != reflect.Ptr || ptrValue.IsNil() {
		return fmt.Errorf("unmarshal target must be a non-nil pointer")
	}

	// Get the struct.
	structValue := ptrValue.Elem()
	if structValue.Kind() != reflect.Struct {
		return fmt.Errorf("unmarshal target must be struct")
	}
	structType := structValue.Type()

	argFields := map[int]argField{}
	optionFields := map[string]optionField{}

	// Iterate over struct fields.
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		// Get tags.
		argTag, argTagFound := fieldType.Tag.Lookup("arg")
		optionTag, optionTagFound := fieldType.Tag.Lookup("option")
		helpTag, helpTagFound := fieldType.Tag.Lookup("help")

		// Validate.
		if !(argTagFound || optionTagFound) {
			continue
		}

		if !helpTagFound || strings.TrimSpace(helpTag) == "" {
			return fmt.Errorf("missing help tag for field %s", fieldType.Name)
		}

		if argTagFound && optionTagFound {
			return fmt.Errorf("field %s cannot have both arg and option tags", fieldType.Name)
		}

		if !field.CanSet() {
			return fmt.Errorf("cannot unmarshal unexported field %s", fieldType.Name)
		}

		if !isSupportedKind(field) {
			return fmt.Errorf("unsupported field type %s for field %s", fieldType.Type.String(), fieldType.Name)
		}

		// Process.
		if argTagFound {
			arg, err := parseArgTag(argTag, helpTag)
			if err != nil {
				return fmt.Errorf("invalid arg tag for field %s: %w", fieldType.Name, err)
			}

			if _, exists := argFields[arg.index]; exists {
				return fmt.Errorf("duplicate arg index %d for field %s", arg.index, fieldType.Name)
			}

			argFields[arg.index] = argField{arg, field}
			continue
		}

		if optionTagFound {
			option, err := parseOptionTag(optionTag, helpTag)
			if err != nil {
				return fmt.Errorf("invalid option tag for field %s: %w", fieldType.Name, err)
			}

			if option.short != "" {
				if _, exists := optionFields[option.short]; exists {
					return fmt.Errorf("duplicate option %s for field %s", option.short, fieldType.Name)
				}
				optionFields[option.short] = optionField{option, field}
			}

			if option.long != "" {
				if _, exists := optionFields[option.long]; exists {
					return fmt.Errorf("duplicate option %s for field %s", option.long, fieldType.Name)
				}
				optionFields[option.long] = optionField{option, field}
			}

			continue
		}
	}

	// Argument fields must:
	//  - be continuous from index 0 to N
	//  - all required args must come before optional args
	//  - only final can be variadic
	var optionalArgFound bool
	for i := range len(argFields) {
		if _, exists := argFields[i]; !exists {
			return fmt.Errorf("missing arg field for index %d", i)
		}

		if !argFields[i].arg.required {
			optionalArgFound = true
			continue
		}

		if optionalArgFound && argFields[i].arg.required {
			return fmt.Errorf("required arg found after optional arg at index %d", i)
		}

		if isArrayOrSlice(argFields[i].field) && i != len(argFields)-1 {
			return fmt.Errorf("variadic args must be the final arg")
		}
	}

	// Map options to fields.
	consumedArgs := map[int]bool{}

	for i := range args {
		if consumedArgs[i] {
			continue
		}

		if optionField, found := optionFields[args[i]]; found {
			consumedArgs[i] = true

			// Set boolean flags to true.
			if optionField.field.Kind() == reflect.Bool {
				err := setFieldValue(optionField.field, "true")
				if err != nil {
					return fmt.Errorf("invalid value for option %s: %w", args[i], err)
				}
				continue
			}

			// Get the value.
			if i+1 >= len(args) {
				return fmt.Errorf("missing value for option %s", args[i])
			}

			err := setFieldValue(optionField.field, args[i+1])
			if err != nil {
				return fmt.Errorf("invalid value for option %s: %w", args[i], err)
			}

			consumedArgs[i+1] = true
		}
	}

	// Map args to fields.
	var argIndex int
	for i := range len(args) {
		if consumedArgs[i] {
			continue
		}

		if argField, found := argFields[argIndex]; found {
			var exitEarly bool
			value := args[i]

			if isArrayOrSlice(argField.field) {
				// Variadic argument consumes all remaining, unless we reach an arg that
				// has already been consumed (e.g. a recognised option handled above).
				for j := i + 1; j < len(args); j++ {
					if consumedArgs[j] {
						break
					}
					value += "," + args[j]
					consumedArgs[j] = true
				}
				exitEarly = true
			}

			err := setFieldValue(argField.field, value)
			if err != nil {
				return fmt.Errorf("invalid value for arg %s: %w", argField.arg.name, err)
			}

			consumedArgs[i] = true

			if exitEarly {
				break
			}
		}

		argIndex++
	}

	// Validate required args.
	for i, argField := range argFields {
		if argField.arg.required {
			if !consumedArgs[i] {
				return fmt.Errorf("missing required arg %s", argField.arg.name)
			}
		}
	}

	if len(consumedArgs) != len(args) {
		sb := strings.Builder{}
		for i, arg := range args {
			if consumedArgs[i] {
				continue
			}

			sb.WriteString(" ")
			sb.WriteString(arg)
		}
		return fmt.Errorf("unexpected args:%s", sb.String())
	}

	return nil
}

type arg struct {
	name     string
	required bool
	index    int
	help     string
}

type argField struct {
	arg   *arg
	field reflect.Value
}

func parseArgTag(tag, help string) (*arg, error) {
	elements := strings.SplitN(tag, ":", 2)

	if len(elements) != 2 {
		err := fmt.Errorf("invalid arg tag format: index:<required-name> or index:[optional-name]")
		return nil, err
	}

	indexPart := strings.Trim(elements[0], " ")
	index, err := strconv.Atoi(indexPart)
	if err != nil {
		err = fmt.Errorf("invalid arg tag format: index is not a number")
		return nil, err
	}

	var found bool
	var required bool
	name := strings.Trim(elements[1], " ")

	if strings.HasPrefix(name, "<") && strings.HasSuffix(name, ">") {
		required = true
		found = true
	}

	if strings.HasPrefix(name, "[") && strings.HasSuffix(name, "]") {
		required = true
		found = true
	}

	if !found {
		err = fmt.Errorf("invalid arg tag format: missing required (<>) or optional ([]) notation")
		return nil, err
	}

	return &arg{name, required, index, help}, nil
}

const (
	shortOptionRegex = `^-[A-z0-9]{1}$`
	longOptionRegex  = `^--[A-z0-9-]{2,}$`
)

var (
	shortRegex = regexp.MustCompile(shortOptionRegex)
	longRegex  = regexp.MustCompile(longOptionRegex)
)

type option struct {
	short string
	long  string
	help  string
}

type optionField struct {
	option *option
	field  reflect.Value
}

func parseOptionTag(tag, help string) (*option, error) {
	elements := strings.SplitN(tag, "|", 2)

	if len(elements) == 0 {
		err := fmt.Errorf("invalid option tag format: -s or --long or -s|--long")
		return nil, err
	}

	var short string
	var shortFound int
	var long string
	var longFound int

	for _, rawElement := range elements {
		element := strings.Trim(rawElement, " ")

		if shortRegex.MatchString(element) {
			short = element
			shortFound++
		}

		if longRegex.MatchString(element) {
			long = element
			longFound++
		}
	}

	if shortFound > 1 {
		err := fmt.Errorf("invalid option tag format: only 1 short option allowed")
		return nil, err
	}

	if longFound > 1 {
		err := fmt.Errorf("invalid option tag format: only 1 long option allowed")
		return nil, err
	}

	if shortFound+longFound == 0 {
		err := fmt.Errorf("invalid option tag format: at least 1 short or long option required")
		return nil, err
	}

	return &option{short, long, help}, nil
}

func isSupportedKind(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.Bool:
		return true

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true

	case reflect.Float32, reflect.Float64:
		return true

	case reflect.Complex64, reflect.Complex128:
		return true

	case reflect.Array, reflect.Slice:
		elementType := reflect.New(field.Type().Elem()).Elem()
		return isSupportedKind(elementType)

	case reflect.String:
		return true

	default:
		return false
	}
}

func isArrayOrSlice(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.Array, reflect.Slice:
		return true
	default:
		return false
	}
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.Bool:
		field.SetBool(true)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer value: %s", value)
		}
		field.SetInt(num)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid unsigned integer value: %s", value)
		}
		field.SetUint(num)

	case reflect.Float32, reflect.Float64:
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float value: %s", value)
		}
		field.SetFloat(num)

	case reflect.Complex64, reflect.Complex128:
		num, err := strconv.ParseComplex(value, 128)
		if err != nil {
			return fmt.Errorf("invalid complex value: %s", value)
		}
		field.SetComplex(num)

	case reflect.Array, reflect.Slice:
		// For simplicity, assume comma-separated values.
		values := strings.Split(value, ",")
		slice := reflect.MakeSlice(field.Type(), len(values), len(values))
		for i, v := range values {
			err := setFieldValue(slice.Index(i), v)
			if err != nil {
				return err
			}
		}
		field.Set(slice)

	case reflect.String:
		field.SetString(value)

	default:
		return fmt.Errorf("unsupported field type")
	}

	return nil
}
