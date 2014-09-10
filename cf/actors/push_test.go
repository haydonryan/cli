package actors_test

import (
	. "github.com/onsi/ginkgo"
	"os"
	//	. "github.com/onsi/gomega"
)

var _ = Describe("Push Actor", func() {
	Describe("GatherFiles", func() {
		BeforeEach(func() {

		})

		AfterEach(func() {
		})

		Context("when the input is a zipfile", func() {
			It("extracts the zip", func() {

			})

			It("copies the files into the upload dir", func() {

			})

		})

		Context("when the input is a directory full of files", func() {
			It("copies the files into the upload dir", func() {

			})
		})
	})

	Describe(".UploadApp", func() {
		//A lot of stuff I assume
	})
})

// PIt("returns an error when the directory to upload does not exist", func() {
// config := testconfig.NewRepository()
// gateway := net.NewCloudControllerGateway((config), time.Now)
// zipper := &app_files.ApplicationZipper{}

// repo := NewCloudControllerApplicationBitsRepository(config, gateway)

// apiErr := repo.UploadApp("app-guid", "/foo/bar", func(_ string, _, _ int64) {})
// Expect(apiErr).To(HaveOccurred())
// Expect(apiErr.Error()).To(ContainSubstring(filepath.Join("foo", "bar")))
// })

// Context("when uploading a directory", func() {
// 	var appPath string

// 	BeforeEach(func() {
// 		appPath = filepath.Join(fixturesDir, "example-app")

// 		// the executable bit is the only bit we care about here
// 		err := os.Chmod(filepath.Join(appPath, "Gemfile"), 0467)
// 		Expect(err).NotTo(HaveOccurred())
// 	})

// 	AfterEach(func() {
// 		os.Chmod(filepath.Join(appPath, "Gemfile"), 0666)
// 	})

// It("preserves the executable bits when uploading app files", func() {
// 	apiErr := testUploadBits(defaultRequests...)
// 	Expect(apiErr).NotTo(HaveOccurred())

// 	Expect(reportedFilePath).To(Equal(appPath))
// 	Expect(reportedFileCount).To(Equal(int64(len(expectedApplicationContent))))
// 	Expect(reportedUploadSize).To(Equal(int64(532)))
// })

// PContext("when there are no files to upload", func() {
// It("makes a request without a zipfile", func() {
// 	emptyDir := filepath.Join(fixturesDir, "empty-dir")
// 	err := testUploadApp(
// 		testnet.TestRequest{
// 			Method:  "PUT",
// 			Path:    "/v2/resource_match",
// 			Matcher: testnet.RequestBodyMatcher("[]"),
// 			Response: testnet.TestResponse{
// 				Status: http.StatusOK,
// 				Body:   "[]",
// 			},
// 		},
// 		testapi.NewCloudControllerTestRequest(testnet.TestRequest{
// 			Method: "PUT",
// 			Path:   "/v2/apps/my-cool-app-guid/bits",
// 			Matcher: func(request *http.Request) {
// 				err := request.ParseMultipartForm(maxMultipartResponseSizeInBytes)
// 				Expect(err).NotTo(HaveOccurred())
// 				defer request.MultipartForm.RemoveAll()

// 				Expect(len(request.MultipartForm.Value)).To(Equal(1), "Should have 1 value")
// 				valuePart, ok := request.MultipartForm.Value["resources"]

// 				Expect(ok).To(BeTrue(), "Resource manifest not present")
// 				Expect(valuePart).To(Equal([]string{"[]"}))
// 				Expect(request.MultipartForm.File).To(BeEmpty())
// 			},
// 			Response: testnet.TestResponse{
// 				Status: http.StatusCreated,
// 				Body: `
// 					{
// 						"metadata":{
// 							"guid": "my-job-guid",
// 							"url": "/v2/jobs/my-job-guid"
// 						}
// 					}`,
// 			}}),
// 		createProgressEndpoint("running"),
// 		createProgressEndpoint("finished"),
// 	)

// 	Expect(err).NotTo(HaveOccurred())
// 	Expect(reportedFileCount).To(Equal(int64(0)))
// 	Expect(reportedUploadSize).To(Equal(int64(0)))
// 	Expect(reportedZipCount).To(Equal(-1))
// 	Expect(reportedFilePath).To(Equal(emptyDir))
// })
// })
// })

// PContext("when excluding a default ignored item", func() {
// var appPath string

// BeforeEach(func() {
// 	appPath = filepath.Join(fixturesDir, "exclude-a-default-cfignore")
// })

// It("includes the ignored item", func() {
// 	err := testUploadApp(appPath,
// 		matchExcludedResourceRequest,
// 		uploadApplicationRequest(func(zipReader *zip.Reader) {
// 			Expect(len(zipReader.File)).To(Equal(3), "Wrong number of files in zip")
// 		}),
// 		createProgressEndpoint("running"),
// 		createProgressEndpoint("finished"),
// 	)

// 	Expect(err).NotTo(HaveOccurred())
// 	Expect(reportedFileCount).To(Equal(int64(3)))
// 	Expect(reportedUploadSize).To(Equal(int64(354)))
// 	Expect(reportedZipCount).To(Equal(3))
// })
// })

func executableBits(mode os.FileMode) os.FileMode {
	return mode & 0111
}
