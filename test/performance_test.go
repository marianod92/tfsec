package test

import (
	"fmt"
	"testing"

	"github.com/aquasecurity/tfsec/internal/pkg/parser"
	"github.com/aquasecurity/tfsec/internal/pkg/scanner"
	"github.com/aquasecurity/tfsec/internal/pkg/testutil/filesystem"
)

func BenchmarkCalculate(b *testing.B) {
	fs, err := filesystem.New()
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	createBadBlocks(fs)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		blocks, err := parser.New(fs.RealPath("/project"), parser.OptionStopOnHCLError()).ParseDirectory()
		if err != nil {
			panic(err)
		}
		_, _ = scanner.New().Scan(blocks)
	}
}

func createBadBlocks(fs *filesystem.FileSystem) {
	_ = fs.WriteTextFile("/project/main.tf", `
		module "something" {
			source = "../modules/problem"
		}
		`)

	for _, rule := range scanner.GetRegisteredRules() {
		for i, bad := range rule.Base.Rule().Terraform.BadExamples {
			_ = fs.WriteTextFile(fmt.Sprintf("/modules/problem/%s-%d.tf", rule.ID(), i), bad)
		}
	}
}