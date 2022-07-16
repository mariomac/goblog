package tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

func CopyDir(src string, dest string) error {

	//var s is the path-separator. It will be changed, if the platform is Windows.
	s := "/"

	if runtime.GOOS == "windows" {
		s = `\`
	}

	dirinfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if dest[:len(src)] == src {
		return fmt.Errorf("Cannot copy a folder into the folder itself!")
	}

	f, err := os.Open(src)
	if err != nil {
		return err
	}

	file, err := f.Stat()
	if err != nil {
		return err
	}
	if !file.IsDir() {
		return fmt.Errorf("Source " + file.Name() + " is not a directory!")
	}

	err = os.Mkdir(dest, dirinfo.Mode())
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {

		if f.IsDir() {

			err = CopyDir(src+s+f.Name(), dest+s+f.Name())
			if err != nil {
				return err
			}

		}

		if !f.IsDir() {

			content, err := ioutil.ReadFile(src + s + f.Name())
			if err != nil {
				return err

			}

			err = ioutil.WriteFile(dest+s+f.Name(), content, f.Mode())
			if err != nil {
				return err

			}

		}

	}

	return nil
}
