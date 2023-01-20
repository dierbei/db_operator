package builders

import (
	"fmt"
	"testing"
)

func TestDeployBuilder_Build(t *testing.T) {
	b, err := NewDeployBuilder("test", "default")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(b.Replicas(3).Build())
}
