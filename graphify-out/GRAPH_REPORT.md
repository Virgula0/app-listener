# Graph Report - app-listener  (2026-10-20)

## Corpus Check
- 505 files · ~1,723,856 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6614 nodes · 28953 edges · 235 communities (195 shown, 40 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2677 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `45e22f1c`
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
- BinaryEntry
- sync.Mutex
- protectedPage
- selection.go
- render.go
- recorder_test.go
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
- .resKey
- github.com/charmbracelet/bubbles/textarea.Model
- codeedit.go
- shellrc.go
- .FilterExistingWhitelist
- TestMonitorUseCaseLifecycle
- UI
- sanitizeTerminalText
- picker
- dialLiveSession
- time.Duration
- Message
- controlServer
- TempBinary
- NetEventType
- IntegrationSuite
- absPath
- monitorUnitTest
- desktop.go
- fileedit_model_test.go
- CandidateDir
- GuardInodeKey
- fadeScroll
- install/install.go
- cursor.go
- proto.go
- Config
- navItem
- Plan
- EventType
- mutedText
- engine
- os.File
- New
- Recorder
- exe_supersede.h
- installPage
- catalogRefresher
- uninstall_system.go
- context.Context
- IntegrationSuite
- auditGroupCoverage
- update_test.go
- Network
- WithSuspectsFrom
- Guard
- fcntl
- bpfstats/main.go
- networkPage
- codecheck/main.go
- journal
- time.Time
- process_vm_readv.c
- GuardEvent
- LibraryClosure
- gui/session.go
- changelog.go
- net_tester/main.go
- ValidateTargets
- recordFlow
- daemon_multicall_test.go
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- eventUnitTest
- Monitor
- prctl_exe_swap.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- pulse_test.go
- ConfinedEntry
- fakeGuardRepo
- Well-known bypass classes
- hashCache
- jit_provenance.c
- inspector_probe.c
- recordops.c
- protected_test.go
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- mayUpdate
- mountTable
- github.com/charmbracelet/bubbletea.Cmd
- NetGuardEvent
- IntegrationSuite
- io.Writer
- TrustGuard
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
- vaultfs_test.go
- trace-app-libs.sh
- IntegrationSuite
- vaultfs.go
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- IntegrationSuite
- Serve
- parseGuardEvent
- Open
- Ledger
- buildTrustedSet
- .step
- changelog_test.go
- runNetworkGuard
- fakeStore
- inodeChain
- IntegrationSuite
- buildOneGuard
- NetworkMonitor
- downloadFile
- runDaemon
- Installation
- backend_test.go
- graphify-refresh.sh
- layout.go
- .decryptFileWithProgress
- .populateWatchInodes
- desktop_test.go
- shell
- .showSized
- safeio.go
- WithBottomBar
- NewDumpHook
- wm.go
- Limitations
- Troubleshooting
- IntegrationSuite
- RawTarget
- netinject.c
- CopyTreeWithProgress
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
- fileedit_symlink_test.go
- NewEventFanout
- edit-protected
- monitor: see what touches a path
- kernel-prof.sh
- update/update.go
- pathCache
- IntegrationSuite
- DaemonEvent
- batchInodeRows
- fakeDaemon
- uninstall/uninstall.go
- Desktop GUI
- IntegrationSuite
- github.com/stretchr/testify/suite.Suite
- wm_test.go
- .eventTypeAllowed
- bleedLayout
- AddServeFlags
- vouchedDevs
- .awaitDaemonUpOrDead
- IntegrationSuite
- KernelDev
- tempGrantSetup
- CleanupStalePins
- entryOf
- ptrace_race.c
- prctl_gate
- Choice
- formatEntry
- supersede_probe.c
- SuperblockDev
- .runEditorHarness
- Execute
- coverage_test.go
- github.com/Virgula0/app-listener/tools/codecheck
- taint_victim.c
- configedit_test.go
- TestMain
- Daemon configuration

## God Nodes (most connected - your core abstractions)
1. `Guard` - 118 edges
2. `Config` - 109 edges
3. `GuardInodeKey` - 100 edges
4. `shell` - 91 edges
5. `Load()` - 89 edges
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

## Communities (235 total, 40 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.11
Nodes (27): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion, multicallCacheKey (+19 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (144): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity() (+136 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): Hash(), configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec() (+2 more)

### Community 3 - "fscrypt.go"
Cohesion: 0.05
Nodes (63): encryptDirectories(), secureResources(), verifyEncryptionState(), askEncryption(), deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), checkKeyLen() (+55 more)

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
Cohesion: 0.06
Nodes (27): addErrKind, ReplacementCheck, TestWithChmodDropWriteNeedsReadOnly(), classifyAddErr(), guardModeKey(), modeString(), NewGuard(), addRecorder() (+19 more)

### Community 9 - "mkdirs"
Cohesion: 0.12
Nodes (40): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), newWatchPattern(), systemBinaryPatterns(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped() (+32 more)

### Community 10 - "BinaryEntry"
Cohesion: 0.10
Nodes (14): deferredBinary, keptEntries, BinaryEntry, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), checkBtrfsKeys(), canonicalBinaryPath(), canonicalPaths(), confinedEntry() (+6 more)

