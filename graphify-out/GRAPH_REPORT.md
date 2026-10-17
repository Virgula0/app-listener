# Graph Report - app-listener  (2026-10-17)

## Corpus Check
- 379 files · ~1,573,928 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4778 nodes · 20024 edges · 219 communities (172 shown, 47 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1681 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `cb4e9630`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- GenerateConf
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- rcUser
- guardUnitTest
- mkdirs
- render.go
- runDaemon
- SelectOrphans
- update.go
- planUpdaters
- os.FileMode
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- TestHashFileLifecycle
- TestMonitorUseCaseLifecycle
- __always_inline
- SanitizeText
- buildTrustedSet
- inodeChain
- codeedit.go
- ListUsers
- serve.go
- bpfstats/main.go
- time.Duration
- fileedit_model_test.go
- github.com/spf13/cobra.Command
- Pulse
- hashFileShared
- sync.Mutex
- expandPlaceholders
- runMonitor
- IntegrationSuite
- absPath
- monitorUnitTest
- IntegrationSuite
- update_test.go
- BinaryStat
- stdio
- golang.org/x/sys/unix.Stat_t
- configedit_test.go
- LibraryClosure
- filevault_test.go
- tempGrantSetup
- daemonconfig.go
- github.com/charmbracelet/bubbletea.Model
- TrustGuard
- daemon_multicall_test.go
- systemd.go
- engine
- Read
- buildOneGuard
- MountVouch
- exe_supersede.h
- GuardEvent
- PlanDaemonUnitUpdate
- pinstate.go
- Run
- RawTarget
- changelog.go
- openBunTmp
- fcntl
- copyRegularAt
- fileedit_symlink_test.go
- time.Time
- btrfsLayoutFrom
- User
- lockAndDeprovision
- .runEditorHarness
- NetworkMonitor
- process_vm_readv.c
- fakeGuardRepo
- IntegrationSuite
- watchPattern
- runUninstall
- Monitor
- netModel
- github.com/charmbracelet/bubbles/textarea.Model
- highlight_test.go
- highlight.go
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- filevault.go
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- revertSSHAgents
- runOfflineEdit
- IntegrationSuite
- Well-known bypass classes
- auditGroupCoverage
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- runNetworkMonitor
- Build & sign release assets (reusable workflow)
- .populateWatchInodes
- TestBatchInodeRows
- mountTable
- IntegrationSuite
- app-listener Terminal Demo Recording
- IntegrationSuite
- guard.bpf.c
- check_watched
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- RefreshSection
- gui/session.go
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- bottomBar
- make test-integration (rootful Docker suite)
- IntegrationSuite
- supersede_probe.c
- liveSession
- trace-app-libs.sh
- IntegrationSuite
- vaultFS
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- EventType
- IntegrationSuite
- Candidate
- pathCache
- Guard
- net_tester/main.go
- .eventTypeAllowed
- newTestVetter
- Vault
- launchKey
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- eventUnitTest
- rootHandle
- KernelDev
- Installation
- GuardInodeKey
- graphify-refresh.sh
- harnessSuite
- Vault
- main_test.go
- runEditProtected
- ParseEventsFlag
- bpfLsmListed
- IntegrationSuite
- dataRowRenderer
- parseGuardEvent
- find_test.go
- IntegrationSuite
- Limitations
- Troubleshooting
- integrationtests/networkguard_test.go
- mountsUnitTest
- Daemon configuration
- edit-protected
- runNetworkGuard
- github.com/charmbracelet/bubbletea.Cmd
- Compatibility
- How it works
- inode_dev.h
- downloadAndVerify
- changelog_test.go
- find.go
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- runSetPassword
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- launch_rule_at
- IntegrationSuite
- monitor: see what touches a path
- kernel-prof.sh
- ResolveBtrfsLayout
- FileEvent
- testLimiter
- Verify
- findTopLevelWindow
- .exeKeys
- runForward
- io.Reader
- DiscoverInfraBinaries
- golang.org/x/sys/unix.Timespec
- MulticallError

## God Nodes (most connected - your core abstractions)
1. `Guard` - 117 edges
2. `Config` - 96 edges
3. `GuardInodeKey` - 95 edges
4. `Load()` - 90 edges
5. `IntegrationSuite` - 81 edges
6. `absPath()` - 81 edges
7. `writeConfig()` - 69 edges
8. `IntegrationSuite` - 68 edges
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

## Communities (219 total, 47 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.11
Nodes (16): Execute(), tempJournalRow, launchRule, trustEvent, btrfsIoctlFsInfoArgs, Check(), collectDiagnostics(), readSysctl() (+8 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (131): awaitStartupOrSignal(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestAwaitStartupOrSignalAbortsAndLocksDown(), TestAwaitStartupOrSignalAbortWithFailedStartup(), TestAwaitStartupOrSignalCompletes(), atomicWriteAt(), descendNoFollow() (+123 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (11): Hash(), configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec() (+3 more)

### Community 3 - "GenerateConf"
Cohesion: 0.05
Nodes (60): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), confirmRunPrereq(), askEncryption(), sectionsFromCandidates(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse() (+52 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (73): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+65 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.11
Nodes (34): code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image(), fill_inode_key(), get_task_exe_key() (+26 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "rcUser"
Cohesion: 0.12
Nodes (34): setupBunTmpdir(), rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), EnsureBunTmpDir(), ensureRCBlock() (+26 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (22): addErrKind, guardUnitTest, ReplacementCheck, classifyAddErr(), commMatchesGuardedBinary(), guardModeKey(), addRecorder(), copySelf() (+14 more)

### Community 9 - "mkdirs"
Cohesion: 0.12
Nodes (37): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+29 more)

### Community 10 - "render.go"
Cohesion: 0.11
Nodes (13): lineStyles, viewportState, fg(), baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+5 more)

### Community 11 - "runDaemon"
Cohesion: 0.10
Nodes (21): newCatalogRefresher(), newControlManager(), startControlManager(), makeReloadHandler(), markReady(), noteReloadedSettings(), notifySystemdReady(), runDaemon() (+13 more)

### Community 12 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 13 - "update.go"
Cohesion: 0.11
Nodes (17): showChangelog(), applyUpdate(), assetsFor(), checkUpdatePreconditions(), confirmUnit(), confirmUpdate(), isTerminal(), parseVerPart() (+9 more)

### Community 14 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 15 - "os.FileMode"
Cohesion: 0.10
Nodes (28): TestDaemonUnitWithMetadataOutput(), dupFD(), installConfig(), installFile(), installFileAs(), installServices(), installSSHAgent(), installSSHAgentEnv() (+20 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 19 - "fileEditModel"
Cohesion: 0.12
Nodes (8): clipLabel(), humanSize(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind, fileNode

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (49): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+41 more)

### Community 22 - "Resource"
Cohesion: 0.09
Nodes (9): Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots(), vaultOpForGuard() (+1 more)

### Community 23 - "TestHashFileLifecycle"
Cohesion: 0.33
Nodes (10): HashFileExists(), LoadHashFile(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle(), TestRemoveHashFileInPlace(), TestWriteHashFileInPlaceOverPlaceholder(), WriteHashFile() (+2 more)

### Community 24 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.33
Nodes (8): NewGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestNetworkGuardUseCaseLifecycle(), TestNetworkMonitorUseCaseLifecycle(), fakeNetworkGuardRepo

### Community 25 - "__always_inline"
Cohesion: 0.18
Nodes (28): add_inode_to_guard(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), evict_inode_from_guard(), get_dentry_from_path(), get_inode_from_path(), guard_inode_getattr() (+20 more)

### Community 26 - "SanitizeText"
Cohesion: 0.04
Nodes (41): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), drainEphemeral(), newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows() (+33 more)

### Community 27 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 28 - "inodeChain"
Cohesion: 0.21
Nodes (16): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD(), keyOf() (+8 more)

### Community 29 - "codeedit.go"
Cohesion: 0.30
Nodes (14): class, classOf(), dedent(), Edit(), leadingIndent(), lines(), Navigate(), onlySpaces() (+6 more)

### Community 30 - "ListUsers"
Cohesion: 0.08
Nodes (32): askBunTmpdirUsers(), gatherPerUserSetup(), applyDiffAdditions(), configPaths(), pathCovered(), restoreDaemonAfterDiffAbort(), TestPathCovered(), checkRunningBinaryMatchesInstalled() (+24 more)

### Community 31 - "serve.go"
Cohesion: 0.06
Nodes (36): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+28 more)

### Community 32 - "bpfstats/main.go"
Cohesion: 0.09
Nodes (33): CheckBPFLSM(), confirmMasterKeyOverwrite(), runBPFCheck(), runGenKey(), runOneShotMode(), ensureMasterKey(), checkKeyLen(), GenerateMasterKey() (+25 more)

### Community 33 - "time.Duration"
Cohesion: 0.18
Nodes (6): drain(), readEvents(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), catalogRefresher, debouncer

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.17
Nodes (26): OpenRegularNoFollow(), fileModeToUnixMode(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd() (+18 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.08
Nodes (33): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+25 more)

### Community 36 - "Pulse"
Cohesion: 0.28
Nodes (7): New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse, Sub

### Community 37 - "hashFileShared"
Cohesion: 0.23
Nodes (8): HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled(), unchangedSince()

### Community 38 - "sync.Mutex"
Cohesion: 0.10
Nodes (21): swapFile(), readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines() (+13 more)

### Community 39 - "expandPlaceholders"
Cohesion: 0.07
Nodes (26): symlinkTargetGlob(), confSafeMatch(), expandPlaceholders(), BinaryRule, homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestExpandWhitelist() (+18 more)

### Community 40 - "runMonitor"
Cohesion: 0.08
Nodes (24): CheckEBPF(), configureDaemonLogging(), loadDaemonConfig(), newPinGeneration(), prepareDaemonStart(), resolveConfigPath(), runLockdown(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup() (+16 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 45 - "update_test.go"
Cohesion: 0.12
Nodes (18): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel() (+10 more)

### Community 46 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 49 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 50 - "LibraryClosure"
Cohesion: 0.12
Nodes (16): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), isScript(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+8 more)

### Community 51 - "filevault_test.go"
Cohesion: 0.19
Nodes (25): looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault(), sealFileVault(), inode(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsEncryptedDispatchesToFileVaultForRegularFiles() (+17 more)

### Community 52 - "tempGrantSetup"
Cohesion: 0.57
Nodes (7): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure()

### Community 53 - "daemonconfig.go"
Cohesion: 0.15
Nodes (38): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+30 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Model"
Cohesion: 0.12
Nodes (10): NewDaemonModel(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), assertViewFits(), TestViewsFitReportedTerminalSize(), DaemonResourceInfo, noticeModel (+2 more)

### Community 55 - "TrustGuard"
Cohesion: 0.09
Nodes (12): tempGrant, tempMaskOp, tempRow, trustHook, trustRows, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), setSupersedeTrust() (+4 more)

### Community 56 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 57 - "systemd.go"
Cohesion: 0.09
Nodes (39): reloadDaemonForPasswordChange(), applyLiveRefresh(), buildBinaryIfNeeded(), installBinaryOnly(), runUpdateCatalogOnly(), softenAutomatedRefreshErr(), TestApplyLiveRefreshDeliversReload(), TestSoftenAutomatedRefreshErr() (+31 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (30): engine, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded() (+22 more)

### Community 59 - "Read"
Cohesion: 0.07
Nodes (28): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+20 more)

### Community 60 - "buildOneGuard"
Cohesion: 0.07
Nodes (43): backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), reloadOnce(), reloadSlotsNeeded(), startGuardedDaemon() (+35 more)

### Community 61 - "MountVouch"
Cohesion: 0.09
Nodes (17): GlobChild, statFunc, appendInodeRow(), batchInodeRows(), inodeRowsOf(), deleteResKeys(), compileGlobNames(), globText() (+9 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.16
Nodes (18): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+10 more)

### Community 63 - "GuardEvent"
Cohesion: 0.08
Nodes (21): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless() (+13 more)

### Community 64 - "PlanDaemonUnitUpdate"
Cohesion: 0.27
Nodes (10): ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate(), ShippedDaemonUnit(), fakeBinary(), TestPlanDaemonUnitUpdate(), TestPlanDaemonUnitUpdateDropsLocalEdits(), TestShippedDaemonUnit() (+2 more)

### Community 65 - "pinstate.go"
Cohesion: 0.15
Nodes (25): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+17 more)

### Community 66 - "Run"
Cohesion: 0.12
Nodes (9): columnState, dataRow, guiModel, headerWidget, newColumnState(), newDataRow(), newHeaderWidget(), Run() (+1 more)

### Community 67 - "RawTarget"
Cohesion: 0.24
Nodes (11): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+3 more)

### Community 68 - "changelog.go"
Cohesion: 0.29
Nodes (14): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), latestTime(), newerThanInstalled() (+6 more)

### Community 69 - "openBunTmp"
Cohesion: 0.39
Nodes (8): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 72 - "fileedit_symlink_test.go"
Cohesion: 0.17
Nodes (16): AtomicWriteAt(), openRootT(), TestWriteFileKeepMeta(), TestWriteFileKeepMetaRefusesSymlink(), mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent() (+8 more)

### Community 73 - "time.Time"
Cohesion: 0.17
Nodes (9): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, sweepID, sweepMemo (+1 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.17
Nodes (9): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), btrfsTypes(), TestBtrfsLayoutFromNestedMembers() (+1 more)

### Community 75 - "User"
Cohesion: 0.14
Nodes (21): entryWriters(), reserveChain(), TestGlobBuilder_LibBitsPerEntry(), underResource(), bunEntryConfigured(), bunLaunchersFor(), globBuilder, ParseGlobName() (+13 more)

### Community 76 - "lockAndDeprovision"
Cohesion: 0.11
Nodes (14): deprovisionKind, migrationTarget, writeTempWithMetadata(), classifyDeprovision(), classifyDeprovisionErr(), classifyMigrationTarget(), isDeprovisionBusy(), isDeprovisionMissing() (+6 more)

### Community 78 - "NetworkMonitor"
Cohesion: 0.08
Nodes (15): NetBpfEvent, NetEvent, NetEventType, FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType(), eventsetSummary() (+7 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 82 - "watchPattern"
Cohesion: 0.21
Nodes (11): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), addSystemBinary(), catalogPatterns(), catalogWatchPlan(), inHome(), systemBinaryPatterns(), watchPattern (+3 more)

### Community 83 - "runUninstall"
Cohesion: 0.15
Nodes (23): parseMaintenanceFlags(), runMaintenanceMode(), TestValidateMaintenanceFlags(), validateMaintenanceFlags(), deletePostBackups(), restoreBackups(), offerBackupCleanup(), runUninstall() (+15 more)

### Community 84 - "Monitor"
Cohesion: 0.15
Nodes (3): NewMonitor(), Monitor, watchedInode

### Community 85 - "netModel"
Cohesion: 0.21
Nodes (8): formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg, netModel

### Community 86 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.28
Nodes (4): column(), repeatKey(), reposition(), Finder

### Community 87 - "highlight_test.go"
Cohesion: 0.21
Nodes (21): moveTo(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines() (+13 more)

### Community 88 - "highlight.go"
Cohesion: 0.09
Nodes (9): diffModel, editorModel, Highlighter, lexerFor(), shebangLexer(), TestDaemonConfTokens(), TestTokenizeKeepsLineCount(), tokenClassOf() (+1 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "filevault.go"
Cohesion: 0.20
Nodes (16): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+8 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.10
Nodes (23): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+15 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.10
Nodes (13): foreignPinnedLink(), reasonText(), eventTypeKey(), modeLabel(), NewNetGuard(), judgePasses(), parseProcs(), TestJudgePasses() (+5 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "revertSSHAgents"
Cohesion: 0.25
Nodes (9): detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeSSHAgentEnv(), removeSSHAgentUnit(), revertSSHAgents(), staleRCUsers(), TestIsInstallerSSHAgentUnit(), TestRemoveSSHAgentUnit() (+1 more)

### Community 100 - "runOfflineEdit"
Cohesion: 0.10
Nodes (20): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep(), pickDirsToDecrypt() (+12 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "auditGroupCoverage"
Cohesion: 0.18
Nodes (13): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+5 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "runNetworkMonitor"
Cohesion: 0.15
Nodes (10): ParseNetEventsFlag(), runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NetworkMonitorRepository, NewNetModel() (+2 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 110 - "TestBatchInodeRows"
Cohesion: 0.50
Nodes (5): inodeBatchFn, fakeBatches(), TestBatchInodeRows(), countInodeRows(), TestCountInodeRows()

### Community 111 - "mountTable"
Cohesion: 0.16
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+1 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "guard.bpf.c"
Cohesion: 0.12
Nodes (40): exe_refused(), inode_dev(), check_and_emit_args(), code_suspect(), current_is_inspector(), emit_process_denial(), emit_with(), event_is_read_class() (+32 more)

### Community 116 - "check_watched"
Cohesion: 0.24
Nodes (14): admit_refusal(), check_watched(), emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg() (+6 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "RefreshSection"
Cohesion: 0.08
Nodes (30): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+22 more)

### Community 121 - "gui/session.go"
Cohesion: 0.21
Nodes (14): ancestorEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue(), TestOwnedByChecksTypeAndOwner() (+6 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "liveSession"
Cohesion: 0.26
Nodes (6): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), displayEntries(), holdForward(), liveSession

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "vaultFS"
Cohesion: 0.14
Nodes (12): checkFreshDir(), openDirNoFollow(), renewDir(), TestRuntimeClassMarkers(), applyNewFileMeta(), defaultNewFileMode(), isBinaryContent(), listEntries() (+4 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "EventType"
Cohesion: 0.06
Nodes (27): EventType, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps() (+19 more)

### Community 144 - "Candidate"
Cohesion: 0.11
Nodes (27): nextConfig(), appendSectionsAndEdit(), collectDiffAdditions(), runDiffCatalog(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), cleanOrphanedFscrypt() (+19 more)

### Community 145 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 146 - "Guard"
Cohesion: 0.05
Nodes (30): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), binaryVerifyState, deferredBinary, exeRow, BinaryEntry (+22 more)

### Community 147 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 148 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 149 - "newTestVetter"
Cohesion: 0.44
Nodes (10): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+2 more)

### Community 150 - "Vault"
Cohesion: 0.15
Nodes (17): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+9 more)

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.26
Nodes (20): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+12 more)

### Community 154 - "eventUnitTest"
Cohesion: 0.27
Nodes (4): BpfEvent, eventUnitTest, Cstr(), DecodeBpfEvent()

### Community 155 - "rootHandle"
Cohesion: 0.24
Nodes (4): rootHandle, openRoot(), TestNoteRootMountFollowsMountID(), TestReadOnlyWalkToleratesDanglingSymlinks()

### Community 157 - "KernelDev"
Cohesion: 0.18
Nodes (14): KernelDev(), BtrfsDevices(), btrfsSysfsDevices(), TestBtrfsSysfsDevices(), codeVouched(), isTraced(), loaderEnv(), readMountDevs() (+6 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "GuardInodeKey"
Cohesion: 0.03
Nodes (79): setConfinedHomes(), refusedReplacements, VettedInode, Multicall, MulticallKind, supersedeUnnamed(), appletTag(), buildOf() (+71 more)

### Community 163 - "main_test.go"
Cohesion: 0.18
Nodes (5): infraContainerPaths(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

### Community 164 - "runEditProtected"
Cohesion: 0.16
Nodes (16): dialLiveSession(), LiveModeAvailable(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig(), runEditConfig() (+8 more)

### Community 165 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 166 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 169 - "parseGuardEvent"
Cohesion: 0.15
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 170 - "find_test.go"
Cohesion: 0.48
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 174 - "integrationtests/networkguard_test.go"
Cohesion: 0.27
Nodes (5): IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 176 - "Daemon configuration"
Cohesion: 0.22
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 178 - "runNetworkGuard"
Cohesion: 0.10
Nodes (17): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), Mode, NetGuardEvent (+9 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "inode_dev.h"
Cohesion: 0.24
Nodes (8): ctx_ptr(), sb_dev(), btrfs_copy_gate(), guard_file_ioctl(), guard_file_ioctl_compat(), guard_inode_free(), is_btrfs_copy_ioctl(), trust_inode_free()

### Community 183 - "downloadAndVerify"
Cohesion: 0.21
Nodes (9): downloadAndVerify(), downloadFile(), downloadReleaseFiles(), sanityCheckBinary(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestSanityCheckBinaryRefusesOtherMachine(), updateHTTPClient() (+1 more)

### Community 184 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 185 - "find.go"
Cohesion: 0.24
Nodes (8): findButton, span, findAll(), lowerRunes(), newFindInput(), overlay(), runesEqual(), TestFindAllSmartCase()

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - "runSetPassword"
Cohesion: 0.22
Nodes (10): TestValidatePassword(), ValidatePassword(), promptPassword(), authorizeFirstTimeSetPassword(), TestAuthorizeFirstTimeSetPasswordRefusedWhileDaemonRunning(), PromptNewPassword(), runSetPassword(), promptEditPassword() (+2 more)

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

### Community 194 - "launch_rule_at"
Cohesion: 0.40
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "ResolveBtrfsLayout"
Cohesion: 0.22
Nodes (9): checkBtrfsKeys(), ensureBtrfsLayout(), FirstOnBtrfs(), OnBtrfs(), ResolveBtrfsLayout(), SuperblockDev(), SuperblockDevPath(), TestResolveBtrfsLayoutHostBTF() (+1 more)

### Community 200 - "FileEvent"
Cohesion: 0.25
Nodes (4): FileEvent, Run(), newFakeMonitorRepo(), fakeMonitorRepo

### Community 201 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 202 - "Verify"
Cohesion: 0.43
Nodes (7): OriginOf(), parseHash(), TestHashVerifyRoundTrip(), TestVerifyMalformed(), Verify(), Origin, parsedHash

### Community 203 - "findTopLevelWindow"
Cohesion: 0.38
Nodes (4): findTopLevelWindow(), floatWindow(), internAtom(), setFloatHint()

### Community 205 - "runForward"
Cohesion: 0.40
Nodes (6): chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), normalizedEvents(), runForward()

### Community 206 - "io.Reader"
Cohesion: 0.33
Nodes (5): BtrfsMounted(), mountinfoDev(), mountinfoHasFstype(), TestMountinfoDev(), TestMountinfoHasFstype()

### Community 207 - "DiscoverInfraBinaries"
Cohesion: 0.33
Nodes (6): DiscoverInfraBinaries(), infraFromRunning(), runningExecutables(), pathIsRunning(), TestDiscoverInfraBinaries(), TestInfraFromRunning()

## Knowledge Gaps
- **32 isolated node(s):** `IntegrationSuite`, `github.com/Virgula0/app-listener`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 462 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **47 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `integrationtests/networkguard_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `daemon_multicall_test.go`?**
  _High betweenness centrality (0.042) - this node is a cross-community bridge._
- **What connects `IntegrationSuite`, `github.com/Virgula0/app-listener`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10671022870067856 - nodes in this community are weakly interconnected._
- **Why does `Process inspectors` connect `Daemon configuration` to `watchPattern`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02073365231259968 - nodes in this community are weakly interconnected._
- **Why does `watchPattern` connect `watchPattern` to `daemon/daemon.go`, `mkdirs`, `time.Time`, `time.Duration`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0879746835443038 - nodes in this community are weakly interconnected._