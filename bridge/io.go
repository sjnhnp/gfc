package bridge

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	Binary = "Binary"
	Text   = "Text"
)

func (a *App) WriteFile(path string, content string, options IOOptions) FlagResult {
	log.Printf("WriteFile [%s]: %s", options.Mode, path)

	fullPath := GetPath(path)

	if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
		return FlagResult{false, err.Error()}
	}

	var data []byte
	var err error

	switch options.Mode {
	case Text:
		data = []byte(content)
	case Binary:
		data, err = base64.StdEncoding.DecodeString(content)
		if err != nil {
			return FlagResult{false, err.Error()}
		}
	default:
		return FlagResult{false, "Unsupported IO mode: " + options.Mode}
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return FlagResult{false, err.Error()}
	}

	return FlagResult{true, "Success"}
}

func (a *App) ReadFile(path string, options IOOptions) FlagResult {
	log.Printf("ReadFile [%s]: %s", options.Mode, path)

	fullPath := GetPath(path)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return FlagResult{false, err.Error()}
	}

	switch options.Mode {
	case Text:
		return FlagResult{true, string(data)}
	case Binary:
		return FlagResult{true, base64.StdEncoding.EncodeToString(data)}
	default:
		return FlagResult{false, "Unsupported IO mode: " + options.Mode}
	}
}

func (a *App) MoveFile(source string, target string) FlagResult {
	log.Printf("MoveFile: %s -> %s", source, target)

	fullSource := GetPath(source)
	fullTarget := GetPath(target)

	if err := os.MkdirAll(filepath.Dir(fullTarget), os.ModePerm); err != nil {
		return FlagResult{false, err.Error()}
	}

	if err := os.Rename(fullSource, fullTarget); err != nil {
		return FlagResult{false, err.Error()}
	}

	return FlagResult{true, "Success"}
}

func (a *App) RemoveFile(path string) FlagResult {
	log.Printf("RemoveFile: %s", path)

	fullPath := GetPath(path)

	if err := os.RemoveAll(fullPath); err != nil {
		return FlagResult{false, err.Error()}
	}

	return FlagResult{true, "Success"}
}

func (a *App) CopyFile(src string, dst string) FlagResult {
	log.Printf("CopyFile: %s -> %s", src, dst)

	srcPath := GetPath(src)
	dstPath := GetPath(dst)

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return FlagResult{false, err.Error()}
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm); err != nil {
		return FlagResult{false, err.Error()}
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return FlagResult{false, err.Error()}
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return FlagResult{false, err.Error()}
	}

	return FlagResult{true, "Success"}
}

func (a *App) MakeDir(path string) FlagResult {
	log.Printf("MakeDir: %s", path)

	fullPath := GetPath(path)

	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		return FlagResult{false, err.Error()}
	}

	return FlagResult{true, "Success"}
}

func (a *App) ReadDir(path string) FlagResult {
	log.Printf("ReadDir: %s", path)

	fullPath := GetPath(path)

	files, err := os.ReadDir(fullPath)
	if err != nil {
		return FlagResult{false, err.Error()}
	}

	var result []string

	for _, file := range files {
		if info, err := file.Info(); err == nil {
			result = append(result, fmt.Sprintf("%v,%v,%v", info.Name(), info.Size(), info.IsDir()))
		}
	}

	return FlagResult{true, strings.Join(result, "|")}
}

func (a *App) OpenDir(path string) FlagResult {
	log.Printf("OpenDir: %s", path)

	fullPath := GetPath(path)

	err := browser.OpenURL(fullPath)
	if err != nil {
		return FlagResult{false, err.Error()}
	}

	return FlagResult{true, "Success"}
}

func (a *App) OpenURI(uri string) FlagResult {
	log.Printf("OpenURI: %s", uri)

	err := browser.OpenURL(uri)
	if err != nil {
		return FlagResult{false, err.Error()}
	}

	return FlagResult{true, "Success"}
}

func (a *App) AbsolutePath(path string) FlagResult {
	log.Printf("AbsolutePath: %s", path)

	absPath := GetPath(path)

	return FlagResult{true, absPath}
}

func (a *App) UnzipZIPFile(path string, output string) FlagResult {
	log.Printf("UnzipZIPFile: %s -> %s", path, output)

	fullPath := GetPath(path)
	outputPath := GetPath(output)

	archive, err := zip.OpenReader(fullPath)
	if err != nil {
		return FlagResult{false, err.Error()}
	}
	defer archive.Close()

	cleanOutputPath := outputPath + "/"

	for _, f := range archive.File {
		filePath := filepath.ToSlash(filepath.Clean(filepath.Join(outputPath, f.Name)))

		if !strings.HasPrefix(filePath, cleanOutputPath) {
			continue
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			continue
		}

		fileInArchive, err := f.Open()
		if err != nil {
			continue
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			fileInArchive.Close()
			continue
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			fileInArchive.Close()
			dstFile.Close()
			continue
		}

		fileInArchive.Close()
		dstFile.Close()
	}

	return FlagResult{true, "Success"}
}

