package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jlaffaye/ftp"
)

// uploadFile uploads a single file to the specified FTP path.
func uploadFile(ftpConn *ftp.ServerConn, localFile string, ftpPath string) error {
	file, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return ftpConn.Stor(ftpPath, file)
}

// uploadDirectory uploads all files and subdirectories in the specified local directory to the FTP server.
func uploadDirectory(ftpConn *ftp.ServerConn, localDir string, ftpDir string) error {
	files, err := ioutil.ReadDir(localDir)
	if err != nil {
		return err
	}

	// Create the directory on the FTP server.
	err = ftpConn.MakeDir(ftpDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		localFilePath := filepath.Join(localDir, file.Name())
		ftpFilePath := filepath.Join(ftpDir, file.Name())

		if file.IsDir() {
			// Recursively upload subdirectories.
			err := uploadDirectory(ftpConn, localFilePath, ftpFilePath)
			if err != nil {
				return err
			}
		} else {
			// Upload files.
			err := uploadFile(ftpConn, localFilePath, ftpFilePath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	// Connect to the FTP server.
	ftpConn, err := ftp.Dial("192.168.239.124:2121") // Corrected line to dial the FTP server.
	if err != nil {
		log.Fatal(err)
	}

	// Log in with username "123456" and password "elradfmw".
	err = ftpConn.Login("user1", "123456") // Updated with the correct username and password.
	if err != nil {
		log.Fatal(err)
	}
	defer ftpConn.Logout()

	// Specify the local path and the FTP path.
	localPath := "C:/Users/hp/Documents/FTPFiles/WeatherAppdocs.txt" // Change this to your local file or directory.
	ftpPath := "/FTPServerPath"                                      // Change this to your desired remote FTP path.

	// Check if the local path is a file or a directory.
	fileInfo, err := os.Stat(localPath)
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		// If it's a directory, upload it.
		err = uploadDirectory(ftpConn, localPath, ftpPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Directory uploaded successfully.")
	} else {
		// If it's a file, upload it.
		err = uploadFile(ftpConn, localPath, ftpPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("File uploaded successfully.")
	}
}
