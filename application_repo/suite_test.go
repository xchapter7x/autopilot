package application_repo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestApplicationRepo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Suite")
}
