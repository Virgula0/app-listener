# Graph Report - app-listener  (2026-10-16)

## Corpus Check
- 374 files · ~1,568,051 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4708 nodes · 19586 edges · 180 communities (155 shown, 25 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1652 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7b6dca46`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- Vault
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- User
- guardUnitTest
- globreserve_test.go
- filevault_test.go
- time.Duration
- catalogrefresh.go
- binaryVetter
- install.go
- Candidate
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- DaemonUseCase
- auth.go
- IntegrationSuite
- guard.bpf.c
- NetworkMonitor
- github.com/charmbracelet/bubbles/textarea.Model
- LibraryClosure
- Ledger
- runUninstall
- serve.go
- bpfstats/main.go
- catalogRefresher
- fileedit_model_test.go
- liveSession
- uninstall_system.go
- openFileVaultWithMasterKey
- controlServer
- ConfinedEntry
- update.go
- daemon_multicall_test.go
- IntegrationSuite
- Monitor
- IntegrationSuite
- update_test.go
- GuardRepository
- stdio
- loadDaemonConfig
- .FilterExistingWhitelist
- planUpdaters
- inode_dev.h
- guardModel
- Config
- github.com/charmbracelet/bubbletea.Cmd
- TrustGuard
- os.FileMode
- Pulse
- GuardInodeKey
- .decryptFileWithProgress
- CandidateDir
- fileedit_symlink_test.go
- exe_supersede.h
- GuardEvent
- Resource
- pinstate.go
- gui/session.go
- runNetworkMonitor
- githubRelease
- runLockdown
- fcntl
- copyRegularAt
- KernelDev
- time.Time
- btrfsLayoutFrom
- confgen.go
- github.com/spf13/cobra.Command
- .refreshOnce
- process_vm_readv.c
- net_tester/main.go
- IntegrationSuite
- netModel
- netGuardModel
- testLimiter
- downloadAndVerify
- auditGroupCoverage
- EventType
- IntegrationSuite
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- model
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- changelog_test.go
- writeWithin
- IntegrationSuite
- Well-known bypass classes
- fakeDaemon
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- TestMonitorUseCaseLifecycle
- Build & sign release assets (reusable workflow)
- runMonitor
- .readEvent
- sync.Mutex
- verifyReleaseWithKey
- app-listener Terminal Demo Recording
- IntegrationSuite
- guardLSMHooks
- configedit_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- .exeKeys
- NewDumpHook
- WebSocket /ws client
- batchInodeRows
- launch_probe.c
- install.sh
- mc_tag.h
- harnessSuite
- make test-integration (rootful Docker suite)
- integrationtests/guard_test.go
- supersede_probe.c
- trace-app-libs.sh
- IntegrationSuite
- ParseEventsFlag
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- buildTrustedSet
- reloadOnce
- Guard
- render.go
- bottomBar
- newTestVetter
- wipe
- bpfLsmListed
- daemon: permanent protection with encryption at rest
- parseConfig
- formatEntry
- journalTo
- Installation
- os.File
- graphify-refresh.sh
- RawTarget
- parsePasswd
- runNetworkGuard
- launchKey
- Vault
- Limitations
- Troubleshooting
- Daemon configuration
- edit-protected
- absPath
- Compatibility
- How it works
- buildOneGuard
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- monitor: see what touches a path
- kernel-prof.sh
- vaultFS

## God Nodes (most connected - your core abstractions)
1. `Guard` - 114 edges
2. `Config` - 95 edges
3. `GuardInodeKey` - 94 edges
4. `Load()` - 89 edges
5. `IntegrationSuite` - 81 edges
6. `absPath()` - 79 edges
7. `writeConfig()` - 69 edges
8. `IntegrationSuite` - 67 edges
9. `fileEditModel` - 60 edges
10. `engine` - 58 edges

## Surprising Connections (you probably didn't know these)
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `Well-known bypass classes` --references--> `btrfs_search POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `Well-known bypass classes` --references--> `io_uring POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Shared bubbletea TUI layout across monitor and guard modes** — media_demo_monitor_tui, media_demo_guard_whitelist_tui, media_demo_guard_blacklist_tui, media_demo_tui_resource_footer [INFERRED 0.85]
- **Reserved-name trust defenses (trust object #3)** — integrationtests_exploits_readme_glob_plant, integrationtests_exploits_readme_lib_probe, readme_bun_tmpdir_redirect, claude_catalog_refresh [INFERRED 0.85]
- **Release build, sign and verify pipeline** — _github_workflows_release_prerelease_workflow, _github_workflows_promote_release_promote_release_workflow, _github_workflows__build_release_build_release_workflow, _github_workflows__build_release_ed25519_signing, readme_install_sh, readme_update [INFERRED 0.95]
- **Verifier budget measurement and gating** — _github_workflows_verifier_verifier_gate, claude_bpfstats, readme_daemon_check, claude_verifier_complexity_budget, _github_workflows__build_release_ebpf_regeneration [INFERRED 0.95]

## Communities (180 total, 25 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.12
Nodes (21): Execute(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey (+13 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (110): TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestDiffCatalogFlagRegistered(), collectFilesystemPrereqs(), discordCatalogEntry(), groupedDiscordConf(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic() (+102 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.05
Nodes (14): IntegrationSuite, IntegrationSuite, IntegrationSuite, configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix() (+6 more)

### Community 3 - "Vault"
Cohesion: 0.15
Nodes (18): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), encryptDirectories(), secureResources(), verifyEncryptionState(), confirmRunPrereq(), resolveFilesystemPrereqs() (+10 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (77): TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+69 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.14
Nodes (31): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+23 more)

### Community 7 - "User"
Cohesion: 0.10
Nodes (47): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), deploy(), preflightDeployedBinary(), addKeysToAgent() (+39 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (20): addErrKind, guardUnitTest, ReplacementCheck, BinariesSummary(), classifyAddErr(), guardModeKey(), addRecorder(), copySelf() (+12 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.06
Nodes (63): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+55 more)

### Community 10 - "filevault_test.go"
Cohesion: 0.19
Nodes (23): looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsFileVaultRecordShape() (+15 more)

### Community 11 - "time.Duration"
Cohesion: 0.22
Nodes (6): TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), debouncer

### Community 12 - "catalogrefresh.go"
Cohesion: 0.38
Nodes (3): isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 13 - "binaryVetter"
Cohesion: 0.09
Nodes (21): newBinaryVetter(), openBinaryVetter(), makeReloadHandler(), noteReloadedSettings(), startCatalogRefresh(), newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows() (+13 more)

### Community 14 - "install.go"
Cohesion: 0.07
Nodes (43): confirmMasterKeyOverwrite(), runGenKey(), applyLiveRefresh(), checkRunningBinaryMatchesInstalled(), ensureInstalledBinary(), ensureMasterKey(), installBinaryOnly(), parseMaintenanceFlags() (+35 more)

### Community 15 - "Candidate"
Cohesion: 0.08
Nodes (41): appendSectionsAndEdit(), applyDiffAdditions(), collectDiffAdditions(), configPaths(), pathCovered(), restoreDaemonAfterDiffAbort(), runDiffCatalog(), TestDiffMergePreservesExistingAndParses() (+33 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.12
Nodes (3): IntegrationSuite, guardDeltaEvents(), le64HexKey()

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.18
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 19 - "fileEditModel"
Cohesion: 0.13
Nodes (8): clipLabel(), humanSize(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind, fileNode

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.26
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (48): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+40 more)

### Community 22 - "DaemonUseCase"
Cohesion: 0.17
Nodes (3): concurrencyLimit(), DaemonUseCase, unlockRoots()

### Community 23 - "auth.go"
Cohesion: 0.15
Nodes (25): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+17 more)

### Community 24 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (70): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+62 more)

### Community 26 - "NetworkMonitor"
Cohesion: 0.11
Nodes (13): NetEvent, NetEventType, NetEventTypes(), ParseNetEventType(), eventsetSummary(), eventTypeKey(), modeLabel(), NewNetGuard() (+5 more)

### Community 27 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.05
Nodes (68): newChangelogModel(), showChangelog(), class, findButton, span, diffModel, editorModel, classOf() (+60 more)

### Community 28 - "LibraryClosure"
Cohesion: 0.10
Nodes (19): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), isScript(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+11 more)

### Community 29 - "Ledger"
Cohesion: 0.14
Nodes (9): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), Ledger, Pending (+1 more)

### Community 30 - "runUninstall"
Cohesion: 0.07
Nodes (39): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), deletePostBackups(), restoreBackups(), askDecrypt(), decryptDirectories(), decryptStep() (+31 more)

### Community 31 - "serve.go"
Cohesion: 0.07
Nodes (35): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+27 more)

### Community 32 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 33 - "catalogRefresher"
Cohesion: 0.26
Nodes (4): drain(), newCatalogRefresher(), readEvents(), catalogRefresher

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.17
Nodes (26): OpenRegularNoFollow(), fileModeToUnixMode(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd() (+18 more)

### Community 35 - "liveSession"
Cohesion: 0.12
Nodes (20): dialLiveSession(), parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), displayEntries(), eventsLabel() (+12 more)

### Community 36 - "uninstall_system.go"
Cohesion: 0.05
Nodes (64): authorizeFirstTimeSetPassword(), TestAuthorizeFirstTimeSetPasswordRefusedWhileDaemonRunning(), reloadDaemonForPasswordChange(), buildBinaryIfNeeded(), mustCwd(), removeKeyAndEmptyDir(), removeMasterKey(), revertSystemFiles() (+56 more)

### Community 37 - "openFileVaultWithMasterKey"
Cohesion: 0.22
Nodes (17): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+9 more)

### Community 38 - "controlServer"
Cohesion: 0.12
Nodes (18): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines(), streamEvents() (+10 more)

### Community 39 - "ConfinedEntry"
Cohesion: 0.13
Nodes (29): binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats() (+21 more)

### Community 40 - "update.go"
Cohesion: 0.16
Nodes (15): applyUpdate(), assetsFor(), checkUpdatePreconditions(), confirmUnit(), confirmUpdate(), isTerminal(), parseVerPart(), readInstalledVersion() (+7 more)

### Community 41 - "daemon_multicall_test.go"
Cohesion: 0.16
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 43 - "Monitor"
Cohesion: 0.06
Nodes (11): evalSymlinksOrEmpty(), NewMonitor(), newPathCache(), dirTarget(), TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace(), Monitor (+3 more)

### Community 44 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 45 - "update_test.go"
Cohesion: 0.11
Nodes (13): fetchReleases(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel(), testKeyPair(), TestNewerThanInstalled() (+5 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "loadDaemonConfig"
Cohesion: 0.12
Nodes (16): loadDaemonConfig(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds(), TestPlanTrustedDirs_DuplicateWriterAddsNoEmptyLoader(), TestPlanTrustedDirs_LibDirsOnlyForTheirWriters() (+8 more)

### Community 49 - ".FilterExistingWhitelist"
Cohesion: 0.11
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 50 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 51 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 52 - "guardModel"
Cohesion: 0.07
Nodes (28): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+20 more)

### Community 53 - "Config"
Cohesion: 0.11
Nodes (37): sameConfig(), watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector() (+29 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.10
Nodes (8): clamp(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel, serveTestModel, sessionModel

### Community 55 - "TrustGuard"
Cohesion: 0.05
Nodes (24): tempGrant, tempMaskOp, tempRow, trustRows, checkBtrfsKeys(), ensureBtrfsLayout(), deleteResKeys(), dropSupersededLocked() (+16 more)

### Community 56 - "os.FileMode"
Cohesion: 0.07
Nodes (41): applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig(), updateCatalogConfig() (+33 more)

### Community 57 - "Pulse"
Cohesion: 0.28
Nodes (7): New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse, Sub

### Community 58 - "GuardInodeKey"
Cohesion: 0.04
Nodes (43): engine, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+35 more)

### Community 59 - ".decryptFileWithProgress"
Cohesion: 0.16
Nodes (7): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), TestCopyXattrsPreservesUserAttributes(), copyXattrs()

### Community 60 - "CandidateDir"
Cohesion: 0.08
Nodes (23): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), Discover() (+15 more)

### Community 61 - "fileedit_symlink_test.go"
Cohesion: 0.24
Nodes (15): AtomicWriteAt(), openRootT(), TestWriteFileKeepMeta(), TestWriteFileKeepMetaRefusesSymlink(), mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent() (+7 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.21
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+7 more)

### Community 63 - "GuardEvent"
Cohesion: 0.05
Nodes (20): holdHeadless(), HeadlessLine(), runGuardHeadless(), bpfGuardEvent, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), commMatchesGuardedBinary(), decodeGuardEvent() (+12 more)

### Community 64 - "Resource"
Cohesion: 0.23
Nodes (7): backingDeviceUnion(), Resource, BackingDevices(), encryptionRootOf(), uniqueEncryptionRoots(), vaultOpForGuard(), vaultOpOnRoot()

### Community 65 - "pinstate.go"
Cohesion: 0.17
Nodes (23): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+15 more)

### Community 66 - "gui/session.go"
Cohesion: 0.05
Nodes (29): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+21 more)

### Community 67 - "runNetworkMonitor"
Cohesion: 0.15
Nodes (9): runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NetworkMonitorRepository, NewNetModel(), NetworkMonitorUseCase (+1 more)

### Community 68 - "githubRelease"
Cohesion: 0.29
Nodes (14): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), latestTime(), newerThanInstalled() (+6 more)

### Community 69 - "runLockdown"
Cohesion: 0.29
Nodes (8): configureDaemonLogging(), resolveConfigPath(), runLockdown(), dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "KernelDev"
Cohesion: 0.09
Nodes (33): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), openRoot(), physicalParent(), statFD() (+25 more)

### Community 73 - "time.Time"
Cohesion: 0.16
Nodes (9): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, sweepID, sweepMemo (+1 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 75 - "confgen.go"
Cohesion: 0.09
Nodes (35): logRefreshChange(), TestLogRefreshChange(), LibraryBlock, RefreshOptions, findSectionEnd(), findSectionStart(), firstContentLine(), isLibDirectiveText() (+27 more)

### Community 76 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (19): AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+11 more)

### Community 78 - ".refreshOnce"
Cohesion: 0.16
Nodes (14): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), grantNewBinaries(), Lines(), TestLinesKeyOnLinkAndIncludePending() (+6 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 82 - "netModel"
Cohesion: 0.21
Nodes (8): formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg, netModel

### Community 83 - "netGuardModel"
Cohesion: 0.24
Nodes (5): NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 84 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 85 - "downloadAndVerify"
Cohesion: 0.20
Nodes (10): downloadAndVerify(), downloadFile(), downloadReleaseFiles(), sanityCheckBinary(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestSanityCheckBinaryRefusesOtherMachine(), updateHTTPClient() (+2 more)

### Community 86 - "auditGroupCoverage"
Cohesion: 0.18
Nodes (13): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+5 more)

### Community 87 - "EventType"
Cohesion: 0.06
Nodes (29): EventType, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), binaryStatOf(), TempBinary, TemporaryGrant, TempRule (+21 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "model"
Cohesion: 0.26
Nodes (6): formatType(), listenForEvents(), tickStats(), eventLine, eventMsg, model

### Community 92 - "DaemonEvent"
Cohesion: 0.10
Nodes (31): UIDResolver, NewUIDResolver(), awaitStartupOrSignal(), catchLifecycleSignals(), newDaemonModel(), runDaemon(), runDaemonUI(), runHeadless() (+23 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.09
Nodes (16): foreignPinnedLink(), codeVouched(), isTraced(), judgePasses(), loaderEnv(), parseProcs(), roFor(), TestCodeVouched() (+8 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "changelog_test.go"
Cohesion: 0.22
Nodes (9): rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog(), TestNewChangelogModelViewContainsTitleAndNotes() (+1 more)

### Community 100 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.21
Nodes (10): NewGuardUseCase(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestNetworkGuardUseCaseLifecycle(), TestNetworkMonitorUseCaseLifecycle() (+2 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "runMonitor"
Cohesion: 0.11
Nodes (13): MakeDisplayPaths(), runHeadless(), runMonitor(), runTUI(), FileEvent, Run(), MonitorRepository, NewModel() (+5 more)

### Community 110 - ".readEvent"
Cohesion: 0.11
Nodes (10): BpfEvent, eventUnitTest, NetBpfEvent, Cstr(), DecodeBpfEvent(), FormatAddr(), Ntohs(), reasonText() (+2 more)

### Community 111 - "sync.Mutex"
Cohesion: 0.17
Nodes (6): swapFile(), newControlManager(), startControlManager(), startControlServer(), configEditor, controlManager

### Community 112 - "verifyReleaseWithKey"
Cohesion: 0.40
Nodes (5): parseChecksum(), parsePublicKey(), TestParseChecksum(), verifyRelease(), verifyReleaseWithKey()

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "guardLSMHooks"
Cohesion: 0.33
Nodes (3): trustHook, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent()

### Community 116 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 121 - "NewDumpHook"
Cohesion: 0.20
Nodes (12): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+4 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "batchInodeRows"
Cohesion: 0.19
Nodes (12): markReady(), notifySystemdReady(), inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows() (+4 more)

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "integrationtests/guard_test.go"
Cohesion: 0.27
Nodes (6): attrOpCase, eventsForPath(), guardEventTypes(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 144 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 145 - "reloadOnce"
Cohesion: 0.13
Nodes (18): CheckBPFLSM(), CheckEBPF(), buildConcurrency(), buildGuards(), newPinGeneration(), prepareDaemonStart(), reloadOnce(), reloadSlotsNeeded() (+10 more)

### Community 146 - "Guard"
Cohesion: 0.05
Nodes (30): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), deferredBinary, exeRow, rootHandle, BinaryEntry (+22 more)

### Community 147 - "render.go"
Cohesion: 0.13
Nodes (13): lineStyles, viewportState, fg(), baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+5 more)

### Community 148 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 149 - "newTestVetter"
Cohesion: 0.19
Nodes (18): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+10 more)

### Community 150 - "wipe"
Cohesion: 0.14
Nodes (17): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn() (+9 more)

### Community 151 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 155 - "parseConfig"
Cohesion: 0.29
Nodes (6): parseConfig(), BtrfsMounted(), mountinfoDev(), mountinfoHasFstype(), TestMountinfoDev(), TestMountinfoHasFstype()

### Community 156 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 157 - "journalTo"
Cohesion: 0.53
Nodes (6): journalTo(), tempBinary(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure()

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.04
Nodes (78): multicallAdmissible(), ownerPaths(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), drainEphemeral(), systemRule(), vetDecision (+70 more)

### Community 161 - "RawTarget"
Cohesion: 0.27
Nodes (10): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+2 more)

### Community 163 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

### Community 164 - "runNetworkGuard"
Cohesion: 0.10
Nodes (19): ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), DiscoverInfraBinaries() (+11 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.20
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 178 - "absPath"
Cohesion: 0.09
Nodes (8): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, launchCase, absPath(), TestMain(), helperChild(), TestMain()

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 185 - "buildOneGuard"
Cohesion: 0.10
Nodes (33): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards(), resolveGuardConfig() (+25 more)

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "vaultFS"
Cohesion: 0.17
Nodes (9): TestRuntimeClassMarkers(), applyNewFileMeta(), defaultNewFileMode(), isBinaryContent(), listEntries(), TestDefaultNewFileMode(), TestIsBinaryContent(), writeFileInPlace() (+1 more)

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 450 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **25 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `IntegrationSuite`, `IntegrationSuite`, `daemon_multicall_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.036) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11892730140905323 - nodes in this community are weakly interconnected._
- **Why does `watchPattern` connect `globreserve_test.go` to `Daemon configuration`, `catalogRefresher`, `catalogrefresh.go`, `time.Time`?**
  _High betweenness centrality (0.033) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.026859183356621644 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.05202547080341417 - nodes in this community are weakly interconnected._
- **Should `daemonconfig_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.07092907092907093 - nodes in this community are weakly interconnected._