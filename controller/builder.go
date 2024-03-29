package controller

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/thank243/lightsailMon/app/node"
	"github.com/thank243/lightsailMon/common/ddns"
	"github.com/thank243/lightsailMon/common/ddns/cloudflare"
	"github.com/thank243/lightsailMon/common/ddns/google"
	"github.com/thank243/lightsailMon/common/notify"
	"github.com/thank243/lightsailMon/common/notify/pushplus"
	"github.com/thank243/lightsailMon/common/notify/telegram"
)

func (s *Service) buildNodes(isNotify bool, isDDNS bool) []*node.Node {
	// init notifier
	var notifier notify.Notify
	if isNotify {
		switch s.conf.Notify.Provider {
		case "pushplus":
			notifier = &pushplus.PushPlus{Token: s.conf.Notify.Config["pushplus_token"].(string)}
		case "telegram":
			notifier = &telegram.Telegram{
				ChatID: int64(s.conf.Notify.Config["telegram_chatid"].(int)),
				Token:  s.conf.Notify.Config["telegram_token"].(string),
			}
		}
	}

	// init ddnsCli
	var ddnsCli ddns.Client
	if isDDNS {
		var err error
		switch s.conf.DDNS.Provider {
		case "cloudflare":
			if ddnsCli, err = cloudflare.New(s.conf.DDNS.Config); err != nil {
				log.Panicln(err)
			}
		case "google":
			if ddnsCli, err = google.New(s.conf.DDNS.Config); err != nil {
				log.Panicln(err)
			}
		}
	}

	var nodes []*node.Node
	for i := range s.conf.Nodes {
		newNode := node.New(s.conf.Nodes[i])

		// set ddns client
		if isDDNS {
			newNode.DdnsClient = ddnsCli
		}

		// set notifier
		if isNotify {
			newNode.Notifier = notifier
		}

		// set connection timeout
		if s.conf.Timeout > 0 {
			newNode.Timeout = time.Second * time.Duration(s.conf.Timeout)
		}

		nodes = append(nodes, newNode)
	}

	return nodes
}
