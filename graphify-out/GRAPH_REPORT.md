# Graph Report - app-listener  (2026-10-30)

## Corpus Check
- 544 files · ~1,768,084 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7053 nodes · 31794 edges · 228 communities (190 shown, 38 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2918 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9aad6637`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- go_pkg_strings
- testing.T
- IntegrationSuite
- Vault
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- Candidate
- CandidateDir
- protectedPage
- filevault.go
- wm.go
- github.com/charmbracelet/bubbletea.Cmd
- IntegrationSuite
- guard_trust.bpf.c
- tabs.go
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- Finder
- guard.bpf.c
- NetworkMonitor
- changelog.go
- install/install.go
- LaunchRule
- Binary
- TestMonitorUseCaseLifecycle
- sync.Mutex
- Guard
- picker
- desktop.go
- IntegrationSuite
- Message
- controlServer
- github.com/charmbracelet/bubbles/textarea.Model
- pinstate.go
- IntegrationSuite
- absPath
- monitorUnitTest
- FileEvent
- fileedit_model_test.go
- GuardEvent
- stdio
- update_test.go
- fadeScroll
- go_pkg_path_filepath
- editView
- Client
- daemonconfig.go
- Network
- .apply
- githubRelease
- SanitizeText
- engine
- Monitor
- installPage
- Recorder
- exe_supersede.h
- updatePage
- UI
- unix_connect.c
- context.Context
- IntegrationSuite
- buildTrustedSet
- liveSession
- catalogRefresher
- RefreshSection
- KernelDev
- fcntl
- bpfstats/main.go
- binaryVetter
- codecheck/main.go
- gui/session.go
- EventType
- process_vm_readv.c
- os.File
- maintainer/plan.go
- LibraryClosure
- shellrc.go
- vaultfs.go
- NetGuardEvent
- runBackend
- startCatalogRefresh
- GenerateConf
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- runView
- Active
- prctl_exe_swap.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- time.Time
- Plan
- Well-known bypass classes
- New
- jit_provenance.c
- inspector_probe.c
- recordops.c
- Vault
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- deploy
- BinaryEntry
- newTestVetter
- promptPassword
- IntegrationSuite
- systemd.go
- github.com/spf13/cobra.Command
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- parseDaemonEvents
- Config
- WebSocket /ws client
- safeio.go
- install.sh
- mc_tag.h
- .filterEvent
- make test-integration (rootful Docker suite)
- IntegrationSuite
- fileedit_symlink_test.go
- RawTarget
- trace-app-libs.sh
- IntegrationSuite
- maintainerKeeper
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- mustSubkey
- Serve
- NetEventType
- procstats.go
- Ledger
- planUpdaters
- runNetworkGuard
- GuardInodeKey
- changelog_test.go
- fakeStore
- testLimiter
- IntegrationSuite
- buildOneGuard
- refreshAdmits
- CopyTreeWithProgress
- cursor.go
- Installation
- backend_test.go
- graphify-refresh.sh
- runMonitor
- .exeKeys
- newTestTarget
- lockAndDeprovision
- networkPage
- shell
- layout.go
- watchSet
- Question
- fscrypt/orphans.go
- Limitations
- net_tester/main.go
- IntegrationSuite
- TrustGuard
- Choice
- tempGrantSetup
- Daemon configuration
- .IsRuntime
- Compatibility
- How it works
- inodeChain
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- recordFlow
- .populateWatchInodes
- edit-protected
- TestCleanRel
- kernel-prof.sh
- fakeDaemon
- pathCache
- NewEventFanout
- DaemonEvent
- io.Writer
- hashCache
- runLockdown
- Desktop GUI
- IntegrationSuite
- LiveEmptyWhitelistRejected
- parseGuardEvent
- .eventTypeAllowed
- loadDaemonConfig
- netinject.c
- IntegrationSuite
- .awaitDaemonUpOrDead
- VouchedDevs
- shQuote
- downloadAndVerify
- configedit_test.go
- writeWithin
- supersede_probe.c
- mountTable
- github.com/Virgula0/app-listener/tools/codecheck
- inode_dev.h
- Troubleshooting
- ptrace_seize.c
- monitor.md

## God Nodes (most connected - your core abstractions)
1. `Config` - 134 edges
2. `Guard` - 124 edges
3. `GuardInodeKey` - 105 edges
4. `shell` - 97 edges
5. `Load()` - 95 edges
6. `absPath()` - 92 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 74 edges
10. `IntegrationSuite` - 67 edges

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

## Communities (228 total, 38 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.09
Nodes (33): tempJournalRow, Progress, User, classCacheKey, launchRule, MulticallError, trustEvent, textAnswer (+25 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (163): plantBun(), TestOpenBunTmp_KeepsVettedDir(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestEveryInstallModeMapped() (+155 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (7): configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig()

### Community 3 - "Vault"
Cohesion: 0.19
Nodes (12): isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn(), readKey(), TestNewBoundedKeyFnFirstCall() (+4 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (78): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+70 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.09
Nodes (51): prctl_gate(), guard_task_prctl(), admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader() (+43 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (23): addErrKind, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+15 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "Candidate"
Cohesion: 0.19
Nodes (17): appendSectionsAndEdit(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories(), pickFromCandidates() (+9 more)

### Community 11 - "CandidateDir"
Cohesion: 0.04
Nodes (52): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, GlobChild, compileGlobNames() (+44 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (33): backendLink, checkGrid, configView, forwardForm, forwardView, installUI, protectedPage, recordReview (+25 more)

### Community 13 - "filevault.go"
Cohesion: 0.24
Nodes (15): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+7 more)

### Community 14 - "wm.go"
Cohesion: 0.19
Nodes (17): invokingIDs(), invokingUID(), detectWM(), Detect(), Session, hyprlandSocket(), i3RuntimeSocket(), socketOwnedBy() (+9 more)

### Community 15 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.04
Nodes (48): recordTickCmd(), recordModel, diffModel, editorModel, formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision() (+40 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (64): argv_env(), argv_matched(), argv_rule_at(), argv_word(), caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key() (+56 more)

### Community 18 - "tabs.go"
Cohesion: 0.03
Nodes (24): appTheme, bleedLayout, clickPanel, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer (+16 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (48): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+40 more)

### Community 22 - "Resource"
Cohesion: 0.06
Nodes (12): ownerPaths(), Resource, GuardRepository, Vault, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, runtimeReleased() (+4 more)

### Community 23 - "auth_test.go"
Cohesion: 0.12
Nodes (41): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+33 more)

### Community 24 - "Finder"
Cohesion: 0.09
Nodes (30): findButton, span, column(), findAll(), Finder, lowerRunes(), NewFinder(), newFindInput() (+22 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.08
Nodes (78): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+70 more)

### Community 26 - "NetworkMonitor"
Cohesion: 0.14
Nodes (7): BinariesSummary(), TestBinariesSummary(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 27 - "changelog.go"
Cohesion: 0.09
Nodes (40): class, classOf(), dedent(), Edit(), leadingIndent(), New(), onlySpaces(), repeatKey() (+32 more)

### Community 28 - "install/install.go"
Cohesion: 0.04
Nodes (86): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), bunLaunchersOffered() (+78 more)

### Community 29 - "LaunchRule"
Cohesion: 0.07
Nodes (37): allowedEnv(), binaryLaunchRules(), launchRulesByPath(), TestAllowedEnvKeepsOnlyRootPlaced(), TestLaunchRulesByPath(), argvEntry, placedWalk, treeJudge (+29 more)

### Community 30 - "Binary"
Cohesion: 0.19
Nodes (15): lines(), pickBinaries(), pickLibs(), reviewRecording(), beyondEvents(), Classify(), finishBinary(), Binary (+7 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 32 - "sync.Mutex"
Cohesion: 0.05
Nodes (32): binaryList, filesPage, netGuardTab, policyEditor, policyHooks, policyTab, Target, newBinaryList() (+24 more)

### Community 33 - "Guard"
Cohesion: 0.05
Nodes (21): deferredBinary, exeRow, rootHandle, supersededBinary, tempMaskOp, tempRow, openRoot(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted() (+13 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (30): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+22 more)

### Community 35 - "desktop.go"
Cohesion: 0.04
Nodes (107): installHooks(), releaseBuild(), Run(), staged(), TestPickEncryptedDirSingleSkipsPicker(), pickDir(), pickEncryptedDir(), runOfflineEdit() (+99 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (13): checkSecret(), trustRows(), File, TrustItem, backend, backendOp, backendState, BinaryCheck (+5 more)

### Community 38 - "controlServer"
Cohesion: 0.08
Nodes (26): readConfigPut(), TestReadGrantRequestConfig(), authSelfPeer(), newControlManager(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest() (+18 more)

### Community 39 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.10
Nodes (21): lineStyles, pos, viewportState, lines(), moveTo(), Navigate(), reposition(), baseStyle() (+13 more)

### Community 40 - "pinstate.go"
Cohesion: 0.18
Nodes (22): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+14 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.08
Nodes (14): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest, IntegrationSuite (+6 more)

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (7): IntegrationSuite, IntegrationSuite, IntegrationSuite, absPath(), TestMain(), helperChild(), TestMain()

### Community 44 - "FileEvent"
Cohesion: 0.12
Nodes (6): snapMonitor, FileEvent, MonitorRepository, newFakeMonitorRepo(), fakeSource, fakeMonitorRepo

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "GuardEvent"
Cohesion: 0.04
Nodes (26): snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), commMatchesGuardedBinary(), GuardEvent, logBacklogDenial(), TempBinary, TemporaryGrant (+18 more)

### Community 48 - "update_test.go"
Cohesion: 0.10
Nodes (17): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFilterChannel(), testKeyPair(), TestNewerThanInstalled(), TestNewerThanInstalledAcrossTagFormats() (+9 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (29): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+21 more)

### Community 51 - "editView"
Cohesion: 0.12
Nodes (7): codeEntry, editView, Entry, fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 52 - "Client"
Cohesion: 0.08
Nodes (16): TestBackendServeNeedsHelloFirst(), Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend() (+8 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.11
Nodes (46): watchGroup, addResource(), applyAllowLib(), applyAnyGroupDirective(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector() (+38 more)

### Community 54 - "Network"
Cohesion: 0.12
Nodes (8): netGuardPill(), NetworkGuardRepository, clonePolicy(), NetEngines, NetExtra, NetPolicy, Network, NetGuardConfig

### Community 55 - ".apply"
Cohesion: 0.15
Nodes (8): Apply(), Outcome, ProgressFunc, applier, ledgerOps, PlannedBinary, Store, vaultOps

### Community 56 - "githubRelease"
Cohesion: 0.25
Nodes (17): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+9 more)

### Community 57 - "SanitizeText"
Cohesion: 0.07
Nodes (49): column, feedAction, feedSpec, multiAnswer, panel, procSpec, recordForm, newLogo() (+41 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (33): engine, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded() (+25 more)

### Community 59 - "Monitor"
Cohesion: 0.17
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 60 - "installPage"
Cohesion: 0.06
Nodes (16): actionGate, binaryRow, binSource, handButton, handCheck, homePage, iconTile, installAction (+8 more)

### Community 61 - "Recorder"
Cohesion: 0.17
Nodes (14): drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start(), TestMonitorTargetsNameFilesByAbsolutePath() (+6 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.21
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+7 more)

### Community 63 - "updatePage"
Cohesion: 0.28
Nodes (4): latestRelease, updatePage, boolIndex(), isPreRelease()

### Community 64 - "UI"
Cohesion: 0.12
Nodes (18): channelOf(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestRunWithRefusesConcurrentRun(), InstallMode, InstallRequest, Options (+10 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "context.Context"
Cohesion: 0.18
Nodes (10): reachable(), centered(), FitSize(), poll(), TestFitSize(), hyprClient, hyprMonitor, Rect (+2 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "buildTrustedSet"
Cohesion: 0.19
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 69 - "liveSession"
Cohesion: 0.10
Nodes (26): auditAfterEdit(), coveredByWatch(), TestWithin(), within(), dialLiveSession(), parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger() (+18 more)

### Community 70 - "catalogRefresher"
Cohesion: 0.14
Nodes (9): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), catalogRefresher (+1 more)

### Community 71 - "RefreshSection"
Cohesion: 0.09
Nodes (28): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+20 more)

### Community 72 - "KernelDev"
Cohesion: 0.13
Nodes (21): BinaryStat, binaryStatOf(), binaryStatOf(), KernelDev(), StatBinary(), StatBinaryFile(), BtrfsDevices(), btrfsSysfsDevices() (+13 more)

### Community 73 - "fcntl"
Cohesion: 0.08
Nodes (11): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), hold_forever() (+3 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 75 - "binaryVetter"
Cohesion: 0.08
Nodes (25): newBinaryVetter(), openBinaryVetter(), awaitStartupOrSignal(), backingDeviceUnion(), buildGuards(), catchLifecycleSignals(), reloadOnce(), reloadSlotsNeeded() (+17 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "gui/session.go"
Cohesion: 0.24
Nodes (13): ancestorEnv(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue(), TestOwnedByChecksTypeAndOwner(), TestReadOwnedFileRefusesSymlink() (+5 more)

### Community 78 - "EventType"
Cohesion: 0.09
Nodes (20): systemRule(), BpfEvent, EventType, eventUnitTest, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), specFromForm() (+12 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.23
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "os.File"
Cohesion: 0.05
Nodes (65): scriptIdentity(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched() (+57 more)

### Community 81 - "maintainer/plan.go"
Cohesion: 0.05
Nodes (53): TestGeneratedLibraryBlockRoundTrips(), SystemTrusted(), clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist() (+45 more)

### Community 82 - "LibraryClosure"
Cohesion: 0.07
Nodes (31): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), IsScript(), ldConfInclude(), ldSoConfDirs() (+23 more)

### Community 83 - "shellrc.go"
Cohesion: 0.11
Nodes (52): listFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative(), TestListFacultativesIgnoresDaemonConf(), TestListFacultativesMarksGuardedSSH(), removeFacultative(), removeFacultatives(), TestRemoveSSHAgentFacultativeKeepsGuardedSSHConfig() (+44 more)

### Community 84 - "vaultfs.go"
Cohesion: 0.14
Nodes (19): DefaultNewFileMode(), IsBinaryContent(), IsOctalMode(), ParseMode(), inodeOf(), mustRead(), mustWrite(), openRootT() (+11 more)

### Community 85 - "NetGuardEvent"
Cohesion: 0.10
Nodes (9): snapNetGuard, reasonText(), NetGuardEvent, TestReasonText(), fakeGuard, fakeMonitor, fakeNet, fakeNetGuard (+1 more)

### Community 86 - "runBackend"
Cohesion: 0.36
Nodes (6): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive()

### Community 87 - "startCatalogRefresh"
Cohesion: 0.11
Nodes (17): markReady(), markReadyWithLogin(), notifySystemdReady(), startCatalogRefresh(), hasSentinel(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel() (+9 more)

### Community 88 - "GenerateConf"
Cohesion: 0.06
Nodes (56): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), discordCatalogEntry(), groupedDiscordConf(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), TestGroupedSectionAddressedByEncryptionRoot() (+48 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "runView"
Cohesion: 0.05
Nodes (19): daemonPage, logGrid, logHook, logLine, runHost, runView, stageItem, stageState (+11 more)

### Community 92 - "Active"
Cohesion: 0.18
Nodes (12): nextConfig(), showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), ShowDiff(), EditText(), Active() (+4 more)

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.16
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.08
Nodes (19): foreignPinnedLink(), NewMountVouch(), PumpRingbuf(), eventTypeKey(), codeVouched(), isTraced(), judgePasses(), loaderEnv() (+11 more)

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
Cohesion: 0.15
Nodes (13): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+5 more)

### Community 100 - "time.Time"
Cohesion: 0.24
Nodes (5): LockedError, newWatchSet(), fakeInotify, sweepID, sweepMemo

### Community 101 - "Plan"
Cohesion: 0.20
Nodes (6): dedupe(), Plan, libOfChosen(), sectionFor(), Existing, ExistingRule

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "New"
Cohesion: 0.15
Nodes (21): askEncryption(), TestAskEncryptionSkipsNeedEncryptionFalse(), decryptDirectories(), TestDecryptDirectoriesEmpty(), harnessSuite, looksLikeFileVaultRecord(), inode(), TestIsEncryptedDispatchesToFileVaultForRegularFiles() (+13 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.20
Nodes (13): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+5 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "deploy"
Cohesion: 0.11
Nodes (27): archAssets(), assetsFor(), checkUpdatePreconditions(), confirmUnit(), deploy(), deployGUI(), fetchGUI(), guiAssetsFor() (+19 more)

### Community 111 - "BinaryEntry"
Cohesion: 0.08
Nodes (38): binaryVerifyState, HashStore, keptEntries, sharedHashEntry, BinaryEntry, checkBtrfsKeys(), canonicalBinaryPath(), ConfinedEntry() (+30 more)

### Community 112 - "newTestVetter"
Cohesion: 0.44
Nodes (10): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+2 more)

### Community 113 - "promptPassword"
Cohesion: 0.22
Nodes (11): applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig(), runEditConfig(), promptPassword(), promptPasswordWith() (+3 more)

### Community 115 - "systemd.go"
Cohesion: 0.05
Nodes (61): anyDaemonRunning(), offerDaemonStart(), applyLiveRefresh(), buildBinaryIfNeeded(), installBinaryOnly(), runRefreshCatalogOnly(), softenAutomatedRefreshErr(), confirmRunPrereq() (+53 more)

### Community 116 - "github.com/spf13/cobra.Command"
Cohesion: 0.09
Nodes (36): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+28 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "parseDaemonEvents"
Cohesion: 0.15
Nodes (8): IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite, IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "Config"
Cohesion: 0.07
Nodes (48): makeReloadHandler(), noteReloadedSettings(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), startTrustGuard(), trustStartupError(), auditFindings(), auditFindingsWithConfig() (+40 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "safeio.go"
Cohesion: 0.20
Nodes (16): entryOf(), writeUserFile(), createOwned(), removeCheckedFile(), AtomicWriteAt(), CreateExclNoFollow(), descend(), DescendCreateOwned() (+8 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "IntegrationSuite"
Cohesion: 0.24
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "fileedit_symlink_test.go"
Cohesion: 0.36
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 131 - "RawTarget"
Cohesion: 0.31
Nodes (9): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+1 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "maintainerKeeper"
Cohesion: 0.18
Nodes (12): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), startMaintainers(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), maintainedApp (+4 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "mustSubkey"
Cohesion: 0.31
Nodes (9): newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape(), TestOpenFileVaultRejectsTampering(), TestSealFileVaultNoncesNeverRepeat() (+1 more)

### Community 143 - "Serve"
Cohesion: 0.08
Nodes (29): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+21 more)

### Community 144 - "NetEventType"
Cohesion: 0.08
Nodes (17): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), newNetEventChecks(), FormatAddr(), NetEventTypes() (+9 more)

### Community 145 - "procstats.go"
Cohesion: 0.30
Nodes (9): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+1 more)

### Community 146 - "Ledger"
Cohesion: 0.07
Nodes (28): checkTrustLine(), candidates(), confirmAll(), describe(), record(), run(), applySchema(), Ledger (+20 more)

### Community 147 - "planUpdaters"
Cohesion: 0.14
Nodes (18): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+10 more)

### Community 148 - "runNetworkGuard"
Cohesion: 0.12
Nodes (20): CheckBPFLSM(), confirmMasterKeyOverwrite(), runBPFCheck(), runGenKey(), runOneShotMode(), computeGuardBinaries(), confirmUnsafe(), discoverInfra() (+12 more)

### Community 149 - "GuardInodeKey"
Cohesion: 0.06
Nodes (39): inodeBatchFn, refusedReplacements, ensureBtrfsLayout(), appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows() (+31 more)

### Community 150 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 152 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 154 - "buildOneGuard"
Cohesion: 0.06
Nodes (46): multicallAdmissible(), buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards() (+38 more)

### Community 155 - "refreshAdmits"
Cohesion: 0.33
Nodes (6): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), BinaryRule, materialized(), Extension

### Community 156 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 157 - "cursor.go"
Cohesion: 0.08
Nodes (22): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+14 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.24
Nodes (21): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+13 more)

### Community 161 - "runMonitor"
Cohesion: 0.10
Nodes (21): DrainUntilSignal(), CheckEBPF(), MakeDisplayPaths(), ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), newPinGeneration(), prepareDaemonStart() (+13 more)

### Community 163 - "newTestTarget"
Cohesion: 0.22
Nodes (17): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 164 - "lockAndDeprovision"
Cohesion: 0.09
Nodes (17): auditEntry(), auditTree(), relOrPath(), TestAuditTree(), deprovisionKind, migrationTarget, writeTempWithMetadata(), classifyDeprovision() (+9 more)

### Community 165 - "networkPage"
Cohesion: 0.06
Nodes (34): actionBar, cell, counterStrip, dataTable, eventFeed, networkPage, processView, targetView (+26 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (52): fakeLink, pageDef, pageID, reel, resizeWatch, shell, sizedDialog, snapEngines (+44 more)

### Community 167 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 168 - "watchSet"
Cohesion: 0.18
Nodes (11): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchDir, watchPattern (+3 more)

### Community 169 - "Question"
Cohesion: 0.40
Nodes (3): Confirm(), Question, recordUI

### Community 170 - "fscrypt/orphans.go"
Cohesion: 0.19
Nodes (20): modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors() (+12 more)

### Community 172 - "Limitations"
Cohesion: 0.18
Nodes (11): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Secrets served over IPC (+3 more)

### Community 173 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 175 - "TrustGuard"
Cohesion: 0.05
Nodes (23): lsmHook, statFunc, tempGrant, trustHook, trustRows, Multicall, MulticallKind, multicallFamily() (+15 more)

### Community 176 - "Choice"
Cohesion: 0.50
Nodes (3): MultiSelectKeymap(), Choice, MultiSelect()

### Community 177 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 178 - "Daemon configuration"
Cohesion: 0.20
Nodes (10): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Launch rules: `[launch_rules]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Runtime sections (+2 more)

### Community 179 - ".IsRuntime"
Cohesion: 0.50
Nodes (3): isAbsoluteCandidate(), TestCatalogSanity(), TestRuntimeEntriesNeverDeleteRoot()

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "inodeChain"
Cohesion: 0.23
Nodes (15): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD(), keyOf() (+7 more)

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

### Community 191 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 194 - "recordFlow"
Cohesion: 0.14
Nodes (14): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+6 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "NewEventFanout"
Cohesion: 0.24
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 202 - "DaemonEvent"
Cohesion: 0.10
Nodes (29): ServeConfig, UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+21 more)

