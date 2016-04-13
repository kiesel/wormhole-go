package wormhole

import (
	"gopkg.in/op/go-logging.v1"
	"io"
	"os/exec"
)

var log = logging.MustGetLogger("wormhole")

func transcribeOutput(prefix string, stream io.ReadCloser) {
	var buf = make([]byte, 1024)

	for {
		n, err := stream.Read(buf)

		if n > 0 {
			log.Info("%s: %s", prefix, buf)
		}

		if err != nil {
			return
		}
	}
}

func ExecuteCommand(executable string, args ...string) (err Error) {
	cmd := exec.Command(executable, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer stderr.Close()

	if err := cmd.Start(); err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("Started '%s' w/ PID %d", executable, cmd.Process.Pid)

	go transcribeOutput("out", stdout)
	go transcribeOutput("err", stderr)
	cmd.Wait()

	log.Info("PID %d has quit.", cmd.Process.Pid)

	return nil
}
