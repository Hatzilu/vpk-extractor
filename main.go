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
    rootPath, logPath := getFlags()	
	logfile, err  := os.Create(filepath.Join(logPath,"vpk-extractor-"+time.Now().Format("2006-01-02-15-04-05")+".log"))
	defer logfile.Close()
	
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	logger.SetOutput(logfile)
	
	if logPath == "" {
		logger.SetOutput(logfile)
		} else {
		logger.SetOutput(nil)
	}

	if err != nil {
		logger.Fatal(err)
	}

	pak, err := vpk.OpenDir(rootPath)

	defer pak.Close()

	// Iterate through all files in the VPK
	for _, file := range pak.Entries() {
		entry, err := file.Open()
		if err != nil {
			logger.Fatal(err)
		}
		
		path := filepath.Join("output",file.Filename())

		// Ensure the directories exist by using os.MkdirAll
		dir := filepath.Dir(path)
		if dir_err := os.MkdirAll(dir, 0755); dir_err != nil {
			logger.Fatal(dir_err)
		}
		logger.Println("Writing file",path)
		fmt.Println("Writing file",path)
		writeErr := WriteVpkFile(entry, path)
		if writeErr != nil {
			logger.Fatal(err)
		}	
	}
	logger.Println("Successfully extracted content to "+rootPath)
	fmt.Println("Successfully extracted content to "+rootPath)
}



func WriteVpkFile(file vpk.FileReader, path string) error {
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

func getFlags() (string, string) {
	path := flag.String("p", "", "the full path to the vpk file.")
	logPath := flag.String("logPath", "", "the path where the log will be written to, leave blank to disable the creation of a logfile.")
	
	flag.Parse()

	fmt.Println("path is", *path)

	if *path == "" {
		panic("No .vpk file path was specified.")
	}

	_, err := os.Stat(*path)

	if err != nil {
		panic(err)
	}

	return *path, *logPath
}