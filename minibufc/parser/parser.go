package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

/**
#structure
# filename: mellon.mb

config float_precision = 3;

Vector {
  x: float;
  y: float;
  z: float;
}

Config {
  auto_restart: bool;
  id: number;
  user_name: string;
  score: float = 0.0;      // default value for bachward compatibility
}
*/

type Schema struct {
	Name     string
	Fields   []Field
	Defaults map[string]interface{}
}

type Field struct {
	Name     string
	DataType string
}

type Parser struct {
	Schemas []Schema
	Config  map[string]interface{}
}

func Initialize() {}

func ParseInputFiles(inputFiles []string) *Parser {
	p := &Parser{
		Config: make(map[string]interface{}),
	}

	for _, file := range inputFiles {
		if !verifyFilePath(file) {
			color.New(color.FgRed).Fprintf(os.Stderr, "Invalid file path: %s\n", file)
			return nil
			// continue
		}
		f, err := os.Open(file)
		if err != nil {
			color.New(color.FgRed).Fprintf(os.Stderr, "Error opening file %s: %v\n", file, err)
			return nil
			// continue // or handle error
		}
		scanner := bufio.NewScanner(f)
		inStruct := false
		var currentSchema *Schema
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}
			if !inStruct {
				if strings.Contains(line, "{") {
					name := strings.TrimSpace(strings.Split(line, "{")[0])
					currentSchema = &Schema{Name: name, Defaults: make(map[string]interface{})}
					inStruct = true
				} else if strings.Contains(line, "=") {
					parts := strings.Split(line, "=")
					if len(parts) == 2 {
						key := strings.TrimSpace(parts[0])
						valStr := strings.TrimSpace(parts[1])
						if val, err := strconv.Atoi(valStr); err == nil {
							p.Config[key] = val
						} else if val, err := strconv.ParseFloat(valStr, 64); err == nil {
							p.Config[key] = val
						} else {
							p.Config[key] = valStr
						}
					}
				}
			} else {
				if strings.Contains(line, "}") {
					p.Schemas = append(p.Schemas, *currentSchema)
					inStruct = false
					currentSchema = nil
				} else {
					// parse field
					parts := strings.Split(line, ":")
					if len(parts) != 2 {
						continue
					}
					name := strings.TrimSpace(parts[0])
					dataTypeAndDefault := strings.TrimSpace(parts[1])
					dataTypeAndDefault = strings.TrimSuffix(dataTypeAndDefault, ";")
					var dataType, defaultValue string
					if idx := strings.Index(dataTypeAndDefault, "="); idx != -1 {
						dataType = strings.TrimSpace(dataTypeAndDefault[:idx])
						defaultValue = strings.TrimSpace(dataTypeAndDefault[idx+1:])
					} else {
						dataType = dataTypeAndDefault
					}
					field := Field{Name: name, DataType: dataType}
					currentSchema.Fields = append(currentSchema.Fields, field)
					if defaultValue != "" {
						//remove everything after ; in defaultValue
						defaultValue = strings.SplitN(defaultValue, ";", 2)[0]
						val, err := parseValue(dataType, defaultValue)
						if err != nil {
							color.New(color.FgRed).Fprintf(os.Stderr, "Error parsing default value for %s: %v\n", name, err)
							return nil
						} else {
							currentSchema.Defaults[name] = val
						}
					}
				}
			}
		}
		if err := scanner.Err(); err != nil {
			color.New(color.FgRed).Fprintf(os.Stderr, "Error reading file %s: %v\n", file, err)
			return nil
		}
		f.Close()
	}

	if err := verifySchemas(p); err != nil {
		color.New(color.FgRed).Fprintf(os.Stderr, "Invalid schema: %v\n", err)
		return nil
	}

	return p
}

func verifySchemas(s *Parser) error {
	// make sure no two schemas have the same name
	// make sure all field data types are valid
	// make sure no two fields in the same schema have the same name
	schemaNames := make(map[string]bool)
	validDataTypes := map[string]bool{
		"bool":   true,
		"number": true,
		"float":  true,
		"string": true,
	}
	for _, schema := range s.Schemas {
		if schemaNames[schema.Name] {
			return fmt.Errorf("duplicate schema name: %s", schema.Name)
		}
		schemaNames[schema.Name] = true
		fieldNames := make(map[string]bool)
		for _, field := range schema.Fields {
			if fieldNames[field.Name] {
				return fmt.Errorf("duplicate field name '%s' in schema '%s'", field.Name, schema.Name)
			}
			fieldNames[field.Name] = true
			if !validDataTypes[field.DataType] {
				return fmt.Errorf("invalid data type '%s' for field '%s' in schema '%s'", field.DataType, field.Name, schema.Name)
			}
		}
	}
	return nil
}

func verifyFilePath(filePath string) bool {
	if !strings.HasSuffix(filePath, ".mb") {
		return false
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func parseValue(dataType, value string) (interface{}, error) {
	switch dataType {
	case "bool":
		if value != "true" && value != "false" {
			return nil, fmt.Errorf("invalid bool value: %s", value)
		}
		return strings.ToLower(value) == "true", nil
	case "number":
		i, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid number value: %s", value)
		}
		return i, nil
	case "float":
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float value: %s", value)
		}
		return f, nil
	case "string":
		return value, nil
	default:
		return nil, fmt.Errorf("unknown data type: %s", dataType)
	}
}
