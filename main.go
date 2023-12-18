package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/NublyBR/go-vpk"
)

func main() {
    vpkFilePath, output, logPath := getFlags()	
	enableLogging := len(logPath) > 0
	
	logger := CustomLogger(os.Stdout, "INFO: ", log.Ldate|log.Ltime, enableLogging)
	
	
	if enableLogging {
		logfile, err  := os.Create(filepath.Join(logPath,"vpk-extractor-"+time.Now().Format("2006-01-02-15-04-05")+".log"))
	
		if err != nil {
			fmt.Println(err.Error())
		}

		defer logfile.Close()
		logger.SetOutput(logfile)
	} else {
		logger.SetOutput(nil)
	}
	
	pak, err := vpk.OpenDir(vpkFilePath)

	if err != nil {
		logger.Fatal(err)
	}

	defer pak.Close()

	// Iterate through all files in the VPK
	for _, file := range pak.Entries() {
		entry, err := file.Open()
		if err != nil {
			logger.Fatal(err)
		}
		
		path := filepath.Join(output,fileNameWithoutExtension(filepath.Base(vpkFilePath)),file.Filename())

		// Ensure the directories exist by using os.MkdirAll
		dir := filepath.Dir(path)
		if dir_err := os.MkdirAll(dir, 0755); dir_err != nil {
			logger.Fatal(dir_err)
		}
		logger.Println("Extract file",path)
		writeErr := ExtractVpkFile(entry, path)
		if writeErr != nil {
			logger.Fatal(err)
		}	
	}
	logger.Println("Successfully extracted content to "+vpkFilePath)
}




func ExtractVpkFile(file vpk.FileReader, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	io.Copy(f, file)

	f.Sync()
	closeErr := f.Close()
	if closeErr != nil {
		return closeErr
	}
	return nil
}

func fileNameWithoutExtension(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func getFlags() (string, string, string) {
	path := flag.String("p", "", "The full path to the vpk file.")
	output := flag.String("o", ".", "The path to output the files in, leave empty to generate the files in the same directory as the executable.")
	logPath := flag.String("logPath", "", "The path where the log will be written to, leave blank to disable the creation of a logfile.")
	
	flag.Parse()


	if *path == "" {
		panic("No .vpk file path was specified.")
	}

	_, err := os.Stat(*path)

	if err != nil {
		panic(err)
	}

	if filepath.Ext(*path) != ".vpk" {
		panic("File is not vpk format")
	}

	return *path, *output, *logPath
}