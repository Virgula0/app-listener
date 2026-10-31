# Graph Report - app-listener  (2026-10-31)

## Corpus Check
- 564 files · ~1,794,942 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7348 nodes · 33651 edges · 252 communities (209 shown, 43 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 3004 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `0aa430de`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
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
- editView
- Network
- runNetworkMonitor
- buildOneGuard
- sanitizeTerminalText
- IntegrationSuite
- guard_trust.bpf.c
- navItem
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- sync.Mutex
- guard.bpf.c
- LibraryClosure
- github.com/charmbracelet/bubbletea.Cmd
- backend.go
- runDaemon
- codeedit.go
- TestMonitorUseCaseLifecycle
- fileedit_model_test.go
- find.go
- picker
- secret_probe/main.go
- IntegrationSuite
- Message
- controlServer
- .posAt
- time.Time
- .showSized
- cursor.go
- monitorUnitTest
- CandidateDir
- .FilterExistingWhitelist
- Guard
- stdio
- journal
- fadeScroll
- IntegrationSuite
- protectedPage
- Client
- installPage
- session
- .apply
- FindDaemonProcesses
- mutedText
- GuardInodeKey
- go_pkg_strings
- runMonitor
- runView
- LaunchRule
- policyEditor
- unix_connect.c
- IntegrationSuite
- IntegrationSuite
- Recorder
- os.FileMode
- Vault
- trustManager
- planUpdaters
- fcntl
- find_test.go
- update_test.go
- codecheck/main.go
- recordFlow
- newTestVetter
- process_vm_readv.c
- os.File
- Binary
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
- SecretStoreNotes
- Gate
- Well-known bypass classes
- NetGuardEvent
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- expandPlaceholders
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- .refreshOnce
- BinaryEntry
- Active
- runOfflineEdit
- IntegrationSuite
- maintainerKeeper
- systemd.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- NetworkMonitor
- IntegrationSuite
- WebSocket /ws client
- inet_connect.c
- install.sh
- KernelDev
- classify
- make test-integration (rootful Docker suite)
- RunWith
- logLine
- FileEvent
- trace-app-libs.sh
- IntegrationSuite
- .judgeProcess
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- daemon: permanent protection with encryption at rest
- serve.go
- liveSession
- github.com/charmbracelet/bubbles/textarea.Model
- Ledger
- EventType
- TempBinary
- dialLiveSession
- Daemon configuration
- fakeStore
- mc_tag.h
- IntegrationSuite
- searchDirs
- startCatalogRefresh
- copyRegularAt
- dialog_test.go
- Installation
- backend_test.go
- graphify-refresh.sh
- layout.go
- NetEventType
- newTestTarget
- absPath
- feed.go
- shell
- inodeChain
- plan_test.go
- appTheme
- SelectOrphans
- batchInodeRows
- Limitations
- fetchReleases
- parseDaemonEvents
- TrustGuard
- net_tester/main.go
- debouncer
- githubRelease
- Run
- bufio.Reader
- IntegrationSuite
- tempGrantSetup
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- github.com/spf13/cobra.Command
- Self-updating apps and the catalog
- fileedit_symlink_test.go
- CopyTreeWithProgress
- downloadAndVerify
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- New
- pathCache
- pinstate_test.go
- DaemonEvent
- procstats.go
- sources
- monitor
- IntegrationSuite
- audit.go
- changelog_test.go
- .eventTypeAllowed
- parseGuardEvent
- WithBottomBar
- run
- .awaitDaemonUpOrDead
- .TestRecord_GuardParity
- linkView
- runNetworkGuard
- refuseOutsideHardlinks
- IntegrationSuite
- Use
- serverAuth
- VouchedDevs
- app-listener documentation
- supersede_probe.c
- launchKey
- bleedLayout
- mountTable
- parsePasswd
- github.com/Virgula0/app-listener/tools/codecheck
- inode_dev.h
- ValidateTargets
- Troubleshooting
- deploy
- GenerateConf
- IntegrationSuite
- taint_victim.c
- User
- install/install.go
- TestMain
- monitor.md
- TestCleanRel
- .netIntegrityContainer
- network monitor: see what a program talks to
- daemonconfig.go
- NewEventFanout
- IntegrationSuite
- .runEditorHarness
- go_pkg_os

## God Nodes (most connected - your core abstractions)
1. `Config` - 147 edges
2. `Guard` - 125 edges
3. `GuardInodeKey` - 112 edges
4. `shell` - 97 edges
5. `Load()` - 96 edges
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

## Communities (252 total, 43 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.18
Nodes (33): tempJournalRow, argvMatch, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion (+25 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (179): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+171 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (8): configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote(), watchPathsInConfig()

### Community 3 - "op"
Cohesion: 0.16
Nodes (13): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+5 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.06
Nodes (84): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+76 more)

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
Cohesion: 0.06
Nodes (17): addErrKind, Start(), verifiedBinary(), classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest (+9 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (34): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+26 more)

### Community 10 - "GuardEvent"
Cohesion: 0.06
Nodes (14): holdHeadless(), HeadlessLine(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), commMatchesGuardedBinary(), GuardEvent, logBacklogDenial() (+6 more)

### Community 11 - "editView"
Cohesion: 0.11
Nodes (9): codeEntry, editView, handButton, Entry, selectable(), fileSize(), joinRel(), newCodeEntry() (+1 more)

### Community 12 - "Network"
Cohesion: 0.17
Nodes (5): netGuardPill(), clonePolicy(), NetExtra, NetPolicy, Network

### Community 13 - "runNetworkMonitor"
Cohesion: 0.09
Nodes (15): DrainUntilSignal(), ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), runHeadless(), runHeadless(), runNetworkMonitor() (+7 more)