func (a *App) UnzipTarGZFile(path string, output string) FlagResult {
	log.Printf("UnzipTarGZFile: %s -> %s", path, output)

	fullPath := GetPath(path)
	outputPath := GetPath(output)

	gzipFile, err := os.Open(fullPath)
	if err != nil {
		return FlagResult{false, err.Error()}
	}
	defer gzipFile.Close()

	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return FlagResult{false, err.Error()}
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	cleanOutputPath := outputPath + "/"

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return FlagResult{false, err.Error()}
		}

		filePath := filepath.ToSlash(filepath.Clean(filepath.Join(outputPath, header.Name)))

		if !strings.HasPrefix(filePath, cleanOutputPath) {
			continue
		}

		if header.Typeflag == tar.TypeDir {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			continue
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, header.FileInfo().Mode())
		if err != nil {
			continue
		}

		if _, err := io.Copy(dstFile, tarReader); err != nil {
			dstFile.Close()
			continue
		}

		dstFile.Close()
	}

	return FlagResult{true, "Success"}
}

func (a *App) UnzipGZFile(path string, output string) FlagResult {
	log.Printf("UnzipGZFile: %s -> %s", path, output)

	fullPath := GetPath(path)
	outputPath := GetPath(output)

	gzipFile, err := os.Open(fullPath)
	if err != nil {
		return FlagResult{false, err.Error()}
	}
	defer gzipFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return FlagResult{false, err.Error()}
	}
	defer outputFile.Close()

	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return FlagResult{false, err.Error()}
	}
	defer gzipReader.Close()

	if _, err := io.Copy(outputFile, gzipReader); err != nil {
		return FlagResult{false, err.Error()}
	}

	return FlagResult{true, "Success"}
}

func (a *App) FileExists(path string) FlagResult {
	log.Printf("FileExists: %s", path)

	path = GetPath(path)

	_, err := os.Stat(path)
	if err == nil {
		return FlagResult{true, "true"}
	}

	if os.IsNotExist(err) {
		return FlagResult{true, "false"}
	}

	return FlagResult{false, err.Error()}
}

// OpenFileDialog opens a file selection dialog and returns the selected file path
func (a *App) OpenFileDialog(title string, filters string) FlagResult {
	log.Printf("OpenFileDialog: %s, filters: %s", title, filters)

	// Parse filters: "YAML Files:*.yaml,*.yml|JSON Files:*.json"
	var dialogFilters []runtime.FileFilter
	if filters != "" {
		filterPairs := strings.Split(filters, "|")
		for _, pair := range filterPairs {
			parts := strings.Split(pair, ":")
			if len(parts) == 2 {
				dialogFilters = append(dialogFilters, runtime.FileFilter{
					DisplayName: parts[0],
					Pattern:     parts[1],
				})
			}
		}
	}

	opts := runtime.OpenDialogOptions{
		Title:   title,
		Filters: dialogFilters,
	}

	selectedFile, err := runtime.OpenFileDialog(a.Ctx, opts)
	if err != nil {
		return FlagResult{false, err.Error()}
	}

	if selectedFile == "" {
		return FlagResult{false, "cancelled"}
	}

	return FlagResult{true, filepath.ToSlash(selectedFile)}
}

// SaveFileDialog opens a save file dialog and returns the selected file path
func (a *App) SaveFileDialog(title string, defaultFilename string, filters string) FlagResult {
	log.Printf("SaveFileDialog: %s, default: %s, filters: %s", title, defaultFilename, filters)

	// Parse filters: "JSON Files:*.json|All Files:*.*"
	var dialogFilters []runtime.FileFilter
	if filters != "" {
		filterPairs := strings.Split(filters, "|")
		for _, pair := range filterPairs {
			parts := strings.Split(pair, ":")
			if len(parts) == 2 {
				dialogFilters = append(dialogFilters, runtime.FileFilter{
					DisplayName: parts[0],
					Pattern:     parts[1],
				})
			}
		}
	}

	opts := runtime.SaveDialogOptions{
		Title:           title,
		DefaultFilename: defaultFilename,
		Filters:         dialogFilters,
	}

	selectedFile, err := runtime.SaveFileDialog(a.Ctx, opts)
	if err != nil {
		return FlagResult{false, err.Error()}
	}

	if selectedFile == "" {
		return FlagResult{false, "cancelled"}
	}

	return FlagResult{true, filepath.ToSlash(selectedFile)}
}

// ReadExternalFile reads a file from an absolute path (outside the app's data directory)
func (a *App) ReadExternalFile(path string) FlagResult {
	log.Printf("ReadExternalFile: %s", path)

	data, err := os.ReadFile(path)
	if err != nil {
		return FlagResult{false, err.Error()}
	}

	return FlagResult{true, string(data)}
}

// WriteExternalFile writes content to an absolute path (outside the app's data directory)
func (a *App) WriteExternalFile(path string, content string) FlagResult {
	log.Printf("WriteExternalFile: %s", path)

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return FlagResult{false, err.Error()}
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return FlagResult{false, err.Error()}
	}

	return FlagResult{true, "Success"}
}
