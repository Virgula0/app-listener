# Graph Report - app-listener  (2026-11-06)

## Corpus Check
- 574 files · ~1,819,627 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7537 nodes · 34929 edges · 269 communities (219 shown, 50 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 3097 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `1ab868ab`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- go_pkg_strings
- testing.T
- IntegrationSuite
- User
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- fakeGuardRepo
- ConfinedEntry
- RunCmd
- Guard
- buildOneGuard
- DaemonEvent
- IntegrationSuite
- guard_trust.bpf.c
- widgets.go
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- op
- guard.bpf.c
- github.com/charmbracelet/bubbles/textarea.Model
- inode_dev
- gui/session.go
- rcUser
- loadDaemonConfig
- SampleContent
- CandidateDir
- journal
- picker
- secret_probe/main.go
- IntegrationSuite
- Message
- editView
- NetEventType
- runView
- watchSet
- fileedit_model_test.go
- monitorUnitTest
- pinstate.go
- cursor.go
- time.Time
- stdio
- .showSized
- fadeScroll
- EventType
- protectedPage
- Client
- highlight_test.go
- .admitGated
- .apply
- backend_test.go
- FileTarget
- GuardInodeKey
- Network
- GuardEvent
- serve.go
- IntegrationSuite
- golang.org/x/sys/unix.Stat_t
- FileEvent
- unix_connect.c
- IntegrationSuite
- vaultfs.go
- globBuilder
- parseDaemonEvents
- wipe
- render.go
- Config
- fcntl
- policyEditor
- headless/headless.go
- codecheck/main.go
- NetGuard
- liveSession
- process_vm_readv.c
- .filterEvent
- controlServer
- daemon.conf reference
- daemon_multicall_test.go
- classifyGroup
- bpfstats/main.go
- exe_supersede.h
- loginServer
- UI
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- LaunchRule
- DumpHook
- netinject.c
- README.md
- .populateWatchInodes
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- openLedger
- tabs.go
- Well-known bypass classes
- main_test.go
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- github.com/charmbracelet/bubbletea.Model
- Build & sign release assets (reusable workflow)
- FindDaemonProcesses
- newTestTarget
- run.go
- Monitor
- Prereq
- IntegrationSuite
- RefreshSection
- Recorder
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- NetworkMonitor
- IntegrationSuite
- WebSocket /ws client
- installPage
- inet_connect.c
- install.sh
- uninstallAll
- SanitizeText
- make test-integration (rootful Docker suite)
- SelectOrphans
- IntegrationSuite
- github.com/spf13/cobra.Command
- trace-app-libs.sh
- IntegrationSuite
- safeio.go
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- session
- Binary
- runOfflineEdit
- Ledger
- TestMonitorUseCaseLifecycle
- maintainerKeeper
- btrfsLayoutFrom
- changelog.go
- recordFlow
- fakeStore
- RunWith
- IntegrationSuite
- os.FileInfo
- verdict
- io/fs.FileInfo
- wire.go
- Installation
- NetGuardEvent
- graphify-refresh.sh
- .Layout
- update_test.go
- .judgeProcess
- app-listener documentation
- feed.go
- shell
- .netIntegrityContainer
- network monitor: see what a program talks to
- activity
- monitor: see what touches a path
- Vault
- Limitations
- startCatalogRefresh
- plan_test.go
- net_tester/main.go
- PlanDaemonUnitUpdate
- sources
- tempGrantSetup
- shQuote
- filevault_test.go
- sync.Mutex
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- fakeDaemon
- os.File
- Self-updating apps and the catalog
- downloadFile
- WithBottomBar
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- Gate
- pathCache
- NewDumpHook
- classify
- BusClient
- RawTarget
- IntegrationSuite
- monitor
- IntegrationSuite
- absPath
- writeWithin
- configedit_test.go
- github.com/stretchr/testify/suite.Suite
- update/update.go
- options.go
- .awaitDaemonUpOrDead
- .TestRecord_GuardParity
- runGenKey
- bpfLsmListed
- find_test.go
- linkView
- searchDirs
- inodeChain
- Development
- runNetworkGuard
- supersede_probe.c
- planUpdaters
- runWizard
- launch_probe.c
- IntegrationSuite
- github.com/Virgula0/app-listener/tools/codecheck
- runBackend
- taint_victim.c
- DiscoverInfraBinaries
- mutedText
- promptPassword
- Troubleshooting
- mountTable
- GenerateConf
- logLine
- maintainer/plan.go
- record
- runLockdown
- TestCleanRel
- fileedit_symlink_test.go
- parseKernelRelease
- .handle
- icons.go
- golang.org/x/sys/unix.Timespec
- AddServeFlags
- Plan
- fileID
- MulticallError
- ParseServeFlags
- .eventTypeAllowed
- applyLaunchRule
- sources_test.go
- logGrid
- isRegularFileTarget
- CleanupStalePins
- bleedLayout
- launchKey
- Vault
- archiveMembers
- watchGroup
- Execute

## God Nodes (most connected - your core abstractions)
1. `Config` - 149 edges
2. `Guard` - 125 edges
3. `GuardInodeKey` - 113 edges
4. `Load()` - 101 edges
5. `absPath()` - 97 edges
6. `shell` - 97 edges
7. `IntegrationSuite` - 83 edges
8. `writeConfig()` - 79 edges
9. `fileEditModel` - 77 edges
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

## Communities (269 total, 50 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.09
Nodes (21): tempJournalRow, Progress, User, invalidPassword, recordTick, linkedFile, argvMatch, launchRule (+13 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (186): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+178 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (8): Hash(), configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig()

### Community 3 - "User"
Cohesion: 0.07
Nodes (49): entryWriters(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), gatherPerUserSetup(), addedResources(), applyDiffAdditions() (+41 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.06
Nodes (87): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+79 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.09
Nodes (51): prctl_gate(), guard_task_prctl(), admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader() (+43 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.23
Nodes (5): guardBinaryFlag(), IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent()

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (18): addErrKind, TestWithChmodDropWriteNeedsReadOnly(), classifyAddErr(), Mode, guardModeKey(), modeString(), addRecorder(), copySelf() (+10 more)

### Community 9 - "mkdirs"
Cohesion: 0.10
Nodes (46): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), newWatchPattern(), newFakeInotify() (+38 more)

### Community 10 - "fakeGuardRepo"
Cohesion: 0.06
Nodes (21): binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps(), ResolveTempBinary(), stripPinnedAllow(), StripPinnedTempAllows() (+13 more)

### Community 11 - "ConfinedEntry"
Cohesion: 0.13
Nodes (29): binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats() (+21 more)

### Community 12 - "RunCmd"
Cohesion: 0.05
Nodes (59): offerDaemonStart(), runDiffCatalog(), installFacultatives(), installGUI(), buildBinaryIfNeeded(), certifyInstall(), deploy(), deployFiles() (+51 more)

### Community 13 - "Guard"
Cohesion: 0.04
Nodes (31): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), deferredBinary, exeRow, keptEntries, rootHandle (+23 more)

### Community 14 - "buildOneGuard"
Cohesion: 0.08
Nodes (42): backingDeviceUnion(), buildGuards(), buildOneGuard(), eventFilterOptions(), reloadOnce(), reloadSlotsNeeded(), startGuardedDaemon(), startGuardedDaemonAbortable() (+34 more)

### Community 15 - "DaemonEvent"
Cohesion: 0.10
Nodes (33): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), runDaemon(), runDaemonUI(), runHeadless() (+25 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.08
Nodes (69): argv_env(), argv_matched(), argv_not_root(), argv_rule_at(), argv_short(), argv_word(), caller_updates(), creds_overridden() (+61 more)

### Community 18 - "widgets.go"
Cohesion: 0.03
Nodes (21): appTheme, binaryRow, clickPanel, handButton, iconTile, iconTileRenderer, navItem, navItemRenderer (+13 more)

### Community 19 - "fileEditModel"
Cohesion: 0.07
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.13
Nodes (26): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+18 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.13
Nodes (52): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+44 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (11): Resource, GuardRepository, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, TestRepeatLogOncePerErrorAndOnRecovery(), uniqueEncryptionRoots(), unlockRoots() (+3 more)

### Community 23 - "auth_test.go"
Cohesion: 0.12
Nodes (40): CheckPassword(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf(), parseHash() (+32 more)

### Community 24 - "op"
Cohesion: 0.17
Nodes (13): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+5 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.07
Nodes (98): exe_refused(), add_inode_to_guard(), admit_port(), btrfs_copy_gate(), bus_key(), check_and_emit_args(), chmod_only_drops_write(), code_suspect() (+90 more)

### Community 26 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.08
Nodes (30): class, pos, diffModel, editorModel, classOf(), column(), dedent(), Edit() (+22 more)

### Community 27 - "inode_dev"
Cohesion: 0.43
Nodes (4): ctx_ptr(), inode_dev(), sb_dev(), guard_inode_free()

### Community 28 - "gui/session.go"
Cohesion: 0.07
Nodes (40): invokingIDs(), ancestorEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+32 more)

### Community 29 - "rcUser"
Cohesion: 0.10
Nodes (49): setupBunTmpdir(), listFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative(), TestListFacultativesIgnoresDaemonConf(), TestListFacultativesMarksGuardedSSH(), rcTarget, BunTmpDir() (+41 more)

### Community 30 - "loadDaemonConfig"
Cohesion: 0.08
Nodes (23): loadDaemonConfig(), noteReloadedSettings(), ensureAnchor(), ensurePortAnchors(), ensureRootDir(), ephemeralRange(), parseReservedPorts(), readReservedPorts() (+15 more)

### Community 31 - "SampleContent"
Cohesion: 0.20
Nodes (14): installFileAs(), installServices(), installUnit(), upsertFile(), isInstallerSSHAgentUnit(), TestIsInstallerSSHAgentUnit(), TestDaemonUnitCmdPrintsBundledUnit(), SampleContent() (+6 more)

### Community 32 - "CandidateDir"
Cohesion: 0.03
Nodes (83): swapFile(), appendSectionsAndEdit(), collectDiffAdditions(), TestDiffMergePreservesExistingAndParses(), TestUncoveredCandidates(), uncoveredCandidates(), editConfig(), groupCandidates() (+75 more)

### Community 33 - "journal"
Cohesion: 0.14
Nodes (7): fakeEngines, fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 34 - "picker"
Cohesion: 0.07
Nodes (33): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+25 more)

### Community 35 - "secret_probe/main.go"
Cohesion: 0.18
Nodes (31): busPath(), connexec(), connFD(), createItem(), dial(), dupkey(), execInto(), failf() (+23 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (14): checkSecret(), checkTrustLine(), trustRows(), File, TrustItem, backend, backendOp, backendState (+6 more)

### Community 38 - "editView"
Cohesion: 0.12
Nodes (7): codeEntry, editView, Entry, fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 39 - "NetEventType"
Cohesion: 0.08
Nodes (16): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), FormatAddr(), NetEventTypes(), Ntohs() (+8 more)

### Community 40 - "runView"
Cohesion: 0.09
Nodes (7): daemonPage, runHost, runView, stageItem, stageState, captureLogs(), orDefault()

### Community 41 - "watchSet"
Cohesion: 0.13
Nodes (15): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+7 more)

