# Graph Report - app-listener  (2026-10-19)

## Corpus Check
- 496 files · ~1,722,120 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6538 nodes · 28378 edges · 235 communities (194 shown, 41 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2621 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `0cbfcd0e`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- .Unlock
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- GuardEvent
- Network
- protectedPage
- Finder
- NetGuardEvent
- editView
- IntegrationSuite
- guard_trust.bpf.c
- runView
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- pinstate.go
- guard.bpf.c
- TrustGuard
- selection.go
- highlight_test.go
- rcUser
- codeedit.go
- TestMonitorUseCaseLifecycle
- UI
- NetEventType
- picker
- decryptStep
- IntegrationSuite
- Message
- controlServer
- TempBinary
- NetworkMonitor
- main_test.go
- absPath
- monitorUnitTest
- RunWith
- fileedit_model_test.go
- CandidateDir
- stdio
- classifyMigrationTarget
- fadeScroll
- .showSized
- cursor.go
- Client
- watchGroup
- tabs.go
- Plan
- FileEvent
- mutedText
- GuardInodeKey
- os.File
- wipe
- Recorder
- exe_supersede.h
- procstats.go
- catalogRefresher
- FileTarget
- wm.go
- IntegrationSuite
- dialLiveSession
- update_test.go
- os.FileMode
- .FilterExistingWhitelist
- Guard
- fcntl
- bpfstats/main.go
- RefreshSection
- journal
- time.Time
- process_vm_readv.c
- singleBinaryPaths
- deploy
- planUpdaters
- update/update.go
- net_tester/main.go
- PlanDaemonUnitUpdate
- fakeGuardRepo
- buildTrustedSet
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- Vault
- Monitor
- prctl_exe_swap.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- gui/session.go
- github.com/charmbracelet/bubbles/textarea.Model
- Well-known bypass classes
- RunCmd
- jit_provenance.c
- inspector_probe.c
- recordops.c
- sysStatus
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- runMaintenanceMode
- mountTable
- github.com/charmbracelet/bubbletea.Cmd
- EventType
- IntegrationSuite
- fetchPair
- CheckBPFLSM
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
- Candidate
- trace-app-libs.sh
- IntegrationSuite
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- serveWebSocket
- BinaryEntry
- Vault
- Open
- Ledger
- editControlSession
- New
- SelectOrphans
- runMonitor
- fakeStore
- KernelDev
- shQuote
- buildOneGuard
- prepareDaemonStart
- NewDumpHook
- inodeChain
- Installation
- backend_test.go
- graphify-refresh.sh
- .Layout
- changelog_test.go
- vaultfs_test.go
- .populateWatchInodes
- sync.Mutex
- shell
- revertSSHAgents
- mustSubkey
- Config
- io.Writer
- fileedit_symlink_test.go
- Limitations
- Troubleshooting
- IntegrationSuite
- parseGuardEvent
- netinject.c
- copyRegularAt
- inode_dev.h
- SanitizeText
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
- refuseOutsideHardlinks
- eventUnitTest
- edit-protected
- monitor: see what touches a path
- kernel-prof.sh
- run
- pathCache
- bpfLsmListed
- DaemonEvent
- parseResize
- runNetworkMonitor
- Desktop GUI
- IntegrationSuite
- User
- tempGrantSetup
- .eventTypeAllowed
- bleedLayout
- NewEventFanout
- .awaitDaemonUpOrDead
- searchDirs
- options.go
- IntegrationSuite
- TestCleanRel
- Daemon configuration
- prctl_gate
- run
- Prereq
- ptrace_race.c
- downloadFile
- DiscoverInfraBinaries
- checkRunningBinaryMatchesInstalled
- taint_victim.c
- parseKernelRelease
- golang.org/x/sys/unix.Timespec
- serve.go
- configedit_test.go
- MulticallError

## God Nodes (most connected - your core abstractions)
1. `Guard` - 118 edges
2. `Config` - 109 edges
3. `GuardInodeKey` - 99 edges
4. `Load()` - 89 edges
5. `absPath()` - 88 edges
6. `shell` - 87 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 69 edges
10. `IntegrationSuite` - 68 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `resolveAssets()` --calls--> `asset()`  [INFERRED]
  cmd/functions/update/update.go → internal/gui/icons.go

## Import Cycles
- None detected.

## Communities (235 total, 41 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.09
Nodes (18): tempJournalRow, Progress, User, invalidPassword, recordTick, linkedFile, launchRule, trustEvent (+10 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (166): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), atomicWriteAt(), descendNoFollow(), refuseExistingSymlink() (+158 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): Hash(), configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec() (+2 more)

### Community 3 - ".Unlock"
Cohesion: 0.11
Nodes (18): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), hasEncryptionPolicy(), isLockedRegularFileErr(), isNotEncryptedErr() (+10 more)

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
Nodes (21): addErrKind, ReplacementCheck, classifyAddErr(), addRecorder(), copySelf(), guardUnitTest, runTool(), walkEntries() (+13 more)

