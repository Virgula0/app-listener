# Graph Report - app-listener  (2026-10-30)

## Corpus Check
- 562 files · ~1,787,552 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7299 nodes · 33359 edges · 237 communities (198 shown, 39 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2992 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `78d3fac3`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- go_pkg_strings
- testing.T
- IntegrationSuite
- filevault.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- Open
- Config
- protectedPage
- GuardInodeKey
- Network
- github.com/charmbracelet/bubbletea.Cmd
- IntegrationSuite
- guard_trust.bpf.c
- navItem
- fileEditModel
- networkmonitor.bpf.c
- newFakeVault
- Resource
- setpassword_authz_test.go
- GenerateConf
- guard.bpf.c
- SanitizeText
- install/install.go
- plan_test.go
- desktop_test.go
- recordFlow
- TestMonitorUseCaseLifecycle
- runNetworkMonitor
- Guard
- picker
- secret_probe/main.go
- IntegrationSuite
- Message
- controlServer
- render.go
- pinstate.go
- IntegrationSuite
- absPath
- monitorUnitTest
- fyne.io/fyne/v2/driver/desktop.Cursor
- planUpdaters
- fakeGuardRepo
- stdio
- update_test.go
- fadeScroll
- go_pkg_path_filepath
- install_ui.go
- proto.go
- watchGroup
- sync.Mutex
- .apply
- changelog.go
- mutedText
- engine
- Monitor
- daemon_multicall_test.go
- .FilterExistingWhitelist
- inode_dev
- .filterEvent
- FileTarget
- unix_connect.c
- wm.go
- IntegrationSuite
- liveSession
- runForward
- time.Time
- RefreshCatalog
- NetGuardEvent
- fcntl
- runMonitor
- LaunchRule
- codecheck/main.go
- runNetworkGuard
- CandidateDir
- process_vm_readv.c
- os.File
- Plan
- LibraryClosure
- shellrc.go
- golang.org/x/sys/unix.Stat_t
- bpfstats/main.go
- EventType
- loginServer
- RefreshSection
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- runView
- highlight_test.go
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- BinaryEntry
- auditGroupCoverage
- Well-known bypass classes
- SampleContent
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- RunCmd
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- github.com/charmbracelet/bubbles/textarea.Model
- ConfinedEntry
- installPage
- .netIntegrityContainer
- IntegrationSuite
- protected_test.go
- fadescroll_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- UI
- maintainerKeeper
- WebSocket /ws client
- inet_connect.c
- install.sh
- mc_tag.h
- classify
- make test-integration (rootful Docker suite)
- IntegrationSuite
- catalogRefresher
- runBackend
- trace-app-libs.sh
- IntegrationSuite
- selection.go
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- vaultfs.go
- serve.go
- NetEventType
- find.go
- Ledger
- newTestTarget
- GuardEvent
- buildTrustedSet
- Daemon configuration
- fakeStore
- debouncer
- IntegrationSuite
- buildOneGuard
- configedit_test.go
- CopyTreeWithProgress
- IntegrationSuite
- Installation
- backend_test.go
- graphify-refresh.sh
- serverAuth
- WithSuspectsFrom
- RootPlacedName
- .populateWatchInodes
- feed.go
- shell
- .Layout
- boundedGlob
- changelog_test.go
- fscrypt/orphans.go
- Limitations
- RunWith
- IntegrationSuite
- TrustGuard
- net_tester/main.go
- ValidateTargets
- deploy
- IntegrationSuite
- prctl_gate
- How it works
- inodeChain
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- github.com/spf13/cobra.Command
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- probeSystem
- maintainer/plan.go
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- migrate_test.go
- pathCache
- Vault
- DaemonEvent
- WithBottomBar
- FileEvent
- eventUnitTest
- monitor
- IntegrationSuite
- main_test.go
- Run
- .eventTypeAllowed
- KernelDev
- parseKernelRelease
- fileedit_symlink_test.go
- .awaitDaemonUpOrDead
- .TestRecord_GuardParity
- controlManager
- parseDaemonEvents
- downloadAndVerify
- NewEventFanout
- linkView
- supersede_probe.c
- guardModel
- mountTable
- github.com/Virgula0/app-listener/tools/codecheck
- Troubleshooting
- Candidate
- shQuote
- taint_victim.c
- app-listener documentation
- monitor: see what touches a path
- bleedLayout
- sources_test.go
- parsePasswd

## God Nodes (most connected - your core abstractions)
1. `Config` - 146 edges
2. `Guard` - 125 edges
3. `GuardInodeKey` - 112 edges
4. `shell` - 97 edges
5. `Load()` - 95 edges
6. `absPath()` - 93 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 74 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `eBPF at three levels` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go
- `Why not an AppArmor profile?` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go
- `Modes` --references--> `monitor()`  [INFERRED]
  docs/README.md → integrationtests/secret_probe/main.go
- `Distributions` --references--> `monitor()`  [INFERRED]
  docs/compatibility.md → integrationtests/secret_probe/main.go
- `Files` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go

## Import Cycles
- None detected.

## Communities (237 total, 39 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.12
Nodes (15): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey, fileID (+7 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (125): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), awaitStartupOrSignal(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestAwaitStartupOrSignalAbortsAndLocksDown(), TestAwaitStartupOrSignalAbortWithFailedStartup() (+117 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (8): Hash(), configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig()

### Community 3 - "filevault.go"
Cohesion: 0.06
Nodes (47): deprovisionKind, migrationTarget, classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), isFileVaultCiphertext() (+39 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.06
Nodes (84): TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+76 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (28): addErrKind, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+20 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (36): TestOpenBunTmp_SymlinkReplacedTargetUntouched(), newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds() (+28 more)

### Community 10 - "Open"
Cohesion: 0.14
Nodes (22): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+14 more)

### Community 11 - "Config"
Cohesion: 0.06
Nodes (63): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), addedResources(), applyDiffAdditions() (+55 more)

### Community 12 - "protectedPage"
Cohesion: 0.04
Nodes (25): backendLink, checkGrid, codeEntry, configView, editView, forwardForm, forwardView, handButton (+17 more)

### Community 13 - "GuardInodeKey"
Cohesion: 0.06
Nodes (41): inodeBatchFn, refusedReplacements, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows(), walkInodeBatches() (+33 more)

### Community 14 - "Network"
Cohesion: 0.17
Nodes (5): netGuardPill(), clonePolicy(), NetExtra, NetPolicy, Network

### Community 15 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.07
Nodes (17): recordTickCmd(), recordModel, diffModel, editorModel, listenForDaemonEvents(), tickDaemonStats(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast() (+9 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (61): argv_env(), argv_matched(), argv_rule_at(), argv_word(), caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key() (+53 more)

### Community 18 - "navItem"
Cohesion: 0.03
Nodes (23): appTheme, binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panelRenderer (+15 more)

### Community 19 - "fileEditModel"
Cohesion: 0.06
Nodes (37): newFindInput(), clipLabel(), humanSize(), listEntries(), clamp(), newFileEditModel(), paneTitle(), sortNodes() (+29 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "newFakeVault"
Cohesion: 0.10
Nodes (54): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+46 more)

### Community 22 - "Resource"
Cohesion: 0.05
Nodes (14): fakeDaemon, Resource, GuardRepository, Vault, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots() (+6 more)

### Community 23 - "setpassword_authz_test.go"
Cohesion: 0.12
Nodes (38): CheckPassword(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf(), parseHash() (+30 more)

### Community 24 - "GenerateConf"
Cohesion: 0.13
Nodes (22): discordCatalogEntry(), groupedDiscordConf(), TestGroupedSectionAddressedByEncryptionRoot(), TestPatchCatalogSectionGroupedConfig(), TestSectionsFromCandidatesSkipsMissingWatchSubPaths(), SecretServiceBlock(), EnsureSecretServiceBlock(), firstContentLine() (+14 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.08
Nodes (81): add_inode_to_guard(), admit_port(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), discover_guarded_parent() (+73 more)

### Community 26 - "SanitizeText"
Cohesion: 0.06
Nodes (32): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), systemPatterns(), systemRule(), newHashCache(), ds() (+24 more)

### Community 27 - "install/install.go"
Cohesion: 0.05
Nodes (78): installHooks(), releaseBuild(), Run(), staged(), installGUI(), asRoot(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled() (+70 more)

### Community 28 - "plan_test.go"
Cohesion: 0.40
Nodes (11): scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent(), TestScanAndValidateAcceptsLibDirContent(), TestScanRefusesSpecialFiles(), TestScanUntouchedLibDir(), TestTouches() (+3 more)

### Community 29 - "desktop_test.go"
Cohesion: 0.08
Nodes (21): buttonCursor, pageDef, stageItem, stageState, handCursors(), handCursorsOn(), newButtonCursor(), showConfirm() (+13 more)

### Community 30 - "recordFlow"
Cohesion: 0.04
Nodes (69): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), lines(), offerBackupDeletion(), pickBinaries() (+61 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 32 - "runNetworkMonitor"
Cohesion: 0.12
Nodes (11): DrainUntilSignal(), runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NetworkMonitorRepository, NewNetModel() (+3 more)

### Community 33 - "Guard"
Cohesion: 0.06
Nodes (19): deferredBinary, exeRow, rootHandle, supersededBinary, openRoot(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), SharedPinDegraded(), TestNoteRootMountFollowsMountID() (+11 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (29): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+21 more)

### Community 35 - "secret_probe/main.go"
Cohesion: 0.31
Nodes (18): busPath(), connexec(), connFD(), dial(), execInto(), failf(), get(), handshake() (+10 more)

### Community 37 - "Message"
Cohesion: 0.10
Nodes (10): checkSecret(), checkTrustLine(), trustRows(), backend, backendOp, backendState, Message, Request (+2 more)

### Community 38 - "controlServer"
Cohesion: 0.10
Nodes (23): readConfigPut(), TestReadGrantRequestConfig(), authSelfPeer(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest() (+15 more)

### Community 39 - "render.go"
Cohesion: 0.14
Nodes (13): lineStyles, viewportState, fg(), baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+5 more)

### Community 40 - "pinstate.go"
Cohesion: 0.18
Nodes (22): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+14 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (5): IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "fyne.io/fyne/v2/driver/desktop.Cursor"
Cohesion: 0.24
Nodes (3): handRadioGroup, handSelect, handCursor()

### Community 45 - "planUpdaters"
Cohesion: 0.14
Nodes (18): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+10 more)

### Community 46 - "fakeGuardRepo"
Cohesion: 0.06
Nodes (17): TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant() (+9 more)

### Community 48 - "update_test.go"
Cohesion: 0.10
Nodes (17): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFilterChannel(), testKeyPair(), TestNewerThanInstalled(), TestNewerThanInstalledAcrossTagFormats() (+9 more)

### Community 49 - "fadeScroll"
Cohesion: 0.07
Nodes (19): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+11 more)

### Community 51 - "install_ui.go"
Cohesion: 0.09
Nodes (16): binaryList, installUI, multiAnswer, textAnswer, scrollable(), selectable(), newRadioGroup(), encodeTable() (+8 more)

### Community 52 - "proto.go"
Cohesion: 0.06
Nodes (24): entryOf(), Client, File, Progress, TrustItem, User, childEnv(), launch() (+16 more)

### Community 53 - "watchGroup"
Cohesion: 0.06
Nodes (63): SecretScope, watchGroup, addResource(), applyAllowLib(), applyAnyGroupDirective(), applyBlockLine(), applyConfigLine(), applyDirective() (+55 more)

### Community 54 - "sync.Mutex"
Cohesion: 0.08
Nodes (15): placedKey(), Caller, Gate, Role, Listen(), peerCred(), TestReadFrame(), decodeFrame() (+7 more)

### Community 55 - ".apply"
Cohesion: 0.11
Nodes (15): Apply(), DeleteBackup(), Outcome, ProgressFunc, TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal() (+7 more)

### Community 56 - "changelog.go"
Cohesion: 0.25
Nodes (17): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+9 more)

### Community 57 - "mutedText"
Cohesion: 0.08
Nodes (40): filesPage, panel, newBinaryList(), newBinaryRow(), newCheck(), newSelect(), TestChoiceWidgetsPointingHand(), TestHandControlsShowArrowWhenDisabled() (+32 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (33): engine, portOwner, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+25 more)

### Community 59 - "Monitor"
Cohesion: 0.17
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 60 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 61 - ".FilterExistingWhitelist"
Cohesion: 0.14
Nodes (13): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+5 more)

### Community 62 - "inode_dev"
Cohesion: 0.10
Nodes (31): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+23 more)