### Community 42 - "fileedit_model_test.go"
Cohesion: 0.21
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 44 - "pinstate.go"
Cohesion: 0.18
Nodes (22): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+14 more)

### Community 45 - "cursor.go"
Cohesion: 0.08
Nodes (17): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+9 more)

### Community 46 - "time.Time"
Cohesion: 0.04
Nodes (49): LockedError, sweepID, sweepMemo, Stats, parseKBValue(), parseStat(), parseStatus(), Read() (+41 more)

### Community 48 - ".showSized"
Cohesion: 0.10
Nodes (14): binaryList, installUI, multiAnswer, sizedDialog, scrollable(), selectable(), newFadeScroll(), ask() (+6 more)

### Community 49 - "fadeScroll"
Cohesion: 0.05
Nodes (30): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+22 more)

### Community 50 - "EventType"
Cohesion: 0.07
Nodes (22): selfEventsFor(), BpfEvent, EventType, eventUnitTest, Target, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction() (+14 more)

### Community 51 - "protectedPage"
Cohesion: 0.05
Nodes (23): backendLink, configView, forwardForm, forwardView, protectedPage, recordForm, recordView, trustPage (+15 more)

### Community 52 - "Client"
Cohesion: 0.10
Nodes (10): Client, launch(), TestHelperBackend(), TestReaderCapsLines(), Reader, Writer, NewReader(), NewWriter() (+2 more)

