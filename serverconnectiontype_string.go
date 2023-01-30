// Code generated by "stringer -type=ServerConnectionType -trimprefix ServerConnection"; DO NOT EDIT.

package dvls

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ServerConnectionUndefined-0]
	_ = x[ServerConnectionRDPConfigured-1]
	_ = x[ServerConnectionRDPFilename-2]
	_ = x[ServerConnectionCommandLine-3]
	_ = x[ServerConnectionVNC-4]
	_ = x[ServerConnectionWebBrowser-5]
	_ = x[ServerConnectionLogMeIn-6]
	_ = x[ServerConnectionTeamViewer-7]
	_ = x[ServerConnectionPutty-8]
	_ = x[ServerConnectionFtp-9]
	_ = x[ServerConnectionVirtualPC-10]
	_ = x[ServerConnectionRadmin-11]
	_ = x[ServerConnectionDameware-12]
	_ = x[ServerConnectionVMWare-13]
	_ = x[ServerConnectionPCAnywhere-14]
	_ = x[ServerConnectionICA-15]
	_ = x[ServerConnectionXWindow-16]
	_ = x[ServerConnectionHyperV-17]
	_ = x[ServerConnectionAddOn-18]
	_ = x[ServerConnectionRemoteAssistance-19]
	_ = x[ServerConnectionVPN-20]
	_ = x[ServerConnectionVirtualBox-21]
	_ = x[ServerConnectionVMRC-22]
	_ = x[ServerConnectionXenServer-23]
	_ = x[ServerConnectionWindowsVirtualPC-24]
	_ = x[ServerConnectionGroup-25]
	_ = x[ServerConnectionCredential-26]
	_ = x[ServerConnectionHpRgs-27]
	_ = x[ServerConnectionDesktone-28]
	_ = x[ServerConnectionApplicationTool-29]
	_ = x[ServerConnectionSessionTool-30]
	_ = x[ServerConnectionContact-31]
	_ = x[ServerConnectionDataEntry-32]
	_ = x[ServerConnectionDataReport-33]
	_ = x[ServerConnectionAgent-34]
	_ = x[ServerConnectionComputer-35]
	_ = x[ServerConnectionDropBox-36]
	_ = x[ServerConnectionS3-37]
	_ = x[ServerConnectionAzureStorage-38]
	_ = x[ServerConnectionCitrixWeb-39]
	_ = x[ServerConnectionPowerShell-40]
	_ = x[ServerConnectionHostSessionTool-41]
	_ = x[ServerConnectionShortcut-42]
	_ = x[ServerConnectionIntelAMT-43]
	_ = x[ServerConnectionAzure-44]
	_ = x[ServerConnectionDocument-45]
	_ = x[ServerConnectionVMWareConsole-46]
	_ = x[ServerConnectionInventoryReport-47]
	_ = x[ServerConnectionSkyDrive-48]
	_ = x[ServerConnectionScreenConnect-49]
	_ = x[ServerConnectionAzureTableStorage-50]
	_ = x[ServerConnectionAzureQueueStorage-51]
	_ = x[ServerConnectionTemplateGroup-52]
	_ = x[ServerConnectionHost-53]
	_ = x[ServerConnectionDatabase-54]
	_ = x[ServerConnectionCustomer-55]
	_ = x[ServerConnectionADConsole-56]
	_ = x[ServerConnectionAws-57]
	_ = x[ServerConnectionSNMPReport-58]
	_ = x[ServerConnectionSync-59]
	_ = x[ServerConnectionGateway-60]
	_ = x[ServerConnectionPlayList-61]
	_ = x[ServerConnectionTerminalConsole-62]
	_ = x[ServerConnectionPSExec-63]
	_ = x[ServerConnectionAppleRemoteDesktop-64]
	_ = x[ServerConnectionSpiceworks-65]
	_ = x[ServerConnectionDeskRoll-66]
	_ = x[ServerConnectionSecureCRT-67]
	_ = x[ServerConnectionIterm-68]
	_ = x[ServerConnectionSheet-69]
	_ = x[ServerConnectionSplunk-70]
	_ = x[ServerConnectionPortForward-71]
	_ = x[ServerConnectionTeamViewerConsole-72]
	_ = x[ServerConnectionScreenHero-73]
	_ = x[ServerConnectionTelnet-74]
	_ = x[ServerConnectionSerial-75]
	_ = x[ServerConnectionSSHTunnel-76]
	_ = x[ServerConnectionSSHShell-77]
	_ = x[ServerConnectionResetPassword-78]
	_ = x[ServerConnectionWayk-79]
	_ = x[ServerConnectionControlUp-80]
	_ = x[ServerConnectionDataSource-81]
	_ = x[ServerConnectionChromeRemoteDesktop-82]
	_ = x[ServerConnectionRDCommander-83]
	_ = x[ServerConnectionIDrac-84]
	_ = x[ServerConnectionIlo-85]
	_ = x[ServerConnectionWebDav-86]
	_ = x[ServerConnectionBeyondTrustPasswordSafeConsole-87]
	_ = x[ServerConnectionDevolutionsProxy-88]
	_ = x[ServerConnectionFtpNative-89]
	_ = x[ServerConnectionPowerShellRemoteConsole-90]
	_ = x[ServerConnectionProxyTunnel-91]
	_ = x[ServerConnectionRoot-92]
	_ = x[ServerConnectionBeyondTrustPasswordSafe-93]
	_ = x[ServerConnectionFileExplorer-94]
	_ = x[ServerConnectionScp-95]
	_ = x[ServerConnectionSftp-96]
	_ = x[ServerConnectionAzureBlobStorage-97]
	_ = x[ServerConnectionTFtp-98]
	_ = x[ServerConnectionGoToAssist-99]
	_ = x[ServerConnectionIPTable-100]
	_ = x[ServerConnectionHub-101]
	_ = x[ServerConnectionGoogleDrive-102]
	_ = x[ServerConnectionGoogleCloud-103]
	_ = x[ServerConnectionNoVNC-104]
	_ = x[ServerConnectionSplashtop-105]
	_ = x[ServerConnectionJumpDesktop-106]
	_ = x[ServerConnectionBoxNet-107]
	_ = x[ServerConnectionMSPAnywhere-108]
	_ = x[ServerConnectionRepository-109]
	_ = x[ServerConnectionCyberArkPSM-110]
	_ = x[ServerConnectionCloudBerryRemoteAssistant-111]
	_ = x[ServerConnectionITGlue-112]
	_ = x[ServerConnectionSmartFolder-113]
	_ = x[ServerConnectionCyberArkJump-114]
	_ = x[ServerConnectionWindowsAdminCenter-115]
	_ = x[ServerConnectionDevolutionsGateway-116]
	_ = x[ServerConnectionWaykDenConsole-117]
	_ = x[ServerConnectionRDGatewayConsole-118]
	_ = x[ServerConnectionCyberArkDashboard-119]
	_ = x[ServerConnectionDVLSPamDashboard-120]
	_ = x[ServerConnectionSMB-121]
	_ = x[ServerConnectionAppleRemoteManagement-122]
	_ = x[ServerConnectionRustDesk-123]
	_ = x[ServerConnectionPAM-124]
	_ = x[ServerConnectionITManager-125]
	_ = x[ServerConnectionCustomImage-126]
}

