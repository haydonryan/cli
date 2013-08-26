package api

import (
	"archive/zip"
	"bytes"
	"cf"
	"cf/configuration"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

type ApplicationRepository interface {
	FindByName(config *configuration.Configuration, name string) (app cf.Application, err error)
	SetEnv(config *configuration.Configuration, app cf.Application, name string, value string) (err error)
	Create(config *configuration.Configuration, newApp cf.Application) (createdApp cf.Application, err error)
	Upload(config *configuration.Configuration, app cf.Application) (err error)
}

type CloudControllerApplicationRepository struct {
}

func (repo CloudControllerApplicationRepository) FindByName(config *configuration.Configuration, name string) (app cf.Application, err error) {
	apps, err := findApplications(config)
	lowerName := strings.ToLower(name)
	if err != nil {
		return
	}

	for _, a := range apps {
		if strings.ToLower(a.Name) == lowerName {
			return a, nil
		}
	}

	err = errors.New("Application not found")
	return
}

func (repo CloudControllerApplicationRepository) SetEnv(config *configuration.Configuration, app cf.Application, name string, value string) (err error) {
	path := fmt.Sprintf("%s/v2/apps/%s", config.Target, app.Guid)
	data := fmt.Sprintf(`{"environment_json":{"%s":"%s"}}`, name, value)
	request, err := NewAuthorizedRequest("PUT", path, config.AccessToken, strings.NewReader(data))
	if err != nil {
		return
	}

	err = PerformRequest(request)
	return
}

func (repo CloudControllerApplicationRepository) Create(config *configuration.Configuration, newApp cf.Application) (createdApp cf.Application, err error) {
	path := fmt.Sprintf("%s/v2/apps", config.Target)
	data := fmt.Sprintf(
		`{"space_guid":"%s","name":"%s","instances":1,"buildpack":null,"command":null,"memory":256,"stack_guid":null}`,
		config.Space.Guid, newApp.Name,
	)
	request, err := NewAuthorizedRequest("POST", path, config.AccessToken, strings.NewReader(data))
	if err != nil {
		return
	}

	resource := new(Resource)
	err = PerformRequestForBody(request, resource)

	if err != nil {
		return
	}

	createdApp.Guid = resource.Metadata.Guid
	createdApp.Name = resource.Entity.Name
	return
}

func (repo CloudControllerApplicationRepository) Upload(config *configuration.Configuration, app cf.Application) (err error) {
	url := fmt.Sprintf("%s/v2/apps/%s/bits", config.Target, app.Guid)
	dir, err := os.Getwd()
	if err != nil {
		return
	}

	zipBuffer, err := zipApplication(dir)
	if err != nil {
		return
	}

	body, boundary, err := createApplicationUploadBody(zipBuffer)
	if err != nil {
		return
	}

	request, err := NewAuthorizedRequest("PUT", url, config.AccessToken, body)
	contentType := fmt.Sprintf("multipart/form-data; boundary=%s", boundary)
	request.Header.Set("Content-Type", contentType)
	if err != nil {
		return
	}

	err = PerformRequest(request)
	return
}

func zipApplication(dirToZip string) (zipBuffer *bytes.Buffer, err error) {
	zipBuffer = new(bytes.Buffer)
	writer := zip.NewWriter(zipBuffer)

	addFileToZip := func(path string, f os.FileInfo, inErr error) (err error) {
		err = inErr
		if err != nil {
			return
		}

		if f.IsDir() {
			return
		}

		fileName := strings.TrimPrefix(path, dirToZip+"/")
		zipFile, err := writer.Create(fileName)
		if err != nil {
			return
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return
		}

		_, err = zipFile.Write(content)
		if err != nil {
			return
		}

		return
	}

	err = filepath.Walk(dirToZip, addFileToZip)

	if err != nil {
		return
	}

	err = writer.Close()
	return
}

func createApplicationUploadBody(zipBuffer *bytes.Buffer) (body *bytes.Buffer, boundary string, err error) {
	body = new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	boundary = writer.Boundary()

	part, err := writer.CreateFormField("resources")
	if err != nil {
		return
	}

	_, err = io.Copy(part, bytes.NewBufferString("[]"))
	if err != nil {
		return
	}

	part, err = createZipPartWriter(zipBuffer, writer)
	if err != nil {
		return
	}

	_, err = io.Copy(part, zipBuffer)
	if err != nil {
		return
	}

	err = writer.Close()
	return
}

func createZipPartWriter(zipBuffer *bytes.Buffer, writer *multipart.Writer) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="application"; filename="application.zip"`)
	h.Set("Content-Type", "application/zip")
	h.Set("Content-Length", fmt.Sprintf("%d", zipBuffer.Len()))
	h.Set("Content-Transfer-Encoding", "binary")
	return writer.CreatePart(h)
}

func findApplications(config *configuration.Configuration) (apps []cf.Application, err error) {
	path := fmt.Sprintf("%s/v2/spaces/%s/apps", config.Target, config.Space.Guid)
	request, err := NewAuthorizedRequest("GET", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	response := new(ApiResponse)
	err = PerformRequestForBody(request, response)
	if err != nil {
		return
	}

	for _, r := range response.Resources {
		apps = append(apps, cf.Application{r.Entity.Name, r.Metadata.Guid})
	}

	return
}
