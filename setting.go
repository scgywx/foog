package foog

import(
	"fmt"
	"database/sql"
	"strings"
	"errors"
	"os"
	"log"
	"sync"
	_ "github.com/go-sql-driver/mysql"
)

type SettingRule struct{
	Name string
	Table string
	Index string
	Value string
	Field string
	Where string
	Order string
	FieldSkip string
}

type column struct {
	Field string
	Type string
	Table string
	Comment string
}

type desc struct{
	Name string
	Table string
	Where string
	Order string
	Fields []string
	Indexs []string
	Values []string
	FieldSkip []string
	Columns []column
	BaseType string
	BaseTypeCode string
	IndexType string
	IndexTypeList []string
}

var typesMap = map[string]string{
	"int":                "int",
	"integer":            "int",
	"tinyint":            "int",
	"smallint":           "int",
	"mediumint":          "int",
	"bigint":             "int",
	"int unsigned":       "int",
	"integer unsigned":   "int",
	"tinyint unsigned":   "int",
	"smallint unsigned":  "int",
	"mediumint unsigned": "int",
	"bigint unsigned":    "int",
	"bit":                "int",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"date":               "string",
	"datetime":           "string",
	"timestamp":          "string",
	"time":               "string",
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string",
	"varbinary":          "string",
}

var (
	dbHost string
	dbUser string
	dbPass string
	dbName string
	dbInst *sql.DB
	mutex sync.RWMutex
	indexs map[string]interface{}
)

func SetSetting(host, user, pass, dbname string){
	dbHost = host
	dbUser = user
	dbPass = pass
	dbName = dbname
	indexs = make(map[string]interface{})
}

func GetSettingDB()(*sql.DB, error){
	if dbInst != nil{
		return dbInst, nil
	}

	var (
		dsn string
		err error
	)

	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbUser, dbPass, dbHost, dbName)
	dbInst, err = sql.Open("mysql", dsn)
	if err != nil{
		return nil, err
	}

	return dbInst, nil
}

func CleanSetting(){
	mutex.Lock()
	defer mutex.Unlock()

	indexs = make(map[string]interface{})
}

func StoreSetting(key string, val interface{}){
	mutex.Lock()
	defer mutex.Unlock()

	indexs[key] = val
}

func LoadSetting(key string)(interface{}, bool){
	mutex.RLock()
	defer mutex.RUnlock()

	v, found := indexs[key]
	if !found {
		log.Println("setting miss", key)
	}

	return v, found
}

func MakeSetting(path string, rules []SettingRule)bool{
	db, err := GetSettingDB()
	if err != nil{
		log.Println("[ERROR]", err)
		return false
	}

	for _, r := range rules{
		//query table schema
		query := fmt.Sprintf(`SELECT COLUMN_NAME,DATA_TYPE,TABLE_NAME,COLUMN_COMMENT 
			FROM information_schema.COLUMNS 
			WHERE table_schema = DATABASE() AND TABLE_NAME='%s'`, r.Table)
		rows, err := db.Query(query)
		if err != nil{
			log.Println("[ERROR]", "query setting failed, name=" + r.Name + ", sql=" + query)
			return false
		}

		defer rows.Close()

		//init desc
		d := &desc{}
		d.Name = r.Name
		d.Table = r.Table
		d.Where = r.Where
		d.Order = r.Order

		if len(r.Field) > 0{
			d.Fields = strings.Split(r.Field, ",")
		}

		if len(r.Index) > 0{
			d.Indexs = strings.Split(r.Index, ",")
		}

		if len(r.Value) > 0{
			d.Values = strings.Split(r.Value, ",")
		}

		if len(r.FieldSkip) > 0{
			d.FieldSkip = strings.Split(r.FieldSkip, ",")
		}

		//fetch columns
		d.Columns = make([]column, 0, 10)
		for rows.Next() {
			col := column{}
			if err := rows.Scan(&col.Field, &col.Type, &col.Table, &col.Comment); err != nil {
				log.Println("[ERROR]", err)
				return false
			}

			if len(d.Fields) > 0 && !inArray(d.Fields, col.Field){
				continue
			}

			if len(d.FieldSkip) > 0 && inArray(d.FieldSkip, col.Field){
				continue
			}

			d.Columns = append(d.Columns, col)
		}

		//make base struct
		err = genStruct(d)
		if err != nil{
			log.Println("[ERROR] gen struct failed", err)
			return false
		}

		//make index struct
		err = genIndex(d)
		if err != nil{
			log.Println("[ERROR] gen index failed", err)
			return false
		}

		//make go file
		code, err := genCode(d)
		if err != nil{
			log.Println("[ERROR] gan code failed", err)
			return false
		}

		f, err := os.OpenFile(path + "/" + r.Name + ".go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil{
			log.Println("[ERROR] write file failed", err)
			return false
		}

		defer f.Close()
		f.WriteString(code)

		log.Println("generate succed", r.Name)
	}

	return  true
}

