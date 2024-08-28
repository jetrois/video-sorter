package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter the path to your movie files: ")
    inputDir, _ := reader.ReadString('\n')
    inputDir = strings.TrimSpace(inputDir) // Remove any surrounding whitespace or newlines
    inputDir = strings.Trim(inputDir, "'") // Remove surrounding single quotes if any

    sdDir := filepath.Join(inputDir, "SD")
    hd720pDir := filepath.Join(inputDir, "720p")
    fullHD1080pDir := filepath.Join(inputDir, "1080p")
    ultraHD2160pDir := filepath.Join(inputDir, "2160p")

    // Create directories if they don't exist
    os.MkdirAll(sdDir, os.ModePerm)
    os.MkdirAll(hd720pDir, os.ModePerm)
    os.MkdirAll(fullHD1080pDir, os.ModePerm)
    os.MkdirAll(ultraHD2160pDir, os.ModePerm)

    err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !info.IsDir() {
            width, height, err := getVideoResolution(path)
            if err != nil {
                log.Printf("Failed to get resolution for %s: %v", path, err)
                return nil
            }

            switch {
            case height < 718:
                moveFile(path, sdDir)
            case height == 720:
                moveFile(path, hd720pDir)
            case height == 1080:
                moveFile(path, fullHD1080pDir)
            case height == 2160:
                moveFile(path, ultraHD2160pDir)
            default:
                log.Printf("Unhandled resolution for %s: %dx%d", path, width, height)
            }
        }
        return nil
    })

    if err != nil {
        log.Fatalf("Error walking the path %q: %v\n", inputDir, err)
    }

    fmt.Println("Sorting completed successfully!")
}

func getVideoResolution(filePath string) (int, int, error) {
    cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", filePath)
    output, err := cmd.Output()
    if err != nil {
        return 0, 0, err
    }

    // Regex to extract resolution
    re := regexp.MustCompile(`(\d+)x(\d+)`)
    matches := re.FindStringSubmatch(string(output))

    if len(matches) < 3 {
        return 0, 0, fmt.Errorf("could not determine resolution for %s", filePath)
    }

    width, _ := strconv.Atoi(matches[1])
    height, _ := strconv.Atoi(matches[2])

    return width, height, nil
}

func moveFile(sourcePath, destDir string) {
    destPath := filepath.Join(destDir, filepath.Base(sourcePath))
    err := os.Rename(sourcePath, destPath)
    if err != nil {
        log.Printf("Failed to move file %s to %s: %v", sourcePath, destDir, err)
    } else {
        fmt.Printf("Moved %s to %s\n", sourcePath, destDir)
    }
}