### Community 14 - "buildOneGuard"
Cohesion: 0.06
Nodes (57): CheckBPFLSM(), CheckEBPF(), buildGuards(), buildOneGuard(), eventFilterOptions(), newPinGeneration(), prepareDaemonStart(), reloadOnce() (+49 more)

### Community 15 - "sanitizeTerminalText"
Cohesion: 0.07
Nodes (26): listenForDaemonEvents(), tickDaemonStats(), formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), clipToWidth(), formatResourceBar() (+18 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (6): eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.08
Nodes (65): argv_env(), argv_matched(), argv_rule_at(), argv_word(), caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key() (+57 more)

### Community 18 - "navItem"
Cohesion: 0.03
Nodes (23): binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pill (+15 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (11): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (51): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+43 more)

### Community 22 - "Resource"
Cohesion: 0.05
Nodes (18): ownerPaths(), newControlManager(), startControlManager(), startControlServer(), makeReloadHandler(), configEditor, controlManager, fakeDaemon (+10 more)

### Community 23 - "auth.go"
Cohesion: 0.13
Nodes (40): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+32 more)

### Community 24 - "sync.Mutex"
Cohesion: 0.07
Nodes (20): Target, newFilesPage(), newActivity(), Files, FileSpec, FileTarget, guardRow(), newFileTarget() (+12 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.07
Nodes (95): exe_refused(), inode_dev(), add_inode_to_guard(), admit_port(), btrfs_copy_gate(), bus_key(), check_and_emit_args(), chmod_only_drops_write() (+87 more)

### Community 26 - "LibraryClosure"
Cohesion: 0.10
Nodes (21): hostView(), IsScript(), LibraryClosure(), mountVouchesOwnership(), rootOwnedInode(), rootOwnedSafe(), TestLibraryClosureRejectsUserWritableOriginDir(), runsNatively() (+13 more)

### Community 27 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.07
Nodes (14): recordTickCmd(), recordModel, diffModel, editorModel, listenForNetEvents(), tickNetStats(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast() (+6 more)

### Community 29 - "runDaemon"
Cohesion: 0.16
Nodes (13): awaitStartupOrSignal(), catchLifecycleSignals(), confirmMasterKeyOverwrite(), runDaemon(), runGenKey(), runOneShotMode(), startGuardedDaemonAbortable(), startPprof() (+5 more)

### Community 30 - "codeedit.go"
Cohesion: 0.10
Nodes (39): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), Navigate() (+31 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 32 - "fileedit_model_test.go"
Cohesion: 0.19
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 33 - "find.go"
Cohesion: 0.16
Nodes (11): findButton, span, findAll(), lowerRunes(), newFindInput(), overlay(), runesEqual(), TestFindAllSmartCase() (+3 more)

### Community 34 - "picker"
Cohesion: 0.09
Nodes (20): keyEntry, listOpts, pickEntry, picker, treeSearch, TestSearchTreeFollowsNoSymlink(), errText(), isFindShortcut() (+12 more)

### Community 35 - "secret_probe/main.go"
Cohesion: 0.22
Nodes (25): busPath(), connexec(), connFD(), createItem(), dial(), execInto(), failf(), fakebus() (+17 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (13): checkSecret(), checkTrustLine(), trustRows(), File, TrustItem, backend, backendOp, backendState (+5 more)

### Community 38 - "controlServer"
Cohesion: 0.16
Nodes (7): authSelfPeer(), readClientLines(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), controlServer, editControlSession

### Community 39 - ".posAt"
Cohesion: 0.11
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+6 more)

### Community 40 - "time.Time"
Cohesion: 0.14
Nodes (10): LockedError, isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, sweepID (+2 more)

### Community 41 - ".showSized"
Cohesion: 0.09
Nodes (13): binaryList, installUI, multiAnswer, sizedDialog, fadeScrollOf(), newFadeScroll(), ask(), diffText() (+5 more)

### Community 42 - "cursor.go"
Cohesion: 0.08
Nodes (21): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+13 more)

### Community 43 - "monitorUnitTest"
Cohesion: 0.15
Nodes (4): sameInode(), dirTarget(), newMountTable(), monitorUnitTest

### Community 44 - "CandidateDir"
Cohesion: 0.09
Nodes (38): appendSectionsAndEdit(), collectDiffAdditions(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates() (+30 more)

### Community 45 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 46 - "Guard"
Cohesion: 0.05
Nodes (28): backingDeviceUnion(), lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), deferredBinary, rootHandle, openRoot() (+20 more)

### Community 48 - "journal"
Cohesion: 0.15
Nodes (5): fakeGuard, fakeMonitor, fakeNet, fakeNetMonitor, journal

### Community 49 - "fadeScroll"
Cohesion: 0.07
Nodes (14): fadeBar, fadeBarRenderer, fadeScroll, fadeScrollRenderer, newTestFadeScroll(), settle(), TestFadeBarDragEndsGlide(), TestFadeBarDragMovesContent() (+6 more)

### Community 50 - "IntegrationSuite"
Cohesion: 0.16
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 51 - "protectedPage"
Cohesion: 0.05
Nodes (20): backendLink, configView, forwardForm, forwardView, protectedPage, recordReview, recordView, trustPage (+12 more)

### Community 52 - "Client"
Cohesion: 0.10
Nodes (10): Client, launch(), TestHelperBackend(), TestReaderCapsLines(), Reader, Writer, NewReader(), NewWriter() (+2 more)

### Community 53 - "installPage"
Cohesion: 0.05
Nodes (36): actionGate, binSource, handCheck, installAction, installedBinary, InstallMode, installPage, InstallRequest (+28 more)

### Community 54 - "session"
Cohesion: 0.11
Nodes (9): peerCred(), TestReadFrame(), decodeFrame(), encodeWithSerial(), header(), readFrame(), queued, senderQueue (+1 more)

### Community 55 - ".apply"
Cohesion: 0.14
Nodes (10): TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan(), applier, ledgerOps, PlannedBinary (+2 more)

### Community 56 - "FindDaemonProcesses"
Cohesion: 0.13
Nodes (15): applyLiveRefresh(), runRefreshCatalogOnly(), softenAutomatedRefreshErr(), TestApplyLiveRefreshDeliversReload(), TestSoftenAutomatedRefreshErr(), FindDaemonProcesses(), isDaemonCmdline(), TestFindDaemonProcesses() (+7 more)

### Community 57 - "mutedText"
Cohesion: 0.07
Nodes (45): checkGrid, filesPage, panel, recordForm, TrustRow, newLogo(), newBinaryList(), newBinaryRow() (+37 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (53): engine, portOwner, refusedReplacements, resInfo, taintOwner, RegisterBusDir(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+45 more)

### Community 59 - "go_pkg_strings"
Cohesion: 0.06
Nodes (20): parseReservedPorts(), readReservedPorts(), TestParseReservedPorts(), atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget() (+12 more)

### Community 60 - "runMonitor"
Cohesion: 0.19
Nodes (14): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+6 more)

### Community 62 - "runView"
Cohesion: 0.07
Nodes (10): daemonPage, daemonState, logGrid, runHost, runView, readDaemonState(), newLogGrid(), captureLogs() (+2 more)

### Community 63 - "LaunchRule"
Cohesion: 0.12
Nodes (20): allowedEnv(), binaryLaunchRules(), launchRulesByPath(), TestAllowedEnvKeepsOnlyRootPlaced(), TestLaunchRulesByPath(), argvEntry, applyLaunchRule(), envName() (+12 more)

### Community 64 - "policyEditor"
Cohesion: 0.13
Nodes (13): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), newNetGuardTab() (+5 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "IntegrationSuite"
Cohesion: 0.06
Nodes (29): exploitTest, IntegrationSuite, newEventTypes(), IntegrationSuite, detectWM(), reachable(), centered(), Detect() (+21 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.21
Nodes (6): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent()

### Community 68 - "Recorder"
Cohesion: 0.21
Nodes (12): Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start(), TestMonitorTargetsNameFilesByAbsolutePath(), TestProgressNamesMulticallApplets() (+4 more)

### Community 69 - "os.FileMode"
Cohesion: 0.13
Nodes (21): installFileAs(), installUnit(), readAtNoFollow(), upsertFile(), upsertUnitAt(), writeAtNoFollow(), isInstallerSSHAgentUnit(), TestIsInstallerSSHAgentUnit() (+13 more)

### Community 70 - "Vault"
Cohesion: 0.06
Nodes (42): TestDiffMergePreservesExistingAndParses(), encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption() (+34 more)

### Community 71 - "trustManager"
Cohesion: 0.11
Nodes (20): inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree() (+12 more)

### Community 72 - "planUpdaters"
Cohesion: 0.06
Nodes (34): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+26 more)

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "find_test.go"
Cohesion: 0.18
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 75 - "update_test.go"
Cohesion: 0.10
Nodes (17): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFilterChannel(), testKeyPair(), TestNewerThanInstalled(), TestNewerThanInstalledAcrossTagFormats() (+9 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "recordFlow"
Cohesion: 0.13
Nodes (17): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+9 more)

### Community 78 - "newTestVetter"
Cohesion: 0.29
Nodes (11): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+3 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.21
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "os.File"
Cohesion: 0.04
Nodes (75): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), scriptIdentity(), placedKey(), startTrustGuard(), trustStartupError(), confirmUserPlaced() (+67 more)

