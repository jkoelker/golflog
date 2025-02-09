//
// Copyright 2022 FloatMe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package golflog_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	// nolint:gci // This is the correct order
	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/floatme-corp/golflog"
	"github.com/floatme-corp/golflog/mocks"
)

func bufferLogger() (*bytes.Buffer, logr.Logger) {
	buf := new(bytes.Buffer)
	log := funcr.New(
		func(prefix, args string) {
			fmt.Fprint(buf, prefix, args)
		},
		funcr.Options{
			Verbosity: golflog.MaxLevel,
		},
	)

	return buf, log
}

type NewLoggerSuite struct {
	suite.Suite

	configurator *mocks.Configurator
}

func (suite *NewLoggerSuite) SetupTest() {
	suite.configurator = &mocks.Configurator{}
}

func (suite *NewLoggerSuite) TearDownTest() {
	suite.configurator.AssertExpectations(suite.T())
}

func (suite *NewLoggerSuite) TestVerbosityMaxClamping() {
	buf, logger := bufferLogger()

	suite.configurator.On(
		"Verbosity",
		mock.AnythingOfType("int"),
	).Return(func(verbosity int) error {
		suite.Equal(golflog.MaxLevel, verbosity)

		return nil
	}).Once()

	suite.configurator.On("Build").Return(logger, nil)

	log, err := golflog.NewLogger(suite.configurator, "test", golflog.MaxLevel+1)
	if suite.NoError(err) {
		suite.NotNil(log)
		log.Info("test")
		suite.Equal(`test"level"=0 "msg"="test"`, buf.String())
	}
}

func (suite *NewLoggerSuite) TestVerbosityMinClamping() {
	buf, logger := bufferLogger()

	suite.configurator.On(
		"Verbosity",
		mock.AnythingOfType("int"),
	).Return(func(verbosity int) error {
		suite.Equal(golflog.MinLevel, verbosity)

		return nil
	}).Once()

	suite.configurator.On("Build").Return(logger, nil)

	log, err := golflog.NewLogger(suite.configurator, "test", golflog.MinLevel-1)
	if suite.NoError(err) {
		suite.NotNil(log)
		log.Info("test")
		suite.Equal(`test"level"=0 "msg"="test"`, buf.String())
	}
}

func (suite *NewLoggerSuite) TestVerbosityError() {
	mockErr := errors.New("mock error")

	suite.configurator.On("Verbosity", mock.AnythingOfType("int")).Return(mockErr).Once()

	_, err := golflog.NewLogger(suite.configurator, "test", 0)
	suite.Error(err)
}

func (suite *NewLoggerSuite) TestBuildError() {
	mockErr := errors.New("mock error")

	suite.configurator.On("Verbosity", mock.AnythingOfType("int")).Return(nil).Once()
	suite.configurator.On("Build").Return(logr.Logger{}, mockErr)

	_, err := golflog.NewLogger(suite.configurator, "test", 0)
	if suite.Error(err) {
		suite.ErrorIs(err, mockErr)
	}
}

func TestNewLogger(t *testing.T) {
	suite.Run(t, new(NewLoggerSuite))
}

type NewLoggerWithBuildInfoSuite struct {
	suite.Suite

	configurator *mocks.Configurator
	buildInfo    *mocks.BuildInfo
}

func (suite *NewLoggerWithBuildInfoSuite) SetupTest() {
	suite.configurator = &mocks.Configurator{}
	suite.buildInfo = &mocks.BuildInfo{}
}

func (suite *NewLoggerWithBuildInfoSuite) TearDownTest() {
	suite.configurator.AssertExpectations(suite.T())
	suite.buildInfo.AssertExpectations(suite.T())
}

func (suite *NewLoggerWithBuildInfoSuite) TestNewLoggerError() {
	mockErr := errors.New("mock error")

	suite.configurator.On("Verbosity", mock.AnythingOfType("int")).Return(mockErr).Once()

	_, err := golflog.NewLoggerWithBuildInfo(suite.configurator, nil, "test", 0)
	if suite.Error(err) {
		suite.ErrorIs(err, mockErr)
	}
}

func (suite *NewLoggerWithBuildInfoSuite) TestBuildInfo() {
	buf, logger := bufferLogger()

	suite.configurator.On("Verbosity", mock.AnythingOfType("int")).Return(nil).Once()
	suite.configurator.On("Build").Return(logger, nil)
	suite.buildInfo.On("Version").Return("version")
	suite.buildInfo.On("Commit").Return("commit")
	suite.buildInfo.On("Date").Return("date")
	suite.buildInfo.On("Time").Return("time")

	log, err := golflog.NewLoggerWithBuildInfo(suite.configurator, suite.buildInfo, "test", 0)
	if suite.NoError(err) {
		suite.NotNil(log)
		log.Info("test")
		suite.Equal(
			`test"level"=0 "msg"="test"`+
				` "build_version"="version"`+
				` "build_commit"="commit"`+
				` "build_date"="date"`+
				` "build_time"="time"`,
			buf.String(),
		)
	}
}

func TestNewLoggerWithBuildInfo(t *testing.T) {
	suite.Run(t, new(NewLoggerWithBuildInfoSuite))
}
