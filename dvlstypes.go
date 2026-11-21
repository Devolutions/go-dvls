package dvls

// ServerConnectionSubType is the subtype of an entry; for a website entry, the browser that opens it.
type ServerConnectionSubType string

const (
	ServerConnectionSubTypeDefault          ServerConnectionSubType = "Default"
	ServerConnectionSubTypeCertificate      ServerConnectionSubType = "Certificate"
	ServerConnectionSubTypeMicrosoftEdge    ServerConnectionSubType = "Edge"
	ServerConnectionSubTypeFirefox          ServerConnectionSubType = "FireFox"
	ServerConnectionSubTypeGoogleChrome     ServerConnectionSubType = "GoogleChrome"
	ServerConnectionSubTypeInternetExplorer ServerConnectionSubType = "IE"
	ServerConnectionSubTypeOpera            ServerConnectionSubType = "Opera"
	ServerConnectionSubTypeAppleSafari      ServerConnectionSubType = "Safari"
)
