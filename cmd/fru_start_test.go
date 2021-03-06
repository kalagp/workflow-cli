//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package cmd_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/dellemc-symphony/workflow-cli/resources"
	homedir "github.com/mitchellh/go-homedir"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FruStart", func() {
	var binLocation string
	var endpointLocation string
	var StateFile string
	var target string
	var tableString string
	var nodeSelection string

	BeforeEach(func() {
		binLocation = fmt.Sprintf("../bin/%s/workflow-cli", runtime.GOOS)
		endpointLocation = fmt.Sprintf("../bin/%s/endpoint.yaml", runtime.GOOS)

		tableString = `+--------+----------+---------------+------------+----------+
| SELECT | HOSTNAME | SERIAL NUMBER |  MGMT IP   |  STATUS  |
+--------+----------+---------------+------------+----------+
|      1 | node01   |       1234567 | 10.10.10.1 | online   |
|      2 | node02   |      98765432 | 10.10.10.2 | degraded |
|      3 | node03   |      91827465 | 10.10.10.3 | online   |
+--------+----------+---------------+------------+----------+
`

		nodeSelection = `+----------+---------------+------------+----------+
| HOSTNAME | SERIAL NUMBER |  MGMT IP   |  STATUS  |
+----------+---------------+------------+----------+
| node02   |      98765432 | 10.10.10.2 | degraded |
+----------+---------------+------------+----------+
`

		dir, err := homedir.Dir()
		Expect(err).ToNot(HaveOccurred())
		StateFile = fmt.Sprintf("%s/.cli", dir)

		if https {
			target = "https://localhost:8080"
		} else {
			target = "http://localhost:8080"
		}

		cmd := exec.Command(binLocation, "target", target)
		err = cmd.Run()
		Expect(err).To(BeNil())

		// Remove any endpoint file previously in bin
		os.Remove(endpointLocation)
	})
	AfterEach(func() {
		os.Remove(StateFile)
	})

	Context("When the endpoint file has all fields filled", func() {
		It("UNIT should run the 'fru start' without user input", func() {
			err := resources.WriteEndpointsFile("AllFields", endpointLocation)
			Expect(err).To(BeNil())

			startTime := time.Now()

			cmd := exec.Command(binLocation, "fru", "start")

			stdin, err := cmd.StdinPipe()
			Expect(err).To(BeNil())
			defer stdin.Close()

			stdout, err := cmd.StdoutPipe()
			Expect(err).To(BeNil())
			defer stdout.Close()

			stderr, err := cmd.StderrPipe()
			Expect(err).To(BeNil())
			defer stderr.Close()

			cmd.Start()

			// Select node 2 for removal
			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm selection
			io.WriteString(stdin, "Y\n")
			time.Sleep(500 * time.Millisecond)

			// Select node 2 for add
			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm selection
			io.WriteString(stdin, "Y\n")
			time.Sleep(500 * time.Millisecond)

			errBuf := new(bytes.Buffer)
			errBuf.ReadFrom(stderr)
			outBuf := new(bytes.Buffer)
			outBuf.ReadFrom(stdout)

			Expect(err).To(BeNil())
			Expect(outBuf.String()).To(ContainSubstring(tableString))
			Expect(outBuf.String()).To(ContainSubstring(nodeSelection))
			Expect(errBuf.String()).To(ContainSubstring("Workflow complete"))

			elapsedTime := time.Since(startTime)
			Expect(elapsedTime.Seconds()).To(BeNumerically(">", 5))
		})
	})

	Context("When the endpoint file has endpoints but not credentials", func() {
		It("UNIT should run the 'fru start' and prompt for input", func() {
			err := resources.WriteEndpointsFile("MissingCredentials", endpointLocation)
			Expect(err).To(BeNil())

			cmd := exec.Command(binLocation, "fru", "start")

			stdin, err := cmd.StdinPipe()
			Expect(err).To(BeNil())
			defer stdin.Close()

			stdout, err := cmd.StdoutPipe()
			Expect(err).To(BeNil())
			defer stdout.Close()

			stderr, err := cmd.StderrPipe()
			Expect(err).To(BeNil())
			defer stderr.Close()

			cmd.Start()

			//Username for RackHD
			io.WriteString(stdin, "RackHDUsername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for RackHD
			io.WriteString(stdin, "RackHDPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Username for CoprHD
			io.WriteString(stdin, "CoprHDusername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for CoprHD
			io.WriteString(stdin, "CoprHDpassword\n")
			time.Sleep(500 * time.Millisecond)

			//Username for vCenter
			io.WriteString(stdin, "vCenterUsername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for vCenter
			io.WriteString(stdin, "vCenterPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Username for ScaleIO
			io.WriteString(stdin, "ScaleIOUsername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for ScaleIO
			io.WriteString(stdin, "ScaleIOPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Select node 2 for removal
			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			// Confirm selection
			io.WriteString(stdin, "Y\n")
			time.Sleep(500 * time.Millisecond)

			//Username for ScaleIO MDM
			io.WriteString(stdin, "ScaleIOMDMUsername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for ScaleIO MDM
			io.WriteString(stdin, "ScaleIOMDMPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Select node 2 for addition
			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm selection
			io.WriteString(stdin, "Y\n")
			time.Sleep(500 * time.Millisecond)

			errBuf := new(bytes.Buffer)
			errBuf.ReadFrom(stderr)
			outBuf := new(bytes.Buffer)
			outBuf.ReadFrom(stdout)

			Expect(outBuf.String()).To(ContainSubstring(tableString))
			Expect(outBuf.String()).To(ContainSubstring(nodeSelection))
			Expect(errBuf.String()).To(ContainSubstring("Workflow complete"))

			cmd.Wait()

		})
	})

	Context("When the endpoint file is missing", func() {
		It("UNIT should prompt for endpoints and creds", func() {
			cmd := exec.Command(binLocation, "fru", "start")

			stdin, err := cmd.StdinPipe()
			Expect(err).To(BeNil())
			defer stdin.Close()

			stdout, err := cmd.StdoutPipe()
			Expect(err).To(BeNil())
			defer stdout.Close()

			stderr, err := cmd.StderrPipe()
			Expect(err).To(BeNil())
			defer stderr.Close()

			cmd.Start()

			//Endpoint for RackHD
			io.WriteString(stdin, "RackHDEndpoint\n")
			time.Sleep(500 * time.Millisecond)

			//Username for RackHD
			io.WriteString(stdin, "RackHDUsername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for RackHD
			io.WriteString(stdin, "RackHDPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for CoprHD
			io.WriteString(stdin, "CoprHDEndpoint\n")
			time.Sleep(500 * time.Millisecond)

			//Username for CoprHD
			io.WriteString(stdin, "CoprHDusername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for CoprHD
			io.WriteString(stdin, "CoprHDpassword\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for vCenter
			io.WriteString(stdin, "vCenterEndpoint\n")
			time.Sleep(500 * time.Millisecond)

			//Username for vCenter
			io.WriteString(stdin, "vCenterUsername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for vCenter
			io.WriteString(stdin, "vCenterPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for ScaleIO
			io.WriteString(stdin, "ScaleIOEndpoint\n")
			time.Sleep(500 * time.Millisecond)

			//Username for ScaleIO
			io.WriteString(stdin, "ScaleIOUsername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for ScaleIO
			io.WriteString(stdin, "ScaleIOPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Select node 2 for removal
			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			// Confirm selection
			io.WriteString(stdin, "Y\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for ScaleIO MDM
			io.WriteString(stdin, "ScaleIOMDMEndpoint\n")
			time.Sleep(500 * time.Millisecond)

			//Username for ScaleIO MDM
			io.WriteString(stdin, "ScaleIOMDMUsername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for ScaleIO MDM
			io.WriteString(stdin, "ScaleIOMDMPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Select node 2 for addition
			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm selection
			io.WriteString(stdin, "Y\n")
			time.Sleep(500 * time.Millisecond)

			errBuf := new(bytes.Buffer)
			errBuf.ReadFrom(stderr)
			outBuf := new(bytes.Buffer)
			outBuf.ReadFrom(stdout)

			Expect(errBuf.String()).To(ContainSubstring("endpoint.yaml not found, will prompt user"))
			Expect(outBuf.String()).To(ContainSubstring(tableString))
			Expect(outBuf.String()).To(ContainSubstring(nodeSelection))
			Expect(errBuf.String()).To(ContainSubstring("Workflow complete"))

		})
	})

})
