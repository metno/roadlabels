package exttools

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var AnalysisDir = "/lustre/storeB/immutable/archive/projects/metproduction/yr_short"

func GetTemp(reftime time.Time, lat float32, lon float32) (float64, error) {

	datestr := reftime.Format("2006010215")
	datelong, err := strconv.ParseInt(datestr, 10, 64)
	if err != nil {
		return -273.15, fmt.Errorf("ParseInt: %v", err)
	}
	pythonCode := fmt.Sprintf(`import ncvars; ncvars.analysis_dir="%s"; import sys; ncvars.print_t2m(%d, %f, %f);`, AnalysisDir,
		datelong, lat, lon)

	log.Printf("pythoncode: %s", pythonCode)
	cout, cerr, err := RunCommand("python3", "-c", pythonCode)
	if err != nil {
		log.Printf("%s %s\n", err, cerr)
		return -273.15, fmt.Errorf("%v. %s", err, cerr)
	}

	v, err := strconv.ParseFloat(strings.TrimSuffix(cout, "\n"), 64)
	return v, err
}

func RunCommand(name string, args ...string) (string, string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	var outbuf, errbuf bytes.Buffer
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("RunCommand timed out")
		return "", "", err

	}

	stdout := outbuf.String()
	stderr := errbuf.String()

	exitCode := 63
	exitError, castOk := err.(*exec.ExitError)

	if castOk {
		ws := exitError.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	} else {
		if err == nil {
			exitCode = 0
		}
	}
	d, castOk := err.(*exec.Error)
	if castOk {
		return stdout, stderr, fmt.Errorf("cmd.Run(): %v", d)
	}

	//log.Printf("command result stdout: %v, stderr: %v, exitCode: %v", stdout, stderr, exitCode)
	if exitCode != 0 {
		return stdout, stderr, fmt.Errorf("exit code: %d", exitCode)
	}
	if err != nil {
		return stdout, stderr, err
	}
	return stdout, stderr, err
}
