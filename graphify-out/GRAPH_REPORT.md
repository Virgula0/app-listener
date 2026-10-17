# Graph Report - app-listener  (2026-10-17)

## Corpus Check
- 423 files · ~1,606,854 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 5373 nodes · 22904 edges · 203 communities (172 shown, 31 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 1963 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `430314ba`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- confgen.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- globreserve_test.go
- render.go
- Config
- ui.go
- guardModel
- mayUpdate
- fscrypt/orphans.go
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- NetGuardEvent
- guard.bpf.c
- SanitizeText
- networkPage
- inodeChain
- codeedit.go
- journal
- serve.go
- bpfstats/main.go
- go_pkg_path_filepath
- fileedit_model_test.go
- github.com/spf13/cobra.Command
- IntegrationSuite
- hashFileShared
- controlServer
- .objs
- Guard
- IntegrationSuite
- absPath
- monitorUnitTest
- IntegrationSuite
- update_test.go
- daemon_multicall_test.go
- stdio
- golang.org/x/sys/unix.Stat_t
- catalogRefresher
- LibraryClosure
- uninstall_system.go
- tempGrantSetup
- daemonconfig.go
- github.com/charmbracelet/bubbletea.Model
- TrustGuard
- runView
- install/install.go
- engine
- netModel
- buildOneGuard
- .Refresh
- exe_supersede.h
- NetworkMonitor
- installPage
- pinstate.go
- gui/session.go
- runMonitor
- changelog.go
- sync.Mutex
- fcntl
- copyRegularAt
- fileedit_symlink_test.go
- time.Time
- btrfsLayoutFrom
- Network
- parseGuardEvent
- icon
- NetEventType
- process_vm_readv.c
- GuardEvent
- .decryptFileWithProgress
- newTestTarget
- wizard.go
- Monitor
- runNetworkGuard
- github.com/charmbracelet/bubbles/textarea.Model
- liveSession
- highlight_test.go
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- GlobReservations
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- github.com/charmbracelet/bubbletea.Cmd
- policyEditor
- IntegrationSuite
- Well-known bypass classes
- auditGroupCoverage
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- time.Duration
- Build & sign release assets (reusable workflow)
- widgets.go
- KernelDev
- mountTable
- TestMonitorUseCaseLifecycle
- app-listener Terminal Demo Recording
- IntegrationSuite
- IntegrationSuite
- RefreshCatalog
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- runEditProtected
- New
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- CandidateDir
- make test-integration (rootful Docker suite)
- IntegrationSuite
- supersede_probe.c
- .policyTools
- trace-app-libs.sh
- IntegrationSuite
- vaultFS
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TempBinary
- io.Reader
- unit.go
- pathCache
- SuperblockDev
- net_tester/main.go
- .eventTypeAllowed
- newTestVetter
- Vault
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- EventType
- main_test.go
- watchPattern
- IntegrationSuite
- Installation
- GuardInodeKey
- graphify-refresh.sh
- Pulse
- Highlighter
- lockAndDeprovision
- openBunTmp
- feed.go
- shell
- writeWithin
- github.com/stretchr/testify/suite.Suite
- Execute
- find_test.go
- IntegrationSuite
- Limitations
- Troubleshooting
- newWatchPattern
- mountsUnitTest
- Daemon configuration
- edit-protected
- changelog_test.go
- configedit_test.go
- Compatibility
- How it works
- .netIntegrityContainer
- IntegrationSuite
- lockRootRecovering
- stripTempJournal
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- .runEditorHarness
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- dialLiveSession
- runForward
- TestBatchInodeRows
- monitor: see what touches a path
- kernel-prof.sh
- launchKey
- .populateWatchInodes
- runLockdown
- Desktop GUI

## God Nodes (most connected - your core abstractions)
1. `Guard` - 117 edges
2. `Config` - 96 edges
3. `GuardInodeKey` - 96 edges
4. `Load()` - 90 edges
5. `IntegrationSuite` - 82 edges
6. `absPath()` - 81 edges
7. `writeConfig()` - 69 edges
8. `IntegrationSuite` - 68 edges
9. `BinaryEntry` - 60 edges
10. `fileEditModel` - 60 edges

## Surprising Connections (you probably didn't know these)
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `runConfigEditor()` --calls--> `EditText()`  [INFERRED]
  cmd/functions/install/steps.go → internal/install/editor.go
- `resolveAssets()` --calls--> `asset()`  [INFERRED]
  cmd/functions/update/update.go → internal/gui/icons.go

## Import Cycles
- None detected.

## Communities (203 total, 31 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.13
Nodes (18): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey, attrOpCase (+10 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (109): loadDaemonConfig(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), TestDiffMergePreservesExistingAndParses(), TestDiffCatalogFlagRegistered(), discordCatalogEntry() (+101 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "confgen.go"
Cohesion: 0.04
Nodes (80): logRefreshChange(), TestLogRefreshChange(), appendSectionsAndEdit(), collectDiffAdditions(), configPaths(), pathCovered(), runDiffCatalog(), TestPathCovered() (+72 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (78): TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+70 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "shellrc.go"
Cohesion: 0.05
Nodes (85): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), restoreDaemonAfterDiffAbort(), deploy() (+77 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (22): addErrKind, ReplacementCheck, BinariesSummary(), classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest (+14 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.26
Nodes (24): buildGlobReservations(), bitOf(), discordLibConfig(), hasChild(), mkdirs(), steamLibDirConfig(), TestBuildGlobReservations_BunRootsListed(), TestBuildGlobReservations_BunTmpdir() (+16 more)

### Community 10 - "render.go"
Cohesion: 0.10
Nodes (14): findButton, lineStyles, span, findAll(), lowerRunes(), newFindInput(), overlay(), runesEqual() (+6 more)

### Community 11 - "Config"
Cohesion: 0.08
Nodes (41): ServeConfig, newCatalogRefresher(), catchLifecycleSignals(), makeReloadHandler(), markReady(), newDaemonModel(), noteReloadedSettings(), notifySystemdReady() (+33 more)

### Community 12 - "ui.go"
Cohesion: 0.07
Nodes (29): installHooks(), releaseBuild(), Run(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary() (+21 more)

### Community 13 - "guardModel"
Cohesion: 0.07
Nodes (28): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+20 more)

### Community 14 - "mayUpdate"
Cohesion: 0.20
Nodes (11): hardLinked(), mayUpdate(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning(), TestIsGeneralTool(), TestIsGeneralToolJudgesTheInode(), TestMayUpdate_HardLinkWaivedForUserWritableLibBinary(), TestPlanUpdaters_ScopedToOwnResources() (+3 more)

### Community 15 - "fscrypt/orphans.go"
Cohesion: 0.23
Nodes (19): appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors(), orphanedPolicies() (+11 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (6): eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (57): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+49 more)

### Community 19 - "fileEditModel"
Cohesion: 0.11
Nodes (11): clipLabel(), humanSize(), isBinaryContent(), listEntries(), sortNodes(), TestIsBinaryContent(), chownUser, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (49): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+41 more)

### Community 22 - "Resource"
Cohesion: 0.05
Nodes (16): libBinaries(), planUpdaters(), resourceBinaries(), resourceLinks(), TestPlanUpdaters_LinkOwnsItsResourceBit(), fakeDaemon, Resource, GuardRepository (+8 more)

### Community 23 - "auth.go"
Cohesion: 0.20
Nodes (17): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+9 more)

### Community 24 - "NetGuardEvent"
Cohesion: 0.10
Nodes (10): snapNetGuard, NetGuardEvent, NetworkGuardRepository, protoString(), NetworkGuardUseCase, NewNetworkGuardUseCase(), TestNetworkGuardUseCaseStartError(), netGuardEventLine (+2 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "SanitizeText"
Cohesion: 0.05
Nodes (33): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), grantNewBinaries() (+25 more)

### Community 27 - "networkPage"
Cohesion: 0.14
Nodes (7): TestValidateForwardFlags(), binaryList, networkPage, watchTab, TestUnchangedPerShared(), newBinaryList(), newWatchTab()

### Community 28 - "inodeChain"
Cohesion: 0.23
Nodes (15): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD(), keyOf() (+7 more)

### Community 29 - "codeedit.go"
Cohesion: 0.20
Nodes (20): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), Navigate() (+12 more)

### Community 30 - "journal"
Cohesion: 0.14
Nodes (6): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 31 - "serve.go"
Cohesion: 0.07
Nodes (35): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+27 more)

### Community 32 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 33 - "go_pkg_path_filepath"
Cohesion: 0.26
Nodes (3): isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.17
Nodes (25): fileModeToUnixMode(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod() (+17 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.08
Nodes (33): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+25 more)

### Community 37 - "hashFileShared"
Cohesion: 0.21
Nodes (9): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+1 more)

### Community 38 - "controlServer"
Cohesion: 0.08
Nodes (23): readConfigPut(), TestReadGrantRequestConfig(), newControlManager(), readAuthRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines() (+15 more)

### Community 39 - ".objs"
Cohesion: 0.10
Nodes (8): exeRow, tempGrant, tempMaskOp, tempRow, TestAppletKeyTag(), TestWithRealKeysAddsFileKeyOfApplets(), withRealKeys(), liveTrust()

### Community 40 - "Guard"
Cohesion: 0.05
Nodes (35): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), deferredBinary, rootHandle, BinaryEntry, SectionScan (+27 more)

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
Cohesion: 0.07
Nodes (27): downloadFile(), InstalledVersion(), parseChecksum(), parsePublicKey(), readInstalledVersion(), rsaPublicKeyPEM(), signedChecksum(), TestAssetsFor() (+19 more)

### Community 46 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 49 - "catalogRefresher"
Cohesion: 0.13
Nodes (14): drain(), readEvents(), sameConfig(), swapFile(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails() (+6 more)

### Community 50 - "LibraryClosure"
Cohesion: 0.07
Nodes (30): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+22 more)

### Community 51 - "uninstall_system.go"
Cohesion: 0.05
Nodes (49): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), removeKeyAndEmptyDir(), removeMasterKey() (+41 more)

