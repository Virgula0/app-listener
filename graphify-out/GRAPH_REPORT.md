# Graph Report - app-listener  (2026-10-29)

## Corpus Check
- 534 files · ~1,754,051 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6939 nodes · 30800 edges · 238 communities (201 shown, 37 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2863 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `721bb35a`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- Vault
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- tempGrantSetup
- New
- protectedPage
- os.File
- daemonconfig.go
- root.go
- IntegrationSuite
- guard_trust.bpf.c
- navItem
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- trustManager
- guard.bpf.c
- TrustGuard
- find_test.go
- sync.Mutex
- golang.org/x/sys/unix.Stat_t
- github.com/charmbracelet/bubbletea.Cmd
- TestMonitorUseCaseLifecycle
- policyEditor
- Guard
- picker
- SelectOrphans
- IntegrationSuite
- Message
- controlServer
- .posAt
- go_pkg_os
- IntegrationSuite
- absPath
- monitorUnitTest
- BinaryEntry
- fileedit_model_test.go
- EventType
- stdio
- update/update.go
- fadeScroll
- vaultfs_test.go
- editView
- Client
- Config
- hashCache
- applier
- github.com/spf13/cobra.Command
- mutedText
- GuardInodeKey
- guardModel
- audit.go
- Recorder
- exe_supersede.h
- installPage
- UI
- unix_connect.c
- wm_test.go
- IntegrationSuite
- User
- liveSession
- DaemonUseCase
- RefreshSection
- pickfs.go
- fcntl
- bpfstats/main.go
- CandidateDir
- codecheck/main.go
- fakeGuardRepo
- update_test.go
- process_vm_readv.c
- GuardEvent
- op
- KernelDev
- rcUser
- Monitor
- newTestTarget
- os.FileMode
- GenerateConf
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- runView
- Candidate
- netinject.c
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
- openFileVaultWithMasterKey
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- Finder
- mountTable
- New
- recordFlow
- IntegrationSuite
- recoverFileInPlace
- newTestVetter
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- .TestRecord_GuardParity
- .apply
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- .filterEvent
- make test-integration (rootful Docker suite)
- IntegrationSuite
- backend.go
- updatePage
- trace-app-libs.sh
- IntegrationSuite
- maintainerKeeper
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- catalogRefresher
- serve.go
- NetworkMonitor
- gui/session.go
- Ledger
- runOfflineEdit
- go_pkg_strings
- githubRelease
- readGrantRequest
- fakeStore
- github.com/charmbracelet/bubbles/textarea.Model
- IntegrationSuite
- buildOneGuard
- runNetworkGuard
- planUpdaters
- NetworkMonitorRepository
- Installation
- backend_test.go
- graphify-refresh.sh
- layout.go
- gui.go
- IntegrationSuite
- time.Time
- feed.go
- shell
- binaryList
- watchSet
- WithBottomBar
- Vault
- downloadFile
- Limitations
- net_tester/main.go
- IntegrationSuite
- compileGlobNames
- lockAndDeprovision
- NewEventFanout
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
- dialog_test.go
- edit-protected
- changelog_test.go
- kernel-prof.sh
- CopyTreeWithProgress
- pathCache
- .FilterExistingWhitelist
- DaemonEvent
- Daemon configuration
- fileedit_symlink_test.go
- NetGuardEvent
- Desktop GUI
- IntegrationSuite
- WithSuspectsFrom
- mustSubkey
- .eventTypeAllowed
- configedit_test.go
- inode_dev.h
- procstats.go
- .awaitDaemonUpOrDead
- bpfLsmListed
- VouchedDevs
- io.Writer
- resolveFilesystemPrereqs
- shQuote
- .exeKeys
- CleanupStalePins
- launchKey
- newBinary
- supersede_probe.c
- checkRunningBinaryMatchesInstalled
- inspector_probe.c
- Choice
- github.com/Virgula0/app-listener/tools/codecheck
- Vault
- Plan
- taint_victim.c
- Troubleshooting
- appTheme
- bleedLayout
- TestCleanRel
- monitor.md

## God Nodes (most connected - your core abstractions)
1. `Guard` - 123 edges
2. `Config` - 118 edges
3. `GuardInodeKey` - 103 edges
4. `shell` - 97 edges
5. `absPath()` - 90 edges
6. `Load()` - 90 edges
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

## Communities (238 total, 37 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.16
Nodes (45): discoverInfra(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion (+37 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (149): plantBun(), TestOpenBunTmp_KeepsVettedDir(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestEveryInstallModeMapped() (+141 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "Vault"
Cohesion: 0.16
Nodes (15): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+7 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (80): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+72 more)

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
Cohesion: 0.06
Nodes (24): addErrKind, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+16 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "tempGrantSetup"
Cohesion: 0.21
Nodes (10): tempJournalRows(), journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure() (+2 more)

### Community 11 - "New"
Cohesion: 0.19
Nodes (19): looksLikeFileVaultRecord(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceWrongKeyFails(), TestUnlockLockFileInPlacePreservesInode(), TestUnlockLockFileInPlaceRoundTrip(), TestVaultUnlockLockDispatchForRegularFile() (+11 more)

### Community 12 - "protectedPage"
Cohesion: 0.04
Nodes (32): backendLink, checkGrid, configView, forwardForm, forwardView, installUI, protectedPage, recordReview (+24 more)

### Community 13 - "os.File"
Cohesion: 0.05
Nodes (65): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), scriptIdentity(), setConfinedHomes(), startTrustGuard(), trustStartupError() (+57 more)

### Community 14 - "daemonconfig.go"
Cohesion: 0.35
Nodes (6): recordTick, rejectSymlinkedComponents(), validateWatchTarget(), isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 15 - "root.go"
Cohesion: 0.10
Nodes (15): releaseBuild(), Run(), closeDumpFile(), Execute(), initLogger(), postRun(), InstallMode, InstallRequest (+7 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (6): eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (60): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+52 more)

### Community 18 - "navItem"
Cohesion: 0.03
Nodes (18): binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pill (+10 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (11): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (48): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+40 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (7): fakeDaemon, Resource, GuardRepository, encryptionRootOf(), uniqueEncryptionRoots(), vaultOpForGuard(), vaultOpOnRoot()

### Community 23 - "auth_test.go"
Cohesion: 0.12
Nodes (41): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+33 more)

### Community 24 - "trustManager"
Cohesion: 0.09
Nodes (30): fileID(), openBunTmp(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), makeReloadHandler(), noteReloadedSettings(), startCatalogRefresh(), inspectorPaths() (+22 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.08
Nodes (79): exe_refused(), inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade() (+71 more)

### Community 26 - "TrustGuard"
Cohesion: 0.06
Nodes (24): lsmHook, trustHook, trustRows, deleteResKeys(), guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), TestUnixSocketHooksAreRequired(), dropSupersededLocked() (+16 more)

### Community 27 - "find_test.go"
Cohesion: 0.21
Nodes (18): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+10 more)

### Community 28 - "sync.Mutex"
Cohesion: 0.04
Nodes (31): refusedReplacements, TestRefusedReplacements_ReportsEachInodeOnce(), specFromForm(), guardPill(), netGuardPill(), NetworkGuardRepository, newActivity(), Files (+23 more)

### Community 29 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.12
Nodes (22): placedWalk, treeJudge, confineBelow(), inUserHome(), vouchedFd(), GlobSystemPlaced(), placedDir(), placedEntry() (+14 more)

### Community 30 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.08
Nodes (19): recordTickCmd(), recordModel, listenForDaemonEvents(), tickDaemonStats(), listenForGuardEvents(), tickGuardStats(), syncViewport(), listenForNetEvents() (+11 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 32 - "policyEditor"
Cohesion: 0.07
Nodes (22): handCheckGroup, handRadioGroup, handSelect, netGuardTab, policyEditor, policyHooks, policyTab, handCursor() (+14 more)

### Community 33 - "Guard"
Cohesion: 0.06
Nodes (20): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), exeRow, rootHandle, supersededBinary, tempGrant (+12 more)

### Community 34 - "picker"
Cohesion: 0.15
Nodes (5): keyEntry, picker, errText(), isFindShortcut(), newKeyEntry()

### Community 35 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (13): checkSecret(), checkTrustLine(), trustRows(), TestLiveSessionWatchAndPinger(), File, TrustItem, backend, backendOp (+5 more)

### Community 38 - "controlServer"
Cohesion: 0.17
Nodes (7): readClientLines(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), controlServer, editControlSession, grantRequest

### Community 39 - ".posAt"
Cohesion: 0.11
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+6 more)

### Community 40 - "go_pkg_os"
Cohesion: 0.08
Nodes (21): atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync(), writeWithin() (+13 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.07
Nodes (14): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest, IntegrationSuite (+6 more)

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, IntegrationSuite, IntegrationSuite, absPath()

### Community 43 - "monitorUnitTest"
Cohesion: 0.15
Nodes (4): sameInode(), dirTarget(), newMountTable(), monitorUnitTest

### Community 44 - "BinaryEntry"
Cohesion: 0.06
Nodes (29): deferredBinary, keptEntries, snapEngines, BinaryEntry, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), TestWithChmodDropWriteNeedsReadOnly(), checkBtrfsKeys() (+21 more)

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.21
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "EventType"
Cohesion: 0.05
Nodes (24): snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), specSummary() (+16 more)

### Community 48 - "update/update.go"
Cohesion: 0.11
Nodes (38): applyUpdate(), archAssets(), assetsFor(), channelOf(), checkUpdatePreconditions(), confirmUnit(), deploy(), deployGUI() (+30 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (27): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+19 more)

### Community 50 - "vaultfs_test.go"
Cohesion: 0.13
Nodes (21): DefaultNewFileMode(), IsBinaryContent(), IsOctalMode(), ParseMode(), inodeOf(), mustRead(), mustWrite(), openRootT() (+13 more)

### Community 51 - "editView"
Cohesion: 0.11
Nodes (9): codeEntry, editView, handButton, Entry, selectable(), fileSize(), joinRel(), newCodeEntry() (+1 more)

### Community 52 - "Client"
Cohesion: 0.10
Nodes (12): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+4 more)

