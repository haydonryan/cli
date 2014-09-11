package actors_test

import (
	"errors"
	"github.com/cloudfoundry/cli/cf/actors"
	fakeBits "github.com/cloudfoundry/cli/cf/api/application_bits/fakes"
	"github.com/cloudfoundry/cli/cf/api/resources"
	"github.com/cloudfoundry/cli/cf/app_files/fakes"
	"github.com/cloudfoundry/cli/cf/models"
	"github.com/cloudfoundry/gofileutils/fileutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("Push Actor", func() {
	var (
		appBitsRepo  *fakeBits.FakeApplicationBitsRepository
		appFiles     *fakes.FakeAppFiles
		zipper       *fakes.FakeZipper
		actor        actors.PushActor
		fixturesDir  string
		appDir       string
		allFiles     []models.AppFileFields
		presentFiles []resources.AppFileResource
	)

	BeforeEach(func() {
		appBitsRepo = &fakeBits.FakeApplicationBitsRepository{}
		appFiles = &fakes.FakeAppFiles{}
		zipper = &fakes.FakeZipper{}
		actor = actors.NewPushActor(appBitsRepo, zipper, appFiles)
		fixturesDir = filepath.Join("..", "..", "fixtures", "applications")
	})

	Describe("GatherFiles", func() {
		BeforeEach(func() {
			allFiles = []models.AppFileFields{
				models.AppFileFields{Path: "example-app/.cfignore"},
				models.AppFileFields{Path: "example-app/app.rb"},
				models.AppFileFields{Path: "example-app/config.ru"},
				models.AppFileFields{Path: "example-app/Gemfile"},
				models.AppFileFields{Path: "example-app/Gemfile.lock"},
				models.AppFileFields{Path: "example-app/ignore-me"},
				models.AppFileFields{Path: "example-app/manifest.yml"},
			}

			presentFiles = []resources.AppFileResource{
				resources.AppFileResource{Path: "example-app/ignore-me"},
			}

			appDir = filepath.Join(fixturesDir, "example-app.zip")
			zipper.UnzipReturns(nil)
			appFiles.AppFilesInDirReturns(allFiles, nil)
			appBitsRepo.GetApplicationFilesReturns(presentFiles, nil)
		})

		AfterEach(func() {
		})

		Context("when the input is a zipfile", func() {
			BeforeEach(func() {
				zipper.IsZipFileReturns(true)
			})

			It("extracts the zip", func() {
				fileutils.TempDir("gather-files", func(tmpDir string, err error) {
					files, err := actor.GatherFiles(appDir, tmpDir)
					Expect(zipper.UnzipCallCount()).To(Equal(1))
					Expect(err).NotTo(HaveOccurred())
					Expect(files).To(Equal(presentFiles))
				})
			})

		})

		Context("when the input is a directory full of files", func() {
			BeforeEach(func() {
				zipper.IsZipFileReturns(false)
			})

			It("does not try to unzip the directory", func() {
				fileutils.TempDir("gather-files", func(tmpDir string, err error) {
					files, err := actor.GatherFiles(appDir, tmpDir)
					Expect(zipper.UnzipCallCount()).To(Equal(0))
					Expect(err).NotTo(HaveOccurred())
					Expect(files).To(Equal(presentFiles))
				})
			})
		})

		Context("when errors occur", func() {
			It("returns an error if it cannot unzip the files", func() {
				fileutils.TempDir("gather-files", func(tmpDir string, err error) {
					zipper.IsZipFileReturns(true)
					zipper.UnzipReturns(errors.New("error"))
					_, err = actor.GatherFiles(appDir, tmpDir)
					Expect(err).To(HaveOccurred())
				})
			})

			It("returns an error if it cannot walk the files", func() {
				fileutils.TempDir("gather-files", func(tmpDir string, err error) {
					appFiles.AppFilesInDirReturns(nil, errors.New("error"))
					_, err = actor.GatherFiles(appDir, tmpDir)
					Expect(err).To(HaveOccurred())
				})
			})

			It("returns an error if we cannot reach the cc", func() {
				fileutils.TempDir("gather-files", func(tmpDir string, err error) {
					appBitsRepo.GetApplicationFilesReturns(nil, errors.New("error"))
					_, err = actor.GatherFiles(appDir, tmpDir)
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("when using .cfignore", func() {
			BeforeEach(func() {
				appBitsRepo.GetApplicationFilesReturns(nil, nil)
				appDir = filepath.Join(fixturesDir, "exclude-a-default-cfignore")
			})

			It("includes the .cfignore file in the upload directory", func() {
				fileutils.TempDir("gather-files", func(tmpDir string, err error) {
					files, err := actor.GatherFiles(appDir, tmpDir)
					Expect(err).NotTo(HaveOccurred())

					_, err = os.Stat(filepath.Join(tmpDir, ".cfignore"))
					Expect(os.IsNotExist(err)).To(BeFalse())
					Expect(len(files)).To(Equal(0))
				})
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

// })

func executableBits(mode os.FileMode) os.FileMode {
	return mode & 0111
}
