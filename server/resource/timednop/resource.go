package timednop

import (
	"log"
	"os"
	"time"

	"github.com/p0pr0ck5/volchestrator/lease"
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
func (m *Manager) Associate(lease *lease.Lease) error {
	m.log.Printf("Associating %+v\n", lease)
	time.Sleep(time.Second * 10)
	return nil
}

// Disassociate implements resource.Manager
func (m *Manager) Disassociate(lease *lease.Lease) error {
	m.log.Printf("Disassociating %+v\n", lease)
	time.Sleep(time.Second * 10)
	return nil
}
