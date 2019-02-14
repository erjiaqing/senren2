package taskworker

import (
	"archive/zip"
	"io"
	"os"

	"github.com/erjiaqing/senren2/pkg/types/base"
)

func compress(baseDir string, file *base.HomeworkArchiveDescriptor, output *zip.Writer) error {
	if file.Type == "file" {
		f, err := os.Open(file.Source)
		if err != nil {
			return err
		}
		defer f.Close()
		info, err := f.Stat()
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		header.Name = baseDir + "/" + file.Name
		if err != nil {
			return err
		}
		writer, err := output.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, f)
		if err != nil {
			return err
		}
	} else if file.Type == "archive" {
		for _, c := range file.Content {
			if err := compress(baseDir+"/"+file.Name, c, output); err != nil {
				return err
			}
		}
	} else if file.Type == "root" {
		for _, c := range file.Content {
			if err := compress(file.Name, c, output); err != nil {
				return err
			}
		}
	}
	return nil
}

func homeworkArchiveTask(uid string, t *base.HomeworkArchiveTask) error {
	os.MkdirAll("output/"+uid, os.ModePerm)
	d, _ := os.Create("output/" + uid + "/" + t.OutputFileName)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	return compress("", t.Desc, w)
}
