# Graph Report - app-listener  (2026-11-01)

## Corpus Check
- 564 files · ~1,795,531 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7352 nodes · 33724 edges · 272 communities (221 shown, 51 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 3007 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `6879a2ac`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- go_pkg_strings
- testing.T
- IntegrationSuite
- op
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- GuardEvent
- time.Time
- Network
- runNetworkMonitor
- buildOneGuard
- guardModel
- absPath
- guard_trust.bpf.c
- navItem
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- FileTarget
- guard.bpf.c
- LibraryClosure
- github.com/charmbracelet/bubbletea.Cmd
- SanitizeText
- prepareDaemonStart
- highlight_test.go
- MonitorRepository
- fileedit_model_test.go
- validateConfigText
- picker
- secret_probe/main.go
- IntegrationSuite
- Message
- controlServer
- github.com/charmbracelet/bubbles/textarea.Model
- watchSet
- .showSized
- cursor.go
- monitorUnitTest
- Candidate
- .ExpandLibDirWriters
- Guard
- stdio
- journal
- fadeScroll
- main_test.go
- protectedPage
- Client
- installPage
- session
- .apply
- updatePage
- mutedText
- GuardInodeKey
- writeWithin
- RawTarget
- runView
- LaunchRule
- policyEditor
- unix_connect.c
- wm.go
- IntegrationSuite
- Recorder
- safeio.go
- Vault
- trustManager
- planUpdaters
- fcntl
- selection_test.go
- update_test.go
- codecheck/main.go
- recordFlow
- newTestVetter
- process_vm_readv.c
- os.File
- OpenConfined
- .filterEvent
- rcUser
- golang.org/x/sys/unix.Stat_t
- bpfstats/main.go
- exe_supersede.h
- loginServer
- gui/session.go
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- Monitor
- Vault
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- runWizard
- sync.Mutex
- Well-known bypass classes
- NetGuardEvent
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- User
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- Serve
- BinaryEntry
- UI
- runOfflineEdit
- IntegrationSuite
- ConfinedEntry
- RunCmd
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- NetworkMonitor
- IntegrationSuite
- WebSocket /ws client
- inet_connect.c
- install.sh
- hashCache
- classify
- make test-integration (rootful Docker suite)
- find
- logLine
- recoverFileInPlace
- trace-app-libs.sh
- IntegrationSuite
- KernelDev
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- daemon: permanent protection with encryption at rest
- serveWebSocket
- liveSession
- Finder
- Ledger
- EventType
- fakeGuardRepo
- dialLiveSession
- Daemon configuration
- fakeStore
- mc_tag.h
- IntegrationSuite
- rootHandle
- reloadOnce
- copyRegularAt
- dialog_test.go
- Installation
- backend_test.go
- graphify-refresh.sh
- layout.go
- NetEventType
- newTestTarget
- IntegrationSuite
- feed.go
- shell
- inodeChain
- plan_test.go
- vaultfs.go
- SelectOrphans
- batchInodeRows
- Limitations
- downloadFile
- parseDaemonEvents
- TrustGuard
- net_tester/main.go
- debouncer
- changelog.go
- run.go
- bufio.Reader
- IntegrationSuite
- tempGrantSetup
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- NewDumpHook
- Self-updating apps and the catalog
- fileedit_symlink_test.go
- os.FileInfo
- github.com/spf13/cobra.Command
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- filevault_test.go
- pathCache
- pinstate.go
- DaemonEvent
- decryptStep
- vaultfs_test.go
- sources
- monitor
- IntegrationSuite
- .populateWatchInodes
- uninstallAll
- .eventTypeAllowed
- BusClient
- WithBottomBar
- revertSystemFiles
- .awaitDaemonUpOrDead
- .TestRecord_GuardParity
- linkView
- runNetworkGuard
- IntegrationSuite
- IntegrationSuite
- SystemTrusted
- configedit_test.go
- StatInode
- app-listener documentation
- supersede_probe.c
- launchKey
- bleedLayout
- mountTable
- PlanDaemonUnitUpdate
- github.com/Virgula0/app-listener/tools/codecheck
- inode_dev.h
- runBackend
- maintainer/plan.go
- ValidateTargets
- AddServeFlags
- Troubleshooting
- update/update.go
- GenerateConf
- guardLSMHooks
- IntegrationSuite
- taint_victim.c
- ListUsers
- ReadInstallCert
- sources_test.go
- runGenKey
- logGrid
- Multicall
- bpfLsmListed
- linkFixture
- runMonitor
- monitor: see what touches a path
- TestCleanRel
- .exeKeys
- .netIntegrityContainer
- formatEntry
- network monitor: see what a program talks to
- parseReservedPorts
- parseKernelRelease
- Execute
- golang.org/x/sys/unix.Timespec
- fileID
- MulticallError
- .TestDaemon_Bypass_PrctlExeFileSwap
- Config
- NewEventFanout
- IntegrationSuite
- .runEditorHarness
- RunWith

## God Nodes (most connected - your core abstractions)
1. `Config` - 147 edges
2. `Guard` - 125 edges
3. `GuardInodeKey` - 112 edges
4. `Load()` - 97 edges
5. `shell` - 97 edges
6. `absPath()` - 93 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 75 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Modes` --references--> `monitor()`  [INFERRED]
  docs/README.md → integrationtests/secret_probe/main.go
- `Distributions` --references--> `monitor()`  [INFERRED]
  docs/compatibility.md → integrationtests/secret_probe/main.go
- `Files` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `Workspaces` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `eBPF at three levels` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go

## Import Cycles
- None detected.

## Communities (272 total, 51 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.09
Nodes (20): tempJournalRow, Progress, User, invalidPassword, recordTick, linkedFile, argvMatch, launchRule (+12 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (161): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity() (+153 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (8): configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote(), watchPathsInConfig()

### Community 3 - "op"
Cohesion: 0.17
Nodes (12): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+4 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (82): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+74 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.09
Nodes (51): prctl_gate(), guard_task_prctl(), admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader() (+43 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.07
Nodes (15): addErrKind, classifyAddErr(), addRecorder(), copySelf(), guardUnitTest, runTool(), walkEntries(), walkInodes() (+7 more)

### Community 9 - "mkdirs"
Cohesion: 0.12
Nodes (40): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), newWatchPattern(), systemBinaryPatterns(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped() (+32 more)

### Community 10 - "GuardEvent"
Cohesion: 0.06
Nodes (19): holdHeadless(), HeadlessLine(), bpfGuardEvent, snapEngines, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestDispatchWithoutOwnAllowedEvents() (+11 more)

### Community 11 - "time.Time"
Cohesion: 0.08
Nodes (14): LockedError, sweepID, sweepMemo, codeEntry, editView, handButton, Entry, TrustRow (+6 more)

### Community 12 - "Network"
Cohesion: 0.11
Nodes (9): netGuardPill(), NetworkGuardRepository, clonePolicy(), NetEngines, NetExtra, NetPolicy, Network, verifyHashes() (+1 more)

### Community 13 - "runNetworkMonitor"
Cohesion: 0.10
Nodes (13): ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString() (+5 more)

### Community 14 - "buildOneGuard"
Cohesion: 0.06
Nodes (43): DrainUntilSignal(), buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards() (+35 more)

### Community 15 - "guardModel"
Cohesion: 0.05
Nodes (39): fg(), markStyle(), Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace() (+31 more)

### Community 16 - "absPath"
Cohesion: 0.08
Nodes (10): attrOpCase, IntegrationSuite, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent (+2 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.08
Nodes (65): argv_env(), argv_matched(), argv_rule_at(), argv_word(), caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key() (+57 more)

### Community 18 - "navItem"
Cohesion: 0.03
Nodes (16): appTheme, binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panelRenderer (+8 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (51): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+43 more)

### Community 22 - "Resource"
Cohesion: 0.05
Nodes (11): fakeDaemon, Resource, GuardRepository, Vault, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots() (+3 more)

### Community 23 - "auth_test.go"
Cohesion: 0.13
Nodes (39): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+31 more)

### Community 24 - "FileTarget"
Cohesion: 0.11
Nodes (13): guardPill(), Files, FileSpec, FileTarget, Policy, guardRow(), newFileTarget(), NormalizeSpec() (+5 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.07
Nodes (95): exe_refused(), inode_dev(), add_inode_to_guard(), admit_port(), btrfs_copy_gate(), bus_key(), check_and_emit_args(), chmod_only_drops_write() (+87 more)

### Community 26 - "LibraryClosure"
Cohesion: 0.07
Nodes (27): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), IsScript(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+19 more)

### Community 27 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.07
Nodes (12): recordTickCmd(), recordModel, diffModel, editorModel, renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), noticeModel (+4 more)

### Community 28 - "SanitizeText"
Cohesion: 0.06
Nodes (30): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), scriptIdentity(), logRefreshChange(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly() (+22 more)

### Community 29 - "prepareDaemonStart"
Cohesion: 0.14
Nodes (16): CheckBPFLSM(), CheckEBPF(), configureDaemonLogging(), newPinGeneration(), prepareDaemonStart(), resolveConfigPath(), runBPFCheck(), runLockdown() (+8 more)

### Community 30 - "highlight_test.go"
Cohesion: 0.12
Nodes (26): New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines(), TestSetTextStartsAtFirstLine() (+18 more)

### Community 31 - "MonitorRepository"
Cohesion: 0.13
Nodes (13): MonitorRepository, NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+5 more)

### Community 32 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 33 - "validateConfigText"
Cohesion: 0.10
Nodes (17): encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption(), TestAskEncryptionSkipsNeedEncryptionFalse() (+9 more)

### Community 34 - "picker"
Cohesion: 0.10
Nodes (19): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, errText(), isFindShortcut() (+11 more)

### Community 35 - "secret_probe/main.go"
Cohesion: 0.22
Nodes (25): busPath(), connexec(), connFD(), createItem(), dial(), execInto(), failf(), fakebus() (+17 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (15): checkSecret(), checkTrustLine(), trustRows(), File, TrustItem, backend, backendOp, backendState (+7 more)

### Community 38 - "controlServer"
Cohesion: 0.17
Nodes (9): authSelfPeer(), readClientLines(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), configEditor, controlServer (+1 more)

### Community 39 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.09
Nodes (34): class, lineStyles, pos, viewportState, classOf(), column(), dedent(), Edit() (+26 more)

### Community 40 - "watchSet"
Cohesion: 0.19
Nodes (7): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, Process inspectors

### Community 41 - ".showSized"
Cohesion: 0.09
Nodes (12): binaryList, installUI, resizeWatch, sizedDialog, scrollable(), newBinaryList(), selectable(), newFadeScroll() (+4 more)

### Community 42 - "cursor.go"
Cohesion: 0.05
Nodes (25): buttonCursor, handCheckGroup, handRadioGroup, handSelect, tabBar, tabButton, tabButtonRenderer, tabItem (+17 more)

### Community 44 - "Candidate"
Cohesion: 0.06
Nodes (51): systemPatterns(), systemRule(), nextConfig(), appendSectionsAndEdit(), collectDiffAdditions(), TestDiffMergePreservesExistingAndParses(), TestUncoveredCandidates(), uncoveredCandidates() (+43 more)

### Community 45 - ".ExpandLibDirWriters"
Cohesion: 0.14
Nodes (13): homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelistSymlinkEscape(), TestSteamLibDirsCoverProtonCodeOnly(), TestSteamLibDirsExcludeScoutRuntime(), TestSteamLibDirWritersCoverItsOwnBinaries(), boundedGlob() (+5 more)

### Community 46 - "Guard"
Cohesion: 0.06
Nodes (19): binaryVerifyState, deferredBinary, exeRow, supersededBinary, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), SharedPinDegraded(), confinedEntry(), deleteResKeys() (+11 more)

### Community 48 - "journal"
Cohesion: 0.14
Nodes (6): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 49 - "fadeScroll"
Cohesion: 0.07
Nodes (15): fadeBar, fadeBarRenderer, fadeScroll, fadeScrollRenderer, fadeScrollOf(), newTestFadeScroll(), settle(), TestFadeBarDragEndsGlide() (+7 more)

### Community 50 - "main_test.go"
Cohesion: 0.06
Nodes (19): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest, IntegrationSuite (+11 more)

### Community 51 - "protectedPage"
Cohesion: 0.05
Nodes (20): backendLink, configView, forwardForm, forwardView, protectedPage, recordReview, recordView, trustPage (+12 more)

### Community 52 - "Client"
Cohesion: 0.12
Nodes (10): Client, childEnv(), launch(), Start(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines(), Writer (+2 more)

### Community 53 - "installPage"
Cohesion: 0.09
Nodes (17): actionGate, binSource, installAction, InstallMode, installPage, InstallRequest, Options, stagedBinary (+9 more)

### Community 54 - "session"
Cohesion: 0.10
Nodes (9): peerCred(), TestReadFrame(), decodeFrame(), encodeWithSerial(), header(), readFrame(), queued, senderQueue (+1 more)

### Community 55 - ".apply"
Cohesion: 0.11
Nodes (14): TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan(), TestPlanAdmitsOnlyWhatWasRecorded(), TestPlanExtendsAnExistingSection(), testRecording() (+6 more)

### Community 56 - "updatePage"
Cohesion: 0.16
Nodes (12): anyDaemonRunning(), installedBinary, latestRelease, updatePage, hashFile(), hashReader(), readInstalled(), newUpdatePage() (+4 more)

### Community 57 - "mutedText"
Cohesion: 0.06
Nodes (48): checkGrid, filesPage, handCheck, homePage, multiAnswer, panel, recordForm, statusTile (+40 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (46): engine, portOwner, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel() (+38 more)

### Community 59 - "writeWithin"
Cohesion: 0.25
Nodes (8): atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync(), writeWithin()

### Community 60 - "RawTarget"
Cohesion: 0.27
Nodes (10): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+2 more)

### Community 62 - "runView"
Cohesion: 0.09
Nodes (7): daemonPage, runHost, runView, stageItem, stageState, captureLogs(), orDefault()

### Community 63 - "LaunchRule"
Cohesion: 0.13
Nodes (18): allowedEnv(), binaryLaunchRules(), launchRulesByPath(), TestAllowedEnvKeepsOnlyRootPlaced(), TestLaunchRulesByPath(), argvEntry, envName(), LaunchRule (+10 more)

### Community 64 - "policyEditor"
Cohesion: 0.13
Nodes (12): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), newNetGuardTab() (+4 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "wm.go"
Cohesion: 0.11
Nodes (25): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+17 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.23
Nodes (5): guardBinaryFlag(), IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent()

### Community 68 - "Recorder"
Cohesion: 0.14
Nodes (16): drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start(), record() (+8 more)

### Community 69 - "safeio.go"
Cohesion: 0.21
Nodes (15): readAtNoFollow(), upsertUnitAt(), writeAtNoFollow(), invokingIDs(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), createOwned(), AtomicWriteAt() (+7 more)

### Community 70 - "Vault"
Cohesion: 0.11
Nodes (23): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr() (+15 more)

### Community 71 - "trustManager"
Cohesion: 0.11
Nodes (21): inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree() (+13 more)

### Community 72 - "planUpdaters"
Cohesion: 0.08
Nodes (27): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+19 more)

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "selection_test.go"
Cohesion: 0.25
Nodes (10): clipboardText(), setClipboard(), cellOf(), repeat(), selKey(), TestSelectionEditing(), TestSelectionLargeDelete(), TestSelectionMouseDrag() (+2 more)

### Community 75 - "update_test.go"
Cohesion: 0.07
Nodes (25): fetchReleases(), parseChecksum(), rsaPublicKeyPEM(), signedChecksum(), TestAssetsFor(), TestFetchGUIWithoutAsset(), TestFetchReleases(), TestFetchReleasesHTTPError() (+17 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "recordFlow"
Cohesion: 0.11
Nodes (22): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), lines(), offerBackupDeletion(), pickBinaries() (+14 more)

### Community 78 - "newTestVetter"
Cohesion: 0.44
Nodes (10): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+2 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.21
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "os.File"
Cohesion: 0.06
Nodes (44): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+36 more)

### Community 81 - "OpenConfined"
Cohesion: 0.08
Nodes (39): swapFile(), loaderName(), LoaderRefusal(), soNames(), hostLoader(), TestLoaderRefusal(), MulticallRefusal(), RefuseIdentityLines() (+31 more)

### Community 83 - "rcUser"
Cohesion: 0.07
Nodes (65): installFacultatives(), installSSHAgentFacultative(), listFacultatives(), runFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative(), TestListFacultativesIgnoresDaemonConf(), TestListFacultativesMarksGuardedSSH() (+57 more)

### Community 84 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.09
Nodes (28): placedWalk, treeJudge, FirstOnBtrfs(), mountinfoDev(), OnBtrfs(), SuperblockDev(), TestMountinfoDev(), TestSuperblockDevMatchesStatOffBtrfs() (+20 more)

### Community 85 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 86 - "exe_supersede.h"
Cohesion: 0.22
Nodes (14): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+6 more)

### Community 87 - "loginServer"
Cohesion: 0.06
Nodes (30): markReady(), markReadyWithLogin(), notifySystemdReady(), hasSentinel(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel(), runtimePathsFor() (+22 more)

### Community 88 - "gui/session.go"
Cohesion: 0.20
Nodes (15): ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+7 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "Monitor"
Cohesion: 0.18
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 92 - "Vault"
Cohesion: 0.24
Nodes (3): Vault, Open(), TestOpenVaultRefusesSymlink()

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

### Community 95 - "NetGuard"
Cohesion: 0.09
Nodes (17): foreignPinnedLink(), eventTypeKey(), codeVouched(), isTraced(), judgePasses(), loaderEnv(), parseProcs(), roFor() (+9 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "time.Duration"
Cohesion: 0.14
Nodes (14): procStats, statsSample, duration(), cpuTime(), humanBytes(), rssBytes(), New(), recv() (+6 more)

### Community 100 - "runWizard"
Cohesion: 0.06
Nodes (50): askBunTmpdirUsers(), gatherPerUserSetup(), setupBunTmpdir(), addedResources(), applyDiffAdditions(), TestAddedResourcesSkipsAlreadyGuardedSSH(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled() (+42 more)

### Community 101 - "sync.Mutex"
Cohesion: 0.14
Nodes (8): placedKey(), fyneLogFilter, Caller, Gate, Role, Listen(), Bus, Judge

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "NetGuardEvent"
Cohesion: 0.14
Nodes (8): snapNetGuard, reasonText(), NetGuardEvent, TestReasonText(), protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "User"
Cohesion: 0.08
Nodes (28): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor() (+20 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.17
Nodes (9): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), btrfsTypes(), TestBtrfsLayoutFromNestedMembers() (+1 more)

### Community 110 - "Serve"
Cohesion: 0.17
Nodes (17): basicAuth(), isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), securityHeaders(), Serve() (+9 more)

### Community 111 - "BinaryEntry"
Cohesion: 0.13
Nodes (16): HashStore, keptEntries, sharedHashEntry, BinaryEntry, checkBtrfsKeys(), canonicalBinaryPath(), claimHash(), hashFileShared() (+8 more)

### Community 112 - "UI"
Cohesion: 0.09
Nodes (19): showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), cancelUI, TestEditTextMapsUICancel(), MultiSelectKeymap() (+11 more)

### Community 113 - "runOfflineEdit"
Cohesion: 0.10
Nodes (14): auditAfterEdit(), auditFindings(), auditFindingsWithConfig(), TestPickEncryptedDirSingleSkipsPicker(), chooseResource(), pickDir(), pickEncryptedDir(), runOfflineEdit() (+6 more)

### Community 115 - "ConfinedEntry"
Cohesion: 0.26
Nodes (20): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+12 more)

### Community 116 - "RunCmd"
Cohesion: 0.07
Nodes (41): offerDaemonStart(), runDiffCatalog(), installGUI(), buildBinaryIfNeeded(), certifyInstall(), installBinaryOnly(), markInstallComplete(), preflightDeployedBinary() (+33 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "NetworkMonitor"
Cohesion: 0.12
Nodes (8): BinariesSummary(), TestBinariesSummary(), PumpRingbuf(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "inet_connect.c"
Cohesion: 0.26
Nodes (15): attack_connect(), attack_fastopen(), attack_uring(), denied(), dial(), fill_addr(), handoff(), listen_on() (+7 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "hashCache"
Cohesion: 0.12
Nodes (17): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, BinaryStat, DurableStat, HashCache (+9 more)

### Community 127 - "classify"
Cohesion: 0.22
Nodes (12): call(), TestClassify(), TestClassifyCreateItemReplace(), TestClassifyOpenCallsExact(), arg(), attributesOf(), classify(), classifyCollection() (+4 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "find"
Cohesion: 0.22
Nodes (16): deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup(), TestMovedFlagError(), movedFlagError(), Delete(), Find() (+8 more)

### Community 130 - "logLine"
Cohesion: 0.17
Nodes (7): logHook, logLine, journalLine(), TestJournalLine(), logRow(), logStyle(), DumpHook

### Community 131 - "recoverFileInPlace"
Cohesion: 0.21
Nodes (16): classifyRegularFileTarget(), clearRecovery(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink(), sealFileVaultWithMasterKey() (+8 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "KernelDev"
Cohesion: 0.31
Nodes (9): fdStat, physicalParent(), statFD(), KernelDev(), BtrfsDevices(), btrfsSysfsDevices(), TestBtrfsSysfsDevices(), readMountDevs() (+1 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.13
Nodes (15): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+7 more)

### Community 143 - "serveWebSocket"
Cohesion: 0.15
Nodes (12): configureWebSocket(), decodeEndsAtEOF(), parseResize(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestParseResize(), TestSameOrigin() (+4 more)

### Community 144 - "liveSession"
Cohesion: 0.16
Nodes (13): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), forwardCheck() (+5 more)

### Community 145 - "Finder"
Cohesion: 0.13
Nodes (19): findButton, span, findAll(), Finder, lowerRunes(), NewFinder(), newFindInput(), overlay() (+11 more)

### Community 146 - "Ledger"
Cohesion: 0.07
Nodes (25): candidates(), confirmAll(), describe(), record(), run(), applySchema(), Ledger, Pending (+17 more)

### Community 147 - "EventType"
Cohesion: 0.07
Nodes (20): selfEventsFor(), snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting() (+12 more)

### Community 148 - "fakeGuardRepo"
Cohesion: 0.06
Nodes (19): binaryStatOf(), TempBinary, TemporaryGrant, TempRule, ResolveTempBinary(), TestResolveTempBinaryClassesRuntime(), applyTempGrants(), closeTempBinaries() (+11 more)

### Community 149 - "dialLiveSession"
Cohesion: 0.13
Nodes (19): TestValidatePassword(), ValidatePassword(), dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig() (+11 more)

### Community 150 - "Daemon configuration"
Cohesion: 0.17
Nodes (12): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Launch rules: `[launch_rules]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Port sections, Reloading (+4 more)

### Community 152 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 154 - "rootHandle"
Cohesion: 0.23
Nodes (3): rootHandle, openRoot(), TestReadOnlyWalkToleratesDanglingSymlinks()

### Community 155 - "reloadOnce"
Cohesion: 0.07
Nodes (27): drain(), newCatalogRefresher(), readEvents(), sameConfig(), newControlManager(), startControlManager(), backingDeviceUnion(), buildGuards() (+19 more)

### Community 156 - "copyRegularAt"
Cohesion: 0.36
Nodes (9): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), fchownFromInfo(), futimesFromInfo(), readDirNames(), readlinkFD() (+1 more)

