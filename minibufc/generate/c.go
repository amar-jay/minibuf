package generate

import (
	"fmt"
	"strings"

	"github.com/amar-jay/minibuf/parser"
)

/*
**
```c
#include "minibuf.h"

MiniBuf_t mb = mb_init();
char buf[256];

// parse Vector
Vector_t v;
int err = mb_vector_parse(buf, &v);

	if (err != MB_OK) {
	    // handle errors
	}

// parse Config
Config_t conf;
err = mb_config_parse(buf, &conf);

	if (err != MB_OK) {
	    // handle errors
	}

```
*/
func C(p *parser.Parser) (string, string) {
	header := generateHeader(p)
	c := generateCCode(p)
	return header, c
}

func generateHeader(p *parser.Parser) string {
	var sb strings.Builder
	sb.WriteString("#ifndef MINIBUF_H\n#define MINIBUF_H\n\n")
	sb.WriteString("#include <stdbool.h>\n#include <stdint.h>\n#include <stddef.h>\n\n")
	sb.WriteString("#define MB_OK 0\n#define MB_ERR_INVALID_FORMAT 1\n#define MB_ERR_BUFFER_TOO_SMALL 2\n\n")
	sb.WriteString("extern int mb_float_precision;\n\n")
	for _, schema := range p.Schemas {
		sb.WriteString("typedef struct {\n")
		for _, field := range schema.Fields {
			ctype := cType(field.DataType)
			if field.DataType == "string" {
				sb.WriteString(fmt.Sprintf("    %s %s[256];\n", ctype, field.Name))
			} else {
				sb.WriteString(fmt.Sprintf("    %s %s;\n", ctype, field.Name))
			}
		}
		sb.WriteString(fmt.Sprintf("} %s_t;\n\n", strings.ToLower(schema.Name)))
	}
	for _, schema := range p.Schemas {
		name := strings.ToLower(schema.Name)
		sb.WriteString(fmt.Sprintf("int mb_%s_parse(const char* buf, %s_t* out);\n", name, name))
		sb.WriteString(fmt.Sprintf("int mb_%s_serialize(const %s_t* in, char* buf, size_t buf_size);\n\n", name, name))
	}
	sb.WriteString("#endif\n")
	return sb.String()
}

func cType(dt string) string {
	switch dt {
	case "bool":
		return "bool"
	case "number":
		return "int32_t"
	case "float":
		return "float"
	case "string":
		return "char"
	}
	return "void*"
}