### Community 81 - "Binary"
Cohesion: 0.13
Nodes (20): formatElapsed(), lines(), orNone(), pickBinaries(), pickLibs(), printBinary(), printRecording(), reviewRecording() (+12 more)

### Community 83 - "rcUser"
Cohesion: 0.07
Nodes (66): setupBunTmpdir(), addedResources(), TestAddedResourcesSkipsAlreadyGuardedSSH(), installFacultatives(), listFacultatives(), runFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative() (+58 more)

### Community 84 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.11
Nodes (23): placedWalk, treeJudge, confineBelow(), inUserHome(), PathUnder(), vouchedFd(), GlobSystemPlaced(), placedDir() (+15 more)

### Community 85 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 86 - "exe_supersede.h"
Cohesion: 0.22
Nodes (14): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+6 more)

### Community 87 - "loginServer"
Cohesion: 0.07
Nodes (25): markReadyWithLogin(), hasSentinel(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel(), runtimePathsFor(), startLoginServer(), loggedInUsers() (+17 more)

### Community 88 - "gui/session.go"
Cohesion: 0.23
Nodes (13): ancestorEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestOwnedByChecksTypeAndOwner(), TestReadOwnedFileRefusesSymlink() (+5 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "Monitor"
Cohesion: 0.15
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 92 - "Vault"
Cohesion: 0.07
Nodes (37): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+29 more)

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

### Community 95 - "NetGuard"
Cohesion: 0.10
Nodes (10): eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), hook (+2 more)

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
Cohesion: 0.12
Nodes (15): procStats, statsSample, duration(), formatAgo(), cpuTime(), humanBytes(), rssBytes(), New() (+7 more)