### Community 53 - "Config"
Cohesion: 0.06
Nodes (68): addedResources(), appendSectionsAndEdit(), collectDiffAdditions(), TestAddedResourcesSkipsAlreadyGuardedSSH(), TestUncoveredCandidates(), uncoveredCandidates(), runConfigEditor(), watchGroup (+60 more)

### Community 54 - "hashCache"
Cohesion: 0.12
Nodes (17): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, BinaryStat, DurableStat, HashCache (+9 more)

### Community 55 - "applier"
Cohesion: 0.22
Nodes (4): applier, ledgerOps, PlannedBinary, vaultOps

### Community 56 - "github.com/spf13/cobra.Command"
Cohesion: 0.10
Nodes (32): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+24 more)

### Community 57 - "mutedText"
Cohesion: 0.07
Nodes (39): filesPage, multiAnswer, panel, recordForm, newLogo(), newBinaryList(), newBinaryRow(), newCheck() (+31 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (53): engine, inodeBatchFn, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+45 more)

### Community 59 - "guardModel"
Cohesion: 0.07
Nodes (26): formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), appendCapped(), clipToWidth(), formatResourceBar(), formatType() (+18 more)

### Community 60 - "audit.go"
Cohesion: 0.09
Nodes (31): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource() (+23 more)

