# Graph Report - app-listener  (2026-10-18)

## Corpus Check
- 486 files · ~1,711,574 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6455 nodes · 27858 edges · 228 communities (192 shown, 36 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2587 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `835002e1`
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
- binaryList
- protectedPage
- github.com/charmbracelet/bubbletea.Cmd
- NetworkMonitor
- editView
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- go_pkg_path_filepath
- guard.bpf.c
- TrustGuard
- Finder
- netModel
- shellrc.go
- Config
- TestMonitorUseCaseLifecycle
- UI
- NetEventType
- picker
- selection.go
- IntegrationSuite
- backend
- controlServer
- fileedit_model_test.go
- .FilterExistingWhitelist
- IntegrationSuite
- absPath
- monitorUnitTest
- uninstall/uninstall.go
- GuardEvent
- daemon_multicall_test.go
- stdio
- User
- fadeScroll
- backend.go
- cursor.go
- Client
- daemonconfig.go
- navItem
- .apply
- github.com/spf13/cobra.Command
- mutedText
- GuardInodeKey
- os.File
- expandPlaceholders
- Recorder
- exe_supersede.h
- Network
- CandidateDir
- bpfstats/main.go
- gui/session.go
- newTestTarget
- recordFlow
- vouchedDevs
- fcntl
- proto.go
- Guard
- FileTarget
- pickfs.go
- catalogRefresher
- update_test.go
- GlobReservations
- process_vm_readv.c
- integrationtests/networkguard_test.go
- journal
- planUpdaters
- changelog.go
- net_tester/main.go
- highlight_test.go
- github.com/charmbracelet/bubbles/textarea.Model
- loadDaemonConfig
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- configedit_test.go
- Monitor
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- globBuilder
- codeedit.go
- downloadAndVerify
- Well-known bypass classes
- updatePage
- jit_provenance.c
- inspector_probe.c
- recordops.c
- tempGrantSetup
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- hashCache
- mountTable
- NetGuardEvent
- .exeKeys
- IntegrationSuite
- WithSuspectsFrom
- FileEvent
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- .TestRecord_GuardParity
- RefreshSection
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- .filterEvent
- make test-integration (rootful Docker suite)
- .launch
- supersede_probe.c
- install/install.go
- trace-app-libs.sh
- IntegrationSuite
- dialLiveSession
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- fakeGuardRepo
- time.Duration
- copyRegularAt
- Open
- find_test.go
- runView
- Binary
- orphans_test.go
- EventType
- fakeStore
- New
- shQuote
- SanitizeText
- .Layout
- buildTrustedSet
- changelog_test.go
- Installation
- backend_test.go
- graphify-refresh.sh
- vaultfs.go
- inodeChain
- controlManager
- RawTarget
- sync.Mutex
- shell
- github.com/stretchr/testify/suite.Suite
- Ledger
- taint_victim.c
- selection_test.go
- dialog_test.go
- Limitations
- Troubleshooting
- IntegrationSuite
- IntegrationSuite
- markReady
- main_test.go
- Vault
- time.Time
- Compatibility
- How it works
- NewDumpHook
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- downloadFile
- testLimiter
- edit-protected
- monitor: see what touches a path
- kernel-prof.sh
- securewrite_test.go
- pathCache
- DaemonEvent
- GenerateConf
- io.Writer
- audit.go
- Desktop GUI
- IntegrationSuite
- WithBottomBar
- .refreshOnce
- .eventTypeAllowed
- resolveBits
- diffModel
- bleedLayout
- IntegrationSuite
- LibraryClosure
- grantNewBinaries
- fakeDaemon
- confirmUpdate
- update/update.go
- Daemon configuration
- watchSet
- bpfLsmListed
- parseGuardEvent
- unit.go
- Prereq
- serve.go

## God Nodes (most connected - your core abstractions)
1. `Guard` - 118 edges
2. `Config` - 107 edges
3. `GuardInodeKey` - 99 edges
4. `Load()` - 89 edges
5. `shell` - 87 edges
6. `IntegrationSuite` - 82 edges
7. `absPath()` - 82 edges
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
- `runConfigEditor()` --calls--> `EditText()`  [INFERRED]
  cmd/functions/install/steps.go → internal/install/editor.go

## Import Cycles
- None detected.

## Communities (228 total, 36 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.15
Nodes (31): atomicWriteAt(), writeAndSync(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs (+23 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (130): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), awaitStartupOrSignal(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestAwaitStartupOrSignalAbortsAndLocksDown(), TestAwaitStartupOrSignalAbortWithFailedStartup() (+122 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "Vault"
Cohesion: 0.06
Nodes (37): deprovisionKind, harnessSuite, migrationTarget, isFileVaultCiphertext(), isRegularFileTarget(), writeTempWithMetadata(), classifyDeprovision(), classifyDeprovisionErr() (+29 more)

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
Nodes (23): addErrKind, ReplacementCheck, classifyAddErr(), commMatchesGuardedBinary(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest (+15 more)

### Community 9 - "mkdirs"
Cohesion: 0.11
Nodes (42): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), newWatchPattern(), newFakeInotify() (+34 more)

### Community 10 - "BinaryEntry"
Cohesion: 0.07
Nodes (39): binaryVerifyState, HashStore, keptEntries, sharedHashEntry, BinaryEntry, checkBtrfsKeys(), ensureBtrfsLayout(), BinariesSummary() (+31 more)

### Community 11 - "binaryList"
Cohesion: 0.09
Nodes (14): binaryList, netGuardTab, policyEditor, policyHooks, policyTab, newBinaryList(), newPolicyTab(), verbFor() (+6 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (22): backendLink, configView, forwardForm, forwardView, installUI, protectedPage, recordForm, recordReview (+14 more)

### Community 13 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.11
Nodes (10): formatElapsed(), recordTickCmd(), recordModel, renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel (+2 more)

### Community 14 - "NetworkMonitor"
Cohesion: 0.08
Nodes (14): BpfEvent, eventUnitTest, NetBpfEvent, Cstr(), DecodeBpfEvent(), FormatAddr(), Ntohs(), reasonText() (+6 more)

### Community 15 - "editView"
Cohesion: 0.12
Nodes (7): codeEntry, editView, Entry, fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (57): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+49 more)

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
Cohesion: 0.08
Nodes (9): Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots(), vaultOpForGuard() (+1 more)

### Community 23 - "auth.go"
Cohesion: 0.14
Nodes (37): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+29 more)

### Community 24 - "go_pkg_path_filepath"
Cohesion: 0.34
Nodes (3): isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "TrustGuard"
Cohesion: 0.06
Nodes (21): tempGrant, trustHook, trustRows, deleteResKeys(), guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), dropSupersededLocked(), PruneSuperseded() (+13 more)

### Community 27 - "Finder"
Cohesion: 0.16
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 28 - "netModel"
Cohesion: 0.07
Nodes (28): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+20 more)

### Community 29 - "shellrc.go"
Cohesion: 0.16
Nodes (34): setupBunTmpdir(), rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), EnsureBunTmpDir(), ensureRCBlock() (+26 more)

