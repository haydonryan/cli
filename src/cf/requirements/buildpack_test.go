package requirements

import (
	"cf"
	"github.com/stretchr/testify/assert"
	testapi "testhelpers/api"
	testassert "testhelpers/assert"
	testterm "testhelpers/terminal"
	"testing"
)

func TestBuildpackReqExecute(t *testing.T) {
	buildpack := cf.Buildpack{}
	buildpack.Name = "my-buildpack"
	buildpack.Guid = "my-buildpack-guid"
	buildpackRepo := &testapi.FakeBuildpackRepository{FindByNameBuildpack: buildpack}
	ui := new(testterm.FakeUI)

	buildpackReq := newBuildpackRequirement("foo", ui, buildpackRepo)
	success := buildpackReq.Execute()

	assert.True(t, success)
	assert.Equal(t, buildpackRepo.FindByNameName, "foo")
	assert.Equal(t, buildpackReq.GetBuildpack(), buildpack)
}

func TestBuildpackReqExecuteWhenBuildpackNotFound(t *testing.T) {
	buildpackRepo := &testapi.FakeBuildpackRepository{FindByNameNotFound: true}
	ui := new(testterm.FakeUI)

	buildpackReq := newBuildpackRequirement("foo", ui, buildpackRepo)

	testassert.AssertPanic(t, testterm.FailedWasCalled, func() {
		buildpackReq.Execute()
	})
}
