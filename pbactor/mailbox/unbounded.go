package mailbox

import (
	"keywea.com/cloud/pblib/pbutils/queue/goring"
	"keywea.com/cloud/pblib/pbutils/queue/mpsc"
)

type unboundedMailboxQueue struct {
	userMailbox *goring.Queue
}

func (q *unboundedMailboxQueue) Push(m interface{}) {
	q.userMailbox.Push(m)
}

func (q *unboundedMailboxQueue) Pop() interface{} {
	m, o := q.userMailbox.Pop()
	if o {
		return m
	}
	return nil
}

// Unbounded returns a producer which creates an unbounded mailbox
func Unbounded(mailboxStats ...Statistics) Producer {
	return func() Mailbox {
		q := &unboundedMailboxQueue{
			userMailbox: goring.New(10),
		}
		return &defaultMailbox{
			systemMailbox: mpsc.New(),
			userMailbox:   q,
			mailboxStats:  mailboxStats,
		}
	}
}