### Community 9 - "mkdirs"
Cohesion: 0.12
Nodes (37): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+29 more)

### Community 10 - "GuardEvent"
Cohesion: 0.05
Nodes (27): displayEntries(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI(), snapGuard (+19 more)

### Community 11 - "Network"
Cohesion: 0.06
Nodes (16): binaryList, netGuardTab, policyEditor, policyHooks, policyTab, verbFor(), netWarnings(), netGuardPill() (+8 more)

### Community 12 - "protectedPage"
Cohesion: 0.06
Nodes (26): backendLink, checkGrid, configView, forwardForm, forwardView, protectedPage, recordReview, recordView (+18 more)

### Community 13 - "Finder"
Cohesion: 0.12
Nodes (20): findButton, span, column(), findAll(), Finder, lowerRunes(), NewFinder(), newFindInput() (+12 more)

### Community 14 - "NetGuardEvent"
Cohesion: 0.10
Nodes (17): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), snapNetGuard, Mode (+9 more)

### Community 15 - "editView"
Cohesion: 0.12
Nodes (8): codeEntry, editView, Entry, selectable(), fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "runView"
Cohesion: 0.06
Nodes (16): logGrid, logHook, logLine, runView, stageItem, stageState, asset(), coloredIcon() (+8 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (11): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (49): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+41 more)

### Community 22 - "Resource"
Cohesion: 0.06
Nodes (14): ownerPaths(), newControlManager(), startControlManager(), fakeDaemon, Resource, GuardRepository, Vault, concurrencyLimit() (+6 more)

### Community 23 - "auth_test.go"
Cohesion: 0.13
Nodes (39): CheckPassword(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf(), parseHash() (+31 more)

### Community 24 - "pinstate.go"
Cohesion: 0.18
Nodes (22): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+14 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): exe_refused(), inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade() (+64 more)

### Community 26 - "TrustGuard"
Cohesion: 0.06
Nodes (21): exeRow, tempGrant, tempMaskOp, tempRow, trustHook, trustRows, deleteResKeys(), setSupersedeTrust() (+13 more)

### Community 27 - "selection.go"
Cohesion: 0.11
Nodes (15): pos, cellToRune(), clipboardText(), cursorPos(), Selection, offset(), setClipboard(), cellOf() (+7 more)

### Community 28 - "highlight_test.go"
Cohesion: 0.12
Nodes (28): moveTo(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines() (+20 more)

### Community 29 - "rcUser"
Cohesion: 0.14
Nodes (33): rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), EnsureBunTmpDir(), ensureRCBlock(), EnsureSSHAgentEnv() (+25 more)