### Community 61 - "Recorder"
Cohesion: 0.17
Nodes (14): Target, drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start() (+6 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.21
Nodes (14): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+6 more)

### Community 63 - "installPage"
Cohesion: 0.09
Nodes (16): actionGate, binSource, handCheck, installAction, installPage, uninstallCard, UninstallHooks, UninstallRequest (+8 more)

### Community 64 - "UI"
Cohesion: 0.09
Nodes (18): installHooks(), staged(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestRunWithRefusesConcurrentRun(), UpdateHooks, cancelUI, InstallHooks (+10 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.41
Nodes (12): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), main(), read_secret() (+4 more)

### Community 66 - "wm_test.go"
Cohesion: 0.11
Nodes (24): reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket(), poll() (+16 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "User"
Cohesion: 0.06
Nodes (65): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), bunLaunchersOffered() (+57 more)

### Community 69 - "liveSession"
Cohesion: 0.17
Nodes (12): parseErr(), TestLiveSessionConfigAndPut(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), forwardCheck(), holdForward() (+4 more)

### Community 70 - "DaemonUseCase"
Cohesion: 0.12
Nodes (6): newControlManager(), startControlManager(), startControlServer(), configEditor, controlManager, DaemonUseCase

### Community 71 - "RefreshSection"
Cohesion: 0.08
Nodes (31): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+23 more)

