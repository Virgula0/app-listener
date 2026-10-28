# Graph Report - app-listener  (2026-10-28)

## Corpus Check
- 532 files · ~1,751,055 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6907 nodes · 30675 edges · 235 communities (193 shown, 42 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2859 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `0ac50459`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- go_pkg_strings
- testing.T
- IntegrationSuite
- fscrypt.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- fakeGuardRepo
- New
- protectedPage
- os.File
- go_pkg_path_filepath
- Network
- IntegrationSuite
- guard_trust.bpf.c
- fadescroll.go
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- Guard
- guard.bpf.c
- TrustGuard
- find_test.go
- FileTarget
- golang.org/x/sys/unix.Stat_t
- github.com/charmbracelet/bubbletea.Cmd
- MonitorRepository
- binaryList
- EventType
- picker
- fscrypt/orphans.go
- IntegrationSuite
- Message
- controlServer
- github.com/charmbracelet/bubbles/textarea.Model
- Prereq
- IntegrationSuite
- absPath
- monitorUnitTest
- UnifiedDiff
- fileedit_model_test.go
- journal
- stdio
- update/update.go
- fadeScroll
- vaultfs.go
- cursor.go
- Client
- daemonconfig.go
- hashCache
- .apply
- NewDumpHook
- mutedText
- GuardInodeKey
- runNetworkGuard
- serveWebSocket
- Recorder
- exe_supersede.h
- plan_test.go
- UI
- render.go
- gui/session.go
- IntegrationSuite
- uninstall_system.go
- liveSession
- editView
- .refreshOnce
- FileEvent
- string
- bpfstats/main.go
- CandidateDir
- codecheck/main.go
- shellrc.go
- update_test.go
- process_vm_readv.c
- GuardEvent
- op
- Candidate
- updatePage
- runGuard
- Monitor
- planUpdaters
- dialLiveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- .populateWatchInodes
- newTestTarget
- prctl_exe_swap.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- ConfinedEntry
- pinstate.go
- Well-known bypass classes
- NetEventType
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- downloadAndVerify
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- NetworkMonitor
- mountTable
- highlight_test.go
- recordFlow
- IntegrationSuite
- procstats.go
- newTestVetter
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- .TestRecord_GuardParity
- GenerateConf
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- .filterEvent
- make test-integration (rootful Docker suite)
- IntegrationSuite
- install/install.go
- LibraryClosure
- trace-app-libs.sh
- IntegrationSuite
- maintainerKeeper
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- catalogRefresher
- Serve
- github.com/spf13/cobra.Command
- sources_test.go
- Ledger
- KernelDev
- auditGroupCoverage
- changelog.go
- Config
- fakeStore
- options.go
- IntegrationSuite
- buildOneGuard
- fakeDaemon
- sync.Mutex
- runView
- Installation
- backend_test.go
- graphify-refresh.sh
- layout.go
- headless/headless.go
- IntegrationSuite
- BinaryStat
- feed.go
- shell
- BinaryEntry
- watchSet
- WithBottomBar
- Run
- fetchReleases
- Limitations
- net_tester/main.go
- IntegrationSuite
- hardlinks_test.go
- netinject.c
- copyRegularAt
- runMonitor
- TestMain
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
- Binary
- sources
- edit-protected
- NewEventFanout
- kernel-prof.sh
- testLimiter
- pathCache
- github.com/stretchr/testify/suite.Suite
- DaemonEvent
- Daemon configuration
- NetGuardEvent
- uninstall/uninstall.go
- Desktop GUI
- IntegrationSuite
- time.Time
- IntegrationSuite
- .eventTypeAllowed
- configedit_test.go
- inode_dev
- prctl_gate
- .awaitDaemonUpOrDead
- linkView
- SanitizeText
- WithSuspectsFrom
- ParseEventsFlag
- IntegrationSuite
- ResolvePinBase
- Vault
- recordForm
- supersede_probe.c
- inspector_probe.c
- links_test.go
- github.com/Virgula0/app-listener/tools/codecheck
- ValidateTargets
- taint_victim.c
- Troubleshooting
- VouchedDevs
- bleedLayout
- TestCleanRel
- monitor.md

## God Nodes (most connected - your core abstractions)
1. `Guard` - 123 edges
2. `Config` - 118 edges
3. `GuardInodeKey` - 103 edges
4. `shell` - 97 edges
5. `Load()` - 90 edges
6. `absPath()` - 89 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 69 edges
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

## Communities (235 total, 42 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.09
Nodes (31): tempJournalRow, Progress, User, classCacheKey, launchRule, MulticallError, trustEvent, stageItem (+23 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (156): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), awaitStartupOrSignal(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestAwaitStartupOrSignalAbortsAndLocksDown(), TestAwaitStartupOrSignalAbortWithFailedStartup() (+148 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "fscrypt.go"
Cohesion: 0.05
Nodes (47): confirmMasterKeyOverwrite(), runGenKey(), askEncryption(), deprovisionKind, migrationTarget, isFileVaultCiphertext(), isRegularFileTarget(), writeTempWithMetadata() (+39 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (75): TestBuildTrustedSetIncludesPendingBinaries(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+67 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (45): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_user(), fill_identity() (+37 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (22): addErrKind, ReplacementCheck, classifyAddErr(), addRecorder(), copySelf(), guardUnitTest, runTool(), walkEntries() (+14 more)

### Community 9 - "mkdirs"
Cohesion: 0.09
Nodes (41): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+33 more)

### Community 10 - "fakeGuardRepo"
Cohesion: 0.05
Nodes (26): forwardCheck(), tempMaskOp, tempRow, binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps() (+18 more)

### Community 11 - "New"
Cohesion: 0.08
Nodes (46): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault() (+38 more)

### Community 12 - "protectedPage"
Cohesion: 0.04
Nodes (29): actionGate, backendLink, configView, forwardForm, forwardView, installUI, multiAnswer, protectedPage (+21 more)

### Community 13 - "os.File"
Cohesion: 0.04
Nodes (70): setConfinedHomes(), VettedInode, Multicall, MulticallKind, childEnv(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted() (+62 more)

### Community 14 - "go_pkg_path_filepath"
Cohesion: 0.36
Nodes (4): recordTick, isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 15 - "Network"
Cohesion: 0.14
Nodes (5): NetworkGuardRepository, clonePolicy(), NetExtra, NetPolicy, Network

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (6): eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (57): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+49 more)

### Community 18 - "fadescroll.go"
Cohesion: 0.03
Nodes (22): appTheme, clickPanel, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, tabBar (+14 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.12
Nodes (54): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+46 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (14): backingDeviceUnion(), buildGuards(), Resource, BackingDevices(), GuardRepository, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase (+6 more)

### Community 23 - "auth_test.go"
Cohesion: 0.10
Nodes (46): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+38 more)

### Community 24 - "Guard"
Cohesion: 0.08
Nodes (9): exeRow, rootHandle, supersededBinary, openRoot(), SharedPinDegraded(), TestNoteRootMountFollowsMountID(), Guard, TestReadOnlyWalkToleratesDanglingSymlinks() (+1 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.11
Nodes (59): add_inode_to_guard(), btrfs_copy_gate(), chmod_only_drops_write(), code_suspect(), count_degrade(), discover_guarded_parent(), emit_process_denial(), emit_with() (+51 more)

### Community 26 - "TrustGuard"
Cohesion: 0.04
Nodes (34): lsmHook, statFunc, tempGrant, trustHook, trustRows, checkBtrfsKeys(), ensureBtrfsLayout(), deleteResKeys() (+26 more)

### Community 27 - "find_test.go"
Cohesion: 0.18
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 28 - "FileTarget"
Cohesion: 0.10
Nodes (16): Target, guardPill(), newExeResolver(), Files, FileSpec, FileTarget, Policy, guardRow() (+8 more)

### Community 29 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.11
Nodes (24): placedWalk, treeJudge, FirstOnBtrfs(), OnBtrfs(), SuperblockDev(), TestSuperblockDevMatchesStatOffBtrfs(), confineBelow(), inUserHome() (+16 more)

### Community 30 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.06
Nodes (20): recordTickCmd(), recordModel, diffModel, editorModel, listenForDaemonEvents(), tickDaemonStats(), protoString(), renderBullets() (+12 more)

### Community 31 - "MonitorRepository"
Cohesion: 0.13
Nodes (13): MonitorRepository, NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+5 more)

### Community 32 - "binaryList"
Cohesion: 0.10
Nodes (14): binaryList, netGuardTab, policyEditor, policyHooks, policyTab, newBinaryList(), newPolicyTab(), verbFor() (+6 more)

### Community 33 - "EventType"
Cohesion: 0.12
Nodes (13): BpfEvent, EventType, eventUnitTest, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), specFromForm(), specSummary(), Cstr() (+5 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (30): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+22 more)

### Community 35 - "fscrypt/orphans.go"
Cohesion: 0.19
Nodes (20): modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors() (+12 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (15): checkSecret(), checkTrustLine(), trustRows(), newRecordFlow(), File, TrustItem, backend, backendOp (+7 more)

### Community 38 - "controlServer"
Cohesion: 0.08
Nodes (25): readConfigPut(), TestReadGrantRequestConfig(), newControlManager(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest() (+17 more)

### Community 39 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.08
Nodes (33): class, findButton, pos, span, classOf(), column(), dedent(), Edit() (+25 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.07
Nodes (14): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest, IntegrationSuite (+6 more)

### Community 42 - "absPath"
Cohesion: 0.06
Nodes (8): IntegrationSuite, IntegrationSuite, shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite, IntegrationSuite, absPath()

### Community 44 - "UnifiedDiff"
Cohesion: 0.18
Nodes (13): applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig(), ConfirmOverwrite() (+5 more)

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "journal"
Cohesion: 0.17
Nodes (8): fakeEngines, fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal, NetGuardConfig

### Community 48 - "update/update.go"
Cohesion: 0.16
Nodes (25): applyUpdate(), channelOf(), checkUpdatePreconditions(), confirmUnit(), deploy(), deployGUI(), FetchLatest(), fetchRelease() (+17 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (27): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+19 more)

### Community 50 - "vaultfs.go"
Cohesion: 0.07
Nodes (37): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+29 more)

### Community 51 - "cursor.go"
Cohesion: 0.08
Nodes (23): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+15 more)

### Community 52 - "Client"
Cohesion: 0.10
Nodes (11): Client, launch(), TestHelperBackend(), TestReaderCapsLines(), Reader, Writer, NewReader(), NewWriter() (+3 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.13
Nodes (42): systemPatterns(), systemRule(), watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp() (+34 more)

### Community 54 - "hashCache"
Cohesion: 0.21
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 55 - ".apply"
Cohesion: 0.18
Nodes (5): applier, ledgerOps, PlannedBinary, Store, vaultOps

### Community 56 - "NewDumpHook"
Cohesion: 0.12
Nodes (16): applyVerbosity(), logHook, VerboseLevel, BridgeStdLog(), formatEntry(), NewDumpHook(), sortedKeys(), sortStrings() (+8 more)

### Community 57 - "mutedText"
Cohesion: 0.05
Nodes (49): binaryRow, binSource, checkGrid, filesPage, handCheck, homePage, iconTile, installAction (+41 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (54): engine, inodeBatchFn, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel() (+46 more)

### Community 59 - "runNetworkGuard"
Cohesion: 0.11
Nodes (24): CheckBPFLSM(), CheckEBPF(), prepareDaemonStart(), runBPFCheck(), runOneShotMode(), computeGuardBinaries(), confirmUnsafe(), discoverInfra() (+16 more)

### Community 60 - "serveWebSocket"
Cohesion: 0.16
Nodes (12): configureWebSocket(), decodeEndsAtEOF(), parseResize(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestParseResize(), TestSameOrigin() (+4 more)

### Community 61 - "Recorder"
Cohesion: 0.14
Nodes (16): drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start(), record() (+8 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.20
Nodes (16): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+8 more)

### Community 63 - "plan_test.go"
Cohesion: 0.27
Nodes (14): isWhiteout(), opaque(), scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent(), TestScanAndValidateAcceptsLibDirContent(), TestScanRefusesSpecialFiles() (+6 more)

### Community 64 - "UI"
Cohesion: 0.09
Nodes (19): showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), cancelUI, TestEditTextMapsUICancel(), MultiSelectKeymap() (+11 more)

### Community 65 - "render.go"
Cohesion: 0.13
Nodes (12): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+4 more)

### Community 66 - "gui/session.go"
Cohesion: 0.07
Nodes (40): ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+32 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "uninstall_system.go"
Cohesion: 0.05
Nodes (74): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), entryWriters(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor() (+66 more)

### Community 69 - "liveSession"
Cohesion: 0.16
Nodes (13): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), holdForward() (+5 more)

### Community 70 - "editView"
Cohesion: 0.12
Nodes (8): codeEntry, editView, Entry, fileSize(), joinRel(), newCodeEntry(), newEditView(), relOf()

### Community 71 - ".refreshOnce"
Cohesion: 0.43
Nodes (6): sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad()

### Community 72 - "FileEvent"
Cohesion: 0.16
Nodes (5): snapMonitor, FileEvent, newFakeMonitorRepo(), fakeSource, fakeMonitorRepo

### Community 73 - "string"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 75 - "CandidateDir"
Cohesion: 0.05
Nodes (42): reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), confSafeMatch() (+34 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "shellrc.go"
Cohesion: 0.09
Nodes (56): setupBunTmpdir(), readAtNoFollow(), upsertUnitAt(), writeAtNoFollow(), rcTarget, invokingIDs(), TestWriteUserFileRefusesSymlinks(), writeUserFile() (+48 more)

### Community 78 - "update_test.go"
Cohesion: 0.09
Nodes (22): newerThanInstalled(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleaseByTag(), TestFilterChannel(), testKeyPair() (+14 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.07
Nodes (18): holdHeadless(), HeadlessLine(), bpfGuardEvent, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths(), commMatchesGuardedBinary() (+10 more)

### Community 81 - "op"
Cohesion: 0.17
Nodes (13): clearExcept(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir(), removeAll() (+5 more)

### Community 82 - "Candidate"
Cohesion: 0.09
Nodes (34): appendSectionsAndEdit(), collectDiffAdditions(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), cleanOrphanedFscrypt(), collectFilesystemPrereqs(), addManualDirectories() (+26 more)

### Community 83 - "updatePage"
Cohesion: 0.26
Nodes (5): latestRelease, updatePage, newUpdatePage(), boolIndex(), isPreRelease()

### Community 84 - "runGuard"
Cohesion: 0.22
Nodes (10): DrainUntilSignal(), displayEntries(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI(), runHeadless(), runHeadless() (+2 more)

### Community 86 - "Monitor"
Cohesion: 0.18
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 87 - "planUpdaters"
Cohesion: 0.14
Nodes (18): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+10 more)

### Community 88 - "dialLiveSession"
Cohesion: 0.08
Nodes (26): auditAfterEdit(), auditFindings(), reserveStdout(), runBackend(), dialLiveSession(), LiveModeAvailable(), requireRoot(), runEditProtected() (+18 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 92 - "newTestTarget"
Cohesion: 0.24
Nodes (15): NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded(), TestApplyFromWhitelistStopsFirst() (+7 more)

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.18
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.09
Nodes (13): foreignPinnedLink(), PumpRingbuf(), reasonText(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses() (+5 more)

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

### Community 100 - "ConfinedEntry"
Cohesion: 0.10
Nodes (33): markReady(), notifySystemdReady(), binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared() (+25 more)

### Community 101 - "pinstate.go"
Cohesion: 0.10
Nodes (36): configureDaemonLogging(), loadDaemonConfig(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), resolveConfigPath(), runLockdown(), ensurePinStateFilePlaceholder() (+28 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "NetEventType"
Cohesion: 0.07
Nodes (16): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), FormatAddr(), NetEventTypes(), Ntohs() (+8 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.45
Nodes (13): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+5 more)

### Community 107 - "downloadAndVerify"
Cohesion: 0.22
Nodes (15): archAssets(), assetsFor(), downloadAndVerify(), fetchGUI(), fetchPair(), guiAssetsFor(), hasAsset(), resolveAssets() (+7 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.17
Nodes (9): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), btrfsTypes(), TestBtrfsLayoutFromNestedMembers() (+1 more)

### Community 110 - "NetworkMonitor"
Cohesion: 0.15
Nodes (7): BinariesSummary(), TestBinariesSummary(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 111 - "mountTable"
Cohesion: 0.18
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo(), userMountable() (+1 more)

### Community 112 - "highlight_test.go"
Cohesion: 0.13
Nodes (26): New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines(), TestSetTextStartsAtFirstLine() (+18 more)

### Community 113 - "recordFlow"
Cohesion: 0.13
Nodes (16): encryptBlocker(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord(), runRecordView() (+8 more)

### Community 115 - "procstats.go"
Cohesion: 0.38
Nodes (8): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic()

### Community 116 - "newTestVetter"
Cohesion: 0.29
Nodes (11): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+3 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "GenerateConf"
Cohesion: 0.05
Nodes (68): logRefreshChange(), TestLogRefreshChange(), TestDiffMergePreservesExistingAndParses(), refreshCatalogConfig(), groupedDiscordConf(), TestGroupedSectionAddressedByEncryptionRoot(), TestPatchCatalogSectionGroupedConfig(), TestSetSectionWhitelistPreservesGroupStructure() (+60 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "launch_probe.c"
Cohesion: 0.48
Nodes (5): copy(), is_verb(), main(), read_once(), usage()

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
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "install/install.go"
Cohesion: 0.04
Nodes (122): releaseBuild(), Run(), anyDaemonRunning(), applyDiffAdditions(), offerDaemonStart(), runDiffCatalog(), installGUI(), needsBinary() (+114 more)

### Community 131 - "LibraryClosure"
Cohesion: 0.07
Nodes (28): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+20 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "maintainerKeeper"
Cohesion: 0.19
Nodes (10): appRunning(), libDirSig(), grantNewBinaries(), maintainedApp, maintainerKeeper, Line(), Lines(), TestLinesKeyOnLinkAndIncludePending() (+2 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "catalogRefresher"
Cohesion: 0.14
Nodes (8): drain(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), catalogRefresher, debouncer

### Community 143 - "Serve"
Cohesion: 0.16
Nodes (17): basicAuth(), isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), securityHeaders(), Serve() (+9 more)

### Community 144 - "github.com/spf13/cobra.Command"
Cohesion: 0.12
Nodes (25): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+17 more)

### Community 145 - "sources_test.go"
Cohesion: 0.20
Nodes (9): SystemTrusted(), systemTool(), archiveMembers(), decompressed(), newSources(), tarMembers(), TestArchiveMembersIgnoresNonArchives(), TestSourcesHoldFilesAndArchiveMembers() (+1 more)

### Community 146 - "Ledger"
Cohesion: 0.07
Nodes (25): candidates(), confirmAll(), describe(), record(), run(), applySchema(), Ledger, Pending (+17 more)

### Community 147 - "KernelDev"
Cohesion: 0.18
Nodes (14): KernelDev(), BtrfsDevices(), btrfsSysfsDevices(), TestBtrfsSysfsDevices(), codeVouched(), isTraced(), loaderEnv(), readMountDevs() (+6 more)

### Community 148 - "auditGroupCoverage"
Cohesion: 0.22
Nodes (11): auditEntry(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths(), relOrPath() (+3 more)

### Community 149 - "changelog.go"
Cohesion: 0.17
Nodes (22): changelogText(), cumulativeChangelog(), newChangelogModel(), newReleaseWindow(), releasesSince(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit() (+14 more)

### Community 150 - "Config"
Cohesion: 0.08
Nodes (36): newCatalogRefresher(), makeReloadHandler(), noteReloadedSettings(), startCatalogRefresh(), inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted() (+28 more)

### Community 152 - "options.go"
Cohesion: 0.26
Nodes (11): installHooks(), staged(), InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks (+3 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.07
Nodes (39): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards(), lockVaultFully() (+31 more)

### Community 156 - "sync.Mutex"
Cohesion: 0.18
Nodes (7): fyneLogFilter, newActivity(), TestRingKeepsNewest(), Counts, newRing(), activity, ring

### Community 157 - "runView"
Cohesion: 0.06
Nodes (12): daemonPage, daemonState, logGrid, logLine, runHost, runView, journalLine(), TestJournalLine() (+4 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 161 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 162 - "headless/headless.go"
Cohesion: 0.29
Nodes (8): stagedBinary, stageRelease(), Check(), CheckGUI(), inspect(), buildTiny(), TestCheck(), TestCheckGUI()

### Community 164 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 165 - "feed.go"
Cohesion: 0.06
Nodes (43): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+35 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (56): fakeLink, pageDef, pageID, reel, resizeWatch, shell, sizedDialog, snapEngines (+48 more)

### Community 167 - "BinaryEntry"
Cohesion: 0.07
Nodes (20): deferredBinary, keptEntries, BinaryEntry, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), canonicalBinaryPath(), confinedEntry(), eventMask() (+12 more)

### Community 168 - "watchSet"
Cohesion: 0.15
Nodes (13): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+5 more)

### Community 169 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 170 - "Run"
Cohesion: 0.33
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 171 - "fetchReleases"
Cohesion: 0.22
Nodes (8): downloadFile(), downloadReleaseFiles(), fetchReleases(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleases(), TestFetchReleasesHTTPError(), progressReader

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 175 - "hardlinks_test.go"
Cohesion: 0.33
Nodes (8): lstatOf(), outsideHardlinks(), refuseOutsideHardlinks(), TestRefuseOutsideHardlinks_FileRoot(), TestRefuseOutsideHardlinks_InTreeLinksAccepted(), TestRefuseOutsideHardlinks_MissingRootFails(), TestRefuseOutsideHardlinks_OutsideLinkRefused(), writeFile()

### Community 176 - "netinject.c"
Cohesion: 0.50
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 177 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 178 - "runMonitor"
Cohesion: 0.19
Nodes (14): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+6 more)

### Community 179 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

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
Cohesion: 0.25
Nodes (8): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, The dynamic loader, trust

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
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 194 - "Binary"
Cohesion: 0.12
Nodes (27): lines(), pickBinaries(), pickLibs(), reviewRecording(), appletLine(), beyondEvents(), Classify(), classifyGroup() (+19 more)

### Community 195 - "sources"
Cohesion: 0.27
Nodes (4): pathByte(), pathRecord(), TestPathRecord(), sources

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "NewEventFanout"
Cohesion: 0.24
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 199 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 202 - "DaemonEvent"
Cohesion: 0.12
Nodes (28): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), newDaemonModel(), runDaemon(), runDaemonUI(), runHeadless() (+20 more)

### Community 203 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 204 - "NetGuardEvent"
Cohesion: 0.29
Nodes (3): snapNetGuard, NetGuardEvent, netGuardEventLine

### Community 205 - "uninstall/uninstall.go"
Cohesion: 0.13
Nodes (29): deletePostBackups(), listBackups(), restoreBackups(), askDecrypt(), cleanOrphanedMetadata(), confirmKeyDeletion(), decryptDirectories(), decryptStep() (+21 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.13
Nodes (15): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+7 more)

### Community 209 - "time.Time"
Cohesion: 0.05
Nodes (42): LockedError, runNetworkMonitor(), runTUI(), sweepID, sweepMemo, Progress, formatGuardEventLine(), sanitizeTerminalText() (+34 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 213 - "inode_dev"
Cohesion: 0.15
Nodes (18): exe_is_superseded(), exe_refused(), ctx_ptr(), inode_dev(), sb_dev(), check_and_emit_args(), current_is_inspector(), exe_tag_bits() (+10 more)

### Community 214 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 216 - "linkView"
Cohesion: 0.43
Nodes (3): linkTarget(), entryKind, linkView

### Community 217 - "SanitizeText"
Cohesion: 0.09
Nodes (22): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), scriptIdentity(), drainEphemeral(), newPinGeneration(), reloadOnce() (+14 more)

### Community 218 - "WithSuspectsFrom"
Cohesion: 0.43
Nodes (3): WithSuspectsFrom(), netguardHarness, Option

### Community 219 - "ParseEventsFlag"
Cohesion: 0.33
Nodes (4): ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents()

### Community 221 - "ResolvePinBase"
Cohesion: 0.50
Nodes (5): dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 227 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 229 - "links_test.go"
Cohesion: 0.54
Nodes (7): newLinkView(), checkLink(), linkFixture(), TestCheckLinksAccepts(), TestCheckLinksJudgesEachTarget(), TestCheckLinksRefusesUserFiles(), TestCheckLinksSeesThePlannedTree()

### Community 232 - "ValidateTargets"
Cohesion: 0.10
Nodes (18): BinaryRule, dedupe(), existingSection(), Plan, guardedMatch(), libOfChosen(), materialized(), overlaps() (+10 more)

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 239 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 244 - "monitor.md"
Cohesion: 0.50
Nodes (3): Events, Flags, monitor: see what touches a path

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 639 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **42 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Guard` connect `Guard` to `go_pkg_strings`, `TrustGuard`, `GuardInodeKey`, `time.Duration`, `controlServer`, `ConfinedEntry`, `guardUnitTest`, `BinaryEntry`, `os.File`, `GuardEvent`, `time.Time`, `daemonconfig.go`, `SanitizeText`, `buildOneGuard`, `fakeDaemon`, `sync.Mutex`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.0928391804427588 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `go_pkg_strings`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `TestMain`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.01802620616526547 - nodes in this community are weakly interconnected._
- **Why does `fadeScroll` connect `fadeScroll` to `mutedText`, `fadescroll.go`, `time.Time`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09125085440874915 - nodes in this community are weakly interconnected._