### Community 30 - "codeedit.go"
Cohesion: 0.15
Nodes (17): class, classOf(), dedent(), Edit(), leadingIndent(), lines(), Navigate(), onlySpaces() (+9 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.24
Nodes (8): newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestMonitorUseCaseLifecycle(), TestNetworkGuardUseCaseLifecycle(), TestNetworkMonitorUseCaseLifecycle(), fakeNetworkGuardRepo, fakeNetworkMonitorRepo

### Community 32 - "UI"
Cohesion: 0.10
Nodes (18): showChangelog(), confirmUpdate(), isTerminal(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), TestRunWithRefusesConcurrentRun(), cancelUI (+10 more)

### Community 33 - "NetEventType"
Cohesion: 0.15
Nodes (14): NetEventType, eventsOrAll(), newNetEventChecks(), NetEventTypes(), ParseNetEventType(), WithSuspectsFrom(), eventsetSummary(), modeLabel() (+6 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (30): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+22 more)

### Community 35 - "decryptStep"
Cohesion: 0.25
Nodes (7): askDecrypt(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), TestDecryptDirectoriesEmpty(), VerifyEncryptedKeys(), vaultKeyCheck

### Community 37 - "Message"
Cohesion: 0.09
Nodes (12): checkSecret(), File, backend, backendOp, backendState, fakeLink, Message, PwState (+4 more)

### Community 38 - "controlServer"
Cohesion: 0.22
Nodes (5): swapFile(), startControlServer(), TestRecordFailureLockout(), configEditor, controlServer

### Community 39 - "TempBinary"
Cohesion: 0.12
Nodes (16): TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant() (+8 more)

### Community 40 - "NetworkMonitor"
Cohesion: 0.10
Nodes (11): snapNetMonitor, NetBpfEvent, NetEvent, Cstr(), FormatAddr(), Ntohs(), NewNetworkMonitor(), statInodeKey() (+3 more)

### Community 41 - "main_test.go"
Cohesion: 0.06
Nodes (20): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+12 more)

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (5): IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "RunWith"
Cohesion: 0.16
Nodes (15): needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), stagedBinary, Mode, Request (+7 more)

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.18
Nodes (27): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+19 more)

### Community 46 - "CandidateDir"
Cohesion: 0.11
Nodes (16): symlinkTargetGlob(), expandPlaceholders(), CandidateDir, TestExpandWhitelist(), TestExtraWatchPathsForSkipsMissing(), TestPathsFor(), catalogEntryMatch(), findCatalogRoot() (+8 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (26): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+18 more)

### Community 50 - ".showSized"
Cohesion: 0.08
Nodes (22): installUI, multiAnswer, BinaryCheck, scrollable(), encodeTable(), invokingIDs(), TestEncodeTable(), TestWriteUserFileRefusesSymlinks() (+14 more)

### Community 51 - "cursor.go"
Cohesion: 0.14
Nodes (13): buttonCursor, handCursors(), handCursorsOn(), newButtonCursor(), showConfirm(), showCustomConfirm(), showDialog(), showErrorDialog() (+5 more)

### Community 52 - "Client"
Cohesion: 0.08
Nodes (17): Client, parseConfig(), childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend() (+9 more)

### Community 53 - "watchGroup"
Cohesion: 0.11
Nodes (36): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+28 more)

### Community 54 - "tabs.go"
Cohesion: 0.03
Nodes (17): appTheme, clickPanel, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, tabBar (+9 more)

### Community 55 - "Plan"
Cohesion: 0.09
Nodes (24): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), BinaryRule, Apply(), DeleteBackup(), Outcome, ProgressFunc (+16 more)

### Community 56 - "FileEvent"
Cohesion: 0.18
Nodes (5): snapMonitor, FileEvent, newFakeMonitorRepo(), fakeSource, fakeMonitorRepo

### Community 57 - "mutedText"
Cohesion: 0.04
Nodes (54): actionGate, binaryRow, binSource, counterStrip, filesPage, handButton, homePage, iconTile (+46 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (54): engine, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+46 more)