### Community 157 - "dialog_test.go"
Cohesion: 0.41
Nodes (11): openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder(), TestPickerSubfolderSearchPicksFile() (+3 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 161 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 162 - "NetEventType"
Cohesion: 0.12
Nodes (11): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), FormatAddr(), NetEventTypes(), Ntohs() (+3 more)

### Community 163 - "newTestTarget"
Cohesion: 0.22
Nodes (17): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (45): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+37 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (55): fakeLink, pageDef, pageID, reel, shell, snapNet, sysStatus, newShell() (+47 more)

### Community 167 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 168 - "plan_test.go"
Cohesion: 0.30
Nodes (12): fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent(), TestPathRecord(), TestScanAndValidateAcceptsLibDirContent(), TestScanRefusesSpecialFiles(), TestScanUntouchedLibDir(), TestTouches() (+4 more)

### Community 169 - "vaultfs.go"
Cohesion: 0.15
Nodes (13): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), IsBinaryContent(), IsOctalMode(), ParseMode(), TestDefaultNewFileMode(), TestFileModeToUnixModeRoundTrip() (+5 more)

### Community 170 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 171 - "batchInodeRows"
Cohesion: 0.22
Nodes (11): inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows(), walkInodeBatches(), countInodeRows() (+3 more)

### Community 172 - "Limitations"
Cohesion: 0.18
Nodes (11): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Secrets served over IPC (+3 more)

