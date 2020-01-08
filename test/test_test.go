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
	routeHostsForTests := []string{"localhost", "localhost2", "localhost3"}
	routeHostForSingleTest := routeHostsForTests[0]
	configYAML := getDefaultConfigYAMLvalues()

	BeforeEach(func() {
		configYAML = getDefaultConfigYAMLvalues()
		err := writeDefaultTestConfig(configStoreLocation, configYAML)
		Expect(err).To(BeNil())
	})

	It("should redirect from a page to an url", func() {
		By("visiting written redirect")
		// go through each host
		for _, host := range routeHostsForTests {
			// go through each route
			for key, _ := range configYAML[host].Routes {
				resp, err := httpGetWithHeader(fmt.Sprintf("http://localhost:8080/%v", key), host)
				Expect(err).To(BeNil(), "Request should not return errors")
				Expect(fmt.Sprintf("%v://%v", resp.Request.URL.Scheme, resp.Request.URL.Host)).To(Equal(configYAML[host].Routes[key]))
			}
		}
	})

	It("should redirect to wildcard if is defined", func() {
		configYAML = getDefaultConfigYAMLvalues()
		err := writeDefaultTestConfig(configStoreLocation, configYAML)
		Expect(err).To(BeNil())
		pageRef := "aaaaa"
		resp, err := httpGetWithHeader(fmt.Sprintf("http://localhost:8080/%v", pageRef), routeHostForSingleTest)
		Expect(err).To(BeNil(), "Request should not return errors")
		Expect(resp.StatusCode).To(Equal(200))
		Expect(fmt.Sprintf("%v://%v", resp.Request.URL.Scheme, resp.Request.URL.Host)).To(Equal(configYAML[routeHostForSingleTest].Wildcard))
	})

	It("should redirect from root to a page if is defined", func() {
		configYAML = getDefaultConfigYAMLvalues()
		err := writeDefaultTestConfig(configStoreLocation, configYAML)
		Expect(err).To(BeNil())
		resp, err := httpGetWithHeader(fmt.Sprintf("http://localhost:8080/"), routeHostForSingleTest)
		Expect(err).To(BeNil(), "Request should not return errors")
		Expect(resp.StatusCode).To(Equal(200))
		Expect(fmt.Sprintf("%v://%v", resp.Request.URL.Scheme, resp.Request.URL.Host)).To(Equal(configYAML[routeHostForSingleTest].Root))
	})

	It("should respond with 404 for non-existent url if a wildcard doesn't exist", func() {
		err := removeRootAndWildcardFromConfigYAML(configStoreLocation, configYAML, routeHostForSingleTest)
		Expect(err).To(BeNil())
		By("visiting an non-existent url")
		pageRef := "aaaaa"
		resp, err := httpGetWithHeader(fmt.Sprintf("http://localhost:8080/%v", pageRef), routeHostForSingleTest)
		Expect(err).To(BeNil(), "Request should not return errors")
		Expect(resp.StatusCode).To(Equal(404))
		Expect(resp.Request.URL.Host).To(Equal("localhost:8080"))
	})

	It("should respond with 404 from root if it is not defined", func() {
		err := removeRootAndWildcardFromConfigYAML(configStoreLocation, configYAML, routeHostForSingleTest)
		Expect(err).To(BeNil())
		resp, err := httpGetWithHeader("http://localhost:8080/", routeHostForSingleTest)
		Expect(err).To(BeNil(), "Request should not return errors")
		Expect(resp.StatusCode).To(Equal(404))
		Expect(resp.Request.URL.Host).To(Equal("localhost:8080"))
	})

	It("should redirect any host and any URL if configured to", func() {
		wildcardSite := "https://github.com"
		err := writeDefaultTestConfig(configStoreLocation, types.RouteHosts{
			"*": types.RouteHost{
				Root:     wildcardSite,
				Wildcard: wildcardSite,
			},
		})
		Expect(err).To(BeNil())
		attemptURLs := []string{"a", "abc", "test"}
		for _, url := range attemptURLs {
			resp, err := httpGetWithHeader(fmt.Sprintf("http://localhost:8080/%v", url), routeHostForSingleTest)
			Expect(err).To(BeNil(), "Request should not return errors")
			Expect(resp.StatusCode).To(Equal(200))
			Expect(fmt.Sprintf("%v://%v", resp.Request.URL.Scheme, resp.Request.URL.Host)).To(Equal(wildcardSite))
		}
	})
})

func httpGetWithHeader(url string, host string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Host = host
	client := &http.Client{}
	resp, err = client.Do(req)
	return resp, err
}

func removeRootAndWildcardFromConfigYAML(configStoreLocation string, configYAML types.RouteHosts, routeHostForSingleTest string) (err error) {
	// remove root in config
	file, err := os.Create(configStoreLocation)
	Expect(err).To(BeNil())
	configYAMLnew := configYAML

	// remove wildcard from host localhost
	configYAMLnewHostOne := types.RouteHost(configYAMLnew[routeHostForSingleTest])
	configYAMLnewHostOne.Root = ""
	configYAMLnewHostOne.Wildcard = ""
	configYAMLnew[routeHostForSingleTest] = configYAMLnewHostOne

	// remove wildcard from host wildcard
	configYAMLnewHostOne = types.RouteHost(configYAMLnew["*"])
	configYAMLnewHostOne.Root = ""
	configYAMLnewHostOne.Wildcard = ""
	configYAMLnew["*"] = configYAMLnewHostOne

	configYAMLFmt, _ := yaml.Marshal(configYAMLnew)
	_, err = file.Write(configYAMLFmt)
	return err
}

func writeDefaultTestConfig(configStoreLocation string, configYAML types.RouteHosts) (err error) {
	file, err := os.Create(configStoreLocation)
	Expect(err).To(BeNil())
	configYAMLFmt, _ := yaml.Marshal(configYAML)
	_, err = file.Write(configYAMLFmt)
	return err
}

func getDefaultConfigYAMLvalues() (configYAML types.RouteHosts) {
	return types.RouteHosts{
		"localhost": types.RouteHost{
			Routes: types.Routes{
				"a": "https://duckduckgo.com",
				"b": "https://about.gitlab.com",
				"c": "https://github.com",
			},
			Root:     "https://about.gitlab.com",
			Wildcard: "https://github.com",
		},
		"localhost2": types.RouteHost{
			Routes: types.Routes{
				"a": "https://about.gitlab.com",
				"b": "https://duckduckgo.com",
				"c": "https://github.com",
			},
			Root:     "https://about.gitlab.com",
			Wildcard: "https://github.com",
		},
		"*": types.RouteHost{
			Routes: types.Routes{
				"a": "https://about.gitlab.com",
				"b": "https://github.com",
				"c": "https://duckduckgo.com",
			},
			Root:     "https://about.gitlab.com",
			Wildcard: "https://github.com",
		},
	}
}