### Community 59 - "os.File"
Cohesion: 0.03
Nodes (80): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+72 more)

### Community 60 - "wipe"
Cohesion: 0.18
Nodes (20): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+12 more)

### Community 61 - "Recorder"
Cohesion: 0.15
Nodes (16): Target, drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), sortedKeys() (+8 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.19
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+7 more)

### Community 63 - "procstats.go"
Cohesion: 0.08
Nodes (24): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+16 more)

### Community 64 - "catalogRefresher"
Cohesion: 0.13
Nodes (11): drain(), newCatalogRefresher(), readEvents(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), awaitStartupOrSignal(), TestAwaitStartupOrSignalAbortsAndLocksDown(), TestAwaitStartupOrSignalAbortWithFailedStartup() (+3 more)

### Community 65 - "FileTarget"
Cohesion: 0.11
Nodes (14): guardPill(), newExeResolver(), Files, FileSpec, FileTarget, Policy, guardRow(), newFileTarget() (+6 more)

### Community 66 - "wm.go"
Cohesion: 0.11
Nodes (24): reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket(), poll() (+16 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "dialLiveSession"
Cohesion: 0.06
Nodes (39): ParseEventsFlag(), auditAfterEdit(), auditFindings(), reserveStdout(), runBackend(), dialLiveSession(), LiveModeAvailable(), TestValidateFlagsEditConfig() (+31 more)

### Community 69 - "update_test.go"
Cohesion: 0.10
Nodes (17): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel(), testKeyPair() (+9 more)

### Community 70 - "os.FileMode"
Cohesion: 0.10
Nodes (27): TestDaemonUnitWithMetadataOutput(), dupFD(), installConfig(), installFile(), installFileAs(), installServices(), installSSHAgent(), installSSHAgentEnv() (+19 more)

### Community 71 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 72 - "Guard"
Cohesion: 0.06
Nodes (26): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), deferredBinary, rootHandle, openRoot(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted() (+18 more)

### Community 73 - "fcntl"
Cohesion: 0.09
Nodes (13): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), answer() (+5 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 76 - "RefreshSection"
Cohesion: 0.08
Nodes (33): logRefreshChange(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange(), grantNewBinaries() (+25 more)

### Community 77 - "journal"
Cohesion: 0.15
Nodes (6): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 78 - "time.Time"
Cohesion: 0.08
Nodes (24): LockedError, addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors() (+16 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.19
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "singleBinaryPaths"
Cohesion: 0.16
Nodes (11): GlobChild, statFunc, multicallFamily(), singleBinaryPaths(), compileGlobNames(), globText(), resolveBits(), resolveChildren() (+3 more)

### Community 81 - "deploy"
Cohesion: 0.14
Nodes (20): applyUpdate(), archAssets(), assetsFor(), checkUpdatePreconditions(), confirmUnit(), deploy(), deployGUI(), downloadAndVerify() (+12 more)

### Community 82 - "planUpdaters"
Cohesion: 0.12
Nodes (20): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+12 more)

### Community 83 - "update/update.go"
Cohesion: 0.19
Nodes (23): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), channelOf(), compareStableVersions(), FetchLatest(), fetchReleases() (+15 more)

### Community 84 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 85 - "PlanDaemonUnitUpdate"
Cohesion: 0.31
Nodes (9): ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate(), ShippedDaemonUnit(), fakeBinary(), TestPlanDaemonUnitUpdate(), TestPlanDaemonUnitUpdateDropsLocalEdits(), TestShippedDaemonUnit() (+1 more)

### Community 87 - "buildTrustedSet"
Cohesion: 0.14
Nodes (17): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+9 more)