### Community 173 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 174 - "parseDaemonEvents"
Cohesion: 0.16
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite

### Community 175 - "TrustGuard"
Cohesion: 0.05
Nodes (30): tempGrant, tempMaskOp, tempRow, trustRows, ensureBtrfsLayout(), droppableFreed(), dropSupersededLocked(), freedSuperseded() (+22 more)

### Community 176 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 177 - "debouncer"
Cohesion: 0.33
Nodes (4): TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), debouncer

### Community 178 - "changelog.go"
Cohesion: 0.14
Nodes (27): changelogText(), cumulativeChangelog(), newChangelogModel(), newReleaseWindow(), releasesSince(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit() (+19 more)

### Community 179 - "run.go"
Cohesion: 0.42
Nodes (7): planFrom(), newStage(), Run(), runSandbox(), tail(), Spec, stage

### Community 180 - "bufio.Reader"
Cohesion: 0.11
Nodes (19): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+11 more)

### Community 181 - "IntegrationSuite"
Cohesion: 0.23
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 182 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 186 - "Binary trust"
Cohesion: 0.22
Nodes (9): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch rules, Launch scan, Multicall binaries, The dynamic loader (+1 more)

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "NewDumpHook"
Cohesion: 0.19
Nodes (15): TestSetupDumpLogRefusesSymlink(), applyVerbosity(), effectiveVerbose(), preRun(), setupDumpLog(), VerboseLevel, BridgeStdLog(), NewDumpHook() (+7 more)

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "fileedit_symlink_test.go"
Cohesion: 0.36
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 194 - "os.FileInfo"
Cohesion: 0.10
Nodes (16): auditEntry(), auditTree(), coveredByWatch(), relOrPath(), TestAuditTree(), TestWithin(), within(), migrationTarget (+8 more)

