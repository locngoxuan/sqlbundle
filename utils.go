package sqlbundle

import (
	"archive/tar"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const TIME_FORMAT = "20060102150405"
const (
	PARSER_START      int = iota // 0
	PARSER_UP_BEGIN              // 1
	PARSER_UP_END                // 2
	PARSER_DOWN_BEGIN            // 3
	PARSER_DOWN_END              // 4
)

func makeTimeSequence() string {
	t := time.Now()
	return t.Format(TIME_FORMAT)
}

func isEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func isNillOrEmpty(s *string) bool {
	if s == nil {
		return true
	}
	return isEmpty(*s)
}

func downloadDependency(depDir, link string) (string, error) {
	printInfo("download dependency", link)
	// Create the file
	fileName := path.Base(link)
	dest := filepath.Join(depDir, fileName)
	out, err := os.Create(dest)
	if err != nil {
		return dest, err
	}
	defer func() {
		_ = out.Close()
	}()

	// Get the data
	resp, err := http.Get(link)
	if err != nil {
		return dest, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return dest, errors.New(fmt.Sprintf("bad status: %s", resp.Status))
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return dest, err
	}

	return dest, nil
}

func tarFile(tarPath string, paths []string) error {
	tarFile, err := os.Create(tarPath)
	if err != nil {
		return err
	}
	defer func() {
		err = tarFile.Close()
	}()

	absTar, err := filepath.Abs(tarPath)
	if err != nil {
		return err
	}

	tw := tar.NewWriter(tarFile)
	defer func() {
		_ = tw.Close()
	}()

	// walk each specified path and add encountered file to tar
	for _, path := range paths {
		// validate path
		path = filepath.Clean(path)
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		if absPath == absTar {
			return errors.New(fmt.Sprintf("tar file %s cannot be the source", tarPath))
		}
		if absPath == filepath.Dir(absTar) {
			return errors.New(fmt.Sprintf("tar file %s cannot be in source %s", tarPath, absPath))
		}

		walker := func(file string, finfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// fill in header info using func FileInfoHeader
			hdr, err := tar.FileInfoHeader(finfo, finfo.Name())
			if err != nil {
				return err
			}

			relFilePath := file
			if filepath.IsAbs(path) {
				relFilePath, err = filepath.Rel(path, file)
				if err != nil {
					return err
				}
			}
			// ensure header has relative file path
			hdr.Name = relFilePath

			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
			// if path is a dir, dont continue
			if finfo.Mode().IsDir() {
				return nil
			}

			// add file to tar
			srcFile, err := os.Open(file)
			if err != nil {
				return err
			}
			defer func() {
				_ = srcFile.Close()
			}()
			_, err = io.Copy(tw, srcFile)
			if err != nil {
				return err
			}
			return nil
		}

		// build tar
		if err := filepath.Walk(path, walker); err != nil {
			return err
		}
	}
	return nil
}

func untarFile(tarPath, dest string) error {
	tarUnwrite := func(file string) error {
		tarFile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer func() {
			_ = tarFile.Close()
		}()
		tr := tar.NewReader(tarFile)
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break // End of archive
			}
			if err != nil {
				return err
			}
			fileInfo := hdr.FileInfo()
			fileName := hdr.Name
			absFileName := filepath.Join(dest, fileName)
			// if a dir, create it, then go to next segment
			if fileInfo.Mode().IsDir() {
				if err := os.MkdirAll(absFileName, 0755); err != nil {
					return err
				}
				continue
			}
			// create new file with original file mode
			file, err := os.OpenFile(
				absFileName,
				os.O_RDWR|os.O_CREATE|os.O_TRUNC,
				fileInfo.Mode().Perm(),
			)
			if err != nil {
				return err
			}
			printInfo(fmt.Sprintf("extract %s", absFileName))
			n, cpErr := io.Copy(file, tr)
			if closeErr := file.Close(); closeErr != nil {
				return err
			}
			if cpErr != nil {
				return cpErr
			}
			if n != fileInfo.Size() {
				return fmt.Errorf("wrote %d, want %d", n, fileInfo.Size())
			}
		}
		return nil
	}
	if err := tarUnwrite(tarPath); err != nil {
		return err
	}
	return nil
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = source.Close()
	}()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		_ = destination.Close()
	}()
	_, err = io.Copy(destination, source)
	return err
}

func copyDirectory(scrDir, dest string) error {
	entries, err := ioutil.ReadDir(scrDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(scrDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := createIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := copyDirectory(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := copyFile(sourcePath, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func createIfNotExists(dir string, perm os.FileMode) error {
	if exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func isDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

//deploy
func uploadFile(link, source, user, pass string) error {
	data, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		_ = data.Close()
	}()
	req, err := http.NewRequest("PUT", link, data)
	if err != nil {
		return err
	}
	md5, err := md5sum(source)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("X-CheckSum-MD5", md5)
	req.SetBasicAuth(user, pass)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode != http.StatusCreated {
		return errors.New(res.Status)
	}
	return nil
}

func md5sum(file string) (string, error) {
	hasher := md5.New()
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
