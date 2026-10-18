# Graph Report - app-listener  (2026-10-18)

## Corpus Check
- 486 files · ~1,711,346 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6453 nodes · 27852 edges · 240 communities (197 shown, 43 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2585 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `589047e6`
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
- ConfinedEntry
- policyEditor
- protectedPage
- github.com/charmbracelet/bubbletea.Cmd
- NetEventType
- editView
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- plan.go
- guard.bpf.c
- TrustGuard
- .showSized
- procstats.go
- rcUser
- Config
- TestMonitorUseCaseLifecycle
- runView
- NetworkMonitorRepository
- picker
- github.com/spf13/cobra.Command
- IntegrationSuite
- backend
- controlServer
- fileedit_model_test.go
- .ExpandLibDirWriters
- IntegrationSuite
- absPath
- monitorUnitTest
- RunWith
- GuardEvent
- IntegrationSuite
- stdio
- os.FileMode
- fadeScroll
- UI
- showConfirm
- Client
- daemonconfig.go
- navItem
- .apply
- runMonitor
- mutedText
- GuardInodeKey
- os.File
- CandidateDir
- Recorder
- exe_supersede.h
- sync.Mutex
- New
- bpfstats/main.go
- wm.go
- newTestTarget
- recordFlow
- KernelDev
- fcntl
- desktop.go
- Guard
- FileTarget
- .objs
- catalogRefresher
- update_test.go
- runNetworkGuard
- process_vm_readv.c
- appTheme
- journal
- planUpdaters
- githubRelease
- fileedit_symlink_test.go
- highlight_test.go
- render.go
- pinstate_test.go
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- configedit_test.go
- changelog_test.go
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- gui/session.go
- github.com/charmbracelet/bubbles/textarea.Model
- update/update.go
- Well-known bypass classes
- updatePage
- jit_provenance.c
- inspector_probe.c
- recordops.c
- tempGrantSetup
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- fakeDaemon
- mountTable
- NetGuardEvent
- SanitizeText
- IntegrationSuite
- IntegrationSuite
- classifyGroup
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
- RunCmd
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
- recoverFileInPlace
- DumpHook
- Binary
- SelectOrphans
- EventType
- fakeStore
- openFileVaultWithMasterKey
- hashCache
- buildOneGuard
- layout.go
- buildTrustedSet
- Monitor
- Installation
- backend_test.go
- graphify-refresh.sh
- Vault
- inodeChain
- runDaemon
- parseGuardEvent
- feed.go
- shell
- .populateWatchInodes
- Ledger
- taint_victim.c
- find.go
- inode_dev.h
- Limitations
- Troubleshooting
- IntegrationSuite
- IntegrationSuite
- Daemon configuration
- readKeyFrom
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
- BinaryStat
- lockAndDeprovision
- edit-protected
- monitor: see what touches a path
- kernel-prof.sh
- go_pkg_strings
- pathCache
- installPage
- DaemonEvent
- GenerateConf
- io.Writer
- audit.go
- Desktop GUI
- IntegrationSuite
- findSectionStart
- downloadFile
- .eventTypeAllowed
- BinaryEntry
- WithBottomBar
- bleedLayout
- IntegrationSuite
- LibraryClosure
- NewEventFanout
- fakeVault
- ParseEventsFlag
- runBackend
- run
- IntegrationSuite
- logGrid
- watchSet
- bpfLsmListed
- ValidateTargets
- go_pkg_os
- fakeLink
- headless.go
- PinPrefix
- .runEditorHarness
- TestCleanRel
- NetworkGuardRepository
- stripImmutableFlags
- Prereq
- VerifyInstalledResourcesLocked
- Serve
- applyLiveRefresh
- partitionEncryptionRoots

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

## Communities (240 total, 43 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.17
Nodes (40): AddServeFlags(), TestParseServeFlagsRequiresCredentials(), init(), init(), init(), init(), init(), tempJournalRow (+32 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (132): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), awaitStartupOrSignal(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestAwaitStartupOrSignalAbortsAndLocksDown(), TestAwaitStartupOrSignalAbortWithFailedStartup() (+124 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "Vault"
Cohesion: 0.14
Nodes (17): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), verifyEncryptionState(), isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy() (+9 more)

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
Nodes (28): addErrKind, bpfGuardEvent, ReplacementCheck, TestWithChmodDropWriteNeedsReadOnly(), canonicalPaths(), classifyAddErr(), guardModeKey(), modeString() (+20 more)

### Community 9 - "mkdirs"
Cohesion: 0.10
Nodes (43): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), newWatchPattern(), newFakeInotify() (+35 more)

### Community 10 - "ConfinedEntry"
Cohesion: 0.12
Nodes (28): binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats() (+20 more)

### Community 11 - "policyEditor"
Cohesion: 0.14
Nodes (11): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), newNetGuardTab() (+3 more)