### Community 195 - "github.com/spf13/cobra.Command"
Cohesion: 0.22
Nodes (12): credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), asRoot(), runInstall() (+4 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 199 - "filevault_test.go"
Cohesion: 0.17
Nodes (26): deriveFileVaultSubkey(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault(), sealFileVault(), inode(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated() (+18 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "pinstate.go"
Cohesion: 0.16
Nodes (24): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+16 more)

### Community 202 - "DaemonEvent"
Cohesion: 0.08
Nodes (40): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless() (+32 more)

### Community 203 - "decryptStep"
Cohesion: 0.15
Nodes (13): askDecrypt(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), TestDecryptDirectoriesEmpty(), CleanOrphanedMetadata(), LivePolicyDescriptors() (+5 more)

### Community 204 - "vaultfs_test.go"
Cohesion: 0.28
Nodes (10): inodeOf(), mustRead(), mustWrite(), openRootT(), TestWriteFileInPlace(), TestWriteFileKeepMeta(), TestWriteFileKeepMetaRefusesPlantedTempSibling(), TestWriteFileKeepMetaRefusesSymlink() (+2 more)

### Community 205 - "sources"
Cohesion: 0.21
Nodes (5): TestGeneratedLibraryBlockRoundTrips(), lstatOf(), outsideHardlinks(), libDirDev(), sources

### Community 206 - "monitor"
Cohesion: 0.07
Nodes (33): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks (+25 more)

### Community 210 - "uninstallAll"
Cohesion: 0.20
Nodes (12): releaseBuild(), Run(), confirmKeyDeletion(), run(), runUninstall(), RunWith(), removeKeyAndEmptyDir(), removeMasterKey() (+4 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "BusClient"
Cohesion: 0.19
Nodes (8): BusClient, AttestBusClient(), busClientKey(), execStamp(), imageOf(), pidfdPid(), procIdentity(), suspectLabel()

### Community 213 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 214 - "revertSystemFiles"
Cohesion: 0.20
Nodes (12): removeSessionHooks(), revertSystemFiles(), RemoveLoginDropIn(), installFile(), removeFile(), SyncSecretServiceFiles(), LegacyCatalogRefreshFiles(), RemoveBinSymlink() (+4 more)

### Community 216 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 217 - "linkView"
Cohesion: 0.36
Nodes (4): linkTarget(), entryKind, libPlan, linkView

### Community 218 - "runNetworkGuard"
Cohesion: 0.11
Nodes (22): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), WithSuspectsFrom() (+14 more)

### Community 221 - "SystemTrusted"
Cohesion: 0.20
Nodes (11): hostView(), mountVouchesOwnership(), rootOwnedInode(), rootOwnedSafe(), SystemTrusted(), TestHostView(), systemTool(), archiveMembers() (+3 more)

### Community 222 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 223 - "StatInode"
Cohesion: 0.07
Nodes (34): setConfinedHomes(), ReplacementCheck, ExeKey(), IsMulticallRefusal(), registerBuild(), singleBinaryPaths(), TestExeKeyRefusesGNUCoreutilsSingle(), TestExeKeyRefusesOpaqueMulticall() (+26 more)

### Community 224 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 228 - "mountTable"
Cohesion: 0.13
Nodes (10): MountInfo, mountinfoUnitTest, TestNoteRootMountFollowsMountID(), MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+2 more)

### Community 229 - "PlanDaemonUnitUpdate"
Cohesion: 0.27
Nodes (10): ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate(), ShippedDaemonUnit(), fakeBinary(), TestPlanDaemonUnitUpdate(), TestPlanDaemonUnitUpdateDropsLocalEdits(), TestShippedDaemonUnit() (+2 more)

### Community 231 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 232 - "runBackend"
Cohesion: 0.31
Nodes (7): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), changePasswordLive()

### Community 233 - "maintainer/plan.go"
Cohesion: 0.39
Nodes (7): isWhiteout(), opaque(), pathByte(), pathRecord(), scanLib(), opKind, walker

### Community 234 - "ValidateTargets"
Cohesion: 0.16
Nodes (12): existingSection(), guardedMatch(), overlaps(), refusedPath(), sectionFor(), validateTarget(), ValidateTargets(), mustParse() (+4 more)

### Community 235 - "AddServeFlags"
Cohesion: 0.36
Nodes (8): AddHeadlessFlag(), AddServeFlags(), TestParseServeFlagsRequiresCredentials(), init(), init(), init(), init(), init()

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "update/update.go"
Cohesion: 0.13
Nodes (35): installHooks(), staged(), applyUpdate(), archAssets(), assetsFor(), channelOf(), checkUpdatePreconditions(), confirmUnit() (+27 more)

### Community 238 - "GenerateConf"
Cohesion: 0.05
Nodes (68): refreshCatalogConfig(), discordCatalogEntry(), groupedDiscordConf(), TestGroupedSectionAddressedByEncryptionRoot(), TestPatchCatalogSectionGroupedConfig(), TestSectionsFromCandidatesSkipsMissingWatchSubPaths(), TestSetSectionWhitelistPreservesGroupStructure(), LibraryBlock (+60 more)

### Community 239 - "guardLSMHooks"
Cohesion: 0.25
Nodes (5): lsmHook, trustHook, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), TestSocketHooksAreRequired()

