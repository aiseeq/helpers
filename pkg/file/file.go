package file

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func PutContents(filename string, content string, mode int) error {
	// os.O_RDWR | os.O_CREATE
	fp, err := os.OpenFile(filename, mode, 0666)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.WriteString(content)

	return err
}

func Download(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// Copies file source to destination dest.
func Copy(source string, dest string) (err error) {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err == nil {
			err = os.Chmod(dest, si.Mode())
		}
	}
	return err
}

// Recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
func CopyDir(source string, dest string) (err error) {
	// get properties of source dir
	fi, err := os.Stat(source)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return errors.New("source is not a directory")
	}

	// create dest dir
	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(source)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sfp := source + "/" + entry.Name()
		dfp := dest + "/" + entry.Name()
		if entry.IsDir() {
			if err = CopyDir(sfp, dfp); err != nil {
				return err
			}
		} else {
			// perform copy
			if err = Copy(sfp, dfp); err != nil {
				return err
			}
		}

	}
	return err
}

// Попробовать установить защёлку
// Если старый демон жив и ttk = true, убить его
func TryToLock(lockFile string, ttk bool) (bool, error) {
	if Exists(lockFile) {
		// Защёлка существует
		data, err := ioutil.ReadFile(lockFile)
		if err != nil {
			// Что-то с файлом не так
			return false, err
		}
		pid, err := strconv.Atoi(string(data))
		if err != nil {
			// Не смогли достать pid
			return false, err
		}
		// В доке пишут, что ошибок не возвращает никогда
		process, _ := os.FindProcess(pid)
		// Проверим процесс на существование
		if err = process.Signal(syscall.Signal(0)); err == nil && !ttk {
			// Если процесс есть, и нам нельзя его убивать,
			// Не убиваем приложение
			return false, nil
		}
		process.Signal(syscall.Signal(syscall.SIGTERM))
		for x := 0; ; x++ {
			// Если тут ошибка, значит такого процесса уже нет
			if err = process.Signal(syscall.Signal(0)); err != nil {
				break
			}
			// Цилкл выполнится десяток раз, пока старый процесс не умрёт
			time.Sleep(1 * time.Microsecond)
			if x == 1000000 {
				// Умирает больше секунды, что-то не так
				return false, errors.New("previous deamon didn't died in 1 second after sigrerm")
			}
		}
	}
	// Убираем мёртвую защёлку, если есть
	os.Remove(lockFile)
	if err := PutContents(lockFile, strconv.Itoa(os.Getpid()), os.O_RDWR|os.O_CREATE); err != nil {
		return false, err
	}
	return true, nil
}