### Community 30 - "Config"
Cohesion: 0.08
Nodes (38): newCatalogRefresher(), makeReloadHandler(), startCatalogRefresh(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), startTrustGuard(), trustStartupError(), appendSectionsAndEdit() (+30 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.19
Nodes (14): NewGuardUseCase(), NewNetworkGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError() (+6 more)

### Community 32 - "UI"
Cohesion: 0.06
Nodes (30): installHooks(), releaseBuild(), Run(), staged(), nextConfig(), showChangelog(), TestConfirmUpdateRoutesToActiveUI(), InstallMode (+22 more)

### Community 33 - "NetEventType"
Cohesion: 0.09
Nodes (12): snapNetMonitor, NetEvent, NetEventType, eventsOrAll(), newNetEventChecks(), NetEventTypes(), ParseNetEventType(), eventsetSummary() (+4 more)

### Community 34 - "picker"
Cohesion: 0.14
Nodes (7): keyEntry, picker, pickMode, errText(), isFindShortcut(), newKeyEntry(), newPickRow()

### Community 35 - "selection.go"
Cohesion: 0.15
Nodes (6): pos, cellToRune(), cursorPos(), Selection, offset(), changelogModel

### Community 37 - "backend"
Cohesion: 0.10
Nodes (12): checkSecret(), TestLiveSessionWatchAndPinger(), holdForward(), pingInterval(), File, backend, backendOp, backendState (+4 more)

### Community 38 - "controlServer"
Cohesion: 0.11
Nodes (21): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines() (+13 more)

### Community 39 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 40 - ".FilterExistingWhitelist"
Cohesion: 0.13
Nodes (17): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+9 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 43 - "monitorUnitTest"
Cohesion: 0.15
Nodes (4): sameInode(), dirTarget(), newMountTable(), monitorUnitTest

### Community 44 - "uninstall/uninstall.go"
Cohesion: 0.06
Nodes (45): runOfflineEdit(), anyDaemonRunning(), deletePostBackups(), listBackups(), restoreBackups(), askDecrypt(), confirmKeyDeletion(), decryptDirectories() (+37 more)

### Community 45 - "GuardEvent"
Cohesion: 0.08
Nodes (18): holdHeadless(), HeadlessLine(), runGuardHeadless(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestDispatchWithoutOwnAllowedEvents(), GuardEvent (+10 more)

### Community 46 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 48 - "User"
Cohesion: 0.04
Nodes (70): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), collectDiffAdditions(), configPaths(), pathCovered(), TestPathCovered() (+62 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (27): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+19 more)

### Community 50 - "backend.go"
Cohesion: 0.07
Nodes (29): confirmMasterKeyOverwrite(), runGenKey(), Requested(), TestRequested(), TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), TestSubcommandOutranksGUI(), Execute() (+21 more)

