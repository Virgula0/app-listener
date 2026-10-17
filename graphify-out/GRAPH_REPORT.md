# Graph Report - app-listener  (2026-10-17)

## Corpus Check
- 441 files · ~1,620,935 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 5624 nodes · 23874 edges · 207 communities (172 shown, 35 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2133 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `dd3a277d`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- os.FileMode
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- rcUser
- guardUnitTest
- mkdirs
- codeedit.go
- auditGroupCoverage
- install_ui.go
- github.com/charmbracelet/bubbletea.Cmd
- go_pkg_github_com_testcontainers_testcontainers_go
- Network
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- runNetworkGuard
- guard.bpf.c
- loadDaemonConfig
- vaultFS
- icon
- github.com/charmbracelet/bubbles/textarea.Model
- NetGuardEvent
- Serve
- bpfstats/main.go
- render.go
- picker
- github.com/spf13/cobra.Command
- time.Duration
- Guard
- controlServer
- fileedit_model_test.go
- GuardInodeKey
- IntegrationSuite
- absPath
- monitorUnitTest
- hashFileShared
- update_test.go
- Ledger
- stdio
- liveSession
- catalogRefresher
- buildTrustedSet
- GenerateConf
- fileedit_symlink_test.go
- watchGroup
- procstats.go
- TrustGuard
- DeployInstalledBinary
- applyDiffAdditions
- engine
- Vault
- buildOneGuard
- .Refresh
- exe_supersede.h
- singleBinaryPaths
- runView
- pinstate.go
- gui/session.go
- runMonitor
- netModel
- FileTarget
- fcntl
- copyRegularAt
- BinaryEntry
- time.Time
- btrfsLayoutFrom
- Exe
- bufio.Reader
- daemon_multicall_test.go
- NetEventType
- process_vm_readv.c
- GuardEvent
- fakeGuardRepo
- planUpdaters
- RunWith
- filevault_test.go
- Config
- Finder
- IntegrationSuite
- highlight_test.go
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- find_test.go
- inodeChain
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- User
- policyEditor
- NewDumpHook
- Well-known bypass classes
- Execute
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- RefreshCatalog
- Build & sign release assets (reusable workflow)
- installPage
- KernelDev
- mountTable
- NetworkGuardRepository
- app-listener Terminal Demo Recording
- IntegrationSuite
- IntegrationSuite
- io.Writer
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- DaemonEvent
- IntegrationSuite
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- expandPlaceholders
- make test-integration (rootful Docker suite)
- writeWithin
- supersede_probe.c
- newTestVetter
- trace-app-libs.sh
- IntegrationSuite
- configedit_test.go
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TempBinary
- io.Reader
- options.go
- fakeVault
- openBunTmp
- net_tester/main.go
- edit-protected
- Monitor
- wipe
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- FileEvent
- layout.go
- runNetworkMonitor
- newTestTarget
- Installation
- StatInode
- graphify-refresh.sh
- unixModeToFileMode
- changelog.go
- main_test.go
- NetworkMonitorRepository
- feed.go
- shell
- hashCache
- formatEntry
- golang.org/x/sys/unix.Stat_t
- IntegrationSuite
- IntegrationSuite
- Limitations
- Troubleshooting
- WithBottomBar
- Vault
- Daemon configuration
- sync.Mutex
- runEditConfig
- Compatibility
- How it works
- tempGrantSetup
- SanitizeText
- bleedLayout
- golang.org/x/sys/unix.Timespec
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- .runEditorHarness
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- mountsUnitTest
- github.com/stretchr/testify/suite.Suite
- monitor: see what touches a path
- kernel-prof.sh
- pathCache
- .eventTypeAllowed
- ParseEventsFlag
- .TestCheckHardlinkByInode
- MulticallError
- Desktop GUI

## God Nodes (most connected - your core abstractions)
1. `Guard` - 117 edges
2. `Config` - 96 edges
3. `GuardInodeKey` - 96 edges
4. `Load()` - 90 edges
5. `IntegrationSuite` - 82 edges
6. `absPath()` - 81 edges
7. `fileEditModel` - 77 edges
8. `writeConfig()` - 69 edges
9. `IntegrationSuite` - 68 edges
10. `shell` - 66 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `runConfigEditor()` --calls--> `EditText()`  [INFERRED]
  cmd/functions/install/steps.go → internal/install/editor.go

## Import Cycles
- None detected.

## Communities (207 total, 35 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.10
Nodes (13): tempJournalRow, launchRule, trustEvent, btrfsIoctlFsInfoArgs, daemonStatsMsg, guardStatsMsg, hint, netGuardStatsMsg (+5 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (121): plantBun(), TestOpenBunTmp_KeepsVettedDir(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestDiffCatalogFlagRegistered() (+113 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "os.FileMode"
Cohesion: 0.09
Nodes (30): TestDaemonUnitWithMetadataOutput(), dupFD(), installFile(), installFileAs(), installServices(), installSSHAgent(), installSSHAgentEnv(), installUnit() (+22 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (77): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+69 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "rcUser"
Cohesion: 0.12
Nodes (36): setupBunTmpdir(), revertBunLaunchers(), rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), EnsureBunTmpDir() (+28 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.07
Nodes (14): addErrKind, classifyAddErr(), eventMask(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+6 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (35): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+27 more)

### Community 10 - "codeedit.go"
Cohesion: 0.13
Nodes (28): class, diffModel, editorModel, classOf(), column(), dedent(), Edit(), leadingIndent() (+20 more)

### Community 11 - "auditGroupCoverage"
Cohesion: 0.16
Nodes (15): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+7 more)

### Community 12 - "install_ui.go"
Cohesion: 0.08
Nodes (24): nextConfig(), installUI, multiAnswer, textAnswer, cancelUI, scrollable(), ask(), captureLogs() (+16 more)

### Community 13 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.08
Nodes (10): NewDaemonModel(), clamp(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), DaemonResourceInfo, noticeModel, servedModel (+2 more)

### Community 15 - "Network"
Cohesion: 0.16
Nodes (5): netGuardPill(), clonePolicy(), NetExtra, NetPolicy, Network

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (57): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+49 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.08
Nodes (14): guardBinaryFlag(), netMonitorEvent, IntegrationSuite, IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent() (+6 more)

### Community 19 - "fileEditModel"
Cohesion: 0.10
Nodes (9): clipLabel(), humanSize(), paneTitle(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind (+1 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.19
Nodes (45): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+37 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (14): ownerPaths(), backingDeviceUnion(), Resource, BackingDevices(), GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase (+6 more)

### Community 23 - "auth.go"
Cohesion: 0.12
Nodes (28): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+20 more)

### Community 24 - "runNetworkGuard"
Cohesion: 0.11
Nodes (19): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), DiscoverInfraBinaries(), Mode (+11 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "loadDaemonConfig"
Cohesion: 0.22
Nodes (9): loadDaemonConfig(), resolveConfigPath(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds(), SetUserHomes() (+1 more)

### Community 27 - "vaultFS"
Cohesion: 0.14
Nodes (9): TestRuntimeClassMarkers(), applyNewFileMeta(), defaultNewFileMode(), isBinaryContent(), listEntries(), TestDefaultNewFileMode(), TestIsBinaryContent(), writeFileInPlace() (+1 more)

### Community 28 - "icon"
Cohesion: 0.08
Nodes (18): binaryList, counterStrip, filesPage, targetView, newBinaryList(), newBinaryRow(), specFromForm(), newCounterStrip() (+10 more)

### Community 29 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.15
Nodes (18): pos, repeatKey(), reposition(), viewportOf(), cellToRune(), clipboardText(), cursorPos(), Selection (+10 more)

### Community 30 - "NetGuardEvent"
Cohesion: 0.11
Nodes (8): snapNetGuard, NetGuardEvent, fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 31 - "Serve"
Cohesion: 0.05
Nodes (41): New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), basicAuth(), configureWebSocket(), decodeEndsAtEOF() (+33 more)

### Community 32 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 33 - "render.go"
Cohesion: 0.12
Nodes (12): lineStyles, viewportState, fg(), baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+4 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (29): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+21 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (18): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+10 more)

### Community 36 - "time.Duration"
Cohesion: 0.15
Nodes (8): procStats, statsSample, IntegrationSuite, pooledContainer, cpuTime(), humanBytes(), rssBytes(), Ingested()

### Community 37 - "Guard"
Cohesion: 0.07
Nodes (12): systemPatterns(), systemRule(), rootHandle, BinaryRule, openRoot(), SharedPinDegraded(), TestNoteRootMountFollowsMountID(), deleteResKeys() (+4 more)

### Community 38 - "controlServer"
Cohesion: 0.19
Nodes (8): readClientLines(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), controlServer, editControlSession, grantRequest

### Community 39 - "fileedit_model_test.go"
Cohesion: 0.19
Nodes (28): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+20 more)

### Community 40 - "GuardInodeKey"
Cohesion: 0.05
Nodes (49): checkFreshDir(), openDirNoFollow(), renewDir(), refusedReplacements, VettedInode, appendInodeRow(), batchInodeRows(), inodeRowsOf() (+41 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (5): IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "hashFileShared"
Cohesion: 0.21
Nodes (9): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+1 more)

### Community 45 - "update_test.go"
Cohesion: 0.06
Nodes (44): assetsFor(), checkUpdatePreconditions(), downloadAndVerify(), downloadFile(), downloadReleaseFiles(), FetchLatest(), fetchRelease(), fetchReleases() (+36 more)

### Community 46 - "Ledger"
Cohesion: 0.14
Nodes (4): newBinaryVetter(), Ledger, Pending, Source

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "liveSession"
Cohesion: 0.13
Nodes (19): dialLiveSession(), parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), putConfig(), chooseResources(), displayEntries(), eventsLabel() (+11 more)

### Community 49 - "catalogRefresher"
Cohesion: 0.12
Nodes (12): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), awaitStartupOrSignal(), TestAwaitStartupOrSignalAbortsAndLocksDown() (+4 more)

### Community 50 - "buildTrustedSet"
Cohesion: 0.05
Nodes (45): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+37 more)

