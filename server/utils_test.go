package server

import (
	"testing"
)

func TestValidateProjectName(t *testing.T) {
	assert(validateProjectName("hello"), "hello works")
	assert(validateProjectName("taskcluster-auth"), "taskcluster-auth works")
	assert(validateProjectName("foobar"), "foobar works")
	assert(validateProjectName("abc_123"), "abc_123 works")
	assert(validateProjectName("abc-123"), "abc-123 works")
	assert(validateProjectName("ec2-manager"), "ec2-manager works")

	assert(!validateProjectName("^f"), "^f doesn't work")
	assert(!validateProjectName("qed abc"), "qed abc doesn't work")
	assert(!validateProjectName("🐊"), "🐊 doesn't work")
	assert(!validateProjectName("*"), "* doesn't work")
}
