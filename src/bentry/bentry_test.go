package bentry

import "testing"

const TEST_RESOURCES = "../../testresources/testentries"

func TestBlogContent_Load(t *testing.T) {
	blog := BlogContent{}
	blog.Load(TEST_RESOURCES)
}
