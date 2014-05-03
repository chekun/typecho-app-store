package ziputil

//nroe/ziputil
//forked from yanolab/ziputil

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ZipFile struct {
	writer *zip.Writer
}

func Create(filename string) (*ZipFile, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &ZipFile{writer: zip.NewWriter(file)}, nil
}

func (z *ZipFile) Close() error {
	return z.writer.Close()
}

func (z *ZipFile) AddEntryN(path string, names ...string) error {
	for _, name := range names {
		zipPath := filepath.Join(path, name)
		err := z.AddEntry(zipPath, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (z *ZipFile) AddEntry(path, name string) error {
	fi, err := os.Stat(name)
	if err != nil {
		return err
	}

	fh, err := zip.FileInfoHeader(fi)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		if len(path) > 0 {
			path = path + "/"
		} else {
			path = "./"
		}
	}

	fh.Name = path

	entry, err := z.writer.CreateHeader(fh)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil
	}

	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(entry, file)

	return err
}

func (z *ZipFile) AddDirectoryN(path string, names ...string) error {
	for _, name := range names {
		err := z.AddDirectory(path, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (z *ZipFile) AddDirectory(path, dirName string) error {
	z.AddEntry(path, dirName)

	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return err
	}

	for _, file := range files {
		localPath := filepath.Join(dirName, file.Name())
		zipPath := filepath.Join(path, file.Name())

		err = nil
		if file.IsDir() {
			err = z.AddDirectory(zipPath, localPath)
		} else {
			err = z.AddEntry(zipPath, localPath)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func Zip(zipFile, directory string) error {

	os.Mkdir(filepath.Dir(zipFile), 0755)

	zip, err := Create(zipFile)
	if err != nil {
		return err
	}
	err = zip.AddDirectory("./", directory)
	if err != nil {
		return err
	}
	err = zip.AddEntry("./.app_store.txt", filepath.Dir(filepath.Dir(zipFile))+"/../source/.app_store.txt")
	if err != nil {
		return err
	}
	err = zip.Close()
	if err != nil {
		return err
	}
	return nil
}
