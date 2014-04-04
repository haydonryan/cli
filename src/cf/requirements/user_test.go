/*
                       WARNING WARNING WARNING

                Attention all potential contributors

   This testfile is not in the best state. We've been slowly transitioning
   from the built in "testing" package to using Ginkgo. As you can see, we've
   changed the format, but a lot of the setup, test body, descriptions, etc
   are either hardcoded, completely lacking, or misleading.

   For example:

   Describe("Testing with ginkgo"...)      // This is not a great description
   It("TestDoesSoemthing"...)              // This is a horrible description

   Describe("create-user command"...       // Describe the actual object under test
   It("creates a user when provided ..."   // this is more descriptive

   For good examples of writing Ginkgo tests for the cli, refer to

   src/cf/commands/application/delete_app_test.go
   src/cf/terminal/ui_test.go
   src/github.com/cloudfoundry/loggregator_consumer/consumer_test.go
*/

package requirements_test

import (
	"cf/models"
	. "cf/requirements"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	testapi "testhelpers/api"
	testassert "testhelpers/assert"
	testterm "testhelpers/terminal"
)

var _ = Describe("Testing with ginkgo", func() {
	It("TestUserReqExecute", func() {
		user := models.UserFields{}
		user.Username = "my-user"
		user.Guid = "my-user-guid"

		userRepo := &testapi.FakeUserRepository{FindByUsernameUserFields: user}
		ui := new(testterm.FakeUI)

		userReq := NewUserRequirement("foo", ui, userRepo)
		success := userReq.Execute()

		Expect(success).To(BeTrue())
		Expect(userRepo.FindByUsernameUsername).To(Equal("foo"))
		Expect(userReq.GetUser()).To(Equal(user))
	})

	It("TestUserReqWhenUserDoesNotExist", func() {
		userRepo := &testapi.FakeUserRepository{FindByUsernameNotFound: true}
		ui := new(testterm.FakeUI)

		testassert.AssertPanic(testterm.FailedWasCalled, func() {
			NewUserRequirement("foo", ui, userRepo).Execute()
		})

		testassert.SliceContains(ui.Outputs, testassert.Lines{
			{"FAILED"},
			{"not found"},
		})
	})
})