### Community 100 - "SecretStoreNotes"
Cohesion: 0.39
Nodes (7): dockerKeyringHelper(), gitCredentialHelpers(), readOwnedRegular(), SecretStoreNotes(), TestReadOwnedRegular_RefusesSymlinkAndOtherOwner(), TestSecretStoreNotes(), WarnSecretStores()

### Community 101 - "Gate"
Cohesion: 0.15
Nodes (6): Caller, Gate, Role, Listen(), Bus, Judge

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "NetGuardEvent"
Cohesion: 0.11
Nodes (11): snapNetGuard, NetGuardEvent, NetworkGuardRepository, appendCapped(), protoString(), NetEngines, netGuardEventLine, netGuardEventMsg (+3 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "expandPlaceholders"
Cohesion: 0.10
Nodes (15): reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), expandPlaceholders() (+7 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - ".refreshOnce"
Cohesion: 0.43
Nodes (6): sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad()

### Community 111 - "BinaryEntry"
Cohesion: 0.08
Nodes (40): binaryVerifyState, HashStore, keptEntries, sharedHashEntry, BinaryEntry, checkBtrfsKeys(), ensureBtrfsLayout(), canonicalBinaryPath() (+32 more)

### Community 112 - "Active"
Cohesion: 0.15
Nodes (15): nextConfig(), refreshCatalogConfig(), installConfig(), showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), ConfirmOverwrite() (+7 more)

