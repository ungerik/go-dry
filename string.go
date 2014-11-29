package dry

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
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

// StringMD5Hex returns the hex encoded MD5 hash of data
func StringMD5Hex(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// StringSHA1Base64 returns the base64 encoded SHA1 hash of data
func StringSHA1Base64(data string) string {
	hash := sha1.Sum([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:])
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

func StringFormatBigInt(mem uint64) string {
	switch {
	case mem >= 10e12:
		return fmt.Sprintf("%dT", mem/1e12)
	case mem >= 1e12:
		return strings.TrimSuffix(fmt.Sprintf("%.1fT", float64(mem)/1e12), ".0")

	case mem >= 10e9:
		return fmt.Sprintf("%dG", mem/1e9)
	case mem >= 1e9:
		return strings.TrimSuffix(fmt.Sprintf("%.1fG", float64(mem)/1e9), ".0")

	case mem >= 10e6:
		return fmt.Sprintf("%dM", mem/1e6)
	case mem >= 1e6:
		return strings.TrimSuffix(fmt.Sprintf("%.1fM", float64(mem)/1e6), ".0")

	case mem >= 10e3:
		return fmt.Sprintf("%dk", mem/1e3)
	case mem >= 1e3:
		return strings.TrimSuffix(fmt.Sprintf("%.1fk", float64(mem)/1e3), ".0")
	}
	return fmt.Sprintf("%d", mem)
}

func StringFormatMemory(mem uint64) string {
	return StringFormatBigInt(mem) + "B"
}

func StringReplaceMulti(str string, fromTo ...string) string {
	if len(fromTo)%2 != 0 {
		panic("Need even number of fromTo arguments")
	}
	for i := 0; i < len(fromTo); i += 2 {
		str = strings.Replace(str, fromTo[i], fromTo[i+1], -1)
	}
	return str
}

func StringToUpperCamelCase(str string) string {
	var buf bytes.Buffer
	var last byte = '_'
	for _, c := range []byte(str) {
		if c != '_' {
			if last == '_' {
				c = byte(unicode.ToUpper(rune(c)))
			} else {
				c = byte(unicode.ToLower(rune(c)))
			}
			buf.WriteByte(c)
		}
		last = c
	}
	return buf.String()
}

func StringToLowerCamelCase(str string) string {
	var buf bytes.Buffer
	var last byte
	for _, c := range []byte(str) {
		if c != '_' {
			if last == '_' {
				c = byte(unicode.ToUpper(rune(c)))
			} else {
				c = byte(unicode.ToLower(rune(c)))
			}
			buf.WriteByte(c)
		}
		last = c
	}
	return buf.String()
}

func StringMapSortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
