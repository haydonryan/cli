package manifest_test

import (
	"cf/manifest"
	"generic"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	mr "github.com/tjarratt/mr_t"
	"runtime"
	"strings"
	testassert "testhelpers/assert"
)

func testManifestWithAbsolutePathOnPosix(t mr.TestingT) {
	m, err := manifest.NewManifest("/some/path", generic.NewMap(map[string]interface{}{
		"applications": []interface{}{
			map[string]interface{}{
				"path": "/another/path-segment",
			},
		},
	}))

	assert.NoError(t, err)
	assert.Equal(t, *m.Applications[0].Path, "/another/path-segment")
}

func testManifestWithAbsolutePathOnWindows(t mr.TestingT) {
	m, err := manifest.NewManifest(`C:\some\path`, generic.NewMap(map[string]interface{}{
		"applications": []interface{}{
			map[string]interface{}{
				"path": `C:\another\path`,
			},
		},
	}))

	assert.NoError(t, err)
	assert.Equal(t, *m.Applications[0].Path, `C:\another\path`)
}

func init() {
	Describe("Testing with ginkgo", func() {
		It("TestManifestWithGlobalAndAppSpecificProperties", func() {
			m, err := manifest.NewManifest("/some/path", generic.NewMap(map[string]interface{}{
				"instances": "3",
				"memory":    "512M",
				"applications": []interface{}{
					map[string]interface{}{
						"name":     "bitcoin-miner",
						"no-route": true,
					},
				},
			}))
			assert.NoError(mr.T(), err)

			apps := m.Applications
			assert.Equal(mr.T(), *apps[0].InstanceCount, 3)
			assert.Equal(mr.T(), *apps[0].Memory, uint64(512))
			assert.True(mr.T(), *apps[0].NoRoute)
		})

		It("TestManifestWithInvalidMemory", func() {
			_, err := manifest.NewManifest("/some/path", generic.NewMap(map[string]interface{}{
				"instances": "3",
				"memory":    "512",
				"applications": []interface{}{
					map[string]interface{}{
						"name": "bitcoin-miner",
					},
				},
			}))

			assert.Error(mr.T(), err)
			assert.Contains(mr.T(), err.Error(), "memory")
		})

		It("TestManifestWithTimeoutSetsHealthCheckTimeout", func() {
			m, err := manifest.NewManifest("/some/path", generic.NewMap(map[string]interface{}{
				"applications": []interface{}{
					map[string]interface{}{
						"name":    "bitcoin-miner",
						"timeout": "360",
					},
				},
			}))

			assert.NoError(mr.T(), err)
			assert.Equal(mr.T(), *m.Applications[0].HealthCheckTimeout, 360)
		})

		It("TestManifestWithEmptyEnvVarIsInvalid", func() {
			_, err := manifest.NewManifest("/some/path", generic.NewMap(map[string]interface{}{
				"env": generic.NewMap(map[string]interface{}{
					"bar": nil,
				}),
				"applications": []interface{}{
					map[string]interface{}{
						"name": "bad app",
					},
				},
			}))

			assert.Error(mr.T(), err)
			assert.Contains(mr.T(), err.Error(), "env var 'bar' should not be null")
		})

		It("TestManifestWithAbsolutePath", func() {
			if runtime.GOOS == "windows" {
				testManifestWithAbsolutePathOnWindows(mr.T())
			} else {
				testManifestWithAbsolutePathOnPosix(mr.T())
			}
		})

		It("TestManifestWithRelativePath", func() {
			m, err := manifest.NewManifest("/some/path", generic.NewMap(map[string]interface{}{
				"applications": []interface{}{
					map[string]interface{}{
						"path": "../another/path-segment",
					},
				},
			}))

			assert.NoError(mr.T(), err)
			if runtime.GOOS == "windows" {
				assert.Equal(mr.T(), *m.Applications[0].Path, "\\some\\another\\path-segment")
			} else {
				assert.Equal(mr.T(), *m.Applications[0].Path, "/some/another/path-segment")
			}
		})

		It("TestParsingManifestWithNulls", func() {
			_, errs := manifest.NewManifest("/some/path", generic.NewMap(map[string]interface{}{
				"applications": []interface{}{
					map[string]interface{}{
						"buildpack":  nil,
						"disk_quota": nil,
						"domain":     nil,
						"host":       nil,
						"name":       nil,
						"path":       nil,
						"stack":      nil,
						"memory":     nil,
						"instances":  nil,
						"timeout":    nil,
						"no-route":   nil,
						"services":   nil,
						"env":        nil,
					},
				},
			}))

			assert.Error(mr.T(), errs)
			errorSlice := strings.Split(errs.Error(), "\n")
			manifestKeys := []string{"buildpack", "disk_quota", "domain", "host", "name", "path", "stack",
				"memory", "instances", "timeout", "no-route", "services", "env"}

			for _, key := range manifestKeys {
				testassert.SliceContains(mr.T(), errorSlice, testassert.Lines{{key, "not be null"}})
			}
		})

		It("TestParsingManifestWithPropertiesReturnsErrors", func() {
			_, err := manifest.NewManifest("/some/path", generic.NewMap(map[string]interface{}{
				"applications": []interface{}{
					map[string]interface{}{
						"env": map[string]interface{}{
							"bar": "many-${foo}-are-cool",
						},
					},
				},
			}))

			assert.Error(mr.T(), err)
			assert.Contains(mr.T(), err.Error(), "Properties are not supported. Found property '${foo}'")
		})

		It("TestParsingManifestWithNullCommand", func() {
			m, err := manifest.NewManifest("/some/path", generic.NewMap(map[string]interface{}{
				"applications": []interface{}{
					map[string]interface{}{
						"command": nil,
					},
				},
			}))

			assert.NoError(mr.T(), err)
			assert.Equal(mr.T(), *m.Applications[0].Command, "")
		})

		It("TestParsingEmptyManifestDoesNotSetCommand", func() {
			m, err := manifest.NewManifest("/some/path", generic.NewMap(map[string]interface{}{
				"applications": []interface{}{
					map[string]interface{}{},
				},
			}))

			assert.NoError(mr.T(), err)
			assert.Nil(mr.T(), m.Applications[0].Command)
		})
	})
}
