# Graph Report - app-listener  (2026-10-27)

## Corpus Check
- 525 files · ~1,741,490 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6835 nodes · 30000 edges · 250 communities (197 shown, 53 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2771 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `d9fc9b82`
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
- fakeGuardRepo
- go_pkg_path_filepath
- protectedPage
- find_test.go
- github.com/charmbracelet/bubbletea.Cmd
- installPage
- IntegrationSuite
- guard_trust.bpf.c
- runView
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- Guard
- guard.bpf.c
- TrustGuard
- github.com/charmbracelet/bubbles/textarea.Model
- sync.Mutex
- .step
- daemon_multicall_test.go
- TestMonitorUseCaseLifecycle
- policyEditor
- highlight_test.go
- picker
- Candidate
- time.Duration
- Message
- controlServer
- runDaemon
- NetworkMonitor
- IntegrationSuite
- absPath
- monitorUnitTest
- runNetworkGuard
- fileedit_model_test.go
- NetGuardEvent
- os.File
- fadeScroll
- vaultfs.go
- cursor.go
- editbackend_test.go
- daemonconfig.go
- fyne.io/fyne/v2/widget.Label
- .apply
- KernelDev
- mutedText
- GuardInodeKey
- auth_test.go
- .netIntegrityContainer
- Recorder
- exe_supersede.h
- NetEventType
- DurableStat
- system.go
- gui/session.go
- IntegrationSuite
- expandPlaceholders
- liveSession
- codeEntry
- journal
- changelog_test.go
- fcntl
- bpfstats/main.go
- BinaryStat
- codecheck/main.go
- shellrc.go
- deploy
- process_vm_readv.c
- GuardEvent
- maintainer/plan.go
- orphans_test.go
- maintainerKeeper
- net_tester/main.go
- render.go
- finishBinary
- planUpdaters
- runBackend
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- EventType
- Monitor
- prctl_exe_swap.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- pulse_test.go
- BinaryEntry
- pinstate.go
- Well-known bypass classes
- configedit_test.go
- jit_provenance.c
- inspector_probe.c
- recordops.c
- stripTempJournal
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- batchInodeRows
- mountTable
- dialLiveSession
- recordFlow
- IntegrationSuite
- NewDumpHook
- editauth_test.go
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
- RefreshSection
- LibraryClosure
- trace-app-libs.sh
- IntegrationSuite
- Binary
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- Config
- Serve
- catalogRefresher
- UI
- Ledger
- eventUnitTest
- IntegrationSuite
- readGrantRequest
- buildTrustedSet
- fakeStore
- inodeChain
- IntegrationSuite
- buildOneGuard
- safeio.go
- procstats.go
- runMonitor
- Installation
- backend_test.go
- graphify-refresh.sh
- layout.go
- IntegrationSuite
- .populateWatchInodes
- feed.go
- shell
- newTestTarget
- watchSet
- SanitizeText
- fileedit_symlink_test.go
- Active
- Limitations
- github.com/stretchr/testify/suite.Suite
- IntegrationSuite
- testLimiter
- netinject.c
- CopyTreeWithProgress
- main_test.go
- Compatibility
- How it works
- securewrite_test.go
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- changelog.go
- RawTarget
- edit-protected
- github.com/spf13/cobra.Command
- kernel-prof.sh
- update_test.go
- pathCache
- .exeKeys
- DaemonEvent
- Plan
- ValidateTargets
- uninstall/uninstall.go
- Desktop GUI
- IntegrationSuite
- time.Time
- downloadAndVerify
- .eventTypeAllowed
- desktop_test.go
- WithBottomBar
- root.go
- .awaitDaemonUpOrDead
- NetworkMonitorUseCase
- promptPassword
- tempGrantSetup
- Question
- shQuote
- ptrace_race.c
- prctl_gate
- globBuilder
- RunWith
- supersede_probe.c
- WithSuspectsFrom
- ParseEventsFlag
- inode_dev.h
- formatEntry
- github.com/Virgula0/app-listener/tools/codecheck
- parsePasswd
- taint_victim.c
- TestCatalogSanity
- Troubleshooting
- bpfLsmListed
- monitor.md

## God Nodes (most connected - your core abstractions)
1. `Guard` - 122 edges
2. `Config` - 118 edges
3. `shell` - 91 edges
4. `Load()` - 90 edges
5. `absPath()` - 88 edges
6. `IntegrationSuite` - 83 edges
7. `fileEditModel` - 77 edges
8. `GuardInodeKey` - 72 edges
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
- `vaultOpUnderGuard()` --calls--> `op`  [INFERRED]
  cmd/functions/install/steps.go → internal/maintainer/plan.go

## Import Cycles
- None detected.

## Communities (250 total, 53 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.12
Nodes (34): atomicWriteAt(), writeAndSync(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs (+26 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (143): plantBun(), TestOpenBunTmp_KeepsVettedDir(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), loadDaemonConfig(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod() (+135 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "filevault.go"
Cohesion: 0.05
Nodes (54): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), deprovisionKind, migrationTarget, classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey() (+46 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (73): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+65 more)

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
Cohesion: 0.05
Nodes (21): addErrKind, Mode, TestWithChmodDropWriteNeedsReadOnly(), canonicalPaths(), classifyAddErr(), guardModeKey(), modeString(), NewGuard() (+13 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (34): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+26 more)

### Community 10 - "fakeGuardRepo"
Cohesion: 0.06
Nodes (21): forwardRequest, SupersededBinary, binaryStatOf(), TempBinary, TemporaryGrant, TempRule, planTempBlock(), ResolveTempBinary() (+13 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (23): backendLink, configView, editView, forwardForm, forwardView, installUI, protectedPage, recordView (+15 more)

### Community 13 - "find_test.go"
Cohesion: 0.18
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 14 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.05
Nodes (37): recordTickCmd(), recordModel, diffModel, editorModel, listenForDaemonEvents(), NewDaemonModel(), tickDaemonStats(), listenForGuardEvents() (+29 more)

### Community 15 - "installPage"
Cohesion: 0.07
Nodes (16): actionGate, binSource, installAction, installedBinary, installPage, latestRelease, stagedBinary, updatePage (+8 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (59): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+51 more)

### Community 18 - "runView"
Cohesion: 0.06
Nodes (12): daemonPage, daemonState, logHook, logLine, runHost, runView, journalLine(), TestJournalLine() (+4 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (49): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+41 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (10): Resource, GuardRepository, Vault, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots() (+2 more)

### Community 23 - "auth.go"
Cohesion: 0.13
Nodes (24): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), MigrateLegacyHash(), openLedger(), OriginOf(), parseHash() (+16 more)

### Community 24 - "Guard"
Cohesion: 0.05
Nodes (21): binaryVerifyState, deferredBinary, exeRow, rootHandle, tempMaskOp, tempRow, openRoot(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted() (+13 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "TrustGuard"
Cohesion: 0.05
Nodes (28): tempGrant, trustRows, deleteResKeys(), supersedeUnnamed(), dropSupersededLocked(), exeGone(), freedSuperseded(), liftSuperseded() (+20 more)

### Community 27 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.09
Nodes (31): class, findButton, pos, span, classOf(), column(), dedent(), Edit() (+23 more)

### Community 28 - "sync.Mutex"
Cohesion: 0.04
Nodes (34): fyneLogFilter, specFromForm(), guardPill(), specSummary(), netGuardPill(), EventTypes(), maskFor(), newActivity() (+26 more)

### Community 29 - ".step"
Cohesion: 0.18
Nodes (13): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+5 more)

### Community 30 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 32 - "policyEditor"
Cohesion: 0.14
Nodes (10): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), confirmWarnings() (+2 more)

### Community 33 - "highlight_test.go"
Cohesion: 0.13
Nodes (26): New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines(), TestSetTextStartsAtFirstLine() (+18 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (30): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+22 more)

### Community 35 - "Candidate"
Cohesion: 0.12
Nodes (26): appendSectionsAndEdit(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories(), pickFromCandidates() (+18 more)

### Community 36 - "time.Duration"
Cohesion: 0.25
Nodes (3): IntegrationSuite, pooledContainer, formatClock()

### Community 37 - "Message"
Cohesion: 0.10
Nodes (11): checkSecret(), removeSSHAgentEnv(), File, backend, backendOp, backendState, fakeLink, BinaryCheck (+3 more)

### Community 38 - "controlServer"
Cohesion: 0.13
Nodes (10): swapFile(), readClientLines(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), configEditor, controlServer (+2 more)

### Community 39 - "runDaemon"
Cohesion: 0.09
Nodes (24): newControlManager(), startControlManager(), awaitStartupOrSignal(), backingDeviceUnion(), buildGuards(), catchLifecycleSignals(), reloadOnce(), reloadSlotsNeeded() (+16 more)

### Community 40 - "NetworkMonitor"
Cohesion: 0.13
Nodes (8): BinariesSummary(), TestBinariesSummary(), PumpRingbuf(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "runNetworkGuard"
Cohesion: 0.22
Nodes (13): DrainUntilSignal(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI() (+5 more)

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "NetGuardEvent"
Cohesion: 0.13
Nodes (7): snapNetGuard, NetGuardEvent, protoString(), NetworkGuardRepository, netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 48 - "os.File"
Cohesion: 0.03
Nodes (88): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), setConfinedHomes() (+80 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (28): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+20 more)

### Community 50 - "vaultfs.go"
Cohesion: 0.10
Nodes (21): DefaultNewFileMode(), Vault, IsBinaryContent(), IsOctalMode(), Open(), ParseMode(), inodeOf(), mustRead() (+13 more)

### Community 51 - "cursor.go"
Cohesion: 0.09
Nodes (15): buttonCursor, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor(), showConfirm() (+7 more)

### Community 52 - "editbackend_test.go"
Cohesion: 0.08
Nodes (18): Client, childEnv(), launch(), Start(), helper(), TestChildEnvIsAllowlisted(), TestCleanRel(), TestClientCallsAndEvents() (+10 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.11
Nodes (45): systemPatterns(), systemRule(), watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp() (+37 more)

### Community 54 - "fyne.io/fyne/v2/widget.Label"
Cohesion: 0.03
Nodes (19): appTheme, binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panel (+11 more)

### Community 55 - ".apply"
Cohesion: 0.11
Nodes (14): TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan(), TestPlanAdmitsOnlyWhatWasRecorded(), TestPlanExtendsAnExistingSection(), testRecording() (+6 more)

### Community 56 - "KernelDev"
Cohesion: 0.11
Nodes (23): fdStat, physicalParent(), statFD(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed() (+15 more)

### Community 57 - "mutedText"
Cohesion: 0.06
Nodes (53): binaryList, checkGrid, filesPage, handButton, handCheck, handCheckGroup, homePage, recordForm (+45 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.04
Nodes (43): engine, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+35 more)

### Community 59 - "auth_test.go"
Cohesion: 0.31
Nodes (17): LoadHash(), mustHash(), TestCheckPasswordSharesGuessBudget(), TestMigrateLegacyHashImportFailureRemoves(), TestMigrateLegacyHashImports(), TestMigrateLegacyHashMalformedKeepsStored(), TestMigrateLegacyHashNeverReplaces(), TestMigrateLegacyHashRefusesOddFiles() (+9 more)

### Community 60 - ".netIntegrityContainer"
Cohesion: 0.28
Nodes (5): IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 61 - "Recorder"
Cohesion: 0.16
Nodes (15): Target, drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start() (+7 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.19
Nodes (16): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+8 more)

### Community 63 - "NetEventType"
Cohesion: 0.10
Nodes (13): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), FormatAddr(), NetEventTypes(), Ntohs() (+5 more)

### Community 64 - "DurableStat"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 65 - "system.go"
Cohesion: 0.09
Nodes (34): TestDaemonUnitWithMetadataOutput(), dupFD(), installFileAs(), installServices(), installSSHAgent(), installSSHAgentEnv(), installUnit(), readAtNoFollow() (+26 more)

### Community 66 - "gui/session.go"
Cohesion: 0.07
Nodes (40): invokingIDs(), ancestorEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+32 more)

### Community 68 - "expandPlaceholders"
Cohesion: 0.06
Nodes (29): symlinkTargetGlob(), CandidateDir, confSafeMatch(), expandPlaceholders(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestExpandWhitelist() (+21 more)

### Community 69 - "liveSession"
Cohesion: 0.16
Nodes (12): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), forwardCheck() (+4 more)

### Community 71 - "journal"
Cohesion: 0.15
Nodes (6): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 72 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 73 - "fcntl"
Cohesion: 0.09
Nodes (11): main(), report(), denied(), main(), put(), wait_then_read(), dump_via_debugfs(), main() (+3 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 75 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "shellrc.go"
Cohesion: 0.09
Nodes (57): setupBunTmpdir(), removeFacultative(), removeFacultatives(), TestRemoveSSHAgentFacultativeKeepsGuardedSSHConfig(), detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeBunLauncherEnv(), removeSSHAgentUnit() (+49 more)

### Community 78 - "deploy"
Cohesion: 0.12
Nodes (25): applyUpdate(), archAssets(), assetsFor(), checkUpdatePreconditions(), deploy(), deployGUI(), fetchGUI(), guiAssetsFor() (+17 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.19
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.07
Nodes (24): displayEntries(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI(), bpfGuardEvent (+16 more)

### Community 81 - "maintainer/plan.go"
Cohesion: 0.06
Nodes (51): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+43 more)

### Community 82 - "orphans_test.go"
Cohesion: 0.21
Nodes (19): appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors(), orphanedPolicies() (+11 more)

### Community 83 - "maintainerKeeper"
Cohesion: 0.19
Nodes (10): appRunning(), libDirSig(), grantNewBinaries(), maintainedApp, maintainerKeeper, Line(), Lines(), TestLinesKeyOnLinkAndIncludePending() (+2 more)

### Community 84 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 85 - "render.go"
Cohesion: 0.11
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+6 more)

### Community 86 - "finishBinary"
Cohesion: 0.19
Nodes (12): beyondEvents(), finishBinary(), isShellHelper(), newBinary(), noteWhitelisted(), opNames(), sortedKeys(), TestFileReadersStartDeselected() (+4 more)

### Community 87 - "planUpdaters"
Cohesion: 0.09
Nodes (26): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+18 more)

### Community 88 - "runBackend"
Cohesion: 0.36
Nodes (6): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive()

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "EventType"
Cohesion: 0.10
Nodes (10): snapMonitor, EventType, FileEvent, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), newFakeMonitorRepo(), accessor, fakeSource (+2 more)

### Community 92 - "Monitor"
Cohesion: 0.17
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.18
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.09
Nodes (14): lsmHook, trustHook, foreignPinnedLink(), NewMountVouch(), reasonText(), eventTypeKey(), judgePasses(), parseProcs() (+6 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "pulse_test.go"
Cohesion: 0.42
Nodes (6): New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Sub

### Community 100 - "BinaryEntry"
Cohesion: 0.08
Nodes (37): HashStore, keptEntries, sharedHashEntry, BinaryEntry, checkBtrfsKeys(), canonicalBinaryPath(), ConfinedEntry(), claimHash() (+29 more)

### Community 101 - "pinstate.go"
Cohesion: 0.29
Nodes (14): pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID(), TestPinStateRoundTrip(), TestReadPinStateMalformed() (+6 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.20
Nodes (14): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+6 more)

### Community 106 - "recordops.c"
Cohesion: 0.45
Nodes (13): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+5 more)

### Community 107 - "stripTempJournal"
Cohesion: 0.24
Nodes (10): ensurePinStateFilePlaceholder(), ensurePlaceholder(), writeInPlace(), bootstrapSelfProtectPlaceholders(), readTempJournal(), stripTempJournal(), TestTempJournalRoundTripInPlace(), writeTempJournal() (+2 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 110 - "batchInodeRows"
Cohesion: 0.20
Nodes (13): markReady(), notifySystemdReady(), inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows() (+5 more)

### Community 111 - "mountTable"
Cohesion: 0.18
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+1 more)

### Community 112 - "dialLiveSession"
Cohesion: 0.08
Nodes (27): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditTree(), coveredByWatch(), relOrPath(), TestAuditTree() (+19 more)

### Community 113 - "recordFlow"
Cohesion: 0.16
Nodes (15): encryptBlocker(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord(), runRecordView() (+7 more)

### Community 115 - "NewDumpHook"
Cohesion: 0.15
Nodes (16): applyVerbosity(), Execute(), initLogger(), preRun(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog() (+8 more)

### Community 116 - "editauth_test.go"
Cohesion: 0.13
Nodes (24): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+16 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "GenerateConf"
Cohesion: 0.05
Nodes (59): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles() (+51 more)

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

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "IntegrationSuite"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "RefreshSection"
Cohesion: 0.10
Nodes (26): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+18 more)

### Community 131 - "LibraryClosure"
Cohesion: 0.08
Nodes (26): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+18 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "Binary"
Cohesion: 0.36
Nodes (8): lines(), pickBinaries(), pickLibs(), reviewRecording(), Classify(), Binary, Recording, TestClassifyVerifiesTheInodeThatRan()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "Config"
Cohesion: 0.06
Nodes (60): newCatalogRefresher(), makeReloadHandler(), newDaemonModel(), noteReloadedSettings(), runTUI(), startCatalogRefresh(), entryWriters(), resolveInspectors() (+52 more)

### Community 143 - "Serve"
Cohesion: 0.06
Nodes (34): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+26 more)

### Community 144 - "catalogRefresher"
Cohesion: 0.18
Nodes (6): drain(), readEvents(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), catalogRefresher, debouncer

### Community 145 - "UI"
Cohesion: 0.12
Nodes (18): channelOf(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestRunWithRefusesConcurrentRun(), InstallMode, InstallRequest, Options (+10 more)

### Community 146 - "Ledger"
Cohesion: 0.11
Nodes (5): applySchema(), Ledger, Pending, Source, fakeLedger

### Community 147 - "eventUnitTest"
Cohesion: 0.27
Nodes (4): BpfEvent, eventUnitTest, Cstr(), DecodeBpfEvent()

### Community 149 - "readGrantRequest"
Cohesion: 0.19
Nodes (12): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+4 more)

### Community 150 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 152 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.10
Nodes (30): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards(), lockVaultFully() (+22 more)

### Community 155 - "safeio.go"
Cohesion: 0.20
Nodes (16): entryOf(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), createOwned(), AtomicWriteAt(), CreateExclNoFollow(), descend(), DescendCreateOwned() (+8 more)

### Community 156 - "procstats.go"
Cohesion: 0.38
Nodes (8): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic()

### Community 157 - "runMonitor"
Cohesion: 0.11
Nodes (21): CheckBPFLSM(), CheckEBPF(), MakeDisplayPaths(), configureDaemonLogging(), newPinGeneration(), prepareDaemonStart(), resolveConfigPath(), runBPFCheck() (+13 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.27
Nodes (19): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+11 more)

### Community 161 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 165 - "feed.go"
Cohesion: 0.06
Nodes (43): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+35 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (45): pageID, procStats, reel, resizeWatch, shell, sizedDialog, snapEngines, snapNet (+37 more)

### Community 167 - "newTestTarget"
Cohesion: 0.24
Nodes (16): NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded(), TestApplyFromWhitelistStopsFirst() (+8 more)

### Community 168 - "watchSet"
Cohesion: 0.08
Nodes (27): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), addSystemBinary(), catalogPatterns(), catalogWatchPlan(), inHome(), isSymlink(), newWatchSet() (+19 more)

### Community 169 - "SanitizeText"
Cohesion: 0.09
Nodes (21): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), scriptIdentity(), formatElapsed(), orNone(), printBinary() (+13 more)

### Community 170 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 171 - "Active"
Cohesion: 0.15
Nodes (13): nextConfig(), showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), EditText(), MultiSelectKeymap(), Active() (+5 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 175 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 176 - "netinject.c"
Cohesion: 0.50
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 177 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 179 - "main_test.go"
Cohesion: 0.20
Nodes (5): infraContainerPaths(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "securewrite_test.go"
Cohesion: 0.43
Nodes (6): refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeWithin()

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

### Community 194 - "changelog.go"
Cohesion: 0.25
Nodes (17): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+9 more)

### Community 195 - "RawTarget"
Cohesion: 0.31
Nodes (9): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+1 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "github.com/spf13/cobra.Command"
Cohesion: 0.16
Nodes (19): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+11 more)

### Community 199 - "update_test.go"
Cohesion: 0.10
Nodes (20): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel() (+12 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 202 - "DaemonEvent"
Cohesion: 0.10
Nodes (21): ServeConfig, UIDResolver, NewUIDResolver(), drainEphemeral(), runDaemonUI(), runHeadless(), runServedTUI(), writeEvent() (+13 more)

### Community 203 - "Plan"
Cohesion: 0.22
Nodes (6): dedupe(), Plan, libOfChosen(), sectionFor(), Existing, ExistingRule

### Community 204 - "ValidateTargets"
Cohesion: 0.38
Nodes (6): existingSection(), guardedMatch(), overlaps(), refusedPath(), validateTarget(), ValidateTargets()

### Community 205 - "uninstall/uninstall.go"
Cohesion: 0.07
Nodes (45): TestPickEncryptedDirSingleSkipsPicker(), pickDir(), pickEncryptedDir(), runOfflineEdit(), deletePostBackups(), listBackups(), restoreBackups(), askDecrypt() (+37 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.14
Nodes (14): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+6 more)

### Community 209 - "time.Time"
Cohesion: 0.07
Nodes (23): LockedError, runNetworkMonitor(), runTUI(), sweepID, sweepMemo, Progress, formatGuardEventLine(), sanitizeTerminalText() (+15 more)

### Community 210 - "downloadAndVerify"
Cohesion: 0.13
Nodes (19): installHooks(), releaseBuild(), Run(), staged(), downloadAndVerify(), downloadFile(), downloadReleaseFiles(), FetchLatest() (+11 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "desktop_test.go"
Cohesion: 0.06
Nodes (28): TestEveryInstallModeMapped(), Progress, User, bleedLayout, logGrid, multiAnswer, pageDef, stageItem (+20 more)

### Community 213 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 214 - "root.go"
Cohesion: 0.07
Nodes (44): asRoot(), checkInstalledMatchesSource(), encryptDirectories(), fileSHA256(), runInstall(), sameContent(), secureResources(), verifyEncryptionState() (+36 more)

### Community 216 - "NetworkMonitorUseCase"
Cohesion: 0.29
Nodes (3): NetworkMonitorUseCase, NewNetworkMonitorUseCase(), StreamUseCase

### Community 217 - "promptPassword"
Cohesion: 0.40
Nodes (6): promptPassword(), promptPasswordWith(), PromptNewPassword(), promptNewPassword(), PromptNewPasswordUntilValid(), runSetPassword()

### Community 218 - "tempGrantSetup"
Cohesion: 0.23
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

### Community 219 - "Question"
Cohesion: 0.40
Nodes (3): Confirm(), Question, recordUI

### Community 220 - "shQuote"
Cohesion: 0.19
Nodes (4): IntegrationSuite, shQuote(), IntegrationSuite, IntegrationSuite

### Community 221 - "ptrace_race.c"
Cohesion: 0.52
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 222 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 223 - "globBuilder"
Cohesion: 0.10
Nodes (19): reserveChain(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, GlobChild, statFunc, compileGlobNames(), globText() (+11 more)

### Community 224 - "RunWith"
Cohesion: 0.03
Nodes (86): confirmMasterKeyOverwrite(), runGenKey(), anyDaemonRunning(), offerDaemonStart(), installGUI(), needsBinary(), RunWith(), TestNeedsBinary() (+78 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 226 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 227 - "ParseEventsFlag"
Cohesion: 0.33
Nodes (4): ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents()

### Community 228 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 229 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 231 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 244 - "monitor.md"
Cohesion: 0.50
Nodes (3): Events, Flags, monitor: see what touches a path

## Knowledge Gaps
- **33 isolated node(s):** `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 650 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **53 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Guard` connect `Guard` to `go_pkg_strings`, `TrustGuard`, `GuardInodeKey`, `BinaryEntry`, `runDaemon`, `guardUnitTest`, `DaemonEvent`, `installPage`, `os.File`, `time.Time`, `GuardEvent`, `daemonconfig.go`, `buildOneGuard`, `sync.Mutex`?**
  _High betweenness centrality (0.046) - this node is a cross-community bridge._
- **What connects `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.12308511509789465 - nodes in this community are weakly interconnected._
- **Why does `IntegrationSuite` connect `IntegrationSuite` to `.newGuardTestContainer`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.021810154525386313 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`, `.netIntegrityContainer`, `daemon_multicall_test.go`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09125085440874915 - nodes in this community are weakly interconnected._