### Community 53 - "highlight_test.go"
Cohesion: 0.13
Nodes (26): New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines(), TestSetTextStartsAtFirstLine() (+18 more)

### Community 54 - ".admitGated"
Cohesion: 0.33
Nodes (3): Caller, Role, Judge

### Community 55 - ".apply"
Cohesion: 0.13
Nodes (12): Apply(), ProgressFunc, TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan(), applier (+4 more)

### Community 56 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 57 - "FileTarget"
Cohesion: 0.13
Nodes (11): guardPill(), Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec(), overlaps() (+3 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (55): engine, inodeBatchFn, portOwner, refusedReplacements, resInfo, taintOwner, confineConnect(), isConnectConfined() (+47 more)

### Community 59 - "Network"
Cohesion: 0.09
Nodes (13): netGuardPill(), NetworkGuardRepository, newExeResolver(), guardRow(), clonePolicy(), NetEngines, NetExtra, NetPolicy (+5 more)

### Community 60 - "GuardEvent"
Cohesion: 0.07
Nodes (21): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), holdHeadless(), HeadlessLine(), bpfGuardEvent, snapGuard, TestDispatchFullChannelNeverDropsDenials() (+13 more)

### Community 61 - "serve.go"
Cohesion: 0.09
Nodes (30): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+22 more)