### Community 72 - "pickfs.go"
Cohesion: 0.21
Nodes (15): listOpts, pickEntry, treeSearch, lstatOf(), outsideHardlinks(), TestSearchTreeFollowsNoSymlink(), plural(), entryIsDir() (+7 more)

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 75 - "CandidateDir"
Cohesion: 0.08
Nodes (24): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName() (+16 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "fakeGuardRepo"
Cohesion: 0.06
Nodes (21): binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps(), planTempBlock(), ResolveTempBinary(), stripPinnedAllow() (+13 more)

### Community 78 - "update_test.go"
Cohesion: 0.09
Nodes (20): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel() (+12 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.08
Nodes (18): holdHeadless(), HeadlessLine(), bpfGuardEvent, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths(), commMatchesGuardedBinary() (+10 more)

### Community 81 - "op"
Cohesion: 0.05
Nodes (52): SystemTrusted(), clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath() (+44 more)

### Community 82 - "KernelDev"
Cohesion: 0.16
Nodes (17): fdStat, physicalParent(), statFD(), KernelDev(), BtrfsDevices(), btrfsSysfsDevices(), TestBtrfsSysfsDevices(), codeVouched() (+9 more)

### Community 83 - "rcUser"
Cohesion: 0.11
Nodes (43): listFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative(), TestListFacultativesIgnoresDaemonConf(), TestListFacultativesMarksGuardedSSH(), rcTarget, BunTmpDir(), TestHasRCBlocksTrackEnsureAndRemove() (+35 more)

### Community 84 - "Monitor"
Cohesion: 0.16
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 86 - "newTestTarget"
Cohesion: 0.22
Nodes (17): newExeResolver(), NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 87 - "os.FileMode"
Cohesion: 0.08
Nodes (34): checkAndApply(), TestDaemonUnitWithMetadataOutput(), refreshCatalogConfig(), installFileAs(), installUnit(), readAtNoFollow(), upsertFile(), upsertUnitAt() (+26 more)

### Community 88 - "GenerateConf"
Cohesion: 0.05
Nodes (57): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic() (+49 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "runView"
Cohesion: 0.05
Nodes (19): daemonPage, logGrid, logHook, logLine, runHost, runView, stageItem, stageState (+11 more)

### Community 92 - "Candidate"
Cohesion: 0.17
Nodes (18): selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories(), pickFromCandidates(), pickUsers() (+10 more)

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

### Community 95 - "NetGuard"
Cohesion: 0.10
Nodes (10): reasonText(), eventTypeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), TestReasonText(), hook (+2 more)

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
Nodes (33): binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats() (+25 more)

### Community 101 - "pinstate.go"
Cohesion: 0.13
Nodes (27): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+19 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "NetEventType"
Cohesion: 0.11
Nodes (15): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), FormatAddr(), NetEventTypes(), Ntohs() (+7 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "openFileVaultWithMasterKey"
Cohesion: 0.13
Nodes (12): entryOf(), migrationTarget, openFileVaultWithMasterKey(), sealFileVaultWithMasterKey(), withFileVaultSubkey(), writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile() (+4 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 110 - "Finder"
Cohesion: 0.16
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 111 - "mountTable"
Cohesion: 0.15
Nodes (10): MountInfo, mountinfoUnitTest, TestNoteRootMountFollowsMountID(), MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+2 more)

### Community 112 - "New"
Cohesion: 0.10
Nodes (30): diffModel, editorModel, New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation() (+22 more)

### Community 113 - "recordFlow"
Cohesion: 0.19
Nodes (11): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+3 more)

### Community 115 - "recoverFileInPlace"
Cohesion: 0.31
Nodes (13): classifyRegularFileTarget(), clearRecovery(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), recoverFileInPlace(), rejectSidecarSymlink(), stageRecovery(), inode() (+5 more)

### Community 116 - "newTestVetter"
Cohesion: 0.29
Nodes (11): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+3 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - ".apply"
Cohesion: 0.18
Nodes (11): Apply(), DeleteBackup(), Outcome, ProgressFunc, TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal() (+3 more)

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

### Community 130 - "backend.go"
Cohesion: 0.05
Nodes (93): offerDaemonStart(), runDiffCatalog(), applyLiveRefresh(), asRoot(), buildBinaryIfNeeded(), installBinaryOnly(), markInstallComplete(), preflightDeployedBinary() (+85 more)

### Community 131 - "updatePage"
Cohesion: 0.23
Nodes (7): installedBinary, latestRelease, updatePage, readInstalled(), newUpdatePage(), boolIndex(), isPreRelease()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "maintainerKeeper"
Cohesion: 0.13
Nodes (16): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), startMaintainers(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), grantNewBinaries() (+8 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "catalogRefresher"
Cohesion: 0.14
Nodes (9): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), catalogRefresher (+1 more)