### Community 88 - "liveSession"
Cohesion: 0.18
Nodes (10): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), eventsLabel(), forwardBinaries(), holdForward(), normalizedEvents() (+2 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "Vault"
Cohesion: 0.06
Nodes (40): newBinaryVetter(), openBinaryVetter(), backingDeviceUnion(), buildConcurrency(), buildGuards(), lockOneRoot(), lockRootRecovering(), makeReloadHandler() (+32 more)

### Community 92 - "Monitor"
Cohesion: 0.16
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.18
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.08
Nodes (19): foreignPinnedLink(), reasonText(), eventTypeKey(), codeVouched(), isTraced(), judgePasses(), loaderEnv(), parseProcs() (+11 more)

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
Cohesion: 0.13
Nodes (14): procStats, statsSample, formatAgo(), cpuTime(), humanBytes(), rssBytes(), New(), recv() (+6 more)

### Community 100 - "gui/session.go"
Cohesion: 0.20
Nodes (14): ancestorEnv(), cleanupSessionEnv(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue(), TestOwnedByChecksTypeAndOwner() (+6 more)

### Community 101 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.17
Nodes (12): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+4 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "RunCmd"
Cohesion: 0.05
Nodes (49): installGUI(), applyLiveRefresh(), buildBinaryIfNeeded(), certifyInstall(), installBinaryOnly(), markComplete(), preflightDeployedBinary(), runRefreshCatalogOnly() (+41 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.20
Nodes (14): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+6 more)

### Community 106 - "recordops.c"
Cohesion: 0.45
Nodes (13): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+5 more)

### Community 107 - "sysStatus"
Cohesion: 0.18
Nodes (9): anyDaemonRunning(), sysStatus, probeSystem(), Check(), collectDiagnostics(), readSysctl(), requireKernel51(), DaemonRunning() (+1 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 110 - "runMaintenanceMode"
Cohesion: 0.18
Nodes (19): parseMaintenanceFlags(), runMaintenanceMode(), TestValidateMaintenanceFlags(), validateMaintenanceFlags(), deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup() (+11 more)

### Community 111 - "mountTable"
Cohesion: 0.13
Nodes (10): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+2 more)

### Community 112 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.06
Nodes (20): recordTickCmd(), recordModel, diffModel, editorModel, formatType(), listenForEvents(), syncViewport(), tickStats() (+12 more)

### Community 113 - "EventType"
Cohesion: 0.16
Nodes (14): EventType, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), specFromForm(), specSummary(), EventTypes(), ParseEventType(), maskFor() (+6 more)

### Community 115 - "fetchPair"
Cohesion: 0.15
Nodes (14): installHooks(), releaseBuild(), Requested(), Run(), staged(), TestRequested(), fetchPair(), FetchVerified() (+6 more)