### Community 12 - "protectedPage"
Cohesion: 0.06
Nodes (17): backendLink, checkGrid, configView, forwardForm, forwardView, protectedPage, recordForm, recordReview (+9 more)

### Community 13 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.06
Nodes (23): recordTickCmd(), recordModel, diffModel, editorModel, newDiffModel(), newEditorModel(), formatType(), listenForEvents() (+15 more)

### Community 14 - "NetEventType"
Cohesion: 0.07
Nodes (18): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, eventsOrAll(), newNetEventChecks(), Cstr(), FormatAddr() (+10 more)

### Community 15 - "editView"
Cohesion: 0.12
Nodes (7): codeEntry, editView, Entry, fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (7): eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (11): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.20
Nodes (44): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+36 more)

### Community 22 - "Resource"
Cohesion: 0.08
Nodes (9): Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots(), vaultOpForGuard() (+1 more)

### Community 23 - "auth_test.go"
Cohesion: 0.12
Nodes (40): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+32 more)

### Community 24 - "plan.go"
Cohesion: 0.38
Nodes (4): recordTick, isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (73): exe_refused(), inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade() (+65 more)

### Community 26 - "TrustGuard"
Cohesion: 0.05
Nodes (25): GlobChild, statFunc, trustHook, trustRows, setSupersedeTrust(), cStr(), TestPlanTrustedDirs_DuplicateWriterAddsNoEmptyLoader(), TestPlanTrustedDirs_LibDirsOnlyForTheirWriters() (+17 more)

### Community 27 - ".showSized"
Cohesion: 0.18
Nodes (10): installUI, multiAnswer, scrollable(), selectable(), fadeScrollOf(), newFadeScroll(), ask(), orDefault() (+2 more)

### Community 28 - "procstats.go"
Cohesion: 0.08
Nodes (24): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+16 more)

### Community 29 - "rcUser"
Cohesion: 0.13
Nodes (34): setupBunTmpdir(), rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), EnsureBunTmpDir(), ensureRCBlock() (+26 more)