const _ServerConnectionType_name = "UndefinedRDPConfiguredRDPFilenameCommandLineVNCWebBrowserLogMeInTeamViewerPuttyFtpVirtualPCRadminDamewareVMWarePCAnywhereICAXWindowHyperVAddOnRemoteAssistanceVPNVirtualBoxVMRCXenServerWindowsVirtualPCGroupCredentialHpRgsDesktoneApplicationToolSessionToolContactDataEntryDataReportAgentComputerDropBoxS3AzureStorageCitrixWebPowerShellHostSessionToolShortcutIntelAMTAzureDocumentVMWareConsoleInventoryReportSkyDriveScreenConnectAzureTableStorageAzureQueueStorageTemplateGroupHostDatabaseCustomerADConsoleAwsSNMPReportSyncGatewayPlayListTerminalConsolePSExecAppleRemoteDesktopSpiceworksDeskRollSecureCRTItermSheetSplunkPortForwardTeamViewerConsoleScreenHeroTelnetSerialSSHTunnelSSHShellResetPasswordWaykControlUpDataSourceChromeRemoteDesktopRDCommanderIDracIloWebDavBeyondTrustPasswordSafeConsoleDevolutionsProxyFtpNativePowerShellRemoteConsoleProxyTunnelRootBeyondTrustPasswordSafeFileExplorerScpSftpAzureBlobStorageTFtpGoToAssistIPTableHubGoogleDriveGoogleCloudNoVNCSplashtopJumpDesktopBoxNetMSPAnywhereRepositoryCyberArkPSMCloudBerryRemoteAssistantITGlueSmartFolderCyberArkJumpWindowsAdminCenterDevolutionsGatewayWaykDenConsoleRDGatewayConsoleCyberArkDashboardDVLSPamDashboardSMBAppleRemoteManagementRustDeskPAMITManagerCustomImage"

var _ServerConnectionType_index = [...]uint16{0, 9, 22, 33, 44, 47, 57, 64, 74, 79, 82, 91, 97, 105, 111, 121, 124, 131, 137, 142, 158, 161, 171, 175, 184, 200, 205, 215, 220, 228, 243, 254, 261, 270, 280, 285, 293, 300, 302, 314, 323, 333, 348, 356, 364, 369, 377, 390, 405, 413, 426, 443, 460, 473, 477, 485, 493, 502, 505, 515, 519, 526, 534, 549, 555, 573, 583, 591, 600, 605, 610, 616, 627, 644, 654, 660, 666, 675, 683, 696, 700, 709, 719, 738, 749, 754, 757, 763, 793, 809, 818, 841, 852, 856, 879, 891, 894, 898, 914, 918, 928, 935, 938, 949, 960, 965, 974, 985, 991, 1002, 1012, 1023, 1048, 1054, 1065, 1077, 1095, 1113, 1127, 1143, 1160, 1176, 1179, 1200, 1208, 1211, 1220, 1231}

func (i ServerConnectionType) String() string {
	if i < 0 || i >= ServerConnectionType(len(_ServerConnectionType_index)-1) {
		return "ServerConnectionType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ServerConnectionType_name[_ServerConnectionType_index[i]:_ServerConnectionType_index[i+1]]
}
