package maintner

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/build/maintner/maintpb"
)

type MutationLogger interface {
	Log(m *maintpb.Mutation) error
}

type DiskMutationLogger struct {
	root string
}

func NewDiskMutationLogger(root string) *DiskMutationLogger {
	return &DiskMutationLogger{root: root}
}

func (d *DiskMutationLogger) filename() string {
	now := time.Now().UTC()
	return filepath.Join(d.root, fmt.Sprintf("maintner-%s.proto", now.Format("2006-01-02")))
}

func (d *DiskMutationLogger) Log(m *maintpb.Mutation) error {
	f, err := os.OpenFile(d.filename(), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		// TODO be more graceful about failure
		return err
	}
	data, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	return f.Close()
}