### Community 51 - "GenerateConf"
Cohesion: 0.05
Nodes (64): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), sectionsFromCandidates(), discordCatalogEntry(), groupedDiscordConf(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), TestGeneratedLibraryBlockRoundTrips() (+56 more)

### Community 52 - "fileedit_symlink_test.go"
Cohesion: 0.27
Nodes (14): openRootT(), TestWriteFileKeepMeta(), TestWriteFileKeepMetaRefusesSymlink(), mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink() (+6 more)

### Community 53 - "watchGroup"
Cohesion: 0.11
Nodes (34): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+26 more)

### Community 54 - "procstats.go"
Cohesion: 0.07
Nodes (30): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+22 more)

### Community 55 - "TrustGuard"
Cohesion: 0.05
Nodes (24): exeRow, tempGrant, tempMaskOp, tempRow, trustHook, trustRows, EventType, parseEvents() (+16 more)

### Community 56 - "DeployInstalledBinary"
Cohesion: 0.05
Nodes (59): checkAndApply(), reloadDaemonForPasswordChange(), applyLiveRefresh(), buildBinaryIfNeeded(), installBinaryOnly(), preflightDeployedBinary(), runUpdateCatalogOnly(), softenAutomatedRefreshErr() (+51 more)

### Community 57 - "applyDiffAdditions"
Cohesion: 0.06
Nodes (41): CheckBPFLSM(), confirmMasterKeyOverwrite(), runBPFCheck(), runGenKey(), runOneShotMode(), askBunTmpdirUsers(), gatherPerUserSetup(), applyDiffAdditions() (+33 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (32): engine, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded() (+24 more)

### Community 59 - "Vault"
Cohesion: 0.08
Nodes (38): buildConcurrency(), buildGuards(), lockOneRoot(), lockRootRecovering(), newPinGeneration(), prepareDaemonStart(), reloadOnce(), reloadSlotsNeeded() (+30 more)

### Community 60 - "buildOneGuard"
Cohesion: 0.09
Nodes (37): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards(), resolveGuardConfig() (+29 more)

### Community 61 - ".Refresh"
Cohesion: 0.04
Nodes (12): appTheme, clickPanel, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, newNavItem() (+4 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.14
Nodes (20): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+12 more)

### Community 63 - "singleBinaryPaths"
Cohesion: 0.11
Nodes (17): GlobChild, statFunc, Multicall, MulticallKind, multicallFamily(), singleBinaryPaths(), compileGlobNames(), globText() (+9 more)

### Community 64 - "runView"
Cohesion: 0.12
Nodes (6): logHook, logLine, runView, stageItem, stageState, newRunView()

### Community 65 - "pinstate.go"
Cohesion: 0.12
Nodes (30): configureDaemonLogging(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode() (+22 more)

### Community 66 - "gui/session.go"
Cohesion: 0.07
Nodes (41): invokingIDs(), ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField() (+33 more)

### Community 67 - "runMonitor"
Cohesion: 0.10
Nodes (21): BuildEBPFTargets(), CheckEBPF(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair() (+13 more)

### Community 68 - "netModel"
Cohesion: 0.21
Nodes (8): formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg, netModel

### Community 69 - "FileTarget"
Cohesion: 0.15
Nodes (9): guardPill(), Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec(), overlaps() (+1 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 72 - "BinaryEntry"
Cohesion: 0.07
Nodes (22): deferredBinary, BinaryEntry, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), checkBtrfsKeys(), BinariesSummary(), canonicalBinaryPath() (+14 more)

### Community 73 - "time.Time"
Cohesion: 0.09
Nodes (22): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+14 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "Exe"
Cohesion: 0.07
Nodes (22): sysStatus, probeSystem(), Check(), collectDiagnostics(), readSysctl(), requireKernel51(), DaemonRunning(), systemctlActive() (+14 more)

### Community 76 - "bufio.Reader"
Cohesion: 0.20
Nodes (11): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest(), TestReadGrantRequestForward() (+3 more)

### Community 77 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 78 - "NetEventType"
Cohesion: 0.07
Nodes (19): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, eventsOrAll(), newNetEventChecks(), Cstr(), FormatAddr() (+11 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.06
Nodes (23): holdHeadless(), HeadlessLine(), runGuardHeadless(), bpfGuardEvent, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestDispatchWithoutOwnAllowedEvents() (+15 more)

### Community 82 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 83 - "RunWith"
Cohesion: 0.06
Nodes (46): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary() (+38 more)

### Community 84 - "filevault_test.go"
Cohesion: 0.09
Nodes (45): harnessSuite, classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD() (+37 more)

### Community 85 - "Config"
Cohesion: 0.09
Nodes (35): ServeConfig, openBinaryVetter(), catchLifecycleSignals(), makeReloadHandler(), markReady(), newDaemonModel(), noteReloadedSettings(), notifySystemdReady() (+27 more)

### Community 86 - "Finder"
Cohesion: 0.17
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 88 - "highlight_test.go"
Cohesion: 0.16
Nodes (16): Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens(), TestGutterKeepsTextColumnAligned() (+8 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "find_test.go"
Cohesion: 0.48
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 92 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.10
Nodes (11): foreignPinnedLink(), reasonText(), eventTypeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), TestReasonText() (+3 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "User"
Cohesion: 0.07
Nodes (48): entryWriters(), reserveChain(), TestGlobBuilder_LibBitsPerEntry(), underResource(), bunEntryConfigured(), bunLaunchersFor(), cleanOrphanedFscrypt(), addManualDirectories() (+40 more)

### Community 100 - "policyEditor"
Cohesion: 0.12
Nodes (13): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), newNetGuardTab() (+5 more)

### Community 101 - "NewDumpHook"
Cohesion: 0.20
Nodes (12): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+4 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "Execute"
Cohesion: 0.18
Nodes (10): installHooks(), releaseBuild(), Requested(), Run(), TestRequested(), TestSubcommandOutranksGUI(), Execute(), initLogger() (+2 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "RefreshCatalog"
Cohesion: 0.14
Nodes (17): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+9 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "installPage"
Cohesion: 0.08
Nodes (22): binaryRow, binSource, handButton, homePage, iconTile, installedBinary, installPage, panel (+14 more)

### Community 110 - "KernelDev"
Cohesion: 0.08
Nodes (35): fdStat, BinaryStat, physicalParent(), statFD(), binaryStatOf(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding() (+27 more)

### Community 111 - "mountTable"
Cohesion: 0.24
Nodes (7): MountInfo, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), mountTable

### Community 112 - "NetworkGuardRepository"
Cohesion: 0.12
Nodes (15): NetworkGuardRepository, NewGuardUseCase(), NetworkGuardUseCase, NewNetworkGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo() (+7 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 116 - "io.Writer"
Cohesion: 0.27
Nodes (8): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), fyneLogFilter, candidate

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "DaemonEvent"
Cohesion: 0.14
Nodes (25): UIDResolver, NewUIDResolver(), drainEphemeral(), runHeadless(), writeEvent(), admitEvent(), foldable(), newGateLogLimiter() (+17 more)

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

### Community 127 - "expandPlaceholders"
Cohesion: 0.07
Nodes (26): symlinkTargetGlob(), confSafeMatch(), expandPlaceholders(), BinaryRule, homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestExpandWhitelist() (+18 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "newTestVetter"
Cohesion: 0.15
Nodes (19): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+11 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "TempBinary"
Cohesion: 0.09
Nodes (16): fakeDaemon, TempBinary, TemporaryGrant, TempRule, applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant() (+8 more)

### Community 143 - "io.Reader"
Cohesion: 0.20
Nodes (8): stagedBinary, stageLocal(), TestStageLocal(), TestStageLocalRefusedIsRemoved(), BtrfsMounted(), mountinfoHasFstype(), TestMountinfoHasFstype(), progressReader

### Community 144 - "options.go"
Cohesion: 0.31
Nodes (6): installAction, InstallMode, InstallRequest, Options, InstallHooks, Run()

### Community 146 - "openBunTmp"
Cohesion: 0.53
Nodes (6): fileID(), openBunTmp(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 147 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 148 - "edit-protected"
Cohesion: 0.29
Nodes (7): edit-protected, Editing the daemon config, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 149 - "Monitor"
Cohesion: 0.18
Nodes (3): NewMonitor(), Monitor, watchedInode

### Community 150 - "wipe"
Cohesion: 0.06
Nodes (42): deprovisionKind, migrationTarget, isFileVaultCiphertext(), isRegularFileTarget(), writeTempWithMetadata(), classifyDeprovision(), classifyDeprovisionErr(), hasEncryptionPolicy() (+34 more)

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.29
Nodes (19): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+11 more)

### Community 154 - "FileEvent"
Cohesion: 0.18
Nodes (3): snapMonitor, FileEvent, fakeMonitorRepo

### Community 155 - "layout.go"
Cohesion: 0.17
Nodes (12): adaptiveGrid, checkGrid, gapVBox, encodeTable(), TestAdaptiveGridDropsColumns(), TestEncodeTable(), cellMin(), newAdaptiveGrid() (+4 more)

### Community 156 - "runNetworkMonitor"
Cohesion: 0.18
Nodes (8): ParseNetEventsFlag(), runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NewNetModel(), NetworkMonitorUseCase

### Community 157 - "newTestTarget"
Cohesion: 0.20
Nodes (17): newExeResolver(), NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "StatInode"
Cohesion: 0.06
Nodes (42): confirmUserPlaced(), ReplacementCheck, appletTag(), ExeKey(), IsMulticallRefusal(), mcNameKey(), MulticallNameTags(), MulticallRefusal() (+34 more)

### Community 161 - "unixModeToFileMode"
Cohesion: 0.29
Nodes (6): fileModeToUnixMode(), isOctalMode(), TestFileModeToUnixModeRoundTrip(), TestIsOctalMode(), TestUnixModeToFileMode(), unixModeToFileMode()

### Community 162 - "changelog.go"
Cohesion: 0.11
Nodes (29): changelogText(), cumulativeChangelog(), newChangelogModel(), newReleaseWindow(), releasesSince(), showChangelog(), rel(), tags() (+21 more)

### Community 163 - "main_test.go"
Cohesion: 0.18
Nodes (5): infraContainerPaths(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

### Community 164 - "NetworkMonitorRepository"
Cohesion: 0.16
Nodes (7): WithSuspectsFrom(), NetworkMonitorRepository, RealNetEngines(), NetEngines, netguardHarness, Option, NetGuardConfig

### Community 165 - "feed.go"
Cohesion: 0.08
Nodes (33): actionBar, cell, column, dataTable, eventFeed, feedAction, feedSpec, networkPage (+25 more)

### Community 166 - "shell"
Cohesion: 0.07
Nodes (24): pageDef, pageID, resizeWatch, shell, sizedDialog, snapEngines, snapNet, newLogo() (+16 more)

### Community 167 - "hashCache"
Cohesion: 0.18
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 168 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 169 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 174 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 178 - "sync.Mutex"
Cohesion: 0.21
Nodes (4): newControlManager(), startControlManager(), configEditor, controlManager

### Community 179 - "runEditConfig"
Cohesion: 0.67
Nodes (4): applyEdit(), editAgain(), fetchConfig(), runEditConfig()

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 183 - "SanitizeText"
Cohesion: 0.16
Nodes (9): multicallAdmissible(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), binaryVetter, vetDecision, launchKey(), TestLaunchRulesFitTheKey() (+1 more)

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

### Community 195 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.13
Nodes (6): BpfEvent, eventUnitTest, mountinfoUnitTest, DecodeBpfEvent(), unescapeMountinfo(), editorHarness

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 200 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 201 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 206 - "Desktop GUI"
Cohesion: 0.20
Nodes (10): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Tables (+2 more)

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 518 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **35 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Editor keys` connect `liveSession` to `RunWith`, `edit-protected`?**
  _High betweenness centrality (0.059) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09803402979811338 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `liveSession`, `README.md`?**
  _High betweenness centrality (0.059) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.022476190476190476 - nodes in this community are weakly interconnected._
- **Why does `Delete()` connect `RunWith` to `liveSession`, `daemon/daemon.go`?**
  _High betweenness centrality (0.043) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08958130477117819 - nodes in this community are weakly interconnected._