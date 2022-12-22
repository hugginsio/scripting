package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "strings"
)

var parentDirectory string
var newDirectory string
var removeJson bool

func main() {
  fmt.Println("photo-helper 2022.11.0")

  flag.StringVar(&parentDirectory, "directory", "", "full path to the directory containing image directories")
  flag.StringVar(&newDirectory, "newDirectory", "", "name of the new directory to place images in")
  flag.BoolVar(&removeJson, "rmJson", true, "remove all *.JSON files")
  flag.Parse()

  if parentDirectory == "" {
    flag.PrintDefaults()
    os.Exit(4)
  }

  if newDirectory == "" {
    newDirectory = "photo-helper-collection"
  }

  err := os.Mkdir(newDirectory, 0755)

  if err != nil {
    log.Fatalf("Failed to create directory \"%v\"\n", newDirectory)
  }

  pathSeparator := string(os.PathSeparator)
  if string(parentDirectory[len(parentDirectory)-1:]) != pathSeparator {
    parentDirectory = parentDirectory + string(os.PathSeparator)
  }

  results, err := os.ReadDir(parentDirectory)

  if err != nil {
    log.Fatalln(err.Error())
  }

  var directories []string
  for _, dir := range results {
    if dir.IsDir() && dir.Name() != newDirectory {
      directories = append(directories, dir.Name())
    }
  }

  for _, dir := range directories {
    fmt.Printf("Processing directory \"%v\"\n", dir)
    images, err := os.ReadDir(fmt.Sprintf("%v%v", parentDirectory, dir))

    if err != nil {
      fmt.Printf("Error processing directory %v: \"%v\"", dir, err.Error())
    } else {
      var directoryFiles []string
      for _, imageFile := range images {
        directoryFiles = append(directoryFiles, imageFile.Name())
      }

      removedJsonFiles := 0
      movedFiles := 0

      for _, file := range directoryFiles {
        if strings.HasSuffix(strings.ToLower(file), ".json") {
          err = os.Remove(parentDirectory + dir + pathSeparator + file)

          if err != nil {
            fmt.Printf("- Failed to remove file \"%v\": %v\n", file, err.Error())
          } else {
            removedJsonFiles += 1
          }
        } else {
          oldFilePath := parentDirectory + dir + pathSeparator + file
          newFileName := strings.ReplaceAll(dir, " ", "-") + "-" + file
          newFilePath := parentDirectory + newDirectory + pathSeparator + newFileName
          err = os.Rename(oldFilePath, newFilePath)
          if err != nil {
            fmt.Printf("- Failed to rename & move file \"%v\": %v\n", file, err.Error())
          } else {
            movedFiles += 1
          }
        }
      }

      fmt.Printf("- Removed %v JSON file(s), relocated %v other file(s)\n", removedJsonFiles, movedFiles)
    }
  }
}