### Community 116 - "CheckBPFLSM"
Cohesion: 0.17
Nodes (13): CheckBPFLSM(), confirmMasterKeyOverwrite(), runBPFCheck(), runGenKey(), runOneShotMode(), ensureInstalledBinary(), ensureMasterKey(), prepareInstallation() (+5 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "GenerateConf"
Cohesion: 0.05
Nodes (59): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), discordCatalogEntry(), groupedDiscordConf(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), TestGeneratedLibraryBlockRoundTrips(), TestGroupedSectionAddressedByEncryptionRoot() (+51 more)

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
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "IntegrationSuite"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "newTestTarget"
Cohesion: 0.22
Nodes (17): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 131 - "Candidate"
Cohesion: 0.08
Nodes (38): nextConfig(), addedResources(), appendSectionsAndEdit(), collectDiffAdditions(), configPaths(), pathCovered(), runDiffCatalog(), TestAddedResourcesSkipsAlreadyGuardedSSH() (+30 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "serveWebSocket"
Cohesion: 0.20
Nodes (9): configureWebSocket(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestSameOrigin(), TestSnapshotHubSingleViewerAndCoalescing(), writeSnapshot(), writeSnapshotLoop() (+1 more)

### Community 143 - "BinaryEntry"
Cohesion: 0.07
Nodes (40): binaryVerifyState, HashStore, keptEntries, sharedHashEntry, BinaryEntry, checkBtrfsKeys(), BinariesSummary(), canonicalBinaryPath() (+32 more)

### Community 144 - "Vault"
Cohesion: 0.20
Nodes (5): OpenRegularNoFollow(), Vault, Open(), Unopened(), WriteFileKeepMeta()

### Community 145 - "Open"
Cohesion: 0.16
Nodes (20): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+12 more)

### Community 146 - "Ledger"
Cohesion: 0.07
Nodes (12): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, applySchema(), Ledger (+4 more)

### Community 147 - "editControlSession"
Cohesion: 0.13
Nodes (18): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines() (+10 more)

### Community 148 - "New"
Cohesion: 0.16
Nodes (20): harnessSuite, looksLikeFileVaultRecord(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceWrongKeyFails(), TestUnlockLockFileInPlacePreservesInode(), TestUnlockLockFileInPlaceRoundTrip() (+12 more)

### Community 149 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 150 - "runMonitor"
Cohesion: 0.10
Nodes (19): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+11 more)

### Community 152 - "KernelDev"
Cohesion: 0.06
Nodes (39): confirmUserPlaced(), forwardCheck(), fdStat, BinaryStat, placedWalk, physicalParent(), statFD(), binaryStatOf() (+31 more)

### Community 153 - "shQuote"
Cohesion: 0.15
Nodes (4): IntegrationSuite, IntegrationSuite, shQuote(), IntegrationSuite

### Community 154 - "buildOneGuard"
Cohesion: 0.11
Nodes (28): multicallAdmissible(), buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards() (+20 more)

### Community 155 - "prepareDaemonStart"
Cohesion: 0.09
Nodes (25): CheckEBPF(), configureDaemonLogging(), loadDaemonConfig(), markReady(), newPinGeneration(), noteReloadedSettings(), notifySystemdReady(), prepareDaemonStart() (+17 more)

### Community 156 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 157 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.27
Nodes (19): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+11 more)

### Community 161 - ".Layout"
Cohesion: 0.44
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 162 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 163 - "vaultfs_test.go"
Cohesion: 0.12
Nodes (21): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), IsBinaryContent(), IsOctalMode(), inodeOf(), mustRead(), mustWrite() (+13 more)

### Community 165 - "sync.Mutex"
Cohesion: 0.06
Nodes (46): actionBar, cell, column, dataTable, eventFeed, feedAction, feedSpec, fyneLogFilter (+38 more)

### Community 166 - "shell"
Cohesion: 0.06
Nodes (36): pageDef, pageID, reel, resizeWatch, runHost, shell, sizedDialog, snapEngines (+28 more)

### Community 167 - "revertSSHAgents"
Cohesion: 0.22
Nodes (10): detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeSSHAgentEnv(), removeSSHAgentUnit(), revertSSHAgents(), staleRCUsers(), TestIsInstallerSSHAgentUnit(), TestRemoveSSHAgentUnit() (+2 more)

### Community 168 - "mustSubkey"
Cohesion: 0.31
Nodes (9): newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape(), TestOpenFileVaultRejectsTampering(), TestSealFileVaultNoncesNeverRepeat() (+1 more)

### Community 169 - "Config"
Cohesion: 0.07
Nodes (35): sameConfig(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), startTrustGuard(), trustStartupError(), auditEntry(), auditFindingsWithConfig(), auditGroupCoverage() (+27 more)

### Community 170 - "io.Writer"
Cohesion: 0.18
Nodes (12): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), TestProgressReader(), bottomBar (+4 more)

### Community 171 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 175 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 176 - "netinject.c"
Cohesion: 0.50
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 177 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 178 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 179 - "SanitizeText"
Cohesion: 0.08
Nodes (36): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), formatElapsed(), lines(), offerBackupDeletion() (+28 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (19): AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+11 more)