### Community 52 - "tempGrantSetup"
Cohesion: 0.23
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.14
Nodes (36): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+28 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Model"
Cohesion: 0.09
Nodes (13): NewDaemonModel(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), assertViewFits(), TestViewsFitReportedTerminalSize(), RunFileEditorSession(), RunSession() (+5 more)

### Community 55 - "TrustGuard"
Cohesion: 0.08
Nodes (17): trustHook, trustRows, deleteResKeys(), PruneSuperseded(), setSupersedeTrust(), supersedeMaps(), supersedePruneLoop(), warnLostStamps() (+9 more)

### Community 56 - "runView"
Cohesion: 0.07
Nodes (13): installUI, logHook, logLine, multiAnswer, runView, stageItem, stageState, textAnswer (+5 more)

### Community 57 - "install/install.go"
Cohesion: 0.05
Nodes (65): promptEditPassword(), applyLiveRefresh(), buildBinaryIfNeeded(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), cleanupBackups(), encryptDirectories(), ensureInstalledBinary() (+57 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (31): engine, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded() (+23 more)

### Community 59 - "netModel"
Cohesion: 0.17
Nodes (11): runNetworkMonitor(), runTUI(), formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), NewNetModel(), tickNetStats() (+3 more)

### Community 60 - "buildOneGuard"
Cohesion: 0.06
Nodes (52): CheckBPFLSM(), CheckEBPF(), backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), newPinGeneration() (+44 more)

