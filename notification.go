package agollo

import (
	"encoding/json"
	"sync"
)

type notification struct {
	NamespaceName  string `json:"namespaceName,omitempty"`
	NotificationID int    `json:"notificationId,omitempty"`
}

type notificatonRepo struct {
	notifications sync.Map
}

func (n *notificatonRepo) setNotificationID(namespace string, notificationID int) {
	n.notifications.Store(namespace, notificationID)
}

func (n *notificatonRepo) getNotificationID(namespace string) (int, bool) {
	if val, ok := n.notifications.Load(namespace); ok {
		if ret, ok := val.(int); ok {
			return ret, true
		}
	}

	return defaultNotificationID, false
}

func (n *notificatonRepo) toString() string {
	var notifications []*notification
	n.notifications.Range(func(key, val interface{}) bool {
		k, _ := key.(string)
		v, _ := val.(int)
		notifications = append(notifications, &notification{
			NamespaceName:  k,
			NotificationID: v,
		})

		return true
	})

	bts, err := json.Marshal(&notifications)
	if err != nil {
		return ""
	}

	return string(bts)
}