### Community 203 - "io.Writer"
Cohesion: 0.13
Nodes (11): formatElapsed(), orNone(), printBinary(), printRecording(), TestProgressReader(), fyneLogFilter, bottomBar, newBottomBar() (+3 more)

### Community 204 - "hashCache"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 205 - "runLockdown"
Cohesion: 0.14
Nodes (16): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), resolveConfigPath(), runLockdown(), dirIsRootOwnedSafe(), isBpffs() (+8 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.13
Nodes (15): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+7 more)

### Community 210 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "loadDaemonConfig"
Cohesion: 0.29
Nodes (7): loadDaemonConfig(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), SetUserHomes(), TestResolveConfinedRootOwnedParents(), TestOpenSystemPlacedRefusesUserHome()

### Community 213 - "netinject.c"
Cohesion: 0.32
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 217 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 220 - "shQuote"
Cohesion: 0.16
Nodes (6): IntegrationSuite, IntegrationSuite, overlayPlanter(), shQuote(), inNosuidNamespace(), IntegrationSuite

### Community 221 - "downloadAndVerify"
Cohesion: 0.18
Nodes (15): downloadAndVerify(), downloadFile(), downloadReleaseFiles(), FetchLatest(), fetchPair(), fetchRelease(), fetchReleases(), FetchVerified() (+7 more)

