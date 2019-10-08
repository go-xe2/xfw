package xsql

import (
	"fmt"
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/util/gconv"
	"regexp"
	"strings"
)

// 根据规规则，生成sql查询字段列表
type FieldRule string
const (
	FIELD_RULE_ADMIN_DETAIL = "MD"
	FIELD_RULE_ADMIN_LIST	= "ML"
	FIELD_RULE_USER_DETAIL	= "UD"
	FIELD_RULE_USER_LIST	= "UL"
	FIELD_RULE_GUEST_DETAIL = "GD"
	FIELD_RULE_GUEST_LIST	= "GL"

)

type QueryField struct {
	// 表别名
	tableAlias string
	// 字段名称
	fieldName  string
}

func (f *QueryField) String() string {
	if f.tableAlias != "" {
		return f.tableAlias + "." + f.fieldName
	}
	return f.fieldName
}

func (f *QueryField) Parse(field string) {
	reg := regexp.MustCompile(`[\w]+`)
	if reg.MatchString(field) {
		matchItems := reg.FindAllStringSubmatch(field, -1)
		if len(matchItems) == 1 {
			f.fieldName = field
			f.tableAlias = ""
		} else if len(matchItems) == 2 {
			f.fieldName = matchItems[1][0]
			f.tableAlias = matchItems[0][0]
		} else {
			// 表达式
			f.fieldName = field
			f.tableAlias = ""
		}
	} else {
		f.tableAlias = ""
		f.fieldName = field
	}
}

var prefixFieldName = func(fieldName string, tableAlias ...string) string {
	if len(tableAlias) > 0 {
		arr := strings.Split(fieldName, " ")
		if len(arr) > 1 {
			// 表达式的字段不处理
			return fieldName
		}
		if strings.Index(fieldName, "(") >=0 && strings.Index(fieldName, ")") >= 0 {
			// 有函数的表达式不处理
			return fieldName
		}
		return tableAlias[0] + "." + fieldName
	}
	return fieldName
}

// 解析查询字名称
func parseFieldName(field string, fieldRule interface{}, tableAlias ...string) string {
	switch value := fieldRule.(type) {
	case string:
		if value == "" || value == "-" || value == field {
			return prefixFieldName(field, tableAlias...)
		}
		// 取别名返回
		return fmt.Sprintf("%s as %s", prefixFieldName(field, tableAlias...), value)
		break
	case []interface{}:
		valueLen := len(value)
		if valueLen == 0 {
			return prefixFieldName(field, tableAlias...)
		}
		// ["ML,UL,GL,MD,UD,GD,S"]
		if valueLen == 1 {
			// ["ML,UL,GL,MD,UD,GD,S"]
			return prefixFieldName(field, tableAlias...)
		}
		expr := gconv.String(value[1])
		fullFieldName := prefixFieldName(field, tableAlias...)
		if field != fullFieldName {
			expr = strings.Replace(expr, field, fullFieldName, -1)
		}
		if valueLen == 2 {
			// ["ML,UL,GL,MD,UD,GD,S", "from_unixtime(cr_date)"]
			return fmt.Sprintf("%s as %s", expr, field)
		}
		// ["ML,UL,GL,MD,UD,GD,S", "from_unixtime(cr_date)", "cr_date"]
		alias := gconv.String(value[2])
		return fmt.Sprintf("%s as %s", expr, alias)
		break
	}
	return ""
}

// 生成查询字段列表
func QueryFields(tableAlias string, fields ...string) string {
	if len(fields) == 0 {
		if tableAlias == "" {
			return "*"
		} else {
			return tableAlias + ".*"
		}
	}
	var results []string
	for _, s := range fields {
		if strings.Index(s, ".") >= 0 {
			results = append(results,s)
			continue
		}
		itemArr := strings.Split(s, ",")
		for _,fd := range itemArr {
			if strings.Index(fd, ".") >= 0 {
				results = append(results, fd)
			} else {
				// 该处暂未处理表达式中的字段添加前掇问题
				if tableAlias == "" {
					results = append(results,fd)
				} else {
					results = append(results, prefixFieldName(fd, tableAlias))
				}
			}
		}
	}
	return strings.Join(results, ",")
}


