package commConsts

var (
	LanguageZh = "zh"
	LanguageEn = "en"
	Language   = LanguageZh

	AutoCommitResult bool
	AutoCommitBug    bool
	AutoExtract      bool
	BatchCount       = 1
	Verbose          = false
	IsRelease        bool
	ZtfDir           string
	WorkDir          string
	ExecLogDir       string
	LogDir           string
	WithCache        bool

	RequestType string
	ExecFrom    ExecFromDef
)

var (
	ConfigPath         string
	ServerWorkDir      string
	ServerWorkspaceDir string

	CommandArgs []string

	UnitTestType    string
	UnitBuildTool   BuildTool
	UnitTestTool    TestTool
	AllureReportDir string
	JacocoReport    string
	Options         string
	ProductId       string

	ZenTaoVersion string
	Token         = "Token"
	SessionVar    = "zentaosid"
	SessionId     string
	RequestFix    string

	CurrScriptFile string // scripts/tc-001.py
	CurrResultDate string // 2019-08-15T173802
	CurrCaseId     int    // 2019-08-15T173802

	ScreenWidth    int
	ScreenHeight   int
	CurrBugStepIds string
	Interpreter    string

	// server
	RunMode     string
	IP          string
	MAC         string
	AgentLogDir string
)
