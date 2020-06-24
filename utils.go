package sqlbundle

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

const TIME_FORMAT = "20060102150405"

func makeTimeSequence() string {
	t := time.Now()
	return t.Format(TIME_FORMAT)
}

func downloadDependency(depDir, link string) (string, error) {
	fileName := path.Base(link)
	printInfo("download dependency", link)
	// Create the file
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
			//printInfo(fmt.Sprintf("Contents of %s: ", hdr.Name))
			//if _, err := io.Copy(os.Stdout, tr); err != nil {
			//	return err
			//}
			// determine proper file path info
			finfo := hdr.FileInfo()
			fileName := hdr.Name
			absFileName := filepath.Join(dest, fileName)
			// if a dir, create it, then go to next segment
			if finfo.Mode().IsDir() {
				if err := os.MkdirAll(absFileName, 0755); err != nil {
					return err
				}
				continue
			}
			// create new file with original file mode
			file, err := os.OpenFile(
				absFileName,
				os.O_RDWR|os.O_CREATE|os.O_TRUNC,
				finfo.Mode().Perm(),
			)
			if err != nil {
				return err
			}
			printInfo(fmt.Sprintf("x %s\n", absFileName))
			n, cpErr := io.Copy(file, tr)
			if closeErr := file.Close(); closeErr != nil {
				return err
			}
			if cpErr != nil {
				return cpErr
			}
			if n != finfo.Size() {
				return fmt.Errorf("wrote %d, want %d", n, finfo.Size())
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