### Community 242 - "ListUsers"
Cohesion: 0.08
Nodes (29): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), startMaintainers(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), grantNewBinaries() (+21 more)

### Community 243 - "ReadInstallCert"
Cohesion: 0.36
Nodes (8): InstallCert, ReadInstallCert(), RemoveInstallCert(), RequireNoInstallCert(), TestInstallCertRoundTrip(), TestInstallCertSymlinkRefuses(), TestInstallCertUnreadableRefuses(), withCertPath()

### Community 244 - "sources_test.go"
Cohesion: 0.38
Nodes (4): newSources(), TestSourcesHoldFilesAndArchiveMembers(), TestSourcesRefuseOtherFilesystem(), writeTar()

### Community 245 - "runGenKey"
Cohesion: 0.22
Nodes (10): confirmMasterKeyOverwrite(), runGenKey(), ensureInstalledBinary(), ensureMasterKey(), prepareInstallation(), TestEnsureInstalledBinary(), checkKeyLen(), GenerateMasterKey() (+2 more)

### Community 247 - "Multicall"
Cohesion: 0.29
Nodes (7): Multicall, MulticallKind, multicallFamily(), classifyBytes(), fwd(), TestClassifyMulticall(), TestMulticallMarkersReversed()

### Community 248 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 249 - "linkFixture"
Cohesion: 0.43
Nodes (7): newLinkView(), checkLink(), linkFixture(), TestCheckLinksAccepts(), TestCheckLinksJudgesEachTarget(), TestCheckLinksRefusesUserFiles(), TestCheckLinksSeesThePlannedTree()

