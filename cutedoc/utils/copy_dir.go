package utils

import (
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func copyFile(src string, dst string) error {
	log.Println("copying", src, "to", dst)

	err := os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return err
	}

	inFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, inFile)
	if err != nil {
		return err
	}

	return nil
}

func CopyDirContents(srcDir string, dstDir string, predicate func(ext string) bool) error {
	srcDirAbs, err := filepath.Abs(srcDir)
	if err != nil {
		return err
	}

	dstDirAbs, err := filepath.Abs(dstDir)
	if err != nil {
		return err
	}

	return filepath.WalkDir(srcDirAbs, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || !predicate(filepath.Ext(path)) {
			return err
		}

		relativePath := path[len(srcDirAbs):]
		fileDstPath := filepath.Join(dstDirAbs, relativePath)

		err = copyFile(path, fileDstPath)
		if err != nil {
			return err
		}

		return err
	})
}
