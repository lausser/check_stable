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
  //"syscall"
)

type Result struct {
  // Dreck! Die Felder muessen mit Grossbuchstaben anfangen
  // da sie sonst der feine Herr Golang nicht exportiert
  Output string `json:"output"`
  Exitcode int `json:"exitcode"`
  Serial int `json:"serial"`
}

func main() {
  args := os.Args
  cmd := args[1]
  params := args[2:len(args)]
  var statefilesdir string;
  var result Result

  sha1_hash := fmt.Sprintf("%x", sha1.Sum([]byte(strings.Join(args, ""))))
  if (os.Getenv("OMD_ROOT") != "") {
    statefilesdir = filepath.Join(os.Getenv("OMD_ROOT"),
        "var", "tmp", "stabilize");

  } else {
    statefilesdir = filepath.Join("/", "var", "tmp", "stabilize");
  }
  os.MkdirAll(statefilesdir, 0755)
  statefile := filepath.Join(statefilesdir, sha1_hash)

  output, err := exec.Command(cmd, params...).CombinedOutput()

  //if err := cmd.Start(); err != nil {
    //log.Fatalf("cmd.Start: %v")
  //}

  //if err := cmd.Wait(); err != nil {
    //if exiterr, ok := err.(*exec.ExitError); ok {
      // The program has exited with an exit code != 0

      // There is no plattform independent way to retrieve
      // the exit code, but the following will work on Unix
      //if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
        //log.Printf("Exit Status: %d", status.ExitStatus())
      //}
    //} else {
      //log.Fatalf("cmd.Wait: %v", err)
    //}
  //}

  fmt.Printf("The date is %s\n", output)
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
    if strings.Index(string(output), errors[i]) != -1 {
      error_found = true
      fmt.Printf("plugin error ->%s\n", errors[i])
      break
    }
  }

  raw, err := ioutil.ReadFile(statefile)
  if err != nil {
    result = Result{string(output), 0, 0}
    // error
  } else {
    json.Unmarshal(raw, &result)
  }

  // load json, output, serial
  // else record.output = output, recore.serial = max+1
  if error_found {
    fmt.Printf("no error write to %s\n", statefile)
    // if serial < max, print output, save output with serial++
    // else print error output
    if result.Serial > 2 {
      fmt.Printf("%s\n", output)
    }
  } else {
    // save output, 
    fmt.Printf("no plugin error write good %s\n", statefile)
    // save with serial 0
    // record.serial = 0
    json_result, err := json.Marshal(result)
    fmt.Printf("file written %s\n", string(json_result))
    if err != nil {
      fmt.Println("error:", err)
    }
    err = ioutil.WriteFile(statefile, json_result, 0644)
    if err != nil {
      fmt.Printf("file write fail %s\n", err)
    } else {
      fmt.Printf("file written %s\n", statefile)
    }
  }
}