func genStruct(d *desc)error{
	d.BaseType = "Type_" + d.Name
	d.BaseTypeCode = "type " + d.BaseType + " struct{\n"
	for _, col := range d.Columns {
		colName := ucfirst(col.Field)
		d.BaseTypeCode+= "\t" + colName + "\t\t" + typesMap[col.Type] + " //" + col.Comment + "\n"
	}
	d.BaseTypeCode+= "}"

	return nil
}

func genIndex(d *desc)error{
	if len(d.Indexs) < 1{
		return errors.New("Invalid index, name=" + d.Name)
	}

	d.IndexTypeList = make([]string, len(d.Indexs))
	for k, v := range d.Indexs {
		col, ok := getColumnByName(d.Columns, v)
		if !ok {
			return errors.New("field not found by index:" + v)
		}

		d.IndexTypeList[k] = "map[" + typesMap[col.Type] + "]"
		d.IndexType+= "map[" + typesMap[col.Type] + "]"
	}

	d.IndexType+= d.BaseType

	return nil
}

func genCode(d *desc)(string, error){
	field := ""
	sep := ""
	for _, col := range d.Columns{
		field+= sep + "`" + col.Field + "`"
		sep = ","
	}

	where := ""
	if d.Where != ""{
		where = " WHERE " + d.Where
	}

	order := ""
	if d.Order != "" {
		order = " ORDER BY " + d.Order
	}

	indexchecker := ""
	indexlevel := ""
	for k, v := range d.Indexs {
		colName := ucfirst(v)
		indexlevel+= "[row." + colName + "]"

		if k < len(d.Indexs) - 1 {
			childtype := ""
			for i := k + 1; i < len(d.Indexs); i++{
				childtype+= d.IndexTypeList[i]
			}
			childtype+= d.BaseType

			indexchecker+= "\n\t\tif _, found := ret" + indexlevel + "; !found{\n"
			indexchecker+= "\t\t\tret" + indexlevel + " = make(" + childtype + ")\n"
			indexchecker+= "\t\t}\n"
		}
	}

	scanstr := ""
	fieldNullType := ""
	fieldNullFill := ""
	sep = ""
	for _, col := range d.Columns{
		colName := ucfirst(col.Field)
		colType := typesMap[col.Type]
		tmpName := "_" + ucfirst(col.Field)
		tmpType := ucfirst(colType)
		if colType == "int"{
			tmpType+= "64"
		}
		fieldNullType+= "\t\tvar " + tmpName + " sql.Null" + tmpType + "\n"
		scanstr+= sep + "&" + tmpName
		if colType == "int"{
			fieldNullFill+= "\t\trow." + colName + " = int(" + tmpName+"."+tmpType+")\n"
		}else{
			fieldNullFill+= "\t\trow." + colName + " = " + tmpName+"."+tmpType+"\n"
		}
		
		sep = ", "
	}

	tags := map[string]string{
		"name": d.Name,
		"basetype": d.BaseType,
		"basestruct": d.BaseTypeCode,
		"indextype": d.IndexType,
		"table": "`" + d.Table + "`",
		"field": field,
		"where": where,
		"order": order,
		"scanstr": scanstr,
		"indexchecker": indexchecker,
		"indexlevel": indexlevel,
		"fieldnulltype": fieldNullType,
		"fieldnullfill": fieldNullFill,
	}

	temp := 
`//
//WARNING!!!!!Dont modify this file, it is auto generated by tool
//

package setting

import(
	"github.com/scgywx/foog"
	"database/sql"
)

${basestruct}
	
func Load_${name}()(${indextype}, error){
	v, found := foog.LoadSetting("${name}")
	if found {
		return v.(${indextype}), nil
	}

	db, err := foog.GetSettingDB()
	if err != nil{
		return nil, err
	}

	rows, err := db.Query("SELECT ${field} FROM ${table} ${where} ${order}")
	if err != nil{
		return nil, err
	}

	var ret = make(${indextype})
	for rows.Next() {
${fieldnulltype}
		if err := rows.Scan(${scanstr}); err != nil {
			return nil, err
		}

		var row ${basetype}
${fieldnullfill}
${indexchecker}
		ret${indexlevel} = row
	}

	foog.StoreSetting("${name}", ret)
	
	return ret, nil
}`

	for k, v := range tags{
		temp = strings.ReplaceAll(temp, "${"+k+"}", v)
	}

	return temp, nil
}

func getColumnByName(fields []column, name string)(column, bool){
	for _, v := range fields{
		if name == v.Field{
			return v, true
		}
	}

	return column{}, false
}

func inArray(arr []string, s string)bool{
	for _,v := range arr{
		if v == s{
			return true
		}
	}

	return false
}

func ucfirst(str string)string{
	if len(str) < 1{
		return ""
	}

	bytes := []byte(str)
	if bytes[0] >= 97 && bytes[0] <= 122{
		bytes[0]-= 32
		str = string(bytes)
	}

	return str
}