### Community 223 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 224 - "writeWithin"
Cohesion: 0.25
Nodes (8): atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync(), writeWithin()

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 228 - "mountTable"
Cohesion: 0.13
Nodes (10): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo(), userMountable() (+2 more)

### Community 233 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 241 - "ptrace_seize.c"
Cohesion: 0.50
Nodes (6): main(), peek_buffer(), run_tracer(), run_victim(), touch(), wait_for_file()

### Community 244 - "monitor.md"
Cohesion: 0.50
Nodes (3): Events, Flags, monitor: see what touches a path

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 651 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **38 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `backend` connect `Message` to `sync.Mutex`, `recordFlow`, `time.Duration`, `Vault`, `TestCleanRel`, `liveSession`, `fakeDaemon`, `Vault`, `go_pkg_path_filepath`, `Ledger`, `Client`, `runBackend`, `backend_test.go`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.09425499725535874 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `go_pkg_strings`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `parseDaemonEvents`, `IntegrationSuite`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.017956656346749224 - nodes in this community are weakly interconnected._
- **Why does `shell` connect `shell` to `sync.Mutex`, `UI`, `picker`, `time.Duration`, `go_pkg_strings`, `networkPage`, `protectedPage`, `wm.go`, `tabs.go`, `Network`, `SanitizeText`, `runView`, `installPage`, `cursor.go`, `updatePage`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09441441441441441 - nodes in this community are weakly interconnected._