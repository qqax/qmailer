package qmailer

import (
	"fmt"
	mail "github.com/xhit/go-simple-mail/v2"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type result struct {
	Err   error
	Value interface{}
}

var (
	poolInstance *smtpPool
	once         sync.Once
)

type smtpPool struct {
	id, Pending, Active, GoMax, attempts uint32
	timeout                              time.Duration
}

func (sp *smtpPool) Connect(in chan *mail.Email, out chan *result) {
	atomic.AddUint32(&sp.Active, 1)
	defer atomic.AddUint32(&sp.Active, ^uint32(0))

	id := atomic.AddUint32(&sp.id, 1)

	smtp, err := getSmtpClient()
	if err != nil {
		out <- &result{Err: fmt.Errorf("smtpConnection id%v got getSmtpClient error: %s", id, err)}
		return
	}

	for {
		atomic.AddUint32(&sp.Pending, 1)
		timer := time.NewTimer(sp.timeout)

		select {
		case <-timer.C:

			atomic.AddUint32(&sp.Pending, ^uint32(0))
			if poolInstance.Active == 1 {
				close(out)
			}

			return

		case email, ok := <-in:

			atomic.AddUint32(&sp.Pending, ^uint32(0))
			timer.Stop()

			if ok {
				var i uint32
				for i = 0; i < sp.attempts; i++ {
					time.Sleep(time.Duration(i) * time.Second)
					err = email.Send(smtp)

					if err == nil {
						out <- &result{
							Value: fmt.Sprintf("SEND attempt #%v, smtpConnection id#%v, email %s",
								i, id, strings.Join(email.GetRecipients(), ", ")),
						}
						break
					} else if i+1 == sp.attempts {
						out <- &result{
							Err: fmt.Errorf("ERROR attempt #%v, smtpConnection id#%v, mail %s, err: %s",
								i, id, strings.Join(email.GetRecipients(), ", "), err),
						}
						break
					}
				}
				continue
			}

			if poolInstance.Active == 1 {
				close(out)
			}

			return
		}
	}
}
