//%: {{ if .Centrifuge }}
package centrifuge

import (
	"gitlab.com/proemergotech/uuid-go"
)

const (
	ChannelNamespaceBoundary = ":"
	ChannelUserBoundary      = "#"
	ChannelGuestBoundary     = "?"
	ChannelSiteGroupBoundary = "|"
)

// GlobalChannel is here so it can be copied and used in caller services, DO NOT DELETE.
func GlobalChannel(namespace string, siteGroupCode string, name string) string {
	return namespace + ChannelNamespaceBoundary + siteGroupCode + ChannelSiteGroupBoundary + name
}

// GuestChannel is here so it can be copied and used in caller services, DO NOT DELETE.
func GuestChannel(namespace string, siteGroupCode string, name string, visitorID string) string {
	return GlobalChannel(namespace, siteGroupCode, name) + ChannelGuestBoundary + visitorID
}

// ProfileChannel is here so it can be copied and used in caller services, DO NOT DELETE.
func ProfileChannel(namespace string, siteGroupCode string, name string, profileUUID uuid.UUID) string {
	return GlobalChannel(namespace, siteGroupCode, name) + ChannelUserBoundary + profileUUID.String()
}

//%: {{ end }}
