//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package cmd_test

import (
	"flag"
	"time"

	"github.com/dellemc-symphony/workflow-cli/mock"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"gopkg.in/gin-gonic/gin.v1"

	"testing"
)

var https bool

func init() {
	flag.BoolVar(&https, "https", false, "Set 'true' to enable HTTPS for mock REST endpoint")
}

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Cmd Suite", []Reporter{junitReporter})
}

var router *gin.Engine

var _ = BeforeSuite(func() {
	mock.CreateMock(https)
	time.Sleep(50 * time.Millisecond)
})

var _ = AfterSuite(func() {
	mock.StopMock()
})
