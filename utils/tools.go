package utils

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
)

// FolderExists checks if a folder exists at the given path.
func FolderExists(foldername string) bool {
	foldername = NormalizePath(foldername)
	if _, err := os.Stat(foldername); os.IsNotExist(err) {
		return false
	}
	return true
}

// NormalizePath the path
func NormalizePath(path string) string {
	if strings.HasPrefix(path, "~") {
		path, _ = homedir.Expand(path)
	}
	return path
}

// MakeDir create a folder
func MakeDir(folder string) {
	folder = NormalizePath(folder)
	os.MkdirAll(folder, 0750)
}

// GetOSEnv 获取环境变量
func GetOSEnv(name string, defaultValue string) string {
	variable, ok := os.LookupEnv(name)
	if !ok {
		if defaultValue != "" {
			return defaultValue
		}
		return name
	}
	return variable
}

func GenHash(text string) string {
	h := sha1.New()
	h.Write([]byte(text))
	hashed := h.Sum(nil)
	return fmt.Sprintf("%x", hashed)
}
func GetTS() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func GetPublicIP() string {
	var PublicIP string
	url := "https://ipinfo.io/ip"
	// 创建带有超时的客户端
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err == nil {
		defer resp.Body.Close()
		if ip, ok := io.ReadAll(resp.Body); ok == nil {
			PublicIP = string(ip)
			return PublicIP
		}
	}
	// 如果出现错误（如超时、网络问题），将IP设置为 localhost
	PublicIP = "localhost"
	return PublicIP
}

func FileExists(filename string) bool {
	filename = NormalizePath(filename)
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// ReadingFileUnique
func ReadingFileUnique(filename string) []string {
	var result []string
	if strings.Contains(filename, "~") {
		filename, _ = homedir.Expand(filename)
	}
	file, err := os.Open(filename)
	if err != nil {
		return result
	}
	defer file.Close()

	seen := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val := strings.TrimSpace(scanner.Text())
		// unique stuff
		if val == "" {
			continue
		}
		if seen[val] {
			continue
		}

		seen[val] = true
		result = append(result, val)
	}

	if err := scanner.Err(); err != nil {
		return result
	}
	return result
}

// CleanPath 从目标中提取出域名 https://example.com/path/to/resource → 输出 example_com
// /var/log/app.log（存在） → 输出 app_log。
func CleanPath(raw string) string {
	var out string
	raw = NormalizePath(raw)
	base := raw
	if FileExists(base) {
		base = filepath.Base(raw)
	}

	if strings.Count(base, "/") > 2 {
		base = base[strings.LastIndex(base, "/")+1:]
		if strings.TrimSpace(base) == "" {
			domain, err := GetDomain(raw)
			if err == nil {
				base = domain
			} else {
				base = RandomString(8)
			}
		}
	}

	out = strings.ReplaceAll(base, "/", "_")
	out = strings.ReplaceAll(out, ":", "_")
	// DebugF("CleanPath: %s -- %s", raw, out)
	return out
}

// GetDomain 从rul中提取域名
func GetDomain(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err == nil {
		return u.Hostname(), nil
	}
	return raw, err
}

// RandomString 返回指定长度的随机字符串
func RandomString(n int) string {
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	var letter = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[seededRand.Intn(len(letter))]
	}
	return string(b)
}

func WriteToFile(filename string, data string) (string, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.WriteString(file, data+"\n")
	if err != nil {
		return "", err
	}
	return filename, file.Sync()
}