### Community 113 - "runOfflineEdit"
Cohesion: 0.11
Nodes (18): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), TestPickEncryptedDirSingleSkipsPicker(), chooseResource() (+10 more)

### Community 115 - "maintainerKeeper"
Cohesion: 0.30
Nodes (6): appRunning(), libDirSig(), startMaintainers(), maintainedApp, maintainerKeeper, CleanStages()

### Community 116 - "systemd.go"
Cohesion: 0.11
Nodes (34): installGUI(), buildBinaryIfNeeded(), certifyInstall(), installBinaryOnly(), markInstallComplete(), preflightDeployedBinary(), mustCwd(), revertSystemFiles() (+26 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "NetworkMonitor"
Cohesion: 0.14
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

### Community 126 - "KernelDev"
Cohesion: 0.09
Nodes (24): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, fdStat, BinaryStat, DurableStat (+16 more)

### Community 127 - "classify"
Cohesion: 0.22
Nodes (12): call(), TestClassify(), TestClassifyCreateItemReplace(), TestClassifyOpenCallsExact(), arg(), attributesOf(), classify(), classifyCollection() (+4 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "RunWith"
Cohesion: 0.11
Nodes (28): RunWith(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), deletePostBackups(), listBackups(), restoreBackups(), askDecrypt(), cleanOrphanedMetadata() (+20 more)

### Community 130 - "logLine"
Cohesion: 0.14
Nodes (10): logHook, logLine, journalLine(), TestJournalLine(), logRow(), logStyle(), formatEntry(), sortedKeys() (+2 more)

### Community 131 - "FileEvent"
Cohesion: 0.12
Nodes (6): snapMonitor, FileEvent, MonitorRepository, newFakeMonitorRepo(), fakeSource, fakeMonitorRepo

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - ".judgeProcess"
Cohesion: 0.16
Nodes (12): reasonText(), codeVouched(), isTraced(), loaderEnv(), readMountDevs(), roFor(), TestCodeVouched(), TestLoaderEnv() (+4 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.13
Nodes (15): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+7 more)

### Community 143 - "serve.go"
Cohesion: 0.08
Nodes (31): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+23 more)

### Community 144 - "liveSession"
Cohesion: 0.15
Nodes (14): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), eventsLabel(), forwardBinaries(), holdForward(), normalizedEvents() (+6 more)

