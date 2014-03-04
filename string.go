package dry

import (
	"bytes"
	"crypto/md5"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func StringPrettifyJSON(compactJSON string) string {
	var buf bytes.Buffer
	err := json.Indent(&buf, []byte(compactJSON), "", "\t")
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

func StringEscapeJSON(jsonString string) string {
	jsonString = strings.Replace(jsonString, `\`, `\\`, -1)
	jsonString = strings.Replace(jsonString, `"`, `\"`, -1)
	return jsonString
}

func StringStripHTMLTags(text string) (plainText string) {
	chars := []byte(text)
	tagStart := -1
	for i := 0; i < len(chars); i++ {
		if chars[i] == '<' {
			tagStart = i
		} else if chars[i] == '>' && tagStart != -1 {
			chars = append(chars[:tagStart], chars[i+1:]...)
			i, tagStart = tagStart-1, -1
		}
	}
	return string(chars)
}

func StringMD5(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func StringAddURLParam(url, name, value string) string {
	var separator string
	if strings.IndexRune(url, '?') == -1 {
		separator = "?"
	} else {
		separator = "&"
	}
	return url + separator + name + "=" + value
}

func StringConvertTime(timeString, formatIn, formatOut string) (resultTime string, err error) {
	if timeString == "" {
		return "", nil
	}
	t, err := time.Parse(formatIn, timeString)
	if err != nil {
		return "", err
	}
	return t.Format(formatOut), nil
}

func StringCSV(records [][]string) string {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	err := writer.WriteAll(records)
	if err != nil {
		return ""
	}
	return buf.String()
}

func StringToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}

func StringToFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func StringToBool(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

func StringInSlice(s string, slice []string) bool {
	for i := range slice {
		if slice[i] == s {
			return true
		}
	}
	return false
}

// StringJoinFormat formats every value in values with format
// and joins the result with sep as separator.
// values must be a slice of a formatable type
func StringJoinFormat(format string, values interface{}, sep string) string {
	v := reflect.ValueOf(values)
	if v.Kind() != reflect.Slice {
		panic("values is not a slice")
	}
	var buffer bytes.Buffer
	for i := 0; i < v.Len(); i++ {
		if i > 0 {
			buffer.WriteString(sep)
		}
		buffer.WriteString(fmt.Sprintf(format, v.Index(i).Interface()))
	}
	return buffer.String()
}

// StringJoin formats every value in values according to its default formatting
// and joins the result with sep as separator.
// values must be a slice of a formatable type
func StringJoin(values interface{}, sep string) string {
	v := reflect.ValueOf(values)
	if v.Kind() != reflect.Slice {
		panic("values is not a slice")
	}
	var buffer bytes.Buffer
	for i := 0; i < v.Len(); i++ {
		if i > 0 {
			buffer.WriteString(sep)
		}
		buffer.WriteString(fmt.Sprint(v.Index(i).Interface()))
	}
	return buffer.String()
}

func StringFormatMemory(mem uint64) string {
	switch {
	case mem >= 10e12:
		return fmt.Sprintf("%.1fTB", float64(mem)/10e12)
	case mem >= 10e9:
		return fmt.Sprintf("%.1fGB", float64(mem)/10e9)
	case mem >= 10e6:
		return fmt.Sprintf("%.1fMB", float64(mem)/10e6)
	case mem >= 10e3:
		return fmt.Sprintf("%.1fkB", float64(mem)/10e3)
	}
	return fmt.Sprintf("%dB", mem)
}
