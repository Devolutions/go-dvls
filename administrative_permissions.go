package dvls

// AdministrativePermission identifies a permission that can be granted through
// a DVLS administrative role.
// Values are grouped per feature area (Users=100, AdministrativeRoles=200,
// Repositories=300, EntryTemplates=600, ServerSettings=700, Logs=800, Gateways=900,
// GlobalEntries=1200, Licenses=1300, GlobalImages=1400, SecurityPolicies=1500,
// Support=1600, UserGroups=1700, ApplicationIdentities=1800, Backups=2000,
// LogRetentionPolicies=2100, UserBehaviorAnalytics=2200, SecurityDashboard=2300,
// PasswordPolicies=2400, ScheduledReports=2500, Notifications=2600, PamSettings=2700,
// Webhooks=2800, Monitoring=2900, PamInfrastructureVault=3000, CustomWidgets=3100,
// DataSourceSettings=3200, DataSourcePermissions=3300, PamProviders=3400, Tools=3500,
// BusinessUnits=3600, SecurityProvider=3700) and are sparse within each group.
// Unknown values returned by newer server versions still round-trip unchanged.
//
//go:generate stringer -type=AdministrativePermission -trimprefix AdministrativePermission
type AdministrativePermission int

const (
	AdministrativePermissionUsersView                       AdministrativePermission = 100
	AdministrativePermissionUsersAdd                        AdministrativePermission = 101
	AdministrativePermissionUsersEdit                       AdministrativePermission = 102
	AdministrativePermissionUsersDelete                     AdministrativePermission = 103
	AdministrativePermissionUsersDeletedView                AdministrativePermission = 104
	AdministrativePermissionUsersRestore                    AdministrativePermission = 105
	AdministrativePermissionUsersDeletePermanent            AdministrativePermission = 106
	AdministrativePermissionUsersResetPassword              AdministrativePermission = 107
	AdministrativePermissionUsersMfaRequestView             AdministrativePermission = 110
	AdministrativePermissionUsersLockedView                 AdministrativePermission = 111
	AdministrativePermissionUsersMfaRequestCancel           AdministrativePermission = 112
	AdministrativePermissionUsersMfaDelete                  AdministrativePermission = 113
	AdministrativePermissionUsersAppKeysView                AdministrativePermission = 114
	AdministrativePermissionUsersAppKeysDelete              AdministrativePermission = 115
	AdministrativePermissionUsersDisconnect                 AdministrativePermission = 116
	AdministrativePermissionUsersUnlock                     AdministrativePermission = 119
	AdministrativePermissionUsersUpdateFromIdentityProvider AdministrativePermission = 120
	AdministrativePermissionUsersCleanup                    AdministrativePermission = 121
	AdministrativePermissionUsersImport                     AdministrativePermission = 124
	AdministrativePermissionUsersStatusEdit                 AdministrativePermission = 125
	AdministrativePermissionUsersViewActivity               AdministrativePermission = 126
	AdministrativePermissionUsersResetMfa                   AdministrativePermission = 127
	AdministrativePermissionUsersRemoveMfaRequest           AdministrativePermission = 128
	AdministrativePermissionUsersReport                     AdministrativePermission = 129
	AdministrativePermissionUsersConnectedView              AdministrativePermission = 130
	AdministrativePermissionUsersSQLoginFix                 AdministrativePermission = 131
	AdministrativePermissionUsersSQLPermissionsView         AdministrativePermission = 132
	AdministrativePermissionUsersInfoHistoryDelete          AdministrativePermission = 133

	AdministrativePermissionAdministrativeRolesView              AdministrativePermission = 200
	AdministrativePermissionAdministrativeRolesAdd               AdministrativePermission = 201
	AdministrativePermissionAdministrativeRolesEdit              AdministrativePermission = 202
	AdministrativePermissionAdministrativeRolesDelete            AdministrativePermission = 203
	AdministrativePermissionAdministrativeRolesAssignmentsView   AdministrativePermission = 204
	AdministrativePermissionAdministrativeRolesAssignmentsManage AdministrativePermission = 205

	AdministrativePermissionRepositoriesView                                AdministrativePermission = 300
	AdministrativePermissionRepositoriesAdd                                 AdministrativePermission = 301
	AdministrativePermissionRepositoriesEdit                                AdministrativePermission = 302
	AdministrativePermissionRepositoriesDelete                              AdministrativePermission = 303
	AdministrativePermissionRepositoriesAssignmentsManage                   AdministrativePermission = 304
	AdministrativePermissionRepositoriesContentExport                       AdministrativePermission = 305
	AdministrativePermissionRepositoriesContentViewDeleted                  AdministrativePermission = 306
	AdministrativePermissionRepositoriesContentRestore                      AdministrativePermission = 307
	AdministrativePermissionRepositoriesContentDeletePermanent              AdministrativePermission = 308
	AdministrativePermissionRepositoriesContentOpenedView                   AdministrativePermission = 309
	AdministrativePermissionRepositoriesContentStatusView                   AdministrativePermission = 310
	AdministrativePermissionRepositoriesContentPermissionView               AdministrativePermission = 312
	AdministrativePermissionRepositoriesContentExpiredView                  AdministrativePermission = 313
	AdministrativePermissionRepositoriesContentLogsView                     AdministrativePermission = 314
	AdministrativePermissionRepositoriesContentLastUsageView                AdministrativePermission = 315
	AdministrativePermissionRepositoriesContentPasswordAnalyzerView         AdministrativePermission = 316
	AdministrativePermissionRepositoriesFlagConnectionAsClosed              AdministrativePermission = 317
	AdministrativePermissionRepositoriesContentDocumentExport               AdministrativePermission = 318
	AdministrativePermissionRepositoriesPermissionView                      AdministrativePermission = 319
	AdministrativePermissionRepositoriesContentAccessRequestView            AdministrativePermission = 320
	AdministrativePermissionRepositoriesContentAccessRequestApprove         AdministrativePermission = 321
	AdministrativePermissionRepositoriesContentToDoManage                   AdministrativePermission = 322
	AdministrativePermissionRepositoriesContentImport                       AdministrativePermission = 323
	AdministrativePermissionRepositoriesContentCheckinOverride              AdministrativePermission = 324
	AdministrativePermissionRepositoriesContentLogsEdit                     AdministrativePermission = 325
	AdministrativePermissionRepositoriesManageGateway                       AdministrativePermission = 326
	AdministrativePermissionRepositoriesContentFullAccess                   AdministrativePermission = 327
	AdministrativePermissionRepositoriesAccessRequestView                   AdministrativePermission = 330
	AdministrativePermissionRepositoriesAccessRequestApprove                AdministrativePermission = 331
	AdministrativePermissionRepositoriesContentView                         AdministrativePermission = 332
	AdministrativePermissionRepositoriesPamForceCheckin                     AdministrativePermission = 336
	AdministrativePermissionRepositoriesPamCheckoutView                     AdministrativePermission = 337
	AdministrativePermissionRepositoriesPamCheckoutApprove                  AdministrativePermission = 338
	AdministrativePermissionRepositoriesContentStatusEdit                   AdministrativePermission = 339
	AdministrativePermissionRepositoriesContentStatusBypass                 AdministrativePermission = 340
	AdministrativePermissionRepositoriesContentTimeBasedAccessBypass        AdministrativePermission = 341
	AdministrativePermissionRepositoriesAssignmentsView                     AdministrativePermission = 342
	AdministrativePermissionRepositoriesListView                            AdministrativePermission = 343
	AdministrativePermissionRepositoriesContentCheckoutView                 AdministrativePermission = 344
	AdministrativePermissionRepositoriesPamAdd                              AdministrativePermission = 345
	AdministrativePermissionRepositoriesPamListView                         AdministrativePermission = 346
	AdministrativePermissionRepositoriesContentTransferIn                   AdministrativePermission = 347
	AdministrativePermissionRepositoriesContentTransferOut                  AdministrativePermission = 348
	AdministrativePermissionRepositoriesHistoryView                         AdministrativePermission = 349
	AdministrativePermissionRepositoriesContentPasswordAnalyzerPasswordView AdministrativePermission = 350
	AdministrativePermissionRepositoriesContentPasswordRotationView         AdministrativePermission = 351
	AdministrativePermissionRepositoriesContentPendingApprovalStatusManage  AdministrativePermission = 352
	AdministrativePermissionRepositoriesContentDocumentationExport          AdministrativePermission = 353
	AdministrativePermissionRepositoriesContentDocumentationHistoryDelete   AdministrativePermission = 354
	AdministrativePermissionRepositoriesContentDocumentationHistoryRestore  AdministrativePermission = 355
	AdministrativePermissionRepositoriesContentLogsDelete                   AdministrativePermission = 356
	AdministrativePermissionRepositoriesContentQRCodeView                   AdministrativePermission = 357
	AdministrativePermissionRepositoriesHistoryDelete                       AdministrativePermission = 358
	AdministrativePermissionRepositoriesContentAttachmentsHistoryDelete     AdministrativePermission = 359
	AdministrativePermissionRepositoriesContentHistoryDelete                AdministrativePermission = 360
	AdministrativePermissionRepositoriesPamCheckoutRequest                  AdministrativePermission = 361

	AdministrativePermissionEntryTemplatesView       AdministrativePermission = 600
	AdministrativePermissionEntryTemplatesAdd        AdministrativePermission = 601
	AdministrativePermissionEntryTemplatesEdit       AdministrativePermission = 602
	AdministrativePermissionEntryTemplatesDelete     AdministrativePermission = 603
	AdministrativePermissionEntryTemplatesViewUsage  AdministrativePermission = 604
	AdministrativePermissionEntryTemplatesSetDefault AdministrativePermission = 605
	AdministrativePermissionEntryTemplatesExport     AdministrativePermission = 606

	AdministrativePermissionServerSettingsGeneralView        AdministrativePermission = 700
	AdministrativePermissionServerSettingsGeneralEdit        AdministrativePermission = 701
	AdministrativePermissionServerSettingsAuthenticationView AdministrativePermission = 702
	AdministrativePermissionServerSettingsAuthenticationEdit AdministrativePermission = 703
	AdministrativePermissionServerSettingsEmailView          AdministrativePermission = 704
	AdministrativePermissionServerSettingsEmailEdit          AdministrativePermission = 705
	AdministrativePermissionServerSettingsLogsView           AdministrativePermission = 708
	AdministrativePermissionServerSettingsLogsEdit           AdministrativePermission = 709
	AdministrativePermissionServerSettingsFeaturesView       AdministrativePermission = 710
	AdministrativePermissionServerSettingsFeaturesEdit       AdministrativePermission = 711
	AdministrativePermissionServerSettingsApprovalsView      AdministrativePermission = 712
	AdministrativePermissionServerSettingsApprovalsEdit      AdministrativePermission = 713
	AdministrativePermissionServerSettingsAdvancedView       AdministrativePermission = 714
	AdministrativePermissionServerSettingsAdvancedEdit       AdministrativePermission = 715
	AdministrativePermissionServerSettingsMfaView            AdministrativePermission = 716
	AdministrativePermissionServerSettingsMfaEdit            AdministrativePermission = 717
	AdministrativePermissionServerSettingsSecurityView       AdministrativePermission = 718
	AdministrativePermissionServerSettingsSecurityEdit       AdministrativePermission = 719
	AdministrativePermissionServerSettingsGeoIPView          AdministrativePermission = 720
	AdministrativePermissionServerSettingsGeoIPEdit          AdministrativePermission = 721
	AdministrativePermissionServerSettingsAIAssistantView    AdministrativePermission = 722
	AdministrativePermissionServerSettingsAIAssistantEdit    AdministrativePermission = 723

	AdministrativePermissionLogsServerView            AdministrativePermission = 803
	AdministrativePermissionLogsAdminView             AdministrativePermission = 804
	AdministrativePermissionLogsAdminPermissionView   AdministrativePermission = 806
	AdministrativePermissionLogsLoginAttemptsView     AdministrativePermission = 807
	AdministrativePermissionLogsLoginHistoryView      AdministrativePermission = 808
	AdministrativePermissionLogsLoginLastView         AdministrativePermission = 809
	AdministrativePermissionLogsSystemPermissionsView AdministrativePermission = 810
	AdministrativePermissionLogsPamServerView         AdministrativePermission = 812
	AdministrativePermissionLogsAdminDelete           AdministrativePermission = 813

	AdministrativePermissionGatewaysFullAccess        AdministrativePermission = 900
	AdministrativePermissionGatewaysUserAccessView    AdministrativePermission = 901
	AdministrativePermissionGatewaysPermissionsView   AdministrativePermission = 902
	AdministrativePermissionGatewaysSessionsTerminate AdministrativePermission = 903

	AdministrativePermissionGlobalEntriesAdd                     AdministrativePermission = 1200
	AdministrativePermissionGlobalEntriesEdit                    AdministrativePermission = 1201
	AdministrativePermissionGlobalEntriesDelete                  AdministrativePermission = 1202
	AdministrativePermissionGlobalEntriesSessionView             AdministrativePermission = 1203
	AdministrativePermissionGlobalEntriesExecute                 AdministrativePermission = 1204
	AdministrativePermissionGlobalEntriesViewPassword            AdministrativePermission = 1205
	AdministrativePermissionGlobalVaultView                      AdministrativePermission = 1206
	AdministrativePermissionGlobalEntriesToolsMacroScriptExecute AdministrativePermission = 1207

	AdministrativePermissionLicensesView           AdministrativePermission = 1300
	AdministrativePermissionLicensesAdd            AdministrativePermission = 1301
	AdministrativePermissionLicensesEdit           AdministrativePermission = 1302
	AdministrativePermissionLicensesDelete         AdministrativePermission = 1303
	AdministrativePermissionLicensesAssign         AdministrativePermission = 1304
	AdministrativePermissionLicensesUserReportView AdministrativePermission = 1305
	AdministrativePermissionLicensesReportView     AdministrativePermission = 1306

	AdministrativePermissionGlobalImagesView    AdministrativePermission = 1400
	AdministrativePermissionGlobalImagesAdd     AdministrativePermission = 1401
	AdministrativePermissionGlobalImagesEdit    AdministrativePermission = 1402
	AdministrativePermissionGlobalImagesDelete  AdministrativePermission = 1403
	AdministrativePermissionGlobalImagesCleanup AdministrativePermission = 1404

	AdministrativePermissionSecurityPoliciesView   AdministrativePermission = 1500
	AdministrativePermissionSecurityPoliciesAdd    AdministrativePermission = 1501
	AdministrativePermissionSecurityPoliciesEdit   AdministrativePermission = 1502
	AdministrativePermissionSecurityPoliciesDelete AdministrativePermission = 1503

	AdministrativePermissionSupportDownloadPackage     AdministrativePermission = 1600
	AdministrativePermissionSupportSendTicket          AdministrativePermission = 1601
	AdministrativePermissionSupportDiagnostic          AdministrativePermission = 1602
	AdministrativePermissionSupportInvalidateCache     AdministrativePermission = 1603
	AdministrativePermissionSupportPackDatabase        AdministrativePermission = 1604
	AdministrativePermissionSupportDataSourceMigration AdministrativePermission = 1605

	AdministrativePermissionUserGroupsView                       AdministrativePermission = 1700
	AdministrativePermissionUserGroupsAdd                        AdministrativePermission = 1701
	AdministrativePermissionUserGroupsEdit                       AdministrativePermission = 1702
	AdministrativePermissionUserGroupsDelete                     AdministrativePermission = 1703
	AdministrativePermissionUserGroupsViewDeleted                AdministrativePermission = 1704
	AdministrativePermissionUserGroupsRestore                    AdministrativePermission = 1705
	AdministrativePermissionUserGroupsDeletePermanent            AdministrativePermission = 1706
	AdministrativePermissionUserGroupsEditMembership             AdministrativePermission = 1707
	AdministrativePermissionUserGroupsUpdateFromIdentityProvider AdministrativePermission = 1708
	AdministrativePermissionUserGroupsImport                     AdministrativePermission = 1709

	AdministrativePermissionApplicationIdentitiesView          AdministrativePermission = 1800
	AdministrativePermissionApplicationIdentitiesAdd           AdministrativePermission = 1801
	AdministrativePermissionApplicationIdentitiesEdit          AdministrativePermission = 1802
	AdministrativePermissionApplicationIdentitiesDelete        AdministrativePermission = 1803
	AdministrativePermissionApplicationIdentitiesResetPassword AdministrativePermission = 1804
	AdministrativePermissionApplicationIdentitiesViewActivity  AdministrativePermission = 1806

	AdministrativePermissionBackupsView     AdministrativePermission = 2000
	AdministrativePermissionBackupsEdit     AdministrativePermission = 2001
	AdministrativePermissionBackupsRun      AdministrativePermission = 2002
	AdministrativePermissionBackupsLogsView AdministrativePermission = 2003

	AdministrativePermissionLogRetentionPoliciesView AdministrativePermission = 2100
	AdministrativePermissionLogRetentionPoliciesEdit AdministrativePermission = 2101
	AdministrativePermissionLogRetentionPoliciesRun  AdministrativePermission = 2102

	AdministrativePermissionUserBehaviorAnalyticsConfigurationView  AdministrativePermission = 2200
	AdministrativePermissionUserBehaviorAnalyticsConfigurationEdit  AdministrativePermission = 2201
	AdministrativePermissionUserBehaviorAnalyticsResultsView        AdministrativePermission = 2202
	AdministrativePermissionUserBehaviorAnalyticsResultsExport      AdministrativePermission = 2203
	AdministrativePermissionUserBehaviorAnalyticsResultsAcknowledge AdministrativePermission = 2204

	AdministrativePermissionSecurityDashboardView                 AdministrativePermission = 2300
	AdministrativePermissionSecurityDashboardIgnoreItems          AdministrativePermission = 2301
	AdministrativePermissionSecurityDashboardServicesStatusDelete AdministrativePermission = 2303

	AdministrativePermissionPasswordPoliciesView                           AdministrativePermission = 2400
	AdministrativePermissionPasswordPoliciesAdd                            AdministrativePermission = 2401
	AdministrativePermissionPasswordPoliciesEdit                           AdministrativePermission = 2402
	AdministrativePermissionPasswordPoliciesDelete                         AdministrativePermission = 2403
	AdministrativePermissionPasswordPoliciesExport                         AdministrativePermission = 2404
	AdministrativePermissionPasswordPoliciesPassphraseDictionariesView     AdministrativePermission = 2405
	AdministrativePermissionPasswordPoliciesPassphraseDictionariesAdd      AdministrativePermission = 2406
	AdministrativePermissionPasswordPoliciesPassphraseDictionariesEdit     AdministrativePermission = 2407
	AdministrativePermissionPasswordPoliciesPassphraseDictionariesDelete   AdministrativePermission = 2408
	AdministrativePermissionPasswordPoliciesPassphraseDictionariesExport   AdministrativePermission = 2409
	AdministrativePermissionPasswordPoliciesPassphraseDictionariesDownload AdministrativePermission = 2410
	AdministrativePermissionPasswordPoliciesForcedDefaultPolicyBypass      AdministrativePermission = 2411

	AdministrativePermissionScheduledReportsView   AdministrativePermission = 2500
	AdministrativePermissionScheduledReportsAdd    AdministrativePermission = 2501
	AdministrativePermissionScheduledReportsEdit   AdministrativePermission = 2502
	AdministrativePermissionScheduledReportsDelete AdministrativePermission = 2503
	AdministrativePermissionScheduledReportsSend   AdministrativePermission = 2504

	AdministrativePermissionNotificationsSubscriberView        AdministrativePermission = 2600
	AdministrativePermissionNotificationsSubscriberAdd         AdministrativePermission = 2601
	AdministrativePermissionNotificationsSubscriberEdit        AdministrativePermission = 2602
	AdministrativePermissionNotificationsSubscriberDelete      AdministrativePermission = 2603
	AdministrativePermissionNotificationsSubscriberGroupView   AdministrativePermission = 2604
	AdministrativePermissionNotificationsSubscriberGroupAdd    AdministrativePermission = 2605
	AdministrativePermissionNotificationsSubscriberGroupEdit   AdministrativePermission = 2606
	AdministrativePermissionNotificationsSubscriberGroupDelete AdministrativePermission = 2607
	AdministrativePermissionNotificationsSubscriptionsView     AdministrativePermission = 2608
	AdministrativePermissionNotificationsSubscriptionsAdd      AdministrativePermission = 2609
	AdministrativePermissionNotificationsSubscriptionsEdit     AdministrativePermission = 2610
	AdministrativePermissionNotificationsSubscriptionsDelete   AdministrativePermission = 2611
	AdministrativePermissionNotificationsUserView              AdministrativePermission = 2612
	AdministrativePermissionNotificationsUserDelete            AdministrativePermission = 2613
	AdministrativePermissionNotificationsWorkspaceView         AdministrativePermission = 2614
	AdministrativePermissionNotificationsWorkspaceEdit         AdministrativePermission = 2615

	AdministrativePermissionPamSettingsPropagationTemplateView        AdministrativePermission = 2700
	AdministrativePermissionPamSettingsPropagationTemplateAdd         AdministrativePermission = 2701
	AdministrativePermissionPamSettingsPropagationTemplateEdit        AdministrativePermission = 2702
	AdministrativePermissionPamSettingsPropagationTemplateDelete      AdministrativePermission = 2703
	AdministrativePermissionPamSettingsPropagationTemplateExport      AdministrativePermission = 2704
	AdministrativePermissionPamSettingsPropagationConfigurationView   AdministrativePermission = 2705
	AdministrativePermissionPamSettingsPropagationConfigurationAdd    AdministrativePermission = 2706
	AdministrativePermissionPamSettingsPropagationConfigurationEdit   AdministrativePermission = 2707
	AdministrativePermissionPamSettingsPropagationConfigurationDelete AdministrativePermission = 2708
	AdministrativePermissionPamSettingsCheckoutPoliciesView           AdministrativePermission = 2711
	AdministrativePermissionPamSettingsCheckoutPoliciesAdd            AdministrativePermission = 2712
	AdministrativePermissionPamSettingsCheckoutPoliciesEdit           AdministrativePermission = 2713
	AdministrativePermissionPamSettingsCheckoutPoliciesDelete         AdministrativePermission = 2714
	AdministrativePermissionPamSettingsLifecyclePoliciesView          AdministrativePermission = 2718
	AdministrativePermissionPamSettingsLifecyclePoliciesAdd           AdministrativePermission = 2719
	AdministrativePermissionPamSettingsLifecyclePoliciesEdit          AdministrativePermission = 2720
	AdministrativePermissionPamSettingsLifecyclePoliciesDelete        AdministrativePermission = 2721
	AdministrativePermissionPamSettingsOTPView                        AdministrativePermission = 2723
	AdministrativePermissionPamSettingsOTPAdd                         AdministrativePermission = 2724
	AdministrativePermissionPamSettingsOTPEdit                        AdministrativePermission = 2725
	AdministrativePermissionPamSettingsOTPDelete                      AdministrativePermission = 2726
	AdministrativePermissionPamSettingsUsagePoliciesView              AdministrativePermission = 2728
	AdministrativePermissionPamSettingsUsagePoliciesEdit              AdministrativePermission = 2729
	AdministrativePermissionPamSettingsCustomProviderTemplatesView    AdministrativePermission = 2730
	AdministrativePermissionPamSettingsCustomProviderTemplatesAdd     AdministrativePermission = 2731
	AdministrativePermissionPamSettingsCustomProviderTemplatesEdit    AdministrativePermission = 2732
	AdministrativePermissionPamSettingsCustomProviderTemplatesDelete  AdministrativePermission = 2733
	AdministrativePermissionPamSettingsCustomProviderTemplatesExport  AdministrativePermission = 2734
	AdministrativePermissionPamSettingsGlobalView                     AdministrativePermission = 2735
	AdministrativePermissionPamSettingsGlobalEdit                     AdministrativePermission = 2736
	AdministrativePermissionPamSettingsJitTemplatesView               AdministrativePermission = 2737
	AdministrativePermissionPamSettingsJitTemplatesAdd                AdministrativePermission = 2738
	AdministrativePermissionPamSettingsJitTemplatesEdit               AdministrativePermission = 2739
	AdministrativePermissionPamSettingsJitTemplatesDelete             AdministrativePermission = 2740
	AdministrativePermissionPamSettingsPrivilegedAccessRiskView       AdministrativePermission = 2741
	AdministrativePermissionPamSettingsPrivilegedAccessRiskAdd        AdministrativePermission = 2742
	AdministrativePermissionPamSettingsPrivilegedAccessRiskEdit       AdministrativePermission = 2743
	AdministrativePermissionPamSettingsPrivilegedAccessRiskDelete     AdministrativePermission = 2744

	AdministrativePermissionWebhooksView   AdministrativePermission = 2800
	AdministrativePermissionWebhooksAdd    AdministrativePermission = 2801
	AdministrativePermissionWebhooksEdit   AdministrativePermission = 2802
	AdministrativePermissionWebhooksDelete AdministrativePermission = 2803

	AdministrativePermissionMonitoringPrivilegedSessionsView   AdministrativePermission = 2900
	AdministrativePermissionMonitoringSessionRecordingView     AdministrativePermission = 2901
	AdministrativePermissionMonitoringSessionRecordingDelete   AdministrativePermission = 2902
	AdministrativePermissionMonitoringSessionRecordingDownload AdministrativePermission = 2903
	AdministrativePermissionMonitoringSessionRecordingPreserve AdministrativePermission = 2904
	AdministrativePermissionMonitoringPrivilegedSessionsClose  AdministrativePermission = 2905
	AdministrativePermissionMonitoringSystemDashboardView      AdministrativePermission = 2906
	AdministrativePermissionMonitoringSystemWarningsView       AdministrativePermission = 2907

	AdministrativePermissionPamInfrastructureVaultView   AdministrativePermission = 3000
	AdministrativePermissionPamInfrastructureVaultManage AdministrativePermission = 3001

	AdministrativePermissionCustomWidgetsView         AdministrativePermission = 3100
	AdministrativePermissionCustomWidgetsAdd          AdministrativePermission = 3101
	AdministrativePermissionCustomWidgetsEdit         AdministrativePermission = 3102
	AdministrativePermissionCustomWidgetsDelete       AdministrativePermission = 3103
	AdministrativePermissionCustomWidgetLayoutsView   AdministrativePermission = 3104
	AdministrativePermissionCustomWidgetLayoutsAdd    AdministrativePermission = 3105
	AdministrativePermissionCustomWidgetLayoutsEdit   AdministrativePermission = 3106
	AdministrativePermissionCustomWidgetLayoutsDelete AdministrativePermission = 3107

	AdministrativePermissionDataSourceSettingsView          AdministrativePermission = 3217
	AdministrativePermissionDataSourceSettingsEdit          AdministrativePermission = 3218
	AdministrativePermissionDataSourceSettingsHistoryView   AdministrativePermission = 3219
	AdministrativePermissionDataSourceSettingsHistoryDelete AdministrativePermission = 3220

	AdministrativePermissionDataSourcePermissionsView AdministrativePermission = 3320
	AdministrativePermissionDataSourcePermissionsEdit AdministrativePermission = 3321

	AdministrativePermissionPamProvidersView                   AdministrativePermission = 3400
	AdministrativePermissionPamProvidersAdd                    AdministrativePermission = 3401
	AdministrativePermissionPamProvidersEdit                   AdministrativePermission = 3402
	AdministrativePermissionPamProvidersDelete                 AdministrativePermission = 3403
	AdministrativePermissionPamProvidersScanView               AdministrativePermission = 3404
	AdministrativePermissionPamProvidersScanAdd                AdministrativePermission = 3405
	AdministrativePermissionPamProvidersScanEdit               AdministrativePermission = 3406
	AdministrativePermissionPamProvidersScanDelete             AdministrativePermission = 3407
	AdministrativePermissionPamProvidersScanExecute            AdministrativePermission = 3408
	AdministrativePermissionPamProvidersScanResultsView        AdministrativePermission = 3409
	AdministrativePermissionPamProvidersScanResultsImport      AdministrativePermission = 3410
	AdministrativePermissionPamProvidersScanResultsAcknowledge AdministrativePermission = 3411
	AdministrativePermissionPamProvidersListGroups             AdministrativePermission = 3412
	AdministrativePermissionPamProvidersRemoveGroupAssignment  AdministrativePermission = 3413
	AdministrativePermissionPamProvidersAssignmentsManage      AdministrativePermission = 3414
	AdministrativePermissionPamProvidersAssignmentsView        AdministrativePermission = 3415
	AdministrativePermissionPamProvidersCreateAccount          AdministrativePermission = 3416

	AdministrativePermissionToolsMacroScriptExecute       AdministrativePermission = 3500
	AdministrativePermissionToolsRemoteExecute            AdministrativePermission = 3501
	AdministrativePermissionToolsMacroScriptEntryExecute  AdministrativePermission = 3502
	AdministrativePermissionToolsConsoleManagementExecute AdministrativePermission = 3503
	AdministrativePermissionToolsWebManagementExecute     AdministrativePermission = 3504

	AdministrativePermissionBusinessUnitsView   AdministrativePermission = 3600
	AdministrativePermissionBusinessUnitsAdd    AdministrativePermission = 3601
	AdministrativePermissionBusinessUnitsEdit   AdministrativePermission = 3602
	AdministrativePermissionBusinessUnitsDelete AdministrativePermission = 3603
	AdministrativePermissionBusinessUnitsAssign AdministrativePermission = 3604

	AdministrativePermissionSecurityProviderView AdministrativePermission = 3700
	AdministrativePermissionSecurityProviderEdit AdministrativePermission = 3701
)