### Community 64 - "FileTarget"
Cohesion: 0.06
Nodes (30): handCheckGroup, netGuardTab, policyEditor, policyHooks, policyTab, newCheckGroup(), newPolicyTab(), verbFor() (+22 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "wm.go"
Cohesion: 0.11
Nodes (25): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+17 more)

### Community 68 - "liveSession"
Cohesion: 0.26
Nodes (6): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), holdForward(), pingInterval(), liveSession

### Community 69 - "runForward"
Cohesion: 0.11
Nodes (23): TestValidatePassword(), ValidatePassword(), dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig() (+15 more)

### Community 70 - "time.Time"
Cohesion: 0.11
Nodes (13): LockedError, isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, Process inspectors (+5 more)

### Community 71 - "RefreshCatalog"
Cohesion: 0.17
Nodes (14): sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), RefreshCatalog(), staleDiscordLibConf() (+6 more)

### Community 72 - "NetGuardEvent"
Cohesion: 0.09
Nodes (11): snapNetGuard, NetGuardEvent, NetworkGuardRepository, NetEngines, fakeGuard, fakeMonitor, fakeNet, fakeNetGuard (+3 more)

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "runMonitor"
Cohesion: 0.13
Nodes (18): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), pathWithinDir() (+10 more)

### Community 75 - "LaunchRule"
Cohesion: 0.12
Nodes (20): allowedEnv(), binaryLaunchRules(), launchRulesByPath(), TestAllowedEnvKeepsOnlyRootPlaced(), TestLaunchRulesByPath(), argvEntry, applyLaunchRule(), envName() (+12 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "runNetworkGuard"
Cohesion: 0.11
Nodes (19): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), DiscoverInfraBinaries(), Mode (+11 more)