### Community 30 - "Config"
Cohesion: 0.08
Nodes (52): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), appendSectionsAndEdit(), applyDiffAdditions(), collectDiffAdditions(), configPaths() (+44 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.27
Nodes (10): NewGuardUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestNetworkGuardUseCaseLifecycle(), TestNetworkGuardUseCaseStartError() (+2 more)

### Community 32 - "runView"
Cohesion: 0.14
Nodes (3): runHost, runView, captureLogs()

### Community 33 - "NetworkMonitorRepository"
Cohesion: 0.20
Nodes (3): NetworkMonitorRepository, NetworkMonitorUseCase, NewNetworkMonitorUseCase()

### Community 34 - "picker"
Cohesion: 0.08
Nodes (30): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+22 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.18
Nodes (17): credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestSetupDumpLogRefusesSymlink() (+9 more)

### Community 37 - "backend"
Cohesion: 0.10
Nodes (12): checkSecret(), TestCheckSecretRefusesLineBreaks(), File, backend, backendOp, backendState, BinaryCheck, Message (+4 more)

### Community 38 - "controlServer"
Cohesion: 0.07
Nodes (29): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines() (+21 more)

### Community 39 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 40 - ".ExpandLibDirWriters"
Cohesion: 0.14
Nodes (14): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestSteamLibDirsCoverProtonCodeOnly(), TestSteamLibDirsExcludeScoutRuntime(), TestSteamLibDirWritersCoverItsOwnBinaries(), ConfSafePath() (+6 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (7): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath(), TestMain(), helperChild(), TestMain()

### Community 44 - "RunWith"
Cohesion: 0.16
Nodes (21): RunWith(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), deletePostBackups(), listBackups(), restoreBackups(), confirmKeyDeletion(), offerBackupCleanup() (+13 more)

### Community 45 - "GuardEvent"
Cohesion: 0.07
Nodes (21): holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI(), snapEngines, snapGuard (+13 more)

### Community 46 - "IntegrationSuite"
Cohesion: 0.16
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 48 - "os.FileMode"
Cohesion: 0.07
Nodes (39): checkAndApply(), runRefreshCatalogOnly(), softenAutomatedRefreshErr(), refreshCatalogConfig(), TestSoftenAutomatedRefreshErr(), dupFD(), installFile(), installFileAs() (+31 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (26): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+18 more)

### Community 50 - "UI"
Cohesion: 0.07
Nodes (25): nextConfig(), showChangelog(), confirmUpdate(), isTerminal(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), TestRunWithRefusesConcurrentRun() (+17 more)

### Community 51 - "showConfirm"
Cohesion: 0.12
Nodes (13): buttonCursor, handCursors(), handCursorsOn(), newButtonCursor(), showConfirm(), showCustomConfirm(), showDialog(), showErrorDialog() (+5 more)

### Community 52 - "Client"
Cohesion: 0.09
Nodes (15): Client, childEnv(), launch(), Start(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines(), Hello (+7 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.12
Nodes (43): grantNewBinaries(), watchGroup, Line(), Lines(), TestLinesKeyOnLinkAndIncludePending(), addResource(), applyAllowLib(), applyConfigLine() (+35 more)

### Community 54 - "navItem"
Cohesion: 0.03
Nodes (22): binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pill (+14 more)

### Community 55 - ".apply"
Cohesion: 0.12
Nodes (14): Plan, TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan(), TestPlanAdmitsOnlyWhatWasRecorded(), testRecording() (+6 more)

### Community 56 - "runMonitor"
Cohesion: 0.11
Nodes (18): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+10 more)

### Community 57 - "mutedText"
Cohesion: 0.11
Nodes (29): homePage, panel, statusTile, newBinaryRow(), newFadeVScroll(), newTargetRow(), newCounterStrip(), newStatusTile() (+21 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.04
Nodes (45): refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestSetInspectorsBounded(), appendInodeRow(), batchInodeRows(), inodeRowsOf() (+37 more)

### Community 59 - "os.File"
Cohesion: 0.04
Nodes (75): checkFreshDir(), openDirNoFollow(), renewDir(), setConfinedHomes(), VettedInode, Multicall, MulticallKind, placedWalk (+67 more)

### Community 60 - "CandidateDir"
Cohesion: 0.06
Nodes (35): systemPatterns(), systemRule(), entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder (+27 more)

### Community 61 - "Recorder"
Cohesion: 0.16
Nodes (15): Target, drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), sortedKeys() (+7 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.19
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+7 more)

### Community 63 - "sync.Mutex"
Cohesion: 0.09
Nodes (14): netGuardPill(), newActivity(), guardRow(), TestRingKeepsNewest(), clonePolicy(), NetExtra, NetPolicy, Network (+6 more)

### Community 64 - "New"
Cohesion: 0.17
Nodes (19): harnessSuite, looksLikeFileVaultRecord(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceWrongKeyFails(), TestUnlockLockFileInPlacePreservesInode(), TestUnlockLockFileInPlaceRoundTrip() (+11 more)

### Community 65 - "bpfstats/main.go"
Cohesion: 0.13
Nodes (23): CheckBPFLSM(), runBPFCheck(), guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity() (+15 more)

### Community 66 - "wm.go"
Cohesion: 0.11
Nodes (24): reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket(), poll() (+16 more)

### Community 67 - "newTestTarget"
Cohesion: 0.21
Nodes (17): newExeResolver(), FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps() (+9 more)

### Community 68 - "recordFlow"
Cohesion: 0.16
Nodes (15): encryptBlocker(), newRecordFlow(), confirmAndApply(), offerBackupDeletion(), pickEncryption(), runRecord(), runRecordView(), recordFlow (+7 more)

### Community 69 - "KernelDev"
Cohesion: 0.09
Nodes (26): fdStat, physicalParent(), statFD(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed() (+18 more)

### Community 70 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "desktop.go"
Cohesion: 0.12
Nodes (11): releaseBuild(), Run(), RunWith(), TestRunWithRefusesConcurrentRun(), stageItem, stageState, textAnswer, encodeTable() (+3 more)

### Community 72 - "Guard"
Cohesion: 0.04
Nodes (24): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), engine, rootHandle, SectionScan, openRoot() (+16 more)

### Community 73 - "FileTarget"
Cohesion: 0.09
Nodes (11): filesPage, guardPill(), Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec() (+3 more)

### Community 74 - ".objs"
Cohesion: 0.11
Nodes (6): exeRow, tempGrant, tempMaskOp, tempRow, eventMask(), liveTrust()

### Community 76 - "catalogRefresher"
Cohesion: 0.15
Nodes (8): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), catalogRefresher, debouncer

### Community 77 - "update_test.go"
Cohesion: 0.10
Nodes (17): fetchReleases(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel(), testKeyPair(), TestNewerThanInstalled() (+9 more)

### Community 78 - "runNetworkGuard"
Cohesion: 0.10
Nodes (24): CheckEBPF(), ParseNetEventsFlag(), newPinGeneration(), prepareDaemonStart(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode() (+16 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 81 - "journal"
Cohesion: 0.14
Nodes (7): fakeEngines, fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 82 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 83 - "githubRelease"
Cohesion: 0.25
Nodes (17): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+9 more)

### Community 84 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 85 - "highlight_test.go"
Cohesion: 0.15
Nodes (17): fg(), Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens() (+9 more)

### Community 86 - "render.go"
Cohesion: 0.13
Nodes (13): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+5 more)

### Community 87 - "pinstate_test.go"
Cohesion: 0.12
Nodes (28): configureDaemonLogging(), loadDaemonConfig(), resolveConfigPath(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState() (+20 more)

### Community 88 - "liveSession"
Cohesion: 0.14
Nodes (15): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), displayEntries(), eventsLabel(), forwardBinaries() (+7 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 92 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.12
Nodes (8): eventTypeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), hook, NetGuard, procRec

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "gui/session.go"
Cohesion: 0.09
Nodes (29): invokingIDs(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), ancestorEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv() (+21 more)

### Community 100 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.11
Nodes (31): class, pos, classOf(), column(), dedent(), Edit(), leadingIndent(), lines() (+23 more)

### Community 101 - "update/update.go"
Cohesion: 0.13
Nodes (31): installHooks(), staged(), applyUpdate(), archAssets(), assetsFor(), channelOf(), checkUpdatePreconditions(), confirmUnit() (+23 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "updatePage"
Cohesion: 0.26
Nodes (4): latestRelease, updatePage, boolIndex(), isPreRelease()

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
Cohesion: 0.17
Nodes (9): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), btrfsTypes(), TestBtrfsLayoutFromNestedMembers() (+1 more)

### Community 111 - "mountTable"
Cohesion: 0.13
Nodes (10): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+2 more)

### Community 112 - "NetGuardEvent"
Cohesion: 0.18
Nodes (6): snapNetGuard, NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 113 - "SanitizeText"
Cohesion: 0.11
Nodes (15): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), startTrustGuard() (+7 more)

### Community 116 - "classifyGroup"
Cohesion: 0.11
Nodes (23): appletLine(), Classify(), classifyGroup(), classifyLibs(), finishBinary(), innermostResource(), libCovered(), newBinary() (+15 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.36
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "RefreshSection"
Cohesion: 0.11
Nodes (25): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+17 more)

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
Cohesion: 0.31
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "RunCmd"
Cohesion: 0.07
Nodes (38): installGUI(), buildBinaryIfNeeded(), certifyInstall(), ensureInstalledBinary(), installBinaryOnly(), markComplete(), preflightDeployedBinary(), prepareInstallation() (+30 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "dialLiveSession"
Cohesion: 0.12
Nodes (17): TestValidatePassword(), ValidatePassword(), dialLiveSession(), applyEdit(), editAgain(), fetchConfig(), putConfig(), runEditConfig() (+9 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "fakeGuardRepo"
Cohesion: 0.06
Nodes (19): SupersededBinary, TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows(), applyTempGrants() (+11 more)

### Community 143 - "time.Duration"
Cohesion: 0.15
Nodes (13): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+5 more)

### Community 144 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 145 - "Open"
Cohesion: 0.13
Nodes (21): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+13 more)

### Community 146 - "recoverFileInPlace"
Cohesion: 0.31
Nodes (13): classifyRegularFileTarget(), clearRecovery(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), recoverFileInPlace(), rejectSidecarSymlink(), stageRecovery(), inode() (+5 more)

### Community 147 - "DumpHook"
Cohesion: 0.16
Nodes (8): logHook, logLine, logRow(), logStyle(), formatEntry(), sortedKeys(), sortStrings(), DumpHook

### Community 148 - "Binary"
Cohesion: 0.20
Nodes (13): formatElapsed(), lines(), pickBinaries(), pickLibs(), printBinary(), printRecording(), reviewRecording(), Binary (+5 more)

### Community 149 - "SelectOrphans"
Cohesion: 0.24
Nodes (16): appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors(), SelectOrphans() (+8 more)

### Community 150 - "EventType"
Cohesion: 0.08
Nodes (20): snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), planTempBlock() (+12 more)

### Community 152 - "openFileVaultWithMasterKey"
Cohesion: 0.18
Nodes (14): deriveFileVaultSubkey(), newFileVaultAEAD(), openFileVault(), openFileVaultWithMasterKey(), sealFileVault(), sealFileVaultWithMasterKey(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated() (+6 more)

### Community 153 - "hashCache"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 154 - "buildOneGuard"
Cohesion: 0.09
Nodes (38): backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), reloadOnce(), reloadSlotsNeeded(), startGuardedDaemon() (+30 more)

### Community 155 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 156 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 157 - "Monitor"
Cohesion: 0.18
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.29
Nodes (18): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+10 more)