### Community 51 - "cursor.go"
Cohesion: 0.14
Nodes (13): buttonCursor, handCursors(), handCursorsOn(), newButtonCursor(), showConfirm(), showCustomConfirm(), showDialog(), showErrorDialog() (+5 more)

### Community 52 - "Client"
Cohesion: 0.08
Nodes (18): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+10 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.14
Nodes (37): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+29 more)

### Community 54 - "navItem"
Cohesion: 0.03
Nodes (16): appTheme, clickPanel, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, tabBar (+8 more)

### Community 55 - ".apply"
Cohesion: 0.11
Nodes (14): Section, dedupe(), Plan, sectionMaterialized(), TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal() (+6 more)

### Community 56 - "github.com/spf13/cobra.Command"
Cohesion: 0.05
Nodes (48): AddServeFlags(), CheckEBPF(), credentialFlagState(), isLoopbackHost(), MakeDisplayPaths(), ParseNetEventsFlag(), ParseServeFlags(), requestedBoolFlag() (+40 more)

### Community 57 - "mutedText"
Cohesion: 0.05
Nodes (46): actionGate, binaryRow, binSource, checkGrid, filesPage, handButton, homePage, iconTile (+38 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (47): engine, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+39 more)

### Community 59 - "os.File"
Cohesion: 0.04
Nodes (67): checkFreshDir(), openDirNoFollow(), renewDir(), setConfinedHomes(), VettedInode, Multicall, MulticallKind, ExeKey() (+59 more)

### Community 60 - "expandPlaceholders"
Cohesion: 0.14
Nodes (8): symlinkTargetGlob(), expandPlaceholders(), TestExpandWhitelist(), TestExtraWatchPathsForSkipsMissing(), LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 61 - "Recorder"
Cohesion: 0.14
Nodes (19): newBinary(), drain(), Recorder, Target, monitorTargets(), newRecorder(), opNames(), opsOf() (+11 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.14
Nodes (20): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+12 more)

### Community 63 - "Network"
Cohesion: 0.14
Nodes (7): NetworkGuardRepository, clonePolicy(), NetEngines, NetExtra, NetPolicy, Network, NetGuardConfig