func generateCCode(p *parser.Parser) string {
	var sb strings.Builder
	sb.WriteString("#include \"minibuf.h\"\n#include <string.h>\n#include <stdio.h>\n#include <stdlib.h>\n#include <math.h>\n\n")
	if val, ok := p.Config["float_precision"]; ok {
		if f, ok := val.(float64); ok {
			sb.WriteString(fmt.Sprintf("int mb_float_precision = %d;\n\n", int(f)))
		} else if i, ok := val.(int); ok {
			sb.WriteString(fmt.Sprintf("int mb_float_precision = %d;\n\n", i))
		} else {
			sb.WriteString("int mb_float_precision = 3;\n\n")
		}
	} else {
		sb.WriteString("int mb_float_precision = 3;\n\n")
	}
	for _, schema := range p.Schemas {
		name := strings.ToLower(schema.Name)
		sb.WriteString(fmt.Sprintf("int mb_%s_parse(const char* buf, %s_t* out) {\n", name, name))
		sb.WriteString("    char* start = strchr(buf, '[');\n")
		sb.WriteString("    if (!start) return MB_ERR_INVALID_FORMAT;\n")
		sb.WriteString("    start++;\n")
		sb.WriteString("    char* end = strchr(start, ']');\n")
		sb.WriteString("    if (!end) return MB_ERR_INVALID_FORMAT;\n")
		sb.WriteString("    *end = '\\0';\n")
		sb.WriteString("    int count = atoi(start);\n")
		sb.WriteString("    *end = ']';\n")
		sb.WriteString("    char* values = end + 1;\n")
		sb.WriteString("    char* vals = strdup(values);\n")
		sb.WriteString("    char* token = strtok(vals, \";\");\n")
		sb.WriteString("    int i = 0;\n")
		// set defaults
		for _, field := range schema.Fields {
			if def, ok := schema.Defaults[field.Name]; ok {
				switch field.DataType {
				case "bool":
					if def.(bool) {
						sb.WriteString(fmt.Sprintf("    out->%s = true;\n", field.Name))
					} else {
						sb.WriteString(fmt.Sprintf("    out->%s = false;\n", field.Name))
					}
				case "number":
					sb.WriteString(fmt.Sprintf("    out->%s = %d;\n", field.Name, def.(int)))
				case "float":
					sb.WriteString(fmt.Sprintf("    out->%s = %.3f;\n", field.Name, def.(float64)))
				case "string":
					sb.WriteString(fmt.Sprintf("    strcpy(out->%s, \"%s\");\n", field.Name, def.(string)))
				}
			}
		}
		sb.WriteString(fmt.Sprintf("    while (token && i < %d) {\n", len(schema.Fields)))
		for idx, field := range schema.Fields {
			sb.WriteString(fmt.Sprintf("        if (i == %d) {\n", idx))
			switch field.DataType {
			case "bool":
				sb.WriteString(fmt.Sprintf("            out->%s = strcmp(token, \"T\") == 0;\n", field.Name))
			case "number":
				sb.WriteString(fmt.Sprintf("            out->%s = atoi(token);\n", field.Name))
			case "float":
				sb.WriteString(fmt.Sprintf("            out->%s = atof(token);\n", field.Name))
			case "string":
				sb.WriteString(fmt.Sprintf("            strcpy(out->%s, token);\n", field.Name))
			}
			sb.WriteString("        }\n")
		}
		sb.WriteString("        i++;\n")
		sb.WriteString("        token = strtok(NULL, \";\");\n")
		sb.WriteString("    }\n")
		sb.WriteString("    free(vals);\n")
		sb.WriteString("    return MB_OK;\n")
		sb.WriteString("}\n\n")
		// serialize
		sb.WriteString(fmt.Sprintf("int mb_%s_serialize(const %s_t* in, char* buf, size_t buf_size) {\n", name, name))
		num := len(schema.Fields)
		sb.WriteString(fmt.Sprintf("    int len = snprintf(buf, buf_size, \"[%d]\", %d);\n", num, num))
		sb.WriteString("    if (len >= buf_size) return MB_ERR_BUFFER_TOO_SMALL;\n")
		sb.WriteString("    char* pos = buf + len;\n")
		for idx, field := range schema.Fields {
			sep := ";"
			if idx == 0 {
				sep = ""
			}
			switch field.DataType {
			case "bool":
				sb.WriteString(fmt.Sprintf("    len += snprintf(pos, buf_size - len, \"%s%%s\", in->%s ? \"T\" : \"F\");\n", sep, field.Name))
			case "number":
				sb.WriteString(fmt.Sprintf("    len += snprintf(pos, buf_size - len, \"%s%%d\", in->%s);\n", sep, field.Name))
			case "float":
				sb.WriteString(fmt.Sprintf("    len += snprintf(pos, buf_size - len, \"%s%%s%%d.%%0*d\", in->%s < 0 ? \"-\" : \"\", abs((int)in->%s), mb_float_precision, (int)((fabsf(in->%s) - abs((int)in->%s)) * powf(10, mb_float_precision) + 0.5f));\n", sep, field.Name, field.Name, field.Name, field.Name))
			case "string":
				sb.WriteString(fmt.Sprintf("    len += snprintf(pos, buf_size - len, \"%s%%s\", in->%s);\n", sep, field.Name))
			}
			sb.WriteString("    if (len >= buf_size) return MB_ERR_BUFFER_TOO_SMALL;\n")
			sb.WriteString("    pos = buf + len;\n")
		}
		sb.WriteString("    return MB_OK;\n")
		sb.WriteString("}\n\n")
	}
	return sb.String()
}