### Community 62 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 63 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.13
Nodes (21): placedWalk, treeJudge, confineBelow(), inUserHome(), vouchedFd(), GlobSystemPlaced(), placedDir(), placedEntry() (+13 more)

### Community 64 - "FileEvent"
Cohesion: 0.18
Nodes (5): snapMonitor, FileEvent, newFakeMonitorRepo(), fakeSource, fakeMonitorRepo

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 67 - "vaultfs.go"
Cohesion: 0.09
Nodes (26): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), Vault, IsBinaryContent(), IsOctalMode(), Open(), ParseMode() (+18 more)

### Community 68 - "globBuilder"
Cohesion: 0.11
Nodes (16): reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, GlobChild, compileGlobNames(), globText() (+8 more)

### Community 69 - "parseDaemonEvents"
Cohesion: 0.15
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite

### Community 70 - "wipe"
Cohesion: 0.13
Nodes (17): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), isLockedRegularFileErr(), isNotEncryptedErr(), newBoundedKeyFn(), readKey(), TestIsNotEncryptedErr() (+9 more)

### Community 71 - "render.go"
Cohesion: 0.09
Nodes (18): findButton, lineStyles, span, viewportState, findAll(), lowerRunes(), newFindInput(), overlay() (+10 more)

### Community 72 - "Config"
Cohesion: 0.06
Nodes (55): sameConfig(), newDaemonModel(), runTUI(), inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), allowedEnv() (+47 more)

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "policyEditor"
Cohesion: 0.13
Nodes (12): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), newNetGuardTab() (+4 more)

