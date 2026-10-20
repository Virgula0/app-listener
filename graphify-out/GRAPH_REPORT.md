# Graph Report - app-listener  (2026-10-20)

## Corpus Check
- 511 files · ~1,726,876 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6664 nodes · 29277 edges · 256 communities (208 shown, 48 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2710 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7bb61bc9`
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
- BinaryEntry
- FileTarget
- protectedPage
- github.com/charmbracelet/bubbles/textarea.Model
- .renderRow
- go_pkg_path_filepath
- absPath
- guard_trust.bpf.c
- runView
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- DaemonUseCase
- auth.go
- go_pkg_strings
- guard.bpf.c
- TrustGuard
- find_test.go
- codeedit.go
- runDaemon
- go_pkg_os
- TestMonitorUseCaseLifecycle
- UI
- runNetworkMonitor
- picker
- runBackend
- IntegrationSuite
- Message
- controlServer
- fakeGuardRepo
- NetEventType
- IntegrationSuite
- IntegrationSuite
- monitorUnitTest
- root.go
- fileedit_model_test.go
- CandidateDir
- GuardInodeKey
- fadeScroll
- install/install.go
- cursor.go
- Client
- daemonconfig.go
- fyne.io/fyne/v2/widget.Label
- Plan
- EventType
- mutedText
- engine
- StatInode
- recoverFileInPlace
- Recorder
- exe_supersede.h
- updatePage
- rcUser
- Config
- context.Context
- IntegrationSuite
- fileedit_model.go
- update_test.go
- .showSized
- WithSuspectsFrom
- Guard
- fcntl
- bpfstats/main.go
- policyEditor
- codecheck/main.go
- journal
- watchSet
- process_vm_readv.c
- GuardEvent
- LibraryClosure
- SelectOrphans
- githubRelease
- net_tester/main.go
- update/update.go
- recordFlow
- serve.go
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- OpenConfined
- Monitor
- prctl_exe_swap.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- ConfinedEntry
- installPage
- Well-known bypass classes
- hashCache
- jit_provenance.c
- inspector_probe.c
- recordops.c
- readGrantRequest
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- Resource
- mountTable
- github.com/charmbracelet/bubbletea.Cmd
- NetGuardEvent
- IntegrationSuite
- sync.Mutex
- daemon_multicall_test.go
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
- newTestTarget
- fetchConfig
- trace-app-libs.sh
- IntegrationSuite
- os.File
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- IntegrationSuite
- serveWebSocket
- unit.go
- appTheme
- Ledger
- buildTrustedSet
- .step
- headless/headless.go
- NetworkMonitorRepository
- fakeStore
- inodeChain
- IntegrationSuite
- buildOneGuard
- .FilterExistingWhitelist
- fetchReleases
- catalogRefresher
- Installation
- backend_test.go
- graphify-refresh.sh
- .Layout
- BinaryStat
- .populateWatchInodes
- feed.go
- shell
- editView
- gui/session.go
- WithBottomBar
- New
- binledger.go
- Limitations
- Troubleshooting
- IntegrationSuite
- runMonitor
- netinject.c
- .decryptFileWithProgress
- inode_dev.h
- Binary
- Compatibility
- How it works
- github.com/spf13/cobra.Command
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- editControlSession
- NewEventFanout
- edit-protected
- monitor: see what touches a path
- kernel-prof.sh
- downloadAndVerify
- pathCache
- confirmUpdate
- DaemonEvent
- batchInodeRows
- RefreshSection
- uninstallAll
- Desktop GUI
- IntegrationSuite
- auth_test.go
- launchKey
- .eventTypeAllowed
- bleedLayout
- time.Time
- lockAndDeprovision
- .awaitDaemonUpOrDead
- TestCleanRel
- KernelDev
- tempGrantSetup
- .markProcess
- newTestVetter
- ptrace_race.c
- prctl_gate
- controlManager
- formatEntry
- supersede_probe.c
- wm.go
- changelog_test.go
- github.com/stretchr/testify/suite.Suite
- copyRegularAt
- github.com/Virgula0/app-listener/tools/codecheck
- mustSubkey
- IntegrationSuite
- taint_victim.c
- fileedit_symlink_test.go
- securewrite_test.go
- setpassword_authz_test.go
- wm_test.go
- configedit_test.go
- Daemon configuration
- TestMain
- IntegrationSuite
- .exeKeys
- ValidateTargets
- promptPassword
- vouchedDevs
- .tokenize
- parseResize
- fakeLink
- Existing
- stripImmutableFlags
- .TestDaemon_Btrfs_CopyIoctlGate
- .TestDaemon_Bypass_PrctlExeFileSwap
- IsBinaryContent

## God Nodes (most connected - your core abstractions)
1. `Guard` - 118 edges
2. `Config` - 113 edges
3. `GuardInodeKey` - 100 edges
4. `shell` - 90 edges
5. `Load()` - 90 edges
6. `absPath()` - 88 edges
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

## Communities (256 total, 48 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.18
Nodes (28): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion, multicallCacheKey (+20 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (153): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), awaitStartupOrSignal(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestAwaitStartupOrSignalAbortsAndLocksDown(), TestAwaitStartupOrSignalAbortWithFailedStartup() (+145 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "Vault"
Cohesion: 0.13
Nodes (20): encryptDirectories(), secureResources(), verifyEncryptionState(), confirmRunPrereq(), resolveFilesystemPrereqs(), isFileVaultCiphertext(), isRegularFileTarget(), Vault (+12 more)

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
Nodes (24): addErrKind, ReplacementCheck, TestWithChmodDropWriteNeedsReadOnly(), canonicalPaths(), classifyAddErr(), guardModeKey(), modeString(), NewGuard() (+16 more)

### Community 9 - "mkdirs"
Cohesion: 0.09
Nodes (42): TestOpenBunTmp_SymlinkReplacedTargetUntouched(), newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds() (+34 more)

### Community 10 - "BinaryEntry"
Cohesion: 0.06
Nodes (25): keptEntries, BinaryEntry, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), checkBtrfsKeys(), canonicalBinaryPath(), confinedEntry() (+17 more)

### Community 11 - "FileTarget"
Cohesion: 0.11
Nodes (14): specFromForm(), guardPill(), specSummary(), EventTypes(), maskFor(), Files, FileSpec, FileTarget (+6 more)

### Community 12 - "protectedPage"
Cohesion: 0.06
Nodes (18): backendLink, configView, forwardForm, forwardView, protectedPage, recordForm, recordReview, recordView (+10 more)

### Community 13 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.09
Nodes (23): findButton, pos, span, viewportState, column(), dedent(), lines(), moveTo() (+15 more)

### Community 14 - ".renderRow"
Cohesion: 0.24
Nodes (5): lineStyles, fg(), lineLen(), markStyle(), stylesOf()

### Community 16 - "absPath"
Cohesion: 0.08
Nodes (8): IntegrationSuite, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, absPath()

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "runView"
Cohesion: 0.06
Nodes (13): daemonPage, daemonState, logGrid, logHook, logLine, runHost, runView, journalLine() (+5 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (47): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+39 more)

### Community 22 - "DaemonUseCase"
Cohesion: 0.07
Nodes (11): GuardRepository, Vault, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles(), uniqueEncryptionRoots() (+3 more)

### Community 23 - "auth.go"
Cohesion: 0.19
Nodes (17): CheckPassword(), Hash(), importLegacyHash(), openLedger(), OriginOf(), parseHash(), StoreHash(), TestHashVerifyRoundTrip() (+9 more)

### Community 24 - "go_pkg_strings"
Cohesion: 0.11
Nodes (4): attrOpCase, guardExploitTest, i3Result, i3Workspace

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (73): exe_refused(), inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade() (+65 more)

### Community 26 - "TrustGuard"
Cohesion: 0.05
Nodes (28): lsmHook, statFunc, tempGrant, trustHook, trustRows, ensureBtrfsLayout(), dropSupersededLocked(), PruneSuperseded() (+20 more)

### Community 27 - "find_test.go"
Cohesion: 0.16
Nodes (22): findAll(), lowerRunes(), NewFinder(), runesEqual(), findEditor(), keys(), TestFindAllSmartCase(), TestFinderLeavesSaveToHost() (+14 more)

### Community 28 - "codeedit.go"
Cohesion: 0.11
Nodes (37): class, classOf(), Edit(), leadingIndent(), New(), onlySpaces(), SetText(), softTabBack() (+29 more)

### Community 29 - "runDaemon"
Cohesion: 0.05
Nodes (72): CheckBPFLSM(), CheckEBPF(), backingDeviceUnion(), buildGuards(), catchLifecycleSignals(), configureDaemonLogging(), loadDaemonConfig(), lockOneRoot() (+64 more)

### Community 30 - "go_pkg_os"
Cohesion: 0.07
Nodes (37): TestPickEncryptedDirSingleSkipsPicker(), pickDir(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep() (+29 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 32 - "UI"
Cohesion: 0.09
Nodes (21): showChangelog(), InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks, InstallHooks (+13 more)

### Community 33 - "runNetworkMonitor"
Cohesion: 0.05
Nodes (42): runNetworkMonitor(), runTUI(), Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace() (+34 more)

### Community 34 - "picker"
Cohesion: 0.07
Nodes (32): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, lstatOf(), outsideHardlinks() (+24 more)

### Community 35 - "runBackend"
Cohesion: 0.15
Nodes (9): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), changePasswordLive(), Reader (+1 more)

### Community 37 - "Message"
Cohesion: 0.12
Nodes (9): File, backend, backendOp, backendState, BinaryCheck, Message, PwState, Request (+1 more)

### Community 38 - "controlServer"
Cohesion: 0.28
Nodes (3): startControlServer(), TestRecordFailureLockout(), controlServer

### Community 39 - "fakeGuardRepo"
Cohesion: 0.06
Nodes (22): confirmUserPlaced(), forwardCheck(), SupersededBinary, binaryStatOf(), TempBinary, TemporaryGrant, TempRule, planTempBlock() (+14 more)

### Community 40 - "NetEventType"
Cohesion: 0.07
Nodes (18): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), BinariesSummary(), TestBinariesSummary(), FormatAddr() (+10 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "IntegrationSuite"
Cohesion: 0.10
Nodes (4): shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite

### Community 44 - "root.go"
Cohesion: 0.08
Nodes (21): installHooks(), releaseBuild(), Run(), staged(), TestEveryInstallModeMapped(), TestMovedFlagError(), closeDumpFile(), Execute() (+13 more)

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.21
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "CandidateDir"
Cohesion: 0.06
Nodes (42): reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), appendSectionsAndEdit(), addManualDirectories(), editConfig(), groupCandidates() (+34 more)

### Community 48 - "GuardInodeKey"
Cohesion: 0.06
Nodes (27): exeRow, refusedReplacements, tempMaskOp, tempRow, eventMask(), supersedeUnnamed(), appletTag(), buildOf() (+19 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (27): Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+19 more)

### Community 50 - "install/install.go"
Cohesion: 0.05
Nodes (75): confirmMasterKeyOverwrite(), runGenKey(), offerDaemonStart(), applyLiveRefresh(), asRoot(), buildBinaryIfNeeded(), certifyInstall(), checkInstalledMatchesSource() (+67 more)

### Community 51 - "cursor.go"
Cohesion: 0.08
Nodes (22): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+14 more)

### Community 52 - "Client"
Cohesion: 0.09
Nodes (14): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+6 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.10
Nodes (50): systemPatterns(), systemRule(), grantNewBinaries(), watchGroup, Line(), Lines(), TestLinesKeyOnLinkAndIncludePending(), addResource() (+42 more)

### Community 54 - "fyne.io/fyne/v2/widget.Label"
Cohesion: 0.03
Nodes (23): binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panel, panelRenderer (+15 more)

### Community 55 - "Plan"
Cohesion: 0.12
Nodes (13): dedupe(), Plan, libOfChosen(), TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan() (+5 more)

### Community 56 - "EventType"
Cohesion: 0.08
Nodes (12): snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, Cstr(), DecodeBpfEvent(), MonitorRepository (+4 more)

### Community 57 - "mutedText"
Cohesion: 0.08
Nodes (38): checkGrid, filesPage, handCheck, homePage, newLogo(), newBinaryRow(), newCheck(), newFadeVScroll() (+30 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (29): engine, resInfo, taintOwner, TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded(), pathWithin() (+21 more)

### Community 59 - "StatInode"
Cohesion: 0.09
Nodes (27): Multicall, MulticallKind, ExeKey(), IsMulticallRefusal(), multicallFamily(), registerBuild(), singleBinaryPaths(), TestExeKeyRefusesGNUCoreutilsSingle() (+19 more)

### Community 60 - "recoverFileInPlace"
Cohesion: 0.20
Nodes (18): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+10 more)

### Community 61 - "Recorder"
Cohesion: 0.16
Nodes (15): drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start(), record() (+7 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.19
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+7 more)

### Community 63 - "updatePage"
Cohesion: 0.23
Nodes (7): installedBinary, latestRelease, updatePage, readInstalled(), newUpdatePage(), boolIndex(), isPreRelease()

### Community 64 - "rcUser"
Cohesion: 0.12
Nodes (41): setupBunTmpdir(), rcTarget, TestHasRCBlocksTrackEnsureAndRemove(), TestRemoveSSHAgentShellEnvKeepsSSHConfig(), applyRC(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv() (+33 more)

### Community 65 - "Config"
Cohesion: 0.05
Nodes (78): newCatalogRefresher(), makeReloadHandler(), noteReloadedSettings(), startCatalogRefresh(), entryWriters(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), startTrustGuard() (+70 more)

### Community 66 - "context.Context"
Cohesion: 0.19
Nodes (9): centered(), FitSize(), poll(), TestFitSize(), hyprClient, hyprMonitor, Rect, Size (+1 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "fileedit_model.go"
Cohesion: 0.10
Nodes (19): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource() (+11 more)

### Community 69 - "update_test.go"
Cohesion: 0.11
Nodes (17): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFilterChannel(), testKeyPair(), TestNewerThanStable(), TestParseChecksum() (+9 more)

### Community 70 - ".showSized"
Cohesion: 0.10
Nodes (9): binaryList, installUI, resizeWatch, sizedDialog, scrollable(), newBinaryList(), newFadeScroll(), ask() (+1 more)

### Community 71 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 72 - "Guard"
Cohesion: 0.07
Nodes (15): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), deferredBinary, rootHandle, openRoot(), SharedPinDegraded() (+7 more)

### Community 73 - "fcntl"
Cohesion: 0.09
Nodes (11): main(), report(), denied(), main(), put(), wait_then_read(), dump_via_debugfs(), main() (+3 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 75 - "policyEditor"
Cohesion: 0.13
Nodes (11): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), confirmWarnings() (+3 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "journal"
Cohesion: 0.15
Nodes (6): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 78 - "watchSet"
Cohesion: 0.15
Nodes (13): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchDir, watchPattern (+5 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.19
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.06
Nodes (31): displayEntries(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI(), bpfGuardEvent (+23 more)

### Community 81 - "LibraryClosure"
Cohesion: 0.06
Nodes (36): hardLinked(), mayUpdate(), TestIsGeneralTool(), TestIsGeneralToolJudgesTheInode(), ldConfParser, defaultLibDirs(), elfInterp(), fileExists() (+28 more)

### Community 82 - "SelectOrphans"
Cohesion: 0.23
Nodes (16): appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors(), SelectOrphans() (+8 more)

### Community 83 - "githubRelease"
Cohesion: 0.23
Nodes (18): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+10 more)

### Community 84 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 85 - "update/update.go"
Cohesion: 0.14
Nodes (17): applyUpdate(), channelOf(), confirmUnit(), deploy(), deployGUI(), InstalledGUIVersion(), InstalledVersion(), Latest() (+9 more)

### Community 86 - "recordFlow"
Cohesion: 0.13
Nodes (17): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+9 more)

### Community 87 - "serve.go"
Cohesion: 0.15
Nodes (18): basicAuth(), isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), securityHeaders(), Serve() (+10 more)

### Community 88 - "liveSession"
Cohesion: 0.12
Nodes (18): checkSecret(), dialLiveSession(), parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), putConfig(), chooseResources(), eventsLabel() (+10 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "OpenConfined"
Cohesion: 0.15
Nodes (19): loaderName(), LoaderRefusal(), soNames(), MulticallRefusal(), OpenConfined(), appletLine(), classifyGroup(), classifyLibs() (+11 more)

### Community 92 - "Monitor"
Cohesion: 0.17
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.18
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.18
Nodes (4): foreignPinnedLink(), PumpRingbuf(), hook, NetGuard

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
Cohesion: 0.12
Nodes (31): binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats() (+23 more)

### Community 101 - "installPage"
Cohesion: 0.08
Nodes (6): actionGate, binSource, installAction, installPage, uninstallCard, daemonStopNote()

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "hashCache"
Cohesion: 0.21
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.20
Nodes (14): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+6 more)

### Community 106 - "recordops.c"
Cohesion: 0.45
Nodes (13): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+5 more)

### Community 107 - "readGrantRequest"
Cohesion: 0.18
Nodes (13): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+5 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "Resource"
Cohesion: 0.06
Nodes (26): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), logRefreshChange(), TestLogRefreshChange(), libBinaries(), planUpdaters() (+18 more)

### Community 111 - "mountTable"
Cohesion: 0.27
Nodes (7): MountInfo, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), mountTable

### Community 112 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.07
Nodes (15): recordTickCmd(), recordModel, diffModel, editorModel, newFindInput(), NewDaemonModel(), TestNoticeViewFitsAndPlacesOKLast(), assertViewFits() (+7 more)

### Community 113 - "NetGuardEvent"
Cohesion: 0.08
Nodes (22): ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode() (+14 more)

### Community 115 - "sync.Mutex"
Cohesion: 0.18
Nodes (7): fyneLogFilter, newActivity(), TestRingKeepsNewest(), Counts, newRing(), activity, ring

### Community 116 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "GenerateConf"
Cohesion: 0.04
Nodes (67): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), groupedDiscordConf() (+59 more)

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

### Community 130 - "newTestTarget"
Cohesion: 0.22
Nodes (17): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 131 - "fetchConfig"
Cohesion: 0.40
Nodes (6): applyEdit(), editAgain(), fetchConfig(), nextConfig(), runEditConfig(), EditText()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "os.File"
Cohesion: 0.07
Nodes (30): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+22 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 143 - "serveWebSocket"
Cohesion: 0.23
Nodes (9): configureWebSocket(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestSameOrigin(), TestSnapshotHubSingleViewerAndCoalescing(), writeSnapshot(), writeSnapshotLoop() (+1 more)

### Community 144 - "unit.go"
Cohesion: 0.15
Nodes (20): installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), TestDaemonUnitWithMetadataOutput(), Mode (+12 more)

### Community 146 - "Ledger"
Cohesion: 0.10
Nodes (11): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), Ledger, Pending (+3 more)

### Community 147 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 148 - ".step"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 149 - "headless/headless.go"
Cohesion: 0.29
Nodes (8): stagedBinary, stageRelease(), Check(), CheckGUI(), inspect(), buildTiny(), TestCheck(), TestCheckGUI()

### Community 150 - "NetworkMonitorRepository"
Cohesion: 0.14
Nodes (7): DrainUntilSignal(), runHeadless(), ProtocolString(), NetworkMonitorRepository, NetworkMonitorUseCase, NewNetworkMonitorUseCase(), StreamUseCase

### Community 152 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.13
Nodes (24): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards(), selfGuards (+16 more)

### Community 155 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 156 - "fetchReleases"
Cohesion: 0.18
Nodes (10): downloadFile(), downloadReleaseFiles(), fetchRelease(), fetchReleases(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), TestFetchReleases() (+2 more)

### Community 157 - "catalogRefresher"
Cohesion: 0.10
Nodes (16): drain(), readEvents(), sameConfig(), swapFile(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails() (+8 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.27
Nodes (19): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+11 more)

### Community 161 - ".Layout"
Cohesion: 0.44
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 163 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 165 - "feed.go"
Cohesion: 0.06
Nodes (42): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+34 more)

### Community 166 - "shell"
Cohesion: 0.06
Nodes (44): anyDaemonRunning(), pageID, reel, shell, snapEngines, snapNet, sysStatus, newShell() (+36 more)

### Community 167 - "editView"
Cohesion: 0.11
Nodes (9): codeEntry, editView, handButton, Entry, selectable(), fileSize(), joinRel(), newCodeEntry() (+1 more)

### Community 168 - "gui/session.go"
Cohesion: 0.05
Nodes (59): entryOf(), checkAndApply(), refreshCatalogConfig(), installFileAs(), installUnit(), readAtNoFollow(), upsertFile(), upsertUnitAt() (+51 more)

### Community 169 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 170 - "New"
Cohesion: 0.19
Nodes (17): harnessSuite, looksLikeFileVaultRecord(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceWrongKeyFails(), TestUnlockLockFileInPlacePreservesInode(), TestUnlockLockFileInPlaceRoundTrip() (+9 more)

### Community 171 - "binledger.go"
Cohesion: 0.14
Nodes (15): applySchema(), EnsurePlaceholders(), JournalPath(), Open(), TestEnsurePlaceholdersNoInstallDir(), TestLedgerKeepsJournalFile(), TestLedgerRecordIfAbsentKeepsExistingVerdict(), TestLedgerRecordLookupPending() (+7 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 175 - "runMonitor"
Cohesion: 0.18
Nodes (15): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+7 more)

### Community 176 - "netinject.c"
Cohesion: 0.50
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 177 - ".decryptFileWithProgress"
Cohesion: 0.24
Nodes (5): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata()

### Community 178 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 179 - "Binary"
Cohesion: 0.20
Nodes (15): lines(), pickBinaries(), pickLibs(), reviewRecording(), beyondEvents(), Classify(), finishBinary(), Binary (+7 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "github.com/spf13/cobra.Command"
Cohesion: 0.10
Nodes (32): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+24 more)

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

### Community 194 - "editControlSession"
Cohesion: 0.30
Nodes (5): readClientLines(), streamEvents(), TestEditControlSessionEndIdempotent(), editControlSession, grantRequest

### Community 195 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "downloadAndVerify"
Cohesion: 0.17
Nodes (20): archAssets(), assetsFor(), checkUpdatePreconditions(), downloadAndVerify(), fetchGUI(), FetchLatest(), fetchPair(), FetchVerified() (+12 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "confirmUpdate"
Cohesion: 0.29
Nodes (7): confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), cancelUI, TestEditTextMapsUICancel(), Use()

### Community 202 - "DaemonEvent"
Cohesion: 0.11
Nodes (30): ServeConfig, UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+22 more)

### Community 203 - "batchInodeRows"
Cohesion: 0.20
Nodes (13): markReady(), notifySystemdReady(), inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows() (+5 more)

### Community 204 - "RefreshSection"
Cohesion: 0.15
Nodes (18): RefreshOptions, SectionScan, admitNew(), SectionChange, keepExisting(), LiveEmptyWhitelistRejected(), missingFrom(), RefreshCatalog() (+10 more)

### Community 205 - "uninstallAll"
Cohesion: 0.13
Nodes (25): deletePostBackups(), listBackups(), restoreBackups(), confirmKeyDeletion(), offerBackupCleanup(), run(), runUninstall(), RunWith() (+17 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.14
Nodes (14): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+6 more)

### Community 209 - "auth_test.go"
Cohesion: 0.36
Nodes (15): HashConfigured(), LoadHash(), MigrateLegacyHash(), mustHash(), TestCheckPasswordSharesGuessBudget(), TestMigrateLegacyHashEmptyPlaceholder(), TestMigrateLegacyHashImportFailureRemoves(), TestMigrateLegacyHashImports() (+7 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 213 - "time.Time"
Cohesion: 0.21
Nodes (7): LockedError, newWatchSet(), fakeInotify, sweepID, sweepMemo, Progress, Seen

### Community 214 - "lockAndDeprovision"
Cohesion: 0.18
Nodes (8): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), applyRawKeyPolicy(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision(), modifiedContextWithSource()

### Community 217 - "KernelDev"
Cohesion: 0.16
Nodes (17): fdStat, physicalParent(), statFD(), KernelDev(), BtrfsDevices(), btrfsSysfsDevices(), TestBtrfsSysfsDevices(), codeVouched() (+9 more)

### Community 218 - "tempGrantSetup"
Cohesion: 0.23
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

### Community 219 - ".markProcess"
Cohesion: 0.18
Nodes (7): reasonText(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), TestReasonText(), procRec

### Community 220 - "newTestVetter"
Cohesion: 0.29
Nodes (11): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+3 more)

### Community 221 - "ptrace_race.c"
Cohesion: 0.52
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 222 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 223 - "controlManager"
Cohesion: 0.24
Nodes (4): newControlManager(), startControlManager(), configEditor, controlManager

### Community 224 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 225 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 226 - "wm.go"
Cohesion: 0.29
Nodes (8): detectWM(), reachable(), Detect(), Session, hyprlandSocket(), i3RuntimeSocket(), socketOwnedBy(), Kind

### Community 227 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 228 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.18
Nodes (4): mountinfoUnitTest, unescapeMountinfo(), harnessSuite, editorHarness

### Community 229 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 231 - "mustSubkey"
Cohesion: 0.31
Nodes (9): newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape(), TestOpenFileVaultRejectsTampering(), TestSealFileVaultNoncesNeverRepeat() (+1 more)

### Community 234 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 235 - "securewrite_test.go"
Cohesion: 0.31
Nodes (8): atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync(), writeWithin()

### Community 236 - "setpassword_authz_test.go"
Cohesion: 0.44
Nodes (8): authorizeFirstTimeSetPassword(), fixed(), setDaemonRunning(), TestAuthorizeFirstTimeSetPasswordRefusedWhileDaemonRunning(), TestSetPasswordChangesInstallOrigin(), TestSetPasswordFirstTimeRefusedWhileRunning(), TestSetPasswordOfflineChangeNeedsCurrent(), TestSetPasswordWeakNewRefused()

### Community 237 - "wm_test.go"
Cohesion: 0.67
Nodes (8): envOf(), serve(), shortDir(), TestDetect(), TestHyprlandFloat(), TestHyprlandFloatTimesOutWithoutWindow(), TestI3Float(), uid()

### Community 238 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 239 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 240 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 243 - "ValidateTargets"
Cohesion: 0.38
Nodes (6): existingSection(), guardedMatch(), overlaps(), refusedPath(), validateTarget(), ValidateTargets()

### Community 244 - "promptPassword"
Cohesion: 0.40
Nodes (6): promptPassword(), promptPasswordWith(), PromptNewPassword(), promptNewPassword(), PromptNewPasswordUntilValid(), runSetPassword()

### Community 245 - "vouchedDevs"
Cohesion: 0.33
Nodes (6): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs()

### Community 246 - ".tokenize"
Cohesion: 0.40
Nodes (3): TestDaemonConfTokens(), TestTokenizeKeepsLineCount(), tokenClassOf()

### Community 247 - "parseResize"
Cohesion: 0.40
Nodes (3): decodeEndsAtEOF(), parseResize(), TestParseResize()

### Community 250 - "stripImmutableFlags"
Cohesion: 0.50
Nodes (4): stripEntryFlags(), stripImmutableFlags(), TestStripImmutableFlagsDoesNotHangOnFIFO(), TestStripImmutableFlagsSkipsDanglingSymlink()

## Knowledge Gaps
- **33 isolated node(s):** `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 626 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **48 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `TestMain`, `daemon_multicall_test.go`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`, `.TestDaemon_Btrfs_CopyIoctlGate`, `.TestDaemon_Bypass_PrctlExeFileSwap`, `go_pkg_strings`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **What connects `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.017737644118358462 - nodes in this community are weakly interconnected._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `TrustGuard`, `engine`, `time.Duration`, `os.File`, `ConfinedEntry`, `guardUnitTest`, `DaemonEvent`, `BinaryEntry`, `Resource`, `GuardInodeKey`, `GuardEvent`, `sync.Mutex`, `time.Time`, `daemonconfig.go`, `buildOneGuard`, `runDaemon`, `controlManager`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09125085440874915 - nodes in this community are weakly interconnected._
- **Why does `shell` connect `shell` to `UI`, `picker`, `time.Duration`, `wm.go`, `feed.go`, `.showSized`, `installPage`, `BinaryEntry`, `policyEditor`, `protectedPage`, `FileTarget`, `root.go`, `runView`, `cursor.go`, `fyne.io/fyne/v2/widget.Label`, `mutedText`, `updatePage`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `Vault` be split into smaller, more focused modules?**
  _Cohesion score 0.12903225806451613 - nodes in this community are weakly interconnected._