package files

import (
	"archive/tar"
    "bufio"
	"compress/gzip"
	"fmt"
	"io"
    "path/filepath"
	"os"
    "runtime"
	"strings"
    "github.com/rickar/props"
)

const (
    PathSeparator = string(os.PathSeparator)
)

func UserHomeDir() string {
    if runtime.GOOS == "windows" {
        home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
        if home == "" {
            home = os.Getenv("USERPROFILE")
        }
        return home
    }
    return os.Getenv("HOME")
}

func IsExistingFileOrDirectory(path string) bool {
    if fi, err := os.Stat(path); err == nil {
        switch mode := fi.Mode(); {
        case mode.IsRegular():
            return true
        case mode.IsDir():
            return true
        }
    }
    return false
}

func IsExistingFile(path string) bool {
    if fi, err := os.Stat(path); err == nil {
        switch mode := fi.Mode(); {
        case mode.IsRegular():
            return true
        case mode.IsDir():
            return false
        }
    }
    return false
}

func IsExistingDirectory(path string) bool {
    if fi, err := os.Stat(path); err == nil {
        switch mode := fi.Mode(); {
        case mode.IsRegular():
            return false
        case mode.IsDir():
            return true
        }
    }
    return false
}

func MkDir(path string) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        os.Mkdir(path, 0700)
    }
}

func LoadPropertiesFile(path string) *props.Properties {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        // path does not exist
        return nil
    }
    // path exists, so load the properties file ...
    file, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    propFile := bufio.NewReader(file)
    props, err := props.Read(propFile)
    if err != nil {
        panic(err)
    }
    return props
}

func CreateTar(destinationfile string, parentDir string, nameOfFileOrDirectoryToArchive string, verbose bool) error {
    prefix := parentDir
    if !strings.HasSuffix(prefix, PathSeparator) {
        prefix = prefix + PathSeparator
    }

    // Create the tar file and the writer 
    tarfile, err := os.Create(destinationfile)
    if err != nil {
        return err
    }

    defer tarfile.Close()
    var fileWriter io.WriteCloser = tarfile
    if strings.HasSuffix(destinationfile, ".gz") {
        // Add a gzip filter if .gz is in the destination filename
        fileWriter = gzip.NewWriter(tarfile) // add a gzip filter
        defer fileWriter.Close()
    }

    tarfileWriter := tar.NewWriter(fileWriter)
    defer tarfileWriter.Close()

    // Define a function to write each file we find
    walkFunc := func(path string, f os.FileInfo, err error) error {
        // skip all directories
        if f.IsDir() {
            return nil
        }

        // Compute the relative path
        relPath := strings.Replace(path, prefix, "", 1)
        if verbose {
            fmt.Println("archiving  " + relPath)
        }

        // prepare the tar header
        header := new(tar.Header)
        header.Name = relPath
        header.Size = f.Size()
        header.Mode = int64(f.Mode())
        header.ModTime = f.ModTime()

        // write the entry for the file
        err = tarfileWriter.WriteHeader(header)
        if err != nil {
            return err
        }

        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()

        // and the file contents
        _, err = io.Copy(tarfileWriter, file)
        if err != nil {
            return err
        }
        return nil
    }

    // Read the directory structure recursively using our function, skipping symbolic links
    targetPath := parentDir + PathSeparator + nameOfFileOrDirectoryToArchive
    return filepath.Walk(targetPath, walkFunc)
}

func ExtractTar(tarFilePath string, parentDirPath string, verbose bool) error {
    prefix := ""
    if len(parentDirPath) != 0 {
        prefix = parentDirPath + PathSeparator
    }

	// Open the tar file
    file, err := os.Open(tarFilePath)
    if err != nil {
    	return err
    }
    defer file.Close()

    // Read the tar file
    var fileReader io.ReadCloser = file
	if strings.HasSuffix(tarFilePath, ".gz") {
    	// we are reading a tar.gz file, so add a filter to handle gzipped file
        if fileReader, err = gzip.NewReader(file); err != nil {
        	return err
        }
        defer fileReader.Close()
    }

    // Extracting tarred files
    tarBallReader := tar.NewReader(fileReader)
    for {
        header, err := tarBallReader.Next()
        if err != nil {
            if err == io.EOF {
                break
            }
            return err
        }

		// get the individual filename and extract to the current directory
        filename := filepath.FromSlash(prefix + header.Name)
        if verbose {
            fmt.Println(filename)
        }
		switch header.Typeflag {
        case tar.TypeDir:
            // handle directory
            if verbose {
                fmt.Println("Making dir:  " + filename)
            }
            err = os.MkdirAll(filename, os.FileMode(header.Mode)) // or use 0755 if you prefer
            if err != nil {
            	return err
            }
        case tar.TypeReg:
            fallthrough
        case tar.TypeRegA:
            // handle normal file
            dirPath := filepath.Dir(filename)
            if _, err := os.Stat(dirPath); os.IsNotExist(err) {
                if verbose {
                    fmt.Println("Making dir:  " + dirPath)
                }
                err = os.MkdirAll(dirPath, 0755) // or use 0755 if you prefer
                if err != nil {
                    return err
                }
            }
            if verbose {
                fmt.Println("Making file: " + filename)
            }
            writer, err := os.Create(filename)
            if err != nil {
            	return err
            }
            io.Copy(writer, tarBallReader)
            err = os.Chmod(filename, os.FileMode(header.Mode))
            writer.Close()
            if err != nil {
            	return err
            }
        default:
            if verbose {
                fmt.Printf("Unable to untar type : %c in file %s", header.Typeflag, filename)
            }
        }
    }
    return nil
}
