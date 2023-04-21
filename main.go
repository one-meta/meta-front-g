package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	mapset "github.com/deckarep/golang-set"
	"github.com/spf13/viper"
)

var (
	conf          config
	path          string
	pathSeparator = string(os.PathSeparator)
)

type config struct {
	Field  Field  `json:"field"`
	Column Column `json:"column"`
	Router Router `json:"router"`
}
type Field struct {
	IgnoreField    []string `json:"ignoreField,omitempty"`
	IgnoreEntity   []string `json:"ignoreEntity,omitempty"`
	ParseWithType  []string `json:"parseWithType,omitempty"`
	ParseWithField []string `json:"parseWithField,omitempty"`
	ExtendField    []string `json:"extendField,omitempty"`
}

type Column struct {
	Templates []Template `json:"templates,omitempty"`
}
type Router struct {
	Templates []Template `json:"templates,omitempty"`
}

type Template struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type FieldData struct {
	Name     string `json:"name,omitempty"`
	DataType string `json:"dataType,omitempty"`
}

// 解析整个typings.d.ts
func main() {
	flag.StringVar(&path, "path", "", "typings.d.ts文件绝对路径")
	flag.Parse()
	if path == "" {
		path = "typings.d.ts"
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("当前路径未发现 typings.d.ts 或未指定path")
		flag.Usage()
		os.Exit(0)
	}
	fmt.Println("starting")
	now := time.Now()

	os.RemoveAll("Column")
	CheckDirAndMk(".", "Column")
	os.RemoveAll("Pages")
	os.RemoveAll("routers.ts")
	CheckDirAndMk(".", "Pages")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	line := bufio.NewReader(file)

	dataMap := make(map[string][]FieldData)
	var entityData []FieldData
	var entityName string
	for {
		content, _, err := line.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		line := strings.TrimSpace(string(content))
		// 去掉注释
		if strings.HasPrefix(line, "/**") {
			continue
		}
		// 获取到类型行开头
		if strings.HasPrefix(line, "type") && strings.Contains(line, "= {") {
			// 类型名称
			split := strings.Split(line, " ")
			entityName = split[1]

			// 存在，从map中取出来
			if _, ok := dataMap[entityName]; ok {
				entityData = dataMap[entityName]
			} else {
				// 不存在
				entityData = []FieldData{}
				dataMap[entityName] = entityData
			}
			continue
		}

		// 字段
		if strings.Contains(line, ":") && strings.Contains(line, ";") {
			fieldData := parseFieldData(line)
			entityData = append(entityData, *fieldData)
		}
		dataMap[entityName] = entityData
	}

	colTemplates := conf.Column.Templates
	routerTemplates := conf.Router.Templates
	templateMap := array2Map(colTemplates)
	routerMap := array2Map(routerTemplates)

	// 忽略的字段
	ignoreFields := conf.Field.IgnoreField
	ignoreFieldSet := array2Set(ignoreFields)
	IgnoreEntities := conf.Field.IgnoreEntity
	IgnoreEntitySet := array2Set(IgnoreEntities)

	// 根据类型解析
	parseWithTypes := conf.Field.ParseWithType
	// 根据字段解析
	parseWithFields := conf.Field.ParseWithField
	parseWithFieldSet := array2Set(parseWithFields)
	// 增加的字段
	extendFields := conf.Field.ExtendField

	// 生成路由文件
	routerName := "routes.ts"
	resultFile, err := os.Create(routerName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("生成 %s\n", routerName)

	routerWriter := bufio.NewWriter(resultFile)
	fmt.Fprintln(routerWriter, "export default [")

	for dataEntityName, v := range dataMap {
		if len(v) != 0 {
			ignoreEntityFlag := false
			lowerEntityName := strings.ToLower(dataEntityName)
			for val := range IgnoreEntitySet.Iterator().C {
				if strings.Contains(lowerEntityName, val.(string)) {
					ignoreEntityFlag = true
					break
				}
			}
			// 忽略小写开头的实体
			ignoreLowerFist := WithLowerCaseFirst(dataEntityName)
			if ignoreEntityFlag || ignoreLowerFist {
				continue
			}

			// 生成列文件
			fileName := dataEntityName + ".tsx"
			resultFile, err := os.Create(fmt.Sprintf("Column%s%s", pathSeparator, fileName))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("生成 %s\n", fileName)

			writer := bufio.NewWriter(resultFile)
			// 生成头
			fmt.Fprintln(writer, "import { ProColumns } from \"@ant-design/pro-components\";")
			fmt.Fprintln(writer, fmt.Sprintf("export const %sColumns: ProColumns<API.%s>[] = [", LowerCaseFirst(dataEntityName), dataEntityName))

			// 生成增加的字段
			if len(extendFields) != 0 {
				for _, extendField := range extendFields {
					columnData := getTemplate(templateMap, extendField, extendField)
					if columnData != "" {
						fmt.Fprintln(writer, columnData)
					}
				}
			}

			parseWithTypeSet := array2Set(parseWithTypes)
			// 遍历字段数组
			for _, fieldData := range v {
				fieldName := strings.ToLower(fieldData.Name)
				dataType := strings.ToLower(fieldData.DataType)
				// 排除字段
				if ignoreFieldSet.Contains(fieldName) {
					continue
				}

				var columnData string
				// 根据字段名处理的数据
				if parseWithFieldSet.Contains(fieldName) {
					columnData = getTemplate(templateMap, fieldName, fieldData.Name)
					// fmt.Println("field columnData ", columnData)
				} else {
					// 根据类型处理的数据
					if parseWithTypeSet.Contains(dataType) {
						columnData = getTemplate(templateMap, dataType, fieldData.Name)
						// fmt.Println("type columnData ", columnData)
					} else {
						columnData = getTemplate(templateMap, "default", fieldData.Name)
						// fmt.Println("type columnData ", columnData)
					}
				}
				if columnData != "" {
					fmt.Fprintln(writer, columnData)
				}
			}
			// 结束
			fmt.Fprintln(writer, "]")
			writer.Flush()

			// 页面
			// 生成index.tsx
			indexFileName := "index"
			indexContent := readFile(indexFileName)
			// 模板替换
			upperCaseFirst := UpperCaseFirst(lowerEntityName)
			data := fmt.Sprintf(indexContent,
				// 列导入
				LowerCaseFirst(dataEntityName), dataEntityName,
				// Service导入
				upperCaseFirst, upperCaseFirst, upperCaseFirst, upperCaseFirst, upperCaseFirst, dataEntityName,
				// 多行选择 useState
				dataEntityName,
				// columns 列
				dataEntityName, LowerCaseFirst(dataEntityName),
				// 详情页面路由
				lowerEntityName,
				// 根据id删除
				upperCaseFirst,
				// Service方法
				upperCaseFirst, upperCaseFirst, upperCaseFirst, upperCaseFirst,
			)
			CheckDirAndMk("Pages", dataEntityName)
			createAndWriteFile(indexFileName, dataEntityName, data)

			// 生成Detail.tsx
			detailFileName := "Detail"
			detailContent := readFile(detailFileName)
			// 模板替换
			detailData := fmt.Sprintf(detailContent,
				// 列导入
				LowerCaseFirst(dataEntityName), dataEntityName,
				upperCaseFirst, dataEntityName,
				dataEntityName, dataEntityName,
				upperCaseFirst,
				LowerCaseFirst(dataEntityName),
			)
			createAndWriteFile(detailFileName, dataEntityName, detailData)

			// 路由文件
			for k := range routerMap {
				data := getTemplate(routerMap, k, dataEntityName)
				fmt.Fprintln(routerWriter, data)
			}
		}
	}

	fmt.Fprintln(routerWriter, "]")
	routerWriter.Flush()

	fmt.Println("done.")
	fmt.Printf("耗时: %v\n", time.Since(now))
}

func array2Map(colTemplates []Template) map[string]string {
	datamap := map[string]string{}
	for _, v := range colTemplates {
		datamap[v.Key] = v.Value
	}
	return datamap
}

func createAndWriteFile(fileName string, dataEntityName string, data string) {
	fileName = fmt.Sprintf("Pages%s%s%s%s.tsx", pathSeparator, dataEntityName, pathSeparator, fileName)
	_, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(fileName, []byte(data), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func readFile(fileName string) string {
	filePath := fmt.Sprintf("BasePage%s%s.tsx", pathSeparator, fileName)
	f, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return string(f)
}

func getTemplate(templateMap map[string]string, name, replaceKey string) string {
	if value, ok := templateMap[name]; ok {
		// 有占位符
		count := strings.Count(value, "%s")
		if count > 0 {
			return strings.ReplaceAll(value, "%s", replaceKey)
		} else {
			return value
		}
	}
	return ""
}

func array2Set(fields []string) mapset.Set {
	set := mapset.NewSet()
	for _, v := range fields {
		set.Add(v)
	}
	return set
}

func parseFieldData(line string) *FieldData {
	param, paramType := trimData(line)
	return &FieldData{
		Name:     param,
		DataType: paramType,
	}
}

func trimData(line string) (string, string) {
	all := strings.ReplaceAll(line, "?", "")
	all = strings.ReplaceAll(all, ": ", ":")
	all = strings.ReplaceAll(all, ";", "")
	split := strings.Split(all, ":")
	param := split[0]
	paramType := split[1]
	return param, paramType
}

func CheckDirAndMk(rootPath, dirName string) {
	_, err := os.Stat(rootPath + pathSeparator + dirName)
	if err != nil {
		os.MkdirAll(rootPath+pathSeparator+dirName, os.ModePerm)
	}
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("read config failed: %v", err)
	}
	err = viper.Unmarshal(&conf)
}

func LowerCaseFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func UpperCaseFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func WithLowerCaseFirst(str string) bool {
	if str[0] >= 97 && str[0] <= 122 {
		return true
	}
	return false
}