### Community 11 - "sync.Mutex"
Cohesion: 0.06
Nodes (25): guardPill(), newActivity(), EntryFor(), Exe, newExeResolver(), processKey(), Files, FileSpec (+17 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (19): backendLink, checkGrid, editView, forwardForm, forwardView, protectedPage, recordForm, recordView (+11 more)

### Community 13 - "selection.go"
Cohesion: 0.10
Nodes (16): pos, Highlighter, cellToRune(), clipboardText(), cursorPos(), Selection, offset(), setClipboard() (+8 more)

### Community 14 - "render.go"
Cohesion: 0.15
Nodes (12): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+4 more)

### Community 15 - "recorder_test.go"
Cohesion: 0.36
Nodes (4): recordTick, isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (6): eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "runView"
Cohesion: 0.06
Nodes (14): daemonPage, daemonState, logGrid, logHook, logLine, runHost, runView, journalLine() (+6 more)

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
Cohesion: 0.06
Nodes (16): ownerPaths(), libBinaries(), planUpdaters(), resourceBinaries(), resourceLinks(), TestPlanUpdaters_LinkOwnsItsResourceBit(), Resource, GuardRepository (+8 more)

### Community 23 - "auth_test.go"
Cohesion: 0.13
Nodes (38): CheckPassword(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf(), parseHash() (+30 more)

### Community 24 - "pinstate.go"
Cohesion: 0.07
Nodes (44): configureDaemonLogging(), loadDaemonConfig(), lockOneRoot(), lockRootRecovering(), noteReloadedSettings(), relockStaleVaults(), resolveConfigPath(), runLockdown() (+36 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (73): exe_refused(), inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade() (+65 more)

### Community 26 - ".resKey"
Cohesion: 0.08
Nodes (10): exeRow, tempGrant, tempMaskOp, tempRow, eventMask(), readPinnedResID(), requirePinnedSelfAllowRoot(), WithPinnedSelfVaultAccess() (+2 more)

### Community 27 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.13
Nodes (20): findButton, span, column(), findAll(), Finder, lowerRunes(), NewFinder(), newFindInput() (+12 more)

### Community 28 - "codeedit.go"
Cohesion: 0.08
Nodes (48): nextConfig(), installConfig(), class, classOf(), dedent(), Edit(), leadingIndent(), lines() (+40 more)

### Community 29 - "shellrc.go"
Cohesion: 0.16
Nodes (35): rcTarget, BunTmpDir(), applyRC(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), EnsureBunTmpDir(), ensureRCBlock() (+27 more)

### Community 30 - ".FilterExistingWhitelist"
Cohesion: 0.12
Nodes (19): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+11 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 32 - "UI"
Cohesion: 0.08
Nodes (25): showChangelog(), confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), InstallMode, InstallRequest, Options (+17 more)

### Community 33 - "sanitizeTerminalText"
Cohesion: 0.07
Nodes (31): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+23 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (30): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+22 more)

