package main

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/loophole/cli/cmd"
	"github.com/tcnksm/go-latest"
)

// Will be filled in during build
var version = "development"
var commit = "unknown"

func main() {
	githubTag := &latest.GithubTag{
		Owner:      "loophole",
		Repository: "cli",
	}

	if version == "development" {
		fmt.Println("Update check disabled while using development version.")
	} else {
		res, err := latest.Check(githubTag, version)
		if err != nil {
			log.Fatal("GithubTag error:" + err.Error())
		}
		if _, err := os.Stat("loophole_version_" + res.Current); err == nil {
			fmt.Println("################")
			fmt.Println("It looks like you recently downloaded a newer version of Loophole, please use it instead of this one!")
			fmt.Println("It should be located in the folder \"loophole_version_" + res.Current + "\"")
			fmt.Println("################")
			fmt.Println()
			res.Outdated = false
		}

		if res.Outdated {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Your version of Loophole is outdated. Do you wish to download version " + res.Current + " now?")

			fmt.Print("Y/n : ")
			text, _ := reader.ReadString('\n')
			// convert CRLF to LF
			text = strings.Replace(text, "\n", "", -1)

			fmt.Println(text)

			if strings.Contains(text, "n") || strings.Contains(text, "N") {
				//skip update
			} else {
				archiveExt := ".tar.gz"
				if runtime.GOOS == "windows" {
					archiveExt = ".zip"
				}
				urlBase := "https://github.com/loophole/cli/releases/download/"
				url := fmt.Sprintf("%s%s%s%s%s%s%s%s%s", urlBase, res.Current, "/loophole_", res.Current, "_", runtime.GOOS, "_", runtime.GOARCH, archiveExt)
				fileName := "loophole_version_" + res.Current
				archiveName := res.Current
				err := download(archiveName+archiveExt, url)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				if runtime.GOOS == "windows" {
					err = extractZip(archiveName+archiveExt, fileName)
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(1)
					}
				} else {
					err = extractTarGz(archiveName+archiveExt, fileName)
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(1)
					}
				}
				err = os.Remove(archiveName + archiveExt)
				if err != nil {
					fmt.Println("Unable to delete downloaded compressed file: " + err.Error())
				}
				fmt.Println("Download finished! Please start the new version located in the folder: loophole_version_" + res.Current)
				os.Exit(0)
			}
		}
	}
	cmd.Execute(version, commit)
}

func download(filepath string, url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func extractTarGz(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}

	gzipReader, err := gzip.NewReader(sourceFile)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzipReader)

	if err := os.Mkdir(dest, 0755); err != nil {
		return err
	}

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if strings.Contains(header.Name, "..") {
			fmt.Println("Illegal path: File name contains '..'")
			os.Exit(1)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(dest + "/" + header.Name)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
			outFile.Close()

		default:
			fmt.Printf("extractTarGz: unknown type: %b in %s", header.Typeflag, header.Name)
			os.Exit(1)
		}

	}
	return nil
}

func extractZip(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	os.MkdirAll(dest, 0755)

	for _, file := range reader.File {
		data, err := file.Open()

		if err != nil {
			return err
		}
		defer data.Close()

		if strings.Contains(file.Name, "..") {
			fmt.Println("Illegal path: File name contains '..'")
			os.Exit(1)
		}

		path := filepath.Join(dest, file.Name)

		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), file.Mode())
			endFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}

			_, err = io.Copy(endFile, data)
			endFile.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