### Community 64 - "CandidateDir"
Cohesion: 0.18
Nodes (13): Discover(), DiscoverSystem(), CandidateDir, TestDiscoverOnlyExisting(), TestDiscoverSystemOnlyAbsolute(), TestPathsFor(), catalogEntryMatch(), findCatalogRoot() (+5 more)

### Community 65 - "bpfstats/main.go"
Cohesion: 0.05
Nodes (61): CheckBPFLSM(), configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), runBPFCheck(), runLockdown(), runOneShotMode() (+53 more)

### Community 66 - "gui/session.go"
Cohesion: 0.07
Nodes (41): ancestorEnv(), cleanupSessionEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField() (+33 more)

### Community 67 - "newTestTarget"
Cohesion: 0.23
Nodes (16): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+8 more)

### Community 68 - "recordFlow"
Cohesion: 0.14
Nodes (15): encryptBlocker(), newRecordFlow(), confirmAndApply(), offerBackupDeletion(), pickEncryption(), runRecord(), runRecordView(), recordFlow (+7 more)

### Community 69 - "vouchedDevs"
Cohesion: 0.33
Nodes (6): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs()

### Community 70 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "proto.go"
Cohesion: 0.06
Nodes (30): Progress, User, fakeLink, multiAnswer, pageDef, stageItem, stageState, textAnswer (+22 more)

### Community 72 - "Guard"
Cohesion: 0.04
Nodes (32): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), deferredBinary, exeRow, rootHandle, tempMaskOp (+24 more)

### Community 73 - "FileTarget"
Cohesion: 0.15
Nodes (9): guardPill(), Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec(), overlaps() (+1 more)

### Community 74 - "pickfs.go"
Cohesion: 0.27
Nodes (13): listOpts, pickEntry, treeSearch, TestSearchTreeFollowsNoSymlink(), plural(), entryIsDir(), filterEntries(), listDir() (+5 more)

### Community 76 - "catalogRefresher"
Cohesion: 0.15
Nodes (7): drain(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), catalogRefresher, debouncer

### Community 77 - "update_test.go"
Cohesion: 0.08
Nodes (20): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel() (+12 more)

### Community 78 - "GlobReservations"
Cohesion: 0.17
Nodes (10): reserveChain(), GlobChild, compileGlobNames(), globText(), GlobReservations, ParseGlobName(), resolveChildren(), TestCompileGlobNames_RepeatedPatternOrsBits() (+2 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "integrationtests/networkguard_test.go"
Cohesion: 0.27
Nodes (5): IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 81 - "journal"
Cohesion: 0.14
Nodes (7): fakeEngines, fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 82 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 83 - "changelog.go"
Cohesion: 0.26
Nodes (15): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), latestTime(), newerThanInstalled(), newerThanStable() (+7 more)

### Community 84 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 85 - "highlight_test.go"
Cohesion: 0.12
Nodes (19): editorModel, fg(), Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo() (+11 more)

### Community 86 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.17
Nodes (12): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+4 more)

### Community 87 - "loadDaemonConfig"
Cohesion: 0.33
Nodes (6): loadDaemonConfig(), resolveConfigPath(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds()

### Community 88 - "liveSession"
Cohesion: 0.26
Nodes (8): parseErr(), TestLiveSessionConfigAndPut(), chooseResources(), eventsLabel(), forwardBinaries(), normalizedEvents(), runForward(), liveSession

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 92 - "Monitor"
Cohesion: 0.14
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.06
Nodes (29): fdStat, physicalParent(), statFD(), foreignPinnedLink(), NewMountVouch(), KernelDev(), LstatInode(), BtrfsDevices() (+21 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "globBuilder"
Cohesion: 0.35
Nodes (4): entryWriters(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder

### Community 100 - "codeedit.go"
Cohesion: 0.19
Nodes (26): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), moveTo() (+18 more)

### Community 101 - "downloadAndVerify"
Cohesion: 0.17
Nodes (20): archAssets(), assetsFor(), checkUpdatePreconditions(), downloadAndVerify(), fetchGUI(), FetchLatest(), fetchPair(), FetchVerified() (+12 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "updatePage"
Cohesion: 0.15
Nodes (5): latestRelease, updatePage, boolIndex(), isPreRelease(), Pulse

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.22
Nodes (13): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+5 more)

### Community 106 - "recordops.c"
Cohesion: 0.45
Nodes (13): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+5 more)

### Community 107 - "tempGrantSetup"
Cohesion: 0.23
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "hashCache"
Cohesion: 0.07
Nodes (28): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, BinaryStat, DurableStat, placedWalk (+20 more)

### Community 111 - "mountTable"
Cohesion: 0.27
Nodes (7): MountInfo, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), mountTable

