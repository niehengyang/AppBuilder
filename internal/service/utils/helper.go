package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"math"
	"math/big"
	rand2 "math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// RoundFloat
//
//	@Description: float64保留小数位，支持保留到小数点后1、2、3、4位
//	@param f	原数据
//	@param n	需要保留的小数位
//	@return float64	结果
//	@return error	错误信息
func RoundFloat(f float64, n int) (float64, error) {
	var fstr string
	switch n {
	case 1:
		fstr = fmt.Sprintf("%.1f", f)
	case 2:
		fstr = fmt.Sprintf("%.2f", f)
	case 3:
		fstr = fmt.Sprintf("%.3f", f)
	case 4:
		fstr = fmt.Sprintf("%.4f", f)
	}
	f64, err := strconv.ParseFloat(fstr, 64)
	return f64, err
}

// RandFloat
//
//	@Description: 生成指定范围内的随机float64
//	@param min	最小值
//	@param max	最大值
//	@param n	要保留的小数位，支持1、2、3
//	@return float64	结果
func RandFloat(min, max float64, n int) float64 {
	rand2.Seed(time.Now().UnixNano())
	randNum := min + rand2.Float64()*(max-min)
	var f float64
	switch n {
	case 1:
		f, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", randNum), 64)
	case 2:
		f, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", randNum), 64)
	case 3:
		f, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", randNum), 64)
	}

	return f
}

// RandInt
//
//	@Description: 生成指定范围内的随机int64
//	@param min	最小值
//	@param max	最大值
//	@return int64	结果
func RandInt(min, max int64) int64 {
	if min > max {
		var str any = "the min is greater than max!"
		panic(str)
	}
	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))

		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}

// TryCatch
//
//	@Description: 异常捕捉
//	@param todo	要执行的业务程序函数体
//	@param catch
//	@return exception	异常结果
func TryCatch(todo func(), catch func()) (exception string) {
	defer func() {
		if err := recover(); err != any(nil) {
			exception = fmt.Sprintf("%v", err)
			catch() //捕捉到异常时的处理函数
		}
	}()
	todo() //要执行的业务程序函数体
	return exception
}

// HasDuplicateEntryError
//
//	@Description: 查找是否存在Duplicate entry（数据库重复写入相同数据）错误
//	@param error	数据库操作的错误源
//	@return bool	是否是数据库重复写入相同数据的错误
func HasDuplicateEntryError(error interface{}) bool {
	errStr := fmt.Sprintf("%s", error)
	index := strings.Index(errStr, "Duplicate entry")
	if index != -1 {
		return true
	}
	return false
}

// StructAssign
//
//	@Description: 将一个结构体的字段值赋给另一个结构体中相同的字段
//	@param target	目标结构体
//	@param source	当前结构体
func StructAssign(target interface{}, source interface{}) {
	bVal := reflect.ValueOf(target).Elem() //获取reflect.Type类型
	vVal := reflect.ValueOf(source).Elem()
	vTypeOfT := vVal.Type()
	for i := 0; i < vVal.NumField(); i++ {
		// 在要修改的结构体中查询有数据结构体中相同属性的字段，有则修改其值
		name := vTypeOfT.Field(i).Name
		//fmt.Println(name)
		//todo basemodel里的字段无法读取,
		//todo source 里的子集 结构体读取
		if ok := bVal.FieldByName(name).IsValid(); ok {
			//fmt.Println("存在，去赋值，",reflect.ValueOf(vVal.Field(i).Interface()))
			bVal.FieldByName(name).Set(reflect.ValueOf(vVal.Field(i).Interface()))
		}
	}
}

// getStrongPasswordString
//
//	@Description: 随机生成指定位数的大写字母和数字的组合
//	@param l 密码长度
//	@return string 密码
func GetStrongPasswordString(l int) string {
	//~!@#$%^&*()_+{}":?><;.,
	str := "123456789ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghijklmnpqrstuvwxyz!@#$%&*"
	bytes := []byte(str)
	result := []byte{}
	r := rand2.New(rand2.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	ok1, _ := regexp.MatchString(".[1|2|3|4|5|6|7|8|9]", string(result))
	ok2, _ := regexp.MatchString(".[Z|X|C|V|B|N|M|A|S|D|F|G|H|J|K|L|Q|W|E|R|T|Y|U|I|P]", string(result))
	ok3, _ := regexp.MatchString(".[z|x|c|v|b|n|m|a|s|d|f|g|h|j|k|l|q|w|e|r|t|y|u|i|p]", string(result))
	ok4, _ := regexp.MatchString(".[!|@|#|$|%|&|*]", string(result))
	if ok1 && ok2 && ok3 && ok4 {
		return string(result)
	} else {
		return GetStrongPasswordString(l)
	}
}

// Base64ToImage
//
//	@Description: 将base64编码数据转换成Image类型
//	@param base64St 图片Base64数据
//	@return image.Image	图片数据
//	@return error	错误处理
func Base64ToImage(base64Str string) (image.Image, error) {

	// 将 Base64 字符串解码为字节数组
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, errors.New("base64 decode failed: " + err.Error())
	}

	// 将字节数组解码为 Image 接口类型
	img, _, err := image.Decode(strings.NewReader(string(data)))
	if err != nil {
		return nil, errors.New("image decode failed: " + err.Error())
	}

	return img, nil
}

// Base64ToImageBuff
//
//	@Description: 将base64编码数据转换成Image字节类型
//	@param target	目标数据
//	@param source	当前图片
func Base64ToImageBuff(base64Str string) ([]byte, error) {
	// 将 Base64 字符串解码为字节数组
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, errors.New("base64 decode failed: " + err.Error())
	}

	return data, nil
}

const (
	fileNameLength = 10 // 文件名长度
)

// ImagePathToBase64
//
//	@Description: 通过图片地址转换成base64编码数据
//	@param target	当前图片
//	@param source	目标类型
func ImagePathToBase64(imagePath string) (string, error) {
	// 读取图像文件
	imageFile, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return "", errors.New("base64 encode failed: " + err.Error())
	}

	// 将图像文件内容编码为Base64字符串
	base64Encoded := base64.StdEncoding.EncodeToString(imageFile)
	return base64Encoded, nil
}

// GenerateRandomFileName
//
//	@Description: 生成不重复的名称
func GenerateRandomFileName() string {
	// 生成随机文件名
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand2.Seed(time.Now().UnixNano())
	b := make([]byte, fileNameLength)
	for i := range b {
		b[i] = letters[rand2.Intn(len(letters))]
	}
	return string(b)
}

// ParseGzipMsg
//
//	@Description: 解压通过gzip压缩的数据
//	@param source	目标类型
func ParseGzipMsg(source []byte) ([]byte, error) {

	// 使用bytes.NewReader创建一个Reader来读取压缩数据
	compressedDataReader := bytes.NewReader(source)

	// 创建一个gzip.Reader来解压数据
	gz, err := gzip.NewReader(compressedDataReader)
	if err != nil {
		fmt.Println("解压gzip数据时出错:", err)
		return nil, err
	}
	defer gz.Close()

	// 读取解压后的数据
	uncompressedData, err := ioutil.ReadAll(gz)
	if err != nil {
		fmt.Println("读取解压数据时出错:", err)
		return nil, err
	}

	return uncompressedData, nil

}