### Community 186 - "Binary trust"
Cohesion: 0.25
Nodes (8): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, The dynamic loader, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 194 - "refuseOutsideHardlinks"
Cohesion: 0.28
Nodes (8): lstatOf(), outsideHardlinks(), refuseOutsideHardlinks(), TestRefuseOutsideHardlinks_FileRoot(), TestRefuseOutsideHardlinks_InTreeLinksAccepted(), TestRefuseOutsideHardlinks_MissingRootFails(), TestRefuseOutsideHardlinks_OutsideLinkRefused(), writeFile()

### Community 195 - "eventUnitTest"
Cohesion: 0.32
Nodes (3): BpfEvent, eventUnitTest, DecodeBpfEvent()

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `--record` (+2 more)

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "run"
Cohesion: 0.36
Nodes (8): InstalledGUIVersion(), InstalledVersion(), readInstalledVersion(), run(), runUpdate(), syncGUI(), TestReadInstalledVersion(), GUIInstalled()

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 202 - "DaemonEvent"
Cohesion: 0.08
Nodes (41): UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemon(), runDaemonUI(), runHeadless() (+33 more)

### Community 203 - "parseResize"
Cohesion: 0.40
Nodes (3): decodeEndsAtEOF(), parseResize(), TestParseResize()

### Community 204 - "runNetworkMonitor"
Cohesion: 0.13
Nodes (10): ParseNetEventsFlag(), runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NetworkMonitorRepository, NewNetModel() (+2 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.15
Nodes (13): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Protected (+5 more)

### Community 209 - "User"
Cohesion: 0.12
Nodes (24): entryWriters(), reserveChain(), underResource(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir() (+16 more)

### Community 210 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 213 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 216 - "searchDirs"
Cohesion: 0.18
Nodes (10): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths(), resolveSoname() (+2 more)

### Community 217 - "options.go"
Cohesion: 0.31
Nodes (9): InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks, InstallHooks, StagedRelease (+1 more)

### Community 221 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 222 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 223 - "run"
Cohesion: 0.14
Nodes (18): confirmKeyDeletion(), run(), runUninstall(), RunWith(), removeKeyAndEmptyDir(), removeMasterKey(), revertBunLaunchers(), TestRemoveKeyAndEmptyDir() (+10 more)

### Community 226 - "Prereq"
Cohesion: 0.20
Nodes (6): classifySetupError(), TestClassifySetupErrorGeneric(), TestClassifySetupErrorNotSetup(), TestClassifySetupErrorNotSupported(), Prereq, fakeVault

### Community 227 - "ptrace_race.c"
Cohesion: 0.52
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 228 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 231 - "DiscoverInfraBinaries"
Cohesion: 0.33
Nodes (6): DiscoverInfraBinaries(), infraFromRunning(), runningExecutables(), pathIsRunning(), TestDiscoverInfraBinaries(), TestInfraFromRunning()

### Community 232 - "checkRunningBinaryMatchesInstalled"
Cohesion: 0.40
Nodes (5): checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), fileSHA256(), sameContent(), TestCheckRunningBinaryMatchesInstalled()

### Community 235 - "parseKernelRelease"
Cohesion: 0.67
Nodes (4): ciliumVersion, parseKernelRelease(), TestParseKernelRelease(), unameKernelVersion()

### Community 237 - "serve.go"
Cohesion: 0.15
Nodes (18): basicAuth(), isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), securityHeaders(), Serve() (+10 more)

### Community 238 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

## Knowledge Gaps
- **32 isolated node(s):** `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener`, `open_by_handle_at POC`, `statonly stat/statx metadata leak` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 597 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **41 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `shQuote`, `IntegrationSuite`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **What connects `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.0934448531264171 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.0164000537706681 - nodes in this community are weakly interconnected._
- **Why does `Editor keys` connect `edit-protected` to `liveSession`, `runMaintenanceMode`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08795845504706264 - nodes in this community are weakly interconnected._