### Community 61 - ".Refresh"
Cohesion: 0.04
Nodes (10): appTheme, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, newNavItem(), rgb() (+2 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.14
Nodes (20): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+12 more)

### Community 63 - "NetworkMonitor"
Cohesion: 0.16
Nodes (5): NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 64 - "installPage"
Cohesion: 0.20
Nodes (5): binSource, installedBinary, installPage, hashReader(), readInstalled()

### Community 65 - "pinstate.go"
Cohesion: 0.27
Nodes (15): pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID(), TestPinStateRoundTrip(), TestReadPinStateMalformed() (+7 more)

### Community 66 - "gui/session.go"
Cohesion: 0.18
Nodes (16): ancestorEnv(), cleanupSessionEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField() (+8 more)

### Community 67 - "runMonitor"
Cohesion: 0.08
Nodes (23): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), ParseEventsFlag(), pathWithinDir(), ResolveTargets(), validatePair() (+15 more)

### Community 68 - "changelog.go"
Cohesion: 0.12
Nodes (34): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), showChangelog(), applyUpdate(), assetsFor(), checkUpdatePreconditions() (+26 more)

### Community 69 - "sync.Mutex"
Cohesion: 0.12
Nodes (13): Target, RealFileEngines(), FileEngines, Files, FileSpec, FileTarget, guardRow(), newFileTarget() (+5 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "fileedit_symlink_test.go"
Cohesion: 0.40
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 73 - "time.Time"
Cohesion: 0.16
Nodes (9): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, sweepID, sweepMemo (+1 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "Network"
Cohesion: 0.11
Nodes (11): snapNet, netGuardPill(), fillNetwork(), snapEntry(), clonePolicy(), NetEngines, NetExtra, NetPolicy (+3 more)

### Community 76 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 77 - "icon"
Cohesion: 0.14
Nodes (8): filesPage, specFromForm(), newTargetView(), specSummary(), asset(), coloredIcon(), icon(), EventTypes()

### Community 78 - "NetEventType"
Cohesion: 0.08
Nodes (12): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, newNetEventChecks(), FormatAddr(), NetEventTypes(), Ntohs() (+4 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.05
Nodes (14): holdHeadless(), HeadlessLine(), runGuardHeadless(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestDispatchWithoutOwnAllowedEvents(), TestLocalizePaths() (+6 more)

### Community 81 - ".decryptFileWithProgress"
Cohesion: 0.21
Nodes (5): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata()

### Community 82 - "newTestTarget"
Cohesion: 0.22
Nodes (16): newExeResolver(), NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+8 more)

### Community 83 - "wizard.go"
Cohesion: 0.10
Nodes (28): runMaintenanceMode(), deletePostBackups(), restoreBackups(), askDecrypt(), decryptDirectories(), decryptStep(), offerBackupCleanup(), pickDirsToDecrypt() (+20 more)

### Community 84 - "Monitor"
Cohesion: 0.15
Nodes (3): NewMonitor(), Monitor, watchedInode

### Community 85 - "runNetworkGuard"
Cohesion: 0.11
Nodes (22): ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), WithSuspectsFrom() (+14 more)

### Community 86 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.19
Nodes (6): viewportState, moveTo(), Finder, baseStyle(), viewportOf(), viewportView()

### Community 87 - "liveSession"
Cohesion: 0.26
Nodes (6): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), displayEntries(), holdForward(), liveSession

### Community 88 - "highlight_test.go"
Cohesion: 0.11
Nodes (28): nextConfig(), New(), SetText(), TestCursorDoesNotBlink(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines(), fg(), lexerFor() (+20 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "GlobReservations"
Cohesion: 0.17
Nodes (10): GlobChild, statFunc, compileGlobNames(), globText(), GlobReservations, resolveBits(), resolveChildren(), statKey() (+2 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.10
Nodes (30): UIDResolver, NewUIDResolver(), drainEphemeral(), runHeadless(), writeEvent(), admitEvent(), foldable(), newGateLogLimiter() (+22 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.09
Nodes (12): reasonText(), eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs() (+4 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 100 - "policyEditor"
Cohesion: 0.11
Nodes (15): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), eventsOrAll(), netWarnings() (+7 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "auditGroupCoverage"
Cohesion: 0.16
Nodes (15): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+7 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "time.Duration"
Cohesion: 0.20
Nodes (6): TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), debouncer

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "widgets.go"
Cohesion: 0.07
Nodes (27): binaryRow, clickPanel, counterStrip, homePage, iconTile, panel, pill, statusTile (+19 more)

### Community 110 - "KernelDev"
Cohesion: 0.11
Nodes (25): BinaryStat, binaryStatOf(), KernelDev(), LstatInode(), StatBinary(), StatBinaryFile(), BtrfsDevices(), btrfsSysfsDevices() (+17 more)

### Community 111 - "mountTable"
Cohesion: 0.27
Nodes (7): MountInfo, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), mountTable

### Community 112 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.24
Nodes (11): NewGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+3 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 116 - "RefreshCatalog"
Cohesion: 0.25
Nodes (9): updateCatalogConfig(), RefreshCatalog(), staleDiscordLibConf(), TestRefreshCatalogAdmitVetsOnlyNewLines(), TestRefreshCatalogKeepsLibLinesOfLockedLookingVault(), TestRefreshCatalogKeepsSymlinkLineNotTarget(), TestRefreshCatalogPrunesVanishedAppVersionLibLines(), TestRefreshCatalogReportsVersionBump() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "runEditProtected"
Cohesion: 0.19
Nodes (14): TestValidatePassword(), ValidatePassword(), LiveModeAvailable(), runEditProtected(), pingInterval(), chooseResource(), promptPassword(), runLiveEdit() (+6 more)

### Community 121 - "New"
Cohesion: 0.10
Nodes (41): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault() (+33 more)

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
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 127 - "CandidateDir"
Cohesion: 0.04
Nodes (53): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), systemPatterns(), systemRule(), entryWriters(), reserveChain(), symlinkTargetGlob() (+45 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "vaultFS"
Cohesion: 0.18
Nodes (10): TestRuntimeClassMarkers(), applyNewFileMeta(), defaultNewFileMode(), openRootT(), TestDefaultNewFileMode(), TestWriteFileKeepMeta(), TestWriteFileKeepMetaRefusesSymlink(), writeFileInPlace() (+2 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "TempBinary"
Cohesion: 0.09
Nodes (21): binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps(), planTempBlock(), ResolveTempBinary(), stripPinnedAllow() (+13 more)

### Community 143 - "io.Reader"
Cohesion: 0.33
Nodes (4): BtrfsMounted(), mountinfoHasFstype(), TestMountinfoHasFstype(), progressReader

### Community 144 - "unit.go"
Cohesion: 0.16
Nodes (17): confirmUnit(), TestUnifiedDiff(), TestUnifiedDiffIdentical(), UnifiedDiff(), DaemonUnitWithMetadataOutput(), ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate() (+9 more)

### Community 145 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 146 - "SuperblockDev"
Cohesion: 0.29
Nodes (7): FirstOnBtrfs(), mountinfoDev(), OnBtrfs(), SuperblockDev(), SuperblockDevPath(), TestMountinfoDev(), TestSuperblockDevMatchesStatOffBtrfs()

### Community 147 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 148 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 149 - "newTestVetter"
Cohesion: 0.19
Nodes (18): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+10 more)

### Community 150 - "Vault"
Cohesion: 0.19
Nodes (11): isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn(), readKey(), TestNewBoundedKeyFnFirstCall() (+3 more)

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.29
Nodes (19): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+11 more)

### Community 154 - "EventType"
Cohesion: 0.10
Nodes (13): snapEngines, snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting() (+5 more)

### Community 155 - "main_test.go"
Cohesion: 0.10
Nodes (7): IntegrationSuite, IntegrationSuite, infraContainerPaths(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

### Community 156 - "watchPattern"
Cohesion: 0.36
Nodes (7): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchPattern, Process inspectors

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "GuardInodeKey"
Cohesion: 0.04
Nodes (63): setConfinedHomes(), refusedReplacements, VettedInode, Multicall, MulticallKind, supersedeUnnamed(), appletTag(), buildOf() (+55 more)

### Community 161 - "Pulse"
Cohesion: 0.28
Nodes (7): New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse, Sub

### Community 162 - "Highlighter"
Cohesion: 0.14
Nodes (4): diffModel, editorModel, Highlighter, changelogModel

### Community 163 - "lockAndDeprovision"
Cohesion: 0.21
Nodes (8): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), applyRawKeyPolicy(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision(), modifiedContextWithSource()

### Community 164 - "openBunTmp"
Cohesion: 0.27
Nodes (11): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+3 more)

### Community 165 - "feed.go"
Cohesion: 0.08
Nodes (39): actionBar, cell, column, dataTable, eventFeed, feedAction, feedSpec, processView (+31 more)

### Community 166 - "shell"
Cohesion: 0.08
Nodes (26): Options, pageDef, pageID, shell, sysStatus, newLogo(), newShell(), Run() (+18 more)

### Community 167 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 168 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.16
Nodes (4): harnessSuite, mountinfoUnitTest, unescapeMountinfo(), editorHarness

### Community 169 - "Execute"
Cohesion: 0.40
Nodes (5): Requested(), TestRequested(), Execute(), initLogger(), main()

### Community 170 - "find_test.go"
Cohesion: 0.48
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 174 - "newWatchPattern"
Cohesion: 0.31
Nodes (11): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+3 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 178 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 179 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 184 - "lockRootRecovering"
Cohesion: 0.29
Nodes (8): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), PinPrefix(), sanitizeGen(), SharedPinPrefix(), TestPinPrefixDistinctAndStable(), TestPinPrefixShape()

### Community 185 - "stripTempJournal"
Cohesion: 0.32
Nodes (8): ensurePinStateFilePlaceholder(), ensurePlaceholder(), bootstrapSelfProtectPlaceholders(), ensureHashFilePlaceholder(), readTempJournal(), stripTempJournal(), TestTempJournalRoundTripInPlace(), writeTempJournal()

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

### Community 194 - "dialLiveSession"
Cohesion: 0.38
Nodes (7): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig(), runEditConfig()

### Community 195 - "runForward"
Cohesion: 0.40
Nodes (6): chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), normalizedEvents(), runForward()

### Community 196 - "TestBatchInodeRows"
Cohesion: 0.40
Nodes (6): inodeBatchFn, fakeBatches(), TestBatchInodeRows(), countInodeRows(), InodeOccupancy(), TestCountInodeRows()

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 204 - "runLockdown"
Cohesion: 0.20
Nodes (11): configureDaemonLogging(), confirmMasterKeyOverwrite(), resolveConfigPath(), runGenKey(), runLockdown(), runOneShotMode(), dirIsRootOwnedSafe(), isBpffs() (+3 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.25
Nodes (8): Actions, Desktop GUI, Files, Home, Install, Network, The daemon binary, Workspaces

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 515 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **31 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Process inspectors` connect `watchPattern` to `Daemon configuration`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12890152861713472 - nodes in this community are weakly interconnected._
- **Why does `Daemon configuration` connect `Daemon configuration` to `watchPattern`, `README.md`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02620802620802621 - nodes in this community are weakly interconnected._
- **Why does `watchPattern` connect `watchPattern` to `time.Time`, `go_pkg_path_filepath`, `newWatchPattern`, `catalogRefresher`?**
  _High betweenness centrality (0.048) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08958130477117819 - nodes in this community are weakly interconnected._