### Community 35 - "dialLiveSession"
Cohesion: 0.09
Nodes (26): TestValidatePassword(), ValidatePassword(), reserveStdout(), runBackend(), dialLiveSession(), LiveModeAvailable(), applyEdit(), checkAndApply() (+18 more)

### Community 36 - "time.Duration"
Cohesion: 0.14
Nodes (9): procStats, statsSample, IntegrationSuite, pooledContainer, formatClock(), cpuTime(), humanBytes(), rssBytes() (+1 more)

### Community 37 - "Message"
Cohesion: 0.12
Nodes (7): checkSecret(), backend, backendOp, backendState, Message, Request, RequireDaemonStopped()

### Community 38 - "controlServer"
Cohesion: 0.10
Nodes (22): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines() (+14 more)

### Community 39 - "TempBinary"
Cohesion: 0.12
Nodes (16): TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant() (+8 more)

### Community 40 - "NetEventType"
Cohesion: 0.07
Nodes (18): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), newNetEventChecks(), FormatAddr(), NetEventTypes() (+10 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (7): IntegrationSuite, IntegrationSuite, shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "desktop.go"
Cohesion: 0.07
Nodes (41): installHooks(), releaseBuild(), Run(), staged(), installGUI(), needsBinary(), RunWith(), TestNeedsBinary() (+33 more)

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "CandidateDir"
Cohesion: 0.06
Nodes (37): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), groupCandidates(), pickDirectories(), pickFromCandidates() (+29 more)

### Community 48 - "GuardInodeKey"
Cohesion: 0.06
Nodes (31): refusedReplacements, verifiedBinary(), deleteInoKeys(), supersedeUnnamed(), appletTag(), buildOf(), GuardInodeKey, mcNameKey() (+23 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (27): Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+19 more)

### Community 50 - "install/install.go"
Cohesion: 0.04
Nodes (91): anyDaemonRunning(), offerDaemonStart(), applyLiveRefresh(), asRoot(), buildBinaryIfNeeded(), certifyInstall(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled() (+83 more)

### Community 51 - "cursor.go"
Cohesion: 0.08
Nodes (19): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+11 more)

### Community 52 - "proto.go"
Cohesion: 0.07
Nodes (23): Client, File, Progress, User, childEnv(), launch(), Start(), helper() (+15 more)

### Community 53 - "Config"
Cohesion: 0.06
Nodes (64): newBinaryVetter(), openBinaryVetter(), newPinGeneration(), reloadOnce(), reloadSlotsNeeded(), setConfinedHomes(), appendSectionsAndEdit(), grantNewBinaries() (+56 more)

### Community 54 - "navItem"
Cohesion: 0.03
Nodes (16): appTheme, clickPanel, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, tabBar (+8 more)

### Community 55 - "Plan"
Cohesion: 0.11
Nodes (18): BinaryRule, Plan, materialized(), TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan() (+10 more)

### Community 56 - "EventType"
Cohesion: 0.09
Nodes (16): snapMonitor, EventType, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), specFromForm(), specSummary() (+8 more)

### Community 57 - "mutedText"
Cohesion: 0.06
Nodes (44): binaryRow, counterStrip, filesPage, homePage, iconTile, panel, pill, statusTile (+36 more)

### Community 58 - "engine"
Cohesion: 0.05
Nodes (29): engine, resInfo, taintOwner, TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded(), pathWithin() (+21 more)

