package main

import (
  "fmt"
  "os"
  "os/exec"
  "strings"
  "path/filepath"
  "crypto/sha1"
  "encoding/json"
  "io/ioutil"
  "syscall"
)

type Result struct {
  // Dreck! Die Felder muessen mit Grossbuchstaben anfangen
  // da sie sonst der feine Herr Golang nicht exportiert
  Output string `json:"output"`
  ExitCode int `json:"exitcode"`
  Serial int `json:"serial"`
}

func match() {
}

func loadResult(resultFile string) Result {
  var result Result;
  raw, err := ioutil.ReadFile(resultFile)
  if err != nil {
    result = Result{"", 0, 0}
  } else {
    json.Unmarshal(raw, &result)
  }
  return result
}

func saveResult(resultFile string, result Result) {
  json_result, err := json.Marshal(result)
  if err != nil {
    fmt.Println("json error:", err)
  }
  err = ioutil.WriteFile(resultFile, json_result, 0644)
  if err != nil {
    fmt.Printf("could not write resultFile %s\n", err)
  }
}

func initResultFile(args string) string {
  var resultDir string;
  sha1_hash := fmt.Sprintf("%x", sha1.Sum([]byte(args)))
  if (os.Getenv("OMD_ROOT") != "") {
    resultDir = filepath.Join(os.Getenv("OMD_ROOT"),
        "var", "tmp", "check_stable");

  } else {
    resultDir = filepath.Join("/", "var", "tmp", "check_stable");
  }
  os.MkdirAll(resultDir, 0755)
  return filepath.Join(resultDir, sha1_hash)
}

func main() {
  args := os.Args
  cmd := args[1]
  params := args[2:len(args)]
  exitCode := 0

  resultFile := initResultFile(strings.Join(args, ""))
  plugin := exec.Command(cmd, params...)

  pluginOutput, err := plugin.CombinedOutput()
  if err != nil {
    if exitError, ok := err.(*exec.ExitError); ok {
      exitCode = exitError.Sys().(syscall.WaitStatus).ExitStatus()
    }
  } else {
    exitCode = plugin.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
  }

  var errors = []string{
      "Return code of 127 is out of bounds - plugin may be missing",
      "Return code of 254 is out of bounds",
      "Return code of 255 is out of bounds",
      "Service Check Timed Out",
      "could not start controlmaster",
      "Could not open pipe: /usr/bin/ssh",
  }
  var error_found = false
  for i:= 0; i < len(errors); i++ {
    if strings.Index(string(pluginOutput), errors[i]) != -1 {
      error_found = true
      break
    }
  }

  if error_found {
    // plugin ended with one of these high-load side effects
    lastResult := loadResult(resultFile)
    if lastResult.Serial > 1 {
      // the error did not go away, time to alert
      fmt.Printf("%s\n", pluginOutput)
      os.Exit(exitCode)
    } else {
      // output a fake result and increment the counter
      saveResult(resultFile, Result{ lastResult.Output, lastResult.ExitCode, lastResult.Serial + 1 })
      fmt.Printf("(stabilized #%d) %s\n", lastResult.Serial, lastResult.Output)
      os.Exit(lastResult.ExitCode)
    }
  } else {
    // save the output, plugin had no problems
    saveResult(resultFile, Result{ string(pluginOutput), exitCode, 0 })
    // exit with the original plugin result
    fmt.Printf("%s\n", pluginOutput)
    os.Exit(exitCode)
  }
}

