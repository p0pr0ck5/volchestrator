package timednop

import (
	"log"
	"os"
	"time"

	"github.com/p0pr0ck5/volchestrator/server"
)

// Manager implements resource.Manager
type Manager struct {
	log *log.Logger
}

// New returns a new Manager
func New() *Manager {
	return &Manager{
		log: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
	}
}

// Associate implements resource.Manager
func (m *Manager) Associate(v *server.Volume) error {
	m.log.Printf("Associating %+v\n", v)
	time.Sleep(time.Second)
	return nil
}

// Disassociate implements resource.Manager
func (m *Manager) Disassociate(v *server.Volume) error {
	m.log.Printf("Disassociating %+v\n", v)
	time.Sleep(time.Second)
	return nil
}
