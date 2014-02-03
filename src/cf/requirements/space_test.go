package requirements

import (
	"cf"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	mr "github.com/tjarratt/mr_t"
	testapi "testhelpers/api"
	testassert "testhelpers/assert"
	testterm "testhelpers/terminal"
)

func init() {
	Describe("Testing with ginkgo", func() {
		It("TestSpaceReqExecute", func() {

			space := cf.Space{}
			space.Name = "awesome-sauce-space"
			space.Guid = "my-space-guid"
			spaceRepo := &testapi.FakeSpaceRepository{Spaces: []cf.Space{space}}
			ui := new(testterm.FakeUI)

			spaceReq := newSpaceRequirement("awesome-sauce-space", ui, spaceRepo)
			success := spaceReq.Execute()

			assert.True(mr.T(), success)
			assert.Equal(mr.T(), spaceRepo.FindByNameName, "awesome-sauce-space")
			assert.Equal(mr.T(), spaceReq.GetSpace(), space)
		})
		It("TestSpaceReqExecuteWhenSpaceNotFound", func() {

			spaceRepo := &testapi.FakeSpaceRepository{FindByNameNotFound: true}
			ui := new(testterm.FakeUI)

			testassert.AssertPanic(mr.T(), testterm.FailedWasCalled, func() {
				newSpaceRequirement("foo", ui, spaceRepo).Execute()
			})
		})
	})
}
