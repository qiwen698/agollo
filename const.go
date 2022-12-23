package agollo

import (
	"time"
)

const (
	defaultConfName  = "agollo.json" // by qiwen698 edited, ori: app.properties
	defaultNamespace = "application"

	longPoolInterval      = time.Second * 2
	longPoolTimeout       = time.Second * 90
	queryTimeout          = time.Second * 2
	defaultNotificationID = -1

	errMissENV = "environment variable not set" // by qiwen698
	//errMissCLI     = "cli arguments not set"        // by qiwen698
	defaultTagName = "config" // by qiwen698
)