### Community 75 - "headless/headless.go"
Cohesion: 0.23
Nodes (10): stagedBinary, hashFile(), hashReader(), stageRelease(), Check(), CheckGUI(), inspect(), buildTiny() (+2 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "NetGuard"
Cohesion: 0.08
Nodes (16): NewMountVouch(), WithSuspectsFrom(), mcNameKey(), modeLabel(), NewNetGuard(), statInodeKey(), judgePasses(), parseProcs() (+8 more)

### Community 78 - "liveSession"
Cohesion: 0.14
Nodes (18): dialLiveSession(), parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), putConfig(), chooseResources(), confirmUserPlaced(), eventsLabel() (+10 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 81 - "controlServer"
Cohesion: 0.08
Nodes (24): readConfigPut(), TestReadGrantRequestConfig(), authSelfPeer(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest() (+16 more)

### Community 82 - "daemon.conf reference"
Cohesion: 0.08
Nodes (25): A complete example, `allow_env`, `allow_lib`, Binary lines and events, daemon.conf reference, `deny_args`, `warn_args`, `deny_env`, `warn_env`, `[electron_apps]` (+17 more)

### Community 83 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 84 - "classifyGroup"
Cohesion: 0.19
Nodes (14): appletLine(), classifyGroup(), classifyLibs(), innermostResource(), libCovered(), newBinary(), verifyExe(), worldWritableDir() (+6 more)

### Community 85 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 86 - "exe_supersede.h"
Cohesion: 0.22
Nodes (14): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+6 more)

### Community 87 - "loginServer"
Cohesion: 0.09
Nodes (20): markReady(), markReadyWithLogin(), notifySystemdReady(), hasSentinel(), lockParent(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel() (+12 more)

### Community 88 - "UI"
Cohesion: 0.10
Nodes (18): nextConfig(), showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), cancelUI, EditText() (+10 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "LaunchRule"
Cohesion: 0.08
Nodes (34): ShortOpts, argvEntry, LaunchRule, LaunchWordForms(), EncodeLaunchRules(), ShortRefused(), canonLaunchRules(), encodeCatalog() (+26 more)

### Community 92 - "DumpHook"
Cohesion: 0.24
Nodes (5): logHook, formatEntry(), sortedKeys(), sortStrings(), DumpHook

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

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

### Community 100 - "openLedger"
Cohesion: 0.33
Nodes (9): ConfirmGiven(), ConfirmPending(), stillHolds(), openLedger(), TestConfirmGivenNeedsTheReviewedContent(), TestConfirmPendingNeverOpensALineTheLedgerDoesNotRefuse(), TestConfirmPendingOnlyWhatWasRefusedAndIsStillThere(), TestPendingRowsSkipChangedContent() (+1 more)

### Community 101 - "tabs.go"
Cohesion: 0.12
Nodes (7): tabBar, tabButton, tabButtonRenderer, tabItem, TestTabBarShowsOnlySelectedPage(), newTabBar(), newTabButton()

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "main_test.go"
Cohesion: 0.12
Nodes (6): IntegrationSuite, infraContainerPaths(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "github.com/charmbracelet/bubbletea.Model"
Cohesion: 0.10
Nodes (9): recordTickCmd(), recordModel, renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel, serveTestModel (+1 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "FindDaemonProcesses"
Cohesion: 0.10
Nodes (22): checkAndApply(), applyLiveRefresh(), runRefreshCatalogOnly(), softenAutomatedRefreshErr(), refreshCatalogConfig(), TestApplyLiveRefreshDeliversReload(), TestSoftenAutomatedRefreshErr(), ConfirmOverwrite() (+14 more)

### Community 110 - "newTestTarget"
Cohesion: 0.24
Nodes (15): NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded(), TestApplyFromWhitelistStopsFirst() (+7 more)

### Community 111 - "run.go"
Cohesion: 0.23
Nodes (12): rootOwnedInode(), rootOwnedSafe(), SystemTrusted(), systemTool(), planFrom(), newStage(), Run(), runSandbox() (+4 more)

### Community 112 - "Monitor"
Cohesion: 0.18
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 113 - "Prereq"
Cohesion: 0.33
Nodes (3): confirmRunPrereq(), Prereq, fakeVault

### Community 115 - "RefreshSection"
Cohesion: 0.08
Nodes (32): logRefreshChange(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange(), RefreshOptions (+24 more)

### Community 116 - "Recorder"
Cohesion: 0.21
Nodes (11): drain(), Recorder, Target, newRecorder(), selfExe(), Start(), TestProgressNamesMulticallApplets(), accessKey (+3 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "NetworkMonitor"
Cohesion: 0.13
Nodes (8): BinariesSummary(), TestBinariesSummary(), PumpRingbuf(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "installPage"
Cohesion: 0.07
Nodes (22): anyDaemonRunning(), actionGate, binSource, installAction, installedBinary, installPage, latestRelease, updatePage (+14 more)

### Community 124 - "inet_connect.c"
Cohesion: 0.25
Nodes (17): attack_connect(), attack_fastopen(), attack_uring(), denied(), dial(), fill_addr(), handoff(), listen_on() (+9 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "uninstallAll"
Cohesion: 0.13
Nodes (20): releaseBuild(), Run(), confirmKeyDeletion(), run(), runUninstall(), RunWith(), removeKeyAndEmptyDir(), removeMasterKey() (+12 more)

### Community 127 - "SanitizeText"
Cohesion: 0.06
Nodes (29): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), scriptIdentity(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome() (+21 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "SelectOrphans"
Cohesion: 0.23
Nodes (16): appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors(), SelectOrphans() (+8 more)

### Community 131 - "github.com/spf13/cobra.Command"
Cohesion: 0.27
Nodes (10): TestSetupDumpLogRefusesSymlink(), asRoot(), runInstall(), closeDumpFile(), effectiveVerbose(), headlessRequested(), postRun(), preRun() (+2 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "safeio.go"
Cohesion: 0.23
Nodes (14): readAtNoFollow(), upsertUnitAt(), writeAtNoFollow(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), createOwned(), AtomicWriteAt(), CreateExclNoFollow() (+6 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "session"
Cohesion: 0.11
Nodes (7): peerCred(), placeableAttributes(), TestPlaceableAttributes(), TestReadFrame(), decodeFrame(), readFrame(), session

### Community 143 - "Binary"
Cohesion: 0.19
Nodes (15): lines(), pickBinaries(), pickLibs(), reviewRecording(), beyondEvents(), Classify(), finishBinary(), Binary (+7 more)

### Community 144 - "runOfflineEdit"
Cohesion: 0.13
Nodes (14): auditAfterEdit(), auditFindings(), TestPickEncryptedDirSingleSkipsPicker(), chooseResource(), readPutContent(), runPut(), pickDir(), pickEncryptedDir() (+6 more)

### Community 145 - "Ledger"
Cohesion: 0.09
Nodes (15): candidates(), confirmAll(), describe(), record(), run(), applySchema(), Ledger, Pending (+7 more)

### Community 146 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 147 - "maintainerKeeper"
Cohesion: 0.18
Nodes (12): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), startMaintainers(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), maintainedApp (+4 more)

### Community 148 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 149 - "changelog.go"
Cohesion: 0.15
Nodes (26): changelogText(), cumulativeChangelog(), newChangelogModel(), newReleaseWindow(), releasesSince(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit() (+18 more)

### Community 150 - "recordFlow"
Cohesion: 0.13
Nodes (15): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+7 more)

### Community 152 - "RunWith"
Cohesion: 0.14
Nodes (23): needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), deletePostBackups(), listBackups(), restoreBackups() (+15 more)

### Community 154 - "os.FileInfo"
Cohesion: 0.21
Nodes (5): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata()

### Community 155 - "verdict"
Cohesion: 0.33
Nodes (8): arg(), attributesOf(), classifyCollection(), classifyProps(), classifyService(), withSchema(), verdict, verdictKind

### Community 156 - "io/fs.FileInfo"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 157 - "wire.go"
Cohesion: 0.16
Nodes (13): TestClientAuth(), TestSASL(), clientAuth(), distinctDictKeys(), header(), sameDict(), sameEntries(), sameList() (+5 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "NetGuardEvent"
Cohesion: 0.18
Nodes (6): snapNetGuard, NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 161 - ".Layout"
Cohesion: 0.44
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 162 - "update_test.go"
Cohesion: 0.09
Nodes (23): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleaseByTag(), TestFetchReleases(), TestFetchReleasesHTTPError() (+15 more)

### Community 163 - ".judgeProcess"
Cohesion: 0.16
Nodes (12): reasonText(), codeVouched(), isTraced(), loaderEnv(), readMountDevs(), roFor(), TestCodeVouched(), TestLoaderEnv() (+4 more)

### Community 164 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 165 - "feed.go"
Cohesion: 0.07
Nodes (38): TestForwardFlagSurface(), actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction (+30 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (53): fakeLink, homePage, pageDef, pageID, reel, resizeWatch, shell, snapEngines (+45 more)

### Community 168 - "network monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 169 - "activity"
Cohesion: 0.15
Nodes (7): newActivity(), processKey(), TestRingKeepsNewest(), Counts, newRing(), activity, ring

### Community 170 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 171 - "Vault"
Cohesion: 0.18
Nodes (15): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), encryptDirectories(), secureResources(), verifyEncryptionState(), resolveFilesystemPrereqs(), Vault (+7 more)

### Community 172 - "Limitations"
Cohesion: 0.17
Nodes (12): Already-running processes, Binaries you can write, Commands that stop working, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root (+4 more)

### Community 173 - "startCatalogRefresh"
Cohesion: 0.10
Nodes (14): drain(), newCatalogRefresher(), readEvents(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), newControlManager(), startControlManager() (+6 more)

### Community 174 - "plan_test.go"
Cohesion: 0.40
Nodes (11): scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent(), TestScanAndValidateAcceptsLibDirContent(), TestScanRefusesSpecialFiles(), TestScanUntouchedLibDir(), TestTouches() (+3 more)

### Community 176 - "net_tester/main.go"
Cohesion: 0.19
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 177 - "PlanDaemonUnitUpdate"
Cohesion: 0.21
Nodes (12): TestDaemonUnitWithMetadataOutput(), DaemonUnitWithMetadataOutput(), ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate(), ShippedDaemonUnit(), fakeBinary(), TestPlanDaemonUnitUpdate() (+4 more)

### Community 178 - "sources"
Cohesion: 0.21
Nodes (4): lstatOf(), outsideHardlinks(), libDirDev(), sources

### Community 179 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 180 - "shQuote"
Cohesion: 0.15
Nodes (6): IntegrationSuite, IntegrationSuite, overlayPlanter(), shQuote(), inNosuidNamespace(), IntegrationSuite

### Community 181 - "filevault_test.go"
Cohesion: 0.09
Nodes (47): decryptDirectories(), TestDecryptDirectoriesEmpty(), classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord() (+39 more)

### Community 182 - "sync.Mutex"
Cohesion: 0.05
Nodes (27): lsmHook, tempGrant, tempMaskOp, tempRow, trustHook, trustRows, fyneLogFilter, ensureBtrfsLayout() (+19 more)

### Community 186 - "Binary trust"
Cohesion: 0.20
Nodes (10): Binary ledger, Binary replacement, Binary trust, Browser profiles, Hash cache, Launch rules, Launch scan, Multicall binaries (+2 more)

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 191 - "os.File"
Cohesion: 0.03
Nodes (106): checkFreshDir(), openDirNoFollow(), renewDir(), setConfinedHomes(), placedKey(), ReplacementCheck, statFunc, VettedInode (+98 more)

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "downloadFile"
Cohesion: 0.32
Nodes (5): downloadFile(), downloadReleaseFiles(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), progressReader

### Community 194 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 199 - "Gate"
Cohesion: 0.19
Nodes (5): connectionGone(), Gate, Listen(), TestConnectionGoneOnlyFromTheBus(), Bus

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 202 - "classify"
Cohesion: 0.19
Nodes (12): call(), dictFrame(), TestClassify(), TestClassifyCreateItemReplace(), TestClassifyItemTypeIsSchema(), TestClassifyOpenCallsExact(), TestDuplicateDictKeyRefused(), TestLearnKeyringFromFirstReply() (+4 more)

### Community 203 - "BusClient"
Cohesion: 0.27
Nodes (7): BusClient, AttestBusClient(), busClientKey(), execStamp(), imageOf(), pidfdPid(), procIdentity()

### Community 204 - "RawTarget"
Cohesion: 0.31
Nodes (9): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+1 more)

### Community 206 - "monitor"
Cohesion: 0.07
Nodes (33): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks (+25 more)

### Community 209 - "absPath"
Cohesion: 0.06
Nodes (8): IntegrationSuite, IntegrationSuite, IntegrationSuite, IntegrationSuite, IntegrationSuite, launchCase, absPath(), TestClientVerifyRunsBeforeAnythingIsSent()

### Community 210 - "writeWithin"
Cohesion: 0.25
Nodes (8): atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync(), writeWithin()

### Community 211 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 212 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.24
Nodes (3): harnessSuite, harnessSuite, editorHarness

### Community 213 - "update/update.go"
Cohesion: 0.11
Nodes (40): installHooks(), staged(), applyUpdate(), archAssets(), assetsFor(), channelOf(), checkUpdatePreconditions(), confirmUnit() (+32 more)

### Community 214 - "options.go"
Cohesion: 0.36
Nodes (9): InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks, InstallHooks, StagedRelease (+1 more)

### Community 216 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 217 - "runGenKey"
Cohesion: 0.22
Nodes (10): confirmMasterKeyOverwrite(), runGenKey(), ensureInstalledBinary(), ensureMasterKey(), prepareInstallation(), TestEnsureInstalledBinary(), checkKeyLen(), GenerateMasterKey() (+2 more)

### Community 218 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 219 - "find_test.go"
Cohesion: 0.18
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 220 - "linkView"
Cohesion: 0.31
Nodes (5): linkTarget(), newLinkView(), TestCheckLinksJudgesEachTarget(), entryKind, linkView

### Community 221 - "searchDirs"
Cohesion: 0.18
Nodes (10): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths(), resolveSoname() (+2 more)

### Community 222 - "inodeChain"
Cohesion: 0.23
Nodes (15): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD(), keyOf() (+7 more)

### Community 223 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 224 - "runNetworkGuard"
Cohesion: 0.07
Nodes (43): DrainUntilSignal(), CheckBPFLSM(), CheckEBPF(), MakeDisplayPaths(), ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), newPinGeneration() (+35 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "planUpdaters"
Cohesion: 0.13
Nodes (20): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+12 more)

### Community 227 - "runWizard"
Cohesion: 0.09
Nodes (24): checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), fileSHA256(), runWizard(), sameContent(), selectAndEditConfig(), cleanOrphanedFscrypt(), addManualDirectories() (+16 more)

### Community 228 - "launch_probe.c"
Cohesion: 0.46
Nodes (6): anchor(), copy(), is_verb(), main(), read_once(), usage()

### Community 231 - "runBackend"
Cohesion: 0.31
Nodes (7): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), changePasswordLive()

### Community 233 - "DiscoverInfraBinaries"
Cohesion: 0.33
Nodes (6): DiscoverInfraBinaries(), infraFromRunning(), runningExecutables(), pathIsRunning(), TestDiscoverInfraBinaries(), TestInfraFromRunning()

### Community 234 - "mutedText"
Cohesion: 0.08
Nodes (42): checkGrid, filesPage, handCheck, recordReview, newBinaryList(), newBinaryRow(), newCheck(), newCheckGroup() (+34 more)

### Community 235 - "promptPassword"
Cohesion: 0.28
Nodes (9): applyEdit(), editAgain(), fetchConfig(), runEditConfig(), promptPassword(), promptPasswordWith(), PromptNewPassword(), promptNewPassword() (+1 more)

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "mountTable"
Cohesion: 0.18
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo(), userMountable() (+1 more)

### Community 238 - "GenerateConf"
Cohesion: 0.04
Nodes (78): collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles() (+70 more)

### Community 239 - "logLine"
Cohesion: 0.25
Nodes (5): logLine, journalLine(), TestJournalLine(), logRow(), logStyle()

### Community 240 - "maintainer/plan.go"
Cohesion: 0.31
Nodes (7): isWhiteout(), opaque(), pathByte(), pathRecord(), TestPathRecord(), opKind, walker

### Community 241 - "record"
Cohesion: 0.22
Nodes (6): fixtureExe(), privateTempRoots(), record(), TestClassifyLibs(), TestClassifyRefusesWorldWritableLocations(), TestObserveKeysByKernelExeAndSkipsSelf()

### Community 242 - "runLockdown"
Cohesion: 0.29
Nodes (8): configureDaemonLogging(), resolveConfigPath(), runLockdown(), dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 244 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 245 - "parseKernelRelease"
Cohesion: 0.67
Nodes (4): ciliumVersion, parseKernelRelease(), TestParseKernelRelease(), unameKernelVersion()

### Community 249 - "AddServeFlags"
Cohesion: 0.43
Nodes (7): AddHeadlessFlag(), AddServeFlags(), init(), init(), init(), init(), init()

### Community 251 - "Plan"
Cohesion: 0.22
Nodes (6): dedupe(), Plan, libOfChosen(), sectionFor(), Existing, ExistingRule

### Community 256 - "ParseServeFlags"
Cohesion: 0.33
Nodes (7): credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials()

### Community 257 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 258 - "applyLaunchRule"
Cohesion: 0.33
Nodes (6): applyLaunchRule(), envName(), ParseAllowEnv(), parseLaunchRule(), BuiltinAllowedEnv(), TestBuiltinAllowedEnvFits()

### Community 259 - "sources_test.go"
Cohesion: 0.38
Nodes (4): newSources(), TestSourcesHoldFilesAndArchiveMembers(), TestSourcesRefuseOtherFilesystem(), writeTar()

### Community 261 - "isRegularFileTarget"
Cohesion: 0.40
Nodes (3): isFileVaultCiphertext(), isRegularFileTarget(), hasEncryptionPolicy()

### Community 262 - "CleanupStalePins"
Cohesion: 0.33
Nodes (6): CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 266 - "archiveMembers"
Cohesion: 0.67
Nodes (4): archiveMembers(), decompressed(), tarMembers(), TestArchiveMembersIgnoresNonArchives()

### Community 267 - "watchGroup"
Cohesion: 0.07
Nodes (56): watchGroup, addResource(), applyAllowLib(), applyAnyGroupDirective(), applyBinaryLine(), applyBlockLine(), applyConfigLine(), applyDirective() (+48 more)

### Community 268 - "Execute"
Cohesion: 0.67
Nodes (3): Execute(), initLogger(), main()

## Knowledge Gaps
- **34 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck` (+29 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 688 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **50 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Guard` connect `Guard` to `go_pkg_strings`, `time.Duration`, `guardUnitTest`, `ConfinedEntry`, `buildOneGuard`, `DaemonEvent`, `time.Time`, `sync.Mutex`, `loadDaemonConfig`, `GuardInodeKey`, `GuardEvent`, `fakeDaemon`, `os.File`?**
  _High betweenness centrality (0.033) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _34 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.0922397247872533 - nodes in this community are weakly interconnected._
- **Why does `monitor()` connect `monitor` to `time.Duration`, `secret_probe/main.go`, `app-listener documentation`, `Installation`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.016318574213311055 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `parseDaemonEvents`, `.netIntegrityContainer`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `daemon_multicall_test.go`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09263157894736843 - nodes in this community are weakly interconnected._