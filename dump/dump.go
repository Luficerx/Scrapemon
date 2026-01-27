package dump

import (
	"log"
	"os"
)

const (
	OUTPUT_FOLDER string = "./output/"
)

func Dump(file_name string, content []byte) {
	err := os.MkdirAll(OUTPUT_FOLDER, os.ModePerm)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	output_path := OUTPUT_FOLDER + file_name

	err = os.WriteFile(output_path, content, 0644)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	} else {
		log.Println("[OUTPUT]", output_path)
	}
}

func Read(file_name string, container *[]byte) error {
	content, err := os.ReadFile(file_name)
	*container = content
	return err
}

func Out(filepath string) string {
	return OUTPUT_FOLDER + filepath
}
