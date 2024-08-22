package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	msgpack "github.com/vmihailenco/msgpack/v5"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func handleFunc(data *Data) *Data {
	switch {
	case data.EncryptionMethod == "nothing":
		path := data.EncryptedPath
		return getContent(Host + "/" + string(path) + "?encryption_method=nothing")
	case data.EncryptionMethod == "encoded as base64":
		path, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(data.EncryptedPath, "task_", ""))
		if err != nil {
			panic(err)
		}
		return getContent(Host + "/task_" + string(path) + "?encryption_method=nothing")
	case data.EncryptionMethod == "swapped every pair of characters":
		path := strings.ReplaceAll(data.EncryptedPath, "task_", "")
		path = swap(path)
		return getContent(Host + "/task_" + string(path) + "?encryption_method=nothing")
	case strings.Contains(data.EncryptionMethod, "circularly rotated left by "):
		path := strings.ReplaceAll(data.EncryptedPath, "task_", "")
		path = rotate(path, strings.ReplaceAll(data.EncryptionMethod, "circularly rotated left by ", ""))
		return getContent(Host + "/task_" + string(path) + "?encryption_method=nothing")
	case strings.Contains(data.EncryptionMethod, "encoded it with custom hex character set "):
		path := strings.ReplaceAll(data.EncryptedPath, "task_", "")
		path = hmacSha256Hex(path, strings.ReplaceAll(data.EncryptionMethod, "encoded it with custom hex character set ", ""))
		return getContent(Host + "/task_" + path + "?encryption_method=nothing")
	case strings.Contains(data.EncryptionMethod, "scrambled! original positions as base64 encoded messagepack: "):
		path := strings.ReplaceAll(data.EncryptedPath, "task_", "")
		path = positon(path, strings.ReplaceAll(data.EncryptionMethod, "scrambled! original positions as base64 encoded messagepack: ", ""))
		return getContent(Host + "/task_" + path + "?encryption_method=nothing")
	case strings.Contains(data.EncryptionMethod, "hashed with sha256, good luck"):
		path := strings.ReplaceAll(data.EncryptedPath, "task_", "")
		path = sha256De(path)
		return getContent(Host + "/task_" + path + "?encryption_method=nothing")
	}

	return nil
}

func sha256De(path string) string {
	message := []byte(path) //字符串转化字节数组
	//创建一个基于SHA256算法的hash.Hash接口的对象
	hash := sha256.New() //sha-256加密
	//hash := sha512.New() //SHA-512加密
	//输入数据
	hash.Write(message)
	//计算哈希值
	bytes := hash.Sum(nil)
	//将字符串编码为16进制格式,返回字符串
	hashCode := hex.EncodeToString(bytes)
	return hashCode
}

func positon(path, position string) string {
	messagePackData, err := base64.StdEncoding.DecodeString(position)
	if err != nil {
		panic(err)
	}

	var positions []int
	err = msgpack.Unmarshal(messagePackData, &positions)
	if err != nil {
		panic(err)
	}

	decryptedPath := make([]byte, len(path))
	for i, pos := range positions {
		decryptedPath[pos] = path[i]
	}

	return string(decryptedPath)
}

func hmacSha256Hex(path, hex string) string {
	standardHex := "0123456789abcdef" // The standard hexadecimal character set
	var decryptedPath strings.Builder

	for _, char := range path {
		index := strings.IndexRune(hex, char)
		if index == -1 {
			panic("character not found in custom hex set")
		}
		decryptedPath.WriteByte(standardHex[index])
	}

	return decryptedPath.String()
}

func rotate(str string, left string) string {
	leftInt, err := strconv.Atoi(left)
	if err != nil {
		panic(err)
	}
	num := len(str) - leftInt
	return str[num:] + str[:num]
}

func getContent(url string) *Data {
	fmt.Println(url)
	body, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(body.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	res := new(Data)
	err = json.Unmarshal(data, res)
	if err != nil {
		panic(err)
	}
	return res

}

func swap(str string) string {
	strArr := strings.Split(str, "")

	for i := 0; i < len(strArr); i = i + 2 {
		if i+1 >= len(strArr) {
			break
		}
		strArr[i], strArr[i+1] = strArr[i+1], strArr[i]
	}
	return strings.Join(strArr, "")
}