### Community 143 - "serve.go"
Cohesion: 0.08
Nodes (31): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+23 more)

### Community 144 - "NetworkMonitor"
Cohesion: 0.13
Nodes (8): BinariesSummary(), TestBinariesSummary(), PumpRingbuf(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 145 - "gui/session.go"
Cohesion: 0.19
Nodes (15): invokingIDs(), ancestorEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestOwnedByChecksTypeAndOwner() (+7 more)

### Community 146 - "Ledger"
Cohesion: 0.07
Nodes (28): candidates(), confirmAll(), describe(), record(), run(), applySchema(), Ledger, Pending (+20 more)

### Community 147 - "runOfflineEdit"
Cohesion: 0.12
Nodes (16): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), TestPickEncryptedDirSingleSkipsPicker(), chooseResource() (+8 more)

### Community 148 - "go_pkg_strings"
Cohesion: 0.07
Nodes (7): Progress, User, attrOpCase, guardExploitTest, fileID, i3Result, i3Workspace

### Community 149 - "githubRelease"
Cohesion: 0.29
Nodes (15): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+7 more)

### Community 150 - "readGrantRequest"
Cohesion: 0.18
Nodes (13): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+5 more)

### Community 152 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.15
Nodes (23): class, pos, classOf(), column(), dedent(), Edit(), leadingIndent(), lines() (+15 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.07
Nodes (44): backingDeviceUnion(), buildGuards(), buildOneGuard(), eventFilterOptions(), newPinGeneration(), reloadOnce(), reloadSlotsNeeded(), startGuardedDaemon() (+36 more)

### Community 155 - "runNetworkGuard"
Cohesion: 0.07
Nodes (41): DrainUntilSignal(), CheckBPFLSM(), CheckEBPF(), configureDaemonLogging(), confirmMasterKeyOverwrite(), loadDaemonConfig(), prepareDaemonStart(), resolveConfigPath() (+33 more)

### Community 156 - "planUpdaters"
Cohesion: 0.14
Nodes (18): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+10 more)

### Community 157 - "NetworkMonitorRepository"
Cohesion: 0.18
Nodes (4): NetworkMonitorRepository, NetworkMonitorUseCase, NewNetworkMonitorUseCase(), StreamUseCase

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 161 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 162 - "gui.go"
Cohesion: 0.11
Nodes (29): installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), deletePostBackups(), listBackups() (+21 more)

### Community 164 - "time.Time"
Cohesion: 0.29
Nodes (4): LockedError, sweepID, sweepMemo, dirEntry

### Community 165 - "feed.go"
Cohesion: 0.06
Nodes (41): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+33 more)

### Community 166 - "shell"
Cohesion: 0.04
Nodes (59): anyDaemonRunning(), fakeLink, homePage, Options, pageDef, pageID, reel, resizeWatch (+51 more)

### Community 167 - "binaryList"
Cohesion: 0.08
Nodes (13): binaryList, buttonCursor, handCursors(), handCursorsOn(), newButtonCursor(), showConfirm(), showCustomConfirm(), showDialog() (+5 more)

### Community 168 - "watchSet"
Cohesion: 0.11
Nodes (19): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), addSystemBinary(), catalogPatterns(), catalogWatchPlan(), inHome(), isSymlink(), newWatchSet() (+11 more)

### Community 169 - "WithBottomBar"
Cohesion: 0.31
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 170 - "Vault"
Cohesion: 0.18
Nodes (6): checkFreshDir(), openDirNoFollow(), renewDir(), Vault, Open(), Unopened()

### Community 171 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 172 - "Limitations"
Cohesion: 0.20
Nodes (10): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+2 more)

