package test_test

import (
	"fmt"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("API redirect tests", func() {
	currentWorkingDirectory, _ := os.Getwd()
	configStoreLocation := fmt.Sprintf("%v/%v/%v", currentWorkingDirectory, "..", "config.yaml")
	configYAML := map[string]string{
		"duck":   "https://duckduckgo.com",
		"gitlab": "https://about.gitlab.com",
		"github": "https://github.com",
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
		//pageRef := "duck"
		for key, _ := range configYAML {
			resp, err := http.Get(fmt.Sprintf("http://localhost:8080/%v", key))
			Expect(err).To(BeNil(), "Request should not return errors")
			Expect(fmt.Sprintf("%v://%v", resp.Request.URL.Scheme, resp.Request.URL.Host)).To(Equal(configYAML[key]))
		}
	})

	It("should respond with 404 for non-existent url", func() {
		By("visiting an non-existent url")
		pageRef := "aaaaa"
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/%v", pageRef))
		Expect(err).To(BeNil(), "Request should not return errors")
		Expect(resp.StatusCode).To(Equal(404))
		Expect(resp.Request.URL.Host).To(Equal("localhost:8080"))
	})
})