### Community 161 - "Vault"
Cohesion: 0.09
Nodes (26): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), Vault, IsBinaryContent(), IsOctalMode(), Open(), ParseMode() (+18 more)

### Community 162 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 163 - "runDaemon"
Cohesion: 0.10
Nodes (20): newControlManager(), startControlManager(), catchLifecycleSignals(), makeReloadHandler(), markReady(), noteReloadedSettings(), notifySystemdReady(), runDaemon() (+12 more)

### Community 164 - "parseGuardEvent"
Cohesion: 0.29
Nodes (8): decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels(), suspectLabel()

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (46): actionBar, binaryList, cell, column, counterStrip, dataTable, eventFeed, feedAction (+38 more)

### Community 166 - "shell"
Cohesion: 0.06
Nodes (41): pageDef, pageID, reel, resizeWatch, shell, sizedDialog, snapNet, newLogo() (+33 more)

### Community 168 - "Ledger"
Cohesion: 0.11
Nodes (5): applySchema(), Ledger, Pending, Source, fakeLedger

### Community 170 - "find.go"
Cohesion: 0.10
Nodes (28): findButton, span, findAll(), lowerRunes(), NewFinder(), newFindInput(), overlay(), runesEqual() (+20 more)

### Community 171 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "readKeyFrom"
Cohesion: 0.18
Nodes (12): confirmMasterKeyOverwrite(), runGenKey(), runOneShotMode(), ensureMasterKey(), checkKeyLen(), GenerateMasterKey(), MasterKeyExists(), readKeyFrom() (+4 more)

