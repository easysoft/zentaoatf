package serverConst

type ResultCode int

const (
	ResultSuccess ResultCode = 1
	ResultFail    ResultCode = 0
)

func (c ResultCode) Int() int {
	return int(c)
}

type BuildProgress string

const (
	ProgressCreated    BuildProgress = "created"
	ProgressLaunchVm   BuildProgress = "launch_vm"
	ProgressPending    BuildProgress = "pending"
	ProgressInProgress BuildProgress = "in_progress"
	ProgressTimeout    BuildProgress = "timeout"
	ProgressCompleted  BuildProgress = "completed"
)

type BuildStatus string

const (
	StatusCreated BuildStatus = "created"
	StatusPass    BuildStatus = "pass"
	StatusFail    BuildStatus = "fail"
)

type HostStatus string

const (
	HostActive  HostStatus = "active"
	HostOffline HostStatus = "offline"
)

type VmStatus string

const (
	VmCreated       VmStatus = "created"
	VmLaunch        VmStatus = "launch"
	VmRunning       VmStatus = "running"
	VmActive        VmStatus = "active"
	VmBusy          VmStatus = "busy"
	VmDestroy       VmStatus = "destroy"
	VmFailToCreate  VmStatus = "fail_to_create"
	VmFailToDestroy VmStatus = "fail_to_destroy"
	VmUnknown       VmStatus = "unknown"
)

type DeviceStatus string

const (
	DeviceOff     DeviceStatus = "off"
	DeviceOn      DeviceStatus = "on"
	DeviceActive  DeviceStatus = "active"
	DeviceBusy    DeviceStatus = "busy"
	DeviceUnknown DeviceStatus = "unknown"
)

type ServiceStatus string

const (
	ServiceOff    ServiceStatus = "off"
	ServiceOn     ServiceStatus = "on"
	ServiceActive ServiceStatus = "active"
	ServiceBusy   ServiceStatus = "busy"
)

type Platform string

const (
	Android Platform = "android"
	Ios     Platform = "ios"
	Host    Platform = "host"
	Vm      Platform = "vm"
)

type BuildType string

const (
	AppiumTest   BuildType = "appium_test"
	SeleniumTest BuildType = "selenium_test"
	UnitTest     BuildType = "unit_test"
)

type OsPlatform string

const (
	OsWindows OsPlatform = "windows"
	OsLinux   OsPlatform = "linux"
	OsMac     OsPlatform = "mac"
	OsDriver  OsPlatform = "driver"
)

type OsType string

const (
	Win10  OsType = "win10"
	Win7   OsType = "win7"
	WinXP  OsType = "winxp"
	Ubuntu OsType = "ubuntu"
	CentOS OsType = "centos"
	Mac    OsType = "mac"
	Virtio OsType = "virtio"
)

type OsLang string

const (
	EN_US OsLang = "en_us"
	ZH_CN OsLang = "zh_cn"
	ZH_TW OsLang = "zh_tw"
	ZH_HK OsLang = "zh_hk"
	DE_DE OsLang = "de_de"
	JA_JP OsLang = "ja_jp"
	FR_FR OsLang = "fr_fr"
	RU_RU OsLang = "ru_ru"
)

type BrowserType string

const (
	Chrome  BrowserType = "chrome"
	Firefox BrowserType = "firefox"
	Edge    BrowserType = "edge"
	IE      BrowserType = "ie"
)
