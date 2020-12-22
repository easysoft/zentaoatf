package serverConst

const (
	HeartBeatInterval    = 5
	CheckUpgradeInterval = 5

	AgentRunTime = 30 * 60
	AgentLogDir  = "log-agent"

	QiNiuURL         = "https://dl.cnezsoft.com/ztf/"
	AgentUpgradeURL  = QiNiuURL + "version.txt"
	AgentDownloadURL = QiNiuURL + "%s/%s/ztf.zip"
)
