package controller

import (
	"github.com/sirupsen/logrus"

	"github.com/thank243/lightsailMon/app/node"
	"github.com/thank243/lightsailMon/common/ddns"
	"github.com/thank243/lightsailMon/common/ddns/cloudflare"
	"github.com/thank243/lightsailMon/common/ddns/google"
	"github.com/thank243/lightsailMon/common/notify"
)

func (s *Service) buildNodes() []*node.Node {
	// init notifier
	var notifier notify.Notify
	if s.conf.Notify != nil && s.conf.Notify.Enable {
		switch s.conf.Notify.Provider {
		case "pushplus":
			notifier = &notify.PushPlus{Token: s.conf.Notify.Config["pushplus_token"].(string)}
		case "telegram":
			notifier = &notify.Telegram{
				ChatID: int64(s.conf.Notify.Config["telegram_chatid"].(int)),
				Token:  s.conf.Notify.Config["telegram_token"].(string),
			}
		}
	}

	var nodes []*node.Node
	for i := range s.conf.Nodes {
		newNode := node.New(s.conf.Nodes[i])

		// set ddns client
		var (
			ddnsCli ddns.Client
			err     error
		)
		if s.conf.DDNS != nil && s.conf.DDNS.Enable {
			switch s.conf.DDNS.Provider {
			case "cloudflare":
				if ddnsCli, err = cloudflare.New(s.conf.DDNS.Config, newNode.Domain()); err != nil {
					logrus.Panicln(err)
				}
			case "google":
				if ddnsCli, err = google.New(s.conf.DDNS.Config, newNode.Domain()); err != nil {
					logrus.Panicln(err)
				}
			}
			newNode.SetDdnsClient(ddnsCli)
		}

		// set connection timeout
		if s.conf.Timeout > 0 {
			newNode.SetTimeout(s.conf.Timeout)
		}

		// set notifier
		if notifier != nil {
			newNode.SetNotifier(notifier)
		}

		nodes = append(nodes, newNode)
	}

	return nodes
}