### Community 59 - "os.File"
Cohesion: 0.05
Nodes (65): multicallAdmissible(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir() (+57 more)

### Community 60 - "New"
Cohesion: 0.09
Nodes (43): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault() (+35 more)

### Community 61 - "Recorder"
Cohesion: 0.16
Nodes (15): Target, Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start(), record() (+7 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.19
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+7 more)

### Community 63 - "installPage"
Cohesion: 0.08
Nodes (13): actionGate, binSource, handButton, handCheck, installAction, installPage, latestRelease, uninstallCard (+5 more)

### Community 64 - "catalogRefresher"
Cohesion: 0.14
Nodes (9): drain(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), catalogRefresher (+1 more)

### Community 65 - "uninstall_system.go"
Cohesion: 0.05
Nodes (68): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir() (+60 more)

### Community 66 - "context.Context"
Cohesion: 0.19
Nodes (9): centered(), FitSize(), poll(), TestFitSize(), hyprClient, hyprMonitor, Rect, Size (+1 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "auditGroupCoverage"
Cohesion: 0.15
Nodes (16): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource() (+8 more)

### Community 69 - "update_test.go"
Cohesion: 0.08
Nodes (22): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestAssetsFor(), TestFetchReleases(), TestFetchReleasesHTTPError() (+14 more)

### Community 70 - "Network"
Cohesion: 0.15
Nodes (7): netGuardPill(), stopPumped(), clonePolicy(), NetExtra, NetPolicy, Network, verifyHashes()

### Community 71 - "WithSuspectsFrom"
Cohesion: 0.43
Nodes (3): WithSuspectsFrom(), netguardHarness, Option

### Community 72 - "Guard"
Cohesion: 0.09
Nodes (12): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), rootHandle, SectionScan, openRoot(), SharedPinDegraded() (+4 more)

### Community 73 - "fcntl"
Cohesion: 0.09
Nodes (11): main(), report(), denied(), main(), put(), wait_then_read(), dump_via_debugfs(), main() (+3 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 75 - "networkPage"
Cohesion: 0.07
Nodes (18): binaryList, netGuardTab, networkPage, policyEditor, policyHooks, policyTab, watchTab, newBinaryList() (+10 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "journal"
Cohesion: 0.15
Nodes (5): fakeGuard, fakeMonitor, fakeNet, fakeNetMonitor, journal

### Community 78 - "time.Time"
Cohesion: 0.12
Nodes (13): LockedError, isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, sweepID (+5 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.19
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.07
Nodes (24): displayEntries(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI(), snapGuard (+16 more)

### Community 81 - "LibraryClosure"
Cohesion: 0.07
Nodes (29): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+21 more)

### Community 82 - "gui/session.go"
Cohesion: 0.21
Nodes (15): ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+7 more)

### Community 83 - "changelog.go"
Cohesion: 0.23
Nodes (18): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+10 more)

### Community 84 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 85 - "ValidateTargets"
Cohesion: 0.12
Nodes (13): drain(), dedupe(), existingSection(), guardedMatch(), libOfChosen(), overlaps(), refusedPath(), sectionFor() (+5 more)

### Community 86 - "recordFlow"
Cohesion: 0.13
Nodes (17): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+9 more)

### Community 87 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 88 - "liveSession"
Cohesion: 0.21
Nodes (10): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), eventsLabel(), forwardBinaries(), holdForward(), normalizedEvents() (+2 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "eventUnitTest"
Cohesion: 0.27
Nodes (4): BpfEvent, eventUnitTest, Cstr(), DecodeBpfEvent()

### Community 92 - "Monitor"
Cohesion: 0.17
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.18
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.08
Nodes (13): lsmHook, FirstOnBtrfs(), reasonText(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses() (+5 more)

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

### Community 100 - "ConfinedEntry"
Cohesion: 0.13
Nodes (29): binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats() (+21 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "hashCache"
Cohesion: 0.19
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

### Community 107 - "protected_test.go"
Cohesion: 0.31
Nodes (7): fakeLink, newFakeLink(), protectedShell(), TestForwardRefusesOwnExecutable(), TestProtectedEditFlow(), TestProtectedLiveCardsGated(), TestRecordReviewExtendsAnExistingSection()

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "mayUpdate"
Cohesion: 0.18
Nodes (12): hardLinked(), mayUpdate(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning(), TestIsGeneralTool(), TestIsGeneralToolJudgesTheInode(), TestMayUpdate_HardLinkWaivedForUserWritableLibBinary(), TestPlanUpdaters_ScopedToOwnResources() (+4 more)

### Community 111 - "mountTable"
Cohesion: 0.18
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+1 more)

### Community 112 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.07
Nodes (16): recordTickCmd(), recordModel, diffModel, editorModel, protoString(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented() (+8 more)

### Community 113 - "NetGuardEvent"
Cohesion: 0.16
Nodes (7): snapNetGuard, NetGuardEvent, NetworkGuardRepository, RealNetEngines(), NetEngines, fakeNetGuard, NetGuardConfig

### Community 115 - "io.Writer"
Cohesion: 0.27
Nodes (8): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), fyneLogFilter, candidate

### Community 116 - "TrustGuard"
Cohesion: 0.05
Nodes (30): statFunc, trustHook, trustRows, ensureBtrfsLayout(), deleteResKeys(), dropSupersededLocked(), PruneSuperseded(), setSupersedeTrust() (+22 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "GenerateConf"
Cohesion: 0.06
Nodes (56): logRefreshChange(), TestLogRefreshChange(), LibraryBlock, RefreshOptions, ParseLibrariesSection(), AppendSectionWhitelist(), ConfSafePath(), findSectionEnd() (+48 more)

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

### Community 131 - "vaultfs_test.go"
Cohesion: 0.20
Nodes (14): IsBinaryContent(), IsOctalMode(), ParseMode(), inodeOf(), mustRead(), mustWrite(), openRootT(), TestIsBinaryContent() (+6 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "vaultfs.go"
Cohesion: 0.16
Nodes (7): DefaultNewFileMode(), Vault, Open(), TestDefaultNewFileMode(), TestOpenVaultRefusesSymlink(), Unopened(), WriteFileInPlace()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 143 - "Serve"
Cohesion: 0.08
Nodes (29): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+21 more)

### Community 144 - "parseGuardEvent"
Cohesion: 0.14
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 145 - "Open"
Cohesion: 0.16
Nodes (20): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+12 more)

### Community 146 - "Ledger"
Cohesion: 0.11
Nodes (5): applySchema(), Ledger, Pending, Source, fakeLedger

### Community 147 - "buildTrustedSet"
Cohesion: 0.15
Nodes (16): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+8 more)

### Community 148 - ".step"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 149 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 150 - "runNetworkGuard"
Cohesion: 0.07
Nodes (35): DrainUntilSignal(), CheckBPFLSM(), CheckEBPF(), MakeDisplayPaths(), ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), confirmMasterKeyOverwrite() (+27 more)

### Community 152 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.09
Nodes (38): backingDeviceUnion(), buildGuards(), buildOneGuard(), drainEphemeral(), eventFilterOptions(), startGuardedDaemon(), systemPatterns(), systemRule() (+30 more)

### Community 155 - "NetworkMonitor"
Cohesion: 0.15
Nodes (7): BinariesSummary(), TestBinariesSummary(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 156 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 157 - "runDaemon"
Cohesion: 0.11
Nodes (17): newControlManager(), startControlManager(), makeReloadHandler(), runDaemon(), startCatalogRefresh(), startPprof(), writePidFile(), resolveInspectors() (+9 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.27
Nodes (19): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+11 more)

### Community 161 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 163 - ".decryptFileWithProgress"
Cohesion: 0.21
Nodes (5): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata()

### Community 165 - "desktop_test.go"
Cohesion: 0.06
Nodes (37): TestEveryInstallModeMapped(), actionBar, cell, codeEntry, column, dataTable, eventFeed, feedAction (+29 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (45): pageDef, pageID, reel, resizeWatch, shell, sizedDialog, snapEngines, snapNet (+37 more)

### Community 167 - ".showSized"
Cohesion: 0.12
Nodes (11): configView, installUI, multiAnswer, recordReview, scrollable(), selectable(), newFadeScroll(), ask() (+3 more)

### Community 168 - "safeio.go"
Cohesion: 0.27
Nodes (12): invokingIDs(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), createOwned(), AtomicWriteAt(), CreateExclNoFollow(), descend(), DescendCreateOwned() (+4 more)

### Community 169 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 170 - "NewDumpHook"
Cohesion: 0.24
Nodes (10): applyVerbosity(), VerboseLevel, NewDumpHook(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity(), TestDumpHookFormat(), TestVerboseToLogrus() (+2 more)

### Community 171 - "wm.go"
Cohesion: 0.29
Nodes (8): detectWM(), reachable(), Detect(), Session, hyprlandSocket(), i3RuntimeSocket(), socketOwnedBy(), Kind

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 175 - "RawTarget"
Cohesion: 0.31
Nodes (9): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+1 more)

### Community 176 - "netinject.c"
Cohesion: 0.50
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 177 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 178 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 179 - "Binary"
Cohesion: 0.10
Nodes (29): formatElapsed(), lines(), orNone(), pickBinaries(), pickLibs(), printBinary(), printRecording(), reviewRecording() (+21 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "github.com/spf13/cobra.Command"
Cohesion: 0.16
Nodes (18): credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestSetupDumpLogRefusesSymlink(), TestMovedFlagError() (+10 more)

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

### Community 194 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 195 - "NewEventFanout"
Cohesion: 0.24
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "update/update.go"
Cohesion: 0.12
Nodes (35): applyUpdate(), archAssets(), assetsFor(), channelOf(), checkUpdatePreconditions(), confirmUnit(), deploy(), deployGUI() (+27 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 202 - "DaemonEvent"
Cohesion: 0.08
Nodes (42): ServeConfig, UIDResolver, NewUIDResolver(), awaitStartupOrSignal(), catchLifecycleSignals(), newDaemonModel(), runDaemonUI(), runHeadless() (+34 more)

### Community 203 - "batchInodeRows"
Cohesion: 0.20
Nodes (13): markReady(), notifySystemdReady(), inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows() (+5 more)

### Community 205 - "uninstall/uninstall.go"
Cohesion: 0.07
Nodes (35): TestPickEncryptedDirSingleSkipsPicker(), pickDir(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), confirmKeyDeletion(), decryptDirectories(), decryptStep() (+27 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.14
Nodes (14): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+6 more)

### Community 210 - "wm_test.go"
Cohesion: 0.67
Nodes (8): envOf(), serve(), shortDir(), TestDetect(), TestHyprlandFloat(), TestHyprlandFloatTimesOutWithoutWindow(), TestI3Float(), uid()

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 213 - "AddServeFlags"
Cohesion: 0.36
Nodes (8): AddHeadlessFlag(), AddServeFlags(), TestParseServeFlagsRequiresCredentials(), init(), init(), init(), init(), init()

### Community 214 - "vouchedDevs"
Cohesion: 0.33
Nodes (6): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs()

### Community 217 - "KernelDev"
Cohesion: 0.10
Nodes (28): fdStat, BinaryStat, physicalParent(), statFD(), binaryStatOf(), binaryStatOf(), KernelDev(), StatBinary() (+20 more)

### Community 218 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 219 - "CleanupStalePins"
Cohesion: 0.33
Nodes (6): CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 220 - "entryOf"
Cohesion: 0.40
Nodes (5): entryOf(), FileModeToUnixMode(), TestFileModeToUnixModeRoundTrip(), TestUnixModeToFileMode(), UnixModeToFileMode()

### Community 221 - "ptrace_race.c"
Cohesion: 0.52
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 222 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 223 - "Choice"
Cohesion: 0.50
Nodes (3): MultiSelectKeymap(), Choice, MultiSelect()

### Community 224 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 225 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 226 - "SuperblockDev"
Cohesion: 0.40
Nodes (5): mountinfoDev(), OnBtrfs(), SuperblockDev(), TestMountinfoDev(), TestSuperblockDevMatchesStatOffBtrfs()

### Community 228 - "Execute"
Cohesion: 0.67
Nodes (3): Execute(), initLogger(), main()

### Community 238 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 240 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 244 - "Daemon configuration"
Cohesion: 0.20
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 615 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **40 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `go_pkg_strings`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `TestMain`, `IntegrationSuite`, `IntegrationSuite`, `daemon_multicall_test.go`, `.TestRecord_GuardParity`, `IntegrationSuite`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.11426340447683499 - nodes in this community are weakly interconnected._
- **Why does `Editor keys` connect `edit-protected` to `desktop.go`, `Message`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.01986214402321785 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08957708957708958 - nodes in this community are weakly interconnected._