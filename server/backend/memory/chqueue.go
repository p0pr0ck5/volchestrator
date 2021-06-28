package memory

import (
	"sync"

	"github.com/p0pr0ck5/volchestrator/server/notification"
	"github.com/pkg/errors"
)

type ChQueue struct {
	ch chan *notification.Notification

	count     uint64
	countLock sync.Mutex

	wg sync.WaitGroup

	shutdownCh chan struct{}
}

func NewChQueue() (*ChQueue, error) {
	c := &ChQueue{
		ch:         make(chan *notification.Notification),
		shutdownCh: make(chan struct{}),
	}

	return c, nil
}

func MustNewChQueue() *ChQueue {
	queue, err := NewChQueue()
	if err != nil {
		panic(err)
	}
	return queue
}

func (c *ChQueue) Read() (<-chan *notification.Notification, error) {
	select {
	case <-c.shutdownCh:
		return nil, errors.New("closed")
	default:
	}

	return c.ch, nil
}

func (c *ChQueue) Write(n *notification.Notification) error {
	select {
	case <-c.shutdownCh:
		return errors.New("closed")
	default:
	}

	c.countLock.Lock()
	c.count++
	n.SetMessageID(c.count)
	c.countLock.Unlock()

	go func(n *notification.Notification) {
		c.wg.Add(1)
		defer c.wg.Done()

		select {
		case <-c.shutdownCh:
			return
		default:
		}

		select {
		case <-c.shutdownCh:
		case c.ch <- n:
		}
	}(n)

	return nil
}

func (c *ChQueue) Close() error {
	select {
	case <-c.shutdownCh:
		return errors.New("closed")
	default:
	}

	close(c.shutdownCh)

	go func() {
		for range c.ch {
		}
	}()

	c.wg.Wait()
	close(c.ch)
	return nil
}