### Community 112 - "NetGuardEvent"
Cohesion: 0.09
Nodes (19): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), runHeadless() (+11 more)

### Community 115 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 116 - "FileEvent"
Cohesion: 0.20
Nodes (3): snapMonitor, FileEvent, fakeSource

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.36
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "RefreshSection"
Cohesion: 0.07
Nodes (43): logRefreshChange(), TestLogRefreshChange(), askEncryption(), LibraryBlock, RefreshOptions, SectionScan, findSectionEnd(), findSectionStart() (+35 more)

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

### Community 129 - ".launch"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "install/install.go"
Cohesion: 0.05
Nodes (80): applyDiffAdditions(), restoreDaemonAfterDiffAbort(), runDiffCatalog(), applyLiveRefresh(), buildBinaryIfNeeded(), certifyInstall(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled() (+72 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "dialLiveSession"
Cohesion: 0.07
Nodes (31): ParseEventsFlag(), reserveStdout(), runBackend(), dialLiveSession(), LiveModeAvailable(), TestValidateFlagsEditConfig(), applyEdit(), checkAndApply() (+23 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "fakeGuardRepo"
Cohesion: 0.05
Nodes (24): confirmUserPlaced(), forwardCheck(), SupersededBinary, binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps() (+16 more)

### Community 143 - "time.Duration"
Cohesion: 0.17
Nodes (12): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+4 more)

### Community 144 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 145 - "Open"
Cohesion: 0.16
Nodes (20): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+12 more)

### Community 146 - "find_test.go"
Cohesion: 0.48
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 147 - "runView"
Cohesion: 0.07
Nodes (12): logGrid, logHook, logLine, runView, logRow(), logStyle(), newLogGrid(), captureLogs() (+4 more)

### Community 148 - "Binary"
Cohesion: 0.17
Nodes (16): lines(), pickBinaries(), pickLibs(), reviewRecording(), Classify(), finishBinary(), Binary, Recording (+8 more)

### Community 149 - "orphans_test.go"
Cohesion: 0.18
Nodes (19): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors(), orphanedPolicies() (+11 more)

### Community 150 - "EventType"
Cohesion: 0.18
Nodes (9): EventType, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), specFromForm(), specSummary(), EventTypes(), ParseEventType(), maskFor() (+1 more)

### Community 152 - "New"
Cohesion: 0.10
Nodes (42): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault() (+34 more)

### Community 153 - "shQuote"
Cohesion: 0.15
Nodes (4): IntegrationSuite, IntegrationSuite, shQuote(), IntegrationSuite

### Community 154 - "SanitizeText"
Cohesion: 0.06
Nodes (53): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), backingDeviceUnion() (+45 more)

### Community 155 - ".Layout"
Cohesion: 0.44
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 156 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 157 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.27
Nodes (19): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+11 more)

### Community 161 - "vaultfs.go"
Cohesion: 0.09
Nodes (26): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), Vault, IsBinaryContent(), IsOctalMode(), Open(), ParseMode() (+18 more)

### Community 162 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 163 - "controlManager"
Cohesion: 0.24
Nodes (4): newControlManager(), startControlManager(), configEditor, controlManager

### Community 164 - "RawTarget"
Cohesion: 0.27
Nodes (10): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+2 more)

### Community 165 - "sync.Mutex"
Cohesion: 0.05
Nodes (51): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+43 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (40): Options, pageID, reel, resizeWatch, runHost, shell, sizedDialog, snapEngines (+32 more)