func BuildSqlField(fields map[string]interface{}, rule string, tableAlias ...string) string {
	if fields == nil {
		return "*"
	}
	var results []string
	for field, v := range fields {
		switch value := v.(type) {
		case string:
			if value == "-" {
				continue
			}
			fieldName := parseFieldName(field, value, tableAlias...)
			if fieldName == "" {
				continue
			}
			results = append(results, fieldName)
			break
		case []interface{}:
			valueLen := len(value)
			fieldName := parseFieldName(field, value, tableAlias...)
			if fieldName == "" {
				continue
			}
			if valueLen == 0 {
				results = append(results, fieldName)
				continue
			}
			srcRule := gconv.String(value[0])
			isInclude := srcRule == "" || srcRule == "S" || strings.Contains(srcRule, rule)
			isIncludeSrc := srcRule == "S" || strings.Contains(srcRule, "S")
			if !isInclude {
				continue
			}
			results = append(results, fieldName)
			if valueLen >= 3 {
				alias := gconv.String(value[2])
				if isIncludeSrc && alias != field {
					// 增加源字段输出
					results = append(results, prefixFieldName(field, tableAlias...))
				}
			}
			break
		}
	}
	if len(results) == 0 {
		return "*"
	}
	return strings.Join(results, ",")
}

// 使用selects选择器生成sql查询字段列表
// selects 可用值：nil, 字段列表字符串，[]string
func BuildSqlFieldBySelect(fields map[string]interface{}, selects interface{}, tableAlias ...string) string {
	if fields == nil {
		return  "*"
	}
	if selects == nil {
		return BuildSqlField(fields, "GL", tableAlias...)
	}
	var selectFields []string
	var results []string
	switch expr := selects.(type) {
	case string:
		selectFields = strings.Split(expr, ",")
		break
	case []string:
		for _, s := range expr {
			arr := strings.Split(s, ",")
			for _, fd := range arr {
				selectFields = append(selectFields, fd)
			}
		}
		break
	}

	for _, field := range selectFields {
		v, ok := fields[field]
		if !ok {
			continue
		}
		fieldName := parseFieldName(field, v)
		if fieldName == "" {
			continue
		}
		results = append(results, fieldName)
	}
	if len(results) == 0 {
		return "*"
	}
	return strings.Join(results, ",")
}

// 检查字段名是否有效
func CheckFieldNameValid(fields map[string]interface{}, fieldName string) bool {
	if fields == nil {
		return false
	}
	_, ok := fields[fieldName]
	return ok
}

// 生成所有字段列表数组
func BuildSqlAllFieldArray(fields map[string]interface{}) []string {
	var results []string
	if fields == nil {
		return []string{}
	}
	for field, v := range fields {
		fieldName := parseFieldName(field, v)
		if fieldName == "" {
			continue
		}
		results = append(results, fieldName)
	}
	return results
}

// 生成所有字段列表
func BuildSqlAllField(fields map[string]interface{}) string {
	arr := BuildSqlAllFieldArray(fields)
	if len(arr) == 0 {
		return "*"
	}
	return strings.Join(arr, ",")
}


// 解析过滤排序字段
// options.alias :表别名
func ParseSortFields(str string, options map[string]interface{}, selects ...map[string]bool) []string {
	var results []string
	items := strings.Split(str, ",")
	selectFields := map[string]bool{}
	hasFilter := len(selects) > 0
	if hasFilter {
		selectFields = selects[0]
	}
	alias := ""
	if options != nil {
		s, ok := options["alias"].(string)
		if ok {
			alias = s
		}
	}
	for _, item := range items {
		fieldVars 	:= strings.Split(item, " ")
		sortField 	:= fieldVars[0]
		sortDir 	:= "asc"
		if sortField == "" {
			continue
		}
		if len(fieldVars) > 1 {
			if fieldVars[1] == "desc" {
				sortDir = "desc"
			}
		}
		if hasFilter {
			b, ok := selectFields[sortField]
			if !ok || !b {
				continue
			}
			results = append(results, prefixFieldName(sortField, alias) + " " + sortDir)
		} else {
			results = append(results, prefixFieldName(sortField, alias) + " " + sortDir)
		}
	}
	return results
}

// 过滤排序字段字符串
func FilterSortFields(sort string, options map[string]interface{}, selects ...map[string]bool) string {
	fields := ParseSortFields(sort, options)
	return strings.Join(fields, ",")
}

// 生成sql时间段条件表达式
func CreateDateBetweenFilterSql(fieldName string, startDate string, endDate string) string {
	var startTime *gtime.Time = nil
	var endTime *gtime.Time = nil
	if startDate != "" {
		startTime = gtime.ParseTimeFromContent(startDate)
	}
	if endDate != "" {
		endTime = gtime.ParseTimeFromContent(endDate)
	}
	var where []string
	if startTime != nil {
		where = append(where, fmt.Sprintf("%s='%v'", fieldName, startTime.String()))
	}
	if endTime != nil {
		where = append(where, fmt.Sprintf("%s='%v'", fieldName, endTime.String()))
	}
	if where != nil {
		return strings.Join(where, " and ")
	}
	return ""
}

