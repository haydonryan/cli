package application_bits_test

import (
	"archive/zip"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	testapi "github.com/cloudfoundry/cli/cf/api/fakes"
	"github.com/cloudfoundry/cli/cf/api/resources"
	"github.com/cloudfoundry/cli/cf/app_files"
	"github.com/cloudfoundry/cli/cf/net"
	testconfig "github.com/cloudfoundry/cli/testhelpers/configuration"
	testnet "github.com/cloudfoundry/cli/testhelpers/net"

	. "github.com/cloudfoundry/cli/cf/api/application_bits"
	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CloudControllerApplicationBitsRepository", func() {
	var (
		fixturesDir string
	)

	BeforeEach(func() {
		cwd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		fixturesDir = filepath.Join(cwd, "../../../fixtures/applications")
	})

	var testUploadBits = func(requests ...testnet.TestRequest) (apiErr error) {
		ts, handler := testnet.NewServer(requests)
		defer ts.Close()

		configRepo := testconfig.NewRepositoryWithDefaults()
		configRepo.SetApiEndpoint(ts.URL)
		gateway := net.NewCloudControllerGateway((configRepo), time.Now)
		gateway.PollingThrottle = time.Duration(0)

		repo := NewCloudControllerApplicationBitsRepository(configRepo, gateway)
		file, err := os.Open(filepath.Join(fixturesDir, "ignored_and_resource_matched_example_app.zip"))
		if err != nil {
			log.Fatal(err)
		}

		file1 := resources.AppFileResource{Path: "app.rb", Sha1: "2474735f5163ba7612ef641f438f4b5bee00127b", Size: 51}
		file2 := resources.AppFileResource{Path: "config.ru", Sha1: "f097424ce1fa66c6cb9f5e8a18c317376ec12e05", Size: 70}
		apiErr = repo.UploadBits("my-cool-app-guid", file, []resources.AppFileResource{file1, file2})

		Expect(handler).To(HaveAllRequestsCalled())
		return
	}

	It("uploads zip files", func() {
		apiErr := testUploadBits(defaultRequests...)

		Expect(apiErr).NotTo(HaveOccurred())
	})

	It("returns a failure when uploading bits fails", func() {
		apiErr := testUploadBits([]testnet.TestRequest{
			testapi.NewCloudControllerTestRequest(testnet.TestRequest{
				Method:  "PUT",
				Path:    "/v2/apps/my-cool-app-guid/bits",
				Matcher: uploadBodyMatcher(defaultZipCheck),
				Response: testnet.TestResponse{
					Status: http.StatusCreated,
					Body: `
					{
						"metadata":{
							"guid": "my-job-guid",
							"url": "/v2/jobs/my-job-guid"
						}
					}`,
				}}),
			createProgressEndpoint("running"),
			createProgressEndpoint("failed"),
		}...)

		Expect(apiErr).To(HaveOccurred())
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
})

var matchedResources = testnet.RemoveWhiteSpaceFromBody(`[
	{
        "fn": "Gemfile",
        "sha1": "d9c3a51de5c89c11331d3b90b972789f1a14699a",
        "size": 59
    },
    {
        "fn": "Gemfile.lock",
        "sha1": "345f999aef9070fb9a608e65cf221b7038156b6d",
        "size": 229
    },
]`)

var unmatchedResources = testnet.RemoveWhiteSpaceFromBody(`[
	{
        "fn": "app.rb",
        "sha1": "2474735f5163ba7612ef641f438f4b5bee00127b",
        "size": 51
    },
    {
        "fn": "config.ru",
        "sha1": "f097424ce1fa66c6cb9f5e8a18c317376ec12e05",
        "size": 70
    }
]`)

func uploadApplicationRequest(zipCheck func(*zip.Reader)) testnet.TestRequest {
	return testapi.NewCloudControllerTestRequest(testnet.TestRequest{
		Method:  "PUT",
		Path:    "/v2/apps/my-cool-app-guid/bits",
		Matcher: uploadBodyMatcher(zipCheck),
		Response: testnet.TestResponse{
			Status: http.StatusCreated,
			Body: `
{
	"metadata":{
		"guid": "my-job-guid",
		"url": "/v2/jobs/my-job-guid"
	}
}
	`},
	})
}

var matchResourceRequest = testnet.TestRequest{
	Method: "PUT",
	Path:   "/v2/resource_match",
	Matcher: testnet.RequestBodyMatcher(testnet.RemoveWhiteSpaceFromBody(`[
    {
        "fn": "Gemfile",
        "sha1": "d9c3a51de5c89c11331d3b90b972789f1a14699a",
        "size": 59
    },
    {
        "fn": "Gemfile.lock",
        "sha1": "345f999aef9070fb9a608e65cf221b7038156b6d",
        "size": 229
    },
    {
        "fn": "app.rb",
        "sha1": "2474735f5163ba7612ef641f438f4b5bee00127b",
        "size": 51
    },
    {
        "fn": "config.ru",
        "sha1": "f097424ce1fa66c6cb9f5e8a18c317376ec12e05",
        "size": 70
    }
]`)),
	Response: testnet.TestResponse{
		Status: http.StatusOK,
		Body:   matchedResources,
	},
}

var defaultZipCheck = func(zipReader *zip.Reader) {
	Expect(len(zipReader.File)).To(Equal(2), "Wrong number of files in zip")

	var expectedPermissionBits os.FileMode
	if runtime.GOOS == "windows" {
		expectedPermissionBits = 0666
	} else {
		expectedPermissionBits = 0755
	}

	Expect(zipReader.File[0].Name).To(Equal("app.rb"))
	Expect(executableBits(zipReader.File[0].Mode())).To(Equal(executableBits(expectedPermissionBits)))

nextFile:
	for _, f := range zipReader.File {
		for _, expected := range expectedApplicationContent {
			if f.Name == expected {
				continue nextFile
			}
		}
		Fail("Expected " + f.Name + " but did not find it")
	}
}

var defaultRequests = []testnet.TestRequest{
	uploadApplicationRequest(defaultZipCheck),
	createProgressEndpoint("running"),
	createProgressEndpoint("finished"),
}

var expectedApplicationContent = []string{"app.rb", "config.ru"}

const maxMultipartResponseSizeInBytes = 4096

func uploadBodyMatcher(zipChecks func(zipReader *zip.Reader)) func(*http.Request) {
	return func(request *http.Request) {
		defer GinkgoRecover()
		err := request.ParseMultipartForm(maxMultipartResponseSizeInBytes)
		if err != nil {
			Fail(fmt.Sprintf("Failed parsing multipart form %v", err))
			return
		}
		defer request.MultipartForm.RemoveAll()

		Expect(len(request.MultipartForm.Value)).To(Equal(1), "Should have 1 value")
		valuePart, ok := request.MultipartForm.Value["resources"]
		Expect(ok).To(BeTrue(), "Resource manifest not present")
		Expect(len(valuePart)).To(Equal(1), "Wrong number of values")

		resourceManifest := valuePart[0]
		chompedResourceManifest := strings.Replace(resourceManifest, "\n", "", -1)
		Expect(chompedResourceManifest).To(Equal(unmatchedResources), "Resources do not match")

		Expect(len(request.MultipartForm.File)).To(Equal(1), "Wrong number of files")

		fileHeaders, ok := request.MultipartForm.File["application"]
		Expect(ok).To(BeTrue(), "Application file part not present")
		Expect(len(fileHeaders)).To(Equal(1), "Wrong number of files")

		applicationFile := fileHeaders[0]
		Expect(applicationFile.Filename).To(Equal("application.zip"), "Wrong file name")

		file, err := applicationFile.Open()
		if err != nil {
			Fail(fmt.Sprintf("Cannot get multipart file %v", err.Error()))
			return
		}

		length, err := strconv.ParseInt(applicationFile.Header.Get("content-length"), 10, 64)
		if err != nil {
			Fail(fmt.Sprintf("Cannot convert content-length to int %v", err.Error()))
			return
		}

		if zipChecks != nil {
			zipReader, err := zip.NewReader(file, length)
			if err != nil {
				Fail(fmt.Sprintf("Error reading zip content %v", err.Error()))
				return
			}

			zipChecks(zipReader)
		}
	}
}

func executableBits(mode os.FileMode) os.FileMode {
	return mode & 0111
}

func createProgressEndpoint(status string) (req testnet.TestRequest) {
	body := fmt.Sprintf(`
	{
		"entity":{
			"status":"%s"
		}
	}`, status)

	req.Method = "GET"
	req.Path = "/v2/jobs/my-job-guid"
	req.Response = testnet.TestResponse{
		Status: http.StatusCreated,
		Body:   body,
	}

	return
}

var matchExcludedResourceRequest = testnet.TestRequest{
	Method: "PUT",
	Path:   "/v2/resource_match",
	Matcher: testnet.RequestBodyMatcher(testnet.RemoveWhiteSpaceFromBody(`[
    {
        "fn": ".svn",
        "sha1": "0",
        "size": 0
    },
    {
        "fn": ".svn/test",
        "sha1": "456b1d3f7cfbadc66d390de79cbbb6e6a10662da",
        "size": 12
    },
    {
        "fn": "_darcs",
        "sha1": "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3",
        "size": 4
    }
]`)),
	Response: testnet.TestResponse{
		Status: http.StatusOK,
		Body:   matchedResources,
	},
}

type countingZipper struct {
	z    app_files.Zipper
	size int
}

func (cz *countingZipper) Zip(dirToZip string, targetFile *os.File) error {
	cz.size = -1
	err := cz.z.Zip(dirToZip, targetFile)
	if err != nil {
		return err
	}

	r, err := zip.OpenReader(targetFile.Name())
	if err != nil {
		return err
	}
	defer r.Close()

	cz.size = len(r.File)
	return nil
}

func (cz *countingZipper) IsZipFile(path string) bool {
	return cz.z.IsZipFile(path)
}