### Community 167 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.18
Nodes (4): mountinfoUnitTest, unescapeMountinfo(), harnessSuite, editorHarness

### Community 168 - "Ledger"
Cohesion: 0.11
Nodes (5): applySchema(), Ledger, Pending, Source, fakeLedger

### Community 170 - "selection_test.go"
Cohesion: 0.23
Nodes (10): clipboardText(), setClipboard(), cellOf(), repeat(), selKey(), TestSelectionEditing(), TestSelectionLargeDelete(), TestSelectionMouseDrag() (+2 more)

### Community 171 - "dialog_test.go"
Cohesion: 0.45
Nodes (10): openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder(), TestPickerSubfolderSearchPicksFile() (+2 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 176 - "markReady"
Cohesion: 0.22
Nodes (10): markReady(), noteReloadedSettings(), notifySystemdReady(), inodeBatchFn, fakeBatches(), TestBatchInodeRows(), countInodeRows(), InodeOccupancy() (+2 more)

### Community 177 - "main_test.go"
Cohesion: 0.20
Nodes (5): infraContainerPaths(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

### Community 179 - "time.Time"
Cohesion: 0.24
Nodes (6): LockedError, sweepID, sweepMemo, Progress, accessor, Seen

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

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

### Community 194 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 195 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 196 - "edit-protected"
Cohesion: 0.22
Nodes (9): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Recording a new path: `--record`, Temporary access for other programs: `--forward` (+1 more)

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "securewrite_test.go"
Cohesion: 0.36
Nodes (7): descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeWithin()

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 202 - "DaemonEvent"
Cohesion: 0.09
Nodes (34): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemon(), runDaemonUI() (+26 more)

### Community 203 - "GenerateConf"
Cohesion: 0.10
Nodes (26): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles() (+18 more)

### Community 204 - "io.Writer"
Cohesion: 0.26
Nodes (10): printBinary(), printRecording(), candidates(), confirmable(), confirmAll(), describe(), hashNow(), run() (+2 more)

### Community 205 - "audit.go"
Cohesion: 0.21
Nodes (16): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource() (+8 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.15
Nodes (13): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Protected (+5 more)

### Community 209 - "WithBottomBar"
Cohesion: 0.31
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 210 - ".refreshOnce"
Cohesion: 0.43
Nodes (6): sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad()

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "resolveBits"
Cohesion: 0.43
Nodes (4): statFunc, resolveBits(), statKey(), UpdaterPlan

### Community 216 - "LibraryClosure"
Cohesion: 0.07
Nodes (30): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+22 more)

### Community 217 - "grantNewBinaries"
Cohesion: 0.50
Nodes (4): grantNewBinaries(), Line(), Lines(), TestLinesKeyOnLinkAndIncludePending()

### Community 219 - "confirmUpdate"
Cohesion: 0.67
Nodes (4): confirmUpdate(), isTerminal(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), ConfirmOnce()

### Community 220 - "update/update.go"
Cohesion: 0.18
Nodes (18): applyUpdate(), channelOf(), confirmUnit(), deploy(), deployGUI(), filterChannel(), findUpdate(), InstalledGUIVersion() (+10 more)

### Community 221 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 224 - "watchSet"
Cohesion: 0.15
Nodes (13): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+5 more)

### Community 225 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 227 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 230 - "unit.go"
Cohesion: 0.09
Nodes (29): installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), TestDaemonUnitWithMetadataOutput(), stagedBinary (+21 more)

### Community 235 - "Prereq"
Cohesion: 0.33
Nodes (3): confirmRunPrereq(), Prereq, fakeVault

### Community 237 - "serve.go"
Cohesion: 0.09
Nodes (30): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+22 more)

## Knowledge Gaps
- **31 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+26 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 586 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **36 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Editor keys` connect `edit-protected` to `uninstall/uninstall.go`, `backend`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _31 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.14682257244534824 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.0210239651416122 - nodes in this community are weakly interconnected._
- **Why does `Delete()` connect `uninstall/uninstall.go` to `go_pkg_path_filepath`, `edit-protected`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08957708957708958 - nodes in this community are weakly interconnected._