package test_test

import (
	"fmt"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/bobymcbobs/url-redirector/src/types"
	"gopkg.in/yaml.v2"
)

var _ = Describe("API redirect tests", func() {
	currentWorkingDirectory, _ := os.Getwd()
	configStoreLocation := fmt.Sprintf("%v/%v/%v", currentWorkingDirectory, "..", "config.yaml")
	configYAML := types.ConfigYAML{
		Routes: types.Routes{
			"duck":   "https://duckduckgo.com",
			"gitlab": "https://about.gitlab.com",
			"github": "https://github.com",
		},
		Root:     "https://about.gitlab.com",
		Wildcard: "https://github.com",
	}

	BeforeEach(func() {
		file, err := os.Create(configStoreLocation)
		Expect(err).To(BeNil())
		configYAMLFmt, _ := yaml.Marshal(configYAML)
		_, err = file.Write(configYAMLFmt)
		Expect(err).To(BeNil())
	})

	It("should redirect from a page to an url", func() {
		By("visiting written redirect")
		for key, _ := range configYAML.Routes {
			resp, err := http.Get(fmt.Sprintf("http://localhost:8080/%v", key))
			Expect(err).To(BeNil(), "Request should not return errors")
			Expect(fmt.Sprintf("%v://%v", resp.Request.URL.Scheme, resp.Request.URL.Host)).To(Equal(configYAML.Routes[key]))
		}
	})

	It("should respond with 404 for non-existent url if a wildcard doesn't exist", func() {
		By("visiting an non-existent url")
		// remove wildcard in config
		file, err := os.Create(configStoreLocation)
		Expect(err).To(BeNil())
		configYAMLnew := configYAML
		configYAMLnew.Wildcard = ""
		configYAMLFmt, _ := yaml.Marshal(configYAMLnew)
		_, err = file.Write(configYAMLFmt)
		Expect(err).To(BeNil())

		pageRef := "aaaaa"
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/%v", pageRef))
		Expect(err).To(BeNil(), "Request should not return errors")
		Expect(resp.StatusCode).To(Equal(404))
		Expect(resp.Request.URL.Host).To(Equal("localhost:8080"))
	})

	It("should redirect to wildcard if is defined", func() {
		pageRef := "aaaaa"
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/%v", pageRef))
		Expect(err).To(BeNil(), "Request should not return errors")
		Expect(resp.StatusCode).To(Equal(200))
		Expect(fmt.Sprintf("%v://%v", resp.Request.URL.Scheme, resp.Request.URL.Host)).To(Equal(configYAML.Wildcard))
	})

	It("should redirect from root to a page if is defined", func() {
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/"))
		Expect(err).To(BeNil(), "Request should not return errors")
		Expect(resp.StatusCode).To(Equal(200))
		Expect(fmt.Sprintf("%v://%v", resp.Request.URL.Scheme, resp.Request.URL.Host)).To(Equal(configYAML.Root))
	})

	It("should respond with 404 from root if it is not defined", func() {
		// remove root in config
		file, err := os.Create(configStoreLocation)
		Expect(err).To(BeNil())
		configYAMLnew := configYAML
		configYAMLnew.Root = ""
		configYAMLFmt, _ := yaml.Marshal(configYAMLnew)
		_, err = file.Write(configYAMLFmt)
		Expect(err).To(BeNil())

		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/"))
		Expect(err).To(BeNil(), "Request should not return errors")
		Expect(resp.StatusCode).To(Equal(404))
		Expect(resp.Request.URL.Host).To(Equal("localhost:8080"))
	})
})