### Community 250 - "runMonitor"
Cohesion: 0.47
Nodes (5): MakeDisplayPaths(), runHeadless(), runMonitor(), runTUI(), NewModel()

### Community 251 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 256 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 257 - "network monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 258 - "parseReservedPorts"
Cohesion: 0.50
Nodes (3): parseReservedPorts(), readReservedPorts(), TestParseReservedPorts()

### Community 259 - "parseKernelRelease"
Cohesion: 0.67
Nodes (4): ciliumVersion, parseKernelRelease(), TestParseKernelRelease(), unameKernelVersion()

### Community 260 - "Execute"
Cohesion: 0.67
Nodes (3): Execute(), initLogger(), main()

### Community 267 - "Config"
Cohesion: 0.05
Nodes (79): loadDaemonConfig(), noteReloadedSettings(), ensureAnchor(), ensurePortAnchors(), ensureRootDir(), ephemeralRange(), ensureSecretAnchors(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup() (+71 more)

### Community 273 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 294 - "RunWith"
Cohesion: 0.13
Nodes (18): needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), applyLiveRefresh(), runRefreshCatalogOnly(), softenAutomatedRefreshErr() (+10 more)

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 664 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **51 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.TestDaemon_Bypass_PrctlExeFileSwap`, `parseDaemonEvents`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.netIntegrityContainer`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.09112613803785274 - nodes in this community are weakly interconnected._
- **Why does `monitor()` connect `monitor` to `app-listener documentation`, `time.Duration`, `secret_probe/main.go`, `Installation`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.017670682730923693 - nodes in this community are weakly interconnected._
- **Why does `Guard` connect `Guard` to `go_pkg_strings`, `time.Duration`, `sync.Mutex`, `rootHandle`, `guardUnitTest`, `DaemonEvent`, `time.Time`, `GuardEvent`, `batchInodeRows`, `buildOneGuard`, `TrustGuard`, `os.File`, `BinaryEntry`, `Resource`, `loginServer`, `GuardInodeKey`, `reloadOnce`, `SanitizeText`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0943859649122807 - nodes in this community are weakly interconnected._