### Community 78 - "CandidateDir"
Cohesion: 0.08
Nodes (24): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), DiscoverSystem() (+16 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.21
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "os.File"
Cohesion: 0.06
Nodes (45): scriptIdentity(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir() (+37 more)

### Community 81 - "Plan"
Cohesion: 0.18
Nodes (10): clearExcept(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir(), removeAll() (+2 more)

### Community 82 - "LibraryClosure"
Cohesion: 0.06
Nodes (32): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), IsScript(), ldConfInclude(), ldSoConfDirs() (+24 more)

### Community 83 - "shellrc.go"
Cohesion: 0.10
Nodes (56): listFacultatives(), runFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative(), TestListFacultativesIgnoresDaemonConf(), TestListFacultativesMarksGuardedSSH(), removeFacultative(), removeFacultatives() (+48 more)

### Community 84 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.12
Nodes (22): placedWalk, treeJudge, confineBelow(), inUserHome(), hostVouchedDevs(), isFuseType(), VouchedDevs(), vouchedFd() (+14 more)

### Community 85 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 86 - "EventType"
Cohesion: 0.09
Nodes (19): selfEventsFor(), BusClient, EventType, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), AttestBusClient(), busClientKey() (+11 more)

### Community 87 - "loginServer"
Cohesion: 0.07
Nodes (24): markReady(), markReadyWithLogin(), notifySystemdReady(), hasSentinel(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel(), runtimePathsFor() (+16 more)

### Community 88 - "RefreshSection"
Cohesion: 0.07
Nodes (40): logRefreshChange(), TestLogRefreshChange(), LibraryBlock, RefreshOptions, SectionScan, AppendSectionWhitelist(), ConfSafePath(), findSectionEnd() (+32 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "runView"
Cohesion: 0.06
Nodes (14): logGrid, logHook, logLine, runView, journalLine(), TestJournalLine(), logRow(), logStyle() (+6 more)

### Community 92 - "highlight_test.go"
Cohesion: 0.18
Nodes (15): lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens(), TestGutterKeepsTextColumnAligned(), TestHighlighterColorsTokens() (+7 more)

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

### Community 95 - "NetGuard"
Cohesion: 0.08
Nodes (14): foreignPinnedLink(), PumpRingbuf(), reasonText(), eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs() (+6 more)

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
Nodes (14): procStats, statsSample, formatAgo(), cpuTime(), humanBytes(), rssBytes(), New(), recv() (+6 more)

### Community 100 - "BinaryEntry"
Cohesion: 0.09
Nodes (20): binaryVerifyState, HashStore, keptEntries, sharedHashEntry, BinaryEntry, checkBtrfsKeys(), canonicalBinaryPath(), claimHash() (+12 more)

### Community 101 - "auditGroupCoverage"
Cohesion: 0.15
Nodes (16): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource() (+8 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "SampleContent"
Cohesion: 0.18
Nodes (15): TestDaemonUnitWithMetadataOutput(), installFileAs(), installUnit(), upsertFile(), isInstallerSSHAgentUnit(), TestIsInstallerSSHAgentUnit(), TestDaemonUnitCmdPrintsBundledUnit(), SampleContent() (+7 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "RunCmd"
Cohesion: 0.03
Nodes (80): anyDaemonRunning(), offerDaemonStart(), runDiffCatalog(), installFacultatives(), installSSHAgentFacultative(), applyLiveRefresh(), buildBinaryIfNeeded(), certifyInstall() (+72 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.17
Nodes (9): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), btrfsTypes(), TestBtrfsLayoutFromNestedMembers() (+1 more)

### Community 110 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.12
Nodes (30): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), moveTo() (+22 more)

### Community 111 - "ConfinedEntry"
Cohesion: 0.20
Nodes (23): forwardCheck(), ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash() (+15 more)

### Community 112 - "installPage"
Cohesion: 0.07
Nodes (21): actionGate, binSource, installAction, installedBinary, InstallMode, installPage, InstallRequest, latestRelease (+13 more)

### Community 113 - ".netIntegrityContainer"
Cohesion: 0.28
Nodes (5): IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 115 - "protected_test.go"
Cohesion: 0.20
Nodes (15): fakeLink, isOwnExecutable(), newFakeLink(), protectedShell(), TestForwardRefusesOwnExecutable(), TestProtectedEditFlow(), TestProtectedLiveCardsGated(), TestRecordReviewExtendsAnExistingSection() (+7 more)

### Community 116 - "fadescroll_test.go"
Cohesion: 0.32
Nodes (10): newTestFadeScroll(), settle(), TestFadeBarDragEndsGlide(), TestFadeBarDragMovesContent(), TestFadeBarTapPages(), TestFadeScrollBarsHiddenWhenContentFits(), TestFadeScrollBothDirections(), TestFadeScrollClampsOffset() (+2 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "UI"
Cohesion: 0.07
Nodes (26): nextConfig(), showChangelog(), channelOf(), confirmUpdate(), isTerminal(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI() (+18 more)

### Community 121 - "maintainerKeeper"
Cohesion: 0.18
Nodes (12): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), startMaintainers(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), maintainedApp (+4 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "inet_connect.c"
Cohesion: 0.26
Nodes (15): attack_connect(), attack_fastopen(), attack_uring(), denied(), dial(), fill_addr(), handoff(), listen_on() (+7 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "classify"
Cohesion: 0.23
Nodes (11): call(), TestClassify(), TestClassifyOpenCallsExact(), arg(), attributesOf(), classify(), classifyCollection(), classifyProps() (+3 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "catalogRefresher"
Cohesion: 0.11
Nodes (16): drain(), newCatalogRefresher(), readEvents(), refreshAdmits(), swapFile(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), addSystemBinary(), catalogPatterns() (+8 more)

### Community 131 - "runBackend"
Cohesion: 0.36
Nodes (6): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "selection.go"
Cohesion: 0.17
Nodes (5): pos, Highlighter, cellToRune(), Selection, offset()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "vaultfs.go"
Cohesion: 0.05
Nodes (52): invokingIDs(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile() (+44 more)

### Community 143 - "serve.go"
Cohesion: 0.09
Nodes (30): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+22 more)

### Community 144 - "NetEventType"
Cohesion: 0.07
Nodes (20): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), BinariesSummary(), TestBinariesSummary(), FormatAddr() (+12 more)

### Community 145 - "find.go"
Cohesion: 0.13
Nodes (27): findButton, span, findAll(), lowerRunes(), NewFinder(), overlay(), runesEqual(), findEditor() (+19 more)

### Community 146 - "Ledger"
Cohesion: 0.07
Nodes (26): candidates(), confirmAll(), describe(), record(), run(), applySchema(), Ledger, Pending (+18 more)

### Community 147 - "newTestTarget"
Cohesion: 0.22
Nodes (17): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 148 - "GuardEvent"
Cohesion: 0.06
Nodes (18): holdHeadless(), HeadlessLine(), bpfGuardEvent, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestDispatchWithoutOwnAllowedEvents(), TestLocalizePaths() (+10 more)

### Community 149 - "buildTrustedSet"
Cohesion: 0.06
Nodes (32): setConfinedHomes(), inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys() (+24 more)

### Community 150 - "Daemon configuration"
Cohesion: 0.17
Nodes (12): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Launch rules: `[launch_rules]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Port sections, Reloading (+4 more)

### Community 152 - "debouncer"
Cohesion: 0.33
Nodes (4): TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), debouncer

### Community 154 - "buildOneGuard"
Cohesion: 0.07
Nodes (39): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), startSelfGuard(), startSelfGuards(), displayEntries(), resolveGuardConfig() (+31 more)

### Community 155 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 156 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 161 - "serverAuth"
Cohesion: 0.32
Nodes (5): TestClientAuth(), TestSASL(), clientAuth(), serverAuth(), saslServer

### Community 162 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 163 - "RootPlacedName"
Cohesion: 0.29
Nodes (6): PathUnder(), RootPlacedName(), rootOnlyTree(), TestRootPlacedName(), userAssignable(), userMountable()

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (45): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+37 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (41): daemonPage, homePage, pageID, reel, resizeWatch, runHost, shell, sizedDialog (+33 more)

### Community 167 - ".Layout"
Cohesion: 0.44
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 168 - "boundedGlob"
Cohesion: 0.40
Nodes (6): boundedGlob(), exists(), matchDir(), mkdirAll(), TestBoundedGlobMatchesFilepathGlob(), TestBoundedGlobStopsAtTheBudget()

### Community 169 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 170 - "fscrypt/orphans.go"
Cohesion: 0.08
Nodes (38): TestPickEncryptedDirSingleSkipsPicker(), pickDir(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep() (+30 more)

### Community 172 - "Limitations"
Cohesion: 0.18
Nodes (11): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Secrets served over IPC (+3 more)

### Community 173 - "RunWith"
Cohesion: 0.09
Nodes (31): needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), deletePostBackups(), listBackups(), restoreBackups() (+23 more)

### Community 175 - "TrustGuard"
Cohesion: 0.05
Nodes (24): lsmHook, tempGrant, tempMaskOp, tempRow, trustHook, trustRows, ensureBtrfsLayout(), deleteResKeys() (+16 more)

### Community 176 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 177 - "ValidateTargets"
Cohesion: 0.19
Nodes (9): existingSection(), guardedMatch(), overlaps(), refusedPath(), sectionFor(), validateTarget(), ValidateTargets(), Existing (+1 more)

### Community 178 - "deploy"
Cohesion: 0.12
Nodes (26): archAssets(), assetsFor(), checkUpdatePreconditions(), deploy(), deployGUI(), fetchGUI(), guiAssetsFor(), hasAsset() (+18 more)

### Community 179 - "IntegrationSuite"
Cohesion: 0.24
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 180 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

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

### Community 191 - "github.com/spf13/cobra.Command"
Cohesion: 0.08
Nodes (38): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+30 more)

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 194 - "probeSystem"
Cohesion: 0.40
Nodes (5): probeSystem(), Check(), collectDiagnostics(), readSysctl(), requireKernel51()

### Community 195 - "maintainer/plan.go"
Cohesion: 0.29
Nodes (9): isWhiteout(), opaque(), pathByte(), pathRecord(), TestPathRecord(), libPlan, op, opKind (+1 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 199 - "migrate_test.go"
Cohesion: 0.07
Nodes (38): harnessSuite, inode(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceRecoversFromInterruptedTransform(), TestUnlockFileInPlaceWrongKeyFails(), TestUnlockLockFileInPlacePreservesInode() (+30 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "Vault"
Cohesion: 0.05
Nodes (66): CheckBPFLSM(), CheckEBPF(), backingDeviceUnion(), buildGuards(), configureDaemonLogging(), confirmMasterKeyOverwrite(), loadDaemonConfig(), lockOneRoot() (+58 more)

### Community 202 - "DaemonEvent"
Cohesion: 0.10
Nodes (36): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless() (+28 more)

### Community 203 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 204 - "FileEvent"
Cohesion: 0.12
Nodes (6): snapMonitor, FileEvent, MonitorRepository, newFakeMonitorRepo(), fakeSource, fakeMonitorRepo

### Community 205 - "eventUnitTest"
Cohesion: 0.27
Nodes (4): BpfEvent, eventUnitTest, Cstr(), DecodeBpfEvent()

### Community 206 - "monitor"
Cohesion: 0.09
Nodes (24): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks (+16 more)

### Community 209 - "main_test.go"
Cohesion: 0.20
Nodes (5): infraContainerPaths(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

### Community 210 - "Run"
Cohesion: 0.33
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "KernelDev"
Cohesion: 0.11
Nodes (27): fdStat, BinaryStat, physicalParent(), statFD(), binaryStatOf(), KernelDev(), StatBinary(), StatBinaryFile() (+19 more)

### Community 213 - "parseKernelRelease"
Cohesion: 0.67
Nodes (4): ciliumVersion, parseKernelRelease(), TestParseKernelRelease(), unameKernelVersion()

### Community 214 - "fileedit_symlink_test.go"
Cohesion: 0.36
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 216 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 217 - "controlManager"
Cohesion: 0.29
Nodes (3): newControlManager(), startControlManager(), controlManager

### Community 218 - "parseDaemonEvents"
Cohesion: 0.15
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite

### Community 221 - "downloadAndVerify"
Cohesion: 0.16
Nodes (16): downloadAndVerify(), downloadFile(), downloadReleaseFiles(), FetchLatest(), fetchPair(), fetchRelease(), fetchReleases(), FetchVerified() (+8 more)

### Community 222 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 224 - "linkView"
Cohesion: 0.31
Nodes (5): linkTarget(), newLinkView(), TestCheckLinksJudgesEachTarget(), entryKind, linkView

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "guardModel"
Cohesion: 0.05
Nodes (39): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+31 more)

### Community 228 - "mountTable"
Cohesion: 0.13
Nodes (10): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), ReadMountinfo(), unescapeMountinfo(), TestParseMountDevMatchesStatEncoding() (+2 more)

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 238 - "Candidate"
Cohesion: 0.09
Nodes (35): appendSectionsAndEdit(), collectDiffAdditions(), TestDiffMergePreservesExistingAndParses(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), addManualDirectories(), editConfig() (+27 more)

### Community 240 - "shQuote"
Cohesion: 0.17
Nodes (5): IntegrationSuite, IntegrationSuite, overlayPlanter(), shQuote(), IntegrationSuite

### Community 243 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 249 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 256 - "sources_test.go"
Cohesion: 0.14
Nodes (10): TestGeneratedLibraryBlockRoundTrips(), systemTool(), archiveMembers(), decompressed(), newSources(), tarMembers(), TestArchiveMembersIgnoresNonArchives(), TestSourcesHoldFilesAndArchiveMembers() (+2 more)

### Community 264 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 665 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **39 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.netIntegrityContainer`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`, `parseDaemonEvents`, `daemon_multicall_test.go`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.12295030371699242 - nodes in this community are weakly interconnected._
- **Why does `Guard` connect `Guard` to `go_pkg_strings`, `time.Duration`, `engine`, `BinaryEntry`, `time.Time`, `SanitizeText`, `guardUnitTest`, `Vault`, `DaemonEvent`, `GuardInodeKey`, `TrustGuard`, `os.File`, `GuardEvent`, `Resource`, `loginServer`, `sync.Mutex`, `buildOneGuard`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.021546403321737177 - nodes in this community are weakly interconnected._
- **Why does `monitor()` connect `monitor` to `time.Duration`, `secret_probe/main.go`, `app-listener documentation`, `How it works`, `Installation`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09263157894736843 - nodes in this community are weakly interconnected._