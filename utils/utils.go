package utils

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// Memory 结构体存储系统内存信息
type Memory struct {
	MemTotal     int // 总内存
	MemFree      int // 空闲内存
	MemAvailable int // 可用内存
}

// processLockFile 用于处理进程锁文件
type processLockFile struct {
	*os.File // 包含文件对象
}

// Unlock 解锁并删除锁文件
func (p processLockFile) Unlock() error {
	path := p.File.Name()
	if err := p.File.Close(); err != nil {
		return err
	}

	return os.Remove(path)
}

func AcquireProcessIDLock(pidFilePath string) (interface{ Unlock() error }, error) {
	if _, err := os.Stat(pidFilePath); !os.IsNotExist(err) {
		raw, err := os.ReadFile(pidFilePath)
		if err != nil {
			return nil, err
		}

		pid, err := strconv.Atoi(string(raw))
		if err != nil {
			return nil, err
		}

		if proc, err := os.FindProcess(int(pid)); err == nil && !errors.Is(proc.Signal(syscall.Signal(0)), os.ErrProcessDone) {
			fmt.Fprintf(os.Stderr, "Process %d is already running!\n", proc.Pid)

		} else if err = os.Remove(pidFilePath); err != nil {
			return nil, err

		}
	}

	f, err := os.OpenFile(pidFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	if _, err := f.Write([]byte(fmt.Sprint(os.Getpid()))); err != nil {
		return nil, err
	}

	return processLockFile{File: f}, nil
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// fill T.fields from V.fields, by same field name
// T.fields can be sub set of V.fields
func FillWith[T any, V any](dst *T, src *V) {
	dst_values := reflect.ValueOf(dst)
	dst_types := reflect.TypeOf(*dst)
	src_values := reflect.ValueOf(src)
	//src_types := src_values.Type()

	for i := 0; i < dst_types.NumField(); i++ {
		dst_field_name := dst_types.Field(i).Name
		src_field_value := reflect.Indirect(src_values).FieldByName(dst_field_name)

		if src_field_value.Kind() == reflect.Ptr && src_field_value.IsNil() {
			//log.Printf("%s is Nil", dst_field_name)
			continue
		}

		if src_field_value.IsValid() && !src_field_value.IsZero() {
			if src_field_value.Kind() == reflect.Ptr {
				reflect.Indirect(dst_values).Field(i).Set(src_field_value.Elem())
			} else {
				reflect.Indirect(dst_values).Field(i).Set(src_field_value)
			}
			//fmt.Printf("set %s to %v\n", dst_field_name, src_field_value.Interface())
		}
	}
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}

	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}

	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}

	//if err = os.Link(src, dst); err == nil {
	//	return
	//}

	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func GetEnvConf() map[string]string {
	confEnv, err := godotenv.Read() // .env in project root path.
	ErrorPanic(err)
	return confEnv
}

func StructToMap[S any](s *S) map[string]interface{} {
	var outInterface map[string]interface{}
	inStruct, _ := json.Marshal(s)
	json.Unmarshal(inStruct, &outInterface)
	return outInterface
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func parseLine(raw string) (key string, value int) {
	//fmt.Println(raw)
	text := strings.ReplaceAll(raw[:len(raw)-2], " ", "")
	keyValue := strings.Split(text, ":")
	return keyValue[0], toInt(keyValue[1])
}

func toInt(raw string) int {
	if raw == "" {
		return 0
	}
	res, err := strconv.Atoi(raw)
	if err != nil {
		panic(err)
	}
	return res
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v", m.NumGC)

	if runtime.GOOS == "linux" {
		var ms Memory = ReadMemoryStats()
		fmt.Printf("\tMemTotal = %v", ms.MemTotal)
		fmt.Printf("\tMemFree = %v", ms.MemFree)
		fmt.Printf("\tMemAvailable = %v\n", ms.MemAvailable)
	} else {
		fmt.Println("")
	}
}

func ReadMemoryStats() Memory {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bufio.NewScanner(file)
	scanner := bufio.NewScanner(file)
	res := Memory{}
	for scanner.Scan() {
		key, value := parseLine(scanner.Text())
		switch key {
		case "MemTotal":
			res.MemTotal = value
		case "MemFree":
			res.MemFree = value
		case "MemAvailable":
			res.MemAvailable = value
		}
	}
	return res
}

func GetMemUsageString() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return fmt.Sprintf("Alloc = %v MiB\tTotalAlloc = %v MiB\tSys = %v MiB\tNumGC = %v", bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
}
