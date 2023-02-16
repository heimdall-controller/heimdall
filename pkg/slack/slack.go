package slack

import (
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	slackclient "github.com/slack-go/slack"
	corev1 "k8s.io/api/core/v1"
	u "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Notification struct {
	Name string
}

// SendEvent SendMessage sends a message using all current senders
func SendEvent(u u.Unstructured, secret *corev1.Secret) {
	token := secret.Data["token"]
	channel := secret.Data["channel"]
	logrus.Infof("Sending event to slack channel %s", channel)
	logrus.Infof("Sending event to slack token %s", token)

	api := slackclient.New(string(token))
	attachment := slackclient.Attachment{
		Fields: []slackclient.AttachmentField{
			{
				Title: "Object Kind: " + u.GetKind(),
			},
			{
				Title: "Object Name: " + u.GetName(),
			},
			{
				Title: "Namespace: " + u.GetNamespace(),
			},
			{
				Title: "Oh no! Please monitor your resource!",
			},
		},
	}

	// Send message to Slack
	channelID, timestamp, err := api.PostMessage(
		string(channel),
		slack.MsgOptionAttachments(attachment),
		slackclient.MsgOptionAsUser(true),
	)
	if err != nil {
		logrus.Errorf("error sending message: %v", err)
	} else {
		logrus.Infof("Message successfully sent to channel %s at %s", channelID, timestamp)
	}
}