### Community 145 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.17
Nodes (7): pos, moveTo(), reposition(), Finder, cursorPos(), Selection, offset()

### Community 146 - "Ledger"
Cohesion: 0.06
Nodes (29): candidates(), confirmAll(), describe(), record(), run(), fyneLogFilter, applySchema(), Ledger (+21 more)

### Community 147 - "EventType"
Cohesion: 0.10
Nodes (18): selfEventsFor(), BpfEvent, EventType, eventUnitTest, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), specFromForm(), specSummary() (+10 more)

### Community 148 - "TempBinary"
Cohesion: 0.05
Nodes (32): ensurePinStateFilePlaceholder(), ensurePlaceholder(), writeInPlace(), bootstrapSelfProtectPlaceholders(), readTempJournal(), stripTempJournal(), TestTempJournalRoundTripInPlace(), writeTempJournal() (+24 more)

### Community 149 - "dialLiveSession"
Cohesion: 0.18
Nodes (14): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig(), runEditConfig(), promptPassword() (+6 more)

### Community 150 - "Daemon configuration"
Cohesion: 0.14
Nodes (13): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Launch rules: `[launch_rules]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Port sections, Process inspectors (+5 more)

### Community 152 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 154 - "searchDirs"
Cohesion: 0.18
Nodes (10): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths(), resolveSoname() (+2 more)

### Community 155 - "startCatalogRefresh"
Cohesion: 0.08
Nodes (22): drain(), newCatalogRefresher(), readEvents(), refreshAdmits(), swapFile(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), addSystemBinary(), catalogPatterns() (+14 more)

### Community 156 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 157 - "dialog_test.go"
Cohesion: 0.39
Nodes (11): pickMode, openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder() (+3 more)

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
Cohesion: 0.11
Nodes (12): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), newNetEventChecks(), FormatAddr(), NetEventTypes() (+4 more)

### Community 163 - "newTestTarget"
Cohesion: 0.22
Nodes (17): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 164 - "absPath"
Cohesion: 0.06
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 165 - "feed.go"
Cohesion: 0.07
Nodes (38): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+30 more)

### Community 166 - "shell"
Cohesion: 0.04
Nodes (58): fakeLink, homePage, pageDef, pageID, reel, resizeWatch, shell, snapEngines (+50 more)

### Community 167 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 168 - "plan_test.go"
Cohesion: 0.20
Nodes (18): isWhiteout(), opaque(), scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent(), TestScanAndValidateAcceptsLibDirContent(), TestScanRefusesSpecialFiles() (+10 more)

### Community 170 - "SelectOrphans"
Cohesion: 0.16
Nodes (20): modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors() (+12 more)

### Community 171 - "batchInodeRows"
Cohesion: 0.20
Nodes (13): markReady(), notifySystemdReady(), inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows() (+5 more)

### Community 172 - "Limitations"
Cohesion: 0.18
Nodes (11): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Secrets served over IPC (+3 more)

### Community 173 - "fetchReleases"
Cohesion: 0.22
Nodes (8): downloadFile(), downloadReleaseFiles(), fetchReleases(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleases(), TestFetchReleasesHTTPError(), progressReader

### Community 174 - "parseDaemonEvents"
Cohesion: 0.15
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite

### Community 175 - "TrustGuard"
Cohesion: 0.06
Nodes (22): setConfinedHomes(), argvMatches, lsmHook, trustHook, trustRows, setSupersedeTrust(), cStr(), syncMap() (+14 more)

### Community 176 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 177 - "debouncer"
Cohesion: 0.33
Nodes (4): TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), debouncer

### Community 178 - "githubRelease"
Cohesion: 0.25
Nodes (17): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+9 more)

### Community 179 - "Run"
Cohesion: 0.33
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 180 - "bufio.Reader"
Cohesion: 0.17
Nodes (14): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+6 more)

### Community 181 - "IntegrationSuite"
Cohesion: 0.23
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 182 - "tempGrantSetup"
Cohesion: 0.23
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

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
Cohesion: 0.29
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "github.com/spf13/cobra.Command"
Cohesion: 0.10
Nodes (33): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+25 more)

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 194 - "CopyTreeWithProgress"
Cohesion: 0.09
Nodes (21): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), CopyTreeWithProgress(), TestCopyTreeRefusesPlantedDestinationSymlink(), TestCopyTreeSourceSwappedAfterClassification() (+13 more)

### Community 195 - "downloadAndVerify"
Cohesion: 0.33
Nodes (10): installHooks(), staged(), downloadAndVerify(), FetchLatest(), fetchPair(), FetchVerified(), Release, sanityCheckBinary() (+2 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 199 - "New"
Cohesion: 0.08
Nodes (47): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault() (+39 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "pinstate_test.go"
Cohesion: 0.13
Nodes (26): configureDaemonLogging(), loadDaemonConfig(), lockOneRoot(), lockRootRecovering(), noteReloadedSettings(), relockStaleVaults(), resolveConfigPath(), runLockdown() (+18 more)

### Community 202 - "DaemonEvent"
Cohesion: 0.12
Nodes (32): ServeConfig, UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+24 more)

### Community 204 - "procstats.go"
Cohesion: 0.30
Nodes (9): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+1 more)

### Community 205 - "sources"
Cohesion: 0.17
Nodes (10): TestGeneratedLibraryBlockRoundTrips(), pathByte(), pathRecord(), TestPathRecord(), archiveMembers(), decompressed(), libDirDev(), tarMembers() (+2 more)

### Community 206 - "monitor"
Cohesion: 0.07
Nodes (33): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks (+25 more)

### Community 209 - "audit.go"
Cohesion: 0.21
Nodes (16): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource() (+8 more)

### Community 210 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "parseGuardEvent"
Cohesion: 0.11
Nodes (16): bpfGuardEvent, BusClient, AttestBusClient(), busClientKey(), execStamp(), imageOf(), pidfdPid(), procIdentity() (+8 more)

### Community 213 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 214 - "run"
Cohesion: 0.27
Nodes (10): fetchRelease(), InstalledGUIVersion(), InstalledVersion(), readInstalledVersion(), run(), runUpdate(), syncGUI(), TestFetchReleaseByTag() (+2 more)

### Community 216 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 217 - "linkView"
Cohesion: 0.29
Nodes (6): linkTarget(), newLinkView(), TestCheckLinksJudgesEachTarget(), entryKind, libPlan, linkView

### Community 218 - "runNetworkGuard"
Cohesion: 0.17
Nodes (15): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), WithSuspectsFrom(), eventsetSummary() (+7 more)

### Community 219 - "refuseOutsideHardlinks"
Cohesion: 0.24
Nodes (8): lstatOf(), outsideHardlinks(), refuseOutsideHardlinks(), TestRefuseOutsideHardlinks_FileRoot(), TestRefuseOutsideHardlinks_InTreeLinksAccepted(), TestRefuseOutsideHardlinks_MissingRootFails(), TestRefuseOutsideHardlinks_OutsideLinkRefused(), writeFile()

### Community 221 - "Use"
Cohesion: 0.22
Nodes (8): channelOf(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestRunWithRefusesConcurrentRun(), cancelUI, TestEditTextMapsUICancel(), Use()

### Community 222 - "serverAuth"
Cohesion: 0.32
Nodes (5): TestClientAuth(), TestSASL(), clientAuth(), serverAuth(), saslServer

### Community 223 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 224 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 228 - "mountTable"
Cohesion: 0.15
Nodes (8): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), ReadMountinfo(), unescapeMountinfo(), mountTable, editorHarness

### Community 229 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

### Community 231 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 234 - "ValidateTargets"
Cohesion: 0.12
Nodes (14): drain(), dedupe(), existingSection(), Plan, guardedMatch(), libOfChosen(), overlaps(), refusedPath() (+6 more)

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "deploy"
Cohesion: 0.14
Nodes (20): archAssets(), assetsFor(), checkUpdatePreconditions(), confirmUnit(), deploy(), deployGUI(), fetchGUI(), guiAssetsFor() (+12 more)

### Community 238 - "GenerateConf"
Cohesion: 0.05
Nodes (72): logRefreshChange(), TestLogRefreshChange(), groupedDiscordConf(), TestGroupedSectionAddressedByEncryptionRoot(), TestPatchCatalogSectionGroupedConfig(), TestSetSectionWhitelistPreservesGroupStructure(), LibraryBlock, RefreshOptions (+64 more)

### Community 242 - "User"
Cohesion: 0.06
Nodes (53): entryWriters(), maintainedApps(), ownedBy(), sameDevice(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers() (+45 more)

### Community 245 - "install/install.go"
Cohesion: 0.10
Nodes (32): asRoot(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), ensureInstalledBinary(), ensureMasterKey(), fileSHA256(), prepareInstallation(), runInstall() (+24 more)

### Community 246 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 251 - "monitor.md"
Cohesion: 0.50
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 257 - "network monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 267 - "daemonconfig.go"
Cohesion: 0.08
Nodes (71): systemPatterns(), systemRule(), ensureAnchor(), ensurePortAnchors(), ensureRootDir(), ephemeralRange(), ensureSecretAnchors(), watchGroup (+63 more)

### Community 273 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 294 - "go_pkg_os"
Cohesion: 0.06
Nodes (28): releaseBuild(), Run(), needsBinary(), TestNeedsBinary(), removeSessionHooks(), closeDumpFile(), Execute(), initLogger() (+20 more)

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 665 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **43 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `parseDaemonEvents`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`, `TestMain`, `go_pkg_strings`, `IntegrationSuite`, `.netIntegrityContainer`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.015872112868358174 - nodes in this community are weakly interconnected._
- **Why does `monitor()` connect `monitor` to `app-listener documentation`, `time.Duration`, `secret_probe/main.go`, `Installation`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0943859649122807 - nodes in this community are weakly interconnected._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `time.Duration`, `guardUnitTest`, `time.Time`, `DaemonEvent`, `GuardEvent`, `daemonconfig.go`, `buildOneGuard`, `BinaryEntry`, `os.File`, `TempBinary`, `Resource`, `loginServer`, `sync.Mutex`, `GuardInodeKey`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **Should `daemonconfig_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06442577030812324 - nodes in this community are weakly interconnected._