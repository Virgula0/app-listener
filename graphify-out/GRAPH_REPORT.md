# Graph Report - app-listener  (2026-11-06)

## Corpus Check
- 567 files · ~1,808,451 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7462 nodes · 34112 edges · 254 communities (212 shown, 42 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 3050 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `249afbc1`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- go_pkg_strings
- testing.T
- IntegrationSuite
- Plan
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- GuardEvent
- User
- FileTarget
- BinaryEntry
- buildOneGuard
- DaemonEvent
- IntegrationSuite
- guard_trust.bpf.c
- navItem
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- plan_test.go
- guard.bpf.c
- go_pkg_github_com_charmbracelet_bubbletea
- Network
- gui/session.go
- fileedit_model_test.go
- go_pkg_path_filepath
- RunCmd
- guardModel
- protected_test.go
- picker
- secret_probe/main.go
- IntegrationSuite
- Message
- editView
- runNetworkMonitor
- New
- watchSet
- deploy
- monitorUnitTest
- pinstate.go
- desktop_test.go
- githubRelease
- stdio
- os.File
- fadeScroll
- LibraryClosure
- protectedPage
- proto.go
- Finder
- session
- .apply
- backend_test.go
- UI
- GuardInodeKey
- NetEventType
- sync.Mutex
- serve.go
- newTestTarget
- golang.org/x/sys/unix.Stat_t
- daemon_multicall_test.go
- unix_connect.c
- IntegrationSuite
- IntegrationSuite
- newTestVetter
- parseDaemonEvents
- Vault
- RootPlacedName
- planUpdaters
- fcntl
- policyEditor
- Guard
- codecheck/main.go
- vaultfs.go
- liveSession
- process_vm_readv.c
- Monitor
- controlServer
- daemon.conf reference
- shellrc.go
- inodeChain
- bpfstats/main.go
- inode_dev
- loginServer
- updatePage
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- LaunchRule
- batchInodeRows
- netinject.c
- README.md
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- run.go
- sources
- Well-known bypass classes
- Config
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- backend.go
- Build & sign release assets (reusable workflow)
- TestCatalogSanity
- github.com/charmbracelet/bubbles/textarea.Model
- FileEvent
- Binary
- recordFlow
- IntegrationSuite
- wire.go
- Recorder
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- time.Time
- IntegrationSuite
- WebSocket /ws client
- CandidateDir
- inet_connect.c
- install.sh
- NetworkMonitor
- BinaryRule
- make test-integration (rootful Docker suite)
- SelectOrphans
- buildTrustedSet
- github.com/spf13/cobra.Command
- trace-app-libs.sh
- IntegrationSuite
- classifyGroup
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- NetGuardEvent
- runOfflineEdit
- Ledger
- TestMonitorUseCaseLifecycle
- maintainerKeeper
- btrfsLayoutFrom
- ConfinedEntry
- audit.go
- fakeStore
- Vault
- IntegrationSuite
- .populateWatchInodes
- hashFileShared
- copytree.go
- KernelDev
- Installation
- serverAuth
- graphify-refresh.sh
- layout.go
- update_test.go
- runLockdown
- absPath
- feed.go
- shell
- findUpdate
- github.com/charmbracelet/bubbletea.Cmd
- changelog_test.go
- classify
- daemon: permanent protection with encryption at rest
- Limitations
- cursor.go
- root.go
- net_tester/main.go
- tempGrantSetup
- .judgeProcess
- fileedit_symlink_test.go
- IntegrationSuite
- journal
- openBunTmp
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- debouncer
- Self-updating apps and the catalog
- .netIntegrityContainer
- io.Writer
- downloadFile
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- secretJudge
- NewMonitor
- lockAndDeprovision
- .admitGated
- Gate
- IntegrationSuite
- runView
- monitor
- IntegrationSuite
- IntegrationSuite
- RefreshSection
- maintainer/plan.go
- mc_tag.h
- latestTime
- securewrite_test.go
- .awaitDaemonUpOrDead
- .TestRecord_GuardParity
- resolveFilesystemPrereqs
- sources_test.go
- find_test.go
- fakeDaemon
- backups.go
- ValidateTargets
- bpfLsmListed
- runNetworkGuard
- supersede_probe.c
- NewEventFanout
- RawTarget
- Daemon configuration
- main_test.go
- github.com/Virgula0/app-listener/tools/codecheck
- safeio.go
- taint_victim.c
- network monitor: see what a program talks to
- mutedText
- options.go
- Troubleshooting
- mountTable
- Candidate
- IntegrationSuite
- IntegrationSuite
- daemon/ports.go
- parseGuardEvent
- namedGeneralTool
- VouchedDevs
- github.com/stretchr/testify/suite.Suite
- CleanupStalePins
- launchKey
- prctl_gate
- .eventTypeAllowed
- bleedLayout
- daemonconfig.go
- .keyringObject

## God Nodes (most connected - your core abstractions)
1. `Config` - 147 edges
2. `Guard` - 125 edges
3. `GuardInodeKey` - 113 edges
4. `Load()` - 97 edges
5. `shell` - 97 edges
6. `absPath()` - 94 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 75 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Distributions` --references--> `monitor()`  [INFERRED]
  docs/compatibility.md → integrationtests/secret_probe/main.go
- `Files` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `Workspaces` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `eBPF at three levels` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go
- `Why not an AppArmor profile?` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go

## Import Cycles
- None detected.

## Communities (254 total, 42 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.14
Nodes (29): tempJournalRow, argvMatch, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion (+21 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (178): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+170 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (7): configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig()

### Community 3 - "Plan"
Cohesion: 0.15
Nodes (12): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+4 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.06
Nodes (84): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+76 more)

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
Cohesion: 0.06
Nodes (15): addErrKind, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool(), walkEntries() (+7 more)

### Community 9 - "mkdirs"
Cohesion: 0.12
Nodes (38): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetHomeDepartureTriggersRefresh(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds() (+30 more)

### Community 10 - "GuardEvent"
Cohesion: 0.04
Nodes (24): holdHeadless(), HeadlessLine(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths(), commMatchesGuardedBinary(), GuardEvent (+16 more)

### Community 11 - "User"
Cohesion: 0.04
Nodes (97): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), gatherPerUserSetup() (+89 more)

### Community 12 - "FileTarget"
Cohesion: 0.08
Nodes (15): filesPage, guardPill(), newExeResolver(), Files, FileSpec, FileTarget, Policy, guardRow() (+7 more)

### Community 13 - "BinaryEntry"
Cohesion: 0.04
Nodes (41): selfEventsFor(), deferredBinary, exeRow, keptEntries, supersededBinary, tempMaskOp, tempRow, BinaryEntry (+33 more)

### Community 14 - "buildOneGuard"
Cohesion: 0.04
Nodes (66): DrainUntilSignal(), CheckEBPF(), MakeDisplayPaths(), backingDeviceUnion(), buildGuards(), buildOneGuard(), drainEphemeral(), eventFilterOptions() (+58 more)

### Community 15 - "DaemonEvent"
Cohesion: 0.11
Nodes (34): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+26 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (64): argv_env(), argv_matched(), argv_not_root(), argv_rule_at(), argv_short(), argv_word(), caller_updates(), current_exe_flags() (+56 more)

### Community 18 - "navItem"
Cohesion: 0.03
Nodes (20): appTheme, binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panelRenderer (+12 more)

### Community 19 - "fileEditModel"
Cohesion: 0.09
Nodes (11): clipLabel(), humanSize(), listEntries(), paneTitle(), sortNodes(), chownUser, dirEntry, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (51): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+43 more)

### Community 22 - "Resource"
Cohesion: 0.06
Nodes (14): ownerPaths(), findResource(), Resource, GuardRepository, Vault, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase (+6 more)

### Community 23 - "auth.go"
Cohesion: 0.12
Nodes (40): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+32 more)

### Community 24 - "plan_test.go"
Cohesion: 0.40
Nodes (11): scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent(), TestScanAndValidateAcceptsLibDirContent(), TestScanRefusesSpecialFiles(), TestScanUntouchedLibDir(), TestTouches() (+3 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.07
Nodes (96): exe_refused(), add_inode_to_guard(), admit_port(), btrfs_copy_gate(), bus_key(), check_and_emit_args(), chmod_only_drops_write(), code_suspect() (+88 more)

### Community 26 - "go_pkg_github_com_charmbracelet_bubbletea"
Cohesion: 0.08
Nodes (49): newChangelogModel(), showChangelog(), class, classOf(), column(), dedent(), Edit(), leadingIndent() (+41 more)

### Community 27 - "Network"
Cohesion: 0.13
Nodes (8): NetworkGuardRepository, clonePolicy(), NetEngines, NetExtra, NetPolicy, Network, verifyHashes(), NetGuardConfig

### Community 28 - "gui/session.go"
Cohesion: 0.07
Nodes (39): ancestorEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue(), TestOwnedByChecksTypeAndOwner() (+31 more)

### Community 29 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 30 - "go_pkg_path_filepath"
Cohesion: 0.36
Nodes (4): recordTick, LiveEmptyWhitelistRejected(), TestLiveEmptyWhitelistRejected(), hint

### Community 31 - "RunCmd"
Cohesion: 0.04
Nodes (67): applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig(), anyDaemonRunning() (+59 more)

### Community 32 - "guardModel"
Cohesion: 0.05
Nodes (39): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+31 more)

### Community 33 - "protected_test.go"
Cohesion: 0.20
Nodes (15): fakeLink, isOwnExecutable(), newFakeLink(), protectedShell(), TestForwardRefusesOwnExecutable(), TestProtectedEditFlow(), TestProtectedLiveCardsGated(), TestRecordReviewExtendsAnExistingSection() (+7 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (31): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+23 more)

### Community 35 - "secret_probe/main.go"
Cohesion: 0.18
Nodes (31): busPath(), connexec(), connFD(), createItem(), dial(), dupkey(), execInto(), failf() (+23 more)

### Community 37 - "Message"
Cohesion: 0.10
Nodes (10): checkSecret(), checkTrustLine(), trustRows(), backend, backendOp, backendState, OpenLedger(), Message (+2 more)

### Community 38 - "editView"
Cohesion: 0.10
Nodes (8): codeEntry, editView, Entry, fileSize(), joinRel(), newCodeEntry(), relOf(), endedBanner()

### Community 39 - "runNetworkMonitor"
Cohesion: 0.10
Nodes (13): ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString() (+5 more)

### Community 40 - "New"
Cohesion: 0.08
Nodes (47): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault() (+39 more)

### Community 41 - "watchSet"
Cohesion: 0.14
Nodes (15): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchDir, watchPattern (+7 more)

### Community 42 - "deploy"
Cohesion: 0.18
Nodes (15): applyUpdate(), confirmUnit(), confirmUpdate(), deploy(), deployGUI(), InstalledGUIVersion(), InstalledVersion(), isTerminal() (+7 more)

### Community 44 - "pinstate.go"
Cohesion: 0.18
Nodes (22): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+14 more)

### Community 45 - "desktop_test.go"
Cohesion: 0.11
Nodes (16): TestEveryInstallModeMapped(), pageDef, stageItem, stageState, textAnswer, TestDaemonPageGates(), encodeTable(), invokingIDs() (+8 more)

### Community 46 - "githubRelease"
Cohesion: 0.20
Nodes (18): archAssets(), assetsFor(), checkUpdatePreconditions(), downloadAndVerify(), fetchGUI(), fetchPair(), guiAssetsFor(), hasAsset() (+10 more)

### Community 48 - "os.File"
Cohesion: 0.05
Nodes (64): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), scriptIdentity(), placedKey(), startTrustGuard(), binaryVetter, vetDecision (+56 more)

### Community 49 - "fadeScroll"
Cohesion: 0.08
Nodes (15): fadeBar, fadeBarRenderer, fadeScroll, fadeScrollRenderer, fadeScrollOf(), newTestFadeScroll(), settle(), TestFadeBarDragEndsGlide() (+7 more)

### Community 50 - "LibraryClosure"
Cohesion: 0.08
Nodes (27): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+19 more)

### Community 51 - "protectedPage"
Cohesion: 0.04
Nodes (30): backendLink, configView, forwardForm, forwardView, installUI, multiAnswer, protectedPage, recordReview (+22 more)

### Community 52 - "proto.go"
Cohesion: 0.06
Nodes (25): Client, File, Progress, TrustItem, User, childEnv(), launch(), helper() (+17 more)

### Community 53 - "Finder"
Cohesion: 0.16
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 54 - "session"
Cohesion: 0.15
Nodes (6): peerCred(), TestReadFrame(), readFrame(), queued, senderQueue, session

### Community 55 - ".apply"
Cohesion: 0.16
Nodes (7): Apply(), ProgressFunc, applier, ledgerOps, PlannedBinary, Store, vaultOps

### Community 56 - "backend_test.go"
Cohesion: 0.19
Nodes (23): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+15 more)

### Community 57 - "UI"
Cohesion: 0.10
Nodes (15): channelOf(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestRunWithRefusesConcurrentRun(), cancelUI, TestEditTextMapsUICancel(), Confirm() (+7 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (47): engine, portOwner, resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded(), nearestRootClaims() (+39 more)

### Community 59 - "NetEventType"
Cohesion: 0.10
Nodes (14): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), newNetEventChecks(), FormatAddr(), NetEventTypes() (+6 more)

### Community 60 - "sync.Mutex"
Cohesion: 0.04
Nodes (36): setConfinedHomes(), argvMatches, lsmHook, refusedReplacements, statFunc, tempGrant, trustHook, trustRows (+28 more)

### Community 61 - "serve.go"
Cohesion: 0.09
Nodes (31): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+23 more)

### Community 62 - "newTestTarget"
Cohesion: 0.24
Nodes (16): NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded(), TestApplyFromWhitelistStopsFirst() (+8 more)

### Community 63 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.16
Nodes (17): placedWalk, treeJudge, confineBelow(), inUserHome(), vouchedFd(), GlobSystemPlaced(), placedDir(), placedEntry() (+9 more)

### Community 64 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 68 - "newTestVetter"
Cohesion: 0.44
Nodes (10): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+2 more)

### Community 69 - "parseDaemonEvents"
Cohesion: 0.15
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite

### Community 70 - "Vault"
Cohesion: 0.16
Nodes (15): secureResources(), verifyEncryptionState(), askEncryption(), isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+7 more)

### Community 71 - "RootPlacedName"
Cohesion: 0.24
Nodes (7): RootPlacedName(), rootOnlyTree(), TestRootPlacedName(), userAssignable(), linkTarget(), entryKind, linkView

### Community 72 - "planUpdaters"
Cohesion: 0.17
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "policyEditor"
Cohesion: 0.14
Nodes (11): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), confirmWarnings() (+3 more)

### Community 75 - "Guard"
Cohesion: 0.09
Nodes (12): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), rootHandle, openRoot(), TestWithChmodDropWriteNeedsReadOnly(), SharedPinDegraded(), TestNoteRootMountFollowsMountID() (+4 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "vaultfs.go"
Cohesion: 0.12
Nodes (24): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), IsBinaryContent(), IsOctalMode(), ParseMode(), inodeOf(), mustRead() (+16 more)

### Community 78 - "liveSession"
Cohesion: 0.11
Nodes (23): dialLiveSession(), parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries() (+15 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.21
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "Monitor"
Cohesion: 0.13
Nodes (3): exeKey, Monitor, watchedInode

### Community 81 - "controlServer"
Cohesion: 0.08
Nodes (26): readConfigPut(), TestReadGrantRequestConfig(), authSelfPeer(), newControlManager(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest() (+18 more)

### Community 82 - "daemon.conf reference"
Cohesion: 0.09
Nodes (22): A complete example, `allow_env`, `allow_lib`, Binary lines and events, daemon.conf reference, `deny_args`, `warn_args`, `deny_env`, `warn_env`, `[electron_apps]` (+14 more)

### Community 83 - "shellrc.go"
Cohesion: 0.15
Nodes (41): rcTarget, TestHasRCBlocksTrackEnsureAndRemove(), TestRemoveSSHAgentShellEnvKeepsSSHConfig(), applyRC(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), ensureRCBlock() (+33 more)

### Community 84 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 85 - "bpfstats/main.go"
Cohesion: 0.12
Nodes (26): CheckBPFLSM(), confirmMasterKeyOverwrite(), runBPFCheck(), runGenKey(), runOneShotMode(), guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection() (+18 more)

### Community 86 - "inode_dev"
Cohesion: 0.11
Nodes (24): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+16 more)

### Community 87 - "loginServer"
Cohesion: 0.10
Nodes (19): markReady(), markReadyWithLogin(), notifySystemdReady(), hasSentinel(), lockParent(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel() (+11 more)

### Community 88 - "updatePage"
Cohesion: 0.21
Nodes (8): installedBinary, latestRelease, updatePage, hashFile(), hashReader(), readInstalled(), boolIndex(), isPreRelease()

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "LaunchRule"
Cohesion: 0.10
Nodes (27): binaryLaunchRules(), launchRulesByPath(), TestLaunchRulesByPath(), ShortOpts, argvEntry, applyLaunchRule(), envName(), LaunchRule (+19 more)

### Community 92 - "batchInodeRows"
Cohesion: 0.24
Nodes (11): inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows(), walkInodeBatches(), countInodeRows() (+3 more)

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

### Community 94 - "README.md"
Cohesion: 0.32
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 95 - "NetGuard"
Cohesion: 0.09
Nodes (13): FirstOnBtrfs(), OnBtrfs(), reasonText(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses() (+5 more)

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
Nodes (15): formatElapsed(), procStats, statsSample, duration(), cpuTime(), humanBytes(), rssBytes(), New() (+7 more)

### Community 100 - "run.go"
Cohesion: 0.36
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "Config"
Cohesion: 0.08
Nodes (27): drain(), newCatalogRefresher(), readEvents(), swapFile(), makeReloadHandler(), noteReloadedSettings(), startCatalogRefresh(), resolveInspectors() (+19 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "backend.go"
Cohesion: 0.06
Nodes (53): asRoot(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), ensureInstalledBinary(), ensureMasterKey(), fileSHA256(), prepareInstallation(), runInstall() (+45 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "TestCatalogSanity"
Cohesion: 0.50
Nodes (3): isAbsoluteCandidate(), TestCatalogSanity(), TestRuntimeEntriesNeverDeleteRoot()

### Community 110 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.10
Nodes (18): lineStyles, pos, viewportState, fg(), baseStyle(), gutterLabel(), lineLen(), markStyle() (+10 more)

### Community 111 - "FileEvent"
Cohesion: 0.11
Nodes (9): snapMonitor, BpfEvent, eventUnitTest, FileEvent, Cstr(), DecodeBpfEvent(), newFakeMonitorRepo(), fakeSource (+1 more)

### Community 112 - "Binary"
Cohesion: 0.23
Nodes (11): lines(), pickBinaries(), pickLibs(), reviewRecording(), Classify(), Binary, Recording, dedupe() (+3 more)

### Community 113 - "recordFlow"
Cohesion: 0.15
Nodes (15): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+7 more)

### Community 115 - "wire.go"
Cohesion: 0.20
Nodes (13): dictFrame(), TestDuplicateDictKeyRefused(), decodeFrame(), distinctDictKeys(), encodeWithSerial(), header(), sameDict(), sameEntries() (+5 more)

### Community 116 - "Recorder"
Cohesion: 0.17
Nodes (14): Target, drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start() (+6 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "time.Time"
Cohesion: 0.21
Nodes (7): LockedError, newWatchSet(), fakeInotify, sweepID, sweepMemo, Progress, Seen

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "CandidateDir"
Cohesion: 0.04
Nodes (53): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, GlobChild, compileGlobNames() (+45 more)

### Community 124 - "inet_connect.c"
Cohesion: 0.25
Nodes (17): attack_connect(), attack_fastopen(), attack_uring(), denied(), dial(), fill_addr(), handoff(), listen_on() (+9 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "NetworkMonitor"
Cohesion: 0.15
Nodes (8): BinariesSummary(), TestBinariesSummary(), PumpRingbuf(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 127 - "BinaryRule"
Cohesion: 0.40
Nodes (5): systemPatterns(), systemRule(), BinaryRule, materialized(), Extension

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 130 - "buildTrustedSet"
Cohesion: 0.16
Nodes (16): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), allowedEnv(), TestAllowedEnvKeepsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree() (+8 more)

### Community 131 - "github.com/spf13/cobra.Command"
Cohesion: 0.09
Nodes (34): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+26 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "classifyGroup"
Cohesion: 0.11
Nodes (24): loaderName(), LoaderRefusal(), soNames(), appletLine(), beyondEvents(), classifyGroup(), classifyLibs(), finishBinary() (+16 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 143 - "NetGuardEvent"
Cohesion: 0.18
Nodes (6): snapNetGuard, NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 144 - "runOfflineEdit"
Cohesion: 0.12
Nodes (16): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), TestPickEncryptedDirSingleSkipsPicker(), chooseResource() (+8 more)

### Community 145 - "Ledger"
Cohesion: 0.08
Nodes (25): candidates(), confirmAll(), describe(), record(), run(), applySchema(), Ledger, Pending (+17 more)

### Community 146 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 147 - "maintainerKeeper"
Cohesion: 0.18
Nodes (12): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), startMaintainers(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), maintainedApp (+4 more)

### Community 148 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 149 - "ConfinedEntry"
Cohesion: 0.29
Nodes (18): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+10 more)

### Community 150 - "audit.go"
Cohesion: 0.21
Nodes (16): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource() (+8 more)

### Community 155 - "hashFileShared"
Cohesion: 0.22
Nodes (9): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+1 more)

### Community 156 - "copytree.go"
Cohesion: 0.21
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 157 - "KernelDev"
Cohesion: 0.07
Nodes (31): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, BusClient, fakeHashStore, fdStat, BinaryStat (+23 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "serverAuth"
Cohesion: 0.32
Nodes (5): TestClientAuth(), TestSASL(), clientAuth(), serverAuth(), saslServer

### Community 161 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 162 - "update_test.go"
Cohesion: 0.12
Nodes (15): newerThanInstalled(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchGUIWithoutAsset(), testKeyPair(), TestNewerThanInstalled() (+7 more)

### Community 163 - "runLockdown"
Cohesion: 0.14
Nodes (16): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), resolveConfigPath(), runLockdown(), dirIsRootOwnedSafe(), isBpffs() (+8 more)

### Community 164 - "absPath"
Cohesion: 0.06
Nodes (8): IntegrationSuite, IntegrationSuite, shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite, IntegrationSuite, absPath()

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (46): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+38 more)

### Community 166 - "shell"
Cohesion: 0.06
Nodes (42): pageID, reel, resizeWatch, shell, sizedDialog, snapEngines, snapNet, sysStatus (+34 more)

### Community 167 - "findUpdate"
Cohesion: 0.19
Nodes (15): installHooks(), staged(), FetchLatest(), fetchReleases(), FetchVerified(), filterChannel(), findUpdate(), Release (+7 more)

### Community 168 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.07
Nodes (15): recordTickCmd(), recordModel, diffModel, editorModel, listenForDaemonEvents(), tickDaemonStats(), clamp(), renderBullets() (+7 more)

### Community 169 - "changelog_test.go"
Cohesion: 0.19
Nodes (12): changelogText(), cumulativeChangelog(), releasesSince(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize() (+4 more)

### Community 170 - "classify"
Cohesion: 0.18
Nodes (15): call(), TestClassify(), TestClassifyCreateItemReplace(), TestClassifyItemTypeIsSchema(), TestClassifyOpenCallsExact(), arg(), attributesOf(), classify() (+7 more)

### Community 171 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.13
Nodes (15): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+7 more)

### Community 172 - "Limitations"
Cohesion: 0.18
Nodes (11): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Secrets served over IPC (+3 more)

### Community 173 - "cursor.go"
Cohesion: 0.08
Nodes (22): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+14 more)

### Community 174 - "root.go"
Cohesion: 0.07
Nodes (38): releaseBuild(), Run(), installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary() (+30 more)

### Community 176 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 177 - "tempGrantSetup"
Cohesion: 0.25
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

### Community 178 - ".judgeProcess"
Cohesion: 0.22
Nodes (10): codeVouched(), isTraced(), loaderEnv(), readMountDevs(), roFor(), TestCodeVouched(), TestLoaderEnv(), TestReadMountDevs() (+2 more)

### Community 179 - "fileedit_symlink_test.go"
Cohesion: 0.36
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 181 - "journal"
Cohesion: 0.14
Nodes (6): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 182 - "openBunTmp"
Cohesion: 0.23
Nodes (11): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+3 more)

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

### Community 191 - "debouncer"
Cohesion: 0.33
Nodes (4): TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), debouncer

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - ".netIntegrityContainer"
Cohesion: 0.28
Nodes (5): IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 194 - "io.Writer"
Cohesion: 0.14
Nodes (10): orNone(), printBinary(), printRecording(), TestProgressReader(), fyneLogFilter, bottomBar, newBottomBar(), renderBar() (+2 more)

### Community 195 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 199 - "secretJudge"
Cohesion: 0.18
Nodes (8): secretJudge, SecretScope, SecretMatch, SecretService, parseSecretMatch(), TestSecretMatch(), LiveGuardAt(), Scope

### Community 200 - "NewMonitor"
Cohesion: 0.23
Nodes (6): evalSymlinksOrEmpty(), NewMonitor(), newPathCache(), putBtrfsLayout(), pathCache, pathCacheEntry

### Community 201 - "lockAndDeprovision"
Cohesion: 0.12
Nodes (12): deprovisionKind, migrationTarget, writeTempWithMetadata(), classifyDeprovision(), classifyDeprovisionErr(), classifyMigrationTarget(), isDeprovisionBusy(), isDeprovisionMissing() (+4 more)

### Community 202 - ".admitGated"
Cohesion: 0.23
Nodes (3): Caller, Role, Judge

### Community 203 - "Gate"
Cohesion: 0.21
Nodes (5): connectionGone(), Gate, Listen(), TestConnectionGoneOnlyFromTheBus(), Bus

### Community 205 - "runView"
Cohesion: 0.05
Nodes (16): daemonPage, daemonState, logGrid, logHook, logLine, runHost, runView, journalLine() (+8 more)

### Community 206 - "monitor"
Cohesion: 0.06
Nodes (39): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks (+31 more)

### Community 209 - "IntegrationSuite"
Cohesion: 0.23
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 210 - "RefreshSection"
Cohesion: 0.08
Nodes (34): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+26 more)

### Community 211 - "maintainer/plan.go"
Cohesion: 0.31
Nodes (8): isWhiteout(), opaque(), pathByte(), pathRecord(), TestPathRecord(), op, opKind, walker

### Community 212 - "mc_tag.h"
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 213 - "latestTime"
Cohesion: 0.29
Nodes (10): newReleaseWindow(), compareStableVersions(), latestTime(), newerThanStable(), parseStableVersion(), parseVerPart(), TestNewerThanStable(), TestParseStableVersion() (+2 more)

### Community 214 - "securewrite_test.go"
Cohesion: 0.31
Nodes (8): atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync(), writeWithin()

### Community 216 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 217 - "resolveFilesystemPrereqs"
Cohesion: 0.16
Nodes (9): encryptDirectories(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), harnessSuite, EnsureSystemSetup(), isKernelAtLeast54(), Prereq (+1 more)

### Community 218 - "sources_test.go"
Cohesion: 0.19
Nodes (10): SystemTrusted(), systemTool(), archiveMembers(), decompressed(), newSources(), tarMembers(), TestArchiveMembersIgnoresNonArchives(), TestSourcesHoldFilesAndArchiveMembers() (+2 more)

### Community 219 - "find_test.go"
Cohesion: 0.19
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 221 - "backups.go"
Cohesion: 0.12
Nodes (28): deletePostBackups(), listBackups(), restoreBackups(), confirmKeyDeletion(), offerBackupCleanup(), run(), runUninstall(), RunWith() (+20 more)

### Community 222 - "ValidateTargets"
Cohesion: 0.19
Nodes (9): existingSection(), guardedMatch(), overlaps(), refusedPath(), sectionFor(), validateTarget(), ValidateTargets(), Existing (+1 more)

### Community 223 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 224 - "runNetworkGuard"
Cohesion: 0.11
Nodes (21): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), WithSuspectsFrom() (+13 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 227 - "RawTarget"
Cohesion: 0.31
Nodes (9): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+1 more)

### Community 228 - "Daemon configuration"
Cohesion: 0.20
Nodes (10): Daemon configuration, Daemon settings, Generic Electron apps, Launch rules, Library trust, Port sections, Reloading, Runtime sections (+2 more)

### Community 229 - "main_test.go"
Cohesion: 0.20
Nodes (5): infraContainerPaths(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

### Community 231 - "safeio.go"
Cohesion: 0.42
Nodes (8): createOwned(), AtomicWriteAt(), CreateExclNoFollow(), descend(), DescendCreateOwned(), DescendNoFollow(), mkdirAt(), openDirNoFollow()

### Community 233 - "network monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 234 - "mutedText"
Cohesion: 0.05
Nodes (49): actionGate, binaryList, binSource, checkGrid, handButton, handCheck, homePage, installAction (+41 more)

### Community 235 - "options.go"
Cohesion: 0.36
Nodes (9): InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks, InstallHooks, StagedRelease (+1 more)

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "mountTable"
Cohesion: 0.18
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo(), userMountable() (+1 more)

### Community 238 - "Candidate"
Cohesion: 0.05
Nodes (75): appendSectionsAndEdit(), TestDiffMergePreservesExistingAndParses(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickFromCandidates(), runConfigEditor() (+67 more)

### Community 241 - "daemon/ports.go"
Cohesion: 0.24
Nodes (11): loadDaemonConfig(), ensureAnchor(), ensurePortAnchors(), ensureRootDir(), ephemeralRange(), parseReservedPorts(), readReservedPorts(), TestParseReservedPorts() (+3 more)

### Community 242 - "parseGuardEvent"
Cohesion: 0.29
Nodes (8): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels()

### Community 243 - "namedGeneralTool"
Cohesion: 0.50
Nodes (4): namedGeneralTool(), TestCatalogWhitelistsNoGeneralTool(), TestFileTransferToolsAreGeneral(), TestNamedGeneralTool()

### Community 244 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 246 - "CleanupStalePins"
Cohesion: 0.33
Nodes (6): CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 248 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 257 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 267 - "daemonconfig.go"
Cohesion: 0.08
Nodes (61): watchGroup, addResource(), applyAllowLib(), applyAnyGroupDirective(), applyBlockLine(), applyConfigLine(), applyDirective(), applyElectronApp() (+53 more)

### Community 280 - ".keyringObject"
Cohesion: 0.24
Nodes (3): placeableAttributes(), TestLearnKeyringFromFirstReply(), TestPlaceableAttributes()

## Knowledge Gaps
- **34 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck` (+29 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 688 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **42 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `monitor()` connect `monitor` to `time.Duration`, `secret_probe/main.go`, `Installation`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _34 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.13602490158381397 - nodes in this community are weakly interconnected._
- **Why does `shQuote()` connect `absPath` to `go_pkg_strings`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.016661037937633683 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `daemon_multicall_test.go`, `.netIntegrityContainer`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `parseDaemonEvents`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09441441441441441 - nodes in this community are weakly interconnected._