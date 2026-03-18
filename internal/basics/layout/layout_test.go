package layout

import (
	"testing"

	"github.com/romangurevitch/go-training/internal/basics/layout/pkg/files"
	"github.com/romangurevitch/go-training/internal/basics/layout/pkg/strings"
	"github.com/romangurevitch/go-training/internal/basics/layout/util"
)

func TestLayout(t *testing.T) {
	_ = util.ToUpper("string")
	_, _ = util.Open("path")
}

func TestBetterLayout(t *testing.T) {
	_ = strings.ToUpper("string")
	_, _ = files.Open("path")
}