### Community 173 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 175 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 176 - "lockAndDeprovision"
Cohesion: 0.22
Nodes (6): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision()

### Community 177 - "NewEventFanout"
Cohesion: 0.24
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 178 - "runMonitor"
Cohesion: 0.13
Nodes (18): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), pathWithinDir() (+10 more)

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
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

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
Cohesion: 0.18
Nodes (16): lines(), pickBinaries(), pickLibs(), reviewRecording(), beyondEvents(), Classify(), finishBinary(), Binary (+8 more)

### Community 195 - "dialog_test.go"
Cohesion: 0.39
Nodes (11): pickMode, openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder() (+3 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "changelog_test.go"
Cohesion: 0.15
Nodes (17): nextConfig(), newChangelogModel(), showChangelog(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize() (+9 more)

### Community 199 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 202 - "DaemonEvent"
Cohesion: 0.08
Nodes (43): ServeConfig, UIDResolver, NewUIDResolver(), awaitStartupOrSignal(), catchLifecycleSignals(), drainEphemeral(), markReady(), newDaemonModel() (+35 more)

### Community 203 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 204 - "fileedit_symlink_test.go"
Cohesion: 0.25
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 205 - "NetGuardEvent"
Cohesion: 0.20
Nodes (5): snapNetGuard, NetGuardEvent, protoString(), netGuardEventLine, netGuardModel

### Community 206 - "Desktop GUI"
Cohesion: 0.13
Nodes (15): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+7 more)

### Community 209 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 210 - "mustSubkey"
Cohesion: 0.27
Nodes (10): deriveFileVaultSubkey(), newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape(), TestOpenFileVaultRejectsTampering() (+2 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 213 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 214 - "procstats.go"
Cohesion: 0.30
Nodes (9): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+1 more)

### Community 216 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 217 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 218 - "io.Writer"
Cohesion: 0.22
Nodes (6): formatElapsed(), orNone(), printBinary(), printRecording(), fyneLogFilter, limitedWriter

### Community 219 - "resolveFilesystemPrereqs"
Cohesion: 0.28
Nodes (6): encryptDirectories(), confirmRunPrereq(), resolveFilesystemPrereqs(), harnessSuite, EnsureSystemSetup(), isKernelAtLeast54()

### Community 220 - "shQuote"
Cohesion: 0.16
Nodes (6): IntegrationSuite, IntegrationSuite, overlayPlanter(), shQuote(), inNosuidNamespace(), IntegrationSuite

### Community 222 - "CleanupStalePins"
Cohesion: 0.33
Nodes (6): CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 224 - "newBinary"
Cohesion: 0.60
Nodes (4): newBinary(), opNames(), sortedKeys(), RawAccessor

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "checkRunningBinaryMatchesInstalled"
Cohesion: 0.40
Nodes (5): checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), fileSHA256(), sameContent(), TestCheckRunningBinaryMatchesInstalled()

### Community 227 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 228 - "Choice"
Cohesion: 0.50
Nodes (3): MultiSelectKeymap(), Choice, MultiSelect()

### Community 232 - "Plan"
Cohesion: 0.20
Nodes (6): dedupe(), Plan, libOfChosen(), sectionFor(), Existing, ExistingRule

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 244 - "monitor.md"
Cohesion: 0.50
Nodes (3): Events, Flags, monitor: see what touches a path

## Knowledge Gaps
- **34 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck` (+29 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 641 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **37 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `TestMain`, `go_pkg_strings`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _34 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.018196856906534328 - nodes in this community are weakly interconnected._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `time.Duration`, `GuardInodeKey`, `time.Time`, `DaemonUseCase`, `ConfinedEntry`, `guardUnitTest`, `DaemonEvent`, `BinaryEntry`, `os.File`, `GuardEvent`, `Resource`, `buildOneGuard`, `sync.Mutex`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09125085440874915 - nodes in this community are weakly interconnected._
- **Why does `backend` connect `Message` to `backend.go`, `time.Duration`, `Vault`, `liveSession`, `Vault`, `recordFlow`, `TestCleanRel`, `runOfflineEdit`, `Client`, `Resource`, `sync.Mutex`, `backend_test.go`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `daemonconfig_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06666666666666667 - nodes in this community are weakly interconnected._