### Community 179 - "time.Time"
Cohesion: 0.13
Nodes (14): LockedError, newWatchSet(), fakeInotify, inodeBatchFn, sweepID, sweepMemo, fakeBatches(), TestBatchInodeRows() (+6 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "NewDumpHook"
Cohesion: 0.21
Nodes (12): applyVerbosity(), effectiveVerbose(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced() (+4 more)

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

### Community 194 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 195 - "lockAndDeprovision"
Cohesion: 0.12
Nodes (12): deprovisionKind, migrationTarget, writeTempWithMetadata(), classifyDeprovision(), classifyDeprovisionErr(), classifyMigrationTarget(), isDeprovisionBusy(), isDeprovisionMissing() (+4 more)

### Community 196 - "edit-protected"
Cohesion: 0.22
Nodes (9): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Recording a new path: `--record`, Temporary access for other programs: `--forward` (+1 more)

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "go_pkg_strings"
Cohesion: 0.05
Nodes (22): Requested(), TestRequested(), atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget() (+14 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "installPage"
Cohesion: 0.08
Nodes (22): actionGate, binSource, installAction, installedBinary, InstallMode, installPage, InstallRequest, Options (+14 more)

### Community 202 - "DaemonEvent"
Cohesion: 0.09
Nodes (36): UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI(), runTUI() (+28 more)

### Community 203 - "GenerateConf"
Cohesion: 0.06
Nodes (51): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic() (+43 more)

### Community 204 - "io.Writer"
Cohesion: 0.31
Nodes (8): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), fyneLogFilter, candidate

### Community 205 - "audit.go"
Cohesion: 0.08
Nodes (31): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource() (+23 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.15
Nodes (13): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Protected (+5 more)

### Community 209 - "findSectionStart"
Cohesion: 0.25
Nodes (8): findSectionStart(), isLibDirectiveText(), ParseSectionHeaderPath(), unquotePath(), IsLibraryDirective(), parseLibDirLine(), ParseSectionWhitelist(), TestParseSectionWhitelist()

### Community 210 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "BinaryEntry"
Cohesion: 0.06
Nodes (28): deferredBinary, keptEntries, BinaryEntry, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), checkBtrfsKeys(), ensureBtrfsLayout(), BinariesSummary(), canonicalBinaryPath() (+20 more)

### Community 213 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 216 - "LibraryClosure"
Cohesion: 0.07
Nodes (30): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+22 more)

### Community 217 - "NewEventFanout"
Cohesion: 0.24
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 218 - "fakeVault"
Cohesion: 0.42
Nodes (3): TestDaemonUseCaseReloadRefusedDuringShutdown(), TestDaemonUseCaseStopBlocksUntilAllLocked(), fakeVault

### Community 219 - "ParseEventsFlag"
Cohesion: 0.22
Nodes (7): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule(), TestUnchangedPerShared()

### Community 220 - "runBackend"
Cohesion: 0.25
Nodes (7): reserveStdout(), runBackend(), LiveModeAvailable(), anyDaemonRunning(), changePasswordLive(), DaemonRunning(), systemctlActive()

### Community 221 - "run"
Cohesion: 0.36
Nodes (8): InstalledGUIVersion(), InstalledVersion(), readInstalledVersion(), run(), runUpdate(), syncGUI(), TestReadInstalledVersion(), GUIInstalled()

### Community 224 - "watchSet"
Cohesion: 0.18
Nodes (11): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchDir, watchPattern (+3 more)

### Community 225 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 226 - "ValidateTargets"
Cohesion: 0.32
Nodes (7): guardedOverlap(), overlaps(), refusedPath(), validateTarget(), ValidateTargets(), mustParse(), TestValidateTargets()

### Community 227 - "go_pkg_os"
Cohesion: 0.09
Nodes (22): needsBinary(), TestNeedsBinary(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), fileSHA256(), sameContent(), TestCheckRunningBinaryMatchesInstalled(), removeKeyAndEmptyDir() (+14 more)

### Community 229 - "headless.go"
Cohesion: 0.60
Nodes (3): Check(), CheckGUI(), inspect()

### Community 230 - "PinPrefix"
Cohesion: 0.40
Nodes (5): PinPrefix(), sanitizeGen(), SharedPinPrefix(), TestPinPrefixDistinctAndStable(), TestPinPrefixShape()

### Community 233 - "NetworkGuardRepository"
Cohesion: 0.18
Nodes (7): WithSuspectsFrom(), NetworkGuardRepository, RealNetEngines(), NetEngines, netguardHarness, Option, NetGuardConfig

### Community 234 - "stripImmutableFlags"
Cohesion: 0.50
Nodes (4): stripEntryFlags(), stripImmutableFlags(), TestStripImmutableFlagsDoesNotHangOnFIFO(), TestStripImmutableFlagsSkipsDanglingSymlink()

### Community 235 - "Prereq"
Cohesion: 0.33
Nodes (3): confirmRunPrereq(), Prereq, fakeVault

### Community 236 - "VerifyInstalledResourcesLocked"
Cohesion: 0.67
Nodes (3): TestVerifyResourcesLocked(), VerifyInstalledResourcesLocked(), VerifyResourcesLocked()

### Community 237 - "Serve"
Cohesion: 0.08
Nodes (29): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+21 more)

## Knowledge Gaps
- **31 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+26 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 586 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **43 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Editor keys` connect `edit-protected` to `liveSession`, `RunWith`?**
  _High betweenness centrality (0.044) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _31 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.0197509660798626 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.044) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08958130477117819 - nodes in this community are weakly interconnected._
- **Why does `Delete()` connect `RunWith` to `plan.go`, `edit-protected`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **Should `Vault` be split into smaller, more focused modules?**
  _Cohesion score 0.13709677419354838 - nodes in this community are weakly interconnected._