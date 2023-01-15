package main

import (
  "bufio"
  "fmt"
  "log"
  "os"
  "path/filepath"
  "regexp"
  "strings"
)

// clocktable
// A simple CLI for aggregating logbooks in org.
func main() {
  fmt.Println("clocktable 2023.01.0")

  arguments := os.Args[1:]

  if len(arguments) == 0 {
    fmt.Println("Missing required argument: directory")
    os.Exit(4)
  }

  directory := arguments[0]

  fmt.Printf("Scanning %v\n", directory)
  processDirectory(directory)
}

func processDirectory(path string) {
  rootDir, err := os.ReadDir(path)
  if err != nil {
    log.Fatal(err)
  }

  for _, fsItem := range rootDir {
    itemPath := filepath.Join(path, fsItem.Name())
    if fsItem.IsDir() {
      processDirectory(itemPath)
    } else {
      readFile, err := os.Open(itemPath)
      if err != nil {
        log.Fatalln(err)
      } else {
        fileScanner := bufio.NewScanner(readFile)
        fileScanner.Split(bufio.ScanLines)

        var currentTask string
        clockMap := make(map[string][]string)
        for fileScanner.Scan() {
          lineContents := strings.TrimSpace(fileScanner.Text())

          regexTask, _ := regexp.Compile("^\\* ([A-Z]).+")
          if regexTask.MatchString(lineContents) {
            // if is task
            if len(clockMap[currentTask]) == 0 {
              // if previous task has no logbook items
              delete(clockMap, currentTask)
            }

            currentTask = lineContents
            clockMap[currentTask] = []string{}
          }

          regexClock, _ := regexp.Compile("^CLOCK:.+=> ")
          if regexClock.MatchString(lineContents) {
            // if is clocked out logbook item
            clockMap[currentTask] = append(clockMap[currentTask], regexClock.ReplaceAllString(lineContents, ""))
          }
        }

        if len(clockMap) > 1 {
          fmt.Printf("#+TITLE: %v\n", fsItem.Name())
        }

        for task, logbook := range clockMap {
          if len(logbook) > 0 {
            fmt.Println(task)
            // TODO: parse logbook
            fmt.Println(logbook)
            fmt.Println()
          }
        }
      }
    }
  }
}
